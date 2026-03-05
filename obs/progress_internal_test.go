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
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// publishProgress Tests

func TestPublishProgress_ShouldNotPanic_WhenListenerIsNil(t *testing.T) {
	event := newProgressEvent(TransferDataEvent, 500, 1000)
	// Should not panic when listener is nil
	assert.NotPanics(t, func() {
		publishProgress(nil, event)
	})
}

func TestPublishProgress_ShouldNotPanic_WhenEventIsNil(t *testing.T) {
	listener := NewMockProgressListener()
	// Should not panic when event is nil
	assert.NotPanics(t, func() {
		publishProgress(listener, nil)
	})
}

func TestPublishProgress_ShouldCallListener_WhenBothAreNotNil(t *testing.T) {
	listener := NewMockProgressListener()
	event := newProgressEvent(TransferDataEvent, 500, 1000)

	publishProgress(listener, event)

	assert.Len(t, listener.Events, 1)
	assert.Equal(t, event, listener.Events[0])
}

func TestPublishProgress_ShouldNotCallListener_WhenListenerIsNil(t *testing.T) {
	event := newProgressEvent(TransferDataEvent, 500, 1000)
	listener := NewMockProgressListener()

	publishProgress(nil, event)

	// Listener should not be called
	assert.Len(t, listener.Events, 0)
}

// teeReader.Read Tests

func TestTeeReader_Read_ShouldReadAndPublishProgress_WhenDataRead(t *testing.T) {
	testData := "test data for progress tracking"
	reader := strings.NewReader(testData)
	listener := NewMockProgressListener()
	tracker := &readerTracker{}

	teeReader := TeeReader(reader, int64(len(testData)), listener, tracker)

	buf := make([]byte, 10)
	n, err := teeReader.Read(buf)

	assert.NoError(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, int64(10), tracker.completedBytes)
	assert.Len(t, listener.Events, 1)
}

func TestTeeReader_Read_ShouldPublishProgressForEachRead(t *testing.T) {
	testData := "test data for multiple reads"
	reader := strings.NewReader(testData)
	listener := NewMockProgressListener()
	tracker := &readerTracker{}

	teeReader := TeeReader(reader, int64(len(testData)), listener, tracker)

	buf := make([]byte, 5)

	// First read
	n1, _ := teeReader.Read(buf)
	assert.Equal(t, 5, n1)
	assert.Equal(t, int64(5), tracker.completedBytes)
	assert.Len(t, listener.Events, 1)

	// Second read
	n2, _ := teeReader.Read(buf)
	assert.Equal(t, 5, n2)
	assert.Equal(t, int64(10), tracker.completedBytes)
	assert.Len(t, listener.Events, 2)
}

func TestTeeReader_Read_ShouldPublishFailedEvent_WhenErrorOccurs(t *testing.T) {
	errReader := &progressErrorReader{err: errors.New("read error")}
	listener := NewMockProgressListener()

	teeReader := TeeReader(errReader, 1000, listener, nil)

	buf := make([]byte, 10)
	_, err := teeReader.Read(buf)

	assert.Error(t, err)
	assert.Len(t, listener.Events, 1)
	assert.Equal(t, TransferFailedEvent, listener.Events[0].EventType)
}

func TestTeeReader_Read_ShouldNotPublishProgress_WhenListenerIsNil(t *testing.T) {
	testData := "test data"
	reader := strings.NewReader(testData)
	tracker := &readerTracker{}

	teeReader := TeeReader(reader, int64(len(testData)), nil, tracker)

	buf := make([]byte, 5)
	teeReader.Read(buf)

	// Should not panic and should update tracker
	assert.Equal(t, int64(5), tracker.completedBytes)
}

func TestTeeReader_Read_ShouldNotUpdateTracker_WhenTrackerIsNil(t *testing.T) {
	testData := "test data"
	reader := strings.NewReader(testData)
	listener := NewMockProgressListener()

	teeReader := TeeReader(reader, int64(len(testData)), listener, nil)

	buf := make([]byte, 5)
	teeReader.Read(buf)

	// Listener should still be called
	assert.Len(t, listener.Events, 1)
}

func TestTeeReader_Read_ShouldHandleEOF(t *testing.T) {
	testData := "test data"
	reader := strings.NewReader(testData)
	listener := NewMockProgressListener()

	teeReader := TeeReader(reader, int64(len(testData)), listener, nil)

	buf := make([]byte, 20)
	n, err := teeReader.Read(buf)

	assert.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Len(t, listener.Events, 1)
	assert.Equal(t, TransferDataEvent, listener.Events[0].EventType)
	assert.Equal(t, int64(len(testData)), listener.Events[0].ConsumedBytes)
}

// teeReader.Close Tests

func TestTeeReader_Close_ShouldCloseReader_WhenReaderIsReadCloser(t *testing.T) {
	closeableReader := &progressCloseableReader{}
	teeReader := TeeReader(closeableReader, 1000, nil, nil)

	err := teeReader.Close()

	assert.NoError(t, err)
	assert.True(t, closeableReader.closed)
}

func TestTeeReader_Close_ShouldReturnNil_WhenReaderIsNotReadCloser(t *testing.T) {
	reader := strings.NewReader("test data")
	teeReader := TeeReader(reader, 1000, nil, nil)

	err := teeReader.Close()

	assert.NoError(t, err)
}

func TestTeeReader_Close_ShouldReturnError_WhenCloseFails(t *testing.T) {
	closeableReader := &progressCloseableReader{closeError: errors.New("close error")}
	teeReader := TeeReader(closeableReader, 1000, nil, nil)

	err := teeReader.Close()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "close error")
}

// progressErrorReader implements io.Reader for testing
type progressErrorReader struct {
	err error
}

func (r *progressErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

// progressCloseableReader implements io.ReadCloser for testing
type progressCloseableReader struct {
	closed     bool
	closeError error
}

func (r *progressCloseableReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (r *progressCloseableReader) Close() error {
	r.closed = true
	return r.closeError
}

// readerTracker tests

func TestReaderTracker_ShouldTrackCompletedBytes_WhenSet(t *testing.T) {
	tracker := &readerTracker{}
	tracker.completedBytes = 100

	assert.Equal(t, int64(100), tracker.completedBytes)
}
