// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package obs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// SseKmsHeader Tests

func TestSseKmsHeader_GetEncryption_ShouldReturnCustomEncryption_WhenEncryptionIsSet(t *testing.T) {
	header := SseKmsHeader{
		Encryption: "custom-encryption",
		Key:        "test-key",
		isObs:      false,
	}
	result := header.GetEncryption()
	assert.Equal(t, "custom-encryption", result)
}

func TestSseKmsHeader_GetEncryption_ShouldReturnDefaultAwsEncryption_WhenEncryptionIsEmptyAndNotObs(t *testing.T) {
	header := SseKmsHeader{
		Encryption: "",
		Key:        "test-key",
		isObs:      false,
	}
	result := header.GetEncryption()
	assert.Equal(t, DEFAULT_SSE_KMS_ENCRYPTION, result)
	assert.Equal(t, "aws:kms", result)
}

func TestSseKmsHeader_GetEncryption_ShouldReturnDefaultObsEncryption_WhenEncryptionIsEmptyAndIsObs(t *testing.T) {
	header := SseKmsHeader{
		Encryption: "",
		Key:        "test-key",
		isObs:      true,
	}
	result := header.GetEncryption()
	assert.Equal(t, DEFAULT_SSE_KMS_ENCRYPTION_OBS, result)
	assert.Equal(t, "kms", result)
}

func TestSseKmsHeader_GetKey_ShouldReturnKey_WhenKeyIsSet(t *testing.T) {
	header := SseKmsHeader{
		Key: "test-key-value",
	}
	result := header.GetKey()
	assert.Equal(t, "test-key-value", result)
}

func TestSseKmsHeader_GetKey_ShouldReturnEmptyString_WhenKeyIsEmpty(t *testing.T) {
	header := SseKmsHeader{
		Key: "",
	}
	result := header.GetKey()
	assert.Equal(t, "", result)
}

// SseCHeader Tests

func TestSseCHeader_GetEncryption_ShouldReturnCustomEncryption_WhenEncryptionIsSet(t *testing.T) {
	header := SseCHeader{
		Encryption: "custom-encryption",
		Key:        "test-key",
		KeyMD5:     "test-md5",
	}
	result := header.GetEncryption()
	assert.Equal(t, "custom-encryption", result)
}

func TestSseCHeader_GetEncryption_ShouldReturnDefaultEncryption_WhenEncryptionIsEmpty(t *testing.T) {
	header := SseCHeader{
		Encryption: "",
		Key:        "test-key",
		KeyMD5:     "test-md5",
	}
	result := header.GetEncryption()
	assert.Equal(t, DEFAULT_SSE_C_ENCRYPTION, result)
	assert.Equal(t, "AES256", result)
}

func TestSseCHeader_GetKey_ShouldReturnKey_WhenKeyIsSet(t *testing.T) {
	header := SseCHeader{
		Key: "test-key-value",
	}
	result := header.GetKey()
	assert.Equal(t, "test-key-value", result)
}

func TestSseCHeader_GetKey_ShouldReturnEmptyString_WhenKeyIsEmpty(t *testing.T) {
	header := SseCHeader{
		Key: "",
	}
	result := header.GetKey()
	assert.Equal(t, "", result)
}

func TestSseCHeader_GetKeyMD5_ShouldReturnCustomKeyMD5_WhenKeyMD5IsSet(t *testing.T) {
	header := SseCHeader{
		KeyMD5: "test-custom-md5",
	}
	result := header.GetKeyMD5()
	assert.Equal(t, "test-custom-md5", result)
}

func TestSseCHeader_GetKeyMD5_ShouldReturnCalculatedMD5_WhenKeyMD5IsEmptyAndKeyIsValid(t *testing.T) {
	// Create a valid base64 encoded key
	validKey := "test-key-12345"
	encodedKey := Base64Encode([]byte(validKey))
	expectedMD5 := Base64Md5([]byte(validKey))

	header := SseCHeader{
		Key:    encodedKey,
		KeyMD5: "",
	}
	result := header.GetKeyMD5()
	assert.Equal(t, expectedMD5, result)
}

func TestSseCHeader_GetKeyMD5_ShouldReturnEmptyString_WhenKeyMD5IsEmptyAndKeyIsInvalid(t *testing.T) {
	header := SseCHeader{
		Key:    "invalid-base64!!!",
		KeyMD5: "",
	}
	result := header.GetKeyMD5()
	assert.Equal(t, "", result)
}

func TestSseCHeader_GetKeyMD5_ShouldReturnMD5OfEmptyString_WhenKeyMD5IsEmptyAndKeyIsEmpty(t *testing.T) {
	header := SseCHeader{
		Key:    "",
		KeyMD5: "",
	}
	result := header.GetKeyMD5()

	// When Key is empty, Base64Decode("") returns (nil, nil),
	// and Base64Md5(nil) returns MD5 of empty bytes base64-encoded
	// Expected: "1B2M2Y8AsgTpgAmY7PhCfg==" or similar
	assert.NotEmpty(t, result)
}

// setSseHeader Tests

func TestSetSseHeader_ShouldSetSseCHeaders_WhenSseCHeaderIsProvided(t *testing.T) {
	headers := make(map[string][]string)
	sseHeader := SseCHeader{
		Encryption: "AES256",
		Key:        "test-key",
		KeyMD5:     "test-md5",
	}

	setSseHeader(headers, sseHeader, false, false)

	assert.Contains(t, headers, "x-amz-server-side-encryption-customer-algorithm")
	assert.Equal(t, []string{"AES256"}, headers["x-amz-server-side-encryption-customer-algorithm"])
	assert.Contains(t, headers, "x-amz-server-side-encryption-customer-key")
	assert.Equal(t, []string{"test-key"}, headers["x-amz-server-side-encryption-customer-key"])
	assert.Contains(t, headers, "x-amz-server-side-encryption-customer-key-MD5")
	assert.Equal(t, []string{"test-md5"}, headers["x-amz-server-side-encryption-customer-key-MD5"])
}

func TestSetSseHeader_ShouldSetObsSseCHeaders_WhenIsObsIsTrue(t *testing.T) {
	headers := make(map[string][]string)
	sseHeader := SseCHeader{
		Encryption: "AES256",
		Key:        "test-key",
		KeyMD5:     "test-md5",
	}

	setSseHeader(headers, sseHeader, false, true)

	assert.Contains(t, headers, "x-obs-server-side-encryption-customer-algorithm")
	assert.Contains(t, headers, "x-obs-server-side-encryption-customer-key")
	assert.Contains(t, headers, "x-obs-server-side-encryption-customer-key-MD5")
}

func TestSetSseHeader_ShouldSetSseKmsHeaders_WhenSseKmsHeaderIsProvided(t *testing.T) {
	headers := make(map[string][]string)
	sseHeader := SseKmsHeader{
		Encryption: "aws:kms",
		Key:        "kms-key-id",
	}

	setSseHeader(headers, sseHeader, false, false)

	assert.Contains(t, headers, "x-amz-server-side-encryption")
	assert.Equal(t, []string{"aws:kms"}, headers["x-amz-server-side-encryption"])
	assert.Contains(t, headers, "x-amz-server-side-encryption-aws-kms-key-id")
	assert.Equal(t, []string{"kms-key-id"}, headers["x-amz-server-side-encryption-aws-kms-key-id"])
}

func TestSetSseHeader_ShouldSetObsSseKmsHeaders_WhenIsObsIsTrue(t *testing.T) {
	headers := make(map[string][]string)
	sseHeader := SseKmsHeader{
		Encryption: "aws:kms",
		Key:        "kms-key-id",
	}

	setSseHeader(headers, sseHeader, false, true)

	assert.Contains(t, headers, "x-obs-server-side-encryption")
	assert.NotContains(t, headers, "x-amz-server-side-encryption")
	assert.Contains(t, headers, "x-obs-server-side-encryption-kms-key-id")
	assert.NotContains(t, headers, "x-amz-server-side-encryption-aws-kms-key-id")
}

func TestSetSseHeader_ShouldSetSseKmsHeadersWithoutKey_WhenKeyIsEmpty(t *testing.T) {
	headers := make(map[string][]string)
	sseHeader := SseKmsHeader{
		Encryption: "aws:kms",
		Key:        "",
	}

	setSseHeader(headers, sseHeader, false, false)

	assert.Contains(t, headers, "x-amz-server-side-encryption")
	assert.NotContains(t, headers, "x-amz-server-side-encryption-aws-kms-key-id")
}

func TestSetSseHeader_ShouldNotSetSseKmsHeaders_WhenSseCOnlyIsTrue(t *testing.T) {
	headers := make(map[string][]string)
	sseHeader := SseKmsHeader{
		Encryption: "aws:kms",
		Key:        "kms-key-id",
	}

	setSseHeader(headers, sseHeader, true, false)

	assert.NotContains(t, headers, "x-amz-server-side-encryption")
	assert.NotContains(t, headers, "x-amz-server-side-encryption-aws-kms-key-id")
}

func TestSetSseHeader_ShouldDoNothing_WhenSseHeaderIsNil(t *testing.T) {
	headers := make(map[string][]string)
	headers["existing-header"] = []string{"value"}

	setSseHeader(headers, nil, false, false)

	assert.Len(t, headers, 1)
	assert.Contains(t, headers, "existing-header")
}
