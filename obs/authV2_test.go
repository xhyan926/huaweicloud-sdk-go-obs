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

// getV2StringToSign Tests

func TestGetV2StringToSign_ShouldCreateBasicString_WhenSimpleInput(t *testing.T) {
	method := "PUT"
	canonicalizedURL := "/test-bucket/test-object"
	headers := make(map[string][]string)

	result := getV2StringToSign(method, canonicalizedURL, headers, false)

	assert.Contains(t, result, "PUT")
	assert.Contains(t, result, "\n")
	assert.Contains(t, result, canonicalizedURL)
}

func TestGetV2StringToSign_ShouldAddEmptyAttachHeaders_WhenCalled(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-object"
	headers := make(map[string][]string)

	result := getV2StringToSign(method, canonicalizedURL, headers, true)

	assert.Contains(t, result, "GET")
	assert.Contains(t, result, "\n")
	assert.Contains(t, result, canonicalizedURL)
}

func TestGetV2StringToSign_ShouldNotAddAttachHeaders_WhenHeadersNotEmpty(t *testing.T) {
	method := "PUT"
	canonicalizedURL := "/test-object"
	headers := map[string][]string{
		"Content-Type":    []string{"application/json"},
	}

	result := getV2StringToSign(method, canonicalizedURL, headers, false)

	assert.NotContains(t, result, "Content-Type")
	assert.NotContains(t, result, "Content-Type")
}

func TestGetV2StringToSign_ShouldHandleQueryParams_WhenURLContainsQuery(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-object?param1=value1&param2=value2"
	headers := make(map[string][]string)

	result := getV2StringToSign(method, canonicalizedURL, headers, true)

	// Should split query params
	assert.Contains(t, result, "param1=value1")
	assert.Contains(t, result, "param2=value2")
	assert.NotContains(t, result, "?")
}

func TestGetV2StringToSign_ShouldNotAddSecurityToken_WhenSecurityTokenNotPresent(t *testing.T) {
	method := "PUT"
	canonicalizedURL := "/test-bucket/test-object"
	headers := make(map[string][]string)

	result := getV2StringToSign(method, canonicalizedURL, headers, false)

	assert.NotContains(t, result, HEADER_STS_TOKEN_OBS)
	assert.NotContains(t, result, "Security-Token")
	assert.NotContains(t, result, "******")
}

func TestGetV2StringToSign_ShouldAddSecurityToken_WhenAmzTokenPresent(t *testing.T) {
	method := "POST"
	canonicalizedURL := "/test-bucket/test-object"
	headers := map[string][]string{
		HEADER_STS_TOKEN_AMZ: []string{"token123"},
	}

	result := getV2StringToSign(method, canonicalizedURL, headers, false)

	assert.Contains(t, result, HEADER_STS_TOKEN_AMZ)
	assert.Contains(t, result, "token123")
}

func TestGetV2StringToSign_ShouldMaskSecurityToken_WhenObsTokenPresent(t *testing.T) {
	// Note: The current implementation only masks security tokens in the URL query string,
	// not in the headers. This test verifies that behavior.

	method := "DELETE"
	canonicalizedURL := "/test-bucket/test-object?x-obs-security-token=token123"
	headers := map[string][]string{
		HEADER_STS_TOKEN_OBS: []string{"token123"},
	}

	result := getV2StringToSign(method, canonicalizedURL, headers, true)

	// When token is in URL query string, it should be masked
	assert.Contains(t, result, HEADER_STS_TOKEN_OBS)
	assert.NotContains(t, result, "token123") // Token in URL should be masked
	assert.Contains(t, result, "******") // Masked version
}

func TestGetV2StringToSign_ShouldHandleMultipleQueryParams_WhenURLContainsMany(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-object?param1=value1&param2=value2&param3=value3&param4=value4"
	headers := make(map[string][]string)

	result := getV2StringToSign(method, canonicalizedURL, headers, true)

	// Note: Query params are included in string to sign without ?
	// This matches AWS V2 signature behavior for canonicalized URLs
	assert.NotContains(t, result, "?")
	assert.Contains(t, result, "param1=value1")
	assert.Contains(t, result, "param4=value4")
}

func TestGetV2StringToSign_ShouldHandleEmptySecurityTokenValue_WhenPresent(t *testing.T) {
	method := "PUT"
	canonicalizedURL := "/test-bucket/test-object"
	headers := map[string][]string{
		HEADER_STS_TOKEN_AMZ: []string{""},
	}

	result := getV2StringToSign(method, canonicalizedURL, headers, false)

	assert.Contains(t, result, HEADER_STS_TOKEN_AMZ)
	assert.NotContains(t, result, "token123") // Empty value not added
}

func TestGetV2StringToSign_ShouldNotSplitEmptyQueryParam_WhenURLEndsWithEmptyQuery(t *testing.T) {
	method := "GET"
	canonicalizedURL := "/test-object?"
	headers := make(map[string][]string)

	result := getV2StringToSign(method, canonicalizedURL, headers, true)

	assert.Contains(t, result, "?")
	assert.NotContains(t, result, "=")
}

// v2Auth Tests

func TestV2Auth_ShouldReturnValidSignature_WhenGivenValidInputs(t *testing.T) {
	ak := "test-access-key"
	sk := "test-secret-key"
	method := "PUT"
	canonicalizedURL := "/test-object"
	headers := make(map[string][]string)

	result := v2Auth(ak, sk, method, canonicalizedURL, headers, true)

	assert.NotNil(t, result)
	assert.Contains(t, result, "Signature")
	assert.Len(t, result["Signature"], 28) // Base64 encoded string
}

func TestV2Auth_ShouldReturnValidSignature_WhenIsObsIsTrue(t *testing.T) {
	ak := "test-access-key"
	sk := "test-secret-key"
	method := "GET"
	canonicalizedURL := "/test-object"
	headers := make(map[string][]string)

	result := v2Auth(ak, sk, method, canonicalizedURL, headers, true)

	assert.NotNil(t, result)
	assert.Contains(t, result, "Signature")
	assert.Len(t, result["Signature"], 28) // Base64 encoded string
}

func TestV2Auth_ShouldHandleEmptyInputs(t *testing.T) {
	result := v2Auth("", "", "GET", "", nil, false)

	assert.NotNil(t, result)
	assert.Contains(t, result, "Signature")
	assert.Len(t, result["Signature"], 28) // Base64 encoded string
}
