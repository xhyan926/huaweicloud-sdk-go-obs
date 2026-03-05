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
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// WithLoggerTimeLoc Tests

func TestWithLoggerTimeLoc_ShouldReturnConfigFunc_WhenGivenLocation(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*3600)
	configFunc := WithLoggerTimeLoc(loc)
	assert.NotNil(t, configFunc)
}

// checkAndLogErr Tests

func TestCheckAndLogErr_ShouldLog_WhenErrorNotNil(t *testing.T) {
	testErr := assert.AnError
	checkAndLogErr(testErr, LEVEL_ERROR, "Test error: %v", testErr)
	// Test passes if no panic occurs
}

func TestCheckAndLogErr_ShouldNotLog_WhenErrorIsNil(t *testing.T) {
	checkAndLogErr(nil, LEVEL_ERROR, "Test error")
	// Test passes if no panic occurs
}

func TestCheckAndLogErr_ShouldHandleDifferentLevels(t *testing.T) {
	levels := []Level{LEVEL_DEBUG, LEVEL_INFO, LEVEL_WARN, LEVEL_ERROR}
	for _, level := range levels {
		checkAndLogErr(assert.AnError, level, "Test at level %v", level)
	}
}

// logResponseHeader Tests

func TestLogResponseHeader_ShouldReturnFormattedHeaders_WhenHeadersAreValid(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Etag", "abc123")
	headers.Set("x-obs-request-id", "req-123")

	result := logResponseHeader(headers)
	assert.Contains(t, result, "Content-Type")
	assert.Contains(t, result, "Etag")
	assert.Contains(t, result, "X-Obs-Request-Id")
}

func TestLogResponseHeader_ShouldHandleEmptyHeaders(t *testing.T) {
	headers := http.Header{}
	result := logResponseHeader(headers)
	assert.Equal(t, "", result)
}

func TestLogResponseHeader_ShouldHandleMultipleValues(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("Content-Type", "text/plain")

	result := logResponseHeader(headers)
	assert.Contains(t, result, "Content-Type")
}

func TestLogResponseHeader_ShouldSkipNonAllowedHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("User-Agent", "test-agent")
	headers.Set("Content-Type", "application/json")

	result := logResponseHeader(headers)
	assert.Contains(t, result, "Content-Type")
	assert.NotContains(t, result, "User-Agent")
}

func TestLogResponseHeader_ShouldHandleHeaderPrefixes(t *testing.T) {
	headers := http.Header{}
	headers.Set("x-reserved-indicator", "true")
	headers.Set("x-obs-request-id", "req-123")

	result := logResponseHeader(headers)
	assert.Contains(t, result, "X-Reserved-Indicator")
	assert.Contains(t, result, "X-Obs-Request-Id")
}

// logRequestHeader Tests

func TestLogRequestHeader_ShouldReturnFormattedHeaders_WhenHeadersAreValid(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Content-Md5", "abc123")

	result := logRequestHeader(headers)
	assert.Contains(t, result, "Content-Type")
	assert.Contains(t, result, "Content-Md5")
}

func TestLogRequestHeader_ShouldHandleEmptyHeaders(t *testing.T) {
	headers := http.Header{}
	result := logRequestHeader(headers)
	assert.Equal(t, "", result)
}

func TestLogRequestHeader_ShouldSkipNonAllowedHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("User-Agent", "test-agent")
	headers.Set("Content-Type", "application/json")

	result := logRequestHeader(headers)
	assert.Contains(t, result, "Content-Type")
	assert.NotContains(t, result, "User-Agent")
}

func TestLogRequestHeader_ShouldHandleWhitespaceKeys(t *testing.T) {
	headers := http.Header{}
	headers.Set(" Content-Type ", "application/json")

	result := logRequestHeader(headers)
	assert.Contains(t, result, "Content-Type")
}

// isErrorLogEnabled Tests

func TestIsErrorLogEnabled_ShouldReturnTrue_WhenLogLevelIsDebug(t *testing.T) {
	_ = InitLogWithCacheCnt("", 0, 0, LEVEL_DEBUG, false, 1)
	defer CloseLog()

	result := isErrorLogEnabled()
	assert.True(t, result)
}

func TestIsErrorLogEnabled_ShouldReturnTrue_WhenLogLevelIsError(t *testing.T) {
	_ = InitLogWithCacheCnt("", 0, 0, LEVEL_ERROR, false, 1)
	defer CloseLog()

	result := isErrorLogEnabled()
	assert.True(t, result)
}

func TestIsErrorLogEnabled_ShouldReturnFalse_WhenLogLevelIsOff(t *testing.T) {
	_ = InitLogWithCacheCnt("", 0, 0, LEVEL_OFF, false, 1)
	defer CloseLog()

	result := isErrorLogEnabled()
	assert.False(t, result)
}

// DoLog Tests

func TestDoLog_ShouldNotPanic_WhenLoggingDisabled(t *testing.T) {
	// Log is disabled by default
	DoLog(LEVEL_DEBUG, "Test log message")
	// Test passes if no panic occurs
}

func TestDoLog_ShouldNotPanic_WhenLogLevelNotMet(t *testing.T) {
	_ = InitLogWithCacheCnt("", 0, 0, LEVEL_ERROR, false, 1)
	defer CloseLog()

	DoLog(LEVEL_INFO, "This should not be logged")
	// Test passes if no panic occurs
}

func TestDoLog_ShouldNotPanic_WhenLogLevelIsMet(t *testing.T) {
	_ = InitLogWithCacheCnt("", 0, 0, LEVEL_DEBUG, false, 1)
	defer CloseLog()

	DoLog(LEVEL_DEBUG, "Test log message")
	// Test passes if no panic occurs
}

func TestDoLog_ShouldHandleMultipleArgs(t *testing.T) {
	_ = InitLogWithCacheCnt("", 0, 0, LEVEL_DEBUG, false, 1)
	defer CloseLog()

	DoLog(LEVEL_DEBUG, "Test message with %s and %d", "arg1", 42)
	// Test passes if no panic occurs
}
