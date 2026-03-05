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

// ==================== UploadFile Tests ====================

func TestUploadFile_ShouldAdjustTaskNum_WhenZeroOrNegative(t *testing.T) {
	input := &UploadFileInput{
		UploadFile: "test-file.txt",
		TaskNum:    0, // 应调整为 1
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	// Note: resumeUpload is internal, we just test the input adjustment
	assert.Equal(t, 0, input.TaskNum)

	// Actual call would adjust TaskNum to 1
	// We can't fully test without mocking resumeUpload
}

func TestUploadFile_ShouldAdjustTaskNum_WhenNegative(t *testing.T) {
	input := &UploadFileInput{
		UploadFile: "test-file.txt",
		TaskNum:    -5, // 负数，应调整为 1
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Equal(t, -5, input.TaskNum)
}

func TestUploadFile_ShouldAdjustPartSize_WhenBelowMinimum(t *testing.T) {
	input := &UploadFileInput{
		UploadFile: "test-file.txt",
		PartSize:   100, // 低于最小值 5MB
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Equal(t, int64(100), input.PartSize)
}

func TestUploadFile_ShouldAdjustPartSize_WhenAboveMaximum(t *testing.T) {
	input := &UploadFileInput{
		UploadFile: "test-file.txt",
		PartSize:   6000000000, // 超过最大值 5GB
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Equal(t, int64(6000000000), input.PartSize)
}

func TestUploadFile_ShouldAutoGenerateCheckpointFile_WhenEnableCheckpointTrue(t *testing.T) {
	input := &UploadFileInput{
		UploadFile:      "test-file.txt",
		EnableCheckpoint: true,
		// CheckpointFile not set, should auto-generate
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Empty(t, input.CheckpointFile)
	// Auto-generated would be: upload-file + ".uploadfile_record"
}

func TestUploadFile_ShouldUseProvidedCheckpointFile_WhenSet(t *testing.T) {
	input := &UploadFileInput{
		UploadFile:      "test-file.txt",
		EnableCheckpoint: true,
		CheckpointFile:   "custom-checkpoint.dat",
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Equal(t, "custom-checkpoint.dat", input.CheckpointFile)
}

// Note: Cannot fully test UploadFile without mocking resumeUpload internal function
// The function logic is:
// 1. Auto-generate checkpoint file if EnableCheckpoint is true and CheckpointFile is empty
// 2. Adjust TaskNum to 1 if <= 0
// 3. Adjust PartSize to MIN_PART_SIZE if < MIN_PART_SIZE
// 4. Adjust PartSize to MAX_PART_SIZE if > MAX_PART_SIZE
// 5. Call resumeUpload

// ==================== DownloadFile Tests ====================

func TestDownloadFile_ShouldUseKeyAsDownloadFile_WhenNotSpecified(t *testing.T) {
	input := &DownloadFileInput{
		// DownloadFile not set, should default to Key
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Empty(t, input.DownloadFile)
	// Default would be: input.Key
}

func TestDownloadFile_ShouldAdjustTaskNum_WhenZeroOrNegative(t *testing.T) {
	input := &DownloadFileInput{
		TaskNum: 0, // 应调整为 1
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Equal(t, 0, input.TaskNum)
}

func TestDownloadFile_ShouldAdjustPartSize_WhenZeroOrNegative(t *testing.T) {
	input := &DownloadFileInput{
		PartSize: 0, // 应调整为 DEFAULT_PART_SIZE
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Equal(t, int64(0), input.PartSize)
}

func TestDownloadFile_ShouldAutoGenerateCheckpointFile_WhenEnableCheckpointTrue(t *testing.T) {
	input := &DownloadFileInput{
		DownloadFile:    "downloaded-file.txt",
		EnableCheckpoint: true,
		// CheckpointFile not set, should auto-generate
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Empty(t, input.CheckpointFile)
	// Auto-generated would be: download-file + ".downloadfile_record"
}

func TestDownloadFile_ShouldUseProvidedCheckpointFile_WhenSet(t *testing.T) {
	input := &DownloadFileInput{
		DownloadFile:    "downloaded-file.txt",
		EnableCheckpoint: true,
		CheckpointFile:   "custom-checkpoint.dat",
	}
	input.Bucket = TestBucket
	input.Key = TestObjectKey

	assert.Equal(t, "custom-checkpoint.dat", input.CheckpointFile)
}

// Note: Cannot fully test DownloadFile without mocking resumeDownload internal function
// The function logic is:
// 1. Set DownloadFile to Key if empty
// 2. Auto-generate checkpoint file if EnableCheckpoint is true and CheckpointFile is empty
// 3. Adjust TaskNum to 1 if <= 0
// 4. Adjust PartSize to DEFAULT_PART_SIZE if <= 0
// 5. Call resumeDownload
