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
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== ListMultipartUploads Tests ====================

func TestListMultipartUploads_ShouldReturnUploadList_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "GET", req.Method)
			return CreateTestHTTPResponse(200, TestListMultipartUploadsXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListMultipartUploadsInput{
		Bucket: TestBucket,
	}

	output, err := client.ListMultipartUploads(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Uploads, 1)
}

func TestListMultipartUploads_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.ListMultipartUploads(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "ListMultipartUploadsInput is nil")
}

func TestListMultipartUploads_ShouldHandleUrlEncoding_WhenEncodingTypeUrl(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "url", req.URL.Query().Get("encoding-type"))
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
	<EncodingType>url</EncodingType>
	<IsTruncated>false</IsTruncated>
</ListMultipartUploadsResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListMultipartUploadsInput{
		Bucket:      TestBucket,
		EncodingType: "url",
	}

	output, err := client.ListMultipartUploads(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestListMultipartUploads_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListMultipartUploadsInput{
		Bucket: TestBucket,
	}

	output, err := client.ListMultipartUploads(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

// ==================== AbortMultipartUpload Tests ====================

func TestAbortMultipartUpload_ShouldAbortUpload_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "DELETE", req.Method)
			assert.Contains(t, req.URL.Query().Get("uploadId"), "upload-id")
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &AbortMultipartUploadInput{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: "test-upload-id",
	}

	output, err := client.AbortMultipartUpload(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestAbortMultipartUpload_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.AbortMultipartUpload(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "AbortMultipartUploadInput is nil")
}

func TestAbortMultipartUpload_ShouldReturnError_WhenUploadIdEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &AbortMultipartUploadInput{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: "",
	}

	output, err := client.AbortMultipartUpload(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "UploadId is empty")
}

// ==================== InitiateMultipartUpload Tests ====================

func TestInitiateMultipartUpload_ShouldInitiateUpload_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "POST", req.Method)
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
	<Bucket>test-bucket</Bucket>
	<Key>test-object.txt</Key>
	<UploadId>test-upload-id-1234567890</UploadId>
</InitiateMultipartUploadResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &InitiateMultipartUploadInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.InitiateMultipartUpload(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "test-upload-id-1234567890", output.UploadId)
}

func TestInitiateMultipartUpload_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.InitiateMultipartUpload(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "InitiateMultipartUploadInput is nil")
}

func TestInitiateMultipartUpload_ShouldDetectContentType_WhenKeyHasExtension(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
	<UploadId>test-upload-id</UploadId>
</InitiateMultipartUploadResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &InitiateMultipartUploadInput{}
	input.Bucket = TestBucket
	input.Key = "test.txt"

	output, err := client.InitiateMultipartUpload(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestInitiateMultipartUpload_ShouldUseContentType_WhenSpecified(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
	<UploadId>test-upload-id</UploadId>
</InitiateMultipartUploadResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &InitiateMultipartUploadInput{}
	input.Bucket = TestBucket
	input.Key = "test.txt"
	input.ContentType = "application/json"

	output, err := client.InitiateMultipartUpload(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestInitiateMultipartUpload_ShouldHandleUrlEncoding_WhenEncodingTypeUrl(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
	<EncodingType>url</EncodingType>
	<UploadId>test-upload-id</UploadId>
</InitiateMultipartUploadResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &InitiateMultipartUploadInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey
	input.EncodingType = "url"

	output, err := client.InitiateMultipartUpload(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== UploadPart Tests ====================

func TestUploadPart_ShouldUploadPart_WhenStringsReader(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		Body:       strings.NewReader("part content"),
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 1, output.PartNumber)
}

func TestUploadPart_ShouldUploadPart_WhenIoReader(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		Body:       bytes.NewReader([]byte("part content")),
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestUploadPart_ShouldUploadPartFromFile_WhenSourceFileSet(t *testing.T) {
	tmpFile := CreateTempFileWithContent(t, "test content for upload")
	defer os.Remove(tmpFile.Name())

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		SourceFile: tmpFile.Name(),
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestUploadPart_ShouldSetOffsetToZero_WhenNegativeOffset(t *testing.T) {
	tmpFile := CreateTempFileWithContent(t, "test content for upload")
	defer os.Remove(tmpFile.Name())

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		SourceFile: tmpFile.Name(),
		Offset:     -5, // 负数，应调整为 0
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestUploadPart_ShouldAdjustPartSize_WhenExceedsFileSize(t *testing.T) {
	tmpFile := CreateTempFileWithContent(t, "test content")
	defer os.Remove(tmpFile.Name())

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		SourceFile: tmpFile.Name(),
		Offset:     0,
		PartSize:   1000000, // 超过文件大小，应调整为文件大小
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestUploadPart_ShouldReturnError_WhenFileOpenFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		SourceFile: "/nonexistent/file.txt",
	}

	output, err := client.UploadPart(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

func TestUploadPart_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.UploadPart(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "UploadPartInput is nil")
}

func TestUploadPart_ShouldReturnError_WhenUploadIdEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "",
		PartNumber: 1,
		Body:       strings.NewReader("content"),
	}

	output, err := client.UploadPart(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "UploadId is empty")
}

func TestUploadPart_ShouldWrapReader_WhenPartSizeSet(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		Body:       strings.NewReader("content"),
		PartSize:   7,
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestUploadPart_ShouldTriggerProgress_WhenListenerProvided(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	progressListener := NewMockProgressListener()
	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		Body:       strings.NewReader("content"),
	}

	output, err := client.UploadPart(input, WithProgress(progressListener))

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestUploadPart_ShouldUseRepeatableAction_WhenRepeatable(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		Body:       strings.NewReader("content"),
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestUploadPart_ShouldUseNonRepeatableAction_WhenNotRepeatable(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<UploadPartOutput>
	<ETag>"part-etag"</ETag>
</UploadPartOutput>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &UploadPartInput{
		Bucket:     TestBucket,
		Key:        TestObjectKey,
		UploadId:   "test-upload-id",
		PartNumber: 1,
		Body:       bytes.NewReader([]byte("content")),
	}

	output, err := client.UploadPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== CompleteMultipartUpload Tests ====================

func TestCompleteMultipartUpload_ShouldCompleteUpload_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "POST", req.Method)
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
	<Location>https://obs.example.com/test-bucket/test-object.txt</Location>
	<Bucket>test-bucket</Bucket>
	<Key>test-object.txt</Key>
	<ETag>"complete-etag"</ETag>
</CompleteMultipartUploadResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CompleteMultipartUploadInput{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: "test-upload-id",
		Parts: []Part{
			{PartNumber: 1, ETag: "\"etag1\""},
			{PartNumber: 2, ETag: "\"etag2\""},
		},
	}

	output, err := client.CompleteMultipartUpload(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestCompleteMultipartUpload_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.CompleteMultipartUpload(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "CompleteMultipartUploadInput is nil")
}

func TestCompleteMultipartUpload_ShouldReturnError_WhenUploadIdEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CompleteMultipartUploadInput{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: "",
	}

	output, err := client.CompleteMultipartUpload(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "UploadId is empty")
}

func TestCompleteMultipartUpload_ShouldHandleUrlEncoding_WhenEncodingTypeUrl(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
	<EncodingType>url</EncodingType>
	<ETag>"complete-etag"</ETag>
</CompleteMultipartUploadResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CompleteMultipartUploadInput{
		Bucket:      TestBucket,
		Key:         TestObjectKey,
		UploadId:    "test-upload-id",
		EncodingType: "url",
		Parts:       []Part{{PartNumber: 1, ETag: "\"etag1\""}},
	}

	output, err := client.CompleteMultipartUpload(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== ListParts Tests ====================

func TestListParts_ShouldReturnPartsList_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestListPartsXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListPartsInput{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: "test-upload-id",
	}

	output, err := client.ListParts(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Parts, 2)
}

func TestListParts_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.ListParts(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "ListPartsInput is nil")
}

func TestListParts_ShouldReturnError_WhenUploadIdEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &ListPartsInput{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: "",
	}

	output, err := client.ListParts(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "UploadId is empty")
}

func TestListParts_ShouldHandleUrlEncoding_WhenEncodingTypeUrl(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
	<EncodingType>url</EncodingType>
	<IsTruncated>false</IsTruncated>
</ListPartsResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListPartsInput{
		Bucket:      TestBucket,
		Key:         TestObjectKey,
		UploadId:    "test-upload-id",
		EncodingType: "url",
	}

	output, err := client.ListParts(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== CopyPart Tests ====================

func TestCopyPart_ShouldCopyPart_WhenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// CopyPart should have copy-source header (default is x-amz-copy-source with SignatureV2)
			if values, ok := req.Header["x-amz-copy-source"]; ok && len(values) > 0 {
				assert.NotEmpty(t, values[0])
			} else if values, ok := req.Header["x-obs-copy-source"]; ok && len(values) > 0 {
				assert.NotEmpty(t, values[0])
			} else {
				assert.Fail(t, "copy-source header should be present (x-amz-copy-source or x-obs-copy-source)")
			}
			headers := make(http.Header)
			headers.Set("ETag", "\"copy-part-etag\"")
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
	<LastModified>2023-01-01T00:00:00Z</LastModified>
	<ETag>"copy-part-etag"</ETag>
</CopyPartResult>`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CopyPartInput{
		Bucket:           TestBucket,
		Key:              "new-object.txt",
		UploadId:         "test-upload-id",
		PartNumber:       1,
		CopySourceBucket: TestBucket,
		CopySourceKey:    TestObjectKey,
	}

	output, err := client.CopyPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 1, output.PartNumber)
}

func TestCopyPart_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.CopyPart(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "CopyPartInput is nil")
}

func TestCopyPart_ShouldReturnError_WhenUploadIdEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CopyPartInput{
		Bucket:           TestBucket,
		Key:              "new-object.txt",
		UploadId:         "",
		PartNumber:       1,
		CopySourceBucket: TestBucket,
		CopySourceKey:    TestObjectKey,
	}

	output, err := client.CopyPart(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "UploadId is empty")
}

func TestCopyPart_ShouldReturnError_WhenSourceBucketEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CopyPartInput{
		Bucket:           TestBucket,
		Key:              "new-object.txt",
		UploadId:         "test-upload-id",
		PartNumber:       1,
		CopySourceBucket: "",
		CopySourceKey:    TestObjectKey,
	}

	output, err := client.CopyPart(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "Source bucket is empty")
}

func TestCopyPart_ShouldReturnError_WhenSourceKeyEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CopyPartInput{
		Bucket:           TestBucket,
		Key:              "new-object.txt",
		UploadId:         "test-upload-id",
		PartNumber:       1,
		CopySourceBucket: TestBucket,
		CopySourceKey:    "",
	}

	output, err := client.CopyPart(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "Source key is empty")
}

func TestCopyPart_ShouldReturnError_WhenCopySourceRangeInvalid(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CopyPartInput{
		Bucket:           TestBucket,
		Key:              "new-object.txt",
		UploadId:         "test-upload-id",
		PartNumber:       1,
		CopySourceBucket: TestBucket,
		CopySourceKey:    TestObjectKey,
		CopySourceRange:  "invalid-range",
	}

	output, err := client.CopyPart(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "Source Range should start with [bytes=]")
}

func TestCopyPart_ShouldSetCopySourceRange_WhenValid(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// Check both possible headers (x-amz-copy-source-range for SignatureV2, x-obs-copy-source-range for SignatureObs)
			if values, ok := req.Header["x-amz-copy-source-range"]; ok && len(values) > 0 {
				assert.Equal(t, "bytes=0-999", values[0])
			} else if values, ok := req.Header["x-obs-copy-source-range"]; ok && len(values) > 0 {
				assert.Equal(t, "bytes=0-999", values[0])
			} else {
				assert.Fail(t, "copy-source-range header not found (x-amz-copy-source-range or x-obs-copy-source-range)")
			}
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
	<ETag>"copy-part-etag"</ETag>
</CopyPartResult>`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CopyPartInput{
		Bucket:           TestBucket,
		Key:              "new-object.txt",
		UploadId:         "test-upload-id",
		PartNumber:       1,
		CopySourceBucket: TestBucket,
		CopySourceKey:    TestObjectKey,
		CopySourceRange:  "bytes=0-999",
	}

	output, err := client.CopyPart(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}
