package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/comrade-coop/apocryph/backend/swarm"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/replication"
)

var ReplicationServiceAccountPrefix = "Replicator-"
var ReplicationCustomTokenIamArn = "arn:minio:iam:::role/idmp-ethauth"

type ReplicationManager struct {
	minio       *minio.Client
	minioAdmin  *madmin.AdminClient
	swarm       *swarm.Swarm
	tokenSigner *TokenSigner
}

func NewReplicationManager(minioAddress string, minioCreds *credentials.Credentials, swarm *swarm.Swarm, tokenSigner *TokenSigner) (*ReplicationManager, error) {
	minioClient, err := minio.New(minioAddress, &minio.Options{
		Creds: minioCreds,
	})
	if err != nil {
		return nil, err
	}
	minioAdmin, err := madmin.NewWithOptions(minioAddress, &madmin.Options{
		Creds: minioCreds,
	})
	if err != nil {
		return nil, err
	}

	return &ReplicationManager{
		minio:       minioClient,
		minioAdmin:  minioAdmin,
		swarm:       swarm,
		tokenSigner: tokenSigner,
	}, nil
}

func (r *ReplicationManager) Run(ctx context.Context) error {
	go func() {
		for {
			err := r.reconcilationLoop(ctx)
			if err != nil {
				log.Printf("Error while reconciling replications: %v", err)
			}

			time.Sleep(1 * time.Minute) // TODO: Maybe change to 5 minutes?
		}
	}()
	return nil
}

func (r *ReplicationManager) reconcilationLoop(ctx context.Context) error {
	errs := []error{}
	buckets, err := r.minio.ListBuckets(ctx)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, bucket := range buckets {
			err = r.UpdateBucket(ctx, bucket.Name)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}

	err = r.UpdateCapacity(ctx)
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type ReplicaStatus int

const (
	Unknown ReplicaStatus = 0
	Alive   ReplicaStatus = 1
	Failing ReplicaStatus = 2
)

// TODO: call UpdateBucket whenever the swarm detects that the peer is down
// TODO: call UpdateBucket whenever we get a new bucket
func (r *ReplicationManager) UpdateBucket(ctx context.Context, bucketId string) error {
	// Get object
	replicationConfig, err := r.minio.GetBucketReplication(ctx, bucketId)
	if err != nil {
		response := minio.ToErrorResponse(err)
		if response.Code != "" {
			println("!!!", response.Code)
			replicationConfig = replication.Config{}
		} else {
			return err
		}
	}
	targetTotalReplicas := 2 // TODO: Make this configurable

	// Get current status
	metrics, err := r.minio.GetBucketReplicationMetrics(ctx, bucketId)
	if err != nil {
		return err
	}

	replicaStatus := map[string]ReplicaStatus{}
	statusCounts := map[ReplicaStatus]int{}

	// Count ourselves
	replicaStatus[r.swarm.OwnName] = Alive
	statusCounts[Alive]++
	syncedReplicas := 1

	for id, stats := range metrics.Stats {
		if stats.Failed.LastHour.Count > 10 {
			replicaStatus[id] = Failing
			statusCounts[Failing]++
		} else {
			replicaStatus[id] = Alive
			statusCounts[Alive]++
		}
		if stats.PendingCount < 10 {
			syncedReplicas += 1
		}
	}

	expectedReplicas, err := r.swarm.FindBucketReplicas(bucketId)
	if err != nil {
		return err
	}

	// Reconcile
	// TODO: Make sure this does what it's supposed to
	for _, node := range expectedReplicas {
		id := node
		if replicaStatus[id] == Unknown {
			rule, err := r.getReplicationRuleForNode(ctx, node, bucketId)
			if err != nil {
				//return err
				continue
			}
			replicaStatus[id] = Alive
			statusCounts[Alive]++
			replicationConfig.AddRule(rule)
		}
	}

	if statusCounts[Alive] < targetTotalReplicas {
		for statusCounts[Alive] < targetTotalReplicas {
			node, err := r.swarm.FindFreeNode()
			if err != nil {
				return err
			}
			if replicaStatus[node] != Unknown {
				break // TODO: avoids infinite loop
			}
			rule, err := r.getReplicationRuleForNode(ctx, node, bucketId)
			if err != nil {
				//return err
				continue
			}
			replicaStatus[node] = Alive
			statusCounts[Alive]++
			replicationConfig.AddRule(rule)
		}
	} else if syncedReplicas > targetTotalReplicas/2+1 { // Cleanup - we have enough alive and synced to start prunning failed
		// TODO: Double-check condition correctness
		for failedId, status := range replicaStatus {
			if status == Failing {
				replicationConfig.RemoveRule(replication.Options{
					ID: failedId,
				})
			}
		}
	}

	r.minio.SetBucketReplication(ctx, bucketId, replicationConfig)

	return nil
}
func (r *ReplicationManager) getReplicationRuleForNode(ctx context.Context, node string, bucketId string) (opt replication.Options, err error) {
	hostname := node

	token, err := r.tokenSigner.GetReplicationToken(bucketId)
	if err != nil {
		return
	}
	cred, err := credentials.NewCustomTokenCredentials(hostname, token, ReplicationCustomTokenIamArn)
	if err != nil {
		return
	}
	adminClient, err := madmin.NewWithOptions(hostname, &madmin.Options{
		Secure: false, // TODO: true
		Creds:  cred,
	})
	if err != nil {
		return
	}
	result, err := adminClient.AddServiceAccount(ctx, madmin.AddServiceAccountReq{
		Name:        ReplicationServiceAccountPrefix + r.swarm.OwnName,
		Description: "Access key used for replicating this bucket to other instances of Apocryph S3, keeping your data safe from any single node failures!",
	})
	if err != nil {
		return
	}
	accessKey := result.AccessKey
	secretKey := result.SecretKey

	// TODO: Create the remote bucket?

	opt = replication.Options{
		ID:                      node,
		ExistingObjectReplicate: "enable",
		ReplicateDeletes:        "enable",
		ReplicateDeleteMarkers:  "enable",
		DestBucket:              fmt.Sprintf("http://%s:%s@%s/%s", accessKey, secretKey, hostname, bucketId), // TODO: https
	}
	return
}

func (r *ReplicationManager) UpdateCapacity(ctx context.Context) error {
	info, err := r.minioAdmin.StorageInfo(ctx)
	if err != nil {
		return err
	}

	available := uint64(0)
	for _, disk := range info.Disks {
		available += disk.AvailableSpace
	}

	r.swarm.UpdateCapacity(available / 1024 / 1024) // Divide, so that we don't end up with a new number every single small change

	return nil
}
