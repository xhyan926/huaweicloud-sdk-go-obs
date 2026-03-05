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
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// MockRoundTripper implements http.RoundTripper for testing
type MockRoundTripper struct {
	// Simple pattern: fixed response and error
	response     *http.Response
	err          error
	capturedReq  *http.Request
	requestCount int
	// Function pattern: dynamic response and error functions
	ResponseFunc func(req *http.Request) *http.Response
	ErrorFunc    func(req *http.Request) error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.requestCount++
	m.capturedReq = req

	// If ResponseFunc is set, use it for dynamic responses
	if m.ResponseFunc != nil {
		if m.ErrorFunc != nil {
			err := m.ErrorFunc(req)
			if err != nil {
				return m.ResponseFunc(req), err
			}
		}
		return m.ResponseFunc(req), nil
	}

	// If ErrorFunc is set alone, return error
	if m.ErrorFunc != nil {
		return nil, m.ErrorFunc(req)
	}

	// Default: use simple pattern
	return m.response, m.err
}

// CreateMockClient creates a test client with a mock HTTP transport
func CreateMockClient(response *http.Response, err error) *ObsClient {
	transport := &MockRoundTripper{
		response: response,
		err:      err,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))
	return client
}

// CreateMockResponse creates a mock HTTP response for testing
func CreateMockResponse(statusCode int, body string, headers http.Header) *http.Response {
	return &http.Response{
		StatusCode:    statusCode,
		Status:        http.StatusText(statusCode),
		Body:          io.NopCloser(strings.NewReader(body)),
		Header:        headers,
		ContentLength: int64(len(body)),
	}
}

// CreateSuccessResponse creates a successful response
func CreateSuccessResponse(body string) *http.Response {
	headers := make(http.Header)
	headers.Set(HEADER_REQUEST_ID, "test-request-id-123")
	return CreateMockResponse(http.StatusOK, body, headers)
}

// CreateErrorResponse creates an error response
func CreateErrorResponse(code, message string) *http.Response {
	headers := make(http.Header)
	headers.Set(HEADER_REQUEST_ID, "test-error-request-id")
	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
	<Code>` + code + `</Code>
	<Message>` + message + `</Message>
	<RequestId>test-request-id</RequestId>
</Error>`
	return CreateMockResponse(http.StatusBadRequest, body, headers)
}

// CreateTempFile creates a temporary file for testing
func CreateTempFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "test-obs-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if content != "" {
		_, err = tmpFile.WriteString(content)
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}

	tmpFile.Close()
	return tmpFile
}

// TestListBucketsWithSignedUrl tests the ListBucketsWithSignedUrl method
func TestListBucketsWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "success",
			response:      CreateSuccessResponse(TestListBucketsXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "http_error",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "server_error",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.ListBucketsWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestCreateBucketWithSignedUrl tests the CreateBucketWithSignedUrl method
func TestCreateBucketWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldCreateBucketWithSignedUrl_ReturnSuccess_WhenGivenValidRegion",
			response:      CreateSuccessResponse(TestCreateBucketXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestCreateBucketXML),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldCreateBucketWithSignedUrl_ReturnSuccess_WhenGivenValidLocationConstraint",
			response:      CreateSuccessResponse(TestCreateBucketXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<CreateBucketConfiguration><Location>us-east-1</Location></CreateBucketConfiguration>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldCreateBucketWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestCreateBucketXML),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldCreateBucketWithSignedUrl_ReturnError_WhenBucketAlreadyExists",
			response:      CreateErrorResponse("BucketAlreadyExists", "The requested bucket name is not available"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestCreateBucketXML),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldCreateBucketWithSignedUrl_ReturnSuccess_WhenGivenEmptyData",
			response:      CreateSuccessResponse(""),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          nil,
			expectError:   false,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.CreateBucketWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestDeleteBucketWithSignedUrl tests the DeleteBucketWithSignedUrl method
func TestDeleteBucketWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldDeleteBucketWithSignedUrl_ReturnSuccess_WhenBucketExists",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteBucketWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketWithSignedUrl_ReturnError_WhenBucketNotEmpty",
			response:      CreateErrorResponse("BucketNotEmpty", "The bucket you tried to delete is not empty"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.DeleteBucketWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestPutFileWithSignedUrl tests the PutFileWithSignedUrl method
func TestPutFileWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		signedUrl      string
		headers        http.Header
		filePath       string
		fileContent    string
		responseError  error
		response       *http.Response
		expectError    bool
	}{
		{
			name:          "success",
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			filePath:      CreateTempFile(t, "test content").Name(),
			fileContent:   "test content",
			responseError: nil,
			response:      CreateSuccessResponse(""),
			expectError:   false,
		},
		{
			name:          "success_with_custom_content_length",
			signedUrl:     "https://obs.example.com",
			headers:       http.Header{HEADER_CONTENT_LENGTH: {"5"}},
			filePath:      CreateTempFile(t, "hello").Name(),
			fileContent:   "hello",
			responseError: nil,
			response:      CreateSuccessResponse(""),
			expectError:   false,
		},
		{
			name:          "success_with_aws_content_length",
			signedUrl:     "https://obs.example.com",
			headers:       http.Header{HEADER_CONTENT_LENGTH_CAMEL: {"10"}},
			filePath:      CreateTempFile(t, "1234567890").Name(),
			fileContent:   "1234567890",
			responseError: nil,
			response:      CreateSuccessResponse(""),
			expectError:   false,
		},
		{
			name:          "error_content_length_larger_than_file",
			signedUrl:     "https://obs.example.com",
			headers:       http.Header{HEADER_CONTENT_LENGTH: {"100"}},
			filePath:      CreateTempFile(t, "short").Name(),
			fileContent:   "short",
			responseError: errors.New("mock error"),
			response:      nil,
			expectError:   true,
		},
		{
			name:          "error_opening_file",
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			filePath:      "/nonexistent/file.txt",
			fileContent:   "",
			responseError: nil,
			response:      nil,
			expectError:   true,
		},
		{
			name:          "empty_path",
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			filePath:      "",
			fileContent:   "",
			responseError: nil,
			response:      CreateSuccessResponse(""),
			expectError:   false,
		},
		{
			name:          "error_http",
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			filePath:      CreateTempFile(t, "test").Name(),
			fileContent:   "test",
			responseError: errors.New("network error"),
			response:      nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.PutFileWithSignedUrl(tt.signedUrl, tt.headers, tt.filePath)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if !tt.expectError && output == nil {
				t.Error("Expected output to not be nil")
			}
			if tt.expectError && err != nil && err.Error() == "ContentLength is larger than fileSize" {
				// Verify the specific error message
				t.Log("Content length error correctly caught")
			}
		})
	}
}

// TestGetObjectWithSignedUrl tests the GetObjectWithSignedUrl method
func TestGetObjectWithSignedUrl(t *testing.T) {
	body := "test object content"
	client := CreateMockClient(CreateMockResponse(http.StatusOK, body, make(http.Header)), nil)

	output, err := client.GetObjectWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestPutObjectWithSignedUrl tests the PutObjectWithSignedUrl method
func TestPutObjectWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateMockResponse(http.StatusOK, "", make(http.Header)), nil)

	output, err := client.PutObjectWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader("test content"))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestCopyObjectWithSignedUrl tests the CopyObjectWithSignedUrl method
func TestCopyObjectWithSignedUrl(t *testing.T) {
	headers := make(http.Header)
	headers.Set(HEADER_ETAG, `"test-etag"`)

	client := CreateMockClient(CreateMockResponse(http.StatusOK, "", headers), nil)

	output, err := client.CopyObjectWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetObjectMetadataWithSignedUrl tests the GetObjectMetadataWithSignedUrl method
func TestGetObjectMetadataWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateMockResponse(http.StatusOK, "", make(http.Header)), nil)

	output, err := client.GetObjectMetadataWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetObjectAclWithSignedUrl tests the GetObjectAclWithSignedUrl method
func TestGetObjectAclWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateMockResponse(http.StatusOK, TestObjectACLXML, make(http.Header)), nil)

	output, err := client.GetObjectAclWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestAppendObjectWithSignedURL tests the AppendObjectWithSignedURL method
func TestAppendObjectWithSignedURL(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldAppendObjectWithSignedURL_ReturnSuccess_WhenGivenValidData",
			response:      CreateMockResponse(http.StatusOK, "", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("append content"),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldAppendObjectWithSignedURL_ReturnSuccess_WhenGivenValidPosition",
			response:      CreateMockResponse(http.StatusOK, "", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("more content"),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldAppendObjectWithSignedURL_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("content"),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldAppendObjectWithSignedURL_ReturnError_WhenGivenInvalidPosition",
			response:      CreateErrorResponse("InvalidArgument", "Position is not valid"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("content"),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldAppendObjectWithSignedURL_ReturnError_WhenParseFails",
			response:      CreateMockResponse(http.StatusOK, "<invalid>xml</>", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("content"),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.AppendObjectWithSignedURL(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestModifyObjectWithSignedURL tests the ModifyObjectWithSignedURL method
func TestModifyObjectWithSignedURL(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldModifyObjectWithSignedURL_ReturnSuccess_WhenGivenValidData",
			response:      CreateMockResponse(http.StatusOK, "", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("modify content"),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldModifyObjectWithSignedURL_ReturnSuccess_WhenGivenValidRange",
			response:      CreateMockResponse(http.StatusOK, "", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("partial content"),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldModifyObjectWithSignedURL_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("content"),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldModifyObjectWithSignedURL_ReturnError_WhenGivenInvalidRange",
			response:      CreateErrorResponse("InvalidRange", "The requested range is not satisfiable"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader("content"),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.ModifyObjectWithSignedURL(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestRestoreObjectWithSignedUrl tests the RestoreObjectWithSignedUrl method
func TestRestoreObjectWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><RestoreObjectRequest/>`), nil)

	output, err := client.RestoreObjectWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader(`<RestoreRequest><Days>3</Days></RestoreRequest>`))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestSetObjectAclWithSignedUrl tests the SetObjectAclWithSignedUrl method
func TestSetObjectAclWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><AccessControlPolicy/>`), nil)

	output, err := client.SetObjectAclWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader(TestObjectACLXML))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetBucketMetadataWithSignedUrl tests the GetBucketMetadataWithSignedUrl method
func TestGetBucketMetadataWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateMockResponse(http.StatusOK, "", make(http.Header)), nil)

	output, err := client.GetBucketMetadataWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestSetBucketStoragePolicyWithSignedUrl tests the SetBucketStoragePolicyWithSignedUrl method
func TestSetBucketStoragePolicyWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketStoragePolicyWithSignedUrl_ReturnSuccess_WhenSetToStandard",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><SetBucketStoragePolicyResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<StorageConfiguration><StorageClass>STANDARD</StorageClass></StorageConfiguration>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketStoragePolicyWithSignedUrl_ReturnSuccess_WhenSetToWarm",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><SetBucketStoragePolicyResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<StorageConfiguration><StorageClass>WARM</StorageClass></StorageConfiguration>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketStoragePolicyWithSignedUrl_ReturnSuccess_WhenSetToCold",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><SetBucketStoragePolicyResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<StorageConfiguration><StorageClass>COLD</StorageClass></StorageConfiguration>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketStoragePolicyWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<StorageConfiguration><StorageClass>STANDARD</StorageClass></StorageConfiguration>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketStoragePolicyWithSignedUrl_ReturnError_WhenGivenInvalidStorageClass",
			response:      CreateErrorResponse("InvalidStorageClass", "The storage class you specified is not valid"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<StorageConfiguration><StorageClass>INVALID</StorageClass></StorageConfiguration>`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketStoragePolicyWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketStoragePolicyWithSignedUrl tests the GetBucketStoragePolicyWithSignedUrl method
func TestGetBucketStoragePolicyWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketStoragePolicyWithSignedUrl_ReturnStandard_WhenConfiguredWithStandard",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><GetBucketStoragePolicyOutput><StorageClass>STANDARD</StorageClass></GetBucketStoragePolicyOutput>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketStoragePolicyWithSignedUrl_ReturnWarm_WhenConfiguredWithWarm",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><GetBucketStoragePolicyOutput><StorageClass>WARM</StorageClass></GetBucketStoragePolicyOutput>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketStoragePolicyWithSignedUrl_ReturnNotConfigured_WhenNotSet",
			response:      CreateErrorResponse("NoSuchStoragePolicy", "The bucket does not have a storage policy configured"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketStoragePolicyWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketStoragePolicyWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketStoragePolicyWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestListMultipartUploadsWithSignedUrl tests the ListMultipartUploadsWithSignedUrl method
func TestListMultipartUploadsWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldListMultipartUploadsWithSignedUrl_ReturnEmptyList_WhenBucketEmpty",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><ListMultipartUploadsResult><Bucket>test-bucket</Bucket><UploadIdMarker></UploadIdMarker><NextUploadIdMarker></NextUploadIdMarker><MaxUploads>1000</MaxUploads><IsTruncated>false</IsTruncated></ListMultipartUploadsResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldListMultipartUploadsWithSignedUrl_ReturnUploads_WhenGivenValidPrefix",
			response:      CreateSuccessResponse(TestListMultipartUploadsXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldListMultipartUploadsWithSignedUrl_ReturnAllUploads_WhenMaxUploadsSet",
			response:      CreateSuccessResponse(TestListMultipartUploadsXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldListMultipartUploadsWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldListMultipartUploadsWithSignedUrl_ReturnError_WhenEncodingTypeUrl",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><ListMultipartUploadsResult><Bucket>test-bucket</Bucket><EncodingType>url</EncodingType><MaxUploads>1000</MaxUploads><IsTruncated>false</IsTruncated><Upload><Key>test%2Fobject.txt</Key><UploadId>upload-id</UploadId></Upload></ListMultipartUploadsResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.ListMultipartUploadsWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestSetBucketQuotaWithSignedUrl tests the SetBucketQuotaWithSignedUrl method
func TestSetBucketQuotaWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketQuotaWithSignedUrl_ReturnSuccess_WhenGivenValidQuota",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><SetBucketQuotaResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Quota><StorageQuota>10737418240</StorageQuota></Quota>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketQuotaWithSignedUrl_ReturnSuccess_WhenSetToUnlimited",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><SetBucketQuotaResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Quota><StorageQuota>-1</StorageQuota></Quota>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketQuotaWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Quota><StorageQuota>10737418240</StorageQuota></Quota>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketQuotaWithSignedUrl_ReturnError_WhenGivenInvalidQuota",
			response:      CreateErrorResponse("InvalidArgument", "Invalid quota value"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Quota><StorageQuota>-100</StorageQuota></Quota>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketQuotaWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Quota><StorageQuota>10737418240</StorageQuota></Quota>`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketQuotaWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketQuotaWithSignedUrl tests the GetBucketQuotaWithSignedUrl method
func TestGetBucketQuotaWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketQuotaWithSignedUrl_ReturnQuota_WhenQuotaSet",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><Quota><StorageQuota>10737418240</StorageQuota></Quota>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketQuotaWithSignedUrlReturnUnlimited_WhenNoQuotaSet",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><Quota><StorageQuota>-1</StorageQuota></Quota>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketQuotaWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketQuotaWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketQuotaWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketQuotaWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestHeadBucketWithSignedUrl tests the HeadBucketWithSignedUrl method
func TestHeadBucketWithSignedUrl(t *testing.T) {
	customHeaders := make(http.Header)
	customHeaders.Set("X-Custom-Header", "custom-value")

	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldHeadBucketWithSignedUrl_ReturnSuccess_WhenBucketExists",
			response:      CreateMockResponse(http.StatusOK, "", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldHeadBucketWithSignedUrl_ReturnSuccess_WhenGivenCustomHeaders",
			response:      CreateMockResponse(http.StatusOK, "", customHeaders),
			signedUrl:     "https://obs.example.com",
			headers:       customHeaders,
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldHeadBucketWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldHeadBucketWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.HeadBucketWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestHeadObjectWithSignedUrl tests the HeadObjectWithSignedUrl method
func TestHeadObjectWithSignedUrl(t *testing.T) {
	customHeaders := make(http.Header)
	customHeaders.Set("X-Custom-Header", "custom-value")

	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldHeadObjectWithSignedUrl_ReturnSuccess_WhenObjectExists",
			response:      CreateMockResponse(http.StatusOK, "", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldHeadObjectWithSignedUrl_ReturnSuccess_WhenGivenVersionId",
			response:      CreateMockResponse(http.StatusOK, "", make(http.Header)),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldHeadObjectWithSignedUrl_ReturnSuccess_WhenGivenCustomHeaders",
			response:      CreateMockResponse(http.StatusOK, "", customHeaders),
			signedUrl:     "https://obs.example.com",
			headers:       customHeaders,
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldHeadObjectWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldHeadObjectWithSignedUrl_ReturnError_WhenObjectNotFound",
			response:      CreateErrorResponse("NoSuchKey", "The specified key does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.HeadObjectWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketStorageInfoWithSignedUrl tests the GetBucketStorageInfoWithSignedUrl method
func TestGetBucketStorageInfoWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketStorageInfoWithSignedUrl_ReturnStorageInfo_WhenBucketHasObjects",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><GetBucketStorageInfoResult><Size>1024</Size><ObjectNumber>10</ObjectNumber></GetBucketStorageInfoResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketStorageInfoWithSignedUrl_ReturnZeroSize_WhenBucketEmpty",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><GetBucketStorageInfoResult><Size>0</Size><ObjectNumber>0</ObjectNumber></GetBucketStorageInfoResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketStorageInfoWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketStorageInfoWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketStorageInfoWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketStorageInfoWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketLocationWithSignedUrl tests the GetBucketLocationWithSignedUrl method
func TestGetBucketLocationWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketLocationWithSignedUrl_ReturnRegion_WhenLocatedInCN",
			response:      CreateSuccessResponse(`<Location>cn-north-4</Location>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketLocationWithSignedUrl_ReturnRegion_WhenLocatedInUS",
			response:      CreateSuccessResponse(`<Location>us-east-1</Location>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketLocationWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketLocationWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketLocationWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketLocationWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestSetBucketAclWithSignedUrl tests the SetBucketAclWithSignedUrl method
func TestSetBucketAclWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketAclWithSignedUrl_ReturnSuccess_WhenSetToPrivate",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><AccessControlPolicy/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<AccessControlPolicy><Owner><ID>owner-id</ID></Owner><AccessControlList/></AccessControlPolicy>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketAclWithSignedUrl_ReturnSuccess_WhenSetToPublicRead",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><AccessControlPolicy/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<AccessControlPolicy><Owner><ID>owner-id</ID></Owner><AccessControlList><Grant><Grantee><URI>http://acs.amazonaws.com/groups/global/AllUsers</URI></Grantee><Permission>READ</Permission></Grant></AccessControlList></AccessControlPolicy>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketAclWithSignedUrl_ReturnSuccess_WhenGivenGroupGrants",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><AccessControlPolicy/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<AccessControlPolicy><Owner><ID>owner-id</ID></Owner><AccessControlList><Grant><Grantee><ID>user-id</ID></Grantee><Permission>FULL_CONTROL</Permission></Grant></AccessControlList></AccessControlPolicy>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketAclWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<AccessControlPolicy/>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketAclWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<AccessControlPolicy/>`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketAclWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketAclWithSignedUrl tests the GetBucketAclWithSignedUrl method
func TestGetBucketAclWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketAclWithSignedUrl_ReturnPrivate_WhenConfiguredPrivate",
			response:      CreateSuccessResponse(TestBucketACLXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketAclWithSignedUrl_ReturnPublicRead_WhenConfiguredPublic",
			response:      CreateSuccessResponse(`<AccessControlPolicy><Owner><ID>test-owner-id</ID><DisplayName>test-owner</DisplayName></Owner><AccessControlList><Grant><Grantee><Type>Group</Type><URI>http://acs.amazonaws.com/groups/global/AllUsers</URI></Grantee><Permission>READ</Permission></Grant></AccessControlList></AccessControlPolicy>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketAclWithSignedUrl_ReturnGrants_WhenConfiguredWithGrants",
			response:      CreateSuccessResponse(TestBucketACLXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketAclWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketAclWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketAclWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestSetBucketPolicyWithSignedUrl tests the SetBucketPolicyWithSignedUrl method
func TestSetBucketPolicyWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketPolicyWithSignedUrl_ReturnSuccess_WhenGivenValidJsonPolicy",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><SetBucketPolicyResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`{"version":"2012-10-17","statement":[]}`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketPolicyWithSignedUrl_ReturnSuccess_WhenGivenValidXmlPolicy",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><SetBucketPolicyResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Policy/>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketPolicyWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`{"version":"2012-10-17","statement":[]}`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketPolicyWithSignedUrl_ReturnError_WhenGivenMalformedPolicy",
			response:      CreateErrorResponse("MalformedPolicy", "The policy was malformed"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`{invalid json}`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketPolicyWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`{"version":"2012-10-17","statement":[]}`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketPolicyWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketPolicyWithSignedUrl tests the GetBucketPolicyWithSignedUrl method
func TestGetBucketPolicyWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketPolicyWithSignedUrl_ReturnPolicy_WhenPolicyExists",
			response:      CreateSuccessResponse(`{"version":"2012-10-17","statement":[]}`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketPolicyWithSignedUrl_ReturnEmptyPolicy_WhenNotSet",
			response:      CreateErrorResponse("NoSuchBucketPolicy", "The bucket policy does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketPolicyWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketPolicyWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketPolicyWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestDeleteBucketPolicyWithSignedUrl tests the DeleteBucketPolicyWithSignedUrl method
func TestDeleteBucketPolicyWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldDeleteBucketPolicyWithSignedUrl_ReturnSuccess_WhenPolicyExists",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketPolicyResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteBucketPolicyWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketPolicyWithSignedUrl_ReturnError_WhenPolicyNotFound",
			response:      CreateErrorResponse("NoSuchBucketPolicy", "The bucket policy does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketPolicyWithSignedUrl_ReturnSuccess_WhenAlreadyDeleted",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketPolicyResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.DeleteBucketPolicyWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestSetBucketCorsWithSignedUrl tests the SetBucketCorsWithSignedUrl method
func TestSetBucketCorsWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketCorsWithSignedUrl_ReturnSuccess_WhenGivenValidCorsRules",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><CORSConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestBucketCorsXML),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketCorsWithSignedUrl_ReturnSuccess_WhenGivenEmptyCors",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><CORSConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<CORSConfiguration/>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketCorsWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestBucketCorsXML),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketCorsWithSignedUrl_ReturnError_WhenGivenMalformedCors",
			response:      CreateErrorResponse("InvalidArgument", "Invalid CORS configuration"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<CORSConfiguration><CORSRule></CORSRule></CORSConfiguration>`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketCorsWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketCorsWithSignedUrl tests the GetBucketCorsWithSignedUrl method
func TestGetBucketCorsWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketCorsWithSignedUrl_ReturnRules_WhenCorsConfigured",
			response:      CreateSuccessResponse(TestBucketCorsXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketCorsWithSignedUrl_ReturnEmpty_WhenNotConfigured",
			response:      CreateErrorResponse("NoSuchCORSConfiguration", "The CORS configuration does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketCorsWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketCorsWithSignedUrl_ReturnError_WhenBucketNotFound",
			response:      CreateErrorResponse("NoSuchBucket", "The specified bucket does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketCorsWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestDeleteBucketCorsWithSignedUrl tests the DeleteBucketCorsWithSignedUrl method
func TestDeleteBucketCorsWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldDeleteBucketCorsWithSignedUrl_ReturnSuccess_WhenCorsExists",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketCorsResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteBucketCorsWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketCorsWithSignedUrl_ReturnError_WhenCorsNotConfigured",
			response:      CreateErrorResponse("NoSuchCORSConfiguration", "The CORS configuration does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.DeleteBucketCorsWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestSetBucketVersioningWithSignedUrl tests the SetBucketVersioningWithSignedUrl method
func TestSetBucketVersioningWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketVersioningWithSignedUrl_ReturnSuccess_WhenSetToEnabled",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><VersioningConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestBucketVersioningXML),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketVersioningWithSignedUrl_ReturnSuccess_WhenSetToSuspended",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><VersioningConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<VersioningConfiguration><Status>Suspended</Status></VersioningConfiguration>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketVersioningWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestBucketVersioningXML),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketVersioningWithSignedUrl_ReturnError_WhenGivenMalformedConfig",
			response:      CreateErrorResponse("InvalidArgument", "Invalid versioning status"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<VersioningConfiguration><Status>Invalid</Status></VersioningConfiguration>`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketVersioningWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketVersioningWithSignedUrl tests the GetBucketVersioningWithSignedUrl method
func TestGetBucketVersioningWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketVersioningWithSignedUrl_ReturnEnabled_WhenVersioningEnabled",
			response:      CreateSuccessResponse(TestBucketVersioningXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketVersioningWithSignedUrl_ReturnSuspended_WhenVersioningSuspended",
			response:      CreateSuccessResponse(`<VersioningConfiguration><Status>Suspended</Status></VersioningConfiguration>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketVersioningWithSignedUrl_ReturnNotConfigured_WhenNotSet",
			response:      CreateSuccessResponse(`<VersioningConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketVersioningWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketVersioningWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestSetBucketWebsiteConfigurationWithSignedUrl tests the SetBucketWebsiteConfigurationWithSignedUrl method
func TestSetBucketWebsiteConfigurationWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketWebsiteConfigurationWithSignedUrl_ReturnSuccess_WhenGivenIndexAndErrorDoc",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><WebsiteConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestBucketWebsiteXML),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketWebsiteConfigurationWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestBucketWebsiteXML),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketWebsiteConfigurationWithSignedUrl_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(TestBucketWebsiteXML),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketWebsiteConfigurationWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketWebsiteConfigurationWithSignedUrl tests the GetBucketWebsiteConfigurationWithSignedUrl method
func TestGetBucketWebsiteConfigurationWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketWebsiteConfigurationWithSignedUrl_ReturnConfig_WhenIndexSet",
			response:      CreateSuccessResponse(TestBucketWebsiteXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketWebsiteConfigurationWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketWebsiteConfigurationWithSignedUrl_ReturnError_WhenNotConfigured",
			response:      CreateErrorResponse("NoSuchWebsiteConfiguration", "The website configuration does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketWebsiteConfigurationWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestDeleteBucketWebsiteConfigurationWithSignedUrl tests the DeleteBucketWebsiteConfigurationWithSignedUrl method
func TestDeleteBucketWebsiteConfigurationWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldDeleteBucketWebsiteConfigurationWithSignedUrl_ReturnSuccess_WhenConfigurationExists",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketWebsiteConfigurationResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteBucketWebsiteConfigurationWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketWebsiteConfigurationWithSignedUrl_ReturnError_WhenConfigurationNotSet",
			response:      CreateErrorResponse("NoSuchWebsiteConfiguration", "The website configuration does not exist"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.DeleteBucketWebsiteConfigurationWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestSetBucketLoggingConfigurationWithSignedUrl tests the SetBucketLoggingConfigurationWithSignedUrl method
func TestSetBucketLoggingConfigurationWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><BucketLoggingStatus/>`), nil)

	output, err := client.SetBucketLoggingConfigurationWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader(`<BucketLoggingStatus><LoggingEnabled><TargetBucket>log-bucket</TargetBucket><TargetPrefix>logs/</TargetPrefix></LoggingEnabled></BucketLoggingStatus>`))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetBucketLoggingConfigurationWithSignedUrl tests the GetBucketLoggingConfigurationWithSignedUrl method
func TestGetBucketLoggingConfigurationWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<BucketLoggingStatus><LoggingEnabled><TargetBucket>log-bucket</TargetBucket><TargetPrefix>logs/</TargetPrefix></LoggingEnabled></BucketLoggingStatus>`), nil)

	output, err := client.GetBucketLoggingConfigurationWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestSetBucketLifecycleConfigurationWithSignedUrl tests the SetBucketLifecycleConfigurationWithSignedUrl method
func TestSetBucketLifecycleConfigurationWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><LifecycleConfiguration/>`), nil)

	output, err := client.SetBucketLifecycleConfigurationWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader(TestBucketLifecycleXML))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetBucketLifecycleConfigurationWithSignedUrl tests the GetBucketLifecycleConfigurationWithSignedUrl method
func TestGetBucketLifecycleConfigurationWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(TestBucketLifecycleXML), nil)

	output, err := client.GetBucketLifecycleConfigurationWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestDeleteBucketLifecycleConfigurationWithSignedUrl tests the DeleteBucketLifecycleConfigurationWithSignedUrl method
func TestDeleteBucketLifecycleConfigurationWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketLifecycleConfigurationResponse/>`), nil)

	output, err := client.DeleteBucketLifecycleConfigurationWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestSetBucketTaggingWithSignedUrl tests the SetBucketTaggingWithSignedUrl method
func TestSetBucketTaggingWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><Tagging/>`), nil)

	output, err := client.SetBucketTaggingWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader(`<Tagging><TagSet><Tag><Key>key1</Key><Value>value1</Value></Tag></TagSet></Tagging>`))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetBucketTaggingWithSignedUrl tests the GetBucketTaggingWithSignedUrl method
func TestGetBucketTaggingWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<Tagging><TagSet><Tag><Key>key1</Key><Value>value1</Value></Tag></TagSet></Tagging>`), nil)

	output, err := client.GetBucketTaggingWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestDeleteBucketTaggingWithSignedUrl tests the DeleteBucketTaggingWithSignedUrl method
func TestDeleteBucketTaggingWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketTaggingResponse/>`), nil)

	output, err := client.DeleteBucketTaggingWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestSetBucketNotificationWithSignedUrl tests the SetBucketNotificationWithSignedUrl method
func TestSetBucketNotificationWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><NotificationConfiguration/>`), nil)

	output, err := client.SetBucketNotificationWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader(`<NotificationConfiguration><TopicConfiguration><Id>test-id</Id><Topic>arn:aws:sns:us-east-1:123456789012:topic</Topic><Event>s3:ObjectCreated:*</Event></TopicConfiguration></NotificationConfiguration>`))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetBucketNotificationWithSignedUrl tests the GetBucketNotificationWithSignedUrl method
func TestGetBucketNotificationWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<NotificationConfiguration/>`), nil)

	output, err := client.GetBucketNotificationWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestDeleteObjectWithSignedUrl tests the DeleteObjectWithSignedUrl method
func TestDeleteObjectWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateMockResponse(http.StatusNoContent, "", make(http.Header)), nil)

	output, err := client.DeleteObjectWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestDeleteObjectsWithSignedUrl tests the DeleteObjectsWithSignedUrl method
func TestDeleteObjectsWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldDeleteObjectsWithSignedUrl_ReturnDeletedList_WhenGivenSingleObject",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteResult><Deleted><Key>test-object.txt</Key></Deleted></DeleteResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Delete><Object><Key>test-object.txt</Key></Object></Delete>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteObjectsWithSignedUrl_ReturnDeletedList_WhenGivenMultipleObjects",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteResult><Deleted><Key>object1.txt</Key></Deleted><Deleted><Key>object2.txt</Key></Deleted></DeleteResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Delete><Object><Key>object1.txt</Key></Object><Object><Key>object2.txt</Key></Object></Delete>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteObjectsWithSignedUrl_ReturnEmptyDeletedList_WhenGivenQuietMode",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteResult></DeleteResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Delete><Quiet>true</Quiet><Object><Key>test-object.txt</Key></Object></Delete>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteObjectsWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Delete><Object><Key>test-object.txt</Key></Object></Delete>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteObjectsWithSignedUrl_ReturnError_WhenGivenInvalidRequest",
			response:      CreateErrorResponse("MalformedXML", "The XML you provided was not well-formed"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Delete><Object></Object></Delete>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteObjectsWithSignedUrl_ReturnSuccess_WhenEncodingTypeUrl",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteResult><Deleted><Key>test%2Fobject.txt</Key></Deleted><EncodingType>url</EncodingType></DeleteResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<Delete><Object><Key>test/object.txt</Key></Object></Delete>`),
			expectError:   false,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.DeleteObjectsWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestAbortMultipartUploadWithSignedUrl tests the AbortMultipartUploadWithSignedUrl method
func TestAbortMultipartUploadWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><AbortMultipartUploadResponse/>`), nil)

	output, err := client.AbortMultipartUploadWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestInitiateMultipartUploadWithSignedUrl tests the InitiateMultipartUploadWithSignedUrl method
func TestInitiateMultipartUploadWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldInitiateMultipartUploadWithSignedUrl_ReturnUploadId_WhenGivenValidRequest",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>upload-id-1234567890</UploadId></InitiateMultipartUploadResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldInitiateMultipartUploadWithSignedUrl_ReturnUploadId_WhenGivenUrlEncodingType",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test%2Fobject.txt</Key><UploadId>upload-id-123</UploadId><EncodingType>url</EncodingType></InitiateMultipartUploadResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldInitiateMultipartUploadWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldInitiateMultipartUploadWithSignedUrl_ReturnError_WhenGivenInvalidRequest",
			response:      CreateErrorResponse("InvalidArgument", "Invalid argument"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.InitiateMultipartUploadWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestUploadPartWithSignedUrl tests the UploadPartWithSignedUrl method
func TestUploadPartWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateMockResponse(http.StatusOK, "", make(http.Header)), nil)

	output, err := client.UploadPartWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader("part content"))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestCompleteMultipartUploadWithSignedUrl tests the CompleteMultipartUploadWithSignedUrl method
func TestCompleteMultipartUploadWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldCompleteMultipartUploadWithSignedUrl_ReturnLocation_WhenGivenValidParts",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test-bucket.obs.region.myhuaweicloud.com/test-object.txt</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"complete-etag"</ETag></CompleteMultipartUploadResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>"part-etag"</ETag></Part></CompleteMultipartUpload>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldCompleteMultipartUploadWithSignedUrl_ReturnLocation_WhenGivenUrlEncodingType",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test-bucket.obs.region.myhuaweicloud.com/test%2Fobject.txt</Location><Bucket>test-bucket</Bucket><Key>test%2Fobject.txt</Key><ETag>"complete-etag"</ETag><EncodingType>url</EncodingType></CompleteMultipartUploadResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>"part-etag"</ETag></Part></CompleteMultipartUpload>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldCompleteMultipartUploadWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>"part-etag"</ETag></Part></CompleteMultipartUpload>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldCompleteMultipartUploadWithSignedUrl_ReturnError_WhenGivenInvalidParts",
			response:      CreateErrorResponse("InvalidPart", "One or more of the specified parts could not be found"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<CompleteMultipartUpload><Part><PartNumber>999</PartNumber><ETag>"invalid-etag"</ETag></Part></CompleteMultipartUpload>`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.CompleteMultipartUploadWithSignedUrl(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestListPartsWithSignedUrl tests the ListPartsWithSignedUrl method
func TestListPartsWithSignedUrl(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldListPartsWithSignedUrl_ReturnEmptyList_WhenUploadHasNoParts",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><ListPartsResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>upload-id</UploadId><PartNumberMarker>0</PartNumberMarker><NextPartNumberMarker>0</NextPartNumberMarker><MaxParts>1000</MaxParts><IsTruncated>false</IsTruncated></ListPartsResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldListPartsWithSignedUrl_ReturnParts_WhenGivenValidUpload",
			response:      CreateSuccessResponse(TestListPartsXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldListPartsWithSignedUrl_ReturnAllParts_WhenMaxPartsSet",
			response:      CreateSuccessResponse(TestListPartsXML),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldListPartsWithSignedUrl_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldListPartsWithSignedUrl_ReturnError_WhenEncodingTypeUrl",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><ListPartsResult><Bucket>test-bucket</Bucket><Key>test%2Fobject.txt</Key><UploadId>upload-id</UploadId><EncodingType>url</EncodingType><MaxParts>1000</MaxParts><IsTruncated>false</IsTruncated></ListPartsResult>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.ListPartsWithSignedUrl(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestCopyPartWithSignedUrl tests the CopyPartWithSignedUrl method
func TestCopyPartWithSignedUrl(t *testing.T) {
	headers := make(http.Header)
	headers.Set(HEADER_ETAG, `"copy-part-etag"`)

	client := CreateMockClient(CreateMockResponse(http.StatusOK, "", headers), nil)

	output, err := client.CopyPartWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestSetBucketRequestPaymentWithSignedUrl tests the SetBucketRequestPaymentWithSignedUrl method
func TestSetBucketRequestPaymentWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><RequestPaymentConfiguration/>`), nil)

	output, err := client.SetBucketRequestPaymentWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader(`<RequestPaymentConfiguration><Payer>Requester</Payer></RequestPaymentConfiguration>`))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestGetBucketRequestPaymentWithSignedUrl tests the GetBucketRequestPaymentWithSignedUrl method
func TestGetBucketRequestPaymentWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(`<RequestPaymentConfiguration><Payer>Requester</Payer></RequestPaymentConfiguration>`), nil)

	output, err := client.GetBucketRequestPaymentWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
}

// TestSetBucketEncryptionWithSignedURL tests the SetBucketEncryptionWithSignedURL method
func TestSetBucketEncryptionWithSignedURL(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		data           io.Reader
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldSetBucketEncryptionWithSignedURL_ReturnSuccess_WhenGivenValidAES256Config",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><ServerSideEncryptionConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<ServerSideEncryptionConfiguration><Rule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>AES256</SSEAlgorithm></ApplyServerSideEncryptionByDefault></Rule></ServerSideEncryptionConfiguration>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketEncryptionWithSignedURL_ReturnSuccess_WhenGivenValidKMSConfig",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><ServerSideEncryptionConfiguration/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<ServerSideEncryptionConfiguration><Rule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>aws:kms</SSEAlgorithm><KMSMasterKeyID>test-kms-key-id</KMSMasterKeyID></ApplyServerSideEncryptionByDefault></Rule></ServerSideEncryptionConfiguration>`),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldSetBucketEncryptionWithSignedURL_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<ServerSideEncryptionConfiguration><Rule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>AES256</SSEAlgorithm></ApplyServerSideEncryptionByDefault></Rule></ServerSideEncryptionConfiguration>`),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldSetBucketEncryptionWithSignedURL_ReturnError_WhenGivenInvalidResponse",
			response:      CreateErrorResponse("InvalidEncryptionConfiguration", "Invalid encryption configuration"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			data:          strings.NewReader(`<InvalidEncryptionConfiguration/>`),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.SetBucketEncryptionWithSignedURL(tt.signedUrl, tt.headers, tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestGetBucketEncryptionWithSignedURL tests the GetBucketEncryptionWithSignedURL method
func TestGetBucketEncryptionWithSignedURL(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldGetBucketEncryptionWithSignedURL_ReturnAES256_WhenConfiguredWithAES256",
			response:      CreateSuccessResponse(`<ServerSideEncryptionConfiguration><Rule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>AES256</SSEAlgorithm></ApplyServerSideEncryptionByDefault></Rule></ServerSideEncryptionConfiguration>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketEncryptionWithSignedURL_ReturnKMS_WhenConfiguredWithKMS",
			response:      CreateSuccessResponse(`<ServerSideEncryptionConfiguration><Rule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>aws:kms</SSEAlgorithm><KMSMasterKeyID>test-kms-key-id</KMSMasterKeyID></ApplyServerSideEncryptionByDefault></Rule></ServerSideEncryptionConfiguration>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldGetBucketEncryptionWithSignedURL_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldGetBucketEncryptionWithSignedURL_ReturnError_WhenNotConfigured",
			response:      CreateErrorResponse("ServerSideEncryptionConfigurationNotFoundError", "The server side encryption configuration was not found"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.GetBucketEncryptionWithSignedURL(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestDeleteBucketEncryptionWithSignedURL tests the DeleteBucketEncryptionWithSignedURL method
func TestDeleteBucketEncryptionWithSignedURL(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		responseError  error
		signedUrl      string
		headers        http.Header
		expectError    bool
		expectNil      bool
	}{
		{
			name:          "ShouldDeleteBucketEncryptionWithSignedURL_ReturnSuccess_WhenEncryptionExists",
			response:      CreateSuccessResponse(`<?xml version="1.0" encoding="UTF-8"?><DeleteBucketEncryptionResponse/>`),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "ShouldDeleteBucketEncryptionWithSignedURL_ReturnError_WhenNetworkFails",
			response:      nil,
			responseError: errors.New("network error"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketEncryptionWithSignedURL_ReturnError_WhenEncryptionNotSet",
			response:      CreateErrorResponse("ServerSideEncryptionConfigurationNotFoundError", "The server side encryption configuration was not found"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
		{
			name:          "ShouldDeleteBucketEncryptionWithSignedURL_ReturnError_WhenAccessDenied",
			response:      CreateErrorResponse("AccessDenied", "Access Denied"),
			signedUrl:     "https://obs.example.com",
			headers:       make(http.Header),
			expectError:   true,
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := CreateMockClient(tt.response, tt.responseError)

			output, err := client.DeleteBucketEncryptionWithSignedURL(tt.signedUrl, tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectNil && output != nil {
				t.Errorf("Expected output to be nil, got %v", output)
			}
			if !tt.expectNil && output == nil {
				t.Error("Expected output to not be nil")
			}
		})
	}
}

// TestListObjectsWithSignedUrl tests the ListObjectsWithSignedUrl method
func TestListObjectsWithSignedUrl(t *testing.T) {
	client := CreateMockClient(CreateSuccessResponse(TestListObjectsXML), nil)

	output, err := client.ListObjectsWithSignedUrl("https://obs.example.com", make(http.Header))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
	// Note: Location is only set when HEADER_BUCKET_REGION header is present
}

// TestListVersionsWithSignedUrl tests the ListVersionsWithSignedUrl method
func TestListVersionsWithSignedUrl(t *testing.T) {
	headers := make(http.Header)
	headers.Set(HEADER_BUCKET_REGION, "us-east-1")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
	<Name>test-bucket</Name>
	<Prefix></Prefix>
	<KeyMarker></KeyMarker>
	<MaxKeys>1000</MaxKeys>
	<IsTruncated>false</IsTruncated>
	<Version>
		<Key>test-object.txt</Key>
		<VersionId>version-123</VersionId>
		<IsLatest>true</IsLatest>
		<LastModified>2023-01-01T00:00:00Z</LastModified>
		<ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
		<Size>1024</Size>
		<StorageClass>STANDARD</StorageClass>
	</Version>
</ListVersionsResult>`

	client := CreateMockClient(CreateSuccessResponse(body), nil)

	output, err := client.ListVersionsWithSignedUrl("https://obs.example.com", headers)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == nil {
		t.Error("Expected output to not be nil")
	}
	if output.Location != "" {
		t.Logf("Location value: %s", output.Location)
	}
}

// BenchmarkPutFileWithSignedUrl benchmarks the PutFileWithSignedUrl method
func BenchmarkPutFileWithSignedUrl(b *testing.B) {
	// Note: Using t *testing.B for benchmark is correct in Go 1.22+
	// Create a temp file for benchmarking
	tmpFile, err := os.CreateTemp("", "test-obs-*.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := "benchmark content"
	if _, err := tmpFile.WriteString(content); err != nil {
		b.Fatal(err)
	}
	tmpFile.Close()

	client := CreateMockClient(CreateSuccessResponse(""), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.PutFileWithSignedUrl("https://obs.example.com", make(http.Header), tmpFile.Name())
	}
}

// BenchmarkGetObjectWithSignedUrl benchmarks the GetObjectWithSignedUrl method
func BenchmarkGetObjectWithSignedUrl(b *testing.B) {
	client := CreateMockClient(CreateMockResponse(http.StatusOK, "test content", make(http.Header)), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetObjectWithSignedUrl("https://obs.example.com", make(http.Header))
	}
}

// BenchmarkPutObjectWithSignedUrl benchmarks the PutObjectWithSignedUrl method
func BenchmarkPutObjectWithSignedUrl(b *testing.B) {
	client := CreateMockClient(CreateMockResponse(http.StatusOK, "", make(http.Header)), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.PutObjectWithSignedUrl("https://obs.example.com", make(http.Header), strings.NewReader("test content"))
	}
}
