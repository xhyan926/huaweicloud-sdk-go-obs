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
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PutObjectOutput.setCallbackReadCloser Tests

func TestPutObjectOutput_SetCallbackReadCloser_ShouldSetData_WhenCalled(t *testing.T) {
	output := &PutObjectOutput{}
	reader := bytes.NewBufferString("callback data")

	output.setCallbackReadCloser(io.NopCloser(reader))
	assert.NotNil(t, output.CallbackBody.data)
}

func TestPutObjectOutput_SetCallbackReadCloser_ShouldReplaceExistingData_WhenCalledAgain(t *testing.T) {
	output := &PutObjectOutput{}
	reader1 := bytes.NewBufferString("first data")
	reader2 := bytes.NewBufferString("second data")

	output.setCallbackReadCloser(io.NopCloser(reader1))
	assert.NotNil(t, output.CallbackBody.data)

	output.setCallbackReadCloser(io.NopCloser(reader2))
	assert.NotNil(t, output.CallbackBody.data)
}

// CompleteMultipartUploadOutput.setCallbackReadCloser Tests

func TestCompleteMultipartUploadOutput_SetCallbackReadCloser_ShouldSetData_WhenCalled(t *testing.T) {
	output := &CompleteMultipartUploadOutput{}
	reader := bytes.NewBufferString("multipart callback data")

	output.setCallbackReadCloser(io.NopCloser(reader))
	assert.NotNil(t, output.CallbackBody.data)
}

func TestCompleteMultipartUploadOutput_SetCallbackReadCloser_ShouldReplaceExistingData_WhenCalledAgain(t *testing.T) {
	output := &CompleteMultipartUploadOutput{}
	reader1 := bytes.NewBufferString("first data")
	reader2 := bytes.NewBufferString("second data")

	output.setCallbackReadCloser(io.NopCloser(reader1))
	assert.NotNil(t, output.CallbackBody.data)

	output.setCallbackReadCloser(io.NopCloser(reader2))
	assert.NotNil(t, output.CallbackBody.data)
}

// CallbackBody.ReadCallbackBody Tests

func TestCallbackBody_ReadCallbackBody_ShouldReturnError_WhenDataIsNil(t *testing.T) {
	body := CallbackBody{data: nil}
	buf := make([]byte, 100)

	n, err := body.ReadCallbackBody(buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "have no callback data")
	assert.Equal(t, 0, n)
}

func TestCallbackBody_ReadCallbackBody_ShouldReadData_WhenDataIsNotNil(t *testing.T) {
	testData := "callback response data"
	body := CallbackBody{data: io.NopCloser(strings.NewReader(testData))}
	buf := make([]byte, 50)

	n, err := body.ReadCallbackBody(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Equal(t, testData, string(buf[:n]))
}

func TestCallbackBody_ReadCallbackBody_ShouldHandlePartialRead_WhenBufferIsSmall(t *testing.T) {
	testData := "long callback data that requires multiple reads"
	body := CallbackBody{data: io.NopCloser(strings.NewReader(testData))}
	buf := make([]byte, 10)

	n, err := body.ReadCallbackBody(buf)
	assert.NoError(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, "long callb", string(buf[:n]))
}

func TestCallbackBody_ReadCallbackBody_ShouldReturnEOF_WhenAllDataRead(t *testing.T) {
	testData := "end of data"
	body := CallbackBody{data: io.NopCloser(strings.NewReader(testData))}
	buf := make([]byte, 100)

	// First read
	n1, err1 := body.ReadCallbackBody(buf)
	assert.NoError(t, err1)
	assert.Equal(t, len(testData), n1)

	// Second read should return EOF
	n2, err2 := body.ReadCallbackBody(buf)
	assert.Error(t, err2)
	assert.Equal(t, io.EOF, err2)
	assert.Equal(t, 0, n2)
}

// CallbackBody.CloseCallbackBody Tests

func TestCallbackBody_CloseCallbackBody_ShouldReturnError_WhenDataIsNil(t *testing.T) {
	body := CallbackBody{data: nil}

	err := body.CloseCallbackBody()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "have no callback data")
}

func TestCallbackBody_CloseCallbackBody_ShouldCloseData_WhenDataIsNotNil(t *testing.T) {
	testData := "data to close"
	body := CallbackBody{data: io.NopCloser(strings.NewReader(testData))}

	err := body.CloseCallbackBody()
	assert.NoError(t, err)
}

func TestCallbackBody_CloseCallbackBody_ShouldBeIdempotent_WhenCalledMultipleTimes(t *testing.T) {
	testData := "data for idempotent close test"
	body := CallbackBody{data: io.NopCloser(strings.NewReader(testData))}

	// First close
	err1 := body.CloseCallbackBody()
	assert.NoError(t, err1)

	// Second close should also succeed
	err2 := body.CloseCallbackBody()
	assert.NoError(t, err2)
}
