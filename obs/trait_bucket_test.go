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

package obs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSetBucketReplicationInput_Trans_ShouldReturnValidParams_WhenValidInput tests trans method
func TestSetBucketReplicationInput_Trans_ShouldReturnValidParams_WhenValidInput(t *testing.T) {
	// Arrange
	input := SetBucketReplicationInput{
		Bucket: "test-bucket",
		ReplicationConfiguration: ReplicationConfiguration{
			Rules: []ReplicationRule{
				{
					ID:     "rule-1",
					Status: RuleStatusEnabled,
					Destination: ReplicationDestination{
						Bucket: "dest-bucket",
					},
				},
			},
		},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	// 检查 replication 子资源参数是否存在
	_, exists := params["replication"]
	assert.True(t, exists, "replication sub-resource parameter should exist")
	assert.NotEmpty(t, headers["Content-MD5"])
}

// TestSetBucketReplicationInput_Trans_ShouldHandleEmptyRules_WhenNoRulesProvided tests empty rules
func TestSetBucketReplicationInput_Trans_ShouldHandleEmptyRules_WhenNoRulesProvided(t *testing.T) {
	// Arrange
	input := SetBucketReplicationInput{
		Bucket:                   "test-bucket",
		ReplicationConfiguration: ReplicationConfiguration{Rules: []ReplicationRule{}},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	_, exists := params["replication"]
	assert.True(t, exists)
}

// TestSetBucketReplicationInput_Trans_ShouldValidateSize_WhenConfigurationTooLarge tests size validation
func TestSetBucketReplicationInput_Trans_ShouldValidateSize_WhenConfigurationTooLarge(t *testing.T) {
	// Arrange - 创建一个超过 50KB 的配置
	rules := make([]ReplicationRule, 200)
	for i := 0; i < 200; i++ {
		longID := strings.Repeat("a", 100) + string(rune(i))
		longPrefix := strings.Repeat("prefix/", 100)
		longBucket := strings.Repeat("bucket-", 50)
		longLocation := strings.Repeat("location-", 50)
		rules[i] = ReplicationRule{
			ID:     longID,
			Status: RuleStatusEnabled,
			Prefix: &ReplicationPrefix{
				PrefixSet: PrefixSet{
					Prefixes: []string{longPrefix},
				},
			},
			Destination: ReplicationDestination{
				Bucket:       longBucket,
				StorageClass: StorageClassStandard,
				Location:     longLocation,
			},
			HistoricalObjectReplication: strings.Repeat("enabled", 100),
		}
	}
	input := SetBucketReplicationInput{
		Bucket:                   "test-bucket",
		ReplicationConfiguration: ReplicationConfiguration{Rules: rules},
	}

	// Act
	_, _, _, err := input.trans(true)

	// Assert
	// 配置过大应该返回错误
	assert.Error(t, err)
}

// TestSetBucketReplicationInput_Trans_ShouldHandleMultipleRules_WhenMultipleRulesProvided tests multiple rules
func TestSetBucketReplicationInput_Trans_ShouldHandleMultipleRules_WhenMultipleRulesProvided(t *testing.T) {
	// Arrange
	rules := make([]ReplicationRule, 10)
	for i := 0; i < 10; i++ {
		rules[i] = ReplicationRule{
			ID:     "rule-" + string(rune('0'+i)),
			Status: RuleStatusEnabled,
			Destination: ReplicationDestination{
				Bucket: "dest-bucket",
			},
		}
	}
	input := SetBucketReplicationInput{
		Bucket:                   "test-bucket",
		ReplicationConfiguration: ReplicationConfiguration{Rules: rules},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	_, exists := params["replication"]
	assert.True(t, exists)
	assert.NotEmpty(t, headers["Content-MD5"])
}

// TestSetBucketDisPolicyInput_Trans_ShouldReturnValidParams_WhenValidInput tests DIS policy trans method
func TestSetBucketDisPolicyInput_Trans_ShouldReturnValidParams_WhenValidInput(t *testing.T) {
	// Arrange
	input := SetBucketDisPolicyInput{
		Bucket: "test-bucket",
		DisPolicyConfiguration: DisPolicyConfiguration{
			Rules: []DisPolicyRule{
				{
					ID:      "rule-1",
					Stream:  "test-stream",
					Project: "test-project-id",
					Events:  []string{"ObjectCreated:*", "ObjectRemoved:*"},
					Prefix:  "images/",
					Suffix:  ".jpg",
					Agency:  "test-agency",
				},
			},
		},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	_, exists := params["disPolicy"]
	assert.True(t, exists, "disPolicy sub-resource parameter should exist")
	assert.Equal(t, []string{"application/json"}, headers["content-type"])
}

// TestSetBucketDisPolicyInput_Trans_ShouldHandleEmptyRules_WhenNoRulesProvided tests empty rules
func TestSetBucketDisPolicyInput_Trans_ShouldHandleEmptyRules_WhenNoRulesProvided(t *testing.T) {
	// Arrange
	input := SetBucketDisPolicyInput{
		Bucket: "test-bucket",
		DisPolicyConfiguration: DisPolicyConfiguration{
			Rules: []DisPolicyRule{},
		},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	_, exists := params["disPolicy"]
	assert.True(t, exists)
	assert.Equal(t, []string{"application/json"}, headers["content-type"])
}

// TestSetBucketDisPolicyInput_Trans_ShouldHandleMultipleRules_WhenMultipleRulesProvided tests multiple rules
func TestSetBucketDisPolicyInput_Trans_ShouldHandleMultipleRules_WhenMultipleRulesProvided(t *testing.T) {
	// Arrange
	input := SetBucketDisPolicyInput{
		Bucket: "test-bucket",
		DisPolicyConfiguration: DisPolicyConfiguration{
			Rules: []DisPolicyRule{
				{
					ID:      "rule-1",
					Stream:  "test-stream-1",
					Project: "test-project-id",
					Events:  []string{"ObjectCreated:*"},
					Agency:  "test-agency",
				},
				{
					ID:      "rule-2",
					Stream:  "test-stream-2",
					Project: "test-project-id",
					Events:  []string{"ObjectRemoved:*"},
					Agency:  "test-agency",
				},
			},
		},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	_, exists := params["disPolicy"]
	assert.True(t, exists)
	assert.Equal(t, []string{"application/json"}, headers["content-type"])
}

// TestSetBucketDisPolicyInput_Trans_ShouldHandleMaxRules_WhenTenRulesProvided tests max 10 rules
func TestSetBucketDisPolicyInput_Trans_ShouldHandleMaxRules_WhenTenRulesProvided(t *testing.T) {
	// Arrange
	rules := make([]DisPolicyRule, 10)
	for i := 0; i < 10; i++ {
		rules[i] = DisPolicyRule{
			ID:      "rule-" + string(rune('1'+i)),
			Stream:  "test-stream",
			Project: "test-project-id",
			Events:  []string{"ObjectCreated:*"},
			Agency:  "test-agency",
		}
	}
	input := SetBucketDisPolicyInput{
		Bucket: "test-bucket",
		DisPolicyConfiguration: DisPolicyConfiguration{
			Rules: rules,
		},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	_, exists := params["disPolicy"]
	assert.True(t, exists)
	assert.Equal(t, []string{"application/json"}, headers["content-type"])
}

// TestSubResourceType_ShouldHaveCorrectValue_WhenReplicationSubResource tests sub-resource constant
func TestSubResourceType_ShouldHaveCorrectValue_WhenReplicationSubResource(t *testing.T) {
	// Assert
	assert.Equal(t, SubResourceType("replication"), SubResourceReplication)
}

// TestSubResourceType_ShouldHaveCorrectValue_WhenDisPolicySubResource tests DIS policy sub-resource constant
func TestSubResourceType_ShouldHaveCorrectValue_WhenDisPolicySubResource(t *testing.T) {
	// Assert
	assert.Equal(t, SubResourceType("disPolicy"), SubResourceDisPolicy)
}

// TestReplicationRule_ShouldHandleDisabledStatus_WhenStatusIsDisabled tests disabled status
func TestReplicationRule_ShouldHandleDisabledStatus_WhenStatusIsDisabled(t *testing.T) {
	// Arrange
	input := SetBucketReplicationInput{
		Bucket: "test-bucket",
		ReplicationConfiguration: ReplicationConfiguration{
			Rules: []ReplicationRule{
				{
					ID:     "rule-1",
					Status: RuleStatusDisabled,
					Destination: ReplicationDestination{
						Bucket: "dest-bucket",
					},
				},
			},
		},
	}

	// Act
	params, headers, data, err := input.trans(true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	_, exists := params["replication"]
	assert.True(t, exists)
}
