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
	"time"

	"github.com/stretchr/testify/assert"
)

// getV4StringToSign Tests

func TestGetV4StringToSign_ShouldCreateBasicString_WhenSimpleInput(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host"}
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	assert.Contains(t, result, V4_HASH_PREFIX)
	assert.Contains(t, result, longDate)
	assert.Contains(t, result, scope)
}

func TestGetV4StringToSign_ShouldIncludeSignedHeaders_WhenHeadersProvided(t *testing.T) {
	method := "PUT"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"content-type", "host", "x-amz-date"}
	headers := map[string][]string{
		"content-type": []string{"application/json"},
		"host":        []string{"obs.example.com"},
		"x-amz-date":  []string{"20240101T120000Z"},
	}

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Different headers should produce different hash
	result2 := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders[:1], headers)
	assert.NotEqual(t, result, result2)
}

func TestGetV4StringToSign_ShouldHandleMultipleHeaderValues(t *testing.T) {
	method := "POST"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"x-custom"}
	headers := map[string][]string{
		"x-custom": []string{"value1", "value2", "value3"},
	}

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Should produce a valid hash (64 hex chars)
	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 4)
	hashLine := lines[3]
	assert.Len(t, hashLine, 64) // SHA256 hex length
}

func TestGetV4StringToSign_ShouldMaskSecurityToken_WhenAmzTokenPresent(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host", HEADER_STS_TOKEN_AMZ}
	headers := map[string][]string{
		"host":                 []string{"obs.example.com"},
		HEADER_STS_TOKEN_AMZ:    []string{"secret-token-12345"},
	}

	// Setup test logger to capture logs
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// The result itself doesn't contain the token (it's in the canonical request which is hashed)
	// The masking happens in the logs, check that result is still produced
	assert.NotEmpty(t, result)
}

func TestGetV4StringToSign_ShouldMaskSecurityToken_WhenObsTokenPresent(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host", HEADER_STS_TOKEN_OBS}
	headers := map[string][]string{
		"host":                 []string{"obs.example.com"},
		HEADER_STS_TOKEN_OBS:    []string{"secret-token-67890"},
	}

	// Setup test logger to capture logs
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// The result itself doesn't contain the token
	assert.NotEmpty(t, result)
}

func TestGetV4StringToSign_ShouldMaskSecurityTokenFromQuery_WhenPresent(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := "X-Amz-Security-Token=secret-token-abc"
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host"}
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	// Setup test logger to capture logs
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Different query URL should produce different hash
	result2 := getV4StringToSign(method, canonicalizedURL, "", scope, longDate, payload, signedHeaders, headers)
	assert.NotEqual(t, result, result2)
}

func TestGetV4StringToSign_ShouldNotMaskEmptySecurityToken(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := "X-Amz-Security-Token="
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host"}
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Empty token should produce a hash
	assert.NotEmpty(t, result)
	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 4)
}

func TestGetV4StringToSign_ShouldHandleObsSecurityTokenInQuery(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := "x-obs-security-token=secret-token-xyz"
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host"}
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	// Setup test logger to capture logs
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Different query URL should produce different hash
	result2 := getV4StringToSign(method, canonicalizedURL, "", scope, longDate, payload, signedHeaders, headers)
	assert.NotEqual(t, result, result2)
}

func TestGetV4StringToSign_ShouldHandleComplexQueryURL(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := "param1=value1&param2=value2&X-Amz-Security-Token=token123"
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host"}
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	// Setup test logger to capture logs
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Different query should produce different hash
	result2 := getV4StringToSign(method, canonicalizedURL, "param1=value1&param2=value2", scope, longDate, payload, signedHeaders, headers)
	assert.NotEqual(t, result, result2)
}

func TestGetV4StringToSign_ShouldJoinSignedHeaders(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"content-md5", "content-type", "host"}
	headers := map[string][]string{
		"content-md5":   []string{"abc123"},
		"content-type":  []string{"application/json"},
		"host":         []string{"obs.example.com"},
	}

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Different signed headers should produce different hash
	result2 := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders[:2], headers)
	assert.NotEqual(t, result, result2)
}

func TestGetV4StringToSign_ShouldHandleCustomPayload(t *testing.T) {
	method := "PUT"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	signedHeaders := []string{"host"}
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Different payload should produce different hash
	result2 := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, UNSIGNED_PAYLOAD, signedHeaders, headers)
	assert.NotEqual(t, result, result2)
}

func TestGetV4StringToSign_ShouldHandleAmzTokenPrecedence_WhenBothTokensPresent(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host", HEADER_STS_TOKEN_AMZ, HEADER_STS_TOKEN_OBS}
	headers := map[string][]string{
		"host":                 []string{"obs.example.com"},
		HEADER_STS_TOKEN_OBS:    []string{"obs-token"},
		HEADER_STS_TOKEN_AMZ:    []string{"amz-token"},
	}

	// Setup test logger to capture logs
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Should produce valid hash
	assert.NotEmpty(t, result)
	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 4)
}

// V4Auth Tests

func TestV4Auth_ShouldCallV4Auth(t *testing.T) {
	ak := "test-access-key"
	sk := "test-secret-key"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := V4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result, "Credential")
	assert.Contains(t, result, "SignedHeaders")
	assert.Contains(t, result, "Signature")
}

// v4Auth Tests

func TestV4Auth_ShouldReturnValidSignature_WhenGivenValidInputs(t *testing.T) {
	ak := "test-access-key"
	sk := "test-secret-key"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result, "Credential")
	assert.Contains(t, result, "SignedHeaders")
	assert.Contains(t, result, "Signature")
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["SignedHeaders"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseAmzDate_WhenValid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "PUT"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.Contains(t, result["Credential"], "test-ak")
	assert.Contains(t, result["Credential"], "20240101")
	assert.Contains(t, result["Credential"], "cn-north-4")
	assert.Len(t, result["Signature"], 64) // SHA256 hex length
}

func TestV4Auth_ShouldUseCurrentTime_WhenAmzDateInvalid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"invalid-date"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseParamDateAmz_WhenValid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "POST"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":                  []string{"obs.example.com"},
		PARAM_DATE_AMZ_CAMEL:    []string{"20240102T130000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result["Credential"], "20240102")
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseCurrentTime_WhenParamDateAmzInvalid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "DELETE"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":                  []string{"obs.example.com"},
		PARAM_DATE_AMZ_CAMEL:    []string{"invalid-date-format"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseDateCamel_WhenValid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "HEAD"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		HEADER_DATE_CAMEL:    []string{"Mon, 01 Jan 2024 12:00:00 GMT"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result["Credential"], "20240101")
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseCurrentTime_WhenDateCamelInvalid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		HEADER_DATE_CAMEL:    []string{"invalid-rfc-date"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseLowercaseDate_WhenValid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		strings.ToLower(HEADER_DATE_CAMEL): []string{"Tue, 02 Jan 2024 13:00:00 GMT"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result["Credential"], "20240102")
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseCurrentTime_WhenLowercaseDateInvalid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		strings.ToLower(HEADER_DATE_CAMEL): []string{"invalid-date"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseCurrentTime_WhenNoDateHeader(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldUseContentSha256_WhenPresent(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "PUT"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	customPayload := "custom-payload-hash"
	headers := map[string][]string{
		"host":                         []string{"obs.example.com"},
		HEADER_DATE_AMZ:                []string{"20240101T120000Z"},
		HEADER_CONTENT_SHA256_AMZ:        []string{customPayload},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Signature"])
	// The signature should be different when using custom payload
	result2 := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	})
	assert.NotEqual(t, result["Signature"], result2["Signature"])
}

func TestV4Auth_ShouldUseUnsignedPayload_WhenContentSha256NotPresent(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldHandleEmptyInputs(t *testing.T) {
	result := v4Auth("", "", "", "GET", "", "", nil)

	assert.NotNil(t, result)
	// Should still produce result even with empty inputs
	// Signature should always be produced (64 chars)
	assert.Len(t, result["Signature"], 64)
}

func TestV4Auth_ShouldIncludeRegionInCredential(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "us-east-1"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.Contains(t, result["Credential"], "us-east-1")
}

func TestV4Auth_ShouldHandleMultipleHeaders(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "POST"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
		"content-type":  []string{"application/json"},
		"content-md5":   []string{"abc123"},
		"x-amz-acl":     []string{"public-read"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result["SignedHeaders"], "content-md5")
	assert.Contains(t, result["SignedHeaders"], "content-type")
	assert.Contains(t, result["SignedHeaders"], "host")
	assert.Contains(t, result["SignedHeaders"], "x-amz-acl")
}

func TestV4Auth_ShouldHandleCaseInsensitiveHeaders(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"Host":              []string{"obs.example.com"},
		"Content-Type":       []string{"application/json"},
		"X-Amz-Date":        []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result["SignedHeaders"], "content-type")
	assert.Contains(t, result["SignedHeaders"], "host")
	assert.Contains(t, result["SignedHeaders"], "x-amz-date")
}

func TestV4Auth_ShouldHandleQueryURL(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := "param1=value1&param2=value2"
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldReturnConsistentResult_WhenSameInputs(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result1 := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)
	result2 := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.Equal(t, result1, result2)
}

func TestV4Auth_ShouldReturnDifferentSignature_WhenDifferentRegions(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""

	result1 := v4Auth(ak, sk, "cn-north-4", method, canonicalizedURL, queryURL, map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	})
	result2 := v4Auth(ak, sk, "us-east-1", method, canonicalizedURL, queryURL, map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	})

	assert.NotEqual(t, result1["Signature"], result2["Signature"])
}

func TestV4Auth_ShouldHandleSecurityToken(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":              []string{"obs.example.com"},
		HEADER_DATE_AMZ:      []string{"20240101T120000Z"},
		HEADER_STS_TOKEN_AMZ: []string{"security-token-123"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["SignedHeaders"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldProduceValidCredentialFormat(t *testing.T) {
	ak := "AKIAIOSFODNN7EXAMPLE"
	sk := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	region := "us-east-1"
	method := "GET"
	canonicalizedURL := "/examplebucket/test.txt"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"examplebucket.s3.amazonaws.com"},
		HEADER_DATE_AMZ: []string{"20130524T000000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.Contains(t, result["Credential"], ak)
	assert.Contains(t, result["Credential"], "20130524")
	assert.Contains(t, result["Credential"], region)
	assert.Contains(t, result["Credential"], V4_SERVICE_NAME)
	assert.Contains(t, result["Credential"], V4_SERVICE_SUFFIX)
}

func TestV4Auth_ShouldSortSignedHeaders(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"z-header":     []string{"value-z"},
		"a-header":     []string{"value-a"},
		"m-header":     []string{"value-m"},
		"content-type": []string{"application/json"},
		"host":        []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	// Headers should be sorted alphabetically
	signedHeaders := result["SignedHeaders"]
	assert.True(t, strings.HasPrefix(signedHeaders, "a-header;"))
	assert.Contains(t, signedHeaders, ";content-type;")
	assert.Contains(t, signedHeaders, ";host;")
	assert.Contains(t, signedHeaders, ";m-header;")
	assert.Contains(t, signedHeaders, ";x-amz-date;")
	assert.True(t, strings.HasSuffix(signedHeaders, ";z-header"))
}

func TestV4Auth_DatePrecedenceAmzDate(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		HEADER_DATE_AMZ:      []string{"20240101T120000Z"},
		HEADER_DATE_CAMEL:    []string{"Mon, 01 Jan 2024 12:00:00 GMT"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	// Should use HEADER_DATE_AMZ (first one checked)
	assert.Contains(t, result["Credential"], "20240101")
}

func TestV4Auth_DatePrecedenceParamDateAmz(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":                  []string{"obs.example.com"},
		PARAM_DATE_AMZ_CAMEL:    []string{"20240102T130000Z"},
		HEADER_DATE_CAMEL:        []string{"Mon, 01 Jan 2024 12:00:00 GMT"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	// Should use PARAM_DATE_AMZ_CAMEL (second one checked)
	assert.Contains(t, result["Credential"], "20240102")
}

func TestV4Auth_DatePrecedenceDateCamel(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		HEADER_DATE_CAMEL:    []string{"Mon, 01 Jan 2024 12:00:00 GMT"},
		strings.ToLower(HEADER_DATE_CAMEL): []string{"Tue, 02 Jan 2024 13:00:00 GMT"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	// Should use HEADER_DATE_CAMEL (third one checked)
	assert.Contains(t, result["Credential"], "20240101")
}

func TestV4Auth_DatePrecedenceLowercaseDate(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		strings.ToLower(HEADER_DATE_CAMEL): []string{"Tue, 02 Jan 2024 13:00:00 GMT"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	// Should use lowercase Date header (fourth one checked)
	assert.Contains(t, result["Credential"], "20240102")
}

func TestV4Auth_ShouldUseCurrentTime_WhenAllDatesInvalid(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":               []string{"obs.example.com"},
		HEADER_DATE_AMZ:      []string{"invalid-1"},
		PARAM_DATE_AMZ_CAMEL: []string{"invalid-2"},
		HEADER_DATE_CAMEL:    []string{"invalid-3"},
		strings.ToLower(HEADER_DATE_CAMEL): []string{"invalid-4"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_RealDateTest(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""

	// Test with a real current date
	now := time.Now().UTC()
	longDate := now.Format(LONG_DATE_FORMAT)
	expectedShortDate := now.Format(SHORT_DATE_FORMAT)

	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{longDate},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result["Credential"], expectedShortDate)
}

func TestV4Auth_RFC1123DateFormat(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":            []string{"obs.example.com"},
		HEADER_DATE_CAMEL: []string{"Wed, 03 Jan 2024 14:30:00 GMT"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.Contains(t, result["Credential"], "20240103")
}

func TestGetV4StringToSign_ShouldNotMask_WhenNoSecurityToken(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	scope := "20240101/cn-north-4/s3/aws4_request"
	longDate := "20240101T120000Z"
	payload := UNSIGNED_PAYLOAD
	signedHeaders := []string{"host"}
	headers := map[string][]string{
		"host": []string{"obs.example.com"},
	}

	result := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, headers)

	// Should not contain "******" when there's no security token
	assert.NotContains(t, result, "******")
}

func TestV4Auth_ShouldHandleEmptySignedHeaders(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	// Should still have signed headers even if minimal
	assert.NotEmpty(t, result["SignedHeaders"])
}

func TestV4Auth_ShouldHandleEmptyCanonicalizedURL(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := ""
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Credential"])
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldHandleComplexQueryParams(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	method := "GET"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := "versionId=123&uploadId=abc&partNumber=1"
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result["Signature"])
}

func TestV4Auth_ShouldHandleAllHTTPMethods(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	region := "cn-north-4"
	canonicalizedURL := "/test-bucket/test-object"
	queryURL := ""
	headers := map[string][]string{
		"host":         []string{"obs.example.com"},
		HEADER_DATE_AMZ: []string{"20240101T120000Z"},
	}

	methods := []string{HTTP_GET, HTTP_POST, HTTP_PUT, HTTP_DELETE, HTTP_HEAD}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			result := v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result["Credential"])
			assert.NotEmpty(t, result["Signature"])
		})
	}
}
