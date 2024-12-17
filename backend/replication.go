package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/replication"
)

type ReplicationManager struct {
	minio      *minio.Client
	minioAdmin *madmin.AdminClient
	swarm      *Swarm
}

func NewReplicationManager(minioClient *minio.Client, minioAdmin *madmin.AdminClient, swarm *Swarm) *ReplicationManager {
	return &ReplicationManager{
		minio:      minioClient,
		minioAdmin: minioAdmin,
		swarm:      swarm,
	}
}

func (r *ReplicationManager) Start(ctx context.Context) {
	go func() {
		for {
			err := r.reconcilationLoop(ctx)
			if err != nil {
				log.Printf("Error while reconciling replications: %v", err)
			}

			time.Sleep(1 * time.Minute) // TODO: Maybe change to 5 minutes?
		}
	}()
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

// TODO: call UpdateBucket whenever the swarm detects that the peer is down
// TODO: call UpdateBucket whenever we get a new bucket
func (r *ReplicationManager) UpdateBucket(ctx context.Context, bucketId string) error {
	// Get object

	replicationConfig, err := r.minio.GetBucketReplication(ctx, bucketId)
	if err != nil {
		return err
	}
	targetTotalReplicas := 2 // TODO: Make this configurable

	// Get current status

	metrics, err := r.minio.GetBucketReplicationMetrics(ctx, bucketId)
	if err != nil {
		return err
	}

	aliveReplicas := 1 // 1, counting ourselves
	syncedReplicas := 0
	deadReplicas := []string{}
	for id, stats := range metrics.Stats {
		if stats.Failed.LastHour.Count > 10 {
			deadReplicas = append(deadReplicas, id)
		} else {
			aliveReplicas += 1
		}
		if stats.PendingCount < 10 {
			syncedReplicas += 1
		}
	}

	// Reconcile
	// TODO: Make sure this does what it's supposed to
	if aliveReplicas < targetTotalReplicas {
		for i := aliveReplicas; i < targetTotalReplicas; i += 1 {
			node, err := r.swarm.FindFreeNode()
			if err != nil {
				return err
			}

			// TODO: Create the remote bucket
			// TODO: Properly authenticate to the remote
			replicationConfig.AddRule(replication.Options{
				ID:                      node.String(),
				ExistingObjectReplicate: "enable",
				ReplicateDeletes:        "enable",
				ReplicateDeleteMarkers:  "enable",
				DestBucket:              "http://" + node.String() + "/" + bucketId,
			})
		}

	} else { // Cleanup
		for _, deadReplicationId := range deadReplicas {
			replicationConfig.RemoveRule(replication.Options{
				ID: deadReplicationId,
			})
		}
	}

	r.minio.SetBucketReplication(ctx, bucketId, replicationConfig)

	return nil
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
