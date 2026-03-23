//go:build integration
// +build integration

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

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

// disConfig DIS 事件通知配置
type disConfig struct {
	Stream  string `json:"stream"`
	Project string `json:"project"`
	Agency  string `json:"agency"`
}

// replicationConfig 跨区域复制配置
type replicationConfig struct {
	DestBucket string `json:"destBucket"`
	Location   string `json:"location"`
}

// testConfigFile 配置文件结构（仅包含特性配置，不含敏感信息）
type testConfigFile struct {
	DIS        disConfig        `json:"dis"`
	Replication replicationConfig `json:"replication"`
}

// testConfig 集成测试配置（运行时）
type testConfig struct {
	ak            string
	sk            string
	endpoint      string
	bucket        string
	region        string
	useTempBucket bool
	dis           disConfig
	replication   replicationConfig
}

// loadTestConfig 加载测试配置，如果环境变量未设置则跳过测试
func loadTestConfig(t *testing.T) *testConfig {
	ak := os.Getenv("OBS_TEST_AK")
	sk := os.Getenv("OBS_TEST_SK")
	endpoint := os.Getenv("OBS_TEST_ENDPOINT")
	bucket := os.Getenv("OBS_TEST_BUCKET")
	region := os.Getenv("OBS_TEST_REGION")

	if ak == "" || sk == "" || endpoint == "" {
		t.Skip("跳过集成测试：未设置必需的环境变量 (OBS_TEST_AK, OBS_TEST_SK, OBS_TEST_ENDPOINT)")
	}

	// 尝试加载配置文件
	cfg := &testConfig{
		ak:            ak,
		sk:            sk,
		endpoint:      endpoint,
		bucket:        bucket,
		region:        region,
		useTempBucket: os.Getenv("OBS_TEST_USE_TEMP_BUCKET") == "true",
	}

	// 加载配置文件（如果存在）
	if configFile := os.Getenv("OBS_TEST_CONFIG_FILE"); configFile != "" {
		loadConfigFile(cfg, configFile, t)
	} else {
		// 尝试默认配置文件路径
		defaultPaths := []string{
			"test.config.json",
			"fixtures/test.config.json",
			filepath.Join("..", "..", "test.config.json"),
		}
		for _, path := range defaultPaths {
			if _, err := os.Stat(path); err == nil {
				loadConfigFile(cfg, path, t)
				break
			}
		}
	}

	// 环境变量优先级高于配置文件
	if v := os.Getenv("OBS_TEST_DIS_STREAM"); v != "" {
		cfg.dis.Stream = v
	}
	if v := os.Getenv("OBS_TEST_DIS_PROJECT"); v != "" {
		cfg.dis.Project = v
	}
	if v := os.Getenv("OBS_TEST_DIS_AGENCY"); v != "" {
		cfg.dis.Agency = v
	}
	if v := os.Getenv("OBS_TEST_REPLICATION_DEST_BUCKET"); v != "" {
		cfg.replication.DestBucket = v
	}
	if v := os.Getenv("OBS_TEST_REPLICATION_LOCATION"); v != "" {
		cfg.replication.Location = v
	}

	return cfg
}

// loadConfigFile 从文件加载配置（仅特性配置，不含认证信息）
func loadConfigFile(cfg *testConfig, filePath string, t *testing.T) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Logf("无法读取配置文件 %s: %v (将使用默认值或环境变量)", filePath, err)
		return
	}

	var config testConfigFile
	if err := json.Unmarshal(data, &config); err != nil {
		t.Logf("解析配置文件失败 %s: %v (将使用默认值或环境变量)", filePath, err)
		return
	}

	// 只填充空值，不覆盖已设置的值（环境变量优先）
	if cfg.dis.Stream == "" && config.DIS.Stream != "" {
		cfg.dis.Stream = config.DIS.Stream
	}
	if cfg.dis.Project == "" && config.DIS.Project != "" {
		cfg.dis.Project = config.DIS.Project
	}
	if cfg.dis.Agency == "" && config.DIS.Agency != "" {
		cfg.dis.Agency = config.DIS.Agency
	}
	if cfg.replication.DestBucket == "" && config.Replication.DestBucket != "" {
		cfg.replication.DestBucket = config.Replication.DestBucket
	}
	if cfg.replication.Location == "" && config.Replication.Location != "" {
		cfg.replication.Location = config.Replication.Location
	}

	t.Logf("已加载配置文件: %s", filePath)
}

// createClient 创建 OBS 客户端
func createClient(t *testing.T, signature obs.SignatureType) *obs.ObsClient {
	config := loadTestConfig(t)

	var client *obs.ObsClient
	var err error

	if config.region != "" {
		client, err = obs.New(config.ak, config.sk, config.endpoint, obs.WithSignature(signature), obs.WithRegion(config.region))
	} else {
		client, err = obs.New(config.ak, config.sk, config.endpoint, obs.WithSignature(signature))
	}

	if err != nil {
		t.Fatalf("创建 OBS 客户端失败: %v", err)
	}

	return client
}

// getTestBucket 获取测试桶名
func getTestBucket(t *testing.T) string {
	config := loadTestConfig(t)
	return config.bucket
}

// createTempBucket 创建临时桶，返回桶名
func createTempBucket(t *testing.T, client *obs.ObsClient) string {
	config := loadTestConfig(t)
	bucketName := fmt.Sprintf("test-temp-bucket-%d", os.Getpid())

	input := &obs.CreateBucketInput{
		Bucket: bucketName,
	}

	if config.region != "" {
		input.BucketLocation = obs.BucketLocation{Location: config.region}
	}

	_, err := client.CreateBucket(input)
	if err != nil {
		t.Fatalf("创建临时桶失败: %v", err)
	}

	return bucketName
}

// deleteTempBucket 删除临时桶
func deleteTempBucket(t *testing.T, client *obs.ObsClient, bucketName string) {
	// 先清空桶中的对象
	input := &obs.ListObjectsInput{
		Bucket: bucketName,
	}

	for {
		output, err := client.ListObjects(input)
		if err != nil {
			t.Logf("列出对象失败: %v", err)
			break
		}

		if len(output.Contents) == 0 {
			break
		}

		// 删除所有对象
		for _, obj := range output.Contents {
			_, _ = client.DeleteObject(&obs.DeleteObjectInput{Bucket: bucketName, Key: obj.Key})
		}

		if !output.IsTruncated {
			break
		}
		input.Marker = output.NextMarker
	}

	// 删除桶
	_, err := client.DeleteBucket(bucketName)
	if err != nil {
		t.Logf("删除临时桶失败: %v", err)
	}
}

// useTempBucket 是否使用临时桶模式
func useTempBucket(t *testing.T) bool {
	config := loadTestConfig(t)
	return config.useTempBucket
}

// setupTestBucket 设置测试桶，根据模式返回桶名和清理函数
func setupTestBucket(t *testing.T, client *obs.ObsClient) (bucketName string) {
	if useTempBucket(t) {
		bucketName = createTempBucket(t, client)
		t.Cleanup(func() {
			deleteTempBucket(t, client, bucketName)
		})
	} else {
		bucketName = getTestBucket(t)
		if bucketName == "" {
			t.Skip("跳过测试：未设置 OBS_TEST_BUCKET 环境变量且未启用临时桶模式")
		}
	}
	return
}

// generateTestID 生成测试用例唯一 ID
func generateTestID(prefix string) string {
	pid := os.Getpid()
	return fmt.Sprintf("%s-%d", prefix, pid)
}

// cleanupReplication 清理跨区域复制配置
func cleanupReplication(t *testing.T, client *obs.ObsClient, bucket string) {
	_, err := client.DeleteBucketReplication(bucket)
	if err != nil {
		t.Logf("清理跨区域复制配置失败: %v", err)
	}
}

// cleanupDisPolicy 清理 DIS 事件通知策略
func cleanupDisPolicy(t *testing.T, client *obs.ObsClient, bucket string) {
	_, err := client.DeleteBucketDisPolicy(bucket)
	if err != nil {
		t.Logf("清理 DIS 事件通知策略失败: %v", err)
	}
}

// skipIfNotOBS 如果不是 OBS 签名则跳过测试
func skipIfNotOBS(t *testing.T, signature obs.SignatureType) {
	if signature != obs.SignatureObs {
		t.Skipf("跳过测试：%s 签名不支持此功能", signature)
	}
}

// waitForBucketDeleted 等待桶删除完成
func waitForBucketDeleted(t *testing.T, client *obs.ObsClient, bucketName string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		_, err := client.HeadBucket(bucketName)
		if err != nil {
			// 桶不存在，说明已删除
			return
		}
		time.Sleep(2 * time.Second)
	}
}

// getTestTimestamp 获取测试时间戳（用于生成唯一资源名）
func getTestTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// getDisConfig 获取 DIS 配置，如果配置缺失则使用默认值
func getDisConfig(t *testing.T) disConfig {
	config := loadTestConfig(t)
	disCfg := config.dis

	// 设置默认值
	if disCfg.Stream == "" {
		disCfg.Stream = "test-dis-stream"
		t.Logf("使用默认 DIS Stream: %s", disCfg.Stream)
	}
	if disCfg.Project == "" {
		disCfg.Project = "test-project-id"
		t.Logf("使用默认 DIS Project: %s", disCfg.Project)
	}
	if disCfg.Agency == "" {
		disCfg.Agency = "test-dis-agency"
		t.Logf("使用默认 DIS Agency: %s", disCfg.Agency)
	}

	return disCfg
}

// getReplicationConfig 获取跨区域复制配置，如果配置缺失则使用默认值
func getReplicationConfig(t *testing.T) replicationConfig {
	config := loadTestConfig(t)
	replCfg := config.replication

	// 如果未配置目标桶，使用源桶作为目标（测试用）
	if replCfg.DestBucket == "" {
		replCfg.DestBucket = config.bucket
		t.Logf("使用源桶作为目标桶进行测试: %s", replCfg.DestBucket)
	}
	if replCfg.Location == "" {
		replCfg.Location = config.region
		if replCfg.Location == "" {
			replCfg.Location = "cn-north-4"
		}
		t.Logf("使用默认复制区域: %s", replCfg.Location)
	}

	return replCfg
}

// getTestDestBucket 获取跨区域复制的目标桶名
func getTestDestBucket(t *testing.T) string {
	return getReplicationConfig(t).DestBucket
}

// getTestReplicationLocation 获取跨区域复制目标区域
func getTestReplicationLocation(t *testing.T) string {
	return getReplicationConfig(t).Location
}
