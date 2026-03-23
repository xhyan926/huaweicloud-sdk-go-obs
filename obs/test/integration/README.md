# 集成测试说明

本目录包含华为云 OBS Go SDK 的集成测试。

## 配置方式

集成测试支持两种配置方式：环境变量和配置文件。

### 方式一：环境变量配置

#### 必需的环境变量

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `OBS_TEST_AK` | 华为云访问密钥 ID | `XXXXXXXXXXXXXXXXXXXX` |
| `OBS_TEST_SK` | 华为云访问密钥 | `YYYYYYYYYYYYYYYYYYYY` |
| `OBS_TEST_ENDPOINT` | OBS 服务终端节点 | `https://obs.cn-north-4.myhuaweicloud.com` |
| `OBS_TEST_BUCKET` | 测试用的桶名称 | `test-bucket` |

#### 可选的环境变量

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `OBS_TEST_REGION` | 测试区域 | `cn-north-4` |
| `OBS_TEST_USE_TEMP_BUCKET` | 是否使用临时桶模式 | `true` 或 `false` |

#### DIS 事件通知测试的可选环境变量

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `OBS_TEST_DIS_STREAM` | DIS 通道名称 | `test-dis-stream` |
| `OBS_TEST_DIS_PROJECT` | DIS 项目 ID | `test-project-id` |
| `OBS_TEST_DIS_AGENCY` | DIS IAM 委托名称 | `test-dis-agency` |

#### 跨区域复制测试的可选环境变量

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `OBS_TEST_REPLICATION_DEST_BUCKET` | 复制目标桶名称 | `test-dest-bucket` |
| `OBS_TEST_REPLICATION_LOCATION` | 复制目标区域 | `cn-south-1` |

### 方式二：配置文件（推荐）

配置文件用于管理**特性相关配置**（如 DIS 通道、复制目标等），**不包含敏感认证信息**。

**注意**：认证信息（ak、sk、endpoint、bucket）必须通过环境变量设置。

#### 1. 设置环境变量（必需）

```bash
export OBS_TEST_AK=your_access_key
export OBS_TEST_SK=your_secret_key
export OBS_TEST_ENDPOINT=your_endpoint
export OBS_TEST_BUCKET=your_test_bucket
```

#### 2. 创建配置文件

复制示例配置文件并修改：

```bash
cp test.config.json.example test.config.json
```

#### 3. 编辑配置文件

配置文件仅包含特性相关配置：

```json
{
  "dis": {
    "stream": "your-dis-stream-name",
    "project": "your-dis-project-id",
    "agency": "your-iam-agency-name"
  },
  "replication": {
    "destBucket": "your-dest-bucket-name",
    "location": "cn-south-1"
  }
}
```

#### 3. 指定配置文件路径

```bash
export OBS_TEST_CONFIG_FILE=test.config.json
```

#### 配置文件搜索路径

如果未设置 `OBS_TEST_CONFIG_FILE`，测试将按以下顺序搜索配置文件：
1. `test.config.json`
2. `fixtures/test.config.json`
3. `../../test.config.json`

### 配置优先级

环境变量的优先级高于配置文件。可以使用环境变量覆盖配置文件中的设置。

## 运行集成测试

### 使用环境变量运行

```bash
export OBS_TEST_AK=your_access_key
export OBS_TEST_SK=your_secret_key
export OBS_TEST_ENDPOINT=your_endpoint
export OBS_TEST_BUCKET=your_test_bucket

go test ./obs/test/integration/ -v -tags=integration
```

### 使用配置文件运行

```bash
# 1. 设置认证环境变量（必需）
export OBS_TEST_AK=your_access_key
export OBS_TEST_SK=your_secret_key
export OBS_TEST_ENDPOINT=your_endpoint
export OBS_TEST_BUCKET=your_test_bucket

# 2. 可选：指定配置文件路径（或使用默认搜索路径）
export OBS_TEST_CONFIG_FILE=test.config.json

# 3. 运行测试
go test ./obs/test/integration/ -v -tags=integration
```

### 运行特定功能的集成测试

```bash
# 只运行跨区域复制集成测试
go test ./obs/test/integration/ -v -tags=integration -run Replication

# 只运行 DIS 事件通知集成测试
go test ./obs/test/integration/ -v -tags=integration -run DisPolicy
```

## 测试覆盖

### 跨区域复制集成测试 (replication_integration_test.go)

- `TestIntegration_SetBucketReplication_ShouldSucceed` - 测试设置跨区域复制配置
- `TestIntegration_GetBucketReplication_ShouldReturnConfig` - 测试获取跨区域复制配置
- `TestIntegration_DeleteBucketReplication_ShouldSucceed` - 测试删除跨区域复制配置
- `TestIntegration_Replication_ShouldSupportOnlyOBSSignature` - 验证仅支持 OBS 签名
- `TestIntegration_Replication_ShouldHandleMultipleRules` - 测试多条规则处理

### DIS 事件通知集成测试 (dis_policy_integration_test.go)

- `TestIntegration_SetBucketDisPolicy_ShouldSucceed` - 测试设置 DIS 事件通知策略
- `TestIntegration_GetBucketDisPolicy_ShouldReturnConfig` - 测试获取 DIS 事件通知配置
- `TestIntegration_DeleteBucketDisPolicy_ShouldSucceed` - 测试删除 DIS 事件通知策略
- `TestIntegration_DisPolicy_ShouldHandleMultipleRules` - 测试多条规则处理
- `TestIntegration_DisPolicy_ShouldSupportOnlyOBSSignature` - 验证仅支持 OBS 签名

## 注意事项

1. **签名类型**：跨区域复制和 DIS 事件通知功能仅支持 OBS 签名类型，不支持 AWS 签名
2. **桶权限**：测试用的桶需要已存在且具有相应的权限
3. **DIS 资源**：运行 DIS 事件通知测试前，需要先创建 DIS 通道和 IAM 委托
4. **临时桶模式**：在配置文件中设置 `useTempBucket: true` 可以自动创建和删除临时测试桶
5. **数据清理**：测试会在完成后自动清理配置数据
6. **默认值**：如果配置文件中某些字段为空，测试会使用默认值或跳过相应测试

## 测试辅助函数 (helper.go)

| 函数 | 说明 |
|------|------|
| `createClient(t, signature)` | 创建 OBS 客户端 |
| `getTestBucket(t)` | 获取测试桶名 |
| `getDisConfig(t)` | 获取 DIS 配置（从配置文件或环境变量） |
| `getReplicationConfig(t)` | 获取跨区域复制配置（从配置文件或环境变量） |
| `setupTestBucket(t, client)` | 设置测试桶（支持临时桶模式） |
| `cleanupReplication(t, client, bucket)` | 清理跨区域复制配置 |
| `cleanupDisPolicy(t, client, bucket)` | 清理 DIS 事件通知策略 |
| `generateTestID(prefix)` | 生成测试用例唯一 ID |

## 故障排除

### 测试跳过

如果环境变量未设置，测试会自动跳过：

```
--- SKIP: TestIntegration_SetBucketReplication_ShouldSucceed (0.00s)
    helper.go:XX: 跳过集成测试：未设置必需的环境变量 (OBS_TEST_AK, OBS_TEST_SK, OBS_TEST_ENDPOINT)
```

### 常见错误

| 错误 | 原因 | 解决方法 |
|------|------|---------|
| `OBS credentials not set` | 环境变量未设置 | 设置必需的环境变量或使用配置文件 |
| `bucket not found` | 测试桶不存在 | 创建测试桶或设置 `useTempBucket: true` |
| `permission denied` | 权限不足 | 检查密钥权限或 IAM 策略 |
| `InvalidAccessKeyId` | 密钥无效 | 检查密钥是否正确 |
| `配置文件解析失败` | JSON 格式错误 | 检查配置文件格式是否正确 |

### 配置文件调试

启用调试日志查看配置加载情况：

```bash
go test ./obs/test/integration/ -v -tags=integration -args -test.v
```

测试会输出类似以下日志：
```
=== RUN   TestIntegration_SetBucketDisPolicy_ShouldSucceed
    helper.go:XX: 已加载配置文件: test.config.json
    helper.go:XX: 使用默认 DIS Stream: test-dis-stream
```
