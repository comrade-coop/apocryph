package main

import (
	"context"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/comrade-coop/apocryph/backend/swarm"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/replication"
	"golang.org/x/crypto/sha3"
)

var ReplicationServiceAccountPrefix = "Replicator-"
var ReplicationCustomTokenIamArn = "arn:minio:iam:::role/idmp-swieauth"

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

			time.Sleep(5 * time.Minute)
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
	Unknown    ReplicaStatus = 0
	Configured ReplicaStatus = 1
	Alive      ReplicaStatus = 2
	Failing    ReplicaStatus = 3
)

// TODO: call UpdateBucket whenever the swarm detects that the peer is down
// TODO: call UpdateBucket whenever we get a new bucket
func (r *ReplicationManager) UpdateBucket(ctx context.Context, bucketId string) error {
	_ = r.swarm.UpdateBucket(bucketId, swarm.Syncing) // TODO: Update syncing/ready correctly
	_ = r.minio.EnableVersioning(ctx, bucketId)       // TODO: Move to frontend/etc.?

	// Get object
	replicationConfig, err := r.minio.GetBucketReplication(ctx, bucketId)
	if err != nil {
		response := minio.ToErrorResponse(err)
		if response.Code != "" {
			println("!!!", response.Code)
			replicationConfig = replication.Config{}
		} else {
			return fmt.Errorf("GetBucketReplication: %w", err)
		}
	}
	targetTotalReplicas := 2 // TODO: Make this configurable
	highestPriority := 0

	var metrics replication.Metrics
	if len(replicationConfig.Rules) > 0 {
		// Get current status
		metrics, err = r.minio.GetBucketReplicationMetrics(ctx, bucketId)
		if err != nil {
			response := minio.ToErrorResponse(err)
			if response.Code != "" {
				println("!!!!", response.Code)
				replicationConfig = replication.Config{}
			} else {
				return fmt.Errorf("GetBucketReplicationMetrics: %w", err)
			}
		}
		a,_:=json.Marshal(metrics)
		println(string(a))
		for _, rule := range replicationConfig.Rules {
			highestPriority = max(highestPriority, rule.Priority)
		}
	}

	replicaStatus := map[string]ReplicaStatus{}
	statusCounts := map[ReplicaStatus]int{}

	// Count ourselves
	replicaStatus[r.swarm.OwnName] = Alive
	statusCounts[Alive]++
	syncedReplicas := 1

	for _, rule := range replicationConfig.Rules {
		replicaStatus[rule.ID] = Configured
		statusCounts[Alive]++ // HACK: Don't assume everyone is down while metrics below are not yet fixed
	}
	// TODO: metrics.Stats is nil for some reason
	/*
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
	}*/

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
				// continue
				return fmt.Errorf("getReplicationRuleForNode: %w", err)
			}
			highestPriority++
			rule.Priority = strconv.FormatInt(int64(highestPriority), 10)

			replicaStatus[id] = Alive
			statusCounts[Alive]++

			err = replicationConfig.AddRule(rule)
			if err != nil {
				return fmt.Errorf("AddRule (expected): %w", err)
			}
		}
	}

	if statusCounts[Alive] < targetTotalReplicas {
		for statusCounts[Alive] < targetTotalReplicas {
			node, err := r.swarm.FindFreeNode()
			if err != nil {
				return fmt.Errorf("FindFreeNode: %w", err)
			}
			if replicaStatus[node] != Unknown {
				break // TODO: avoids infinite loop
			}
			rule, err := r.getReplicationRuleForNode(ctx, node, bucketId)
			if err != nil {
				// continue
				return fmt.Errorf("getReplicationRuleForNode: %w", err)
			}
			highestPriority++
			rule.Priority = strconv.FormatInt(int64(highestPriority), 10)

			replicaStatus[node] = Alive
			statusCounts[Alive]++

			err = replicationConfig.AddRule(rule)
			if err != nil {
				return fmt.Errorf("AddRule (free): %w", err)
			}
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
	fmt.Println("Set replication:", replicationConfig)

	err = r.minio.SetBucketReplication(ctx, bucketId, replicationConfig)
	if err != nil {
		return fmt.Errorf("SetBucketReplication: %w", err)
	}
	fmt.Println("Set replication success wohoo!")

	return nil
}
func (r *ReplicationManager) getReplicationRuleForNode(ctx context.Context, node string, bucketId string) (opt replication.Options, err error) {
	hostname := node

	token, err := r.tokenSigner.GetReplicationToken(bucketId)
	if err != nil {
		err = fmt.Errorf("GetReplicationToken: %w", err)
		return
	}
	stsEndpoint := "http://" + hostname // TODO: https
	cred, err := credentials.NewCustomTokenCredentials(stsEndpoint, token, ReplicationCustomTokenIamArn)
	if err != nil {
		return
	}
	
	// Creating a bucket is handled by the identity plugin.

	// TODO: Find a way to sync bucket policies (below won't work due to only being triggered for/on new nodes)
	// ownPolicy, err := r.minio.GetBucketPolicy(ctx, bucketId)
	// if err != nil {
	// 	err = fmt.Errorf("GetBucketPolicy: %w", err)
	// 	return
	// }
	// err = minioClient.SetBucketPolicy(ctx, bucketId, ownPolicy)
	// if err != nil {
	// 	err = fmt.Errorf("SetBucketPolicy: %w", err)
	// 	return
	// }

	arn, err := r.getRemoteTargetArn(ctx, cred, hostname, bucketId)
	if err != nil {
		err = fmt.Errorf("getRemoteTargetArn: %w", err)
		return
	}

	opt = replication.Options{
		ID:                      node,
		RuleStatus:              "enable",
		ExistingObjectReplicate: "enable",
		ReplicateDeletes:        "enable",
		ReplicateDeleteMarkers:  "enable",
		ReplicaSync:             "enable",
		DestBucket:              arn,
	}
	fmt.Println("Add replication:", opt)
	return
}
func (r *ReplicationManager) getRemoteTargetArn(ctx context.Context, cred *credentials.Credentials, hostname string, bucketId string) (arn string, err error) {
	targets, err := r.minioAdmin.ListRemoteTargets(ctx, bucketId, string(madmin.ReplicationService)) // TODO: O(n^2)
	if err != nil {
		return
	}

	for _, target := range targets {
		if target.Endpoint == hostname && target.TargetBucket == bucketId {
			arn = target.Arn
			return
		}
	}
	adminClient, err := madmin.NewWithOptions(hostname, &madmin.Options{
		Secure: false, // TODO: true
		Creds:  cred,
	})
	if err != nil {
		return
	}
	nameHash := base32.StdEncoding.EncodeToString(sha3.New256().Sum([]byte(r.swarm.OwnName)))
	resultCredentials, err := adminClient.AddServiceAccount(ctx, madmin.AddServiceAccountReq{
		Name:        (ReplicationServiceAccountPrefix + nameHash)[:32],
		Description: fmt.Sprintf("Access key for replicating this bucket to %s, making sure your data is available even if this node fails!", r.swarm.OwnName),
	})
	if err != nil {
		err = fmt.Errorf("AddServiceAccount: %w", err)
		return
	}

	fmt.Printf("!! http://%s:%s@%s/%s \n", resultCredentials.AccessKey, resultCredentials.SecretKey, hostname, bucketId)

	arn, err = r.minioAdmin.SetRemoteTarget(ctx, bucketId, &madmin.BucketTarget{
		SourceBucket: bucketId,
		Endpoint:     hostname,
		Secure:       false, // TODO: true
		Credentials:  &resultCredentials,
		TargetBucket: bucketId,
		Type:         madmin.ReplicationService,
		DisableProxy: false,
	})
	if err != nil {
		err = fmt.Errorf("SetRemoteTarget: %w", err)
		return
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
