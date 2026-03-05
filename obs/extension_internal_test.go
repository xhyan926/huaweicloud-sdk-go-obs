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
	"testing"

	"github.com/stretchr/testify/assert"
)

// setHeaderPrefix Tests

func TestSetHeaderPrefix_ShouldSetHeader_WhenValueIsValid(t *testing.T) {
	headers := make(map[string][]string)
	option := setHeaderPrefix("test-key", "test-value")

	err := option(headers, false)
	assert.NoError(t, err)
	assert.Contains(t, headers, "x-amz-test-key")
	assert.Equal(t, "test-value", headers["x-amz-test-key"][0])
}

func TestSetHeaderPrefix_ShouldSetObsHeader_WhenIsObsIsTrue(t *testing.T) {
	headers := make(map[string][]string)
	option := setHeaderPrefix("test-key", "test-value")

	err := option(headers, true)
	assert.NoError(t, err)
	assert.Contains(t, headers, "x-obs-test-key")
	assert.Equal(t, "test-value", headers["x-obs-test-key"][0])
}

func TestSetHeaderPrefix_ShouldReturnError_WhenValueIsEmpty(t *testing.T) {
	headers := make(map[string][]string)
	option := setHeaderPrefix("test-key", "")

	err := option(headers, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty value")
}

func TestSetHeaderPrefix_ShouldReturnError_WhenValueIsWhitespace(t *testing.T) {
	headers := make(map[string][]string)
	option := setHeaderPrefix("test-key", "   ")

	err := option(headers, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty value")
}

// WithReqPaymentHeader Tests

func TestWithReqPaymentHeader_ShouldSetHeader_WhenCalled(t *testing.T) {
	headers := make(map[string][]string)
	option := WithReqPaymentHeader(Requester)

	err := option(headers, false)
	assert.NoError(t, err)
	assert.Contains(t, headers, "x-amz-request-payer")
	assert.Equal(t, "requester", headers["x-amz-request-payer"][0])
}

func TestWithReqPaymentHeader_ShouldSetObsHeader_WhenIsObsIsTrue(t *testing.T) {
	headers := make(map[string][]string)
	option := WithReqPaymentHeader(Requester)

	err := option(headers, true)
	assert.NoError(t, err)
	assert.Contains(t, headers, "x-obs-request-payer")
}

// WithTrafficLimitHeader Tests

func TestWithTrafficLimitHeader_ShouldSetHeader_WhenCalled(t *testing.T) {
	headers := make(map[string][]string)
	option := WithTrafficLimitHeader(819200)

	err := option(headers, false)
	assert.NoError(t, err)
	assert.Contains(t, headers, "x-amz-traffic-limit")
	assert.Equal(t, "819200", headers["x-amz-traffic-limit"][0])
}

func TestWithTrafficLimitHeader_ShouldHandleNegativeValue(t *testing.T) {
	headers := make(map[string][]string)
	option := WithTrafficLimitHeader(-100)

	err := option(headers, false)
	assert.NoError(t, err)
	assert.Contains(t, headers, "x-amz-traffic-limit")
	assert.Equal(t, "-100", headers["x-amz-traffic-limit"][0])
}

// WithCallbackHeader Tests

func TestWithCallbackHeader_ShouldSetHeader_WhenUrlIsValid(t *testing.T) {
	headers := make(map[string][]string)
	callbackUrl := "http://example.com/callback"
	option := WithCallbackHeader(callbackUrl)

	err := option(headers, false)
	assert.NoError(t, err)
	assert.Contains(t, headers, "x-amz-callback")
	assert.Equal(t, callbackUrl, headers["x-amz-callback"][0])
}

func TestWithCallbackHeader_ShouldReturnError_WhenUrlIsEmpty(t *testing.T) {
	headers := make(map[string][]string)
	option := WithCallbackHeader("")

	err := option(headers, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty value")
}

// WithCustomHeader Tests

func TestWithCustomHeader_ShouldSetHeader_WhenKeyAndValueAreValid(t *testing.T) {
	headers := make(map[string][]string)
	option := WithCustomHeader("X-Custom-Header", "custom-value")

	err := option(headers, false)
	assert.NoError(t, err)
	assert.Equal(t, "custom-value", headers["X-Custom-Header"][0])
}

func TestWithCustomHeader_ShouldReturnError_WhenValueIsEmpty(t *testing.T) {
	headers := make(map[string][]string)
	option := WithCustomHeader("X-Custom-Header", "")

	err := option(headers, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty value")
}

func TestWithCustomHeader_ShouldReturnError_WhenValueIsWhitespace(t *testing.T) {
	headers := make(map[string][]string)
	option := WithCustomHeader("X-Custom-Header", "   ")

	err := option(headers, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty value")
}

func TestWithCustomHeader_ShouldNotAddPrefix_WhenKeyIsProvided(t *testing.T) {
	headers := make(map[string][]string)
	option := WithCustomHeader("Custom-Header", "custom-value")

	err := option(headers, true) // isObs = true
	assert.NoError(t, err)
	assert.Equal(t, "custom-value", headers["Custom-Header"][0])
	assert.NotContains(t, headers, "x-obs-Custom-Header")
}

// removeEndNewlineCharacter Tests

func TestRemoveEndNewlineCharacter_ShouldRemoveLastChar_WhenBufferEnds(t *testing.T) {
	testData := "test data"
	callbackBuffer := bytes.NewBufferString(testData + "\n")

	result := removeEndNewlineCharacter(callbackBuffer)
	assert.Equal(t, []byte(testData), result)
}

func TestRemoveEndNewlineCharacter_ShouldReturnAllButLast_WhenBufferNotEmpty(t *testing.T) {
	testData := "callback body data"
	callbackBuffer := bytes.NewBufferString(testData)

	result := removeEndNewlineCharacter(callbackBuffer)
	assert.Equal(t, []byte(testData)[:len(testData)-1], result)
}

// PreprocessCallbackInputToSHA256 Tests

func TestPreprocessCallbackInputToSHA256_ShouldReturnHash_WhenInputIsValid(t *testing.T) {
	input := &CallbackInput{
		CallbackUrl:  "http://example.com/callback",
		CallbackBody: "{\"bucket\":\"test-bucket\",\"object\":\"test-object\"}",
	}

	hash, err := PreprocessCallbackInputToSHA256(input)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Equal(t, 64, len(hash)) // SHA256 produces 64 hex characters
}

func TestPreprocessCallbackInputToSHA256_ShouldReturnError_WhenInputIsNil(t *testing.T) {
	hash, err := PreprocessCallbackInputToSHA256(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parameter can not be nil")
	assert.Empty(t, hash)
}

func TestPreprocessCallbackInputToSHA256_ShouldReturnError_WhenUrlIsEmpty(t *testing.T) {
	input := &CallbackInput{
		CallbackUrl:  "",
		CallbackBody: "{\"bucket\":\"test-bucket\"}",
	}

	hash, err := PreprocessCallbackInputToSHA256(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CallbackUrl")
	assert.Empty(t, hash)
}

func TestPreprocessCallbackInputToSHA256_ShouldReturnError_WhenBodyIsEmpty(t *testing.T) {
	input := &CallbackInput{
		CallbackUrl:  "http://example.com/callback",
		CallbackBody: "",
	}

	hash, err := PreprocessCallbackInputToSHA256(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CallbackBody")
	assert.Empty(t, hash)
}

func TestPreprocessCallbackInputToSHA256_ShouldProduceSameHash_WhenInputSame(t *testing.T) {
	input := &CallbackInput{
		CallbackUrl:  "http://example.com/callback",
		CallbackBody: "{\"bucket\":\"test\",\"object\":\"file.txt\"}",
	}

	hash1, err1 := PreprocessCallbackInputToSHA256(input)
	hash2, err2 := PreprocessCallbackInputToSHA256(input)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, hash1, hash2)
}

func TestPreprocessCallbackInputToSHA256_ShouldHandleComplexBody(t *testing.T) {
	input := &CallbackInput{
		CallbackUrl:      "http://example.com/callback",
		CallbackBody:     "{\"bucket\":\"test\",\"object\":\"file.txt\",\"etag\":\"abc123\"}",
		CallbackHost:     "example.com",
		CallbackBodyType: "application/json",
	}

	hash, err := PreprocessCallbackInputToSHA256(input)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Equal(t, 64, len(hash))
}

// WithProgress function call tests to achieve 100% coverage

func TestWithProgress_ShouldReturnListener_WhenFunctionCalled(t *testing.T) {
	listener := NewMockProgressListener()
	option := WithProgress(listener)

	// Call the returned function to cover the closure
	result := option()
	assert.Equal(t, listener, result)
}

func TestWithProgress_ShouldReturnNil_WhenListenerIsNil(t *testing.T) {
	option := WithProgress(nil)

	// Call the returned function to cover the closure
	result := option()
	assert.Nil(t, result)
}

func TestWithProgress_ShouldHandleMultipleCalls(t *testing.T) {
	listener := NewMockProgressListener()
	option := WithProgress(listener)

	// Multiple calls should return the same listener
	result1 := option()
	result2 := option()
	assert.Equal(t, listener, result1)
	assert.Equal(t, result1, result2)
}
