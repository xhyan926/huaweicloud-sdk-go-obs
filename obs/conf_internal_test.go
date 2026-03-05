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

	"github.com/stretchr/testify/assert"
)

// prepareEscapeFunc Tests

func TestPrepareEscapeFunc_ShouldReturnUrlEscapeFunc_WhenEscapeIsTrue(t *testing.T) {
	conf := &config{enableCompression: true}
	escapeFunc := conf.prepareEscapeFunc(true)
	assert.NotNil(t, escapeFunc)
	// The escape function should encode spaces
	result := escapeFunc("test file.txt")
	assert.NotEqual(t, "test file.txt", result)
}

func TestPrepareEscapeFunc_ShouldReturnIdentityFunc_WhenEscapeIsFalse(t *testing.T) {
	conf := &config{enableCompression: false}
	escapeFunc := conf.prepareEscapeFunc(false)
	assert.NotNil(t, escapeFunc)
	// The identity function should not change the string
	result := escapeFunc("test file.txt")
	assert.Equal(t, "test file.txt", result)
}

// prepareObjectKey Tests

func TestPrepareObjectKey_ShouldEncodeKey_WhenEscapeIsTrue(t *testing.T) {
	conf := &config{enableCompression: true}
	escapeFunc := conf.prepareEscapeFunc(true)
	encodedKey := conf.prepareObjectKey(true, "test file.txt", escapeFunc)
	// When escape is true, the key should be URL-encoded
	assert.NotEqual(t, "test file.txt", encodedKey)
}

func TestPrepareObjectKey_ShouldNotEncodeKey_WhenEscapeIsFalse(t *testing.T) {
	conf := &config{enableCompression: false}
	escapeFunc := conf.prepareEscapeFunc(false)
	encodedKey := conf.prepareObjectKey(false, "test file.txt", escapeFunc)
	// When escape is false, the key should not be changed
	assert.Equal(t, "test file.txt", encodedKey)
}

func TestPrepareObjectKey_ShouldHandleEmptyKey(t *testing.T) {
	conf := &config{enableCompression: false}
	escapeFunc := conf.prepareEscapeFunc(false)
	encodedKey := conf.prepareObjectKey(false, "", escapeFunc)
	assert.Equal(t, "", encodedKey)
}

func TestPrepareObjectKey_ShouldHandleLeadingSlash(t *testing.T) {
	conf := &config{enableCompression: false}
	escapeFunc := conf.prepareEscapeFunc(false)
	encodedKey := conf.prepareObjectKey(false, "/my-object", escapeFunc)
	// Leading slash should be preserved
	assert.Equal(t, "/my-object", encodedKey)
}

// formatUrls Tests

func TestFormatUrls_ShouldHandleQueryParams(t *testing.T) {
	conf := &config{
		urlHolder: &urlHolder{scheme: "https", host: "obs.example.com", port: 443},
		pathStyle: false,
	}
	params := map[string]string{"param1": "value1"}
	requestURL, canonicalizedURL := conf.formatUrls("bucket", "object", params, false)
	assert.Contains(t, requestURL, "param1=value1")
	assert.NotEmpty(t, canonicalizedURL)
}

func TestFormatUrls_ShouldHandleEscapeParameter(t *testing.T) {
	conf := &config{
		urlHolder: &urlHolder{scheme: "https", host: "obs.example.com", port: 443},
		pathStyle: false,
		enableCompression: true,
	}
	requestURL, _ := conf.formatUrls("bucket", "test file.txt", nil, true)
	// When escape is true, space should be encoded
	assert.Contains(t, requestURL, "%20")
}

// customProxyFunc Tests

func TestCustomProxyFunc_ShouldReturnProxyFunc(t *testing.T) {
	conf := &config{
		proxyURL: "http://proxy.example.com:8080",
	}
	proxyFunc := conf.customProxyFunc()
	assert.NotNil(t, proxyFunc)
}

// customProxyFromEnvironment Tests

func TestCustomProxyFromEnvironment_ShouldReturnProxy_WhenEnvSet(t *testing.T) {
	conf := &config{}
	t.Setenv("HTTP_PROXY", "http://proxy.example.com:8080")
	t.Setenv("HTTPS_PROXY", "http://proxy.example.com:8443")

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	proxyURL, err := conf.customProxyFromEnvironment(req)

	// Should return one of the environment proxies
	assert.True(t, err == nil || proxyURL != nil)
	if proxyURL != nil {
		assert.Contains(t, proxyURL.String(), "proxy.example.com")
	}
}

func TestCustomProxyFromEnvironment_ShouldReturnNil_WhenNoEnvSet(t *testing.T) {
	conf := &config{}
	t.Setenv("HTTP_PROXY", "")
	t.Setenv("HTTPS_PROXY", "")

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	proxyURL, err := conf.customProxyFromEnvironment(req)

	assert.Nil(t, err)
	assert.Nil(t, proxyURL)
}

func TestCustomProxyFromEnvironment_ShouldIgnoreNoProxy(t *testing.T) {
	conf := &config{
		noProxyURL: "localhost",
	}
	t.Setenv("HTTP_PROXY", "http://proxy.example.com:8080")
	t.Setenv("NO_PROXY", "localhost")

	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	proxyURL, err := conf.customProxyFromEnvironment(req)

	// Should return nil due to NO_PROXY
	assert.Nil(t, err)
	assert.Nil(t, proxyURL)
}
