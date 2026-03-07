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
	"testing"

	"github.com/stretchr/testify/assert"
)

// PostPolicy 相关常量测试

func TestPostPolicyConditionKeys_ShouldHaveCorrectValues_GivenConstants(t *testing.T) {
	assert.Equal(t, "$bucket", PostPolicyKeyBucket)
	assert.Equal(t, "$key", PostPolicyKeyKey)
	assert.Equal(t, "$content-type", PostPolicyKeyContentType)
	assert.Equal(t, "$content-length", PostPolicyKeyContentLength)
}

func TestPostPolicyConditionOperators_ShouldHaveCorrectValues_GivenConstants(t *testing.T) {
	assert.Equal(t, "eq", PostPolicyOpEquals)
	assert.Equal(t, "starts-with", PostPolicyOpStartsWith)
	assert.Equal(t, "content-length-range", PostPolicyOpRange)
}

func TestPostPolicyCondition_ShouldSupportEquals_GivenStringValue(t *testing.T) {
	condition := PostPolicyCondition{
		Operator: PostPolicyOpEquals,
		Key:      PostPolicyKeyBucket,
		Value:    "test-bucket",
	}

	assert.Equal(t, PostPolicyOpEquals, condition.Operator)
	assert.Equal(t, PostPolicyKeyBucket, condition.Key)
	assert.Equal(t, "test-bucket", condition.Value)
}

func TestPostPolicyCondition_ShouldSupportStartsWith_GivenPrefix(t *testing.T) {
	condition := PostPolicyCondition{
		Operator: PostPolicyOpStartsWith,
		Key:      PostPolicyKeyKey,
		Value:    "uploads/",
	}

	assert.Equal(t, PostPolicyOpStartsWith, condition.Operator)
	assert.Equal(t, PostPolicyKeyKey, condition.Key)
	assert.Equal(t, "uploads/", condition.Value)
}

func TestPostPolicyCondition_ShouldSupportMultipleTypes_GivenDifferentValues(t *testing.T) {
	// 测试不同类型的值
	cases := []struct {
		name  string
		value interface{}
	}{
		{"Integer value", 12345},
		{"String value", "test"},
		{"Array value", []string{"value1", "value2"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			condition := PostPolicyCondition{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    tc.value,
			}
			assert.Equal(t, tc.value, condition.Value)
		})
	}
}

func TestPostPolicy_ShouldHaveRequiredFields_GivenValidPolicy(t *testing.T) {
	policy := PostPolicy{
		Expiration: "2026-03-07T12:00:00.000Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	assert.Equal(t, "2026-03-07T12:00:00.000Z", policy.Expiration)
	assert.NotEmpty(t, policy.Conditions)
	assert.Equal(t, 1, len(policy.Conditions))
}

func TestPostPolicy_ShouldSerializeToJSON_GivenValidPolicy(t *testing.T) {
	policy := PostPolicy{
		Expiration: "2026-03-07T12:00:00.000Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	jsonData, err := policy.MarshalJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	// 验证 JSON 包含正确的条件格式（数组格式）
	assert.Contains(t, string(jsonData), "conditions")
	assert.Contains(t, string(jsonData), "expiration")
}

func TestPostPolicy_ShouldIncludeConditions_GivenMultipleConditions(t *testing.T) {
	policy := PostPolicy{
		Expiration: "2026-03-07T12:00:00.000Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
			{
				Operator: PostPolicyOpStartsWith,
				Key:      PostPolicyKeyKey,
				Value:    "uploads/",
			},
			{
				Operator: PostPolicyOpRange,
				Key:      PostPolicyKeyContentLength,
				Value:    []interface{}{1, 10 * 1024 * 1024},
			},
		},
	}

	jsonData, err := policy.MarshalJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	// 验证包含所有三个条件
	assert.Contains(t, string(jsonData), "\"eq\"")
	assert.Contains(t, string(jsonData), "\"starts-with\"")
	assert.Contains(t, string(jsonData), "\"content-length-range\"")
}
