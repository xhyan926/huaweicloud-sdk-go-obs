// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of
// License at
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

// StringToInt Tests

func TestStringToInt_ShouldReturnCorrectValue_WhenGivenValidString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      int
		expected int
	}{
		{"Valid positive", "123", 0, 123},
		{"Valid negative", "-456", 0, -456},
		{"Valid zero", "0", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToInt(tt.input, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringToInt_ShouldReturnDefaultValue_WhenGivenInvalidString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      int
		expected int
	}{
		{"Invalid string", "abc", 100, 100},
		{"Empty string", "", 999, 999},
		{"Mixed string", "123abc", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToInt(tt.input, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// StringToInt64 Tests

func TestStringToInt64_ShouldReturnCorrectValue_WhenGivenValidString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      int64
		expected int64
	}{
		{"Valid positive", "123456789", 0, 123456789},
		{"Valid negative", "-9876543210", 0, -9876543210},
		{"Valid zero", "0", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToInt64(tt.input, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringToInt64_ShouldReturnDefaultValue_WhenGivenInvalidString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      int64
		expected int64
	}{
		{"Invalid string", "abc", 999, 999},
		{"Empty string", "", -1, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToInt64(tt.input, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// IntToString Tests

func TestIntToString_ShouldReturnCorrectString_WhenGivenValidInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Zero", 0, "0"},
		{"Positive", 123, "123"},
		{"Negative", -456, "-456"},
		{"Large number", 2147483647, "2147483647"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Int64ToString Tests

func TestInt64ToString_ShouldReturnCorrectString_WhenGivenValidInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Zero", 0, "0"},
		{"Positive", 123456789, "123456789"},
		{"Negative", -9876543210, "-9876543210"},
		{"Large number", 9223372036854775807, "9223372036854775807"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int64ToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// XmlTranscoding Tests

func TestXmlTranscoding_ShouldEscapeSpecialCharacters_WhenGivenValidInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Less than", "<test>", "&lt;test&gt;"},
		{"Ampersand", "a&b", "a&amp;b"},
		{"Apostrophe", "it's", "it&apos;s"},
		{"Quote", "test\"value", "test&quot;value"},
		{"All special", "<>&'\"", "&lt;&gt;&amp;&apos;&quot;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XmlTranscoding(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestXmlTranscoding_ShouldReturnSameString_WhenNoSpecialCharacters(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Normal text", "Hello World"},
		{"Numbers", "123456"},
		{"Mixed", "test123.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XmlTranscoding(tt.input)
			assert.Equal(t, tt.input, result)
		})
	}
}

// Md5 Tests

func TestMd5_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	input := []byte("test data")
	result := Md5(input)
	assert.NotNil(t, result)
	assert.Equal(t, 16, len(result))
	expected := "eb733a00c0c9d336e65691a37ab54293"
	actual := Hex(result)
	assert.Equal(t, expected, actual)
}

func TestMd5_ShouldReturnSameHash_WhenGivenSameInput(t *testing.T) {
	input := []byte("same input")
	result1 := Md5(input)
	result2 := Md5(input)
	assert.Equal(t, result1, result2)
}

// HmacSha1 Tests

func TestHmacSha1_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	key := []byte("test-key")
	value := []byte("test-value")
	result := HmacSha1(key, value)
	assert.NotNil(t, result)
	assert.Equal(t, 20, len(result))
}

func TestHmacSha1_ShouldReturnDifferentHash_WhenGivenDifferentKey(t *testing.T) {
	value := []byte("test-value")
	result1 := HmacSha1([]byte("key1"), value)
	result2 := HmacSha1([]byte("key2"), value)
	assert.NotEqual(t, result1, result2)
}

// HmacSha256 Tests

func TestHmacSha256_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	key := []byte("test-key")
	value := []byte("test-value")
	result := HmacSha256(key, value)
	assert.NotNil(t, result)
	assert.Equal(t, 32, len(result))
}

func TestHmacSha256_ShouldReturnSameHash_WhenGivenSameInput(t *testing.T) {
	key := []byte("test-key")
	value := []byte("test-value")
	result1 := HmacSha256(key, value)
	result2 := HmacSha256(key, value)
	assert.Equal(t, result1, result2)
}

// Base64Encode Tests

func TestBase64Encode_ShouldEncodeCorrectly_WhenGivenValidInput(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"Simple text", []byte("test"), "dGVzdA=="},
		{"Hello World", []byte("Hello World"), "SGVsbG8gV29ybGQ="},
		{"Empty", []byte(""), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Encode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Base64Decode Tests

func TestBase64Decode_ShouldDecodeCorrectly_WhenGivenValidInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		wantErr  bool
	}{
		{"Simple text", "dGVzdA==", []byte("test"), false},
		{"Hello World", "SGVsbG8gV29ybGQ=", []byte("Hello World"), false},
		{"Invalid", "invalid!!", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Base64Decode(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Sha256Hash Tests

func TestSha256Hash_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	input := []byte("test data")
	result := Sha256Hash(input)
	assert.NotNil(t, result)
	assert.Equal(t, 32, len(result))
}

func TestSha256Hash_ShouldReturnSameHash_WhenGivenSameInput(t *testing.T) {
	input := []byte("same input")
	result1 := Sha256Hash(input)
	result2 := Sha256Hash(input)
	assert.Equal(t, result1, result2)
}

// Base64Md5 Tests

func TestBase64Md5_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	input := []byte("test data")
	result := Base64Md5(input)
	assert.NotEmpty(t, result)
	assert.Equal(t, 24, len(result))
}

// Base64Sha256 Tests

func TestBase64Sha256_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	input := []byte("test data")
	result := Base64Sha256(input)
	assert.NotEmpty(t, result)
	assert.Equal(t, 44, len(result))
}

// Base64Md5OrSha256 Tests

func TestBase64Md5OrSha256_ShouldReturnMd5_WhenSha256IsFalse(t *testing.T) {
	input := []byte("test data")
	result := Base64Md5OrSha256(input, false)
	assert.NotEmpty(t, result)
	assert.Equal(t, 24, len(result))
}

func TestBase64Md5OrSha256_ShouldReturnSha256_WhenSha256IsTrue(t *testing.T) {
	input := []byte("test data")
	result := Base64Md5OrSha256(input, true)
	assert.NotEmpty(t, result)
	assert.Equal(t, 44, len(result))
}

// Hex Tests

func TestHex_ShouldEncodeCorrectly_WhenGivenValidInput(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"Simple", []byte("test"), "74657374"},
		{"Hello", []byte("Hello"), "48656c6c6f"},
		{"Empty", []byte(""), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Hex(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// HexMd5 Tests

func TestHexMd5_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	input := []byte("test data")
	result := HexMd5(input)
	assert.NotEmpty(t, result)
}

// HexSha256 Tests

func TestHexSha256_ShouldReturnCorrectHash_WhenGivenValidInput(t *testing.T) {
	input := []byte("test data")
	result := HexSha256(input)
	assert.NotEmpty(t, result)
	assert.Equal(t, 64, len(result))
}

// UrlEncode Tests

func TestUrlEncode_ShouldEncodeAllCharacters_WhenChineseOnlyIsFalse(t *testing.T) {
	input := "test file.txt"
	result := UrlEncode(input, false)
	// url.QueryEscape encodes space as '+', not '%20'
	assert.Contains(t, result, "+")
}

func TestUrlEncode_ShouldEncodeOnlyChinese_WhenChineseOnlyIsTrue(t *testing.T) {
	input := "test测试.txt"
	result := UrlEncode(input, true)
	assert.Contains(t, result, "%")
	assert.Contains(t, result, "test")
	assert.Contains(t, result, ".txt")
}

func TestUrlEncode_ShouldReturnSame_WhenNoSpecialChars(t *testing.T) {
	input := "test.txt"
	result := UrlEncode(input, false)
	assert.NotContains(t, result, "%")
}

// UrlDecode Tests

func TestUrlDecode_ShouldDecodeCorrectly_WhenGivenValidInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"Space", "test%20file", "test file", false},
		{"Plus", "test+file", "test file", false},
		{"Percent", "test%2Ffile", "test/file", false},
		{"Invalid", "%ZZ", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UrlDecode(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUrlDecodeWithoutError_ShouldDecodeCorrectly_WhenGivenValidInput(t *testing.T) {
	result := UrlDecodeWithoutError("test%20file")
	assert.Equal(t, "test file", result)
}

// IsIP Tests

func TestIsIP_ShouldReturnTrue_WhenGivenValidIP(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"IPv4 localhost", "127.0.0.1"},
		{"IPv4 private", "192.168.1.1"},
		{"IPv4 private 2", "10.0.0.1"},
		{"IPv4 max", "255.255.255.255"},
		{"IPv4 min", "0.0.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIP(tt.input)
			assert.True(t, result)
		})
	}
}

func TestIsIP_ShouldReturnFalse_WhenGivenInvalidIP(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty", ""},
		{"Invalid", "invalid"},
		{"Too high", "256.0.0.1"},
		{"Negative", "192.168.-1.1"},
		{"Too few parts", "192.168.1"},
		{"Too many parts", "192.168.1.1.1"},
		{"Domain", "example.com"},
		{"Hostname", "test-server"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIP(tt.input)
			assert.False(t, result)
		})
	}
}

// IsContain Tests

func TestIsContain_ShouldReturnTrue_WhenItemExists(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	result := IsContain(items, "banana")
	assert.True(t, result)
}

func TestIsContain_ShouldReturnFalse_WhenItemNotExists(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	result := IsContain(items, "orange")
	assert.False(t, result)
}

func TestIsContain_ShouldReturnFalse_WhenItemsIsEmpty(t *testing.T) {
	items := []string{}
	result := IsContain(items, "test")
	assert.False(t, result)
}

// StringContains Tests

func TestStringContains_ShouldReplaceAllOccurrences(t *testing.T) {
	result := StringContains("hello world", "o", "O")
	assert.Equal(t, "hellO wOrld", result)
}

func TestStringContains_ShouldNotReplace_WhenSubstrNotExists(t *testing.T) {
	result := StringContains("hello world", "z", "Z")
	assert.Equal(t, "hello world", result)
}

// SignatureType Constants Tests

func TestSignatureTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "v2", string(SignatureV2))
	assert.Equal(t, "v4", string(SignatureV4))
	assert.Equal(t, "OBS", string(SignatureObs))
}

// HttpMethodType Constants Tests

func TestHttpMethodTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, string(HttpMethodGet), HTTP_GET)
	assert.Equal(t, string(HttpMethodPut), HTTP_PUT)
	assert.Equal(t, string(HttpMethodPost), HTTP_POST)
	assert.Equal(t, string(HttpMethodDelete), HTTP_DELETE)
	assert.Equal(t, string(HttpMethodHead), HTTP_HEAD)
	assert.Equal(t, string(HttpMethodOptions), HTTP_OPTIONS)
}

// StorageClassType Constants Tests

func TestStorageClassTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "STANDARD", string(StorageClassStandard))
	assert.Equal(t, "WARM", string(StorageClassWarm))
	assert.Equal(t, "COLD", string(StorageClassCold))
	assert.Equal(t, "DEEP_ARCHIVE", string(StorageClassDeepArchive))
	assert.Equal(t, "INTELLIGENT_TIERING", string(StorageClassIntelligentTiering))
}

// AclType Constants Tests

func TestAclTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "private", string(AclPrivate))
	assert.Equal(t, "public-read", string(AclPublicRead))
	assert.Equal(t, "public-read-write", string(AclPublicReadWrite))
	assert.Equal(t, "authenticated-read", string(AclAuthenticatedRead))
	assert.Equal(t, "bucket-owner-read", string(AclBucketOwnerRead))
	assert.Equal(t, "bucket-owner-full-control", string(AclBucketOwnerFullControl))
}

// PermissionType Constants Tests

func TestPermissionTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "READ", string(PermissionRead))
	assert.Equal(t, "WRITE", string(PermissionWrite))
	assert.Equal(t, "READ_ACP", string(PermissionReadAcp))
	assert.Equal(t, "WRITE_ACP", string(PermissionWriteAcp))
	assert.Equal(t, "FULL_CONTROL", string(PermissionFullControl))
}

// VersioningStatusType Constants Tests

func TestVersioningStatusTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "Enabled", string(VersioningStatusEnabled))
	assert.Equal(t, "Suspended", string(VersioningStatusSuspended))
}

// RuleStatusType Constants Tests

func TestRuleStatusTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "Enabled", string(RuleStatusEnabled))
	assert.Equal(t, "Disabled", string(RuleStatusDisabled))
}

// RestoreTierType Constants Tests

func TestRestoreTierTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "Expedited", string(RestoreTierExpedited))
	assert.Equal(t, "Standard", string(RestoreTierStandard))
	assert.Equal(t, "Bulk", string(RestoreTierBulk))
}

// MetadataDirectiveType Constants Tests

func TestMetadataDirectiveTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "COPY", string(CopyMetadata))
	assert.Equal(t, "REPLACE_NEW", string(ReplaceNew))
	assert.Equal(t, "REPLACE", string(ReplaceMetadata))
}

// PayerType Constants Tests

func TestPayerTypeConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "BucketOwner", string(BucketOwnerPayer))
	assert.Equal(t, "Requester", string(RequesterPayer))
	assert.Equal(t, "requester", string(Requester))
}

// OBS SDK Version Test

func TestObsSdkVersion_ShouldHaveCorrectFormat(t *testing.T) {
	assert.NotEmpty(t, OBS_SDK_VERSION)
	assert.Contains(t, OBS_SDK_VERSION, ".")
}

// User Agent Test

func TestUserAgent_ShouldHaveCorrectFormat(t *testing.T) {
	assert.Contains(t, USER_AGENT, "obs-sdk-go/")
	assert.Contains(t, USER_AGENT, OBS_SDK_VERSION)
}

// ObsError Tests

// func TestObsError_String_ShouldContainStatusAndCode(t *testing.T) {
// 	err := &ObsError{
// 		Status:    "403",
// 		Code:      "AccessDenied",
// 		Message:   "Access Denied",
// 	}
// 	result := err.Error()
// 	assert.Contains(t, result, "Status=403")
// 	assert.Contains(t, result, "Code=AccessDenied")
// }

// Default Constants Tests

func TestDefaultConstants_ShouldHaveReasonableValues(t *testing.T) {
	assert.Equal(t, SignatureV2, DEFAULT_SIGNATURE)
	assert.Equal(t, "region", DEFAULT_REGION)
	assert.Equal(t, 60, DEFAULT_CONNECT_TIMEOUT)
	assert.Equal(t, 60, DEFAULT_SOCKET_TIMEOUT)
	assert.Equal(t, 60, DEFAULT_HEADER_TIMEOUT)
	assert.Equal(t, 30, DEFAULT_IDLE_CONN_TIMEOUT)
	assert.Equal(t, 3, DEFAULT_MAX_RETRY_COUNT)
	assert.Equal(t, 3, DEFAULT_MAX_REDIRECT_COUNT)
	assert.Equal(t, 1000, DEFAULT_MAX_CONN_PER_HOST)
}

// Part Size Constants Tests

func TestPartSizeConstants_ShouldHaveReasonableValues(t *testing.T) {
	assert.True(t, MAX_PART_SIZE > MIN_PART_SIZE)
	assert.True(t, DEFAULT_PART_SIZE >= MIN_PART_SIZE)
	assert.True(t, DEFAULT_PART_SIZE <= MAX_PART_SIZE)
	assert.True(t, MAX_PART_NUM > 0)
}

// Format Constants Tests

func TestFormatConstants_ShouldHaveCorrectFormats(t *testing.T) {
	assert.Equal(t, "20060102T150405Z", LONG_DATE_FORMAT)
	assert.Equal(t, "20060102", SHORT_DATE_FORMAT)
	assert.Contains(t, ISO8601_DATE_FORMAT, "T")
	assert.Contains(t, ISO8601_MIDNIGHT_DATE_FORMAT, "T00:00:00Z")
	assert.Contains(t, RFC1123_FORMAT, "Mon")
	assert.Contains(t, RFC1123_FORMAT, "Jan")
}

// V4 Hash Prefix Test

func TestV4HashPrefix_ShouldHaveCorrectValue(t *testing.T) {
	assert.Equal(t, "AWS4-HMAC-SHA256", V4_HASH_PREFIX)
	assert.Equal(t, "AWS4", V4_HASH_PRE)
}

// V2 Hash Prefix Test

func TestV2HashPrefix_ShouldHaveCorrectValue(t *testing.T) {
	assert.Equal(t, "AWS", V2_HASH_PREFIX)
	assert.Equal(t, "OBS", OBS_HASH_PREFIX)
}

// Service Constants Test

func TestServiceConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "s3", V4_SERVICE_NAME)
	assert.Equal(t, "aws4_request", V4_SERVICE_SUFFIX)
}

// Encryption Constants Tests

func TestEncryptionConstants_ShouldHaveCorrectValues(t *testing.T) {
	assert.Equal(t, "aws:kms", DEFAULT_SSE_KMS_ENCRYPTION)
	assert.Equal(t, "kms", DEFAULT_SSE_KMS_ENCRYPTION_OBS)
	assert.Equal(t, "AES256", DEFAULT_SSE_C_ENCRYPTION)
}
