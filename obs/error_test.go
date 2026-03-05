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

	"github.com/stretchr/testify/assert"
)

// TestObsError_Error_ShouldReturnCorrectFormat_WhenCalled tests ObsError.Error()
func TestObsError_Error_ShouldReturnCorrectFormat_WhenCalled(t *testing.T) {
	err := ObsError{
		BaseModel: BaseModel{RequestId: "test-request-id"},
		Status:    "403",
		Code:      "AccessDenied",
		Message:   "Access Denied",
		Indicator: "",
	}

	result := err.Error()

	assert.Contains(t, result, "Status=403")
	assert.Contains(t, result, "Code=AccessDenied")
	assert.Contains(t, result, "Message=Access Denied")
	assert.Contains(t, result, "RequestId=test-request-id")
}

// TestObsError_Error_ShouldHandleEmptyFields_WhenCalled tests ObsError.Error() with empty fields
func TestObsError_Error_ShouldHandleEmptyFields_WhenCalled(t *testing.T) {
	err := ObsError{
		Status:    "",
		Code:      "",
		Message:   "",
		Indicator: "",
	}

	result := err.Error()

	assert.NotEmpty(t, result)
}
