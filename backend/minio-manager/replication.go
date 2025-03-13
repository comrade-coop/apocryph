package main

import (
	"context"
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/replication"
	"golang.org/x/crypto/sha3"
)

var ReplicationServiceAccountPrefix = "Replicator-"
var ReplicationCustomTokenIamArn = "arn:minio:iam:::role/idmp-swieauth"

type replicationNode struct {
	url *url.URL
	minio *minio.Client
	admin *madmin.AdminClient
	serviceAccount madmin.Credentials
}

func newReplicationNode(ctx context.Context, url *url.URL, tokenSigner *TokenSigner, accountDescription string) (result *replicationNode, err error) {
	token, err := tokenSigner.GetReplicationToken("all")
	if err != nil {
		return
	}
	cred, err := credentials.NewCustomTokenCredentials(url.String(), token, ReplicationCustomTokenIamArn)
	if err != nil {
		return
	}
	minio, err := minio.New(url.Host, &minio.Options{
		Secure: url.Scheme == "https",
		Creds:  cred,
	})
	if err != nil {
		return
	}
	admin, err := madmin.NewWithOptions(url.Host, &madmin.Options{
		Secure: url.Scheme == "https",
		Creds:  cred,
	})
	if err != nil {
		return
	}
	accountDescriptionHash := base32.StdEncoding.EncodeToString(sha3.New256().Sum([]byte(accountDescription)))
	replicationCredential, err := admin.AddServiceAccount(ctx, madmin.AddServiceAccountReq{
		Name:        (ReplicationServiceAccountPrefix + accountDescriptionHash)[:32],
		Description: accountDescription,
	})
	result = &replicationNode{
		url,
		minio,
		admin,
		replicationCredential,
	}
	return
}

func ConfigureAllBucketsReplication(
	ctx context.Context,
	ownUrl *url.URL,
	remoteUrl *url.URL,
	tokenSigner *TokenSigner,
) (err error) {
	source, err := newReplicationNode(ctx, remoteUrl, tokenSigner, fmt.Sprintf("Used for replication to %s", remoteUrl))
	if err != nil {
		return
	}
	destination, err := newReplicationNode(ctx, ownUrl, tokenSigner, fmt.Sprintf("Used for replication to %s", remoteUrl))
	if err != nil {
		return
	}
	
	remoteBuckets, err := source.minio.ListBuckets(ctx)
	if err != nil {
		return err
	}

	errorList := []error{}
	for _, bucket := range remoteBuckets {
		err := configureBucketReplication(ctx, bucket.Name, source, destination)

		if err != nil {
			log.Printf("Setting replication failed for %s: %v\n", bucket.Name, err)
			errorList = append(errorList, err)
		} else {
			log.Printf("Setting replication success for %s!\n", bucket.Name)
		}
	}
	
	return errors.Join(errorList...)
}

func configureBucketReplication(
	ctx context.Context,
	bucketId string,
	source *replicationNode,
	destination *replicationNode,
) (err error) {
	_ = source.minio.EnableVersioning(ctx, bucketId)
	
	err = destination.minio.MakeBucket(ctx, bucketId, minio.MakeBucketOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code != "BucketAlreadyOwnedByYou" {
			err = fmt.Errorf("MakeBucket: %w", err)
			return
		}
	}
	_ = destination.minio.EnableVersioning(ctx, bucketId)
		
	err = syncBucketPolicy(ctx, bucketId, source, destination)
	if err != nil {
		return
	}
	
	err = addBucketReplicationRule(ctx, bucketId, source, destination)
	if err != nil {
		return
	}

	return nil
}

func syncBucketPolicy(ctx context.Context, bucketId string, source *replicationNode, destination *replicationNode) (err error) {
	ownPolicy, err := source.minio.GetBucketPolicy(ctx, bucketId)
	if err != nil {
		return fmt.Errorf("GetBucketPolicy: %w", err)
	}
	err = destination.minio.SetBucketPolicy(ctx, bucketId, ownPolicy)
	if err != nil {
		return fmt.Errorf("SetBucketPolicy: %w", err)
	}
	return
}

func addBucketReplicationRule(ctx context.Context, bucketId string, source *replicationNode, destination *replicationNode) (err error) {
	replicationConfig, err := source.minio.GetBucketReplication(ctx, bucketId)
	if err != nil {
		response := minio.ToErrorResponse(err)
		if response.Code != "" {
			println("If this is not an error, this code to replication.go:", response.Code)
			replicationConfig = replication.Config{}
		} else {
			return fmt.Errorf("GetBucketReplication: %w", err)
		}
	}

	arn, err := getRemoteTargetArn(ctx, bucketId, source, destination)
	if err != nil {
		return
	}
	
	exists := false
	highestPriority := 0
	for _, rule := range replicationConfig.Rules {
		if rule.ID == destination.url.Host {
			exists = true
		}
		if highestPriority < rule.Priority {
			highestPriority = rule.Priority
		}
	}
	
	replicationRule := replication.Options{
		ID:                      destination.url.Host,
		RuleStatus:              "enable",
		ExistingObjectReplicate: "enable",
		ReplicateDeletes:        "enable",
		ReplicateDeleteMarkers:  "enable",
		ReplicaSync:             "enable",
		DestBucket:              arn,
	}
	
	if exists {
		err = replicationConfig.EditRule(replicationRule)
	} else {
		replicationRule.Priority = fmt.Sprint(highestPriority + 1)
		err = replicationConfig.AddRule(replicationRule)
	}
	if err != nil {
		return fmt.Errorf("AddRule: %w", err)
	}

	err = source.minio.SetBucketReplication(ctx, bucketId, replicationConfig)
	if err != nil {
		return fmt.Errorf("SetBucketReplication: %w", err)
	}
	return
}

func getRemoteTargetArn(ctx context.Context, bucketId string, source *replicationNode, destination *replicationNode) (arn string, err error) {
	targets, err := source.admin.ListRemoteTargets(ctx, bucketId, string(madmin.ReplicationService))
	if err != nil {
		err = fmt.Errorf("ListRemoteTargets: %w", err)
		return
	}
	
	expectedTarget := madmin.BucketTarget{
		SourceBucket: bucketId,
		Endpoint:     destination.url.Host,
		Secure:       destination.url.Scheme == "https",
		Credentials:  &destination.serviceAccount,
		TargetBucket: bucketId,
		Type:         madmin.ReplicationService,
		DisableProxy: false,
	}
	
	for _, target := range targets {
		if target.Endpoint == expectedTarget.Endpoint && 
			target.Secure == expectedTarget.Secure && 
			target.TargetBucket == expectedTarget.TargetBucket {
			arn = target.Arn
			return
		}
	}

	arn, err = source.admin.SetRemoteTarget(ctx, bucketId, &expectedTarget)
	if err != nil {
		err = fmt.Errorf("SetRemoteTarget: %w", err)
		return
	}

	return
}
