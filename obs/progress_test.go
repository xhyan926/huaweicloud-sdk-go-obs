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

// TestProgressEventType_ShouldHaveCorrectValues_WhenChecked tests ProgressEventType constants
func TestProgressEventType_ShouldHaveCorrectValues_WhenChecked(t *testing.T) {
	tests := []struct {
		name      string
		eventType ProgressEventType
		expected  int
	}{
		{"TransferStartedEvent", TransferStartedEvent, 1},
		{"TransferDataEvent", TransferDataEvent, 2},
		{"TransferCompletedEvent", TransferCompletedEvent, 3},
		{"TransferFailedEvent", TransferFailedEvent, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, int(tt.eventType))
		})
	}
}

// TestNewProgressEvent_ShouldCreateEvent_WhenGivenValidInputs tests newProgressEvent
func TestNewProgressEvent_ShouldCreateEvent_WhenGivenValidInputs(t *testing.T) {
	eventType := TransferDataEvent
	consumed := int64(500)
	total := int64(1000)

	event := newProgressEvent(eventType, consumed, total)

	assert.NotNil(t, event)
	assert.Equal(t, consumed, event.ConsumedBytes)
	assert.Equal(t, total, event.TotalBytes)
	assert.Equal(t, eventType, event.EventType)
}

// TestNewProgressEvent_ShouldHandleZeroValues_WhenCalled tests newProgressEvent with zero values
func TestNewProgressEvent_ShouldHandleZeroValues_WhenCalled(t *testing.T) {
	event := newProgressEvent(TransferDataEvent, 0, 0)

	assert.NotNil(t, event)
	assert.Equal(t, int64(0), event.ConsumedBytes)
	assert.Equal(t, int64(0), event.TotalBytes)
}

// TestTeeReader_Size_ShouldReturnCorrectSize_WhenCalled tests TeeReader Size method
func TestTeeReader_Size_ShouldReturnCorrectSize_WhenCalled(t *testing.T) {
	totalBytes := int64(1024)
	reader := TeeReader(nil, totalBytes, nil, nil)

	teeReader, ok := reader.(*teeReader)
	assert.True(t, ok)
	assert.Equal(t, totalBytes, teeReader.Size())
}

// TestTeeReader_Close_ShouldNotPanic_WhenCalled tests TeeReader Close method
func TestTeeReader_Close_ShouldNotPanic_WhenCalled(t *testing.T) {
	reader := TeeReader(nil, 1000, nil, nil)

	// Should not panic
	assert.NotPanics(t, func() {
		reader.Close()
	})
}
