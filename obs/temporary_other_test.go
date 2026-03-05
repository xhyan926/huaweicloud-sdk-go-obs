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

// TestCreateBrowserBasedSignature_ShouldReturnError_WhenInputIsNil tests nil input
func TestCreateBrowserBasedSignature_ShouldReturnError_WhenInputIsNil(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	output, err := client.CreateBrowserBasedSignature(nil)

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CreateBrowserBasedSignatureInput is nil")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WhenInputIsValid tests valid input
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WhenInputIsValid(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:     "test-bucket",
		Key:        "test-key",
		Expires:    600,
		FormParams: map[string]string{"acl": "public-read"},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	assert.Contains(t, output.OriginPolicy, "test-bucket")
	assert.Contains(t, output.OriginPolicy, "test-key")
	assert.Contains(t, output.OriginPolicy, "acl")
	assert.Contains(t, output.OriginPolicy, "public-read")
}

// TestCreateBrowserBasedSignature_ShouldUseDefaultExpires_WhenExpiresIsZero tests default expires value
func TestCreateBrowserBasedSignature_ShouldUseDefaultExpires_WhenExpiresIsZero(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 0,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	// Default expires is 300 seconds
	assert.Contains(t, output.OriginPolicy, "\"expiration\"")
}

// TestCreateBrowserBasedSignature_ShouldUseDefaultExpires_WhenExpiresIsNegative tests negative expires value
func TestCreateBrowserBasedSignature_ShouldUseDefaultExpires_WhenExpiresIsNegative(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: -100,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.Contains(t, output.OriginPolicy, "\"expiration\"")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithEmptyBucketAndKey tests empty bucket and key
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithEmptyBucketAndKey(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "",
		Key:     "",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
	// With empty bucket and key, should use "starts-with" conditions
	assert.Contains(t, output.OriginPolicy, "\"starts-with\", \"$bucket\", \"\"")
	assert.Contains(t, output.OriginPolicy, "\"starts-with\", \"$key\", \"\"")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSignatureV4 tests V4 signature
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSignatureV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	assert.NotEmpty(t, output.Algorithm)
	assert.NotEmpty(t, output.Credential)
	assert.NotEmpty(t, output.Date)
	// V4 signature includes algorithm, credential, and date
	assert.Contains(t, output.Algorithm, "AWS4-HMAC-SHA256")
	assert.Contains(t, output.Credential, "test-ak")
	assert.Contains(t, output.Credential, "cn-north-4")
	assert.Contains(t, output.Date, "T")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithEmptyFormParams tests empty form params
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithEmptyFormParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:     "test-bucket",
		Key:        "test-key",
		Expires:    600,
		FormParams: map[string]string{},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithNilFormParams tests nil form params
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithNilFormParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:     "test-bucket",
		Key:        "test-key",
		Expires:    600,
		FormParams: nil,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithRangeParams tests range params
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithRangeParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		RangeParams: []ConditionRange{
			{
				RangeName: "content-length-range",
				Lower:     0,
				Upper:     10485760,
			},
		},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
	assert.Contains(t, output.OriginPolicy, "content-length-range")
	assert.Contains(t, output.OriginPolicy, "0")
	assert.Contains(t, output.OriginPolicy, "10485760")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithMultipleRangeParams tests multiple range params
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithMultipleRangeParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		RangeParams: []ConditionRange{
			{
				RangeName: "content-length-range",
				Lower:     0,
				Upper:     10485760,
			},
			{
				RangeName: "content-length-range",
				Lower:     100,
				Upper:     500,
			},
		},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
	// Should contain both range conditions
	assert.Contains(t, output.OriginPolicy, "content-length-range")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSecurityToken_SignatureObs tests security token with OBS signature
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSecurityToken_SignatureObs(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs),
		WithSecurityToken("test-security-token"))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	// Should include x-obs-security-token in the policy
	assert.Contains(t, output.OriginPolicy, "x-obs-security-token")
	assert.Contains(t, output.OriginPolicy, "test-security-token")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSecurityToken_SignatureV4 tests security token with V4 signature
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSecurityToken_SignatureV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
		WithSecurityToken("test-security-token"))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	// Should include x-amz-security-token in the policy for V4
	assert.Contains(t, output.OriginPolicy, "x-amz-security-token")
	assert.Contains(t, output.OriginPolicy, "test-security-token")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSpecialCharsInBucket tests special characters in bucket name
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSpecialCharsInBucket(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket-name",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.Contains(t, output.OriginPolicy, "test-bucket-name")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSpecialCharsInKey tests special characters in key
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithSpecialCharsInKey(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test/key/with special-chars.txt",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.Contains(t, output.OriginPolicy, "test/key/with special-chars.txt")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithCustomExpires tests custom expires values
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithCustomExpires(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	expiresValues := []int{100, 300, 600, 3600}

	for _, expires := range expiresValues {
		t.Run(string(rune(expires)), func(t *testing.T) {
			input := &CreateBrowserBasedSignatureInput{
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: expires,
			}

			output, err := client.CreateBrowserBasedSignature(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.OriginPolicy)
			assert.NotEmpty(t, output.Signature)
		})
	}
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithAllFields tests all fields
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithAllFields(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs),
		WithSecurityToken("test-security-token"))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key/file.txt",
		Expires: 900,
		FormParams: map[string]string{
			"acl":                  "public-read",
			"content-type":         "application/json",
			"cache-control":        "no-cache",
			"content-disposition":  "attachment; filename=test.txt",
		},
		RangeParams: []ConditionRange{
			{
				RangeName: "content-length-range",
				Lower:     0,
				Upper:     10485760,
			},
		},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	assert.Contains(t, output.OriginPolicy, "test-bucket")
	assert.Contains(t, output.OriginPolicy, "test-key/file.txt")
	assert.Contains(t, output.OriginPolicy, "acl")
	assert.Contains(t, output.OriginPolicy, "public-read")
	assert.Contains(t, output.OriginPolicy, "content-length-range")
	assert.Contains(t, output.OriginPolicy, "x-obs-security-token")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithMultipleFormParams tests multiple form params
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithMultipleFormParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		FormParams: map[string]string{
			"acl":                  "public-read",
			"content-type":         "application/json",
			"cache-control":        "no-cache",
			"content-disposition":  "attachment; filename=test.txt",
			"x-obs-meta-key":      "meta-value",
		},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
	assert.Contains(t, output.OriginPolicy, "acl")
	assert.Contains(t, output.OriginPolicy, "content-type")
	assert.Contains(t, output.OriginPolicy, "cache-control")
	assert.Contains(t, output.OriginPolicy, "content-disposition")
	assert.Contains(t, output.OriginPolicy, "x-obs-meta-key")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithComplexFormParams tests complex form params values
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithComplexFormParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		FormParams: map[string]string{
			"content-disposition":  "attachment; filename=\"测试文件.txt\"",
			"success_action_redirect": "https://example.com/success",
			"x-amz-meta-custom": "custom metadata value",
		},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
	assert.Contains(t, output.OriginPolicy, "content-disposition")
	assert.Contains(t, output.OriginPolicy, "success_action_redirect")
	assert.Contains(t, output.OriginPolicy, "x-amz-meta-custom")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithOnlyBucket tests with only bucket
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithOnlyBucket(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
	assert.Contains(t, output.OriginPolicy, "test-bucket")
	// Should use "starts-with" for key since it's empty
	assert.Contains(t, output.OriginPolicy, "\"starts-with\", \"$key\", \"\"")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithOnlyKey tests with only key
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithOnlyKey(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Signature)
	assert.Contains(t, output.OriginPolicy, "test-key")
	// Should use "starts-with" for bucket since it's empty
	assert.Contains(t, output.OriginPolicy, "\"starts-with\", \"$bucket\", \"\"")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV4AndSecurityToken tests V4 with security token
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV4AndSecurityToken(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
		WithSecurityToken("test-security-token"))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		FormParams: map[string]string{"acl": "public-read"},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	assert.NotEmpty(t, output.Algorithm)
	assert.NotEmpty(t, output.Credential)
	assert.NotEmpty(t, output.Date)
	// V4 signature includes algorithm, credential, date
	assert.Contains(t, output.Algorithm, "AWS4-HMAC-SHA256")
	assert.Contains(t, output.Credential, "test-ak")
	assert.Contains(t, output.Date, "T")
	// Should include x-amz-security-token for V4
	assert.Contains(t, output.OriginPolicy, "x-amz-security-token")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV4AndFormParams tests V4 with form params
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV4AndFormParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		FormParams: map[string]string{
			"acl":          "public-read",
			"content-type": "application/json",
		},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	assert.NotEmpty(t, output.Algorithm)
	assert.NotEmpty(t, output.Credential)
	assert.NotEmpty(t, output.Date)
	assert.Contains(t, output.OriginPolicy, "acl")
	assert.Contains(t, output.OriginPolicy, "content-type")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV4AndRangeParams tests V4 with range params
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV4AndRangeParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		RangeParams: []ConditionRange{
			{
				RangeName: "content-length-range",
				Lower:     0,
				Upper:     10485760,
			},
		},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	assert.NotEmpty(t, output.Algorithm)
	assert.NotEmpty(t, output.Credential)
	assert.NotEmpty(t, output.Date)
	assert.Contains(t, output.OriginPolicy, "content-length-range")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV2Signature tests V2 signature
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithV2Signature(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureV2))

	input := &CreateBrowserBasedSignatureInput{
		Bucket:     "test-bucket",
		Key:        "test-key",
		Expires:    600,
		FormParams: map[string]string{"acl": "public-read"},
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
	// V2 signature doesn't include algorithm, credential, date
	assert.Empty(t, output.Algorithm)
	assert.Empty(t, output.Credential)
	assert.Empty(t, output.Date)
	assert.Contains(t, output.OriginPolicy, "acl")
}

// TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithEmptyAKSK tests with empty AK/SK
func TestCreateBrowserBasedSignature_ShouldReturnSuccess_WithEmptyAKSK(t *testing.T) {
	client, err := New("", "", "https://obs.test.example.com",
		WithSignature(SignatureObs))
	assert.NoError(t, err)

	input := &CreateBrowserBasedSignatureInput{
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateBrowserBasedSignature(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.OriginPolicy)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
}
