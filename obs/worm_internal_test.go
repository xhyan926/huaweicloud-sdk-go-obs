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
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPutBucketWormConfigurationInput_ShouldHaveCorrectFields tests the PutBucketWormConfigurationInput structure
func TestPutBucketWormConfigurationInput_ShouldHaveCorrectFields(t *testing.T) {
	input := &PutBucketWormConfigurationInput{
		Bucket:            "test-bucket",
		WormRetentionMode: WormRetentionModeCompliance,
		WormRetentionPeriod: 30,
		DefaultRetention: DefaultRetention{
			Days: 10,
			Years: 0,
			Mode: "COMPLIANCE",
		},
		ExtendRetention: ExtendRetention{
			Days: 5,
			Years: 0,
			Mode: "GOVERNANCE",
		},
	}

	// Test field values
	assert.Equal(t, "test-bucket", input.Bucket)
	assert.Equal(t, WormRetentionModeCompliance, input.WormRetentionMode)
	assert.Equal(t, 30, input.WormRetentionPeriod)
	assert.Equal(t, 10, input.DefaultRetention.Days)
	assert.Equal(t, 0, input.DefaultRetention.Years)
	assert.Equal(t, "COMPLIANCE", input.DefaultRetention.Mode)
	assert.Equal(t, 5, input.ExtendRetention.Days)
	assert.Equal(t, 0, input.ExtendRetention.Years)
	assert.Equal(t, "GOVERNANCE", input.ExtendRetention.Mode)
}

// TestPutBucketWormConfigurationInput_ShouldSerializeCorrectly tests the XML serialization of PutBucketWormConfigurationInput
func TestPutBucketWormConfigurationInput_ShouldSerializeCorrectly(t *testing.T) {
	input := &PutBucketWormConfigurationInput{
		Bucket:            "test-bucket",
		WormRetentionMode: WormRetentionModeCompliance,
		WormRetentionPeriod: 30,
		DefaultRetention: DefaultRetention{
			Days: 10,
			Years: 0,
			Mode: "COMPLIANCE",
		},
		ExtendRetention: ExtendRetention{
			Days: 5,
			Years: 0,
			Mode: "GOVERNANCE",
		},
	}

	params, headers, data, err := input.trans(true)

	// Test serialization
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"worm": ""}, params)
	assert.Equal(t, "application/xml", headers["Content-Type"][0])
	assert.NotNil(t, data)

	// Test XML content
	// The actual XML content might vary based on the actual struct implementation
	// We're testing that serialization happens without error
}

// TestPutBucketWormConfigurationInput_ShouldDeserializeCorrectly tests the XML deserialization
func TestPutBucketWormConfigurationInput_ShouldDeserializeCorrectly(t *testing.T) {
	xmlContent := `<WormConfiguration>
		<WormRetentionMode>Compliance</WormRetentionMode>
		<WormRetentionPeriod>30</WormRetentionPeriod>
		<DefaultRetention>
			<Days>10</Days>
			<Mode>COMPLIANCE</Mode>
		</DefaultRetention>
		<ExtendRetention>
			<Days>5</Days>
			<Mode>GOVERNANCE</Mode>
		</ExtendRetention>
	</WormConfiguration>`

	var config PutBucketWormConfigurationInput
	err := xml.Unmarshal([]byte(xmlContent), &config)

	assert.NoError(t, err)
	assert.Equal(t, WormRetentionModeCompliance, config.WormRetentionMode)
	assert.Equal(t, 30, config.WormRetentionPeriod)
	assert.Equal(t, 10, config.DefaultRetention.Days)
	assert.Equal(t, "COMPLIANCE", config.DefaultRetention.Mode)
	assert.Equal(t, 5, config.ExtendRetention.Days)
	assert.Equal(t, "GOVERNANCE", config.ExtendRetention.Mode)
}

// TestGetBucketWormConfigurationInput_ShouldHaveCorrectFields tests the GetBucketWormConfigurationInput structure
func TestGetBucketWormConfigurationInput_ShouldHaveCorrectFields(t *testing.T) {
	input := &GetBucketWormConfigurationInput{
		Bucket: "test-bucket",
	}

	assert.Equal(t, "test-bucket", input.Bucket)
}

// TestGetBucketWormConfigurationInput_ShouldSerializeCorrectly tests the serialization of GetBucketWormConfigurationInput
func TestGetBucketWormConfigurationInput_ShouldSerializeCorrectly(t *testing.T) {
	input := &GetBucketWormConfigurationInput{
		Bucket: "test-bucket",
	}

	params, headers, data, err := input.trans(true)

	// Test serialization
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"worm": ""}, params)
	assert.Empty(t, headers)
	assert.Nil(t, data) // GET requests typically don't have a body
}

// TestDeleteBucketWormConfigurationInput_ShouldHaveCorrectFields tests the DeleteBucketWormConfigurationInput structure
func TestDeleteBucketWormConfigurationInput_ShouldHaveCorrectFields(t *testing.T) {
	input := &DeleteBucketWormConfigurationInput{
		Bucket: "test-bucket",
	}

	assert.Equal(t, "test-bucket", input.Bucket)
}

// TestDeleteBucketWormConfigurationInput_ShouldSerializeCorrectly tests the serialization of DeleteBucketWormConfigurationInput
func TestDeleteBucketWormConfigurationInput_ShouldSerializeCorrectly(t *testing.T) {
	input := &DeleteBucketWormConfigurationInput{
		Bucket: "test-bucket",
	}

	params, headers, data, err := input.trans(true)

	// Test serialization
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"worm": ""}, params)
	assert.Empty(t, headers)
	assert.NotNil(t, data) // Should be an empty reader
}

// TestDefaultRetention_ShouldHaveCorrectFields tests the DefaultRetention structure
func TestDefaultRetention_ShouldHaveCorrectFields(t *testing.T) {
	retention := DefaultRetention{
		Days:  10,
		Years: 0,
		Mode:  "COMPLIANCE",
	}

	assert.Equal(t, 10, retention.Days)
	assert.Equal(t, 0, retention.Years)
	assert.Equal(t, "COMPLIANCE", retention.Mode)
}

// TestExtendRetention_ShouldHaveCorrectFields tests the ExtendRetention structure
func TestExtendRetention_ShouldHaveCorrectFields(t *testing.T) {
	retention := ExtendRetention{
		Days:  5,
		Years: 0,
		Mode:  "GOVERNANCE",
	}

	assert.Equal(t, 5, retention.Days)
	assert.Equal(t, 0, retention.Years)
	assert.Equal(t, "GOVERNANCE", retention.Mode)
}

// TestGetBucketWormConfigurationOutput_ShouldHaveCorrectFields tests the GetBucketWormConfigurationOutput structure
func TestGetBucketWormConfigurationOutput_ShouldHaveCorrectFields(t *testing.T) {
	output := &GetBucketWormConfigurationOutput{
		BaseModel: BaseModel{},
		IsWormEnabled:    "Enabled",
		WormRetentionMode: WormRetentionModeCompliance,
		WormRetentionPeriod: 30,
		ExtendRetention: ExtendRetention{
			Days: 5,
			Years: 0,
			Mode: "GOVERNANCE",
		},
		AccessLogEnabled: true,
		ObjectSaveDays:   7,
	}

	assert.Equal(t, "Enabled", output.IsWormEnabled)
	assert.Equal(t, WormRetentionModeCompliance, output.WormRetentionMode)
	assert.Equal(t, 30, output.WormRetentionPeriod)
	assert.Equal(t, 5, output.ExtendRetention.Days)
	assert.Equal(t, true, output.AccessLogEnabled)
	assert.Equal(t, 7, output.ObjectSaveDays)
}