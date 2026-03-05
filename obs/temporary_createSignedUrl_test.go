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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateSignedUrl_ShouldReturnError_WhenInputIsNil tests nil input
func TestCreateSignedUrl_ShouldReturnError_WhenInputIsNil(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com")

	output, err := client.CreateSignedUrl(nil)

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CreateSignedUrlInput is nil")
}

// TestCreateSignedUrl_ShouldReturnSuccess_WhenInputIsValid tests valid input
func TestCreateSignedUrl_ShouldReturnSuccess_WhenInputIsValid(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateSignedUrlInput{
		Method:  HttpMethodGet,
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateSignedUrl(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.SignedUrl)
	assert.Contains(t, output.SignedUrl, "test-bucket")
	assert.NotNil(t, output.ActualSignedRequestHeaders)
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithDifferentHTTPMethods tests different HTTP methods
func TestCreateSignedUrl_ShouldReturnSuccess_WithDifferentHTTPMethods(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	methods := []HttpMethodType{
		HttpMethodGet,
		HttpMethodPut,
		HttpMethodPost,
		HttpMethodDelete,
		HttpMethodHead,
		HttpMethodOptions,
	}

	for _, method := range methods {
		t.Run(string(method), func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:  method,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
			}

			output, err := client.CreateSignedUrl(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.SignedUrl)
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithQueryParams tests query parameters
func TestCreateSignedUrl_ShouldReturnSuccess_WithQueryParams(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	testCases := []struct {
		name        string
		queryParams map[string]string
		contains    []string
	}{
		{
			name:        "single query param",
			queryParams: map[string]string{"response-content-type": "application/json"},
			contains:    []string{"response-content-type"},
		},
		{
			name: "multiple query params",
			queryParams: map[string]string{
				"response-content-type":  "application/json",
				"response-cache-control": "no-cache",
			},
			contains: []string{"response-content-type", "response-cache-control"},
		},
		{
			name:        "query params with special chars",
			queryParams: map[string]string{"response-content-disposition": "attachment; filename=test.txt"},
			contains:    []string{"response-content-disposition"},
		},
		{
			name: "complex query params",
			queryParams: map[string]string{
				"response-content-disposition": "attachment; filename=test.txt",
				"response-content-type":        "application/octet-stream",
			},
			contains: []string{"response-content-disposition", "response-content-type"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:      HttpMethodGet,
				Bucket:      "test-bucket",
				Key:         "test-key",
				Expires:     600,
				QueryParams: tc.queryParams,
			}

			output, err := client.CreateSignedUrl(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.SignedUrl)
			for _, substr := range tc.contains {
				assert.Contains(t, output.SignedUrl, substr)
			}
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithSubResource tests subResource
func TestCreateSignedUrl_ShouldReturnSuccess_WithSubResource(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	subResources := []SubResourceType{
		SubResourceAcl,
		SubResourceCors,
		SubResourceLifecycle,
		SubResourcePolicy,
		SubResourceTagging,
		SubResourceVersioning,
		SubResourceWebsite,
	}

	for _, subResource := range subResources {
		t.Run(string(subResource), func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:      HttpMethodGet,
				Bucket:      "test-bucket",
				Key:         "test-key",
				Expires:     600,
				SubResource: subResource,
			}

			output, err := client.CreateSignedUrl(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.SignedUrl)
			assert.Contains(t, output.SignedUrl, string(subResource))
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithHeaders tests headers
func TestCreateSignedUrl_ShouldReturnSuccess_WithHeaders(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateSignedUrlInput{
		Method:  HttpMethodPut,
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Content-Encoding": "gzip",
		},
	}

	output, err := client.CreateSignedUrl(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.SignedUrl)
	assert.NotNil(t, output.ActualSignedRequestHeaders)
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithDifferentExpiresValues tests different expires values
func TestCreateSignedUrl_ShouldReturnSuccess_WithDifferentExpiresValues(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	testCases := []struct {
		name    string
		expires int
	}{
		{name: "zero expires", expires: 0},
		{name: "negative expires", expires: -100},
		{name: "100 seconds", expires: 100},
		{name: "300 seconds", expires: 300},
		{name: "600 seconds", expires: 600},
		{name: "1 hour", expires: 3600},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: tc.expires,
			}

			output, err := client.CreateSignedUrl(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.SignedUrl)
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithPolicy tests with policy
func TestCreateSignedUrl_ShouldReturnSuccess_WithPolicy(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	policy := `{"expiration":"2025-12-31T00:00:00Z","conditions":[["eq","$bucket","test-bucket"]]}`

	input := &CreateSignedUrlInput{
		Method:  HttpMethodPost,
		Bucket:  "test-bucket",
		Key:     "test-key",
		Policy:  policy,
		Expires: 600,
	}

	output, err := client.CreateSignedUrl(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.SignedUrl)
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithSignatureVersions tests different signature versions
func TestCreateSignedUrl_ShouldReturnSuccess_WithSignatureVersions(t *testing.T) {
	testCases := []struct {
		name        string
		configurers []interface{}
	}{
		{
			name:        "V2 signature",
			configurers: []interface{}{WithSignature(SignatureV2)},
		},
		{
			name:        "V4 signature",
			configurers: []interface{}{WithSignature(SignatureV4), WithRegion("cn-north-4")},
		},
		{
			name:        "OBS signature",
			configurers: []interface{}{WithSignature(SignatureObs)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := CreateTestObsClient("https://obs.test.example.com", tc.configurers...)

			input := &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
			}

			output, err := client.CreateSignedUrl(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.SignedUrl)
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithAllFields tests all fields
func TestCreateSignedUrl_ShouldReturnSuccess_WithAllFields(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateSignedUrlInput{
		Method:  HttpMethodPut,
		Bucket:  "test-bucket",
		Key:     "test-key/file.txt",
		Policy:  `{"expiration":"2025-12-31T00:00:00Z"}`,
		Expires: 900,
		QueryParams: map[string]string{
			"response-content-type":  "application/json",
			"response-cache-control": "no-cache",
		},
		Headers: map[string]string{
			"Content-Type":  "text/plain",
			"Content-MD5":   "d41d8cd98f00b204e9800998ecf8427e",
			"Cache-Control": "no-cache",
		},
	}

	output, err := client.CreateSignedUrl(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.SignedUrl)
	assert.NotNil(t, output.ActualSignedRequestHeaders)
}

// TestCreateSignedUrl_ShouldIncludeHeadersInOutput tests headers in output
func TestCreateSignedUrl_ShouldIncludeHeadersInOutput(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateSignedUrlInput{
		Method:  HttpMethodPut,
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "no-cache",
		},
	}

	output, err := client.CreateSignedUrl(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotNil(t, output.ActualSignedRequestHeaders)
	assert.NotEmpty(t, output.ActualSignedRequestHeaders.Get("Content-Type"))
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithEmptyOrNilValues tests empty and nil values
func TestCreateSignedUrl_ShouldReturnSuccess_WithEmptyOrNilValues(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	testCases := []struct {
		name        string
		input       *CreateSignedUrlInput
		description string
	}{
		{
			name: "empty bucket",
			input: &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "",
				Key:     "test-key",
				Expires: 600,
			},
			description: "should handle empty bucket",
		},
		{
			name: "empty key",
			input: &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "test-bucket",
				Key:     "",
				Expires: 600,
			},
			description: "should handle empty key",
		},
		{
			name: "empty subResource",
			input: &CreateSignedUrlInput{
				Method:      HttpMethodGet,
				Bucket:      "test-bucket",
				Key:         "test-key",
				Expires:     600,
				SubResource: "",
			},
			description: "should handle empty subResource",
		},
		{
			name: "empty QueryParams",
			input: &CreateSignedUrlInput{
				Method:      HttpMethodGet,
				Bucket:      "test-bucket",
				Key:         "test-key",
				Expires:     600,
				QueryParams: map[string]string{},
			},
			description: "should handle empty QueryParams",
		},
		{
			name: "nil QueryParams",
			input: &CreateSignedUrlInput{
				Method:      HttpMethodGet,
				Bucket:      "test-bucket",
				Key:         "test-key",
				Expires:     600,
				QueryParams: nil,
			},
			description: "should handle nil QueryParams",
		},
		{
			name: "empty Headers",
			input: &CreateSignedUrlInput{
				Method:  HttpMethodPut,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
				Headers: map[string]string{},
			},
			description: "should handle empty Headers",
		},
		{
			name: "nil Headers",
			input: &CreateSignedUrlInput{
				Method:  HttpMethodPut,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
				Headers: nil,
			},
			description: "should handle nil Headers",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := client.CreateSignedUrl(tc.input)

			assert.NoError(t, err, tc.description)
			assert.NotNil(t, output, tc.description)
			assert.NotEmpty(t, output.SignedUrl, tc.description)
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithSpecialCharactersInKey tests special characters in key
func TestCreateSignedUrl_ShouldReturnSuccess_WithSpecialCharactersInKey(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	testCases := []struct {
		name string
		key  string
	}{
		{
			name: "key with spaces",
			key:  "test/key/with spaces.txt",
		},
		{
			name: "key with unicode",
			key:  "测试文件.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "test-bucket",
				Key:     tc.key,
				Expires: 600,
			}

			output, err := client.CreateSignedUrl(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.SignedUrl)
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithExtensions tests extension options
func TestCreateSignedUrl_ShouldReturnSuccess_WithExtensions(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	testCases := []struct {
		name        string
		description string
	}{
		{
			name:        "with all extensions",
			description: "should handle all extension options",
		},
		{
			name:        "with custom header",
			description: "should handle custom header",
		},
		{
			name:        "with single extension",
			description: "should handle single extension",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:  HttpMethodPut,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
			}

			var output *CreateSignedUrlOutput
			var err error

			switch tc.name {
			case "with all extensions":
				output, err = client.CreateSignedUrl(input,
					WithReqPaymentHeader(RequesterPayer),
					WithTrafficLimitHeader(819200),
				)
			case "with custom header":
				output, err = client.CreateSignedUrl(input,
					WithCustomHeader("x-obs-test-header", "test-value"),
				)
			case "with single extension":
				output, err = client.CreateSignedUrl(input,
					WithReqPaymentHeader(BucketOwnerPayer),
				)
			}

			assert.NoError(t, err, tc.description)
			assert.NotNil(t, output, tc.description)
			assert.NotEmpty(t, output.SignedUrl, tc.description)
			assert.NotNil(t, output.ActualSignedRequestHeaders, tc.description)
		})
	}
}

// TestCreateSignedUrl_ShouldHandleExtensions_WhenInvalidTypesAreProvided tests handling of invalid extension types
func TestCreateSignedUrl_ShouldHandleExtensions_WhenInvalidTypesAreProvided(t *testing.T) {
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	testCases := []struct {
		name        string
		description string
	}{
		{
			name:        "string extension",
			description: "should handle string extension type",
		},
		{
			name:        "number extension",
			description: "should handle number extension type",
		},
		{
			name:        "multiple invalid extensions",
			description: "should handle multiple invalid extension types",
		},
		{
			name:        "mixed valid and invalid extensions",
			description: "should handle mixed valid and invalid extensions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
			}

			var output *CreateSignedUrlOutput
			var err error

			switch tc.name {
			case "string extension":
				// Pass a string directly - will be ignored
				output, err = client.CreateSignedUrl(input, "unsupported-extension")
			case "number extension":
				// Pass a number directly - will be ignored
				output, err = client.CreateSignedUrl(input, 123)
			case "multiple invalid extensions":
				// Pass multiple invalid types - all will be ignored
				output, err = client.CreateSignedUrl(input, "extension1", 123, true)
			case "mixed valid and invalid extensions":
				// Mix valid and invalid extensions
				output, err = client.CreateSignedUrl(input,
					WithReqPaymentHeader(RequesterPayer),
					"invalid-extension",
					WithTrafficLimitHeader(819200),
					123,
				)
			}

			// Should not fail, just log warnings
			assert.NoError(t, err, tc.description)
			assert.NotNil(t, output, tc.description)
			assert.NotEmpty(t, output.SignedUrl, tc.description)
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithCombinedFields tests combined QueryParams, Headers and SubResource
func TestCreateSignedUrl_ShouldReturnSuccess_WithCombinedFields(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	testCases := []struct {
		name        string
		input       *CreateSignedUrlInput
		contains    []string
		description string
	}{
		{
			name: "QueryParams with SubResource",
			input: &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
				QueryParams: map[string]string{
					"response-cache-control": "no-cache",
				},
				SubResource: SubResourceAcl,
			},
			contains:    []string{"response-cache-control", "acl"},
			description: "should handle QueryParams with SubResource",
		},
		{
			name: "QueryParams, Headers and SubResource",
			input: &CreateSignedUrlInput{
				Method:  HttpMethodPut,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: 600,
				QueryParams: map[string]string{
					"response-content-type": "application/json",
				},
				SubResource: SubResourceCors,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},
			contains:    []string{"response-content-type", "cors"},
			description: "should handle QueryParams, Headers and SubResource together",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := client.CreateSignedUrl(tc.input)

			assert.NoError(t, err, tc.description)
			assert.NotNil(t, output, tc.description)
			assert.NotEmpty(t, output.SignedUrl, tc.description)
			assert.NotNil(t, output.ActualSignedRequestHeaders, tc.description)
			for _, substr := range tc.contains {
				assert.Contains(t, output.SignedUrl, substr)
			}
		})
	}
}

// TestCreateSignedUrl_ShouldReturnSuccess_WithEmptyAKSK tests with empty AK/SK
func TestCreateSignedUrl_ShouldReturnSuccess_WithEmptyAKSK(t *testing.T) {
	client, err := New("", "", "https://obs.test.example.com",
		WithSignature(SignatureObs))
	assert.NoError(t, err)

	input := &CreateSignedUrlInput{
		Method:  HttpMethodGet,
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateSignedUrl(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.SignedUrl)
}

// TestCreateSignedUrl_ShouldReturnSuccess_WhenNoExtensionsProvided tests without extensions
func TestCreateSignedUrl_ShouldReturnSuccess_WhenNoExtensionsProvided(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateSignedUrlInput{
		Method:  HttpMethodGet,
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	output, err := client.CreateSignedUrl(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.SignedUrl)
}

// TestCreateSignedUrl_ShouldLogInfo_WhenExtensionHeaderReturnsError tests extension header error logging
func TestCreateSignedUrl_ShouldLogInfo_WhenExtensionHeaderReturnsError(t *testing.T) {
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	input := &CreateSignedUrlInput{
		Method:  HttpMethodPut,
		Bucket:  "test-bucket",
		Key:     "test-key",
		Expires: 600,
	}

	// Create an extension that returns an error
	extensionWithEmptyValue := setHeaderPrefix("x-obs-test", "")
	extensionAsInterface := extensionHeaders(extensionWithEmptyValue)

	output, err := client.CreateSignedUrl(input, extensionAsInterface)

	// Should not fail, just log the error
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.SignedUrl)
}

// TestCreateSignedUrl_ShouldUseProperIntConversion_WhenNamingSubtests ensures proper integer to string conversion
func TestCreateSignedUrl_ShouldUseProperIntConversion_WhenNamingSubtests(t *testing.T) {
	client := CreateTestObsClient("https://obs.test.example.com",
		WithSignature(SignatureObs))

	expiresValues := []int{100, 300, 600, 3600}

	for _, expires := range expiresValues {
		t.Run(fmt.Sprintf("%d", expires), func(t *testing.T) {
			input := &CreateSignedUrlInput{
				Method:  HttpMethodGet,
				Bucket:  "test-bucket",
				Key:     "test-key",
				Expires: expires,
			}

			output, err := client.CreateSignedUrl(input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.SignedUrl)
		})
	}
}
