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

//go:build unit

package obs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecompressionRule_ShouldHaveCorrectFields(t *testing.T) {
	// 测试 DecompressionRule 结构的字段
	rule := DecompressionRule{
		ID:            "test-id",
		Project:       "test-project",
		Agency:        "test-agency",
		Events:        []string{"ObjectCreated:Put", "ObjectCreated:Post"},
		Prefix:        "test-prefix/",
		Suffix:        "test-suffix",
		Overwrite:     1,
		DecompressPath: "test-path",
		PolicyType:    "test-policy",
	}

	// 验证字段值
	assert.Equal(t, "test-id", rule.ID)
	assert.Equal(t, "test-project", rule.Project)
	assert.Equal(t, "test-agency", rule.Agency)
	assert.Equal(t, []string{"ObjectCreated:Put", "ObjectCreated:Post"}, rule.Events)
	assert.Equal(t, "test-prefix/", rule.Prefix)
	assert.Equal(t, "test-suffix", rule.Suffix)
	assert.Equal(t, 1, rule.Overwrite)
	assert.Equal(t, "test-path", rule.DecompressPath)
	assert.Equal(t, "test-policy", rule.PolicyType)
}

func TestSetBucketDecompressionInput_ShouldSerializeCorrectly(t *testing.T) {
	input := &SetBucketDecompressionInput{
		Bucket: "test-bucket",
		Rules: []DecompressionRule{
			{
				ID:      "rule1",
				Project: "project1",
				Agency:  "agency1",
				Events:  []string{"ObjectCreated:Put"},
				Prefix:  "files/",
			},
		},
	}

	// 测试序列化
	jsonBytes, err := json.Marshal(input)
	require.NoError(t, err)

	// 验证序列化结果 - 实际的序列化结果包含 Bucket 字段，且 decompressPath 为空字符串
	expected := `{"Bucket":"test-bucket","rules":[{"id":"rule1","project":"project1","agency":"agency1","events":["ObjectCreated:Put"],"prefix":"files/","suffix":"","overwrite":0,"policyType":""}]}`
	assert.JSONEq(t, expected, string(jsonBytes))

	// 测试 trans 方法
	params, headers, data, err := input.trans(true)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"obscompresspolicy": ""}, params)
	assert.Equal(t, map[string][]string{"Content-Type": {"application/json"}}, headers)
	assert.NotNil(t, data)
}

func TestSetBucketDecompressionInput_ShouldDeserializeCorrectly(t *testing.T) {
	jsonData := `{"rules":[{"id":"rule1","project":"project1","agency":"agency1","events":["ObjectCreated:Put"],"prefix":"files/","suffix":"","overwrite":0,"policyType":"","decompressPath":""}]}`

	var input SetBucketDecompressionInput
	err := json.Unmarshal([]byte(jsonData), &input)
	require.NoError(t, err)

	// 验证反序列化结果
	assert.Len(t, input.Rules, 1)
	assert.Equal(t, "rule1", input.Rules[0].ID)
	assert.Equal(t, "project1", input.Rules[0].Project)
	assert.Equal(t, "agency1", input.Rules[0].Agency)
	assert.Equal(t, []string{"ObjectCreated:Put"}, input.Rules[0].Events)
	assert.Equal(t, "files/", input.Rules[0].Prefix)
	assert.Equal(t, "", input.Rules[0].Suffix)
	assert.Equal(t, 0, input.Rules[0].Overwrite)
	assert.Equal(t, "", input.Rules[0].PolicyType)
	assert.Equal(t, "", input.Rules[0].DecompressPath)
}

func TestGetBucketDecompressionInput_ShouldHaveCorrectFields(t *testing.T) {
	// 测试 GetBucketDecompressionInput 结构的字段
	input := &GetBucketDecompressionInput{
		Bucket: "test-bucket",
	}

	// 验证字段值
	assert.Equal(t, "test-bucket", input.Bucket)
}

func TestDeleteBucketDecompressionInput_ShouldSerializeCorrectly(t *testing.T) {
	// 测试使用 newSubResourceSerial 创建的序列化器（用于 DELETE 操作）
	serializable := newSubResourceSerial(SubResourceDecompression)
	params, headers, data, err := serializable.trans(true)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, params)
	assert.Equal(t, map[string]string{"obscompresspolicy": ""}, params)
	assert.Nil(t, headers)
	assert.Nil(t, data)
}

func TestDeleteBucketDecompressionInput_ShouldHaveCorrectFields(t *testing.T) {
	// 测试 DeleteBucketDecompressionInput 结构的字段
	input := &DeleteBucketDecompressionInput{
		Bucket: "test-bucket",
	}

	// 验证字段值
	assert.Equal(t, "test-bucket", input.Bucket)
}