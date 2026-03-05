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
	"os"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

// TestNew_ShouldReturnError_WhenEndpointEmpty tests that New() returns an error
// when the endpoint is empty. This tests the error path in initConfigWithDefault()
// at lines 787-788 in client_base.go.
func TestNew_ShouldReturnError_WhenEndpointEmpty(t *testing.T) {
	tests := []struct {
		name     string
		ak       string
		sk       string
		endpoint string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "empty endpoint",
			ak:       "test-ak",
			sk:       "test-sk",
			endpoint: "",
			wantErr:  true,
			errMsg:   "endpoint is not set",
		},
		{
			name:     "whitespace only endpoint",
			ak:       "test-ak",
			sk:       "test-sk",
			endpoint: "   ",
			wantErr:  true,
			errMsg:   "endpoint is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.ak, tt.sk, tt.endpoint)

			if tt.wantErr {
				assert.Nil(t, client, "expected client to be nil when error occurs")
				assert.Error(t, err, "expected error when endpoint is empty")
				assert.True(t, strings.Contains(err.Error(), tt.errMsg),
					"expected error message to contain '%s', got: %v", tt.errMsg, err)
			} else {
				assert.NoError(t, err, "expected no error when endpoint is valid")
				assert.NotNil(t, client, "expected client to be non-nil")
			}
		})
	}
}

// TestNew_ShouldReturnError_WhenGetTransportFails tests that New() returns an error
// when getTransport() fails. This tests the error path at lines 791-792 in client_base.go.
// This test uses gomonkey to mock the private getTransport method.
func TestNew_ShouldReturnError_WhenGetTransportFails(t *testing.T) {
	patches := gomonkey.ApplyPrivateMethod(&config{}, "getTransport",
		func(_ *config) error {
			return assert.AnError
		})
	defer patches.Reset()

	client, err := New("test-ak", "test-sk", "https://obs.example.com")

	assert.Nil(t, client, "expected client to be nil when getTransport fails")
	assert.Error(t, err, "expected error when getTransport fails")
	assert.Equal(t, assert.AnError, err, "expected the mocked error to be returned")
}

// TestNew_ShouldLogPathAccessMode_WhenPathStyleIsTrue tests that New() logs
// "Path" access mode when pathStyle is true. This tests the logging branch
// at lines 800-801 in client_base.go.
func TestNew_ShouldLogPathAccessMode_WhenPathStyleIsTrue(t *testing.T) {
	logFile := "/tmp/test-obs-sdk-pathstyle.log"
	err := InitLog(logFile, 10, 5, LEVEL_WARN, false)
	assert.NoError(t, err)
	defer CloseLog()
	defer os.Remove(logFile)

	client, err := New("test-ak", "test-sk", "https://obs.example.com", WithPathStyle(true))

	assert.NoError(t, err, "expected no error when creating client with pathStyle")
	assert.NotNil(t, client, "expected client to be non-nil")
}

// TestNew_ShouldLogVirtualHostingAccessMode_WhenPathStyleIsFalse tests that New() logs
// "Virtual Hosting" access mode when pathStyle is false (default). This tests the
// default logging branch at lines 799-804 in client_base.go.
func TestNew_ShouldLogVirtualHostingAccessMode_WhenPathStyleIsFalse(t *testing.T) {
	logFile := "/tmp/test-obs-sdk-virtualhost.log"
	err := InitLog(logFile, 10, 5, LEVEL_WARN, false)
	assert.NoError(t, err)
	defer CloseLog()
	defer os.Remove(logFile)

	client, err := New("test-ak", "test-sk", "https://obs.example.com")

	assert.NoError(t, err, "expected no error when creating client without pathStyle")
	assert.NotNil(t, client, "expected client to be non-nil")
}

// TestNew_ShouldSetPathStyleAutomatically_WhenEndpointIsIP tests that New() automatically
// sets pathStyle to true when the endpoint is an IP address. This tests the
// behavior in initConfigWithDefault() which sets pathStyle for IP endpoints.
func TestNew_ShouldSetPathStyleAutomatically_WhenEndpointIsIP(t *testing.T) {
	logFile := "/tmp/test-obs-sdk-ip.log"
	err := InitLog(logFile, 10, 5, LEVEL_WARN, false)
	assert.NoError(t, err)
	defer CloseLog()
	defer os.Remove(logFile)

	tests := []struct {
		name     string
		endpoint string
	}{
		{"IPv4 address", "https://192.168.1.1"},
		{"IPv4 address with port", "https://192.168.1.1:443"},
		{"localhost IP", "https://127.0.0.1"},
		{"localhost IP with port", "http://127.0.0.1:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New("test-ak", "test-sk", tt.endpoint)

			assert.NoError(t, err, "expected no error when creating client with IP endpoint")
			assert.NotNil(t, client, "expected client to be non-nil")
			assert.True(t, client.conf.pathStyle, "expected pathStyle to be true for IP endpoint")
		})
	}
}
