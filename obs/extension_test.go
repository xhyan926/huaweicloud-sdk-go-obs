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
)

// TestWithProgress_ShouldCreateHeaderOption_WhenCalled tests WithProgress function
func TestWithProgress_ShouldCreateHeaderOption_WhenCalled(t *testing.T) {
	// Using a simple listener implementation
	listener := NewMockProgressListener()
	option := WithProgress(listener)

	assert.NotNil(t, option)
}

// TestWithProgress_ShouldNilSafe_WhenNilListener tests WithProgress with nil
func TestWithProgress_ShouldNilSafe_WhenNilListener(t *testing.T) {
	option := WithProgress(nil)

	assert.NotNil(t, option)
}

// TestWithReqPaymentHeader_ShouldCreateHeaderOption_WhenCalled tests WithReqPaymentHeader
func TestWithReqPaymentHeader_ShouldCreateHeaderOption_WhenCalled(t *testing.T) {
	payer := PayerType("requester")
	option := WithReqPaymentHeader(payer)

	assert.NotNil(t, option)
}

// TestWithReqPaymentHeader_ShouldHandleEmptyPayer_WhenCalled tests WithReqPaymentHeader empty
func TestWithReqPaymentHeader_ShouldHandleEmptyPayer_WhenCalled(t *testing.T) {
	option := WithReqPaymentHeader("")

	assert.NotNil(t, option)
}

// TestWithTrafficLimitHeader_ShouldCreateHeaderOption_WhenCalled tests WithTrafficLimitHeader
func TestWithTrafficLimitHeader_ShouldCreateHeaderOption_WhenCalled(t *testing.T) {
	trafficLimit := int64(819200)
	option := WithTrafficLimitHeader(trafficLimit)

	assert.NotNil(t, option)
}

// TestWithTrafficLimitHeader_ShouldHandleZeroLimit_WhenCalled tests WithTrafficLimitHeader with zero
func TestWithTrafficLimitHeader_ShouldHandleZeroLimit_WhenCalled(t *testing.T) {
	option := WithTrafficLimitHeader(0)

	assert.NotNil(t, option)
}

// TestWithHttpTransport_ShouldCreateTransportOption_WhenCalled tests WithHttpTransport
func TestWithHttpTransport_ShouldCreateTransportOption_WhenCalled(t *testing.T) {
	transport := &http.Transport{
		MaxIdleConns: 100,
	}
	conf := &config{}
	configurer := WithHttpTransport(transport)
	configurer(conf)
	assert.Equal(t, transport, conf.transport)
}

// TestWithHttpClient_ShouldCreateClientOption_WhenCalled tests WithHttpClient
func TestWithHttpClient_ShouldCreateClientOption_WhenCalled(t *testing.T) {
	httpClient := &http.Client{}
	conf := &config{}
	configurer := WithHttpClient(httpClient)
	configurer(conf)
	assert.Equal(t, httpClient, conf.httpClient)
}

// TestProgressEvent_ShouldCreateEvent_WhenCalled tests newProgressEvent
func TestProgressEvent_ShouldCreateEvent_WhenCalled(t *testing.T) {
	total := int64(1000)
	consumed := int64(500)

	event := newProgressEvent(TransferDataEvent, consumed, total)

	assert.Equal(t, total, event.TotalBytes)
	assert.Equal(t, consumed, event.ConsumedBytes)
}

// TestProgressEvent_ShouldHandleZeroTotal_WhenCalled tests newProgressEvent with zero total
func TestProgressEvent_ShouldHandleZeroTotal_WhenCalled(t *testing.T) {
	event := newProgressEvent(TransferDataEvent, 0, 0)

	assert.Equal(t, int64(0), event.TotalBytes)
	assert.Equal(t, int64(0), event.ConsumedBytes)
}

// TestPublishProgress_ShouldHandleNilListener_WhenListenerIsNil tests publishProgress with nil
func TestPublishProgress_ShouldHandleNilListener_WhenListenerIsNil(t *testing.T) {
	event := newProgressEvent(TransferDataEvent, 500, 1000)

	// Should not panic
	assert.NotPanics(t, func() {
		publishProgress(nil, event)
	})
}

// TestProgressEventType_ShouldHaveValidValues_WhenCalled tests ProgressEventType constants
func TestProgressEventType_ShouldHaveValidValues_WhenCalled(t *testing.T) {
	tests := []struct {
		name   string
		eventType ProgressEventType
	}{
		{"Transfer started", TransferStartedEvent},
		{"Data transferred", TransferDataEvent},
		{"Transfer completed", TransferCompletedEvent},
		{"Transfer failed", TransferFailedEvent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify event type is not zero
			assert.True(t, tt.eventType > 0)
		})
	}
}
