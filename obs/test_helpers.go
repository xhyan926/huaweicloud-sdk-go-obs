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
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

// CreateTestObsClient creates a test ObsClient with given endpoint
func CreateTestObsClient(endpoint string, configurers ...interface{}) *ObsClient {
	if endpoint == "" {
		endpoint = "https://obs.test.example.com"
	}
	configs := make([]configurer, 0)
	for _, c := range configurers {
		if conf, ok := c.(configurer); ok {
			configs = append(configs, conf)
		}
	}
	client, err := New("test-ak", "test-sk", endpoint, configs...)
	if err != nil {
		panic(err)
	}
	return client
}

// CreateTestHTTPResponse creates a mock HTTP response for testing
func CreateTestHTTPResponse(statusCode int, body string, headers http.Header) *http.Response {
	return &http.Response{
		StatusCode:    statusCode,
		Status:        http.StatusText(statusCode),
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		Header:        headers,
		ContentLength: int64(len(body)),
	}
}

// CreateTestErrorResponse creates a test error response
func CreateTestErrorResponse(code, message string) []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
	<Code>` + code + `</Code>
	<Message>` + message + `</Message>
	<RequestId>test-request-id</RequestId>
</Error>`)
}

// CreateTestXMLResponse creates a test XML response from an object
func CreateTestXMLResponse(input interface{}) []byte {
	data, err := TransToXml(input)
	if err != nil {
		panic(err)
	}
	return data
}

// SetupTestLogger initializes the test logger to a temp file
func SetupTestLogger(t *testing.T) *os.File {
	tmpFile, err := ioutil.TempFile("", "obs-test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	err = InitLog(tmpFile.Name(), 100, 3, LEVEL_DEBUG, false)
	if err != nil {
		t.Fatalf("Failed to init log: %v", err)
	}
	return tmpFile
}

// CleanupTestLogger closes the logger and cleans up the temp file
func CleanupTestLogger(tmpFile *os.File) {
	CloseLog()
	if tmpFile != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
}

// MockFileInfo is a mock implementation of os.FileInfo for testing
type MockFileInfo struct {
	name    string
	size    int64
	mode     os.FileMode
	modTime  time.Time
}

func (m *MockFileInfo) Name() string       { return m.name }
func (m *MockFileInfo) Size() int64        { return m.size }
func (m *MockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *MockFileInfo) ModTime() time.Time { return m.modTime }
func (m *MockFileInfo) IsDir() bool        { return false }
func (m *MockFileInfo) Sys() interface{}    { return nil }

// MockProgressListener tracks progress events for testing
type MockProgressListener struct {
	Events []*ProgressEvent
}

func (m *MockProgressListener) ProgressChanged(event *ProgressEvent) {
	m.Events = append(m.Events, event)
}

// NewMockProgressListener creates a new mock progress listener
func NewMockProgressListener() *MockProgressListener {
	return &MockProgressListener{
		Events: make([]*ProgressEvent, 0),
	}
}

// AssertHTTPMethod checks if the request has the expected method
func AssertHTTPMethod(t *testing.T, req *http.Request, expectedMethod string) {
	if req.Method != expectedMethod {
		t.Errorf("Expected method %s, got %s", expectedMethod, req.Method)
	}
}

// AssertHTTPHeader checks if the request has the expected header
func AssertHTTPHeader(t *testing.T, req *http.Request, key, expectedValue string) {
	if values := req.Header.Get(key); values != expectedValue {
		t.Errorf("Expected header %s=%s, got %s", key, expectedValue, values)
	}
}

// AssertURLPath checks if the request URL has the expected path
func AssertURLPath(t *testing.T, req *http.Request, expectedPath string) {
	if req.URL.Path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, req.URL.Path)
	}
}

// AssertQueryParameter checks if the request URL has the expected query parameter
func AssertQueryParameter(t *testing.T, req *http.Request, key, expectedValue string) {
	values := req.URL.Query()
	if actualValue := values.Get(key); actualValue != expectedValue {
		t.Errorf("Expected query param %s=%s, got %s", key, expectedValue, actualValue)
	}
}

// CreateTestURL creates a test URL for testing
func CreateTestURL(scheme, host string, port int, path string, query map[string]string) *url.URL {
	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	if port > 0 {
		u.Host = host + ":" + IntToString(port)
	}
	if len(query) > 0 {
		q := make(url.Values)
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	return u
}

// CreateTempFileWithContent creates a temp file with the given content
func CreateTempFileWithContent(t *testing.T, content string) *os.File {
	tmpFile, err := ioutil.TempFile("", "obs-test-*.tmp")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to write temp file: %v", err)
	}
	return tmpFile
}

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, s, substr string) {
	if !bytes.Contains([]byte(s), []byte(substr)) {
		t.Errorf("Expected string to contain %q, got %q", substr, s)
	}
}

// AssertNotContains checks if a string does not contain a substring
func AssertNotContains(t *testing.T, s, substr string) {
	if bytes.Contains([]byte(s), []byte(substr)) {
		t.Errorf("Expected string to not contain %q, got %q", substr, s)
	}
}
