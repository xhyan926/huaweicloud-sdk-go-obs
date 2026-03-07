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
	"errors"
	"fmt"
	"time"
)

// buildPostPolicyJSON builds POST policy JSON
func buildPostPolicyJSON(policy *PostPolicy) (string, error) {
	if policy == nil {
		return "", errors.New("policy is nil")
	}

	// 验证过期时间
	if policy.Expiration == "" {
		return "", errors.New("expiration is required")
	}

	// 验证条件
	if len(policy.Conditions) == 0 {
		return "", errors.New("at least one condition is required")
	}

	// 构建 JSON
	jsonData, err := json.Marshal(policy)
	if err != nil {
		return "", fmt.Errorf("failed to marshal policy: %v", err)
	}

	return string(jsonData), nil
}

// BuildPostPolicyExpiration builds expiration time string
func BuildPostPolicyExpiration(expiresIn int64) string {
	expirationTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
	return expirationTime.UTC().Format("2006-01-02T15:04:05.000Z")
}

// ValidatePostPolicy validates POST policy
func ValidatePostPolicy(policy *PostPolicy) error {
	if policy == nil {
		return errors.New("policy is nil")
	}

	// 验证过期时间
	if policy.Expiration == "" {
		return errors.New("expiration is required")
	}

	// 验证条件
	for i, cond := range policy.Conditions {
		if cond.Key == "" {
			return fmt.Errorf("condition at index %d has empty key", i)
		}
		if cond.Operator == "" {
			return fmt.Errorf("condition at index %d has empty operator", i)
		}
	}

	return nil
}

// CreatePostPolicyCondition creates a policy condition
func CreatePostPolicyCondition(operator, key string, value interface{}) PostPolicyCondition {
	return PostPolicyCondition{
		Operator: operator,
		Key:      key,
		Value:    value,
	}
}

// CreateBucketCondition creates a bucket condition
func CreateBucketCondition(bucket string) PostPolicyCondition {
	return CreatePostPolicyCondition(
		PostPolicyOpEquals,
		PostPolicyKeyBucket,
		bucket,
	)
}

// CreateKeyCondition creates a key condition
func CreateKeyCondition(key string) PostPolicyCondition {
	return CreatePostPolicyCondition(
		PostPolicyOpStartsWith,
		PostPolicyKeyKey,
		key,
	)
}

// CalculatePostPolicySignature calculates signature for POST policy
func CalculatePostPolicySignature(policyJSON, secretAccessKey string) (string, error) {
	if policyJSON == "" {
		return "", errors.New("policyJSON is empty")
	}
	if secretAccessKey == "" {
		return "", errors.New("secretAccessKey is empty")
	}

	// Base64 编码 Policy
	encodedPolicy := Base64Encode([]byte(policyJSON))

	// 使用 HMAC-SHA1 计算 POST Policy 签名
	// POST Policy 签名使用 HMAC-SHA1 算法
	signatureBytes := HmacSha1([]byte(secretAccessKey), []byte(encodedPolicy))
	return Base64Encode(signatureBytes), nil
}

// BuildPostPolicyToken builds complete token string
func BuildPostPolicyToken(ak, signature, policy string) string {
	return fmt.Sprintf("%s:%s:%s", ak, signature, policy)
}
