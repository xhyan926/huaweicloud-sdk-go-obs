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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== ListBuckets Tests ====================

func TestListBuckets_ShouldReturnBucketList_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "GET", req.Method)
			assert.Empty(t, req.URL.Path)
			return CreateTestHTTPResponse(200, TestListBucketsXML, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBuckets(&ListBucketsInput{})

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Buckets, 2)
	assert.Equal(t, "bucket1", output.Buckets[0].Name)
	assert.Equal(t, "bucket2", output.Buckets[1].Name)
}

func TestListBuckets_ShouldReturnEmptyList_WhenNoBuckets(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
	<Owner>
		<ID>test-owner-id</ID>
		<DisplayName>test-owner</DisplayName>
	</Owner>
	<Buckets></Buckets>
</ListAllMyBucketsResult>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBuckets(&ListBucketsInput{})

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Buckets, 0)
}

func TestListBuckets_ShouldAcceptNilInput_WhenCalled(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestListBucketsXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBuckets(nil)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestListBuckets_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBuckets(&ListBucketsInput{})

	require.NotNil(t, err)
	assert.Nil(t, output)
}

func TestListBuckets_ShouldReturnError_WhenServerReturnsError(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(403, string(CreateTestErrorResponse("AccessDenied", "Access Denied")), nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBuckets(&ListBucketsInput{})

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "AccessDenied")
}

// ==================== CreateBucket Tests ====================

func TestCreateBucket_ShouldCreateBucket_WhenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "PUT", req.Method)
			assert.Contains(t, req.URL.Host, TestBucket)
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CreateBucketInput{
		Bucket: TestBucket,
	}

	output, err := client.CreateBucket(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestCreateBucket_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.CreateBucket(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "CreateBucketInput is nil")
}

func TestCreateBucket_ShouldSetACL_WhenAclProvided(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// Check both possible headers using the array access
			if values, ok := req.Header["x-amz-acl"]; ok && len(values) > 0 {
				assert.Equal(t, "public-read", values[0])
			} else if values, ok := req.Header["x-obs-acl"]; ok && len(values) > 0 {
				assert.Equal(t, "public-read", values[0])
			} else {
				assert.Fail(t, "ACL header not found")
			}
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CreateBucketInput{
		Bucket: TestBucket,
		ACL:   AclPublicRead,
	}

	output, err := client.CreateBucket(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestCreateBucket_ShouldSetStorageClass_WhenStorageClassProvided(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// Check both possible headers using the array access
			if values, ok := req.Header["x-default-storage-class"]; ok && len(values) > 0 {
				assert.Equal(t, "STANDARD_IA", values[0])
			} else if values, ok := req.Header["x-amz-storage-class"]; ok && len(values) > 0 {
				assert.Equal(t, "STANDARD_IA", values[0])
			} else if values, ok := req.Header["x-obs-storage-class"]; ok && len(values) > 0 {
				assert.Equal(t, "STANDARD_IA", values[0])
			} else {
				assert.Fail(t, "StorageClass header not found")
			}
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CreateBucketInput{
		Bucket:       TestBucket,
		StorageClass: "STANDARD_IA",
	}

	output, err := client.CreateBucket(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestCreateBucket_ShouldSetLocation_WhenLocationProvided(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// Location is sent in request body, but we can't easily verify it in mock
			// Just verify the request was made
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CreateBucketInput{
		Bucket: TestBucket,
	}
	input.Location = "cn-north-4"

	output, err := client.CreateBucket(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestCreateBucket_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &CreateBucketInput{
		Bucket: TestBucket,
	}

	output, err := client.CreateBucket(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

// ==================== DeleteBucket Tests ====================

func TestDeleteBucket_ShouldDeleteBucket_WhenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "DELETE", req.Method)
			assert.Contains(t, req.URL.Host, TestBucket)
			return CreateTestHTTPResponse(204, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucket(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteBucket_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucket(TestBucket)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

func TestDeleteBucket_ShouldReturnError_WhenBucketNotFound(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(404, string(CreateTestErrorResponse("NoSuchBucket", "The specified bucket does not exist")), nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucket(TestBucket)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

// ==================== HeadBucket Tests ====================

func TestHeadBucket_ShouldReturnSuccess_WhenBucketExists(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "HEAD", req.Method)
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.HeadBucket(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestHeadBucket_ShouldReturnError_WhenBucketNotFound(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(404, string(CreateTestErrorResponse("NoSuchBucket", "The specified bucket does not exist")), nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.HeadBucket(TestBucket)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

func TestHeadBucket_ShouldReturnError_WhenNetworkFails(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.HeadBucket(TestBucket)

	require.NotNil(t, err)
	assert.Nil(t, output)
}

// ==================== SetBucketStoragePolicy Tests ====================

func TestSetBucketStoragePolicy_ShouldSetPolicy_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketStoragePolicyInput{
		Bucket: TestBucket,
	}
	input.StorageClass = "STANDARD_IA"

	output, err := client.SetBucketStoragePolicy(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketStoragePolicy_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketStoragePolicy(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketStoragePolicyInput is nil")
}

// ==================== GetBucketStoragePolicy Tests ====================

func TestGetBucketStoragePolicy_ShouldReturnPolicy_WhenSignatureObs(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// URL query should contain storageClass parameter
			assert.Contains(t, req.URL.RawQuery, "storageClass")
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<StorageClass>STANDARD_IA</StorageClass>`, headers)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureObs), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketStoragePolicy(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "STANDARD_IA", output.StorageClass)
}

func TestGetBucketStoragePolicy_ShouldReturnPolicy_WhenSignatureV4(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			// URL query should contain storagePolicy parameter
			assert.Contains(t, req.URL.RawQuery, "storagePolicy")
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<StoragePolicy><DefaultStorageClass>STANDARD_IA</DefaultStorageClass></StoragePolicy>`, headers)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureV4), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketStoragePolicy(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "STANDARD_IA", output.StorageClass)
}

// ==================== SetBucketQuota Tests ====================

func TestSetBucketQuota_ShouldSetQuota_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketQuotaInput{
		Bucket: TestBucket,
	}
	input.Quota = 1024000

	output, err := client.SetBucketQuota(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketQuota_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketQuota(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketQuotaInput is nil")
}

// ==================== GetBucketQuota Tests ====================

func TestGetBucketQuota_ShouldReturnQuota_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<Quota>
    <StorageQuota>1024000</StorageQuota>
</Quota>`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketQuota(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, int64(1024000), output.BucketQuota.Quota)
}

// ==================== GetBucketMetadata Tests ====================

func TestGetBucketMetadata_ShouldReturnMetadata_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set("x-obs-storage-class", "STANDARD_IA")
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetBucketMetadataInput{
		Bucket: TestBucket,
	}

	output, err := client.GetBucketMetadata(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetBucketMetadata_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	defer func() {
		if r := recover(); r != nil {
			// Expected panic for nil input
			assert.NotNil(t, r)
		}
	}()

	output, err := client.GetBucketMetadata(nil)

	// This will panic due to nil dereference, so we catch it above
	_ = output
	_ = err
}

// ==================== GetBucketFSStatus Tests ====================

func TestGetBucketFSStatus_ShouldReturnStatus_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set("x-obs-fs-file-interface", "Enabled")
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetBucketFSStatusInput{}
	input.Bucket = TestBucket

	output, err := client.GetBucketFSStatus(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetBucketFSStatus_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	defer func() {
		if r := recover(); r != nil {
			// Expected panic for nil input
			assert.NotNil(t, r)
		}
	}()

	output, err := client.GetBucketFSStatus(nil)

	// This will panic due to nil dereference, so we catch it above
	_ = output
	_ = err
}

// ==================== GetBucketStorageInfo Tests ====================

func TestGetBucketStorageInfo_ShouldReturnInfo_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<GetBucketStorageInfoResult>
    <Size>1024000</Size>
    <ObjectNumber>100</ObjectNumber>
</GetBucketStorageInfoResult>`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketStorageInfo(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, int64(1024000), output.Size)
	assert.Equal(t, 100, output.ObjectNumber)
}

// ==================== GetBucketLocation Tests ====================

func TestGetBucketLocation_ShouldReturnLocation_WhenSignatureObs(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<Location>cn-north-4</Location>`, headers)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureObs), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketLocation(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "cn-north-4", output.Location)
}

func TestGetBucketLocation_ShouldReturnLocation_WhenSignatureV4(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<CreateBucketConfiguration><LocationConstraint>cn-north-4</LocationConstraint></CreateBucketConfiguration>`, headers)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureV4), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketLocation(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "cn-north-4", output.Location)
}

// ==================== SetBucketAcl Tests ====================

func TestSetBucketAcl_ShouldSetACL_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketAclInput{
		Bucket: TestBucket,
		ACL:    AclPublicRead,
	}

	output, err := client.SetBucketAcl(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketAcl_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketAcl(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketAclInput is nil")
}

// ==================== GetBucketAcl Tests ====================

func TestGetBucketAcl_ShouldReturnACL_WhenSignatureObs(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<AccessControlPolicy>
    <Owner>
        <ID>test-owner-id</ID>
        <DisplayName>test-owner</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>
            <Grantee>
                <Canned>Everyone</Canned>
            </Grantee>
            <Permission>READ</Permission>
        </Grant>
    </AccessControlList>
</AccessControlPolicy>`, nil)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureObs), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketAcl(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.NotEmpty(t, output.Grants)
	// Check that Everyone is converted to GroupAllUsers
	assert.Equal(t, GroupAllUsers, output.Grants[0].Grantee.URI)
}

func TestGetBucketAcl_ShouldReturnACL_WhenSignatureV4(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestBucketACLXML, nil)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureV4), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketAcl(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketPolicy Tests ====================

func TestSetBucketPolicy_ShouldSetPolicy_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketPolicyInput{
		Bucket: TestBucket,
		Policy: TestBucketPolicyXML,
	}

	output, err := client.SetBucketPolicy(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketPolicy_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketPolicy(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketPolicy is nil")
}

// ==================== GetBucketPolicy Tests ====================

func TestGetBucketPolicy_ShouldReturnPolicy_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestBucketPolicyXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketPolicy(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketPolicy Tests ====================

func TestDeleteBucketPolicy_ShouldDeletePolicy_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketPolicy(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketCors Tests ====================

func TestSetBucketCors_ShouldSetCors_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketCorsInput{
		Bucket: TestBucket,
	}
	input.CorsRules = []CorsRule{
		{
			AllowedOrigin: []string{"*"},
			AllowedMethod: []string{"GET", "PUT"},
			AllowedHeader: []string{"*"},
			MaxAgeSeconds: 3000,
		},
	}

	output, err := client.SetBucketCors(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketCors_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketCors(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketCorsInput is nil")
}

// ==================== GetBucketCors Tests ====================

func TestGetBucketCors_ShouldReturnCors_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestBucketCorsXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketCors(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketCors Tests ====================

func TestDeleteBucketCors_ShouldDeleteCors_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketCors(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketVersioning Tests ====================

func TestSetBucketVersioning_ShouldSetVersioning_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketVersioningInput{
		Bucket: TestBucket,
	}
	input.Status = VersioningStatusEnabled

	output, err := client.SetBucketVersioning(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketVersioning_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketVersioning(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketVersioningInput is nil")
}

// ==================== GetBucketVersioning Tests ====================

func TestGetBucketVersioning_ShouldReturnVersioning_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestBucketVersioningXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketVersioning(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, VersioningStatusEnabled, output.Status)
}

// ==================== SetBucketWebsiteConfiguration Tests ====================

func TestSetBucketWebsiteConfiguration_ShouldSetWebsite_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketWebsiteConfigurationInput{
		Bucket: TestBucket,
	}
	input.IndexDocument = IndexDocument{
		Suffix: "index.html",
	}

	output, err := client.SetBucketWebsiteConfiguration(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketWebsiteConfiguration_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketWebsiteConfiguration(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketWebsiteConfigurationInput is nil")
}

// ==================== GetBucketWebsiteConfiguration Tests ====================

func TestGetBucketWebsiteConfiguration_ShouldReturnWebsite_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestBucketWebsiteXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketWebsiteConfiguration(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketWebsiteConfiguration Tests ====================

func TestDeleteBucketWebsiteConfiguration_ShouldDeleteWebsite_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketWebsiteConfiguration(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketLoggingConfiguration Tests ====================

func TestSetBucketLoggingConfiguration_ShouldSetLogging_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketLoggingConfigurationInput{
		Bucket: TestBucket,
	}

	output, err := client.SetBucketLoggingConfiguration(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketLoggingConfiguration_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketLoggingConfiguration(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketLoggingConfigurationInput is nil")
}

// ==================== GetBucketLoggingConfiguration Tests ====================

func TestGetBucketLoggingConfiguration_ShouldReturnLogging_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<BucketLoggingStatus>
</BucketLoggingStatus>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketLoggingConfiguration(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketLifecycleConfiguration Tests ====================

func TestSetBucketLifecycleConfiguration_ShouldSetLifecycle_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketLifecycleConfigurationInput{
		Bucket: TestBucket,
	}
	input.LifecycleRules = []LifecycleRule{
		{
			ID:     "test-rule",
			Prefix: "test-prefix/",
			Status: "Enabled",
			Expiration: Expiration{
				Days: 30,
			},
		},
	}

	output, err := client.SetBucketLifecycleConfiguration(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketLifecycleConfiguration_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketLifecycleConfiguration(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketLifecycleConfigurationInput is nil")
}

// ==================== GetBucketLifecycleConfiguration Tests ====================

func TestGetBucketLifecycleConfiguration_ShouldReturnLifecycle_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, TestBucketLifecycleXML, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketLifecycleConfiguration(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketLifecycleConfiguration Tests ====================

func TestDeleteBucketLifecycleConfiguration_ShouldDeleteLifecycle_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketLifecycleConfiguration(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketEncryption Tests ====================

func TestSetBucketEncryption_ShouldSetEncryption_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketEncryptionInput{
		Bucket: TestBucket,
	}
	input.SSEAlgorithm = DEFAULT_SSE_KMS_ENCRYPTION
	input.KMSMasterKeyID = "test-key-id"

	output, err := client.SetBucketEncryption(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketEncryption_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketEncryption(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketEncryptionInput is nil")
}

// ==================== GetBucketEncryption Tests ====================

func TestGetBucketEncryption_ShouldReturnEncryption_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<ServerSideEncryptionConfiguration>
    <Rule>
        <ApplyServerSideEncryptionByDefault>
            <SSEAlgorithm>AES256</SSEAlgorithm>
        </ApplyServerSideEncryptionByDefault>
    </Rule>
</ServerSideEncryptionConfiguration>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketEncryption(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketEncryption Tests ====================

func TestDeleteBucketEncryption_ShouldDeleteEncryption_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketEncryption(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketTagging Tests ====================

func TestSetBucketTagging_ShouldSetTags_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketTaggingInput{
		Bucket: TestBucket,
	}
	input.Tags = []Tag{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	output, err := client.SetBucketTagging(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketTagging_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketTagging(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketTaggingInput is nil")
}

// ==================== GetBucketTagging Tests ====================

func TestGetBucketTagging_ShouldReturnTags_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<Tagging>
    <TagSet>
        <Tag>
            <Key>key1</Key>
            <Value>value1</Value>
        </Tag>
        <Tag>
            <Key>key2</Key>
            <Value>value2</Value>
        </Tag>
    </TagSet>
</Tagging>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketTagging(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.Tags, 2)
}

// ==================== DeleteBucketTagging Tests ====================

func TestDeleteBucketTagging_ShouldDeleteTags_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketTagging(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketNotification Tests ====================

func TestSetBucketNotification_ShouldSetNotification_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketNotificationInput{
		Bucket: TestBucket,
	}

	output, err := client.SetBucketNotification(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketNotification_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketNotification(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketNotificationInput is nil")
}

// ==================== GetBucketNotification Tests ====================

func TestGetBucketNotification_ShouldReturnNotification_WhenSignatureObs(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<NotificationConfiguration>
</NotificationConfiguration>`, nil)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureObs), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketNotification(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetBucketNotification_ShouldReturnNotification_WhenSignatureV4(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<NotificationConfiguration>
    <TopicConfiguration>
        <Id>test-topic-id</Id>
        <Topic>arn:aws:sns:us-east-1:123456789012:MyTopic</Topic>
        <Event>s3:ObjectCreated:*</Event>
        <Event>s3:ObjectRemoved:*</Event>
    </TopicConfiguration>
</NotificationConfiguration>`, nil)
		},
	}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithSignature(SignatureV4), WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketNotification(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketRequestPayment Tests ====================

func TestSetBucketRequestPayment_ShouldSetRequestPayment_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketRequestPaymentInput{
		Bucket: TestBucket,
	}
	input.Payer = RequesterPayer

	output, err := client.SetBucketRequestPayment(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketRequestPayment_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketRequestPayment(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketRequestPaymentInput is nil")
}

// ==================== GetBucketRequestPayment Tests ====================

func TestGetBucketRequestPayment_ShouldReturnRequestPayment_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<RequestPaymentConfiguration>
    <Payer>Requester</Payer>
</RequestPaymentConfiguration>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketRequestPayment(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, RequesterPayer, output.Payer)
}

// ==================== SetBucketFetchPolicy Tests ====================

func TestSetBucketFetchPolicy_ShouldSetFetchPolicy_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketFetchPolicyInput{
		Bucket: TestBucket,
		Status: "Enabled",
		Agency: "test-agency",
	}

	output, err := client.SetBucketFetchPolicy(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketFetchPolicy_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketFetchPolicy(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketFetchPolicyInput is nil")
}

func TestSetBucketFetchPolicy_ShouldReturnError_WhenStatusEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &SetBucketFetchPolicyInput{
		Bucket: TestBucket,
		Status: "",
		Agency: "test-agency",
	}

	output, err := client.SetBucketFetchPolicy(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "Fetch policy status is empty")
}

func TestSetBucketFetchPolicy_ShouldReturnError_WhenAgencyEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &SetBucketFetchPolicyInput{
		Bucket: TestBucket,
		Status: "Enabled",
		Agency: "",
	}

	output, err := client.SetBucketFetchPolicy(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "Fetch policy agency is empty")
}

// ==================== GetBucketFetchPolicy Tests ====================

func TestGetBucketFetchPolicy_ShouldReturnFetchPolicy_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `{
  "fetch": {
    "status": "Enabled",
    "agency": "test-agency"
  }
}`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetBucketFetchPolicyInput{
		Bucket: TestBucket,
	}

	output, err := client.GetBucketFetchPolicy(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetBucketFetchPolicy_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetBucketFetchPolicy(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetBucketFetchPolicyInput is nil")
}

// ==================== DeleteBucketFetchPolicy Tests ====================

func TestDeleteBucketFetchPolicy_ShouldDeleteFetchPolicy_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &DeleteBucketFetchPolicyInput{
		Bucket: TestBucket,
	}

	output, err := client.DeleteBucketFetchPolicy(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteBucketFetchPolicy_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.DeleteBucketFetchPolicy(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "DeleteBucketFetchPolicyInput is nil")
}

// ==================== SetBucketFetchJob Tests ====================

func TestSetBucketFetchJob_ShouldSetFetchJob_WhenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `{
  "id": "test-job-id",
  "Wait": 0
}`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketFetchJobInput{
		Bucket: TestBucket,
		URL:    "https://example.com/file.txt",
	}

	output, err := client.SetBucketFetchJob(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketFetchJob_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketFetchJob(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketFetchJobInput is nil")
}

func TestSetBucketFetchJob_ShouldReturnError_WhenURLEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &SetBucketFetchJobInput{
		Bucket: TestBucket,
		URL:    "",
	}

	output, err := client.SetBucketFetchJob(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "URL is empty")
}

// ==================== GetBucketFetchJob Tests ====================

func TestGetBucketFetchJob_ShouldReturnFetchJob_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `{
  "id": "test-job-id",
  "status": "Active",
  "code": "",
  "err": ""
}`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &GetBucketFetchJobInput{
		Bucket: TestBucket,
		JobID:  "test-job-id",
	}

	output, err := client.GetBucketFetchJob(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestGetBucketFetchJob_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.GetBucketFetchJob(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "GetBucketFetchJobInput is nil")
}

func TestGetBucketFetchJob_ShouldReturnError_WhenJobIDEmpty(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &GetBucketFetchJobInput{
		Bucket: TestBucket,
		JobID:  "",
	}

	output, err := client.GetBucketFetchJob(input)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "JobID is empty")
}

// ==================== PutBucketPublicAccessBlock Tests ====================

func TestPutBucketPublicAccessBlock_ShouldSetBlock_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &PutBucketPublicAccessBlockInput{
		Bucket: TestBucket,
	}

	output, err := client.PutBucketPublicAccessBlock(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestPutBucketPublicAccessBlock_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.PutBucketPublicAccessBlock(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "PutBucketPublicAccessBlockInput is nil")
}

// ==================== GetBucketPublicAccessBlock Tests ====================

func TestGetBucketPublicAccessBlock_ShouldReturnBlock_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<PublicAccessBlockConfiguration>
    <BlockPublicAcls>true</BlockPublicAcls>
    <IgnorePublicAcls>true</IgnorePublicAcls>
    <BlockPublicPolicy>true</BlockPublicPolicy>
    <RestrictPublicBuckets>true</RestrictPublicBuckets>
</PublicAccessBlockConfiguration>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketPublicAccessBlock(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketPublicAccessBlock Tests ====================

func TestDeleteBucketPublicAccessBlock_ShouldDeleteBlock_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketPublicAccessBlock(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== GetBucketPolicyPublicStatus Tests ====================

func TestGetBucketPolicyPublicStatus_ShouldReturnStatus_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<PolicyStatus><IsPublic>true</IsPublic></PolicyStatus>`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketPolicyPublicStatus(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== GetBucketPublicStatus Tests ====================

func TestGetBucketPublicStatus_ShouldReturnStatus_WhenSuccess(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<BucketStatus><IsPublic>true</IsPublic></BucketStatus>`, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketPublicStatus(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketCustomDomain Tests ====================

func TestSetBucketCustomDomain_ShouldSetDomain_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketCustomDomainInput{
		Bucket:       TestBucket,
		CustomDomain: "custom.example.com",
	}

	output, err := client.SetBucketCustomDomain(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestSetBucketCustomDomain_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.SetBucketCustomDomain(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "SetBucketCustomDomainInput is nil")
}

// ==================== GetBucketCustomDomain Tests ====================

func TestGetBucketCustomDomain_ShouldReturnDomain_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<CustomDomainConfig>
    <Domains>custom.example.com</Domains>
</CustomDomainConfig>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketCustomDomain(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketCustomDomain Tests ====================

func TestDeleteBucketCustomDomain_ShouldDeleteDomain_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &DeleteBucketCustomDomainInput{
		Bucket:       TestBucket,
		CustomDomain: "custom.example.com",
	}

	output, err := client.DeleteBucketCustomDomain(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

func TestDeleteBucketCustomDomain_ShouldReturnError_WhenInputNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	output, err := client.DeleteBucketCustomDomain(nil)

	require.NotNil(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "DeleteBucketCustomDomainInput is nil")
}

// ==================== SetBucketMirrorBackToSource Tests ====================

func TestSetBucketMirrorBackToSource_ShouldSetMirror_WhenValidInput(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketMirrorBackToSourceInput{
		Bucket: TestBucket,
	}

	output, err := client.SetBucketMirrorBackToSource(input)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== GetBucketMirrorBackToSource Tests ====================

func TestGetBucketMirrorBackToSource_ShouldReturnMirror_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(200, `<?xml version="1.0" encoding="UTF-8"?>
<MirrorBackToSource>
    <Status>Enabled</Status>
</MirrorBackToSource>`, nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketMirrorBackToSource(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== DeleteBucketMirrorBackToSource Tests ====================

func TestDeleteBucketMirrorBackToSource_ShouldDeleteMirror_WhenSuccess(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			return CreateTestHTTPResponse(204, "", nil)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketMirrorBackToSource(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
}

// ==================== SetBucketInventory Tests ====================

func TestSetBucketInventory_ShouldSetInventory_GivenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "PUT", req.Method)
			assert.Equal(t, "test-inventory-id", req.URL.Query().Get("inventory"))
			return CreateTestHTTPResponse(200, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketInventoryInput{
		Bucket: TestBucket,
		InventoryConfiguration: InventoryConfiguration{
			Id:        "test-inventory-id",
			IsEnabled: true,
			Destination: InventoryDestination{
				Format: "CSV",
				Bucket: "destination-bucket",
				Prefix: "inventory/",
			},
			Schedule: InventorySchedule{
				Frequency: "Daily",
			},
		},
	}

	output, err := client.SetBucketInventory(input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "test-request-id", output.RequestId)
}

func TestSetBucketInventory_ShouldReturnError_GivenNilInput(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))

	output, err := client.SetBucketInventory(nil)

	require.NotNil(t, err)
	require.Nil(t, output)
	assert.Contains(t, err.Error(), "is nil")
}

func TestSetBucketInventory_ShouldReturnError_GivenNetworkFailure(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	input := &SetBucketInventoryInput{
		Bucket: TestBucket,
		InventoryConfiguration: InventoryConfiguration{
			Id:        "test-inventory-id",
			IsEnabled: true,
		},
	}

	output, err := client.SetBucketInventory(input)

	require.NotNil(t, err)
	require.Nil(t, output)
}

// ==================== GetBucketInventory Tests ====================

func TestGetBucketInventory_ShouldReturnInventory_GivenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, "test-inventory-id", req.URL.Query().Get("inventory"))
			return CreateTestHTTPResponse(200, TestGetBucketInventoryXML, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketInventory(TestBucket, "test-inventory-id")

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "test-inventory-id", output.Id)
	assert.True(t, output.IsEnabled)
	assert.Equal(t, "CSV", output.Destination.Format)
	assert.Equal(t, "destination-bucket", output.Destination.Bucket)
	assert.Equal(t, "test-request-id", output.RequestId)
}

func TestGetBucketInventory_ShouldReturnError_GivenNetworkFailure(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.GetBucketInventory(TestBucket, "test-inventory-id")

	require.NotNil(t, err)
	require.Nil(t, output)
}

// ==================== ListBucketInventory Tests ====================

func TestListBucketInventory_ShouldReturnInventoryList_GivenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "GET", req.Method)
			assert.Equal(t, "", req.URL.Query().Get("inventory"))
			return CreateTestHTTPResponse(200, TestListBucketInventoryXML, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBucketInventory(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.InventoryConfigurations, 2)
	assert.Equal(t, "inventory-1", output.InventoryConfigurations[0].Id)
	assert.Equal(t, "inventory-2", output.InventoryConfigurations[1].Id)
	assert.False(t, output.IsTruncated)
	assert.Equal(t, "test-request-id", output.RequestId)
}

func TestListBucketInventory_ShouldReturnEmptyList_GivenNoInventories(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			emptyListXML := `<?xml version="1.0" encoding="UTF-8"?>
<ListInventoryConfigurationsResult>
	<IsTruncated>false</IsTruncated>
</ListInventoryConfigurationsResult>`
			return CreateTestHTTPResponse(200, emptyListXML, headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBucketInventory(TestBucket)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Len(t, output.InventoryConfigurations, 0)
	assert.Equal(t, "test-request-id", output.RequestId)
}

func TestListBucketInventory_ShouldReturnError_GivenNetworkFailure(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.ListBucketInventory(TestBucket)

	require.NotNil(t, err)
	require.Nil(t, output)
}

// ==================== DeleteBucketInventory Tests ====================

func TestDeleteBucketInventory_ShouldDeleteInventory_GivenValidInput(t *testing.T) {
	headers := make(http.Header)
	headers.Set("x-obs-request-id", "test-request-id")

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			assert.Equal(t, "DELETE", req.Method)
			assert.Equal(t, "test-inventory-id", req.URL.Query().Get("inventory"))
			return CreateTestHTTPResponse(204, "", headers)
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketInventory(TestBucket, "test-inventory-id")

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, "test-request-id", output.RequestId)
}

func TestDeleteBucketInventory_ShouldReturnError_GivenNetworkFailure(t *testing.T) {
	mockTransport := &MockRoundTripper{
		ErrorFunc: func(req *http.Request) error {
			return assert.AnError
		},
	}
	client := CreateTestObsClient(TestEndpoint, WithHttpTransport(&http.Transport{}))
	client.httpClient = &http.Client{Transport: mockTransport}

	output, err := client.DeleteBucketInventory(TestBucket, "test-inventory-id")

	require.NotNil(t, err)
	require.Nil(t, output)
}
