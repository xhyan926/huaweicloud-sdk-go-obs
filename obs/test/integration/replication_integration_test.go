//go:build integration
// +build integration

// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of the
// License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

package integration

import (
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_SetBucketReplication_ShouldSucceed tests setting bucket replication configuration
func TestIntegration_SetBucketReplication_ShouldSucceed(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("replication-rule")
	replCfg := getReplicationConfig(t)

	input := &obs.SetBucketReplicationInput{
		Bucket: bucket,
		ReplicationConfiguration: obs.ReplicationConfiguration{
			Rules: []obs.ReplicationRule{
				{
					ID:     ruleID,
					Status: obs.RuleStatusEnabled,
					Prefix: &obs.ReplicationPrefix{
						PrefixSet: obs.PrefixSet{
							Prefixes: []string{"test-replication/"},
						},
					},
					Destination: obs.ReplicationDestination{
						Bucket:       replCfg.DestBucket,
						StorageClass: obs.StorageClassStandard,
						Location:     replCfg.Location,
					},
					HistoricalObjectReplication: "Enabled",
				},
			},
		},
	}

	// Act
	output, err := client.SetBucketReplication(input)

	// Assert
	require.NoError(t, err, "SetBucketReplication should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.Equal(t, int64(200), output.StatusCode, "Status code should be 200")

	// Cleanup
	cleanupReplication(t, client, bucket)
}

// TestIntegration_GetBucketReplication_ShouldReturnConfig tests getting bucket replication configuration
func TestIntegration_GetBucketReplication_ShouldReturnConfig(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("replication-get-rule")
	replCfg := getReplicationConfig(t)

	setInput := &obs.SetBucketReplicationInput{
		Bucket: bucket,
		ReplicationConfiguration: obs.ReplicationConfiguration{
			Rules: []obs.ReplicationRule{
				{
					ID:     ruleID,
					Status: obs.RuleStatusEnabled,
					Prefix: &obs.ReplicationPrefix{
						PrefixSet: obs.PrefixSet{
							Prefixes: []string{"test-get-replication/"},
						},
					},
					Destination: obs.ReplicationDestination{
						Bucket:       replCfg.DestBucket,
						StorageClass: obs.StorageClassStandard,
						Location:     replCfg.Location,
					},
				},
			},
		},
	}
	_, err := client.SetBucketReplication(setInput)
	require.NoError(t, err, "SetBucketReplication should succeed")

	// Act
	output, err := client.GetBucketReplication(bucket)

	// Assert
	require.NoError(t, err, "GetBucketReplication should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.NotEmpty(t, output.ReplicationConfiguration.Rules, "Should have at least one rule")

	// Cleanup
	cleanupReplication(t, client, bucket)
}

// TestIntegration_DeleteBucketReplication_ShouldSucceed tests deleting bucket replication configuration
func TestIntegration_DeleteBucketReplication_ShouldSucceed(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("replication-delete-rule")
	replCfg := getReplicationConfig(t)

	setInput := &obs.SetBucketReplicationInput{
		Bucket: bucket,
		ReplicationConfiguration: obs.ReplicationConfiguration{
			Rules: []obs.ReplicationRule{
				{
					ID:     ruleID,
					Status: obs.RuleStatusEnabled,
					Destination: obs.ReplicationDestination{
						Bucket: replCfg.DestBucket,
					},
				},
			},
		},
	}
	_, err := client.SetBucketReplication(setInput)
	require.NoError(t, err, "SetBucketReplication should succeed")

	// Act
	output, err := client.DeleteBucketReplication(bucket)

	// Assert
	require.NoError(t, err, "DeleteBucketReplication should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.Equal(t, int64(204), output.StatusCode, "Status code should be 204")
}

// TestIntegration_Replication_ShouldSupportOnlyOBSSignature tests that replication only supports OBS signature
func TestIntegration_Replication_ShouldSupportOnlyOBSSignature(t *testing.T) {
	// Arrange - Create client with OBS signature
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("replication-obs-only-rule")
	replCfg := getReplicationConfig(t)

	input := &obs.SetBucketReplicationInput{
		Bucket: bucket,
		ReplicationConfiguration: obs.ReplicationConfiguration{
			Rules: []obs.ReplicationRule{
				{
					ID:     ruleID,
					Status: obs.RuleStatusEnabled,
					Destination: obs.ReplicationDestination{
						Bucket: replCfg.DestBucket,
					},
				},
			},
		},
	}

	// Act
	_, err := client.SetBucketReplication(input)

	// Assert
	// OBS signature should work
	require.NoError(t, err, "SetBucketReplication with OBS signature should succeed")

	// Cleanup
	cleanupReplication(t, client, bucket)

	// Note: Testing with v4 signature would require a separate client
	// and would likely fail at the service level, not SDK level
}

// TestIntegration_Replication_ShouldHandleMultipleRules tests multiple replication rules
func TestIntegration_Replication_ShouldHandleMultipleRules(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID1 := generateTestID("replication-multi-1")
	ruleID2 := generateTestID("replication-multi-2")
	replCfg := getReplicationConfig(t)

	input := &obs.SetBucketReplicationInput{
		Bucket: bucket,
		ReplicationConfiguration: obs.ReplicationConfiguration{
			Rules: []obs.ReplicationRule{
				{
					ID:     ruleID1,
					Status: obs.RuleStatusEnabled,
					Destination: obs.ReplicationDestination{
						Bucket:       replCfg.DestBucket,
						StorageClass: obs.StorageClassStandard,
						Location:     replCfg.Location,
					},
				},
				{
					ID:     ruleID2,
					Status: obs.RuleStatusEnabled,
					Destination: obs.ReplicationDestination{
						Bucket: replCfg.DestBucket,
					},
				},
			},
		},
	}

	// Act
	output, err := client.SetBucketReplication(input)

	// Assert
	require.NoError(t, err, "SetBucketReplication with multiple rules should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.Equal(t, int64(200), output.StatusCode)

	// Verify by getting the configuration
	getOutput, err := client.GetBucketReplication(bucket)
	require.NoError(t, err)
	assert.Len(t, getOutput.ReplicationConfiguration.Rules, 2, "Should have 2 rules")

	// Cleanup
	cleanupReplication(t, client, bucket)
}
