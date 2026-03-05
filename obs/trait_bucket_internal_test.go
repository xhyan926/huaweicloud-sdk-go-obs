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
	"testing"

	"github.com/stretchr/testify/assert"
)

// DeleteBucketCustomDomainInput trans Tests

func TestDeleteBucketCustomDomainInput_Trans_ShouldReturnParamsWithSubResource(t *testing.T) {
	input := DeleteBucketCustomDomainInput{
		Bucket:       "test-bucket",
		CustomDomain: "example.com",
	}

	params, headers, data, err := input.trans(false)

	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.Len(t, params, 1)
	assert.Nil(t, headers)
	assert.NotNil(t, data)
}

func TestDeleteBucketCustomDomainInput_Trans_ShouldReturnSameResult_WhenIsObs(t *testing.T) {
	input := DeleteBucketCustomDomainInput{
		Bucket:       "test-bucket",
		CustomDomain: "example.com",
	}

	params, headers, data, err := input.trans(true)

	assert.NoError(t, err)
	assert.NotNil(t, params)
	assert.Len(t, params, 1)
	assert.Nil(t, headers)
	assert.NotNil(t, data)
}

// handleDomainConfig Tests

func TestHandleDomainConfig_ShouldReturnHeadersAndData_WhenValidConfiguration(t *testing.T) {
	config := CustomDomainConfiguration{
		Name:          "test-cert",
		CertificateId:  "",
		Certificate:   "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
		PrivateKey:    "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
	}

	headers, data, err := handleDomainConfig(config)

	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	assert.Contains(t, headers, HEADER_MD5_CAMEL)
}

func TestHandleDomainConfig_ShouldReturnError_WhenCertificateIdInvalidLength(t *testing.T) {
	// CERT_ID_SIZE is 32, so using a string of different length
	config := CustomDomainConfiguration{
		Name:          "test-cert",
		CertificateId:  "invalid-id", // Wrong length
		Certificate:   "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
		PrivateKey:    "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
	}

	headers, data, err := handleDomainConfig(config)

	assert.Error(t, err)
	assert.NotNil(t, headers) // Function returns empty map on error, not nil
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "length")
}

func TestHandleDomainConfig_ShouldReturnError_WhenNameTooShort(t *testing.T) {
	// MIN_CERTIFICATE_NAME_LENGTH is 3
	config := CustomDomainConfiguration{
		Name:          "ab", // Too short
		CertificateId: "",   // Empty is valid when Certificate is provided
		Certificate:   "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
		PrivateKey:    "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
	}

	headers, data, err := handleDomainConfig(config)

	assert.Error(t, err)
	assert.NotNil(t, headers) // Function returns empty map on error, not nil
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "length")
}

func TestHandleDomainConfig_ShouldReturnHeadersAndData_WhenCertificateIdIsSet(t *testing.T) {
	// CERT_ID_SIZE is 16, create a 16-char ID
	validId := ""
	for i := 0; i < 16; i++ {
		validId += "a"
	}
	config := CustomDomainConfiguration{
		Name:         "test-cert",
		CertificateId: validId,
	}

	headers, data, err := handleDomainConfig(config)

	assert.NoError(t, err)
	assert.NotNil(t, headers)
	assert.NotNil(t, data)
	assert.Contains(t, headers, HEADER_MD5_CAMEL)
}

func TestHandleDomainConfig_ShouldReturnError_WhenXmlBodyTooLarge(t *testing.T) {
	// Create a very large configuration that will exceed MAX_CERT_XML_BODY_SIZE
	largeCert := "-----BEGIN CERTIFICATE-----\n"
	for i := 0; i < 100000; i++ {
		largeCert += "a"
	}
	largeCert += "\n-----END CERTIFICATE-----"

	largeKey := "-----BEGIN PRIVATE KEY-----\n"
	for i := 0; i < 100000; i++ {
		largeKey += "a"
	}
	largeKey += "\n-----END PRIVATE KEY-----"

	config := CustomDomainConfiguration{
		Name:         "test-cert",
		CertificateId: "",
		Certificate:   largeCert,
		PrivateKey:    largeKey,
	}

	headers, data, err := handleDomainConfig(config)

	assert.Error(t, err)
	assert.NotNil(t, headers)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "length")
}
