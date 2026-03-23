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
	"encoding/json"
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_SetBucketDisPolicy_ShouldSucceed tests setting bucket DIS policy configuration
func TestIntegration_SetBucketDisPolicy_ShouldSucceed(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("dis-rule")
	disCfg := getDisConfig(t)

	input := &obs.SetBucketDisPolicyInput{
		Bucket: bucket,
		DisPolicyConfiguration: obs.DisPolicyConfiguration{
			Rules: []obs.DisPolicyRule{
				{
					ID:      ruleID,
					Stream:  disCfg.Stream,
					Project: disCfg.Project,
					Events:  []string{"ObjectCreated:*", "ObjectRemoved:*"},
					Prefix:  "test-dis/",
					Suffix:  ".jpg",
					Agency:  disCfg.Agency,
				},
			},
		},
	}

	// Act
	output, err := client.SetBucketDisPolicy(input)

	// Assert
	require.NoError(t, err, "SetBucketDisPolicy should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.Equal(t, int64(200), output.StatusCode, "Status code should be 200")

	// Cleanup
	cleanupDisPolicy(t, client, bucket)
}

// TestIntegration_GetBucketDisPolicy_ShouldReturnConfig tests getting bucket DIS policy configuration
func TestIntegration_GetBucketDisPolicy_ShouldReturnConfig(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("dis-get-rule")
	disCfg := getDisConfig(t)

	setInput := &obs.SetBucketDisPolicyInput{
		Bucket: bucket,
		DisPolicyConfiguration: obs.DisPolicyConfiguration{
			Rules: []obs.DisPolicyRule{
				{
					ID:      ruleID,
					Stream:  disCfg.Stream,
					Project: disCfg.Project,
					Events:  []string{"ObjectCreated:*"},
					Prefix:  "test-get-dis/",
					Agency:  disCfg.Agency,
				},
			},
		},
	}
	_, err := client.SetBucketDisPolicy(setInput)
	require.NoError(t, err, "SetBucketDisPolicy should succeed")

	// Act
	output, err := client.GetBucketDisPolicy(bucket)

	// Assert
	require.NoError(t, err, "GetBucketDisPolicy should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.NotEmpty(t, output.DisPolicyConfiguration, "Configuration should not be empty")

	// Parse JSON to verify structure
	var config obs.DisPolicyConfiguration
	err = json.Unmarshal([]byte(output.DisPolicyConfiguration), &config)
	require.NoError(t, err, "Should be able to parse DIS policy configuration")
	assert.NotEmpty(t, config.Rules, "Should have at least one rule")
	assert.Equal(t, ruleID, config.Rules[0].ID)

	// Cleanup
	cleanupDisPolicy(t, client, bucket)
}

// TestIntegration_DeleteBucketDisPolicy_ShouldSucceed tests deleting bucket DIS policy configuration
func TestIntegration_DeleteBucketDisPolicy_ShouldSucceed(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("dis-delete-rule")
	disCfg := getDisConfig(t)

	setInput := &obs.SetBucketDisPolicyInput{
		Bucket: bucket,
		DisPolicyConfiguration: obs.DisPolicyConfiguration{
			Rules: []obs.DisPolicyRule{
				{
					ID:      ruleID,
					Stream:  disCfg.Stream,
					Project: disCfg.Project,
					Events:  []string{"ObjectRemoved:*"},
					Agency:  disCfg.Agency,
				},
			},
		},
	}
	_, err := client.SetBucketDisPolicy(setInput)
	require.NoError(t, err, "SetBucketDisPolicy should succeed")

	// Act
	output, err := client.DeleteBucketDisPolicy(bucket)

	// Assert
	require.NoError(t, err, "DeleteBucketDisPolicy should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.Equal(t, int64(204), output.StatusCode, "Status code should be 204")
}

// TestIntegration_DisPolicy_ShouldHandleMultipleRules tests multiple DIS policy rules
func TestIntegration_DisPolicy_ShouldHandleMultipleRules(t *testing.T) {
	// Arrange
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID1 := generateTestID("dis-multi-1")
	ruleID2 := generateTestID("dis-multi-2")
	disCfg := getDisConfig(t)

	input := &obs.SetBucketDisPolicyInput{
		Bucket: bucket,
		DisPolicyConfiguration: obs.DisPolicyConfiguration{
			Rules: []obs.DisPolicyRule{
				{
					ID:      ruleID1,
					Stream:  disCfg.Stream,
					Project: disCfg.Project,
					Events:  []string{"ObjectCreated:*"},
					Prefix:  "images/",
					Agency:  disCfg.Agency,
				},
				{
					ID:      ruleID2,
					Stream:  disCfg.Stream,
					Project: disCfg.Project,
					Events:  []string{"ObjectRemoved:*"},
					Prefix:  "videos/",
					Agency:  disCfg.Agency,
				},
			},
		},
	}

	// Act
	output, err := client.SetBucketDisPolicy(input)

	// Assert
	require.NoError(t, err, "SetBucketDisPolicy with multiple rules should succeed")
	require.NotNil(t, output, "Output should not be nil")
	assert.Equal(t, int64(200), output.StatusCode)

	// Verify by getting the configuration
	getOutput, err := client.GetBucketDisPolicy(bucket)
	require.NoError(t, err)

	var config obs.DisPolicyConfiguration
	err = json.Unmarshal([]byte(getOutput.DisPolicyConfiguration), &config)
	require.NoError(t, err)
	assert.Len(t, config.Rules, 2, "Should have 2 rules")

	// Cleanup
	cleanupDisPolicy(t, client, bucket)
}

// TestIntegration_DisPolicy_ShouldSupportOnlyOBSSignature tests that DIS policy only supports OBS signature
func TestIntegration_DisPolicy_ShouldSupportOnlyOBSSignature(t *testing.T) {
	// Arrange - Create client with OBS signature
	client := createClient(t, obs.SignatureObs)
	bucket := getTestBucket(t)
	ruleID := generateTestID("dis-obs-only-rule")
	disCfg := getDisConfig(t)

	input := &obs.SetBucketDisPolicyInput{
		Bucket: bucket,
		DisPolicyConfiguration: obs.DisPolicyConfiguration{
			Rules: []obs.DisPolicyRule{
				{
					ID:      ruleID,
					Stream:  disCfg.Stream,
					Project: disCfg.Project,
					Events:  []string{"ObjectCreated:*"},
					Agency:  disCfg.Agency,
				},
			},
		},
	}

	// Act
	_, err := client.SetBucketDisPolicy(input)

	// Assert
	// OBS signature should work
	require.NoError(t, err, "SetBucketDisPolicy with OBS signature should succeed")

	// Cleanup
	cleanupDisPolicy(t, client, bucket)
}
