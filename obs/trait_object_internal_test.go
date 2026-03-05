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

// CopyObjectInput prepareReplaceHeaders Tests

func TestCopyObjectInput_PrepareReplaceHeaders_ShouldSetAllHeaders_WhenAllFieldsAreSet(t *testing.T) {
	input := CopyObjectInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: "source-bucket",
			Key:    "source-key",
		},
		HttpHeader: HttpHeader{
			ContentType:        "application/json",
			CacheControl:       "no-cache",
			ContentDisposition: "attachment; filename=test.txt",
			ContentEncoding:    "gzip",
			ContentLanguage:    "en",
		},
		Expires: "Fri, 31 Dec 2025 23:59:59 GMT",
	}

	headers := make(map[string][]string)
	input.prepareReplaceHeaders(headers)

	assert.Equal(t, "application/json", headers[HEADER_CONTENT_TYPE][0])
	assert.Equal(t, "no-cache", headers[HEADER_CACHE_CONTROL][0])
	assert.Equal(t, "attachment; filename=test.txt", headers[HEADER_CONTENT_DISPOSITION][0])
	assert.Equal(t, "gzip", headers[HEADER_CONTENT_ENCODING][0])
	assert.Equal(t, "en", headers[HEADER_CONTENT_LANGUAGE][0])
	assert.Equal(t, "Fri, 31 Dec 2025 23:59:59 GMT", headers[HEADER_EXPIRES_CAMEL][0])
}

func TestCopyObjectInput_PrepareReplaceHeaders_ShouldSetContentTypeOnly_WhenOnlyContentTypeIsSet(t *testing.T) {
	input := CopyObjectInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: "source-bucket",
			Key:    "source-key",
		},
		HttpHeader: HttpHeader{
			ContentType: "image/jpeg",
		},
	}

	headers := make(map[string][]string)
	input.prepareReplaceHeaders(headers)

	assert.Equal(t, "image/jpeg", headers[HEADER_CONTENT_TYPE][0])
	assert.Len(t, headers[HEADER_CACHE_CONTROL], 0)
	assert.Len(t, headers[HEADER_CONTENT_DISPOSITION], 0)
	assert.Len(t, headers[HEADER_CONTENT_ENCODING], 0)
	assert.Len(t, headers[HEADER_CONTENT_LANGUAGE], 0)
	assert.Len(t, headers[HEADER_EXPIRES_CAMEL], 0)
}

func TestCopyObjectInput_PrepareReplaceHeaders_ShouldSetExpiresHeader_WhenExpiresIsSet(t *testing.T) {
	input := CopyObjectInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: "source-bucket",
			Key:    "source-key",
		},
		HttpHeader: HttpHeader{
			HttpExpires: "Wed, 22 Oct 2025 07:28:00 GMT", // Should be ignored when Expires is set
		},
		Expires: "Wed, 21 Oct 2025 07:28:00 GMT",
	}

	headers := make(map[string][]string)
	input.prepareReplaceHeaders(headers)

	assert.Equal(t, "Wed, 21 Oct 2025 07:28:00 GMT", headers[HEADER_EXPIRES_CAMEL][0])
}

func TestCopyObjectInput_PrepareReplaceHeaders_ShouldSetHttpExpiresHeader_WhenOnlyHttpExpiresIsSet(t *testing.T) {
	input := CopyObjectInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: "source-bucket",
			Key:    "source-key",
		},
		HttpHeader: HttpHeader{
			HttpExpires: "Thu, 23 Oct 2025 07:28:00 GMT",
		},
	}

	headers := make(map[string][]string)
	input.prepareReplaceHeaders(headers)

	assert.Equal(t, "Thu, 23 Oct 2025 07:28:00 GMT", headers[HEADER_EXPIRES_CAMEL][0])
}

func TestCopyObjectInput_PrepareReplaceHeaders_ShouldNotSetAnyHeader_WhenAllFieldsAreEmpty(t *testing.T) {
	input := CopyObjectInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: "source-bucket",
			Key:    "source-key",
		},
		HttpHeader: HttpHeader{},
	}

	headers := make(map[string][]string)
	input.prepareReplaceHeaders(headers)

	assert.Len(t, headers, 0)
}

func TestCopyObjectInput_PrepareReplaceHeaders_ShouldOverrideExistingHeaders_WhenHeadersAlreadyExist(t *testing.T) {
	input := CopyObjectInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: "source-bucket",
			Key:    "source-key",
		},
		HttpHeader: HttpHeader{
			ContentType:  "text/plain",
			CacheControl: "public, max-age=3600",
		},
	}

	headers := map[string][]string{
		HEADER_CONTENT_TYPE:  {"application/octet-stream"},
		HEADER_CACHE_CONTROL: {"no-cache"},
		HEADER_EXPIRES_CAMEL:  {"old-expire-value"},
	}

	input.prepareReplaceHeaders(headers)

	assert.Equal(t, "text/plain", headers[HEADER_CONTENT_TYPE][0])
	assert.Equal(t, "public, max-age=3600", headers[HEADER_CACHE_CONTROL][0])
	// The function does not remove headers not set in input
	assert.Equal(t, "old-expire-value", headers[HEADER_EXPIRES_CAMEL][0])
}

// SetObjectMetadataInput.prepareStorageClass Tests

func TestSetObjectMetadataInputPrepareStorageClass_ShouldSetStandardClass_WhenIsObsIsTrue(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassStandard,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, true)

	assert.Contains(t, headers, "x-obs-storage-class")
	assert.Equal(t, []string{"STANDARD"}, headers["x-obs-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldSetStandardClass_WhenIsObsIsFalse(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassStandard,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, false)

	assert.Contains(t, headers, "x-amz-storage-class")
	assert.Equal(t, []string{"STANDARD"}, headers["x-amz-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldConvertWarmToStandardIA_WhenIsObsIsFalse(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassWarm,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, false)

	assert.Contains(t, headers, "x-amz-storage-class")
	assert.Equal(t, []string{"STANDARD_IA"}, headers["x-amz-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldKeepWarm_WhenIsObsIsTrue(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassWarm,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, true)

	assert.Contains(t, headers, "x-obs-storage-class")
	assert.Equal(t, []string{"WARM"}, headers["x-obs-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldConvertColdToGlacier_WhenIsObsIsFalse(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassCold,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, false)

	assert.Contains(t, headers, "x-amz-storage-class")
	assert.Equal(t, []string{"GLACIER"}, headers["x-amz-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldKeepCold_WhenIsObsIsTrue(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassCold,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, true)

	assert.Contains(t, headers, "x-obs-storage-class")
	assert.Equal(t, []string{"COLD"}, headers["x-obs-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldSetDeepArchive_WhenIsObsIsTrue(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassDeepArchive,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, true)

	assert.Contains(t, headers, "x-obs-storage-class")
	assert.Equal(t, []string{"DEEP_ARCHIVE"}, headers["x-obs-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldSetIntelligentTiering_WhenIsObsIsTrue(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: StorageClassIntelligentTiering,
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, true)

	assert.Contains(t, headers, "x-obs-storage-class")
	assert.Equal(t, []string{"INTELLIGENT_TIERING"}, headers["x-obs-storage-class"])
}

func TestSetObjectMetadataInputPrepareStorageClass_ShouldDoNothing_WhenStorageClassIsEmpty(t *testing.T) {
	input := SetObjectMetadataInput{
		Bucket:       "test-bucket",
		Key:          "test-key",
		StorageClass: "",
	}
	headers := make(map[string][]string)

	input.prepareStorageClass(headers, true)

	assert.NotContains(t, headers, "x-obs-storage-class")
	assert.NotContains(t, headers, "x-amz-storage-class")
}

// ObjectOperationInput.trans Tests

func TestObjectOperationInputTrans_ShouldSetACL_WhenACLIsSet(t *testing.T) {
	input := ObjectOperationInput{
		ACL: AclPublicRead,
	}

	_, headers, _, _ := input.trans(true)

	assert.Contains(t, headers, "x-obs-acl")
	assert.Equal(t, []string{"public-read"}, headers["x-obs-acl"])
}

func TestObjectOperationInputTrans_ShouldSetStorageClass_WhenIsObs(t *testing.T) {
	input := ObjectOperationInput{
		StorageClass: StorageClassWarm,
	}

	_, headers, _, _ := input.trans(true)

	assert.Contains(t, headers, "x-obs-storage-class")
	assert.Equal(t, []string{"WARM"}, headers["x-obs-storage-class"])
}

func TestObjectOperationInputTrans_ShouldConvertWarmToStandardIA_WhenNotIsObs(t *testing.T) {
	input := ObjectOperationInput{
		StorageClass: StorageClassWarm,
	}

	_, headers, _, _ := input.trans(false)

	assert.Contains(t, headers, "x-amz-storage-class")
	assert.Equal(t, []string{"STANDARD_IA"}, headers["x-amz-storage-class"])
}

func TestObjectOperationInputTrans_ShouldConvertColdToGlacier_WhenNotIsObs(t *testing.T) {
	input := ObjectOperationInput{
		StorageClass: StorageClassCold,
	}

	_, headers, _, _ := input.trans(false)

	assert.Contains(t, headers, "x-amz-storage-class")
	assert.Equal(t, []string{"GLACIER"}, headers["x-amz-storage-class"])
}

func TestObjectOperationInputTrans_ShouldSetWebsiteRedirectLocation(t *testing.T) {
	input := ObjectOperationInput{
		WebsiteRedirectLocation: "https://example.com/redirect",
	}

	_, headers, _, _ := input.trans(true)

	assert.Contains(t, headers, "x-obs-website-redirect-location")
	assert.Equal(t, []string{"https://example.com/redirect"}, headers["x-obs-website-redirect-location"])
}

func TestObjectOperationInputTrans_ShouldSetExpires(t *testing.T) {
	input := ObjectOperationInput{
		Expires: 1234567890,
	}

	_, headers, _, _ := input.trans(true)

	assert.Contains(t, headers, "x-obs-expires")
	assert.Equal(t, []string{"1234567890"}, headers["x-obs-expires"])
}

func TestObjectOperationInputTrans_ShouldSetMetadata(t *testing.T) {
	input := ObjectOperationInput{
		Metadata: map[string]string{
			"custom-key": "custom-value",
		},
	}

	_, headers, _, _ := input.trans(true)

	assert.Contains(t, headers, "x-obs-meta-custom-key")
}

func TestObjectOperationInputTrans_ShouldReturnEmptyHeaders_WhenNoFieldsSet(t *testing.T) {
	input := ObjectOperationInput{}

	_, headers, _, _ := input.trans(true)

	assert.Empty(t, headers)
}

func TestObjectOperationInputTrans_ShouldSetGrantHeaders(t *testing.T) {
	input := ObjectOperationInput{
		GrantReadId:       "read-id",
		GrantReadAcpId:    "read-acp-id",
		GrantWriteAcpId:   "write-acp-id",
		GrantFullControlId: "full-control-id",
	}

	_, headers, _, _ := input.trans(true)

	assert.Contains(t, headers, "x-obs-grant-read")
	assert.Contains(t, headers, "x-obs-grant-read-acp")
	assert.Contains(t, headers, "x-obs-grant-write-acp")
	assert.Contains(t, headers, "x-obs-grant-full-control")
}

