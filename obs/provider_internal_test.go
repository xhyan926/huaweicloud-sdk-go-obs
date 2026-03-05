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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// BasicSecurityProvider Tests

func TestBasicSecurityProvider_GetSecurity_ShouldReturnStoredCredentials(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	securityToken := "test-token"

	provider := NewBasicSecurityProvider(ak, sk, securityToken)
	securityHolder := provider.getSecurity()

	assert.Equal(t, ak, securityHolder.ak)
	assert.Equal(t, sk, securityHolder.sk)
	assert.Equal(t, securityToken, securityHolder.securityToken)
}

func TestBasicSecurityProvider_Refresh_ShouldUpdateCredentials(t *testing.T) {
	provider := NewBasicSecurityProvider("old-ak", "old-sk", "old-token")
	provider.refresh("new-ak", "new-sk", "new-token")

	securityHolder := provider.getSecurity()
	assert.Equal(t, "new-ak", securityHolder.ak)
	assert.Equal(t, "new-sk", securityHolder.sk)
	assert.Equal(t, "new-token", securityHolder.securityToken)
}

// EnvSecurityProvider Tests

func TestEnvSecurityProvider_GetSecurity_ShouldReadFromEnv(t *testing.T) {
	provider := NewEnvSecurityProvider("")
	_ = provider.getSecurity()
	// Test passes if no panic occurs
}

func TestNewEnvSecurityProvider_ShouldHandleSuffix(t *testing.T) {
	provider := NewEnvSecurityProvider("test")
	assert.Equal(t, "_test", provider.suffix)
}

func TestNewEnvSecurityProvider_ShouldHandleEmptySuffix(t *testing.T) {
	provider := NewEnvSecurityProvider("")
	assert.Equal(t, "", provider.suffix)
}

// EcsSecurityProvider Tests

func TestEcsSecurityProvider_LoadTemporarySecurityHolder_ShouldReturnEmpty_WhenNotLoaded(t *testing.T) {
	provider := NewEcsSecurityProvider(0)
	holder, success := provider.loadTemporarySecurityHolder()

	assert.False(t, success)
	assert.Equal(t, emptyTemporarySecurityHolder, holder)
}

func TestEcsSecurityProvider_GetAndSetSecurityWithOutLock_ShouldHandleRequestError(t *testing.T) {
	// Create mock transport that returns error
	mockClient := &http.Client{
		Transport: &errorTransport{},
	}

	provider := &EcsSecurityProvider{
		retryCount: 0,
		httpClient: mockClient,
	}

	holder := provider.getAndSetSecurityWithOutLock()
	assert.Equal(t, "", holder.ak)
	assert.Equal(t, "", holder.sk)
	assert.Equal(t, "", holder.securityToken)
}

func TestEcsSecurityProvider_GetAndSetSecurityWithOutLock_ShouldHandleInvalidJSON(t *testing.T) {
	mockClient := &http.Client{
		Transport: &jsonTransport{jsonStr: "invalid json"},
	}

	provider := &EcsSecurityProvider{
		retryCount: 0,
		httpClient: mockClient,
	}

	holder := provider.getAndSetSecurityWithOutLock()
	assert.Equal(t, "", holder.ak)
}

func TestEcsSecurityProvider_GetAndSetSecurityWithOutLock_ShouldParseValidResponse(t *testing.T) {
	credential := struct {
		Credential struct {
			AK            string    `json:"access,omitempty"`
			SK            string    `json:"secret,omitempty"`
			SecurityToken string    `json:"securitytoken,omitempty"`
			ExpireDate    time.Time `json:"expires_at,omitempty"`
		} `json:"credential"`
	}{}

	credential.Credential.AK = "test-ak"
	credential.Credential.SK = "test-sk"
	credential.Credential.SecurityToken = "test-token"
	credential.Credential.ExpireDate = time.Now().Add(1 * time.Hour)

	responseBody, _ := json.Marshal(credential)

	mockClient := &http.Client{
		Transport: &jsonTransport{jsonStr: string(responseBody)},
	}

	provider := &EcsSecurityProvider{
		retryCount: 1,
		httpClient: mockClient,
	}

	holder := provider.getAndSetSecurityWithOutLock()
	assert.Equal(t, "test-ak", holder.ak)
	assert.Equal(t, "test-sk", holder.sk)
	assert.Equal(t, "test-token", holder.securityToken)
}

func TestEcsSecurityProvider_GetAndSetSecurity_ShouldHandleEmptyHolder(t *testing.T) {
	provider := NewEcsSecurityProvider(0)
	holder := provider.getAndSetSecurity()

	assert.Equal(t, "", holder.ak)
	assert.Equal(t, "", holder.sk)
	assert.Equal(t, "", holder.securityToken)
}

func TestEcsSecurityProvider_GetSecurity_ShouldReturnCached_WhenNotExpired(t *testing.T) {
	provider := &EcsSecurityProvider{}

	// Set a cached value
	holder := TemporarySecurityHolder{
		securityHolder: securityHolder{
			ak:            "cached-ak",
			sk:            "cached-sk",
			securityToken: "cached-token",
		},
		expireDate: time.Now().Add(1 * time.Hour),
	}
	provider.val.Store(holder)

	result := provider.getSecurity()
	assert.Equal(t, "cached-ak", result.ak)
	assert.Equal(t, "cached-sk", result.sk)
	assert.Equal(t, "cached-token", result.securityToken)
}

func TestEcsSecurityProvider_GetSecurity_ShouldRefresh_WhenExpired(t *testing.T) {
	// Create a mock response
	credential := struct {
		Credential struct {
			AK            string    `json:"access,omitempty"`
			SK            string    `json:"secret,omitempty"`
			SecurityToken string    `json:"securitytoken,omitempty"`
			ExpireDate    time.Time `json:"expires_at,omitempty"`
		} `json:"credential"`
	}{}

	credential.Credential.AK = "new-ak"
	credential.Credential.SK = "new-sk"
	credential.Credential.SecurityToken = "new-token"
	credential.Credential.ExpireDate = time.Now().Add(1 * time.Hour)

	responseBody, _ := json.Marshal(credential)

	mockClient := &http.Client{
		Transport: &jsonTransport{jsonStr: string(responseBody)},
	}

	provider := &EcsSecurityProvider{
		retryCount: 0,
		httpClient: mockClient,
	}

	// Set an expired cached value
	holder := TemporarySecurityHolder{
		securityHolder: securityHolder{
			ak:            "old-ak",
			sk:            "old-sk",
			securityToken: "old-token",
		},
		expireDate: time.Now().Add(-1 * time.Hour),
	}
	provider.val.Store(holder)

	result := provider.getSecurity()
	assert.Equal(t, "new-ak", result.ak)
	assert.Equal(t, "new-sk", result.sk)
	assert.Equal(t, "new-token", result.securityToken)
}

func TestNewEcsSecurityProvider_ShouldCreateProvider(t *testing.T) {
	retryCount := 3
	provider := NewEcsSecurityProvider(retryCount)

	assert.NotNil(t, provider)
	assert.Equal(t, retryCount, provider.retryCount)
	assert.NotNil(t, provider.httpClient)
}

func TestNewBasicSecurityProvider_ShouldCreateProvider(t *testing.T) {
	ak := "test-ak"
	sk := "test-sk"
	token := "test-token"

	provider := NewBasicSecurityProvider(ak, sk, token)
	assert.NotNil(t, provider)

	holder := provider.getSecurity()
	assert.Equal(t, ak, holder.ak)
	assert.Equal(t, sk, holder.sk)
	assert.Equal(t, token, holder.securityToken)
}

// Mock transports for testing

type errorTransport struct{}

func (et *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, assert.AnError
}

type jsonTransport struct {
	jsonStr string
}

func (jt *jsonTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(jt.jsonStr))),
		Header:     make(http.Header),
	}, nil
}
