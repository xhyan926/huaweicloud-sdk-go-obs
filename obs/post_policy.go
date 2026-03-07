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
	"errors"
	"time"
)

// buildPostPolicyExpiration builds expiration time string for POST policy (internal function)
func buildPostPolicyExpiration(expiresIn int64) string {
	expirationTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
	return expirationTime.UTC().Format("2006-01-02T15:04:05.000Z")
}

// validatePostPolicy validates POST policy structure (internal function)
func validatePostPolicy(bucket, key string) error {
	if bucket == "" {
		return errors.New("bucket is required")
	}
	if key == "" {
		return errors.New("key is required")
	}
	return nil
}
