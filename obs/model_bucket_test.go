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
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReplicationConfiguration_ShouldSerializeToXML_WhenValidConfig tests XML serialization
func TestReplicationConfiguration_ShouldSerializeToXML_WhenValidConfig(t *testing.T) {
	// Arrange
	config := ReplicationConfiguration{
		Rules: []ReplicationRule{
			{
				ID:     "rule-1",
				Status: RuleStatusEnabled,
				Prefix: &ReplicationPrefix{
					PrefixSet: PrefixSet{
						Prefixes: []string{"prefix1", "prefix2"},
					},
				},
				Destination: ReplicationDestination{
					Bucket:       "dest-bucket",
					StorageClass: StorageClassStandard,
					Location:     "region-cn-north-4",
				},
				HistoricalObjectReplication: "Enabled",
			},
		},
	}

	// Act
	data, err := xml.Marshal(config)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "<ReplicationConfiguration>")
	assert.Contains(t, string(data), "<Rule>")
	assert.Contains(t, string(data), "<Id>rule-1</Id>")
	assert.Contains(t, string(data), "<Status>Enabled</Status>")
	assert.Contains(t, string(data), "<Bucket>dest-bucket</Bucket>")
}

// TestReplicationConfiguration_ShouldDeserializeFromXML_WhenValidXML tests XML deserialization
func TestReplicationConfiguration_ShouldDeserializeFromXML_WhenValidXML(t *testing.T) {
	// Arrange
	xmlData := `<ReplicationConfiguration>
		<Rule>
			<Id>rule-1</Id>
			<Status>Enabled</Status>
			<Prefix>
				<PrefixSet>
					<Prefix>prefix1</Prefix>
					<Prefix>prefix2</Prefix>
				</PrefixSet>
			</Prefix>
			<Destination>
				<Bucket>dest-bucket</Bucket>
				<StorageClass>STANDARD</StorageClass>
				<Location>region-cn-north-4</Location>
			</Destination>
			<HistoricalObjectReplication>Enabled</HistoricalObjectReplication>
		</Rule>
	</ReplicationConfiguration>`

	// Act
	var config ReplicationConfiguration
	err := xml.Unmarshal([]byte(xmlData), &config)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, config.Rules, 1)
	assert.Equal(t, "rule-1", config.Rules[0].ID)
	assert.Equal(t, RuleStatusEnabled, config.Rules[0].Status)
	assert.NotNil(t, config.Rules[0].Prefix)
	assert.Len(t, config.Rules[0].Prefix.PrefixSet.Prefixes, 2)
	assert.Equal(t, "prefix1", config.Rules[0].Prefix.PrefixSet.Prefixes[0])
	assert.Equal(t, "dest-bucket", config.Rules[0].Destination.Bucket)
	assert.Equal(t, StorageClassStandard, config.Rules[0].Destination.StorageClass)
	assert.Equal(t, "region-cn-north-4", config.Rules[0].Destination.Location)
	assert.Equal(t, "Enabled", config.Rules[0].HistoricalObjectReplication)
}

// TestReplicationConfiguration_ShouldHandleMultipleRules_WhenMultipleRulesProvided tests multiple rules
func TestReplicationConfiguration_ShouldHandleMultipleRules_WhenMultipleRulesProvided(t *testing.T) {
	cases := []struct {
		name        string
		ruleCount   int
		expectError bool
	}{
		{"Single rule", 1, false},
		{"Two rules", 2, false},
		{"Max rules (100)", 100, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			rules := make([]ReplicationRule, tc.ruleCount)
			for i := 0; i < tc.ruleCount; i++ {
				rules[i] = ReplicationRule{
					ID:     "rule-" + string(rune(i)),
					Status: RuleStatusEnabled,
					Destination: ReplicationDestination{
						Bucket: "dest-bucket",
					},
				}
			}
			config := ReplicationConfiguration{Rules: rules}

			// Act
			data, err := xml.Marshal(config)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, data)
			}
		})
	}
}

// TestReplicationRule_ShouldHandleOptionalFields_WhenFieldsNotSet tests optional fields
func TestReplicationRule_ShouldHandleOptionalFields_WhenFieldsNotSet(t *testing.T) {
	// Arrange
	rule := ReplicationRule{
		Status: RuleStatusEnabled,
		Destination: ReplicationDestination{
			Bucket: "dest-bucket",
		},
	}

	// Act
	data, err := xml.Marshal(rule)

	// Assert
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "<Id>")
	assert.NotContains(t, string(data), "<Prefix>")
	assert.NotContains(t, string(data), "<HistoricalObjectReplication>")
	assert.Contains(t, string(data), "<Status>Enabled</Status>")
}

// TestDisPolicyConfiguration_ShouldSerializeToJSON_WhenValidConfig tests JSON serialization
func TestDisPolicyConfiguration_ShouldSerializeToJSON_WhenValidConfig(t *testing.T) {
	// Arrange
	config := DisPolicyConfiguration{
		Rules: []DisPolicyRule{
			{
				ID:      "rule-1",
				Stream:  "test-stream",
				Project: "test-project-id",
				Events:  []string{"ObjectCreated:*", "ObjectRemoved:*"},
				Prefix:  "images/",
				Suffix:  ".jpg",
				Agency:  "test-agency",
			},
		},
	}

	// Act
	data, err := json.Marshal(config)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), `"rules"`)
	assert.Contains(t, string(data), `"id":"rule-1"`)
	assert.Contains(t, string(data), `"stream":"test-stream"`)
	assert.Contains(t, string(data), `"events":["ObjectCreated:*","ObjectRemoved:*"]`)
}

// TestDisPolicyConfiguration_ShouldDeserializeFromJSON_WhenValidJSON tests JSON deserialization
func TestDisPolicyConfiguration_ShouldDeserializeFromJSON_WhenValidJSON(t *testing.T) {
	// Arrange
	jsonData := `{
		"rules": [{
			"id": "rule-1",
			"stream": "test-stream",
			"project": "test-project-id",
			"events": ["ObjectCreated:*", "ObjectRemoved:*"],
			"prefix": "images/",
			"suffix": ".jpg",
			"agency": "test-agency"
		}]
	}`

	// Act
	var config DisPolicyConfiguration
	err := json.Unmarshal([]byte(jsonData), &config)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, config.Rules, 1)
	assert.Equal(t, "rule-1", config.Rules[0].ID)
	assert.Equal(t, "test-stream", config.Rules[0].Stream)
	assert.Equal(t, "test-project-id", config.Rules[0].Project)
	assert.Contains(t, config.Rules[0].Events, "ObjectCreated:*")
	assert.Contains(t, config.Rules[0].Events, "ObjectRemoved:*")
	assert.Equal(t, "images/", config.Rules[0].Prefix)
	assert.Equal(t, ".jpg", config.Rules[0].Suffix)
	assert.Equal(t, "test-agency", config.Rules[0].Agency)
}

// TestDisPolicyConfiguration_ShouldHandleMultipleRules_WhenMultipleRulesProvided tests multiple rules
func TestDisPolicyConfiguration_ShouldHandleMultipleRules_WhenMultipleRulesProvided(t *testing.T) {
	cases := []struct {
		name        string
		ruleCount   int
		expectError bool
	}{
		{"Single rule", 1, false},
		{"Two rules", 2, false},
		{"Max rules (10)", 10, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			rules := make([]DisPolicyRule, tc.ruleCount)
			for i := 0; i < tc.ruleCount; i++ {
				rules[i] = DisPolicyRule{
					ID:      "rule-" + string(rune('1'+i)),
					Stream:  "test-stream",
					Project: "test-project-id",
					Events:  []string{"ObjectCreated:*"},
					Agency:  "test-agency",
				}
			}
			config := DisPolicyConfiguration{Rules: rules}

			// Act
			data, err := json.Marshal(config)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, data)
			}
		})
	}
}

// TestReplicationDestination_ShouldHandleOptionalFields_WhenOnlyRequiredFieldSet tests optional destination fields
func TestReplicationDestination_ShouldHandleOptionalFields_WhenOnlyRequiredFieldSet(t *testing.T) {
	// Arrange
	dest := ReplicationDestination{
		Bucket: "dest-bucket",
	}

	// Act
	data, err := xml.Marshal(dest)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(data), "<Bucket>dest-bucket</Bucket>")
	assert.NotContains(t, string(data), "<StorageClass>")
	assert.NotContains(t, string(data), "<Location>")
}

// TestReplicationPrefix_ShouldSerializeCorrectly_WhenValidPrefix tests prefix serialization
func TestReplicationPrefix_ShouldSerializeCorrectly_WhenValidPrefix(t *testing.T) {
	// Arrange
	prefix := ReplicationPrefix{
		PrefixSet: PrefixSet{
			Prefixes: []string{"images/", "videos/"},
		},
	}

	// Act
	data, err := xml.Marshal(prefix)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}
