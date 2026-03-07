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
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateBucketInput_ShouldHaveNewFields_GivenCompleteInput tests the new CreateBucketInput fields
func TestCreateBucketInput_ShouldHaveNewFields_GivenCompleteInput(t *testing.T) {
	input := &CreateBucketInput{
		Bucket:                   "test-bucket",
		ACL:                      AclPublicRead,
		StorageClass:             StorageClassStandard,
		BucketType:               "OBJECT",
		SseKmsKeyId:              "test-kms-key-id",
		SseKmsKeyProjectId:       "test-project-id",
		ServerSideDataEncryption:  "AES256",
	}

	assert.Equal(t, "test-bucket", input.Bucket)
	assert.Equal(t, AclPublicRead, input.ACL)
	assert.Equal(t, StorageClassStandard, input.StorageClass)
	assert.Equal(t, "OBJECT", input.BucketType)
	assert.Equal(t, "test-kms-key-id", input.SseKmsKeyId)
	assert.Equal(t, "test-project-id", input.SseKmsKeyProjectId)
	assert.Equal(t, "AES256", input.ServerSideDataEncryption)
}

// TestPutObjectInput_ShouldHavePromotedFields_GivenCompleteInput tests that PutObjectInput correctly promotes fields from embedded types
func TestPutObjectInput_ShouldHavePromotedFields_GivenCompleteInput(t *testing.T) {
	input := &PutObjectInput{
		PutObjectBasicInput: PutObjectBasicInput{
			ObjectOperationInput: ObjectOperationInput{
				Bucket: "test-bucket",
				Key:    "test-key",
				ACL:    AclPrivate,
			},
		},
	}

	assert.Equal(t, "test-bucket", input.Bucket)
	assert.Equal(t, "test-key", input.Key)
	assert.Equal(t, AclPrivate, input.ACL)
}

// TestObjectOperationInput_ShouldHaveObjectLockFields_GivenCompleteInput tests the new ObjectLock fields
func TestObjectOperationInput_ShouldHaveObjectLockFields_GivenCompleteInput(t *testing.T) {
	input := &ObjectOperationInput{
		Bucket:                   "test-bucket",
		Key:                      "test-key",
		ACL:                      AclPrivate,
		ObjectLockMode:           "GOVERNANCE",
		ObjectLockRetainUntilDate: "2026-12-31T12:00:00Z",
	}

	assert.Equal(t, "test-bucket", input.Bucket)
	assert.Equal(t, "test-key", input.Key)
	assert.Equal(t, AclPrivate, input.ACL)
	assert.Equal(t, "GOVERNANCE", input.ObjectLockMode)
	assert.Equal(t, "2026-12-31T12:00:00Z", input.ObjectLockRetainUntilDate)
}

// TestListObjectsInput_ShouldHaveEncodingType_GivenCompleteInput tests the EncodingType field
func TestListObjectsInput_ShouldHaveEncodingType_GivenCompleteInput(t *testing.T) {
	input := &ListObjectsInput{
		Bucket:       "test-bucket",
		Marker:       "test-marker",
		EncodingType: "url",
	}

	assert.Equal(t, "test-bucket", input.Bucket)
	assert.Equal(t, "test-marker", input.Marker)
	assert.Equal(t, "url", input.EncodingType)
}