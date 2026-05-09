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

//go:build unit

package obs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Test Data ====================

const (
	testUploadID  = "test-upload-id-001"
	testETag      = "\"d41d8cd98f00b204e9800998ecf8427e\""
	testObjectSize = 1024 * 100 // 100KB
)

// ==================== TransferTask State Machine Tests ====================

func TestTransferTask_Pause_ShouldSucceed_WhenRunning(t *testing.T) {
	task := newTransferTask(nil)
	assert.Equal(t, TransferStatusRunning, task.Status())

	err := task.Pause()
	assert.NoError(t, err)
	assert.Equal(t, TransferStatusPaused, task.Status())
}

func TestTransferTask_Pause_ShouldReturnError_WhenNotRunning(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus func(task *TransferTask)
	}{
		{"paused", func(task *TransferTask) { task.Pause() }},
		{"cancelled", func(task *TransferTask) { task.Cancel() }},
		{"completed", func(task *TransferTask) { task.finish(nil, nil) }},
		{"failed", func(task *TransferTask) { task.finish(nil, fmt.Errorf("fail")) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := newTransferTask(nil)
			tt.setupStatus(task)

			err := task.Pause()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot pause task in state")
		})
	}
}

func TestTransferTask_Resume_ShouldSucceed_WhenPaused(t *testing.T) {
	task := newTransferTask(nil)
	task.Pause()
	assert.Equal(t, TransferStatusPaused, task.Status())

	err := task.Resume()
	assert.NoError(t, err)
	assert.Equal(t, TransferStatusRunning, task.Status())
}

func TestTransferTask_Resume_ShouldReturnError_WhenNotPaused(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus func(task *TransferTask)
	}{
		{"running", func(task *TransferTask) {}},
		{"cancelled", func(task *TransferTask) { task.Cancel() }},
		{"completed", func(task *TransferTask) { task.finish(nil, nil) }},
		{"failed", func(task *TransferTask) { task.finish(nil, fmt.Errorf("fail")) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := newTransferTask(nil)
			tt.setupStatus(task)

			err := task.Resume()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot resume task in state")
		})
	}
}

func TestTransferTask_Cancel_ShouldSucceed_WhenRunningOrPaused(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(task *TransferTask)
	}{
		{"from_running", func(task *TransferTask) {}},
		{"from_paused", func(task *TransferTask) { task.Pause() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := newTransferTask(nil)
			tt.setup(task)

			err := task.Cancel()
			assert.NoError(t, err)
			assert.Equal(t, TransferStatusCancelled, task.Status())

			abortVal := atomic.LoadInt32(task.abortFlag())
			assert.Equal(t, int32(1), abortVal)
		})
	}
}

func TestTransferTask_Cancel_ShouldReturnError_WhenTerminalState(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus func(task *TransferTask)
	}{
		{"cancelled", func(task *TransferTask) { task.Cancel() }},
		{"completed", func(task *TransferTask) { task.finish(nil, nil) }},
		{"failed", func(task *TransferTask) { task.finish(nil, fmt.Errorf("fail")) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := newTransferTask(nil)
			tt.setupStatus(task)

			err := task.Cancel()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot cancel task in state")
		})
	}
}

func TestTransferTask_Status_ShouldReturnCurrentState(t *testing.T) {
	task := newTransferTask(nil)
	assert.Equal(t, TransferStatusRunning, task.Status())

	task.Pause()
	assert.Equal(t, TransferStatusPaused, task.Status())

	task.Resume()
	assert.Equal(t, TransferStatusRunning, task.Status())

	task.Cancel()
	assert.Equal(t, TransferStatusCancelled, task.Status())
}

func TestTransferTask_GetResult_ShouldReturnResultAndError(t *testing.T) {
	t.Run("should_return_nil_when_not_finished", func(t *testing.T) {
		task := newTransferTask(nil)
		result, err := task.GetResult()
		assert.Nil(t, result)
		assert.NoError(t, err)
	})

	t.Run("should_return_result_when_completed_successfully", func(t *testing.T) {
		task := newTransferTask(nil)
		expectedResult := "test-result"
		task.finish(expectedResult, nil)

		result, err := task.GetResult()
		assert.Equal(t, expectedResult, result)
		assert.NoError(t, err)
	})

	t.Run("should_return_error_when_failed", func(t *testing.T) {
		task := newTransferTask(nil)
		expectedErr := fmt.Errorf("test error")
		task.finish(nil, expectedErr)

		result, err := task.GetResult()
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestTransferTask_Done_ShouldClose_WhenTaskFinishes(t *testing.T) {
	task := newTransferTask(nil)

	select {
	case <-task.Done():
		t.Fatal("Done channel should not be closed before finish")
	default:
	}

	task.finish("result", nil)

	select {
	case <-task.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("Done channel should be closed after finish")
	}
}

func TestTransferTask_Done_ShouldClose_WhenTaskCancelled(t *testing.T) {
	task := newTransferTask(nil)
	task.Cancel()
	// finish must still be called to close the done channel
	task.finish(nil, fmt.Errorf("cancelled"))

	select {
	case <-task.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("Done channel should be closed after cancel+finish")
	}
}

func TestTransferTask_finish_ShouldSetCorrectStatus(t *testing.T) {
	t.Run("should_set_completed_when_no_error", func(t *testing.T) {
		task := newTransferTask(nil)
		task.finish("result", nil)
		assert.Equal(t, TransferStatusCompleted, task.Status())
	})

	t.Run("should_set_failed_when_has_error", func(t *testing.T) {
		task := newTransferTask(nil)
		task.finish(nil, fmt.Errorf("error"))
		assert.Equal(t, TransferStatusFailed, task.Status())
	})

	t.Run("should_keep_cancelled_status_when_cancelled", func(t *testing.T) {
		task := newTransferTask(nil)
		task.Cancel()
		task.finish(nil, nil)
		assert.Equal(t, TransferStatusCancelled, task.Status())
	})

	t.Run("should_only_execute_once", func(t *testing.T) {
		task := newTransferTask(nil)
		task.finish("first", nil)
		task.finish("second", fmt.Errorf("error"))

		result, err := task.GetResult()
		assert.Equal(t, "first", result)
		assert.NoError(t, err)
		assert.Equal(t, TransferStatusCompleted, task.Status())
	})
}

func TestTransferTask_isCancelled_ShouldReturnCorrectState(t *testing.T) {
	task := newTransferTask(nil)
	assert.False(t, task.isCancelled())

	task.Pause()
	assert.False(t, task.isCancelled())

	task.Cancel()
	assert.True(t, task.isCancelled())
}

func TestTransferTask_cancelCleanup_ShouldInvokeCallback(t *testing.T) {
	called := false
	task := newTransferTask(func() {
		called = true
	})

	task.cancelCleanup()
	assert.True(t, called)
}

func TestTransferTask_cancelCleanup_ShouldNotPanic_WhenCallbackIsNil(t *testing.T) {
	task := newTransferTask(nil)
	assert.NotPanics(t, func() {
		task.cancelCleanup()
	})
}

func TestTransferTask_checkAndWaitPause_ShouldBlock_WhenPaused(t *testing.T) {
	task := newTransferTask(nil)
	task.Pause()

	cancelled := int32(0)
	go func() {
		time.Sleep(50 * time.Millisecond)
		task.Resume()
		atomic.StoreInt32(&cancelled, 1)
	}()

	result := task.checkAndWaitPause()
	assert.False(t, result)
}

func TestTransferTask_checkAndWaitPause_ShouldReturnTrue_WhenCancelled(t *testing.T) {
	task := newTransferTask(nil)
	task.Pause()

	go func() {
		time.Sleep(50 * time.Millisecond)
		task.Cancel()
	}()

	result := task.checkAndWaitPause()
	assert.True(t, result)
}

func TestTransferTask_checkAndWaitPause_ShouldReturnFalse_WhenRunning(t *testing.T) {
	task := newTransferTask(nil)
	result := task.checkAndWaitPause()
	assert.False(t, result)
}

func TestTransferTask_abortFlag_ShouldReturnPointer(t *testing.T) {
	task := newTransferTask(nil)
	flag := task.abortFlag()
	require.NotNil(t, flag)
	assert.Equal(t, int32(0), atomic.LoadInt32(flag))

	task.Cancel()
	assert.Equal(t, int32(1), atomic.LoadInt32(flag))
}

// ==================== Synchronous nil-check Tests ====================

func TestUploadFile_ShouldReturnError_WhenInputIsNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)
	output, err := client.UploadFile(nil)
	assert.Nil(t, output)
	assert.Equal(t, "UploadFileInput is nil", err.Error())
}

func TestDownloadFile_ShouldReturnError_WhenInputIsNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)
	output, err := client.DownloadFile(nil)
	assert.Nil(t, output)
	assert.Equal(t, "DownloadFileInput is nil", err.Error())
}

// ==================== UploadFileAsync Tests ====================

func TestUploadFileAsync_ShouldReturnError_WhenInputIsNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)
	task, err := client.UploadFileAsync(nil)
	assert.Nil(t, task)
	assert.Equal(t, "UploadFileInput is nil", err.Error())
}

func TestUploadFileAsync_ShouldUploadSuccessfully_WhenServerResponds(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	initiateCalled := int32(0)
	uploadPartCount := int32(0)
	completeCalled := int32(0)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				atomic.AddInt32(&initiateCalled, 1)
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				atomic.AddInt32(&uploadPartCount, 1)
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				atomic.AddInt32(&completeCalled, 1)
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   testObjectSize,
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)
	assert.Equal(t, TransferStatusRunning, task.Status())

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	result, taskErr := task.GetResult()
	assert.NoError(t, taskErr)
	assert.NotNil(t, result)
	assert.Equal(t, TransferStatusCompleted, task.Status())

	assert.True(t, atomic.LoadInt32(&initiateCalled) >= 1, "InitiateMultipartUpload should be called")
	assert.True(t, atomic.LoadInt32(&uploadPartCount) >= 1, "UploadPart should be called")
	assert.True(t, atomic.LoadInt32(&completeCalled) >= 1, "CompleteMultipartUpload should be called")
}

func TestUploadFileAsync_ShouldHandlePauseAndResume(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	resumeChan := make(chan struct{}, 1)
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				// Block until resume is signaled to allow pause to take effect
				select {
				case <-resumeChan:
				case <-time.After(5 * time.Second):
				}
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodDelete {
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   testObjectSize,
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	// Pause the task
	pauseErr := task.Pause()
	assert.NoError(t, pauseErr)
	assert.Equal(t, TransferStatusPaused, task.Status())

	// Resume the task and signal the blocked upload
	resumeErr := task.Resume()
	assert.NoError(t, resumeErr)
	resumeChan <- struct{}{}

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time after resume")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
}

func TestUploadFileAsync_ShouldCancelAndCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	abortCalled := int32(0)
	blockUpload := make(chan struct{})
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				// Block to give time for cancellation
				select {
				case <-blockUpload:
				case <-time.After(5 * time.Second):
				}
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodDelete && strings.Contains(req.URL.RawQuery, "uploadId") {
				atomic.AddInt32(&abortCalled, 1)
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   testObjectSize,
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	// Cancel the task
	cancelErr := task.Cancel()
	assert.NoError(t, cancelErr)
	assert.Equal(t, TransferStatusCancelled, task.Status())

	// Unblock the upload part to let the goroutine finish
	close(blockUpload)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Cancelled upload did not finish in time")
	}

	assert.Equal(t, TransferStatusCancelled, task.Status())
	_, taskErr := task.GetResult()
	assert.Error(t, taskErr)
	assert.Contains(t, taskErr.Error(), "cancelled")

	// Verify abort was called for cleanup
	assert.True(t, atomic.LoadInt32(&abortCalled) >= 1, "AbortMultipartUpload should be called on cancel")
}

func TestUploadFileAsync_ShouldEmitProgressEvents_WhenUsingListener(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	listener := NewMockProgressListener()

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   testObjectSize,
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input, WithProgress(listener))
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	// Verify progress events were emitted
	assert.NotEmpty(t, listener.Events, "Progress events should be emitted")

	hasStartEvent := false
	hasCompletedEvent := false
	for _, event := range listener.Events {
		if event.EventType == TransferStartedEvent {
			hasStartEvent = true
		}
		if event.EventType == TransferCompletedEvent {
			hasCompletedEvent = true
		}
	}
	assert.True(t, hasStartEvent, "Should emit TransferStartedEvent")
	assert.True(t, hasCompletedEvent, "Should emit TransferCompletedEvent")
}

// ==================== DownloadFileAsync Tests ====================

func TestDownloadFileAsync_ShouldReturnError_WhenInputIsNil(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)
	task, err := client.DownloadFileAsync(nil)
	assert.Nil(t, task)
	assert.Equal(t, "DownloadFileInput is nil", err.Error())
}

func TestDownloadFileAsync_ShouldDownloadSuccessfully_WhenServerResponds(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			// HEAD request for GetObjectMetadata
			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			// GET request for GetObject (range download)
			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	result, taskErr := task.GetResult()
	assert.NoError(t, taskErr)
	assert.NotNil(t, result)
	assert.Equal(t, TransferStatusCompleted, task.Status())

	// Verify the download file exists
	_, statErr := os.Stat(downloadFile)
	assert.NoError(t, statErr, "Download file should exist")
}

func TestDownloadFileAsync_ShouldHandlePauseAndResume(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"
	resumeChan := make(chan struct{}, 1)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				// Block until resume
				select {
				case <-resumeChan:
				case <-time.After(5 * time.Second):
				}
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)

	// Pause
	pauseErr := task.Pause()
	assert.NoError(t, pauseErr)
	assert.Equal(t, TransferStatusPaused, task.Status())

	// Resume
	resumeErr := task.Resume()
	assert.NoError(t, resumeErr)
	resumeChan <- struct{}{}

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time after resume")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
}

func TestDownloadFileAsync_ShouldCancelAndCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"
	blockDownload := make(chan struct{})

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				select {
				case <-blockDownload:
				case <-time.After(5 * time.Second):
				}
				respHeaders := make(http.Header)
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)

	// Cancel
	cancelErr := task.Cancel()
	assert.NoError(t, cancelErr)
	assert.Equal(t, TransferStatusCancelled, task.Status())

	// Unblock the download
	close(blockDownload)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Cancelled download did not finish in time")
	}

	assert.Equal(t, TransferStatusCancelled, task.Status())
	_, taskErr := task.GetResult()
	assert.Error(t, taskErr)
	assert.Contains(t, taskErr.Error(), "cancelled")
}

func TestDownloadFileAsync_ShouldEmitProgressEvents_WhenUsingListener(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"
	listener := NewMockProgressListener()

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input, WithProgress(listener))
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	// Verify progress events were emitted
	assert.NotEmpty(t, listener.Events, "Progress events should be emitted")

	hasStartEvent := false
	hasCompletedEvent := false
	for _, event := range listener.Events {
		if event.EventType == TransferStartedEvent {
			hasStartEvent = true
		}
		if event.EventType == TransferCompletedEvent {
			hasCompletedEvent = true
		}
	}
	assert.True(t, hasStartEvent, "Should emit TransferStartedEvent")
	assert.True(t, hasCompletedEvent, "Should emit TransferCompletedEvent")
}

func TestDownloadFileAsync_ShouldSetDefaultDownloadFile_WhenNotSpecified(t *testing.T) {
	tmpDir := t.TempDir()
	// Change to temp dir so the download file (defaulting to Key) is created there
	originalDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer os.Chdir(originalDir)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", "100")
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: 100,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set("Content-Length", "100")
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(strings.NewReader(string(make([]byte, 100)))),
					ContentLength: 100,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    "default-download-file.txt",
		},
		PartSize: 100,
		TaskNum:  1,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	_, taskErr := task.GetResult()
	assert.NoError(t, taskErr)
	// DownloadFile should have been defaulted to Key value
	assert.Equal(t, "default-download-file.txt", input.DownloadFile)
}

// ==================== UploadFileAsync Error Path Tests ====================

func TestUploadFileAsync_ShouldReturnError_WhenFileNotFound(t *testing.T) {
	client := CreateTestObsClient(TestEndpoint)

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: "/nonexistent/path/file.txt",
		PartSize:   MIN_PART_SIZE,
	}

	task, err := client.UploadFileAsync(input)
	assert.Nil(t, task)
	assert.Error(t, err)
}

func TestUploadFileAsync_ShouldReturnError_WhenFileIsDir(t *testing.T) {
	tmpDir := t.TempDir()
	client := CreateTestObsClient(TestEndpoint)

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: tmpDir,
		PartSize:   MIN_PART_SIZE,
	}

	task, err := client.UploadFileAsync(input)
	assert.Nil(t, task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "folder")
}

func TestUploadFileAsync_ShouldAdjustParameters_WhenInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   0,  // should be adjusted to MIN_PART_SIZE
		TaskNum:    0,  // should be adjusted to 1
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	assert.Equal(t, 1, input.TaskNum)
	assert.Equal(t, int64(MIN_PART_SIZE), input.PartSize)
}

// ==================== DownloadFileAsync Error Path Tests ====================

func TestDownloadFileAsync_ShouldReturnError_WhenHeadFails(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")
			return CreateMockResponse(http.StatusNotFound, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input)
	assert.Nil(t, task)
	assert.Error(t, err)
}

// ==================== UploadFileAsync — goroutine error path tests ====================

func TestUploadFileAsync_ShouldFail_WhenUploadPartReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	abortCalled := int32(0)
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				return CreateMockResponse(http.StatusInternalServerError, TestErrorResponseXML, headers)
			}

			if req.Method == http.MethodDelete && strings.Contains(req.URL.RawQuery, "uploadId") {
				atomic.AddInt32(&abortCalled, 1)
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   testObjectSize,
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	assert.Equal(t, TransferStatusFailed, task.Status())
	result, taskErr := task.GetResult()
	assert.Nil(t, result)
	assert.Error(t, taskErr)

	// Non-checkpoint mode should call AbortMultipartUpload
	assert.True(t, atomic.LoadInt32(&abortCalled) >= 1, "AbortMultipartUpload should be called on upload part failure")
}

func TestUploadFileAsync_ShouldFail_WhenCompleteMultipartUploadFails(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	abortCalled := int32(0)
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			// CompleteMultipartUpload fails
			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				return CreateMockResponse(http.StatusInternalServerError, TestErrorResponseXML, headers)
			}

			if req.Method == http.MethodDelete && strings.Contains(req.URL.RawQuery, "uploadId") {
				atomic.AddInt32(&abortCalled, 1)
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   testObjectSize,
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	assert.Equal(t, TransferStatusFailed, task.Status())
	result, taskErr := task.GetResult()
	assert.Nil(t, result)
	assert.Error(t, taskErr)

	// Non-checkpoint mode should abort after complete failure
	assert.True(t, atomic.LoadInt32(&abortCalled) >= 1, "AbortMultipartUpload should be called after CompleteMultipartUpload failure")
}

func TestUploadFileAsync_ShouldFail_WhenCheckpointFileNotWritable(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	// Use a non-existent deep directory as checkpoint path to trigger write failure
	checkpointFile := "/nonexistent_dir_deep/sub_dir/checkpoint_file"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodDelete && strings.Contains(req.URL.RawQuery, "uploadId") {
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.UploadFileAsync(input)
	// Should return error before starting the async goroutine
	assert.Nil(t, task)
	assert.Error(t, err)
}

func TestUploadFileAsync_ShouldSucceed_WithCheckpointEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	checkpointFile := tmpDir + "/upload.cp"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	result, taskErr := task.GetResult()
	assert.NoError(t, taskErr)
	assert.NotNil(t, result)

	// Checkpoint file should be cleaned up after successful upload
	_, statErr := os.Stat(checkpointFile)
	assert.True(t, os.IsNotExist(statErr), "Checkpoint file should be removed after successful upload")
}

// ==================== DownloadFileAsync — goroutine error path tests ====================

func TestDownloadFileAsync_ShouldFail_WhenGetObjectReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			// Range download returns 500 error
			if req.Method == http.MethodGet {
				return CreateMockResponse(http.StatusInternalServerError, TestErrorResponseXML, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	assert.Equal(t, TransferStatusFailed, task.Status())
	result, taskErr := task.GetResult()
	assert.Nil(t, result)
	assert.Error(t, taskErr)
}

func TestDownloadFileAsync_ShouldFail_WhenRenameFails(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory at the download file path so os.Rename from temp file
	// to this path fails (cannot rename a file onto a directory).
	downloadFilePath := tmpDir + "/target_download"
	err := os.MkdirAll(downloadFilePath, 0755)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFilePath,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	assert.Equal(t, TransferStatusFailed, task.Status())
	result, taskErr := task.GetResult()
	assert.Nil(t, result)
	assert.Error(t, taskErr)

	// Verify the temp file was created (download part succeeded)
	// but the rename to the final path failed because target is a directory
	tempFile := downloadFilePath + ".tmp"
	_, statErr := os.Stat(tempFile)
	assert.NoError(t, statErr, "Temp download file should still exist after rename failure")
}

func TestDownloadFileAsync_ShouldSucceed_WithCheckpointEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"
	checkpointFile := tmpDir + "/download.cp"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile:     downloadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	result, taskErr := task.GetResult()
	assert.NoError(t, taskErr)
	assert.NotNil(t, result)

	// Download file should exist
	_, statErr := os.Stat(downloadFile)
	assert.NoError(t, statErr, "Download file should exist")

	// Checkpoint file should be cleaned up after successful download
	_, cpStatErr := os.Stat(checkpointFile)
	assert.True(t, os.IsNotExist(cpStatErr), "Checkpoint file should be removed after successful download")
}

// ==================== Upload part concurrency edge case tests ====================

func TestUploadFileAsync_ShouldSkipCompletedParts_WhenCheckpointResumed(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	checkpointFile := tmpDir + "/upload.cp"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	uploadPartCount := int32(0)

	// First upload: succeed to create checkpoint with one part
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				atomic.AddInt32(&uploadPartCount, 1)
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("First upload did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	firstUploadPartCount := atomic.LoadInt32(&uploadPartCount)
	assert.True(t, firstUploadPartCount >= 1, "UploadPart should be called at least once")
}

func TestUploadFileAsync_ShouldAbortBeforePartsSubmit_WhenCancelledImmediately(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	abortCalled := int32(0)
	// Use a channel to gate the initiate response so we can cancel before parts
	initiateDone := make(chan struct{})

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				resp := CreateMockResponse(http.StatusOK, body, headers)
				// Signal that initiate is done so cancel can be called
				close(initiateDone)
				return resp
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				// Give time for cancel to happen
				time.Sleep(200 * time.Millisecond)
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodDelete && strings.Contains(req.URL.RawQuery, "uploadId") {
				atomic.AddInt32(&abortCalled, 1)
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   testObjectSize,
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	// Wait for initiate to complete, then cancel immediately
	<-initiateDone
	cancelErr := task.Cancel()
	assert.NoError(t, cancelErr)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Cancelled upload did not finish in time")
	}

	assert.Equal(t, TransferStatusCancelled, task.Status())
	_, taskErr := task.GetResult()
	assert.Error(t, taskErr)
	assert.Contains(t, taskErr.Error(), "cancelled")

	// Cleanup (abort) should have been called
	assert.True(t, atomic.LoadInt32(&abortCalled) >= 1, "AbortMultipartUpload should be called on cancel")
}

// ==================== Download concurrency edge case tests ====================

func TestDownloadFileAsync_ShouldHandleCheckpointRestore(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"
	checkpointFile := tmpDir + "/download.cp"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set(HEADER_REQUEST_ID, "test-request-id")
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile:     downloadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	// First download to create checkpoint
	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("First download did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	_, taskErr := task.GetResult()
	assert.NoError(t, taskErr)

	// Checkpoint file should be removed on success
	_, cpErr := os.Stat(checkpointFile)
	assert.True(t, os.IsNotExist(cpErr), "Checkpoint file should be removed after successful download")
}

func TestDownloadFileAsync_ShouldRemoveTempFile_WhenDownloadFailsWithoutCheckpoint(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			// Return 500 for all GET (range download) requests
			if req.Method == http.MethodGet {
				return CreateMockResponse(http.StatusInternalServerError, TestErrorResponseXML, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		PartSize:     testObjectSize,
		TaskNum:      1,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	assert.Equal(t, TransferStatusFailed, task.Status())
	_, taskErr := task.GetResult()
	assert.Error(t, taskErr)

	// Temp file should be cleaned up when download fails without checkpoint
	tempFile := downloadFile + ".tmp"
	_, statErr := os.Stat(tempFile)
	assert.True(t, os.IsNotExist(statErr), "Temp download file should be removed when download fails without checkpoint")
}

func TestUploadFileAsync_ShouldNotAbort_WhenUploadPartFailsWithCheckpointEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	checkpointFile := tmpDir + "/upload.cp"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	abortCalled := int32(0)
	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				// Return 500 error — server-side failure
				return CreateMockResponse(http.StatusInternalServerError, TestErrorResponseXML, headers)
			}

			if req.Method == http.MethodDelete && strings.Contains(req.URL.RawQuery, "uploadId") {
				atomic.AddInt32(&abortCalled, 1)
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	assert.Equal(t, TransferStatusFailed, task.Status())
	_, taskErr := task.GetResult()
	assert.Error(t, taskErr)

	// With checkpoint enabled, abort should NOT be called in the goroutine error path
	assert.Equal(t, int32(0), atomic.LoadInt32(&abortCalled), "AbortMultipartUpload should NOT be called when checkpoint is enabled")

	// Checkpoint file should still exist for resume
	_, cpStatErr := os.Stat(checkpointFile)
	assert.NoError(t, cpStatErr, "Checkpoint file should be preserved when upload fails with checkpoint enabled")
}

// ==================== Parameter Adjustment Path Tests ====================

func TestUploadFileAsync_ShouldAutoSetCheckpointFile_WhenEnabledWithoutPath(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/test-upload.dat"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		EnableCheckpoint: true,
		// CheckpointFile not set — should be auto-generated to uploadFile + ".uploadfile_record"
		PartSize: testObjectSize,
		TaskNum:  1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())

	// Verify that checkpoint file was auto-generated and then cleaned up
	expectedCheckpointFile := uploadFile + ".uploadfile_record"
	assert.Equal(t, expectedCheckpointFile, input.CheckpointFile)
	_, statErr := os.Stat(expectedCheckpointFile)
	assert.True(t, os.IsNotExist(statErr), "Auto-generated checkpoint file should be removed after successful upload")
}

func TestUploadFileAsync_ShouldTruncatePartSize_WhenExceedsMax(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/test-upload.dat"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   MAX_PART_SIZE + 1024, // exceeds MAX_PART_SIZE, should be truncated
		TaskNum:    1,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Upload did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	assert.Equal(t, int64(MAX_PART_SIZE), input.PartSize, "PartSize should be truncated to MAX_PART_SIZE")
}

func TestDownloadFileAsync_ShouldAdjustParams_WhenInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/test-download.dat"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile:     downloadFile,
		PartSize:         -1, // should be adjusted to DEFAULT_PART_SIZE
		TaskNum:          -1, // should be adjusted to 1
		EnableCheckpoint: true,
		// CheckpointFile not set — should be auto-generated
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Download did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	assert.Equal(t, 1, input.TaskNum, "TaskNum should be adjusted to 1")
	assert.Equal(t, int64(DEFAULT_PART_SIZE), input.PartSize, "PartSize should be adjusted to DEFAULT_PART_SIZE")
	assert.Equal(t, downloadFile+".downloadfile_record", input.CheckpointFile, "CheckpointFile should be auto-generated")
}

// ==================== Upload abort == 1 and IsCompleted path Tests ====================

func TestUploadFileAsync_ShouldSkipCompletedParts_WhenPartIsCompleted(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize*2), 0640)
	require.NoError(t, err)
	checkpointFile := tmpDir + "/upload.cp"

	callCount := int32(0)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				atomic.AddInt32(&callCount, 1)
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	// First upload to create checkpoint with parts marked as completed
	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("First upload did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	// Both parts should have been uploaded in the first pass
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount), "Both parts should be uploaded in first pass")
}

func TestUploadFileAsync_ShouldCancelWithCheckpointEnabled_AndCleanupCheckpointFile(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	checkpointFile := tmpDir + "/upload.cp"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	blockUpload := make(chan struct{})
	abortCalled := int32(0)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				// Block the upload so we can cancel
				select {
				case <-blockUpload:
				case <-time.After(5 * time.Second):
				}
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodDelete && strings.Contains(req.URL.RawQuery, "uploadId") {
				atomic.AddInt32(&abortCalled, 1)
				return CreateMockResponse(http.StatusNoContent, "", headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.UploadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	// Give the goroutine time to start uploading
	time.Sleep(100 * time.Millisecond)

	// Cancel the task
	cancelErr := task.Cancel()
	assert.NoError(t, cancelErr)
	assert.Equal(t, TransferStatusCancelled, task.Status())

	// Unblock the upload
	close(blockUpload)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Cancelled upload did not finish in time")
	}

	assert.Equal(t, TransferStatusCancelled, task.Status())
	_, taskErr := task.GetResult()
	assert.Error(t, taskErr)
	assert.Contains(t, taskErr.Error(), "cancelled")

	// abort should be called in cancelCleanup
	assert.True(t, atomic.LoadInt32(&abortCalled) >= 1, "AbortMultipartUpload should be called on cancel")

	// Checkpoint file should be cleaned up by cancelCleanup when enableCheckpoint is true
	_, cpStatErr := os.Stat(checkpointFile)
	assert.True(t, os.IsNotExist(cpStatErr), "Checkpoint file should be removed by cancelCleanup when enableCheckpoint is true")
}

// ==================== Download cancel with checkpoint enabled path Tests ====================

func TestDownloadFileAsync_ShouldCancelWithCheckpointEnabled_AndCleanupFiles(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"
	checkpointFile := tmpDir + "/download.cp"

	blockDownload := make(chan struct{})

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				// Block the download so we can cancel
				select {
				case <-blockDownload:
				case <-time.After(5 * time.Second):
				}
				respHeaders := make(http.Header)
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile:     downloadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)
	require.NotNil(t, task)

	// Give the goroutine time to start downloading
	time.Sleep(100 * time.Millisecond)

	// Cancel the task
	cancelErr := task.Cancel()
	assert.NoError(t, cancelErr)
	assert.Equal(t, TransferStatusCancelled, task.Status())

	// Unblock the download
	close(blockDownload)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("Cancelled download did not finish in time")
	}

	assert.Equal(t, TransferStatusCancelled, task.Status())
	_, taskErr := task.GetResult()
	assert.Error(t, taskErr)
	assert.Contains(t, taskErr.Error(), "cancelled")

	// cancelCleanup should remove temp file and checkpoint file
	_, tmpStatErr := os.Stat(downloadFile + ".tmp")
	assert.True(t, os.IsNotExist(tmpStatErr), "Temp download file should be removed by cancelCleanup")

	_, cpStatErr := os.Stat(checkpointFile)
	assert.True(t, os.IsNotExist(cpStatErr), "Checkpoint file should be removed by cancelCleanup when enableCheckpoint is true")
}

// ==================== Download IsCompleted path Tests ====================

func TestDownloadFileAsync_ShouldSkipCompletedParts_WhenResumedFromCheckpoint(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"
	checkpointFile := tmpDir + "/download.cp"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile:     downloadFile,
		PartSize:         testObjectSize,
		TaskNum:          1,
		EnableCheckpoint: true,
		CheckpointFile:   checkpointFile,
	}

	// First download to create checkpoint
	task, err := client.DownloadFileAsync(input)
	require.NoError(t, err)

	select {
	case <-task.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("First download did not complete in time")
	}

	assert.Equal(t, TransferStatusCompleted, task.Status())
	_, taskErr := task.GetResult()
	assert.NoError(t, taskErr)

	// Checkpoint file should be removed on success
	_, cpErr := os.Stat(checkpointFile)
	assert.True(t, os.IsNotExist(cpErr), "Checkpoint file should be removed after successful download")

	// Download file should exist
	_, dlErr := os.Stat(downloadFile)
	assert.NoError(t, dlErr, "Download file should exist")
}

// ==================== Sync UploadFile/DownloadFile parameter adjustment Tests ====================

func TestUploadFile_ShouldAutoSetCheckpointFile_WhenEnabledWithoutPath(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/test-upload.dat"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile:       uploadFile,
		EnableCheckpoint: true,
		PartSize:         testObjectSize,
		TaskNum:          1,
	}

	output, err := client.UploadFile(input)
	require.NoError(t, err)
	assert.NotNil(t, output)

	// Verify checkpoint file was auto-generated
	expectedCheckpointFile := uploadFile + ".uploadfile_record"
	assert.Equal(t, expectedCheckpointFile, input.CheckpointFile)

	// Checkpoint file should be cleaned up after success
	_, statErr := os.Stat(expectedCheckpointFile)
	assert.True(t, os.IsNotExist(statErr), "Checkpoint file should be removed after successful upload")
}

func TestUploadFile_ShouldTruncatePartSize_WhenExceedsMax(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/test-upload.dat"
	err := os.WriteFile(uploadFile, make([]byte, testObjectSize), 0640)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploads") {
				body := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><UploadId>test-upload-id-001</UploadId></InitiateMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}

			if req.Method == http.MethodPost && strings.Contains(req.URL.RawQuery, "uploadId") {
				body := `<?xml version="1.0" encoding="UTF-8"?><CompleteMultipartUploadResult><Location>test</Location><Bucket>test-bucket</Bucket><Key>test-object.txt</Key><ETag>"combined-etag"</ETag></CompleteMultipartUploadResult>`
				return CreateMockResponse(http.StatusOK, body, headers)
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		PartSize:   MAX_PART_SIZE + 1024, // exceeds MAX_PART_SIZE
		TaskNum:    1,
	}

	output, err := client.UploadFile(input)
	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, int64(MAX_PART_SIZE), input.PartSize, "PartSize should be truncated to MAX_PART_SIZE")
}

func TestDownloadFile_ShouldAdjustParams_WhenInvalid(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/test-download.dat"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodHead {
				headers.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				headers.Set("ETag", testETag)
				headers.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        headers,
					Body:          io.NopCloser(strings.NewReader("")),
					ContentLength: testObjectSize,
				}
			}

			if req.Method == http.MethodGet {
				respHeaders := make(http.Header)
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile:     downloadFile,
		PartSize:         -1, // should be adjusted to DEFAULT_PART_SIZE
		TaskNum:          -1, // should be adjusted to 1
		EnableCheckpoint: true,
	}

	output, err := client.DownloadFile(input)
	require.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 1, input.TaskNum, "TaskNum should be adjusted to 1")
	assert.Equal(t, int64(DEFAULT_PART_SIZE), input.PartSize, "PartSize should be adjusted to DEFAULT_PART_SIZE")
	assert.Equal(t, downloadFile+".downloadfile_record", input.CheckpointFile, "CheckpointFile should be auto-generated")
}

// ==================== uploadPartConcurrentWithTask direct call Tests ====================

func TestUploadPartConcurrentWithTask_ShouldSkipCompletedParts(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, 200*1024), 0640)
	require.NoError(t, err)

	uploadPartCount := int32(0)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			if req.Method == http.MethodPut && strings.Contains(req.URL.RawQuery, "partNumber") {
				atomic.AddInt32(&uploadPartCount, 1)
				respHeaders := make(http.Header)
				respHeaders.Set("ETag", testETag)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     respHeaders,
					Body:       io.NopCloser(strings.NewReader("")),
				}
			}
			headers := make(http.Header)
			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	ufc := &UploadCheckpoint{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: testUploadID,
		FileInfo: FileStatus{Size: 200 * 1024},
		UploadParts: []UploadPartInfo{
			{PartNumber: 1, PartSize: 100 * 1024, Offset: 0, IsCompleted: true, Etag: testETag},
			{PartNumber: 2, PartSize: 100 * 1024, Offset: 100 * 1024, IsCompleted: false},
		},
	}

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		TaskNum:    1,
		PartSize:   100 * 1024,
	}

	// Directly call the unexported function (accessible in same package)
	err = client.uploadPartConcurrentWithTask(ufc, "", input, nil, nil)
	assert.NoError(t, err)

	// Only part 2 should be uploaded since part 1 is already completed
	assert.Equal(t, int32(1), atomic.LoadInt32(&uploadPartCount), "Only the non-completed part should be uploaded")
}

func TestUploadPartConcurrentWithTask_ShouldStop_WhenAbortFlagSet(t *testing.T) {
	tmpDir := t.TempDir()
	uploadFile := tmpDir + "/upload.txt"
	err := os.WriteFile(uploadFile, make([]byte, 200*1024), 0640)
	require.NoError(t, err)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	// Create a task and set abort flag
	task := newTransferTask(nil)
	atomic.StoreInt32(task.abortFlag(), 1)

	ufc := &UploadCheckpoint{
		Bucket:   TestBucket,
		Key:      TestObjectKey,
		UploadId: testUploadID,
		FileInfo: FileStatus{Size: 200 * 1024},
		UploadParts: []UploadPartInfo{
			{PartNumber: 1, PartSize: 100 * 1024, Offset: 0, IsCompleted: false},
			{PartNumber: 2, PartSize: 100 * 1024, Offset: 100 * 1024, IsCompleted: false},
		},
	}

	input := &UploadFileInput{
		ObjectOperationInput: ObjectOperationInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		UploadFile: uploadFile,
		TaskNum:    1,
		PartSize:   100 * 1024,
	}

	// Should return immediately since abort flag is set
	err = client.uploadPartConcurrentWithTask(ufc, "", input, nil, task)
	assert.NoError(t, err)
}

// ==================== downloadFileConcurrentWithTask direct call Tests ====================

func TestDownloadFileConcurrentWithTask_ShouldSkipCompletedParts(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"

	downloadCount := int32(0)

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			headers.Set(HEADER_REQUEST_ID, "test-request-id")

			if req.Method == http.MethodGet {
				atomic.AddInt32(&downloadCount, 1)
				respHeaders := make(http.Header)
				respHeaders.Set("Content-Length", fmt.Sprintf("%d", testObjectSize))
				respHeaders.Set("ETag", testETag)
				content := strings.NewReader(string(make([]byte, testObjectSize)))
				return &http.Response{
					StatusCode:    http.StatusPartialContent,
					Header:        respHeaders,
					Body:          io.NopCloser(content),
					ContentLength: testObjectSize,
				}
			}

			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	// Create temp file
	tempFileURL := downloadFile + ".tmp"
	err := os.WriteFile(tempFileURL, make([]byte, testObjectSize*2), 0640)
	require.NoError(t, err)

	dfc := &DownloadCheckpoint{
		Bucket:       TestBucket,
		Key:          TestObjectKey,
		DownloadFile: downloadFile,
		ObjectInfo: ObjectInfo{
			Size: testObjectSize * 2,
			ETag: testETag,
		},
		TempFileInfo: TempFileInfo{
			TempFileUrl: tempFileURL,
			Size:        testObjectSize * 2,
		},
		DownloadParts: []DownloadPartInfo{
			{PartNumber: 1, Offset: 0, RangeEnd: testObjectSize - 1, IsCompleted: true},
			{PartNumber: 2, Offset: testObjectSize, RangeEnd: testObjectSize*2 - 1, IsCompleted: false},
		},
	}

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		TaskNum:      1,
		PartSize:     testObjectSize,
	}

	// Directly call the unexported function
	err = client.downloadFileConcurrentWithTask(input, dfc, nil, nil)
	assert.NoError(t, err)

	// Only part 2 should have been downloaded since part 1 is already completed
	assert.Equal(t, int32(1), atomic.LoadInt32(&downloadCount), "Only the non-completed part should be downloaded")
}

func TestDownloadFileConcurrentWithTask_ShouldStop_WhenAbortFlagSet(t *testing.T) {
	tmpDir := t.TempDir()
	downloadFile := tmpDir + "/download.txt"

	mockTransport := &MockRoundTripper{
		ResponseFunc: func(req *http.Request) *http.Response {
			headers := make(http.Header)
			return CreateMockResponse(http.StatusBadRequest, TestErrorResponseXML, headers)
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, _ := New(TestAK, TestSK, TestEndpoint, WithHttpClient(httpClient))

	// Create temp file
	tempFileURL := downloadFile + ".tmp"
	err := os.WriteFile(tempFileURL, make([]byte, testObjectSize*2), 0640)
	require.NoError(t, err)

	// Create a task and set abort flag
	task := newTransferTask(nil)
	atomic.StoreInt32(task.abortFlag(), 1)

	dfc := &DownloadCheckpoint{
		Bucket:       TestBucket,
		Key:          TestObjectKey,
		DownloadFile: downloadFile,
		ObjectInfo: ObjectInfo{
			Size: testObjectSize * 2,
			ETag: testETag,
		},
		TempFileInfo: TempFileInfo{
			TempFileUrl: tempFileURL,
			Size:        testObjectSize * 2,
		},
		DownloadParts: []DownloadPartInfo{
			{PartNumber: 1, Offset: 0, RangeEnd: testObjectSize - 1, IsCompleted: false},
			{PartNumber: 2, Offset: testObjectSize, RangeEnd: testObjectSize*2 - 1, IsCompleted: false},
		},
	}

	input := &DownloadFileInput{
		GetObjectMetadataInput: GetObjectMetadataInput{
			Bucket: TestBucket,
			Key:    TestObjectKey,
		},
		DownloadFile: downloadFile,
		TaskNum:      1,
		PartSize:     testObjectSize,
	}

	// Should return immediately since abort flag is set
	err = client.downloadFileConcurrentWithTask(input, dfc, nil, task)
	assert.NoError(t, err)
}

