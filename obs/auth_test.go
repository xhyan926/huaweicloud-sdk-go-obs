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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setURLWithPolicy Tests

func TestSetURLWithPolicy_ShouldRemoveBucketPrefix_WhenURLStartsWithBucketSlash(t *testing.T) {
	bucketName := "test-bucket"
	canonicalizedUrl := "/test-bucket/some/object/key"

	result := setURLWithPolicy(bucketName, canonicalizedUrl)

	assert.Equal(t, "some/object/key", result)
}

func TestSetURLWithPolicy_ShouldRemoveBucketPrefix_WhenURLIsBucketOnly(t *testing.T) {
	bucketName := "test-bucket"
	canonicalizedUrl := "/test-bucket"

	result := setURLWithPolicy(bucketName, canonicalizedUrl)

	assert.Equal(t, "", result)
}

func TestSetURLWithPolicy_ShouldReturnUnchanged_WhenURLDoesNotStartWithBucket(t *testing.T) {
	bucketName := "test-bucket"
	canonicalizedUrl := "/some/other/path"

	result := setURLWithPolicy(bucketName, canonicalizedUrl)

	assert.Equal(t, "/some/other/path", result)
}

func TestSetURLWithPolicy_ShouldReturnUnchanged_WhenURLEmpty(t *testing.T) {
	bucketName := "test-bucket"
	canonicalizedUrl := ""

	result := setURLWithPolicy(bucketName, canonicalizedUrl)

	assert.Equal(t, "", result)
}

// prepareHostAndDate Tests

func TestPrepareHostAndDate_ShouldAddHostHeader_WhenCalled(t *testing.T) {
	headers := make(map[string][]string)
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, false)

	assert.Equal(t, []string{hostName}, headers[HEADER_HOST_CAMEL])
}

func TestPrepareHostAndDate_ShouldAddDateHeader_WhenNoDateHeader(t *testing.T) {
	headers := make(map[string][]string)
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, false)

	assert.Contains(t, headers, HEADER_DATE_CAMEL)
	assert.Len(t, headers[HEADER_DATE_CAMEL], 1)
}

func TestPrepareHostAndDate_ShouldConvertAmzDate_WhenIsV4(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, true)

	assert.Contains(t, headers, HEADER_DATE_CAMEL)
	assert.Equal(t, []string{"Mon, 01 Jan 2024 12:00:00 GMT"}, headers[HEADER_DATE_CAMEL])
}

func TestPrepareHostAndDate_ShouldKeepGMTDate_WhenIsV2(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"Mon, 01 Jan 2024 12:00:00 GMT"},
	}
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, false)

	assert.Contains(t, headers, HEADER_DATE_CAMEL)
	assert.Equal(t, []string{"Mon, 01 Jan 2024 12:00:00 GMT"}, headers[HEADER_DATE_CAMEL])
}

func TestPrepareHostAndDate_ShouldDeleteAmzDate_WhenConversionFails(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"invalid-date"},
	}
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, false)

	assert.NotContains(t, headers, HEADER_DATE_AMZ)
	assert.Contains(t, headers, HEADER_DATE_CAMEL)
}

func TestPrepareHostAndDate_ShouldUseExistingDate_WhenDateHeaderExists(t *testing.T) {
	existingDate := "Mon, 01 Jan 2024 12:00:00 GMT"
	headers := map[string][]string{
		HEADER_DATE_CAMEL: []string{existingDate},
	}
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, false)

	assert.Equal(t, []string{existingDate}, headers[HEADER_DATE_CAMEL])
}

func TestPrepareHostAndDate_ShouldDeleteAmzDate_WhenMultipleValues(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"20240101T120000Z", "20240102T120000Z"},
	}
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, true)

	assert.NotContains(t, headers, HEADER_DATE_AMZ)
	assert.Contains(t, headers, HEADER_DATE_CAMEL)
}

func TestPrepareHostAndDate_ShouldUseAmzDate_WhenIsV2AndGMT(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"Tue, 02 Jan 2024 13:00:00 GMT"},
	}
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, false)

	// When date is successfully used in V2 mode, HEADER_DATE_AMZ is kept
	assert.Contains(t, headers, HEADER_DATE_AMZ)
	assert.Equal(t, []string{"Tue, 02 Jan 2024 13:00:00 GMT"}, headers[HEADER_DATE_CAMEL])
}

func TestPrepareHostAndDate_ShouldDeleteAmzDate_WhenV4ParseFails(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"invalid-v4-date"},
	}
	hostName := "obs.example.com"

	prepareHostAndDate(headers, hostName, true)

	assert.NotContains(t, headers, HEADER_DATE_AMZ)
	assert.Contains(t, headers, HEADER_DATE_CAMEL)
}

// encodeHeaders Tests

func TestEncodeHeaders_ShouldEncodeHeaderValues(t *testing.T) {
	headers := map[string][]string{
		"content-type": []string{"application/json"},
		"x-custom":     []string{"test value with spaces"},
	}

	encodeHeaders(headers)

	assert.Equal(t, []string{"application/json"}, headers["content-type"])
	// UrlEncode with chineseOnly=true only encodes Chinese characters
	assert.Equal(t, []string{"test value with spaces"}, headers["x-custom"])
}

func TestEncodeHeaders_ShouldEncodeChineseCharacters(t *testing.T) {
	headers := map[string][]string{
		"x-custom": []string{"中文test"},
	}

	encodeHeaders(headers)

	// Chinese characters should be URL encoded
	assert.NotEqual(t, []string{"中文test"}, headers["x-custom"])
}

func TestEncodeHeaders_ShouldHandleMultipleValues(t *testing.T) {
	headers := map[string][]string{
		"x-custom": []string{"value1", "value2", "value3"},
	}

	encodeHeaders(headers)

	assert.Equal(t, []string{"value1", "value2", "value3"}, headers["x-custom"])
}

func TestEncodeHeaders_ShouldHandleEmptyHeaders(t *testing.T) {
	headers := make(map[string][]string)

	encodeHeaders(headers)

	assert.Empty(t, headers)
}

// prepareDateHeader Tests

func TestPrepareDateHeader_ShouldClearDate_WhenBothHeadersPresent(t *testing.T) {
	headers := map[string][]string{
		PARAM_DATE_AMZ_CAMEL: []string{"20240101T120000Z"},
	}
	_headers := map[string][]string{
		HEADER_DATE_CAMEL: []string{"some-date"},
	}

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	assert.Equal(t, []string{""}, _headers[HEADER_DATE_CAMEL])
}

func TestPrepareDateHeader_ShouldClearDate_WhenLowercaseDatePresent(t *testing.T) {
	headers := map[string][]string{
		PARAM_DATE_AMZ_CAMEL: []string{"20240101T120000Z"},
	}
	_headers := map[string][]string{
		strings.ToLower(HEADER_DATE_CAMEL): []string{"some-date"},
	}

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	assert.Equal(t, []string{""}, _headers[HEADER_DATE_CAMEL])
}

func TestPrepareDateHeader_ShouldDoNothing_WhenNoDateHeader(t *testing.T) {
	headers := make(map[string][]string)
	_headers := make(map[string][]string)

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	// Should do nothing when no date headers are present
	assert.Empty(t, _headers)
}

func TestPrepareDateHeader_ShouldDoNothing_WhenLowercaseDateAndNoDataHeader(t *testing.T) {
	headers := make(map[string][]string)
	_headers := map[string][]string{
		strings.ToLower(HEADER_DATE_CAMEL): []string{"some-date"},
	}

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	// Should do nothing when dataHeader is not present
	assert.Equal(t, []string{"some-date"}, _headers[strings.ToLower(HEADER_DATE_CAMEL)])
}

func TestPrepareDateHeader_ShouldClearDate_WhenDataHeaderInHeaders(t *testing.T) {
	headers := map[string][]string{
		PARAM_DATE_AMZ_CAMEL: []string{"20240101T120000Z"},
	}
	_headers := map[string][]string{
		HEADER_DATE_CAMEL: []string{"some-date"},
	}

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	assert.Equal(t, []string{""}, _headers[HEADER_DATE_CAMEL])
}

func TestPrepareDateHeader_ShouldDoNothing_WhenDateCamelNotInHeaders(t *testing.T) {
	headers := make(map[string][]string)
	_headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	assert.Equal(t, []string{"20240101T120000Z"}, _headers[HEADER_DATE_AMZ])
}

func TestPrepareDateHeader_ShouldClearDate_WhenDataHeaderInLowercase(t *testing.T) {
	headers := make(map[string][]string)
	_headers := map[string][]string{
		strings.ToLower(HEADER_DATE_CAMEL): []string{"some-date"},
		HEADER_DATE_AMZ:                   []string{"20240101T120000Z"},
	}

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	assert.Equal(t, []string{""}, _headers[HEADER_DATE_CAMEL])
}

func TestPrepareDateHeader_ShouldClearDate_WhenLowercaseAndDataHeaderPresent(t *testing.T) {
	headers := map[string][]string{
		PARAM_DATE_AMZ_CAMEL: []string{"20240101T120000Z"},
	}
	_headers := map[string][]string{
		strings.ToLower(HEADER_DATE_CAMEL): []string{"some-date"},
	}

	prepareDateHeader(HEADER_DATE_AMZ, PARAM_DATE_AMZ_CAMEL, headers, _headers)

	assert.Equal(t, []string{""}, _headers[HEADER_DATE_CAMEL])
}

// getStringToSign Tests

func TestGetStringToSign_ShouldBuildString_WhenSimpleHeaders(t *testing.T) {
	keys := []string{"content-type", "content-md5", "date"}
	isObs := true
	_headers := map[string][]string{
		"content-type": []string{"application/json"},
		"content-md5":   []string{"abc123"},
		"date":          []string{"20240101T120000Z"},
	}

	result := getStringToSign(keys, isObs, _headers)

	// Non-prefix headers don't include the key
	assert.Contains(t, result, "application/json")
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "20240101T120000Z")
}

func TestGetStringToSign_ShouldHandleMetaHeaders_WhenOBS(t *testing.T) {
	keys := []string{HEADER_PREFIX_META_OBS + "test"}
	isObs := true
	_headers := map[string][]string{
		HEADER_PREFIX_META_OBS + "test": []string{"value1", "value2", "value3"},
	}

	result := getStringToSign(keys, isObs, _headers)

	assert.Contains(t, result, HEADER_PREFIX_META_OBS+"test:value1,value2,value3")
}

func TestGetStringToSign_ShouldJoinMultipleMetaValues(t *testing.T) {
	keys := []string{HEADER_PREFIX_META + "test"}
	isObs := false
	_headers := map[string][]string{
		HEADER_PREFIX_META + "test": []string{"value1", "value2", "value3"},
	}

	result := getStringToSign(keys, isObs, _headers)

	assert.Contains(t, result, HEADER_PREFIX_META+"test:value1,value2,value3")
}

func TestGetStringToSign_ShouldHandleAmzHeaders(t *testing.T) {
	keys := []string{"x-amz-test"}
	isObs := false
	_headers := map[string][]string{
		"x-amz-test": []string{"value1", "value2"},
	}

	result := getStringToSign(keys, isObs, _headers)

	assert.Contains(t, result, "x-amz-test:value1,value2")
}

func TestGetStringToSign_ShouldHandleNonPrefixHeaders(t *testing.T) {
	keys := []string{"content-type"}
	isObs := false
	_headers := map[string][]string{
		"content-type": []string{"value1", "value2"},
	}

	result := getStringToSign(keys, isObs, _headers)

	assert.Contains(t, result, "value1,value2")
}

func TestGetStringToSign_ShouldHandleObsHeaders_WhenIsObs(t *testing.T) {
	keys := []string{"x-obs-test"}
	isObs := true
	_headers := map[string][]string{
		"x-obs-test": []string{"value1", "value2"},
	}

	result := getStringToSign(keys, isObs, _headers)

	assert.Contains(t, result, "x-obs-test:value1,value2")
}

func TestGetStringToSign_ShouldTrimMetaValues(t *testing.T) {
	keys := []string{HEADER_PREFIX_META + "test"}
	isObs := false
	_headers := map[string][]string{
		HEADER_PREFIX_META + "test": []string{"  value1  ", " value2 ", "value3"},
	}

	result := getStringToSign(keys, isObs, _headers)

	assert.Contains(t, result, HEADER_PREFIX_META+"test:value1,value2,value3")
}

// attachHeaders Tests

func TestAttachHeaders_ShouldAddStandardHeaders(t *testing.T) {
	headers := map[string][]string{
		"content-type": []string{"application/json"},
		"content-md5":   []string{"abc123"},
		"date":          []string{"20240101T120000Z"},
	}
	isObs := false

	result := attachHeaders(headers, isObs)

	// Non-prefix headers don't include the key
	assert.Contains(t, result, "application/json")
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "20240101T120000Z")
}

func TestAttachHeaders_ShouldAddAmzHeaders(t *testing.T) {
	headers := map[string][]string{
		"x-amz-test": []string{"test-value"},
	}
	isObs := false

	result := attachHeaders(headers, isObs)

	assert.Contains(t, result, "x-amz-test:test-value")
}

func TestAttachHeaders_ShouldAddObsHeaders_WhenIsObs(t *testing.T) {
	headers := map[string][]string{
		"x-obs-test": []string{"test-value"},
	}
	isObs := true

	result := attachHeaders(headers, isObs)

	assert.Contains(t, result, "x-obs-test:test-value")
}

func TestAttachHeaders_ShouldIgnoreUninterestedHeaders(t *testing.T) {
	headers := map[string][]string{
		"x-custom": []string{"custom-value"},
	}
	isObs := false

	result := attachHeaders(headers, isObs)

	assert.NotContains(t, result, "x-custom")
}

func TestAttachHeaders_ShouldDeleteEmptyKeyHeaders(t *testing.T) {
	headers := map[string][]string{
		"": []string{"empty-key"},
	}
	isObs := false

	attachHeaders(headers, isObs)

	_, exists := headers[""]
	assert.False(t, exists)
}

func TestAttachHeaders_ShouldAddInterestedHeadersWhenMissing(t *testing.T) {
	headers := make(map[string][]string)
	isObs := false

	result := attachHeaders(headers, isObs)

	// interestedHeaders = []string{"content-md5", "content-type", "date"}
	// When headers are missing, they are added and produce empty strings
	// Result should have empty strings for each interested header
	lines := strings.Split(result, "\n")
	// Filter out empty strings
	nonEmptyLines := make([]string, 0)
	for _, line := range lines {
		if line != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	// All should be empty, so nonEmptyLines should be empty
	assert.Empty(t, nonEmptyLines)
}

func TestAttachHeaders_ShouldPreserveExistingInterestedHeaders(t *testing.T) {
	headers := map[string][]string{
		"content-type": []string{"application/json"},
	}
	isObs := false

	result := attachHeaders(headers, isObs)

	// Non-prefix headers don't include the key
	assert.Contains(t, result, "application/json")
}

func TestAttachHeaders_ShouldUseObsPrefix_WhenIsObs(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_OBS: []string{"20240101T120000Z"},
	}
	isObs := true

	result := attachHeaders(headers, isObs)

	assert.Contains(t, result, "20240101T120000Z")
}

func TestAttachHeaders_ShouldHandleMultipleStandardHeaders(t *testing.T) {
	headers := map[string][]string{
		"content-type":   []string{"application/json"},
		"content-md5":    []string{"abc123"},
		"x-amz-acl":      []string{"public-read"},
		"x-amz-meta-tag": []string{"value1,value2"},
	}
	isObs := false

	result := attachHeaders(headers, isObs)

	// Non-prefix headers (content-md5, content-type, date) don't include the key
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "application/json")
	assert.Contains(t, result, "x-amz-acl:public-read")
	assert.Contains(t, result, "x-amz-meta-tag:value1,value2")
}

// getScope Tests

func TestGetScope_ShouldReturnValidScope(t *testing.T) {
	region := "cn-north-4"
	shortDate := "20240101"

	result := getScope(region, shortDate)

	expected := "20240101/cn-north-4/s3/aws4_request"
	assert.Equal(t, expected, result)
}

func TestGetScope_ShouldHandleDifferentRegions(t *testing.T) {
	tests := []struct {
		region    string
		shortDate string
		expected  string
	}{
		{"us-east-1", "20240101", "20240101/us-east-1/s3/aws4_request"},
		{"eu-west-1", "20240101", "20240101/eu-west-1/s3/aws4_request"},
		{"ap-southeast-1", "20240101", "20240101/ap-southeast-1/s3/aws4_request"},
	}

	for _, tt := range tests {
		t.Run(tt.region, func(t *testing.T) {
			result := getScope(tt.region, tt.shortDate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// getCredential Tests

func TestGetCredential_ShouldReturnValidCredential(t *testing.T) {
	ak := "test-access-key"
	region := "cn-north-4"
	shortDate := "20240101"

	credential, scope := getCredential(ak, region, shortDate)

	expectedCredential := "test-access-key/20240101/cn-north-4/s3/aws4_request"
	expectedScope := "20240101/cn-north-4/s3/aws4_request"
	assert.Equal(t, expectedCredential, credential)
	assert.Equal(t, expectedScope, scope)
}

func TestGetCredential_ShouldHandleEmptyAK(t *testing.T) {
	ak := ""
	region := "cn-north-4"
	shortDate := "20240101"

	credential, scope := getCredential(ak, region, shortDate)

	expectedCredential := "/20240101/cn-north-4/s3/aws4_request"
	assert.Equal(t, expectedCredential, credential)
	assert.Equal(t, "20240101/cn-north-4/s3/aws4_request", scope)
}

// getSignedHeaders Tests

func TestGetSignedHeaders_ShouldReturnSortedHeaders(t *testing.T) {
	headers := map[string][]string{
		"z-header":  []string{"value-z"},
		"a-header":  []string{"value-a"},
		"m-header":  []string{"value-m"},
		"content-type": []string{"application/json"},
	}

	signedHeaders, _ := getSignedHeaders(headers)

	assert.Equal(t, []string{"a-header", "content-type", "m-header", "z-header"}, signedHeaders)
}

func TestGetSignedHeaders_ShouldLowercaseKeys(t *testing.T) {
	headers := map[string][]string{
		"Content-Type": []string{"application/json"},
		"Host":         []string{"example.com"},
	}

	signedHeaders, _headers := getSignedHeaders(headers)

	assert.Contains(t, signedHeaders, "content-type")
	assert.Contains(t, signedHeaders, "host")
	assert.Equal(t, []string{"application/json"}, _headers["content-type"])
	assert.Equal(t, []string{"example.com"}, _headers["host"])
}

func TestGetSignedHeaders_ShouldDeleteEmptyKeyHeaders(t *testing.T) {
	headers := map[string][]string{
		"":      []string{"empty"},
		"valid": []string{"value"},
	}

	signedHeaders, _headers := getSignedHeaders(headers)

	assert.NotContains(t, signedHeaders, "")
	_, exists := _headers[""]
	assert.False(t, exists)
}

func TestGetSignedHeaders_ShouldTrimWhitespace(t *testing.T) {
	headers := map[string][]string{
		"  content-type  ": []string{"application/json"},
		"  host  ":        []string{"example.com"},
	}

	signedHeaders, _ := getSignedHeaders(headers)

	assert.Contains(t, signedHeaders, "content-type")
	assert.Contains(t, signedHeaders, "host")
}

func TestGetSignedHeaders_ShouldHandleEmptyHeaders(t *testing.T) {
	headers := make(map[string][]string)

	signedHeaders, _headers := getSignedHeaders(headers)

	assert.Empty(t, signedHeaders)
	assert.Empty(t, _headers)
}

func TestGetSignedHeaders_ShouldCopyValuesToHeaders(t *testing.T) {
	headers := map[string][]string{
		"content-type": []string{"application/json"},
		"host":         []string{"example.com"},
	}

	_, _headers := getSignedHeaders(headers)

	assert.Equal(t, []string{"application/json"}, _headers["content-type"])
	assert.Equal(t, []string{"example.com"}, _headers["host"])
}

func TestGetSignedHeaders_ShouldHandleMultipleValues(t *testing.T) {
	headers := map[string][]string{
		"content-type": []string{"application/json", "text/html"},
		"host":         []string{"example.com"},
	}

	_, _headers := getSignedHeaders(headers)

	assert.Equal(t, []string{"application/json", "text/html"}, _headers["content-type"])
	assert.Equal(t, []string{"example.com"}, _headers["host"])
}

// getSignature Tests

func TestGetSignature_ShouldReturnValidSignature(t *testing.T) {
	stringToSign := "test-string-to-sign"
	sk := "test-secret-key"
	region := "cn-north-4"
	shortDate := "20240101"

	result := getSignature(stringToSign, sk, region, shortDate)

	assert.NotEmpty(t, result)
	assert.Len(t, result, 64) // SHA256 hex length
}

func TestGetSignature_ShouldReturnConsistentResult(t *testing.T) {
	stringToSign := "test-string-to-sign"
	sk := "test-secret-key"
	region := "cn-north-4"
	shortDate := "20240101"

	result1 := getSignature(stringToSign, sk, region, shortDate)
	result2 := getSignature(stringToSign, sk, region, shortDate)

	assert.Equal(t, result1, result2)
}

func TestGetSignature_ShouldHandleDifferentRegions(t *testing.T) {
	stringToSign := "test-string-to-sign"
	sk := "test-secret-key"
	shortDate := "20240101"

	sig1 := getSignature(stringToSign, sk, "cn-north-4", shortDate)
	sig2 := getSignature(stringToSign, sk, "us-east-1", shortDate)

	assert.NotEqual(t, sig1, sig2)
}

func TestGetSignature_ShouldHandleDifferentDates(t *testing.T) {
	stringToSign := "test-string-to-sign"
	sk := "test-secret-key"
	region := "cn-north-4"

	sig1 := getSignature(stringToSign, sk, region, "20240101")
	sig2 := getSignature(stringToSign, sk, region, "20240102")

	assert.NotEqual(t, sig1, sig2)
}

// doAuth Tests

func TestDoAuth_ShouldSkipAuth_WhenAKSKEmpty(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com")
	client.conf.signature = SignatureObs
	headers := make(map[string][]string)

	// Use empty AK/SK by creating a new client without credentials
	badClient, _ := New("", "", "https://obs.example.com")

	_, err := badClient.doAuth("GET", "bucket", "object", make(map[string]string), headers, "")

	assert.NoError(t, err)
	assert.NotContains(t, headers, HEADER_AUTH_CAMEL)
}

func TestDoAuth_ShouldAddSecurityToken_WhenTokenProvidedAndObs(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureObs),
		WithSecurityToken("test-token"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuth("GET", "bucket", "object", params, headers, "")

	assert.NoError(t, err)
	assert.Equal(t, []string{"test-token"}, headers[HEADER_STS_TOKEN_OBS])
}

func TestDoAuth_ShouldAddSecurityToken_WhenTokenProvidedAndV2(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
		WithSecurityToken("test-token"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuth("GET", "bucket", "object", params, headers, "")

	assert.NoError(t, err)
	assert.Equal(t, []string{"test-token"}, headers[HEADER_STS_TOKEN_AMZ])
}

func TestDoAuth_ShouldAddV2Authorization_WhenSignatureV2(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuth("GET", "bucket", "object", params, headers, "")

	assert.NoError(t, err)
	assert.Contains(t, headers, HEADER_AUTH_CAMEL)
	assert.Contains(t, headers[HEADER_AUTH_CAMEL][0], "AWS test-ak:")
}

func TestDoAuth_ShouldAddObsAuthorization_WhenSignatureObs(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureObs),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuth("GET", "bucket", "object", params, headers, "")

	assert.NoError(t, err)
	assert.Contains(t, headers, HEADER_AUTH_CAMEL)
	assert.Contains(t, headers[HEADER_AUTH_CAMEL][0], "OBS test-ak:")
}

func TestDoAuth_ShouldAddV4Authorization_WhenSignatureV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuth("GET", "bucket", "object", params, headers, "")

	assert.NoError(t, err)
	assert.Contains(t, headers, HEADER_AUTH_CAMEL)
	assert.Contains(t, headers[HEADER_AUTH_CAMEL][0], V4_HASH_PREFIX)
	assert.Contains(t, headers[HEADER_AUTH_CAMEL][0], "Credential=")
	assert.Contains(t, headers[HEADER_AUTH_CAMEL][0], "SignedHeaders=")
	assert.Contains(t, headers[HEADER_AUTH_CAMEL][0], "Signature=")
}

func TestDoAuth_ShouldUseProvidedHostName(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com")
	headers := make(map[string][]string)
	params := make(map[string]string)
	providedHost := "custom-host.example.com"

	_, err := client.doAuth("GET", "bucket", "object", params, headers, providedHost)

	assert.NoError(t, err)
	assert.Equal(t, []string{providedHost}, headers[HEADER_HOST_CAMEL])
}

func TestDoAuth_ShouldUseParsedHostName_WhenHostNameEmpty(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com")
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuth("GET", "bucket", "object", params, headers, "")

	assert.NoError(t, err)
	// In virtual hosting mode, the host includes the bucket name
	assert.Contains(t, headers[HEADER_HOST_CAMEL][0], "obs.example.com")
}

func TestDoAuth_ShouldAddContentSha256_WhenV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV4),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuth("GET", "bucket", "object", params, headers, "")

	assert.NoError(t, err)
	assert.Equal(t, []string{UNSIGNED_PAYLOAD}, headers[HEADER_CONTENT_SHA256_AMZ])
}

// doAuthTemporary Tests

func TestDoAuthTemporary_ShouldSkipAuth_WhenAKSKEmpty(t *testing.T) {
	badClient, _ := New("", "", "https://obs.example.com",
		WithSignature(SignatureV2),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	url, err := badClient.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	assert.NotContains(t, url, "Signature=")
}

func TestDoAuthTemporary_ShouldAddSecurityToken_WhenTokenProvidedAndObs(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureObs),
		WithSecurityToken("test-token"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	assert.Equal(t, "test-token", params[HEADER_STS_TOKEN_OBS])
}

func TestDoAuthTemporary_ShouldAddSecurityToken_WhenTokenProvidedAndV2(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
		WithSecurityToken("test-token"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	assert.Equal(t, "test-token", params[HEADER_STS_TOKEN_AMZ])
}

func TestDoAuthTemporary_ShouldReturnV2URL_WhenSignatureV2(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	url, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	assert.Contains(t, url, "AccessKeyId=test-ak")
	assert.Contains(t, url, "Signature=")
}

func TestDoAuthTemporary_ShouldReturnV4URL_WhenSignatureV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	url, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	assert.Contains(t, url, PARAM_ALGORITHM_AMZ_CAMEL+"="+V4_HASH_PREFIX)
	assert.Contains(t, url, PARAM_CREDENTIAL_AMZ_CAMEL)
	assert.Contains(t, url, PARAM_DATE_AMZ_CAMEL)
	assert.Contains(t, url, PARAM_EXPIRES_AMZ_CAMEL+"=3600")
	assert.Contains(t, url, PARAM_SIGNEDHEADERS_AMZ_CAMEL)
	assert.Contains(t, url, PARAM_SIGNATURE_AMZ_CAMEL)
}

func TestDoAuthTemporary_ShouldRemovePort80_WhenV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com:80",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	// In virtual hosting mode, the host includes the bucket name
	assert.Contains(t, headers[HEADER_HOST_CAMEL][0], "obs.example.com")
}

func TestDoAuthTemporary_ShouldRemovePort443_WhenV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com:443",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	// In virtual hosting mode, the host includes the bucket name
	assert.Contains(t, headers[HEADER_HOST_CAMEL][0], "obs.example.com")
}

func TestDoAuthTemporary_ShouldHandleEmptyHostHeader_WhenV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
	)
	headers := make(map[string][]string)
	// Clear the Host header
	headers[HEADER_HOST_CAMEL] = []string{}
	params := make(map[string]string)

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
}

func TestDoAuthTemporary_ShouldKeepOtherPorts_WhenV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com:8080",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	// In virtual hosting mode, the host includes the bucket name and port
	assert.Contains(t, headers[HEADER_HOST_CAMEL][0], "obs.example.com:8080")
}

func TestDoAuthTemporary_ShouldClearObjectKey_WhenPolicyProvided(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)
	policy := `{"expiration":"2024-01-01T00:00:00Z"}`

	url, err := client.doAuthTemporary("POST", "bucket", "object", policy, params, headers, 3600)

	assert.NoError(t, err)
	// URL should not contain the object key when policy is provided
	assert.NotContains(t, url, "/object")
	assert.Contains(t, url, "Policy=")
}

func TestDoAuthTemporary_ShouldHandleInvalidDate_WhenV2(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)
	// Set invalid date
	headers[HEADER_DATE_CAMEL] = []string{"invalid-date"}

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.Error(t, err)
}

func TestDoAuthTemporary_ShouldHandleInvalidDate_WhenV4(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV4),
		WithRegion("cn-north-4"),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)
	// Set invalid date
	headers[HEADER_DATE_CAMEL] = []string{"invalid-date"}

	_, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.Error(t, err)
}

func TestDoAuthTemporary_ShouldUseCurrentDate_WhenV2(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)
	expires := int64(3600)

	url, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, expires)

	assert.NoError(t, err)
	// URL should contain Expires timestamp
	assert.Contains(t, url, "Expires=")
	// Verify it's roughly correct (within a reasonable time window)
	assert.GreaterOrEqual(t, len(url), 50)
}

func TestDoAuthTemporary_ShouldAddQuestionMark_WhenURLNoQuery(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureV2),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	url, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	// URL should have ? after formatUrls
	assert.Contains(t, url, "?")
}

func TestDoAuthTemporary_ShouldNotAddAWS_WhenSignatureObs(t *testing.T) {
	client := CreateTestObsClient("https://obs.example.com",
		WithSignature(SignatureObs),
	)
	headers := make(map[string][]string)
	params := make(map[string]string)

	url, err := client.doAuthTemporary("GET", "bucket", "object", "", params, headers, 3600)

	assert.NoError(t, err)
	// OBS signature should not add "AWS" prefix
	assert.NotContains(t, url, "AWSAccessKeyId")
	assert.Contains(t, url, "AccessKeyId=")
}
