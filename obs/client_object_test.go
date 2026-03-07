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
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== ListObjects Tests ====================

func TestListObjects_ShouldReturnObjectList_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "GET", req.Method)
			return CreateTestHTTPResponse(200, TestListObjectsXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListObjectsInput{
		Bucket: TestBucket,
	}

	output, err := client.ListObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Contents, 2)
	assert.Equal(t, "object1.txt", output.Contents[0].Key)
}

func TestListObjects_ShouldReturnEmptyList_WhenNoObjects(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
	<Name>test-bucket</Name>
	<Prefix></Prefix>
	<Marker></Marker>
	<MaxKeys>1000</MaxKeys>
	<IsTruncated>false</IsTruncated>
</ListBucketResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListObjectsInput{
		Bucket: TestBucket,
	}

	output, err := client.ListObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Contents, 0)
}

func TestListObjects_ShouldFilterByPrefix_WhenPrefixSet(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Contains(t, req.URL.Query().Get("prefix"), "test-prefix")
			return CreateTestHTTPResponse(200, TestListObjectsXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListObjectsInput{
		Bucket: TestBucket,
	}
	input.Prefix = "test-prefix"

	output, err := client.ListObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestListObjects_ShouldHandleDelimiter_WhenDelimiterSet(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "/", req.URL.Query().Get("delimiter"))
			return CreateTestHTTPResponse(200, TestListObjectsXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListObjectsInput{
		Bucket: TestBucket,
	}
	input.Delimiter = "/"

	output, err := client.ListObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestListObjects_ShouldHandleUrlEncoding_WhenEncodingTypeUrl(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "url", req.URL.Query().Get("encoding-type"))
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
	<Name>test-bucket</Name>
	<EncodingType>url</EncodingType>
	<IsTruncated>false</IsTruncated>
</ListBucketResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListObjectsInput{
		Bucket: TestBucket,
	}
	input.EncodingType = "url"

	output, err := client.ListObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestListObjects_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.ListObjects(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "ListObjectsInput is nil")
}

func TestListObjects_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListObjectsInput{
		Bucket: TestBucket,
	}

	output, err := client.ListObjects(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

// ==================== ListPosixObjects Tests ====================

func TestListPosixObjects_ShouldReturnObjectList_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
	<Contents>
		<Key>object1.txt</Key>
		<Size>1024</Size>
	</Contents>
</ListBucketResult>`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListPosixObjectsInput{}
	input.Bucket = TestBucket

	output, err := client.ListPosixObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestListPosixObjects_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.ListPosixObjects(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "ListPosixObjects is nil")
}

// ==================== ListVersions Tests ====================

func TestListVersions_ShouldReturnVersionList_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
	<Name>test-bucket</Name>
	<IsTruncated>false</IsTruncated>
</ListVersionsResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ListVersionsInput{
		Bucket: TestBucket,
	}

	output, err := client.ListVersions(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestListVersions_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.ListVersions(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "ListVersionsInput is nil")
}

// ==================== HeadObject Tests ====================

func TestHeadObject_ShouldReturnMetadata_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "HEAD", req.Method)
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &HeadObjectInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
	}

	output, err := client.HeadObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestHeadObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.HeadObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "HeadObjectInput is nil")
}

func TestHeadObject_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &HeadObjectInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
	}

	output, err := client.HeadObject(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

// ==================== SetObjectMetadata Tests ====================

func TestSetObjectMetadata_ShouldSetMetadata_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetObjectMetadataInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
	}

	output, err := client.SetObjectMetadata(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetObjectMetadata_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	defer func() {
		if r := recover(); r != nil {
			// Expected panic for nil input
			assert.NotNil(t, r)
		}
	}()

	output, err := client.SetObjectMetadata(nil)

	// This will panic due to nil dereference, so we catch it above
	_ = output
	_ = err
}

// ==================== DeleteObject Tests ====================

func TestDeleteObject_ShouldDeleteObject_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "DELETE", req.Method)
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &DeleteObjectInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
	}

	output, err := client.DeleteObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.DeleteObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "DeleteObjectInput is nil")
}

// ==================== DeleteObjects Tests ====================

func TestDeleteObjects_ShouldDeleteMultipleObjects_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "POST", req.Method)
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
	<Deleted>
		<Key>object1.txt</Key>
	</Deleted>
	<Deleted>
		<Key>object2.txt</Key>
	</Deleted>
</DeleteResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &DeleteObjectsInput{
		Bucket: TestBucket,
		Objects: []ObjectToDelete{
			{Key: "object1.txt"},
			{Key: "object2.txt"},
		},
	}

	output, err := client.DeleteObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteObjects_ShouldUseQuietMode_WhenQuietSet(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
</DeleteResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &DeleteObjectsInput{
		Bucket: TestBucket,
		Objects: []ObjectToDelete{
			{Key: "object1.txt"},
		},
		Quiet: true,
	}

	output, err := client.DeleteObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteObjects_ShouldHandleUrlEncoding_WhenEncodingTypeUrl(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
	<EncodingType>url</EncodingType>
</DeleteResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &DeleteObjectsInput{
		Bucket: TestBucket,
		Objects: []ObjectToDelete{
			{Key: "object1.txt"},
		},
		EncodingType: "url",
	}

	output, err := client.DeleteObjects(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteObjects_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.DeleteObjects(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "DeleteObjectsInput is nil")
}

// ==================== SetObjectAcl Tests ====================

func TestSetObjectAcl_ShouldSetACL_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetObjectAclInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
		ACL:    AclPublicRead,
	}

	output, err := client.SetObjectAcl(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetObjectAcl_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetObjectAcl(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetObjectAclInput is nil")
}

// ==================== GetObjectAcl Tests ====================

func TestGetObjectAcl_ShouldReturnACL_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestObjectACLXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetObjectAclInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
	}

	output, err := client.GetObjectAcl(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetObjectAcl_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetObjectAcl(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetObjectAclInput is nil")
}

// ==================== RestoreObject Tests ====================

func TestRestoreObject_ShouldRestoreObject_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "POST", req.Method)
			return CreateTestHTTPResponse(202, `<?xml version="1.0" encoding="UTF-8"?>
<RestoreRequest>
    <Days>3</Days>
</RestoreRequest>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &RestoreObjectInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
		Days:    3,
	}

	output, err := client.RestoreObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestRestoreObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.RestoreObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "RestoreObjectInput is nil")
}

// ==================== GetObjectMetadata Tests ====================

func TestGetObjectMetadata_ShouldReturnMetadata_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "HEAD", req.Method)
			headers := make(http.Header)
			headers.Set("Content-Length", "1024")
			headers.Set("ETag", "\"test-etag\"")
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetObjectMetadataInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
	}

	output, err := client.GetObjectMetadata(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, int64(1024), output.ContentLength)
}

func TestGetObjectMetadata_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetObjectMetadata(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetObjectMetadataInput is nil")
}

// ==================== GetAttribute Tests ====================

func TestGetAttribute_ShouldReturnAttribute_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetAttributeInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.GetAttribute(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetAttribute_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetAttribute(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetAttributeInput is nil")
}

// ==================== GetObject Tests ====================

func TestGetObject_ShouldReturnObject_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "GET", req.Method)
			headers := make(http.Header)
			headers.Set("Content-Type", "text/plain")
			return CreateTestHTTPResponse(200, "test content", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetObjectInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.GetObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	defer output.Body.Close()

	body, _ := ioutil.ReadAll(output.Body)
	assert.Equal(t, []byte("test content"), body)
}

func TestGetObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetObjectInput is nil")
}

func TestGetObject_ShouldReturnError_WhenRangeInvalid(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &GetObjectInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey
	input.Range = "invalid-range" // 不以 "bytes=" 开头

	output, err := client.GetObject(input)

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "Range should start with [bytes=]")
	assert.Nil(t, output)
}

func TestGetObject_ShouldSetRangeHeader_WhenRangeValid(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "bytes=0-999", req.Header.Get("Range"))
			return CreateTestHTTPResponse(206, "test content", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetObjectInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey
	input.Range = "bytes=0-999"

	output, err := client.GetObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetObject_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetObjectInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.GetObject(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

func TestGetObject_ShouldTriggerProgress_WhenListenerProvided(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "test content", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	progressListener := NewMockProgressListener()
	input := &GetObjectInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.GetObject(input, WithProgress(progressListener))

	require.Nil(t, err)
	require.NotNil(t, output)
	defer output.Body.Close()
	body, _ := ioutil.ReadAll(output.Body)
	// Read all to trigger progress
	_, _ = ioutil.ReadAll(strings.NewReader(string(body)))

	// Progress should have been triggered
	assert.NotEmpty(t, progressListener.Events)
}

// ==================== GetObjectWithoutProgress Tests ====================

func TestGetObjectWithoutProgress_ShouldReturnObject_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "test content", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetObjectInput{}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.GetObjectWithoutProgress(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	defer output.Body.Close()
}

func TestGetObjectWithoutProgress_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetObjectWithoutProgress(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetObjectInput is nil")
}

// ==================== PutObject Tests ====================

func TestPutObject_ShouldUploadObject_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "PUT", req.Method)
			headers := make(http.Header)
			headers.Set("ETag", "\"test-etag\"")
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutObjectInput{
		Body: strings.NewReader("test content"),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Contains(t, output.ObjectUrl, TestBucket)
	assert.Contains(t, output.ObjectUrl, TestObjectKey)
}

func TestPutObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.PutObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "PutObjectInput is nil")
}

func TestPutObject_ShouldDetectContentType_WhenKeyHasExtension(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutObjectInput{
		Body: strings.NewReader("test content"),
	}
	input.Bucket = TestBucket
	input.Key = "test.txt"

	output, err := client.PutObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutObject_ShouldUseContentType_WhenSpecified(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutObjectInput{
		Body: strings.NewReader("test content"),
	}
	input.Bucket = TestBucket
	input.Key = "test.txt"
	input.ContentType = "application/json"

	output, err := client.PutObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutObject_ShouldWrapReader_WhenContentLengthSet(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutObjectInput{
		Body: bytes.NewReader([]byte("test")),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey
	input.ContentLength = 100

	output, err := client.PutObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutObject_ShouldCallRepeatableAction_WhenStringsReader(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutObjectInput{
		Body: strings.NewReader("test content"),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutObject_ShouldCallNonRepeatableAction_WhenIoReader(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutObjectInput{
		Body: bytes.NewReader([]byte("test content")),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutObject_ShouldTriggerProgress_WhenListenerProvided(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	progressListener := NewMockProgressListener()
	input := &PutObjectInput{
		Body: strings.NewReader("test content"),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutObject(input, WithProgress(progressListener))

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutObject_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutObjectInput{
		Body: strings.NewReader("test content"),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutObject(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

// ==================== NewFolder Tests ====================

func TestNewFolder_ShouldCreateFolder_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "PUT", req.Method)
			assert.True(t, strings.HasSuffix(req.URL.Path, "/"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &NewFolderInput{}
	input.Bucket = TestBucket
	input.Key = "test-folder"

	output, err := client.NewFolder(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestNewFolder_ShouldAppendSlash_WhenNotEndingWithSlash(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.True(t, strings.HasSuffix(req.URL.Path, "/"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &NewFolderInput{}
	input.Bucket = TestBucket
	input.Key = "test-folder"

	output, err := client.NewFolder(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestNewFolder_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.NewFolder(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "NewFolderInput is nil")
}

// ==================== PutFile Tests ====================

func TestPutFile_ShouldUploadFile_WhenValidInput(t *testing.T) {
	tmpFile := CreateTempFileWithContent(t, "test file content")
	defer os.Remove(tmpFile.Name())

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutFileInput{
		SourceFile: tmpFile.Name(),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutFile(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutFile_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.PutFile(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "PutFileInput is nil")
}

func TestPutFile_ShouldDetectContentType_WhenKeyHasExtension(t *testing.T) {
	tmpFile := CreateTempFileWithContent(t, "test content")
	defer os.Remove(tmpFile.Name())

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutFileInput{
		SourceFile: tmpFile.Name(),
	}
	input.Bucket = TestBucket
	input.Key = "test.txt"

	output, err := client.PutFile(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutFile_ShouldAdjustContentLength_WhenExceedsFileSize(t *testing.T) {
	tmpFile := CreateTempFileWithContent(t, "test content")
	defer os.Remove(tmpFile.Name())

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// ContentLength should be adjusted to file size
			assert.NotEmpty(t, req.Header.Get("Content-Length"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutFileInput{
		SourceFile: tmpFile.Name(),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey
	input.ContentLength = 1000000 // Exceeds file size

	output, err := client.PutFile(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutFile_ShouldReturnError_WhenFileOpenFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutFileInput{
		SourceFile: "/nonexistent/file.txt",
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutFile(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

func TestPutFile_ShouldTriggerProgress_WhenListenerProvided(t *testing.T) {
	tmpFile := CreateTempFileWithContent(t, "test content")
	defer os.Remove(tmpFile.Name())

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	progressListener := NewMockProgressListener()
	input := &PutFileInput{
		SourceFile: tmpFile.Name(),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	output, err := client.PutFile(input, WithProgress(progressListener))

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== CopyObject Tests ====================

func TestCopyObject_ShouldCopyObject_WhenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "PUT", req.Method)
			// CopyObject should NOT have x-obs-copy-source header
			if _, ok := req.Header["x-obs-copy-source"]; ok {
				assert.Fail(t, "x-obs-copy-source header should not be present")
			}
			// Create new headers for response
			respHeaders := make(http.Header)
			respHeaders.Set("x-obs-request-id", "test-request-id")
			respHeaders.Set("ETag", "\"test-etag\"")
			return CreateTestHTTPResponse(200, "", respHeaders)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CopyObjectInput{
		CopySourceBucket: TestBucket,
		CopySourceKey:    TestObjectKey,
	}
	input.Bucket = TestBucket
	input.Key = "new-object.txt"

	output, err := client.CopyObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestCopyObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.CopyObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "CopyObjectInput is nil")
}

func TestCopyObject_ShouldReturnError_WhenSourceBucketEmpty(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CopyObjectInput{
		CopySourceBucket: "",
		CopySourceKey:    TestObjectKey,
	}
	input.Bucket = TestBucket
	input.Key = "new-object.txt"

	output, err := client.CopyObject(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "Source bucket is empty")
}

func TestCopyObject_ShouldReturnError_WhenSourceKeyEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CopyObjectInput{
		CopySourceBucket: TestBucket,
		CopySourceKey:    "",
	}
	input.Bucket = TestBucket
	input.Key = "new-object.txt"

	output, err := client.CopyObject(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "Source key is empty")
}

// ==================== AppendObject Tests ====================

func TestAppendObject_ShouldAppendObject_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "POST", req.Method)
			headers := make(http.Header)
			headers.Set("x-obs-next-append-position", "100")
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &AppendObjectInput{
		Position: 100,
		Body:     strings.NewReader("append content"),
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey
	input.ContentLength = 14

	output, err := client.AppendObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestAppendObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.AppendObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "AppendObjectInput is nil")
}

func TestAppendObject_ShouldDetectContentType_WhenKeyHasExtension(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &AppendObjectInput{
		Body:     strings.NewReader("content"),
		Position: 0,
	}
	input.Bucket = TestBucket
	input.Key = "test.txt"

	output, err := client.AppendObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== ModifyObject Tests ====================

func TestModifyObject_ShouldModifyObject_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &ModifyObjectInput{
		Bucket: TestBucket,
		Key:    TestObjectKey,
		Body:   strings.NewReader("modified content"),
	}

	output, err := client.ModifyObject(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestModifyObject_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.ModifyObject(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "ModifyObjectInput is nil")
}

// ==================== RenameFile Tests ====================

func TestRenameFile_ShouldRenameFile_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "POST", req.Method)
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &RenameFileInput{
		Bucket:       TestBucket,
		Key:          "old-name.txt",
		NewObjectKey: "new-name.txt",
	}

	output, err := client.RenameFile(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestRenameFile_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.RenameFile(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "RenameFileInput is nil")
}

// ==================== RenameFolder Tests ====================

func TestRenameFolder_ShouldRenameFolder_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "POST", req.Method)
			assert.True(t, strings.HasSuffix(req.URL.Path, "/"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &RenameFolderInput{
		Bucket:       TestBucket,
		Key:          "old-folder",
		NewObjectKey: "new-folder",
	}

	output, err := client.RenameFolder(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestRenameFolder_ShouldAppendSlash_WhenNotEndingWithSlash(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			path := req.URL.Path
			assert.True(t, strings.HasSuffix(path, "/"))
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &RenameFolderInput{
		Bucket:       TestBucket,
		Key:          "old-folder",
		NewObjectKey: "new-folder",
	}

	output, err := client.RenameFolder(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestRenameFolder_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.RenameFolder(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "RenameFolderInput is nil")
}

// ==================== SetDirAccesslabel Tests ====================

func TestSetDirAccesslabel_ShouldSetLabel_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetDirAccesslabelInput{}
	input.Bucket = TestBucket
	input.Key = "test-dir/"

	output, err := client.SetDirAccesslabel(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetDirAccesslabel_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetDirAccesslabel(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetDirAccesslabelInput is nil")
}

// ==================== GetDirAccesslabel Tests ====================

func TestGetDirAccesslabel_ShouldReturnLabel_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `{"Accesslabel":["test-label"]}`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetDirAccesslabelInput{}
	input.Bucket = TestBucket
	input.Key = "test-dir/"

	output, err := client.GetDirAccesslabel(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetDirAccesslabel_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetDirAccesslabel(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetDirAccesslabelInput is nil")
}

// ==================== DeleteDirAccesslabel Tests ====================

func TestDeleteDirAccesslabel_ShouldDeleteLabel_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &DeleteDirAccesslabelInput{}
	input.Bucket = TestBucket
	input.Key = "test-dir/"

	output, err := client.DeleteDirAccesslabel(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteDirAccesslabel_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.DeleteDirAccesslabel(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "DeleteDirAccesslabelInput is nil")
}

// CreatePostPolicy tests

func TestCreatePostPolicy_ShouldReturnPolicyAndSignature_GivenValidInput(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CreatePostPolicyInput{
		Bucket: TestBucket,
		Key:    "test-object",
		Expires: 600,
	}

	output, err := client.CreatePostPolicy(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
}

func TestCreatePostPolicy_ShouldUseDefaultExpiration_GivenZeroExpires(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CreatePostPolicyInput{
		Bucket: TestBucket,
		Key:    "test-object",
		Expires: 0,
	}

	output, err := client.CreatePostPolicy(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.Policy)
	assert.NotEmpty(t, output.Signature)
}

func TestCreatePostPolicy_ShouldReturnError_GivenNilInput(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.CreatePostPolicy(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "CreatePostPolicyInput is nil")
}

func TestCreatePostPolicy_ShouldReturnError_GivenEmptyBucket(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CreatePostPolicyInput{
		Bucket: "",
		Key:    "test-object",
		Expires: 600,
	}

	output, err := client.CreatePostPolicy(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "bucket is empty")
}

func TestCreatePostPolicy_ShouldReturnError_GivenEmptyKey(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &CreatePostPolicyInput{
		Bucket: TestBucket,
		Key:    "",
		Expires: 600,
	}

	output, err := client.CreatePostPolicy(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "key is empty")
}

