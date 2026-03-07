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
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// buildPostPolicyJSON tests

func TestBuildPostPolicyJSON_ShouldGenerateValidJSON_GivenValidPolicy(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	result, err := buildPostPolicyJSON(policy)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(result), &parsed)
	assert.NoError(t, err)
	assert.Equal(t, "2024-12-31T23:59:59.000Z", parsed["expiration"])
}

func TestBuildPostPolicyJSON_ShouldIncludeAllConditions_GivenMultipleConditions(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
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
		},
	}

	result, err := buildPostPolicyJSON(policy)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(result), &parsed)
	assert.NoError(t, err)
	assert.NotNil(t, parsed["conditions"])
}

func TestBuildPostPolicyJSON_ShouldReturnError_GivenNilPolicy(t *testing.T) {
	_, err := buildPostPolicyJSON(nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policy is nil")
}

func TestBuildPostPolicyJSON_ShouldReturnError_GivenEmptyExpiration(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	_, err := buildPostPolicyJSON(policy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expiration is required")
}

func TestBuildPostPolicyJSON_ShouldReturnError_GivenEmptyConditions(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{},
	}

	_, err := buildPostPolicyJSON(policy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one condition is required")
}

// BuildPostPolicyExpiration tests

func TestBuildPostPolicyExpiration_ShouldGenerateCorrectFormat_GivenSeconds(t *testing.T) {
	expiresIn := int64(3600) // 1 hour

	result := BuildPostPolicyExpiration(expiresIn)

	assert.NotEmpty(t, result)

	// 验证 ISO 8601 格式
	pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`
	matched, _ := regexp.MatchString(pattern, result)
	assert.True(t, matched, "expiration should match ISO 8601 format")
}

func TestBuildPostPolicyExpiration_ShouldUseUTC_GivenLocalTime(t *testing.T) {
	expiresIn := int64(3600)

	result := BuildPostPolicyExpiration(expiresIn)

	assert.NotEmpty(t, result)

	// 验证以 Z 结尾表示 UTC
	assert.Contains(t, result, "Z")
}

func TestBuildPostPolicyExpiration_ShouldBeFutureTime_GivenPositiveSeconds(t *testing.T) {
	expiresIn := int64(3600)

	result := BuildPostPolicyExpiration(expiresIn)

	expirationTime, err := time.Parse("2006-01-02T15:04:05.000Z", result)
	assert.NoError(t, err)
	assert.True(t, expirationTime.After(time.Now()))
}

func TestBuildPostPolicyExpiration_ShouldGenerateTime_GivenZeroSeconds(t *testing.T) {
	expiresIn := int64(0)

	result := BuildPostPolicyExpiration(expiresIn)

	assert.NotEmpty(t, result)
	pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`
	matched, _ := regexp.MatchString(pattern, result)
	assert.True(t, matched)
}

// ValidatePostPolicy tests

func TestValidatePostPolicy_ShouldReturnNil_GivenValidPolicy(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	err := ValidatePostPolicy(policy)

	assert.NoError(t, err)
}

func TestValidatePostPolicy_ShouldReturnError_GivenNilPolicy(t *testing.T) {
	err := ValidatePostPolicy(nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policy is nil")
}

func TestValidatePostPolicy_ShouldReturnError_GivenEmptyExpiration(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	err := ValidatePostPolicy(policy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expiration is required")
}

func TestValidatePostPolicy_ShouldReturnError_GivenConditionWithEmptyKey(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: PostPolicyOpEquals,
				Key:      "",
				Value:    "test-bucket",
			},
		},
	}

	err := ValidatePostPolicy(policy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has empty key")
}

func TestValidatePostPolicy_ShouldReturnError_GivenConditionWithEmptyOperator(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{
			{
				Operator: "",
				Key:      PostPolicyKeyBucket,
				Value:    "test-bucket",
			},
		},
	}

	err := ValidatePostPolicy(policy)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has empty operator")
}

// CreatePostPolicyCondition tests

func TestCreatePostPolicyCondition_ShouldCreateCorrectCondition_GivenValidParameters(t *testing.T) {
	result := CreatePostPolicyCondition(
		PostPolicyOpEquals,
		PostPolicyKeyBucket,
		"test-bucket",
	)

	assert.Equal(t, PostPolicyOpEquals, result.Operator)
	assert.Equal(t, PostPolicyKeyBucket, result.Key)
	assert.Equal(t, "test-bucket", result.Value)
}

func TestCreatePostPolicyCondition_ShouldSupportDifferentValueTypes_GivenVariousInputs(t *testing.T) {
	testCases := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "String value",
			key:   PostPolicyKeyContentType,
			value: "image/jpeg",
		},
		{
			name:  "Integer value",
			key:   PostPolicyKeyContentLength,
			value: 1024,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CreatePostPolicyCondition(PostPolicyOpEquals, tc.key, tc.value)
			assert.Equal(t, PostPolicyOpEquals, result.Operator)
			assert.Equal(t, tc.key, result.Key)
			assert.Equal(t, tc.value, result.Value)
		})
	}
}

// CreateBucketCondition tests

func TestCreateBucketCondition_ShouldCreateCorrectCondition_GivenBucket(t *testing.T) {
	bucket := "test-bucket"

	result := CreateBucketCondition(bucket)

	assert.Equal(t, PostPolicyOpEquals, result.Operator)
	assert.Equal(t, PostPolicyKeyBucket, result.Key)
	assert.Equal(t, bucket, result.Value)
}

// CreateKeyCondition tests

func TestCreateKeyCondition_ShouldCreateCorrectCondition_GivenKey(t *testing.T) {
	key := "uploads/"

	result := CreateKeyCondition(key)

	assert.Equal(t, PostPolicyOpStartsWith, result.Operator)
	assert.Equal(t, PostPolicyKeyKey, result.Key)
	assert.Equal(t, key, result.Value)
}

// CalculatePostPolicySignature tests

func TestCalculatePostPolicySignature_ShouldCalculateCorrectSignature_GivenValidInput(t *testing.T) {
	policyJSON := `{"expiration":"2024-12-31T23:59:59.000Z","conditions":[["eq","$bucket","test-bucket"]]}`
	secretKey := "test-secret-key"

	result, err := CalculatePostPolicySignature(policyJSON, secretKey)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// 验证结果可以解码（Base64）
	_, err = Base64Decode(result)
	assert.NoError(t, err)
}

func TestCalculatePostPolicySignature_ShouldReturnError_GivenEmptyPolicyJSON(t *testing.T) {
	secretKey := "test-secret-key"

	_, err := CalculatePostPolicySignature("", secretKey)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policyJSON is empty")
}

func TestCalculatePostPolicySignature_ShouldReturnError_GivenEmptySecretKey(t *testing.T) {
	policyJSON := `{"expiration":"2024-12-31T23:59:59.000Z","conditions":[["eq","$bucket","test-bucket"]]}`

	_, err := CalculatePostPolicySignature(policyJSON, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "secretAccessKey is empty")
}

// BuildPostPolicyToken tests

func TestBuildPostPolicyToken_ShouldGenerateCorrectFormat_GivenValidComponents(t *testing.T) {
	ak := "AKIAIOSFODNN7EXAMPLE"
	signature := "test-signature"
	policy := "eyJleHBpcmF0aW9uIjoiMjAyNC0xMi0zMVQyMzo1OTo1OS4wMDBaIiwiY29uZGl0aW9ucyI6W1siZXEiLCIkYnVja2V0IiwidGVzdC1idWNrZXQiXV19"

	result := BuildPostPolicyToken(ak, signature, policy)

	expected := fmt.Sprintf("%s:%s:%s", ak, signature, policy)
	assert.Equal(t, expected, result)
}

func TestBuildPostPolicyToken_ShouldIncludeAllComponents_GivenValidInput(t *testing.T) {
	ak := "AKIAIOSFODNN7EXAMPLE"
	signature := "test-signature"
	policy := "test-policy"

	result := BuildPostPolicyToken(ak, signature, policy)

	assert.Contains(t, result, ak)
	assert.Contains(t, result, signature)
	assert.Contains(t, result, policy)
	assert.Equal(t, 2, strings.Count(result, ":"))
}

func TestBuildPostPolicyToken_ShouldHandleEmptyComponents_GivenEmptyStrings(t *testing.T) {
	result := BuildPostPolicyToken("", "", "")

	assert.Equal(t, "::", result)
}

// Additional test for content-length-range condition

func TestCreateContentLengthRangeCondition_ShouldCreateCorrectCondition_GivenRange(t *testing.T) {
	min := int64(0)
	max := int64(10485760) // 10MB

	condition := CreatePostPolicyCondition(
		PostPolicyOpRange,
		PostPolicyKeyContentLength,
		[]interface{}{min, max},
	)

	assert.Equal(t, PostPolicyOpRange, condition.Operator)
	assert.Equal(t, PostPolicyKeyContentLength, condition.Key)

	// Value should be an array
	valueArray, ok := condition.Value.([]interface{})
	assert.True(t, ok)
	assert.Len(t, valueArray, 2)
	assert.Equal(t, int64(0), valueArray[0])
	assert.Equal(t, int64(10485760), valueArray[1])
}

// Policy JSON marshaling tests

func TestPostPolicy_ShouldMarshalToJSONWithConditions_GivenMultipleConditions(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{
			CreateBucketCondition("test-bucket"),
			CreateKeyCondition("uploads/"),
			CreatePostPolicyCondition(PostPolicyOpEquals, PostPolicyKeyContentType, "image/jpeg"),
		},
	}

	jsonBytes, err := json.Marshal(policy)

	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, `"expiration"`)
	assert.Contains(t, jsonStr, `"conditions"`)
	assert.Contains(t, jsonStr, "test-bucket")
	assert.Contains(t, jsonStr, "uploads/")
	assert.Contains(t, jsonStr, "image/jpeg")
}

func TestPostPolicyCondition_ShouldMarshalToJSONArray_GivenCondition(t *testing.T) {
	condition := PostPolicyCondition{
		Operator: PostPolicyOpEquals,
		Key:      PostPolicyKeyBucket,
		Value:    "test-bucket",
	}

	jsonBytes, err := json.Marshal(condition)

	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "eq")
	assert.Contains(t, jsonStr, "$bucket")
	assert.Contains(t, jsonStr, "test-bucket")
}

// Edge cases and boundary conditions

func TestBuildPostPolicyExpiration_ShouldHandleVeryLargeDuration_GivenLargeSeconds(t *testing.T) {
	expiresIn := int64(86400 * 365) // 1 year

	result := BuildPostPolicyExpiration(expiresIn)

	assert.NotEmpty(t, result)
	pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`
	matched, _ := regexp.MatchString(pattern, result)
	assert.True(t, matched)

	// Verify it's far in the future
	expirationTime, _ := time.Parse("2006-01-02T15:04:05.000Z", result)
	yearFromNow := expirationTime.Sub(time.Now()).Hours() / 24 / 365
	assert.Greater(t, yearFromNow, float64(0.9))
}

func TestValidatePostPolicy_ShouldAcceptSingleCondition_GivenValidCondition(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{
			CreateBucketCondition("test-bucket"),
		},
	}

	err := ValidatePostPolicy(policy)

	assert.NoError(t, err)
}

func TestValidatePostPolicy_ShouldAcceptMultipleConditions_GivenValidConditions(t *testing.T) {
	policy := &PostPolicy{
		Expiration: "2024-12-31T23:59:59.000Z",
		Conditions: []PostPolicyCondition{
			CreateBucketCondition("test-bucket"),
			CreateKeyCondition("uploads/"),
			CreatePostPolicyCondition(PostPolicyOpEquals, PostPolicyKeyContentType, "image/jpeg"),
			CreatePostPolicyCondition(PostPolicyOpEquals, "$acl", "public-read"),
		},
	}

	err := ValidatePostPolicy(policy)

	assert.NoError(t, err)
}


