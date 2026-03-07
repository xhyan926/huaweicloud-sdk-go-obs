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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PostPolicyCondition tests

func TestPostPolicyCondition_ShouldHaveRequiredFields_GivenValidCondition(t *testing.T) {
	condition := PostPolicyCondition{
		Operator: PostPolicyOpEquals,
		Key:      PostPolicyKeyBucket,
		Value:    "test-bucket",
	}

	assert.Equal(t, PostPolicyOpEquals, condition.Operator)
	assert.Equal(t, PostPolicyKeyBucket, condition.Key)
	assert.Equal(t, "test-bucket", condition.Value)
}

func TestPostPolicyCondition_ShouldSupportMultipleTypes_GivenDifferentValues(t *testing.T) {
	testCases := []struct {
		name     string
		operator string
		key      string
		value    interface{}
	}{
		{
			name:     "String value",
			operator: PostPolicyOpEquals,
			key:      PostPolicyKeyContentType,
			value:    "image/jpeg",
		},
		{
			name:     "Integer value",
			operator: PostPolicyOpEquals,
			key:      PostPolicyKeyContentLength,
			value:    1024,
		},
		{
			name:     "String array value",
			operator: PostPolicyOpStartsWith,
			key:      PostPolicyKeyKey,
			value:    "uploads/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			condition := PostPolicyCondition{
				Operator: tc.operator,
				Key:      tc.key,
				Value:    tc.value,
			}

			assert.Equal(t, tc.operator, condition.Operator)
			assert.Equal(t, tc.key, condition.Key)
			assert.Equal(t, tc.value, condition.Value)
		})
	}
}

// PostPolicy tests

func TestPostPolicy_ShouldHaveRequiredFields_GivenValidPolicy(t *testing.T) {
	policy := PostPolicy{
		Expiration: "2024-12-31T23:59:59Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	assert.Equal(t, "2024-12-31T23:59:59Z", policy.Expiration)
	assert.Len(t, policy.Conditions, 1)
	assert.Equal(t, PostPolicyKeyBucket, policy.Conditions[0].Key)
}

func TestPostPolicy_ShouldSerializeToJSON_GivenValidPolicy(t *testing.T) {
	policy := PostPolicy{
		Expiration: "2024-12-31T23:59:59Z",
	}

	data, err := json.Marshal(policy)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, "2024-12-31T23:59:59Z", result["expiration"])
}

func TestPostPolicy_ShouldIncludeConditions_GivenMultipleConditions(t *testing.T) {
	policy := PostPolicy{
		Expiration: "2024-12-31T23:59:59Z",
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
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyContentType,
				Value:    "image/jpeg",
			},
		},
	}

	assert.Len(t, policy.Conditions, 3)
	assert.Equal(t, PostPolicyOpEquals, policy.Conditions[0].Operator)
	assert.Equal(t, PostPolicyOpStartsWith, policy.Conditions[1].Operator)
	assert.Equal(t, PostPolicyOpEquals, policy.Conditions[2].Operator)
}

// CreatePostPolicyInput tests

func TestCreatePostPolicyInput_ShouldHaveRequiredFields_GivenValidInput(t *testing.T) {
	input := CreatePostPolicyInput{
		Bucket:    "test-bucket",
		Key:       "test-key",
		Expires:   1234567890,
		ExpiresIn: 3600,
		Acl:       "public-read",
	}

	assert.Equal(t, "test-bucket", input.Bucket)
	assert.Equal(t, "test-key", input.Key)
	assert.Equal(t, int64(1234567890), input.Expires)
	assert.Equal(t, int64(3600), input.ExpiresIn)
	assert.Equal(t, "public-read", input.Acl)
}

func TestCreatePostPolicyInput_ShouldAllowOptionalFields_GivenMinimalInput(t *testing.T) {
	input := CreatePostPolicyInput{
		Bucket:    "test-bucket",
		Key:       "test-key",
		ExpiresIn: 3600,
	}

	assert.Equal(t, "test-bucket", input.Bucket)
	assert.Equal(t, "test-key", input.Key)
	assert.Equal(t, int64(3600), input.ExpiresIn)
	assert.Equal(t, "", input.Acl)
}

// CreatePostPolicyOutput tests

func TestCreatePostPolicyOutput_ShouldHaveRequiredFields_GivenValidOutput(t *testing.T) {
	output := CreatePostPolicyOutput{
		BaseModel: BaseModel{
			StatusCode: 200,
			RequestId:  "test-request-id",
		},
		Policy:      "eyJleHBpcmF0aW9uIjoiMjAyNC0xMi0zMVQyMzo1OTo1OVoiLCJjb25kaXRpb25zIjpbWyJlcSIsIiRidWNrZXQiLCJ0ZXN0LWJ1Y2tldCJdfX0=",
		Signature:   "test-signature",
		Token:       "AK:test-signature:eyJleHBpcmF0aW9uIjoiMjAyNC0xMi0zMVQyMzo1OTo1OVoiLCJjb25kaXRpb25zIjpbWyJlcSIsIiRidWNrZXQiLCJ0ZXN0LWJ1Y2tldCJdfX0=",
		AccessKeyId: "AK",
	}

	assert.Equal(t, 200, output.StatusCode)
	assert.Equal(t, "test-request-id", output.RequestId)
	assert.Equal(t, "eyJleHBpcmF0aW9uIjoiMjAyNC0xMi0zMVQyMzo1OTo1OVoiLCJjb25kaXRpb25zIjpbWyJlcSIsIiRidWNrZXQiLCJ0ZXN0LWJ1Y2tldCJdfX0=", output.Policy)
	assert.Equal(t, "test-signature", output.Signature)
	assert.Equal(t, "AK:test-signature:eyJleHBpcmF0aW9uIjoiMjAyNC0xMi0zMVQyMzo1OTo1OVoiLCJjb25kaXRpb25zIjpbWyJlcSIsIiRidWNrZXQiLCJ0ZXN0LWJ1Y2tldCJdfX0=", output.Token)
	assert.Equal(t, "AK", output.AccessKeyId)
}

func TestCreatePostPolicyOutput_ShouldSerializeToJSON_GivenValidOutput(t *testing.T) {
	output := CreatePostPolicyOutput{
		Policy:      "test-policy",
		Signature:   "test-signature",
		Token:       "test-token",
		AccessKeyId: "AK",
	}

	data, err := json.Marshal(output)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, "test-policy", result["policy"])
	assert.Equal(t, "test-signature", result["signature"])
	assert.Equal(t, "test-token", result["token"])
	assert.Equal(t, "AK", result["accessKeyId"])
}

// Constants tests

func TestPostPolicyConditionKeys_ShouldHaveCorrectValues_GivenConstants(t *testing.T) {
	assert.Equal(t, "$bucket", PostPolicyKeyBucket)
	assert.Equal(t, "$key", PostPolicyKeyKey)
	assert.Equal(t, "$content-type", PostPolicyKeyContentType)
	assert.Equal(t, "$content-length", PostPolicyKeyContentLength)
}

func TestPostPolicyConditionOperators_ShouldHaveCorrectValues_GivenConstants(t *testing.T) {
	assert.Equal(t, "eq", PostPolicyOpEquals)
	assert.Equal(t, "starts-with", PostPolicyOpStartsWith)
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
