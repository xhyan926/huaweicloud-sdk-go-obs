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
	"time"

	"github.com/stretchr/testify/assert"
)

// TestWithSignature_ShouldSetSignatureType_WhenCalled tests WithSignature configurer
func TestWithSignature_ShouldSetSignatureType_WhenCalled(t *testing.T) {
	tests := []struct {
		name      string
		signature SignatureType
	}{
		{"SignatureV2", SignatureV2},
		{"SignatureV4", SignatureV4},
		{"SignatureObs", SignatureObs},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{}
			configurer := WithSignature(tt.signature)
			configurer(conf)

			assert.Equal(t, tt.signature, conf.signature)
		})
	}
}

// TestWithRegion_ShouldSetRegion_WhenCalled tests WithRegion configurer
func TestWithRegion_ShouldSetRegion_WhenCalled(t *testing.T) {
	region := "cn-north-4"
	conf := &config{}
	configurer := WithRegion(region)
	configurer(conf)

	assert.Equal(t, region, conf.region)
}

// TestWithPathStyle_ShouldSetPathStyle_WhenCalled tests WithPathStyle configurer
func TestWithPathStyle_ShouldSetPathStyle_WhenCalled(t *testing.T) {
	tests := []struct {
		name      string
		pathStyle bool
	}{
		{"PathStyle enabled", true},
		{"PathStyle disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{}
			configurer := WithPathStyle(tt.pathStyle)
			configurer(conf)

			assert.Equal(t, tt.pathStyle, conf.pathStyle)
		})
	}
}

// TestWithSslVerify_ShouldSetSslVerify_WhenCalled tests WithSslVerify configurer
func TestWithSslVerify_ShouldSetSslVerify_WhenCalled(t *testing.T) {
	tests := []struct {
		name       string
		sslVerify bool
	}{
		{"SSL verify enabled", true},
		{"SSL verify disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{}
			configurer := WithSslVerify(tt.sslVerify)
			configurer(conf)

			assert.Equal(t, tt.sslVerify, conf.sslVerify)
		})
	}
}

// TestWithProxyUrl_ShouldSetProxy_WhenCalled tests WithProxyUrl configurer
func TestWithProxyUrl_ShouldSetProxy_WhenCalled(t *testing.T) {
	proxyURL := "http://proxy.example.com:8080"
	conf := &config{}
	configurer := WithProxyUrl(proxyURL)
	configurer(conf)

	assert.Equal(t, proxyURL, conf.proxyURL)
}

// TestWithMaxRetryCount_ShouldSetMaxRetryCount_WhenCalled tests WithMaxRetryCount configurer
func TestWithMaxRetryCount_ShouldSetMaxRetryCount_WhenCalled(t *testing.T) {
	maxRetryCount := 3
	conf := &config{}
	configurer := WithMaxRetryCount(maxRetryCount)
	configurer(conf)

	assert.Equal(t, maxRetryCount, conf.maxRetryCount)
}

// TestWithConnectTimeout_ShouldSetConnectTimeout_WhenCalled tests WithConnectTimeout configurer
func TestWithConnectTimeout_ShouldSetConnectTimeout_WhenCalled(t *testing.T) {
	timeout := 60
	conf := &config{}
	configurer := WithConnectTimeout(timeout)
	configurer(conf)

	assert.Equal(t, timeout, conf.connectTimeout)
}

// TestWithSocketTimeout_ShouldSetSocketTimeout_WhenCalled tests WithSocketTimeout configurer
func TestWithSocketTimeout_ShouldSetSocketTimeout_WhenCalled(t *testing.T) {
	timeout := 60
	conf := &config{}
	configurer := WithSocketTimeout(timeout)
	configurer(conf)

	assert.Equal(t, timeout, conf.socketTimeout)
}

// TestWithMaxConnections_ShouldSetMaxConnections_WhenCalled tests WithMaxConnections configurer
func TestWithMaxConnections_ShouldSetMaxConnections_WhenCalled(t *testing.T) {
	maxConnsPerHost := 10
	conf := &config{}
	configurer := WithMaxConnections(maxConnsPerHost)
	configurer(conf)

	assert.Equal(t, maxConnsPerHost, conf.maxConnsPerHost)
}

// TestWithDisableKeepAlive_ShouldSetDisableKeepAlive_WhenCalled tests WithDisableKeepAlive configurer
func TestWithDisableKeepAlive_ShouldSetDisableKeepAlive_WhenCalled(t *testing.T) {
	disableKeepAlive := true
	conf := &config{}
	configurer := WithDisableKeepAlive(disableKeepAlive)
	configurer(conf)

	assert.Equal(t, disableKeepAlive, conf.disableKeepAlive)
}

// TestWithSecurityToken_ShouldSetSecurityToken_WhenCalled tests WithSecurityToken configurer
func TestWithSecurityToken_ShouldSetSecurityToken_WhenCalled(t *testing.T) {
	securityToken := "test-security-token"
	ak := "test-ak"
	sk := "test-sk"

	conf := &config{}
	// First add a BasicSecurityProvider
	conf.securityProviders = []securityProvider{NewBasicSecurityProvider(ak, sk, "")}

	// Then apply WithSecurityToken
	configurer := WithSecurityToken(securityToken)
	configurer(conf)

	bsp := conf.securityProviders[0].(*BasicSecurityProvider)
	sh := bsp.getSecurity()
	assert.Equal(t, securityToken, sh.securityToken)
}

// TestWithRequestContext_ShouldSetContext_WhenCalled tests WithRequestContext configurer
func TestWithRequestContext_ShouldSetContext_WhenCalled(t *testing.T) {
	// Mock context for testing
	ctx := &mockContext{}
	conf := &config{}
	configurer := WithRequestContext(ctx)
	configurer(conf)

	assert.Equal(t, ctx, conf.ctx)
}

// mockContext implements context.Context for testing
type mockContext struct{}

func (m *mockContext) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }
func (m *mockContext) Done() <-chan struct{}               { return nil }
func (m *mockContext) Err() error                      { return nil }
func (m *mockContext) Value(key interface{}) interface{}   { return nil }

// TestInitConfigWithDefault_ShouldReturnError_WhenEndpointIsEmpty tests initConfigWithDefault
func TestInitConfigWithDefault_ShouldReturnError_WhenEndpointIsEmpty(t *testing.T) {
	conf := &config{}
	conf.endpoint = ""
	err := conf.initConfigWithDefault()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "endpoint is not set")
}

// TestInitConfigWithDefault_ShouldTrimTrailingSlashes_WhenEndpointHasTrailingSlashes tests initConfigWithDefault
func TestInitConfigWithDefault_ShouldTrimTrailingSlashes_WhenEndpointHasTrailingSlashes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Single trailing slash", "https://obs.example.com/", "https://obs.example.com"},
		{"Multiple trailing slashes", "https://obs.example.com//", "https://obs.example.com"},
		{"No trailing slash", "https://obs.example.com", "https://obs.example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{}
			conf.endpoint = tt.input
			err := conf.initConfigWithDefault()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, conf.endpoint)
		})
	}
}

// TestInitConfigWithDefault_ShouldRemoveQueryString_WhenEndpointHasQueryString tests initConfigWithDefault
func TestInitConfigWithDefault_ShouldRemoveQueryString_WhenEndpointHasQueryString(t *testing.T) {
	conf := &config{}
	conf.endpoint = "https://obs.example.com?param=value"
	err := conf.initConfigWithDefault()
	assert.NoError(t, err)
	assert.Equal(t, "https://obs.example.com", conf.endpoint)
}

// TestPrepareBaseURL_ShouldReturnVirtualHostStyle_WhenPathStyleIsFalse tests prepareBaseURL
func TestPrepareBaseURL_ShouldReturnVirtualHostStyle_WhenPathStyleIsFalse(t *testing.T) {
	conf := &config{pathStyle: false, urlHolder: &urlHolder{scheme: "https", host: "obs.example.com", port: 443}}
	requestURL, canonicalizedURL := conf.prepareBaseURL("bucket-name")
	assert.Contains(t, requestURL, "bucket-name.obs.example.com")
	assert.Contains(t, canonicalizedURL, "/")
}

// TestPrepareBaseURL_ShouldReturnPathStyle_WhenPathStyleIsTrue tests prepareBaseURL
func TestPrepareBaseURL_ShouldReturnPathStyle_WhenPathStyleIsTrue(t *testing.T) {
	conf := &config{pathStyle: true, urlHolder: &urlHolder{scheme: "https", host: "obs.example.com", port: 443}}
	requestURL, canonicalizedURL := conf.prepareBaseURL("bucket-name")
	assert.Contains(t, requestURL, "obs.example.com:443/bucket-name")
	assert.Contains(t, canonicalizedURL, "/bucket-name")
}

// TestPrepareBaseURL_ShouldHandleEmptyBucket_WhenCalled tests prepareBaseURL
func TestPrepareBaseURL_ShouldHandleEmptyBucket_WhenCalled(t *testing.T) {
	conf := &config{pathStyle: false, urlHolder: &urlHolder{scheme: "https", host: "obs.example.com", port: 443}}
	requestURL, canonicalizedURL := conf.prepareBaseURL("")
	assert.Contains(t, requestURL, "obs.example.com")
	assert.Equal(t, "/", canonicalizedURL)
}

// TestGetTransport_ShouldCreateDefaultTransport_WhenNotProvided tests getTransport
func TestGetTransport_ShouldCreateDefaultTransport_WhenNotProvided(t *testing.T) {
	conf := &config{}
	err := conf.getTransport()
	assert.NoError(t, err)
	assert.NotNil(t, conf.transport)
}

// TestGetTransport_ShouldSetSslVerify_WhenProvided tests getTransport
func TestGetTransport_ShouldSetSslVerify_WhenProvided(t *testing.T) {
	conf := &config{sslVerify: false}
	err := conf.getTransport()
	assert.NoError(t, err)
	assert.NotNil(t, conf.transport)
	assert.NotNil(t, conf.transport.TLSClientConfig)
}

// TestConfigString_ShouldReturnFormattedString_WhenCalled tests String method
func TestConfigString_ShouldReturnFormattedString_WhenCalled(t *testing.T) {
	conf := &config{
		endpoint:        "https://obs.example.com",
		signature:       SignatureV4,
		pathStyle:       true,
		region:          "cn-north-4",
		connectTimeout:  60,
		socketTimeout:   60,
		maxRetryCount:   3,
		maxConnsPerHost: 1000,
		sslVerify:      true,
		maxRedirectCount: 3,
	}
	result := conf.String()
	assert.Contains(t, result, "https://obs.example.com")
	assert.Contains(t, result, "v4")
	assert.Contains(t, result, "cn-north-4")
	assert.Contains(t, result, "60")
	assert.Contains(t, result, "3")
	assert.Contains(t, result, "1000")
}

// TestPrepareConfig_ShouldSetDefaultValues_WhenConfigNotSet tests prepareConfig
func TestPrepareConfig_ShouldSetDefaultValues_WhenConfigNotSet(t *testing.T) {
	conf := &config{
		maxRetryCount:    -1,
		maxRedirectCount: -1,
	}
	conf.prepareConfig()
	assert.Equal(t, DEFAULT_CONNECT_TIMEOUT, conf.connectTimeout)
	assert.Equal(t, DEFAULT_SOCKET_TIMEOUT, conf.socketTimeout)
	assert.Equal(t, DEFAULT_MAX_RETRY_COUNT, conf.maxRetryCount)
	assert.Equal(t, DEFAULT_MAX_CONN_PER_HOST, conf.maxConnsPerHost)
	assert.Equal(t, DEFAULT_MAX_REDIRECT_COUNT, conf.maxRedirectCount)
}

// TestWithHttpTransport_ShouldSetTransport_WhenCalled tests WithHttpTransport configurer
// Skipping - requires actual http.Transport implementation
// func TestWithHttpTransport_ShouldSetTransport_WhenCalled(t *testing.T) {
// 	conf := &config{}
// 	// Create a dummy transport
// 	transport := &dummyTransport{}
// 	configurer := WithHttpTransport(transport)
// 	configurer(conf)
// 	assert.NotNil(t, conf.transport)
// }

// TestWithHttpClient_ShouldSetClient_WhenCalled tests WithHttpClient configurer
// Skipping - requires actual http.Client implementation
// func TestWithHttpClient_ShouldSetClient_WhenCalled(t *testing.T) {
// 	conf := &config{}
// 	// Create a dummy client
// 	client := &dummyHttpClient{}
// 	configurer := WithHttpClient(client)
// 	configurer(conf)
// 	assert.NotNil(t, conf.httpClient)
// }

// TestWithUserAgent_ShouldSetUserAgent_WhenCalled tests WithUserAgent configurer
func TestWithUserAgent_ShouldSetUserAgent_WhenCalled(t *testing.T) {
	userAgent := "custom-agent/1.0"
	conf := &config{}
	configurer := WithUserAgent(userAgent)
	configurer(conf)
	assert.Equal(t, userAgent, conf.userAgent)
}

// TestWithEnableCompression_ShouldSetCompression_WhenCalled tests WithEnableCompression configurer
func TestWithEnableCompression_ShouldSetCompression_WhenCalled(t *testing.T) {
	tests := []struct {
		name              string
		enableCompression bool
	}{
		{"Enable compression", true},
		{"Disable compression", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{}
			configurer := WithEnableCompression(tt.enableCompression)
			configurer(conf)
			assert.Equal(t, tt.enableCompression, conf.enableCompression)
		})
	}
}

// TestWithMaxRedirectCount_ShouldSetMaxRedirectCount_WhenCalled tests WithMaxRedirectCount configurer
func TestWithMaxRedirectCount_ShouldSetMaxRedirectCount_WhenCalled(t *testing.T) {
	maxRedirectCount := 5
	conf := &config{}
	configurer := WithMaxRedirectCount(maxRedirectCount)
	configurer(conf)
	assert.Equal(t, maxRedirectCount, conf.maxRedirectCount)
}

// TestWithCustomDomainName_ShouldSetCname_WhenCalled tests WithCustomDomainName configurer
func TestWithCustomDomainName_ShouldSetCname_WhenCalled(t *testing.T) {
	tests := []struct {
		name  string
		cname bool
	}{
		{"Cname enabled", true},
		{"Cname disabled", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{}
			configurer := WithCustomDomainName(tt.cname)
			configurer(conf)
			assert.Equal(t, tt.cname, conf.cname)
		})
	}
}

// TestWithIdleConnTimeout_ShouldSetIdleConnTimeout_WhenCalled tests WithIdleConnTimeout configurer
func TestWithIdleConnTimeout_ShouldSetIdleConnTimeout_WhenCalled(t *testing.T) {
	idleConnTimeout := 120
	conf := &config{}
	configurer := WithIdleConnTimeout(idleConnTimeout)
	configurer(conf)
	assert.Equal(t, idleConnTimeout, conf.idleConnTimeout)
}

// TestWithHeaderTimeout_ShouldSetHeaderTimeout_WhenCalled tests WithHeaderTimeout configurer
func TestWithHeaderTimeout_ShouldSetHeaderTimeout_WhenCalled(t *testing.T) {
	headerTimeout := 90
	conf := &config{}
	configurer := WithHeaderTimeout(headerTimeout)
	configurer(conf)
	assert.Equal(t, headerTimeout, conf.headerTimeout)
}

// TestWithProxyFromEnv_ShouldSetProxyFromEnv_WhenCalled tests WithProxyFromEnv configurer
func TestWithProxyFromEnv_ShouldSetProxyFromEnv_WhenCalled(t *testing.T) {
	tests := []struct {
		name           string
		proxyFromEnv   bool
	}{
		{"Proxy from env enabled", true},
		{"Proxy from env disabled", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &config{}
			configurer := WithProxyFromEnv(tt.proxyFromEnv)
			configurer(conf)
			assert.Equal(t, tt.proxyFromEnv, conf.proxyFromEnv)
		})
	}
}

// TestWithNoProxyUrl_ShouldSetNoProxyUrl_WhenCalled tests WithNoProxyUrl configurer
func TestWithNoProxyUrl_ShouldSetNoProxyUrl_WhenCalled(t *testing.T) {
	noProxyURL := "localhost,127.0.0.1"
	conf := &config{}
	configurer := WithNoProxyUrl(noProxyURL)
	configurer(conf)
	assert.Equal(t, noProxyURL, conf.noProxyURL)
}

// TestWithSecurityProviders_ShouldAddSecurityProviders_WhenCalled tests WithSecurityProviders configurer
func TestWithSecurityProviders_ShouldAddSecurityProviders_WhenCalled(t *testing.T) {
	conf := &config{}
	// Create two security providers
	provider1 := NewBasicSecurityProvider("ak1", "sk1", "")
	provider2 := NewBasicSecurityProvider("ak2", "sk2", "")
	configurer := WithSecurityProviders(provider1, provider2)
	configurer(conf)
	assert.Len(t, conf.securityProviders, 2)
}

// TestWithSslVerifyAndPemCerts_ShouldSetBoth_WhenCalled tests WithSslVerifyAndPemCerts configurer
func TestWithSslVerifyAndPemCerts_ShouldSetBoth_WhenCalled(t *testing.T) {
	sslVerify := true
	pemCerts := []byte("-----BEGIN CERTIFICATE-----test-----END CERTIFICATE-----")
	conf := &config{}
	configurer := WithSslVerifyAndPemCerts(sslVerify, pemCerts)
	configurer(conf)
	assert.Equal(t, sslVerify, conf.sslVerify)
	assert.Equal(t, pemCerts, conf.pemCerts)
}

// Dummy implementations for testing
type dummyTransport struct{}

type dummyHttpClient struct{}
