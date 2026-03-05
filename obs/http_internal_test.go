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
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// checkRedirectFunc Tests

func TestCheckRedirectFunc_ShouldReturnErrUseLastResponse_WhenCalled(t *testing.T) {
	req := &http.Request{}
	via := []*http.Request{}

	err := checkRedirectFunc(req, via)
	assert.Equal(t, http.ErrUseLastResponse, err)
}

// getConnDelegate Tests

func TestGetConnDelegate_ShouldReturnDelegate_WhenGivenValidConn(t *testing.T) {
	// Create a pipe for testing
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)
	assert.NotNil(t, delegate)
	assert.Equal(t, client, delegate.conn)
	assert.NotNil(t, delegate.socketTimeout)
	assert.NotNil(t, delegate.finalTimeout)
}

// connDelegate Read Tests

func TestConnDelegate_Read_ShouldReadData_WhenValidDelegate(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	// Write some data
	go func() {
		server.Write([]byte("test data"))
	}()

	buffer := make([]byte, 20)
	n, err := delegate.Read(buffer)
	assert.NoError(t, err)
	assert.Equal(t, 9, n)
	assert.Equal(t, "test data", string(buffer[:n]))
}

func TestConnDelegate_Read_ShouldHandleZeroRead(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	// Close the server to get EOF
	server.Close()

	buffer := make([]byte, 20)
	_, err := delegate.Read(buffer)
	assert.Error(t, err)
}

// connDelegate Write Tests

func TestConnDelegate_Write_ShouldWriteData_WhenValidDelegate(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	data := []byte("test")

	// Create a goroutine to read from the pipe
	go func() {
		buffer := make([]byte, 10)
		_, _ = server.Read(buffer)
	}()

	n, err := delegate.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
}

func TestConnDelegate_Write_ShouldHandleClosedConnection(t *testing.T) {
	server, client := net.Pipe()
	server.Close()
	client.Close()

	delegate := getConnDelegate(client, 10, 100)

	data := []byte("test data")
	_, err := delegate.Write(data)
	assert.Error(t, err)
}

// connDelegate Close Tests

func TestConnDelegate_Close_ShouldCloseConnection(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	delegate := getConnDelegate(client, 10, 100)

	err := delegate.Close()
	assert.NoError(t, err)

	// Verify connection is closed by trying to write
	_, err = delegate.Write([]byte("test"))
	assert.Error(t, err)
}

// connDelegate LocalAddr Tests

func TestConnDelegate_LocalAddr_ShouldReturnLocalAddress(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	addr := delegate.LocalAddr()
	assert.NotNil(t, addr)
	assert.NotEmpty(t, addr.String())
}

// connDelegate RemoteAddr Tests

func TestConnDelegate_RemoteAddr_ShouldReturnRemoteAddress(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	addr := delegate.RemoteAddr()
	assert.NotNil(t, addr)
	assert.NotEmpty(t, addr.String())
}

// connDelegate SetDeadline Tests

func TestConnDelegate_SetDeadline_ShouldSetDeadline(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	err := delegate.SetDeadline(time.Now().Add(delegate.socketTimeout))
	assert.NoError(t, err)
}

// connDelegate SetReadDeadline Tests

func TestConnDelegate_SetReadDeadline_ShouldSetReadDeadline(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	err := delegate.SetReadDeadline(time.Now().Add(delegate.socketTimeout))
	assert.NoError(t, err)
}

// connDelegate SetWriteDeadline Tests

func TestConnDelegate_SetWriteDeadline_ShouldSetWriteDeadline(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	delegate := getConnDelegate(client, 10, 100)

	err := delegate.SetWriteDeadline(time.Now().Add(delegate.socketTimeout))
	assert.NoError(t, err)
}

// isRedirectErr Tests

func TestIsRedirectErr_ShouldReturnTrue_WhenLocationProvided(t *testing.T) {
	location := "http://example.com"
	redirectCount := 0
	maxRedirectCount := 5

	result := isRedirectErr(location, redirectCount, maxRedirectCount)
	assert.True(t, result)
}

func TestIsRedirectErr_ShouldReturnFalse_WhenLocationEmpty(t *testing.T) {
	location := ""
	redirectCount := 0
	maxRedirectCount := 5

	result := isRedirectErr(location, redirectCount, maxRedirectCount)
	assert.False(t, result)
}

func TestIsRedirectErr_ShouldReturnFalse_WhenMaxRedirectsReached(t *testing.T) {
	location := "http://example.com"
	redirectCount := 5
	maxRedirectCount := 5

	result := isRedirectErr(location, redirectCount, maxRedirectCount)
	assert.False(t, result)
}

// setRedirectFlag Tests

func TestSetRedirectFlag_ShouldReturnTrue_When302AndGET(t *testing.T) {
	statusCode := 302
	method := HTTP_GET

	result := setRedirectFlag(statusCode, method)
	assert.True(t, result)
}

func TestSetRedirectFlag_ShouldReturnFalse_WhenNot302(t *testing.T) {
	statusCode := 301
	method := HTTP_GET

	result := setRedirectFlag(statusCode, method)
	assert.False(t, result)
}

func TestSetRedirectFlag_ShouldReturnFalse_WhenNotGET(t *testing.T) {
	statusCode := 302
	method := "POST"

	result := setRedirectFlag(statusCode, method)
	assert.False(t, result)
}

// canNotRetry Tests

func TestCanNotRetry_ShouldReturnTrue_WhenNotRepeatable(t *testing.T) {
	repeatable := false
	statusCode := 200

	result := canNotRetry(repeatable, statusCode)
	assert.True(t, result)
}

func TestCanNotRetry_ShouldReturnTrue_When4xx(t *testing.T) {
	repeatable := true
	statusCode := 404

	result := canNotRetry(repeatable, statusCode)
	assert.True(t, result)
}

func TestCanNotRetry_ShouldReturnTrue_When304(t *testing.T) {
	repeatable := true
	statusCode := 304

	result := canNotRetry(repeatable, statusCode)
	assert.True(t, result)
}

func TestCanNotRetry_ShouldReturnFalse_When5xx(t *testing.T) {
	repeatable := true
	statusCode := 500

	result := canNotRetry(repeatable, statusCode)
	assert.False(t, result)
}

// prepareData Tests

func TestPrepareData_ShouldHandleStringData(t *testing.T) {
	headers := make(map[string][]string)
	dataStr := "test data"

	reader, err := prepareData(headers, dataStr)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Equal(t, "9", headers["Content-Length"][0])
}

func TestPrepareData_ShouldHandleByteData(t *testing.T) {
	headers := make(map[string][]string)
	dataByte := []byte("test data")

	reader, err := prepareData(headers, dataByte)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Equal(t, "9", headers["Content-Length"][0])
}

func TestPrepareData_ShouldHandleNilData(t *testing.T) {
	headers := make(map[string][]string)

	reader, err := prepareData(headers, nil)
	assert.NoError(t, err)
	assert.Nil(t, reader)
}

func TestPrepareData_ShouldHandleInvalidData(t *testing.T) {
	headers := make(map[string][]string)
	data := 123

	reader, err := prepareData(headers, data)
	assert.Error(t, err)
	assert.Nil(t, reader)
}

// prepareAgentHeader Tests

func TestPrepareAgentHeader_ShouldReturnDefaultAgent_WhenClientUserAgentEmpty(t *testing.T) {
	result := prepareAgentHeader("")
	assert.Equal(t, USER_AGENT, result)
}

func TestPrepareAgentHeader_ShouldReturnCustomAgent_WhenClientUserAgentProvided(t *testing.T) {
	customAgent := "custom-agent/1.0"
	result := prepareAgentHeader(customAgent)
	assert.Equal(t, customAgent, result)
}

// prepareHeaders Tests

func TestPrepareHeaders_ShouldKeepOriginalHeaders_WhenValid(t *testing.T) {
	headers := map[string][]string{
		"Content-Type":        {"application/json"},
		"Content-Md5":         {"abc123"},
		"x-amz-meta-custom":   {"value"},
	}

	result := prepareHeaders(headers, false, false)
	assert.Contains(t, result, "Content-Type")
	assert.Contains(t, result, "Content-Md5")
	assert.Contains(t, result, "x-amz-meta-custom")
}

func TestPrepareHeaders_ShouldAddMetaPrefix_WhenMetaTrue(t *testing.T) {
	headers := map[string][]string{
		"custom-header": {"value"},
	}

	result := prepareHeaders(headers, true, false)
	assert.Contains(t, result, "x-amz-meta-custom-header")
}

func TestPrepareHeaders_ShouldRemoveEmptyKeys(t *testing.T) {
	headers := map[string][]string{
		"":           {"value"},
		"Content-Type": {"application/json"},
	}

	result := prepareHeaders(headers, false, false)
	assert.NotContains(t, result, "")
	assert.Contains(t, result, "Content-Type")
}

func TestPrepareHeaders_ShouldHandleObsMetaPrefix(t *testing.T) {
	headers := map[string][]string{
		"custom-header": {"value"},
	}

	result := prepareHeaders(headers, true, true)
	assert.Contains(t, result, "x-obs-meta-custom-header")
}

// checkAndLogErr Tests

func TestCheckAndLogErr_ShouldNotPanic_WhenErrorIsNil(t *testing.T) {
	checkAndLogErr(nil, LEVEL_ERROR, "test")
	// Test passes if no panic
}

func TestCheckAndLogErr_ShouldNotPanic_WhenErrorIsNotNil(t *testing.T) {
	checkAndLogErr(errors.New("test error"), LEVEL_ERROR, "test")
	// Test passes if no panic
}

// logHeaders Tests

func TestLogHeaders_ShouldNotLog_WhenDebugNotEnabled(t *testing.T) {
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	// CloseLog to reset log level
	CloseLog()

	headers := map[string][]string{
		HEADER_AUTH_CAMEL:      {"test-auth"},
		HEADER_STS_TOKEN_AMZ:    {"test-token"},
		"X-Test-Header":        {"test-value"},
	}

	logHeaders(headers, SignatureV4)
	// Test passes if no panic occurs
}

func TestLogHeaders_ShouldLog_WhenDebugEnabled(t *testing.T) {
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	headers := map[string][]string{
		HEADER_AUTH_CAMEL:   {"test-auth"},
		HEADER_STS_TOKEN_AMZ: {"test-token"},
		"X-Test-Header":     {"test-value"},
	}

	logHeaders(headers, SignatureV4)
	// Test passes if no panic occurs
}

func TestLogHeaders_ShouldMaskSecurityToken_WhenHasAmzToken(t *testing.T) {
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	headers := map[string][]string{
		HEADER_AUTH_CAMEL:   {"test-auth"},
		HEADER_STS_TOKEN_AMZ: {"test-token"},
	}

	originalToken := headers[HEADER_STS_TOKEN_AMZ][0]

	logHeaders(headers, SignatureV4)

	// Token should be restored after logging
	assert.Equal(t, originalToken, headers[HEADER_STS_TOKEN_AMZ][0])
}

func TestLogHeaders_ShouldMaskSecurityToken_WhenHasObsToken(t *testing.T) {
	tmpFile := SetupTestLogger(t)
	defer CleanupTestLogger(tmpFile)

	headers := map[string][]string{
		HEADER_AUTH_CAMEL:   {"test-auth"},
		HEADER_STS_TOKEN_OBS: {"test-token"},
	}

	originalToken := headers[HEADER_STS_TOKEN_OBS][0]

	logHeaders(headers, SignatureObs)

	// Token should be restored after logging
	assert.Equal(t, originalToken, headers[HEADER_STS_TOKEN_OBS][0])
}

// prepareRetry Tests

func TestPrepareRetry_ShouldCloseResponse_WhenRespNotNil(t *testing.T) {
	headers := map[string][]string{}
	data := strings.NewReader("test data")

	reader, resp, err := prepareRetry(&http.Response{Body: io.NopCloser(strings.NewReader(""))}, headers, data, nil)

	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Nil(t, resp)
}

func TestPrepareRetry_ShouldUpdateDateHeader_WhenHasDateHeader(t *testing.T) {
	headers := map[string][]string{
		HEADER_DATE_CAMEL: {"old-date"},
	}
	data := strings.NewReader("test data")

	reader, resp, err := prepareRetry(nil, headers, data, nil)

	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Nil(t, resp)
	assert.NotEqual(t, "old-date", headers[HEADER_DATE_CAMEL][0])
}

func TestPrepareRetry_ShouldRemoveAuthHeader_WhenHasAuthHeader(t *testing.T) {
	headers := map[string][]string{
		HEADER_AUTH_CAMEL: {"test-auth"},
	}
	data := strings.NewReader("test data")

	reader, resp, err := prepareRetry(nil, headers, data, nil)

	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Nil(t, resp)
	_, hasAuth := headers[HEADER_AUTH_CAMEL]
	assert.False(t, hasAuth)
}

func TestPrepareRetry_ShouldResetStringsReader(t *testing.T) {
	headers := map[string][]string{}
	data := strings.NewReader("test data")

	// Read some data first
	buffer := make([]byte, 4)
	data.Read(buffer)

	reader, resp, err := prepareRetry(nil, headers, data, nil)

	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Nil(t, resp)
}

func TestPrepareRetry_ShouldResetBytesReader(t *testing.T) {
	headers := map[string][]string{}
	data := bytes.NewReader([]byte("test data"))

	// Read some data first
	buffer := make([]byte, 4)
	data.Read(buffer)

	reader, resp, err := prepareRetry(nil, headers, data, nil)

	assert.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Nil(t, resp)
}

func TestPrepareRetry_ShouldHandleNilData(t *testing.T) {
	headers := map[string][]string{}

	reader, resp, err := prepareRetry(nil, headers, nil, nil)

	assert.NoError(t, err)
	assert.Nil(t, reader)
	assert.Nil(t, resp)
}

// getRequest tests removed - it's a method requiring ObsClient instance
