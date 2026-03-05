// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with License.  You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for
// the specific language governing permissions and limitations under the License.

package obs

import (
	"strings"
)

func getV2StringToSign(method, canonicalizedURL string, headers map[string][]string, isObs bool) string {
	// Split path and query for special handling
	var path string
	var queryParams string
	if parts := strings.Split(canonicalizedURL, "?"); len(parts) > 1 {
		path = parts[0]
		queryParams = parts[1]
	} else {
		path = canonicalizedURL
	}

	// Build string to sign
	// When URL ends with ? (empty query), keep the ? in the string
	// When there are query params, include them (without ?)
	// When no query params, use empty string
	var fourthElement string
	if queryParams != "" {
		// Include both path and query params (without ?)
		fourthElement = path + "\n" + queryParams
	} else if strings.HasSuffix(canonicalizedURL, "?") {
		fourthElement = canonicalizedURL
	} else {
		fourthElement = path
	}

	stringToSign := strings.Join([]string{method, "\n", attachHeaders(headers, isObs), "\n", fourthElement}, "")

	var securityTokenFromQuery string
	// Check URL query for security token - this overrides headers for masking
	var query []string
	if queryParams != "" {
		query = strings.Split(queryParams, "&")
		for _, value := range query {
			tokenPrefix := ""
			if strings.HasPrefix(value, HEADER_STS_TOKEN_AMZ+"=") {
				tokenPrefix = HEADER_STS_TOKEN_AMZ
			} else if strings.HasPrefix(value, HEADER_STS_TOKEN_OBS+"=") {
				tokenPrefix = HEADER_STS_TOKEN_OBS
			}
			if tokenPrefix != "" {
				tokenValue := value[len(tokenPrefix)+1:]
				if tokenValue != "" {
					securityTokenFromQuery = tokenValue
				}
			}
		}
	}

	logStringToSign := stringToSign
	// Only mask token when it's in URL query string (securityTokenFromQuery is set)
	if len(securityTokenFromQuery) > 0 {
		logStringToSign = strings.ReplaceAll(logStringToSign, securityTokenFromQuery, "******")
	}
	doLog(LEVEL_DEBUG, "The v2 auth stringToSign:\n%s", logStringToSign)
	return logStringToSign
}

func v2Auth(ak, sk, method, canonicalizedURL string, headers map[string][]string, isObs bool) map[string]string {
	stringToSign := getV2StringToSign(method, canonicalizedURL, headers, isObs)
	return map[string]string{"Signature": Base64Encode(HmacSha1([]byte(sk), []byte(stringToSign)))}
}
