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
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HandleHttpResponse Tests

func TestHandleHttpResponse_ShouldParseNormalResponse_WhenNoCallbackHeader(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
	}

	output := &BaseModel{}

	err := HandleHttpResponse(PUT_OBJECT, resp.Header, output, resp, true, false)

	assert.NoError(t, err)
}

func TestHandleHttpResponse_ShouldParseNormalResponse_WhenActionNotSupported(t *testing.T) {
	// GetBucket is not in supportCallbackActions, so should use normal response path
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"x-amz-callback": []string{"true"},
		},
		Body: io.NopCloser(strings.NewReader("")),
	}

	output := &BaseModel{}

	err := HandleHttpResponse("GetBucket", resp.Header, output, resp, true, false)

	// Should still work because GetBucket is not in supportCallbackActions
	assert.NoError(t, err)
}

func TestHandleHttpResponse_ShouldParseNormalResponse_WhenActionIsPutFile(t *testing.T) {
	// PUT_FILE is in supportCallbackActions but without callback header should work
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
	}

	output := &BaseModel{}

	err := HandleHttpResponse(PUT_FILE, resp.Header, output, resp, true, false)

	assert.NoError(t, err)
}

// copyHeaders Tests

func TestCopyHeaders_ShouldReturnNewMap_WhenInputIsNotNil(t *testing.T) {
	headers := map[string][]string{
		"Content-Type":   {"application/json"},
		"Cache-Control":  {"no-cache"},
		"User-Agent":     {"test-agent"},
		"X-Amz-Meta-Key": {"value"},
	}

	result := copyHeaders(headers)

	assert.NotNil(t, result)
	assert.Len(t, result, 4)
	assert.Equal(t, []string{"application/json"}, result["content-type"])
	assert.Equal(t, []string{"no-cache"}, result["cache-control"])
	assert.Equal(t, []string{"test-agent"}, result["user-agent"])
	assert.Equal(t, []string{"value"}, result["x-amz-meta-key"])
}

func TestCopyHeaders_ShouldReturnEmptyMap_WhenInputIsNil(t *testing.T) {
	result := copyHeaders(nil)

	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestCopyHeaders_ShouldPreserveAllValues_WhenHeaderHasMultipleValues(t *testing.T) {
	headers := map[string][]string{
		"X-Amz-Custom": {"value1", "value2", "value3"},
	}

	result := copyHeaders(headers)

	assert.Len(t, result["x-amz-custom"], 3)
	assert.Equal(t, "value1", result["x-amz-custom"][0])
	assert.Equal(t, "value2", result["x-amz-custom"][1])
	assert.Equal(t, "value3", result["x-amz-custom"][2])
}

func TestCopyHeaders_ShouldNotModifyOriginal_WhenResultIsModified(t *testing.T) {
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	result := copyHeaders(headers)
	result["new-key"] = []string{"new-value"}
	result["content-type"] = []string{"modified"}

	assert.NotContains(t, headers, "new-key")
	assert.Equal(t, []string{"application/json"}, headers["Content-Type"])
	assert.Equal(t, []string{"modified"}, result["content-type"])
}

// parseHeaders Tests

func TestParseHeaders_ShouldReturnDefaultV2_WhenNoAuthorizationHeader(t *testing.T) {
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	signature, region, signedHeaders := parseHeaders(headers)

	assert.Equal(t, "v2", signature)
	assert.Equal(t, "", region)
	assert.Equal(t, "", signedHeaders)
}

func TestParseHeaders_ShouldReturnV4_WhenAuthorizationHasV4Prefix(t *testing.T) {
	headers := map[string][]string{
		"authorization": {"AWS4-HMAC-SHA256 Credential=test-key/20240101/cn-north-1/obs/aws4_request,SignedHeaders=host;x-amz-date,Signature=abc123"},
	}

	signature, region, signedHeaders := parseHeaders(headers)

	assert.Equal(t, "v4", signature)
	assert.Equal(t, "cn-north-1", region)
	assert.Equal(t, "host;x-amz-date", signedHeaders)
}

func TestParseHeaders_ShouldReturnV2_WhenAuthorizationHasV2Prefix(t *testing.T) {
	headers := map[string][]string{
		"authorization": {"AWS test-ak:signature123"},
	}

	signature, region, signedHeaders := parseHeaders(headers)

	assert.Equal(t, "v2", signature)
	assert.Equal(t, "", region)
	assert.Equal(t, "", signedHeaders)
}

func TestParseHeaders_ShouldHandleCaseInsensitiveAuthorizationHeader(t *testing.T) {
	headers := map[string][]string{
		"authorization": {"AWS4-HMAC-SHA256 Credential=test-key/20240101/us-east-1/obs/aws4_request,SignedHeaders=host,Signature=xyz"},
	}

	signature, region, signedHeaders := parseHeaders(headers)

	assert.Equal(t, "v4", signature)
	assert.Equal(t, "us-east-1", region)
	assert.Equal(t, "host", signedHeaders)
}

func TestParseHeaders_ShouldHandleInvalidV4Format(t *testing.T) {
	headers := map[string][]string{
		"authorization": {"AWS4-HMAC-SHA256 InvalidFormat"},
	}

	signature, region, signedHeaders := parseHeaders(headers)

	assert.Equal(t, "v4", signature)
	assert.Equal(t, "", region)
	assert.Equal(t, "", signedHeaders)
}

// getIsObs Tests

func TestGetIsObs_ShouldReturnTrue_WhenNotTemporaryAndNoAmzHeaders(t *testing.T) {
	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"Cache-Control": {"no-cache"},
	}

	result := getIsObs(false, nil, headers)

	assert.True(t, result)
}

func TestGetIsObs_ShouldReturnFalse_WhenNotTemporaryAndHasAmzHeader(t *testing.T) {
	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"x-amz-meta-key": {"value"},
	}

	result := getIsObs(false, nil, headers)

	assert.False(t, result)
}

func TestGetIsObs_ShouldReturnTrue_WhenTemporaryAndNoAmzQueryOrAccessKey(t *testing.T) {
	querys := []string{
		"Expires=1234567890",
		"test-param=value",
	}

	result := getIsObs(true, querys, nil)

	assert.True(t, result)
}

func TestGetIsObs_ShouldReturnFalse_WhenTemporaryAndHasAmzQueryPrefix(t *testing.T) {
	querys := []string{
		"Expires=1234567890",
		"X-Amz-Date=20240101T000000Z",
	}

	result := getIsObs(true, querys, nil)

	assert.False(t, result)
}

func TestGetIsObs_ShouldReturnFalse_WhenTemporaryAndHasAWSAccessKeyId(t *testing.T) {
	querys := []string{
		"AWSAccessKeyId=AKIAIOSFODNN7EXAMPLE",
		"Expires=1234567890",
	}

	result := getIsObs(true, querys, nil)

	assert.False(t, result)
}

func TestGetIsObs_ShouldHandleCaseInsensitiveAmzPrefix(t *testing.T) {
	querys := []string{
		"x-amz-date=20240101T000000Z",
		"Expires=1234567890",
	}

	result := getIsObs(true, querys, nil)

	assert.False(t, result)
}

// isPathStyle Tests

func TestIsPathStyle_ShouldReturnTrue_WhenHostDoesNotStartWithBucket(t *testing.T) {
	headers := map[string][]string{
		"host": {"obs.example.com"},
	}

	result := isPathStyle(headers, "my-bucket")

	assert.True(t, result)
}

func TestIsPathStyle_ShouldReturnFalse_WhenHostStartsWithBucket(t *testing.T) {
	headers := map[string][]string{
		"host": {"my-bucket.obs.example.com"},
	}

	result := isPathStyle(headers, "my-bucket")

	assert.False(t, result)
}

func TestIsPathStyle_ShouldReturnFalse_WhenNoHostHeader(t *testing.T) {
	headers := map[string][]string{}

	result := isPathStyle(headers, "my-bucket")

	assert.False(t, result)
}

func TestIsPathStyle_ShouldReturnFalse_WhenHostHeaderIsEmpty(t *testing.T) {
	headers := map[string][]string{
		"host": {},
	}

	result := isPathStyle(headers, "my-bucket")

	assert.False(t, result)
}

func TestIsPathStyle_ShouldHandleCaseSensitiveHost(t *testing.T) {
	headers := map[string][]string{
		"host": {"My-Bucket.obs.example.com"},
	}

	result := isPathStyle(headers, "my-bucket")

	// Case matters - "My-Bucket" != "my-bucket"
	assert.True(t, result)
}

func TestIsPathStyle_ShouldHandleSubdomainWithDifferentPort(t *testing.T) {
	headers := map[string][]string{
		"host": {"my-bucket.obs.example.com:8080"},
	}

	result := isPathStyle(headers, "my-bucket")

	assert.False(t, result)
}

// getQuerysResult Tests

func TestGetQuerysResult_ShouldReturnFilteredQueries_WhenValidInput(t *testing.T) {
	querys := []string{
		"key1=value1",
		"key2=value2",
		"key3=value3",
	}

	result := getQuerysResult(querys)

	assert.Len(t, result, 3)
	assert.Contains(t, result, "key1=value1")
	assert.Contains(t, result, "key2=value2")
	assert.Contains(t, result, "key3=value3")
}

func TestGetQuerysResult_ShouldFilterEmptyStrings_WhenInputHasEmptyElements(t *testing.T) {
	querys := []string{
		"key1=value1",
		"",
		"key2=value2",
	}

	result := getQuerysResult(querys)

	assert.Len(t, result, 2)
	assert.Contains(t, result, "key1=value1")
	assert.Contains(t, result, "key2=value2")
}

func TestGetQuerysResult_ShouldFilterEqualSignOnly_WhenInputContainsOnlyEqual(t *testing.T) {
	querys := []string{
		"key1=value1",
		"=",
		"key2=value2",
	}

	result := getQuerysResult(querys)

	assert.Len(t, result, 2)
	assert.Contains(t, result, "key1=value1")
	assert.Contains(t, result, "key2=value2")
}

func TestGetQuerysResult_ShouldReturnEmpty_WhenInputHasOnlyInvalidElements(t *testing.T) {
	querys := []string{
		"",
		"=",
	}

	result := getQuerysResult(querys)

	assert.Len(t, result, 0)
}

func TestGetQuerysResult_ShouldReturnEmpty_WhenInputIsEmpty(t *testing.T) {
	querys := []string{}

	result := getQuerysResult(querys)

	assert.Len(t, result, 0)
}

func TestGetQuerysResult_ShouldHandleMultipleEqualSigns(t *testing.T) {
	querys := []string{
		"key1=value1=value2",
	}

	result := getQuerysResult(querys)

	assert.Len(t, result, 1)
	assert.Equal(t, "key1=value1=value2", result[0])
}

// getParams Tests

func TestGetParams_ShouldReturnParamsMap_WhenValidInput(t *testing.T) {
	querysResult := []string{
		"key1=value1",
		"key2=value2",
		"key3=value3",
	}

	result := getParams(querysResult)

	assert.Len(t, result, 3)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])
	assert.Equal(t, "value3", result["key3"])
}

func TestGetParams_ShouldHandleKeyWithoutValue_WhenInputHasKeyOnly(t *testing.T) {
	querysResult := []string{
		"key1=value1",
		"key2",
		"key3=value3",
	}

	result := getParams(querysResult)

	assert.Len(t, result, 3)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "", result["key2"])
	assert.Equal(t, "value3", result["key3"])
}

func TestGetParams_ShouldHandleMultipleEqualSigns_WhenValueContainsEqual(t *testing.T) {
	querysResult := []string{
		"key1=value1=value2=value3",
	}

	result := getParams(querysResult)

	assert.Equal(t, "value1=value2=value3", result["key1"])
}

func TestGetParams_ShouldReturnEmpty_WhenInputIsEmpty(t *testing.T) {
	querysResult := []string{}

	result := getParams(querysResult)

	assert.Len(t, result, 0)
}

func TestGetParams_ShouldDecodeURL_WhenKeyOrValueIsEncoded(t *testing.T) {
	querysResult := []string{
		"key1=value%201",
		"key%202=value2",
		"key%203=value%203",
	}

	result := getParams(querysResult)

	assert.Equal(t, "value 1", result["key1"])
	assert.Equal(t, "value2", result["key 2"])
	assert.Equal(t, "value 3", result["key 3"])
}

func TestGetParams_ShouldHandleDuplicateKeys_WhenSameKeyAppearsMultipleTimes(t *testing.T) {
	querysResult := []string{
		"key1=value1",
		"key2=value2",
		"key1=value3",
	}

	result := getParams(querysResult)

	// The last value wins for duplicate keys
	assert.Equal(t, "value3", result["key1"])
	assert.Equal(t, "value2", result["key2"])
}

func TestGetParams_ShouldHandleEmptyValue_WhenKeyFollowedByEqual(t *testing.T) {
	querysResult := []string{
		"key1=",
	}

	result := getParams(querysResult)

	assert.Equal(t, "", result["key1"])
}

// GetContentType Tests

func TestGetContentType_ShouldReturnContentType_WhenKnownExtension(t *testing.T) {
	contentType, ok := GetContentType("test.txt")
	assert.True(t, ok)
	assert.Equal(t, "text/plain", contentType)

	contentType, ok = GetContentType("image.jpg")
	assert.True(t, ok)
	assert.Equal(t, "image/jpeg", contentType)

	contentType, ok = GetContentType("document.pdf")
	assert.True(t, ok)
	assert.Equal(t, "application/pdf", contentType)
}

func TestGetContentType_ShouldReturnContentType_WhenLowercaseExtension(t *testing.T) {
	contentType, ok := GetContentType("test.JPEG")
	assert.True(t, ok)
	assert.Equal(t, "image/jpeg", contentType)

	contentType, ok = GetContentType("test.PNG")
	assert.True(t, ok)
	assert.Equal(t, "image/png", contentType)
}

func TestGetContentType_ShouldReturnNotFound_WhenUnknownExtension(t *testing.T) {
	contentType, ok := GetContentType("test.xyz")
	assert.False(t, ok)
	assert.Equal(t, "", contentType)
}

func TestGetContentType_ShouldReturnNotFound_WhenNoExtension(t *testing.T) {
	contentType, ok := GetContentType("test")
	assert.False(t, ok)
	assert.Equal(t, "", contentType)

	contentType, ok = GetContentType("test.")
	assert.False(t, ok)
	assert.Equal(t, "", contentType)
}

func TestGetContentType_ShouldHandleComplexFilenames(t *testing.T) {
	contentType, ok := GetContentType("/path/to/file.tar.gz")
	assert.True(t, ok)
	assert.Equal(t, "application/gzip", contentType)

	contentType, ok = GetContentType("document-v2.docx")
	assert.True(t, ok)
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", contentType)
}

func TestGetContentType_ShouldHandleMultipleDotsInFilename(t *testing.T) {
	contentType, ok := GetContentType("file.name.with.many.dots.html")
	assert.True(t, ok)
	assert.Equal(t, "text/html", contentType)
}

// ObsClient.getContentType Tests

func TestObsClientGetContentType_ShouldReturnKeyContentType_WhenKeyHasExtension(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			ObjectOperationInput: ObjectOperationInput{
				Key: "test.txt",
			},
		},
	}

	client := ObsClient{}
	contentType := client.getContentType(input, "")

	assert.Equal(t, "text/plain", contentType)
}

func TestObsClientGetContentType_ShouldReturnSourceFileContentType_WhenKeyHasNoExtension(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			ObjectOperationInput: ObjectOperationInput{
				Key: "file",
			},
		},
	}

	client := ObsClient{}
	contentType := client.getContentType(input, "document.pdf")

	assert.Equal(t, "application/pdf", contentType)
}

func TestObsClientGetContentType_ShouldReturnKeyContentType_WhenBothHaveExtensions(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			ObjectOperationInput: ObjectOperationInput{
				Key: "image.jpg",
			},
		},
	}

	client := ObsClient{}
	contentType := client.getContentType(input, "file.txt")

	// Key extension takes precedence
	assert.Equal(t, "image/jpeg", contentType)
}

func TestObsClientGetContentType_ShouldReturnEmpty_WhenNeitherHasExtension(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			ObjectOperationInput: ObjectOperationInput{
				Key: "file",
			},
		},
	}

	client := ObsClient{}
	contentType := client.getContentType(input, "document")

	assert.Equal(t, "", contentType)
}

func TestObsClientGetContentType_ShouldReturnEmpty_WhenBothUnknownExtension(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			ObjectOperationInput: ObjectOperationInput{
				Key: "file.xyz",
			},
		},
	}

	client := ObsClient{}
	contentType := client.getContentType(input, "document.abc")

	assert.Equal(t, "", contentType)
}

// ObsClient.isGetContentType Tests

func TestObsClientIsGetContentType_ShouldReturnTrue_WhenContentTypeIsEmptyAndKeyIsNotEmpty(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			HttpHeader: HttpHeader{
				ContentType: "",
			},
			ObjectOperationInput: ObjectOperationInput{
				Key: "test.txt",
			},
		},
	}

	client := ObsClient{}
	result := client.isGetContentType(input)

	assert.True(t, result)
}

func TestObsClientIsGetContentType_ShouldReturnFalse_WhenContentTypeIsNotEmpty(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			HttpHeader: HttpHeader{
				ContentType: "application/json",
			},
			ObjectOperationInput: ObjectOperationInput{
				Key: "test.txt",
			},
		},
	}

	client := ObsClient{}
	result := client.isGetContentType(input)

	assert.False(t, result)
}

func TestObsClientIsGetContentType_ShouldReturnFalse_WhenKeyIsEmpty(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			HttpHeader: HttpHeader{
				ContentType: "",
			},
			ObjectOperationInput: ObjectOperationInput{
				Key: "",
			},
		},
	}

	client := ObsClient{}
	result := client.isGetContentType(input)

	assert.False(t, result)
}

func TestObsClientIsGetContentType_ShouldReturnFalse_WhenBothAreEmpty(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			HttpHeader: HttpHeader{
				ContentType: "",
			},
			ObjectOperationInput: ObjectOperationInput{
				Key: "",
			},
		},
	}

	client := ObsClient{}
	result := client.isGetContentType(input)

	assert.False(t, result)
}
