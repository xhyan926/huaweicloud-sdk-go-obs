# go-sdk-integration

## 技能概述

OBS SDK Go集成测试编写指南。本技能指导用户如何为华为云OBS SDK编写高质量的集成测试，涵盖从创建集成测试框架到编写具体测试用例的完整流程。

## 使用场景

- 在子任务中编写集成测试
- 对存量代码补充集成测试
- 验证核心功能端到端流程
- 集成测试环境配置和资源管理

## 技能调用接口

本技能可以被go-sdk-dev-task技能调用，以实现开发任务的自动化测试生成。

### 调用方式

go-sdk-dev-task可以通过以下方式调用本技能：

```bash
# 基本调用（使用默认参数）
/go-sdk-dev-task --type=feature

# 指定模块列表
/go-sdk-dev-task --test-modules=integration

# 带具体参数调用
/go-sdk-dev-task --integration-modules=bucket,object,auth --test-name=MyFeatureIntegrationTest
```

### 技能协调机制

1. **去重检查**
   - 检查目标函数是否已有集成测试
   - 避免重复生成相同测试用例
   - 采用BDD命名规范

2. **策略一致性**
   - 根据任务类型选择合适的测试策略
   - 集成测试功能完整性优先
   - 考虑测试覆盖率和效率平衡

3. **资源管理**
   - 使用统一的测试数据目录
   - 避免测试资源冲突
   - 确保测试环境清洁

4. **进度同步**
   - 向go-sdk-dev-task报告测试生成进度
   - 协调测试执行时机
   - 汇总测试结果

### 技能调用建议

go-sdk-dev-task在以下情况会调用本技能：

1. **新功能开发**: 自动生成集成测试用例
2. **Bug修复**: 生成回归集成测试
3. **代码重构**: 确保无功能回归
4. **性能优化**: 生成性能对比测试
5. **安全审查**: 重点关注安全相关的集成测试

## 核心功能

### 1. 集成测试文件结构

生成符合规范的集成测试文件，使用build tag `integration`确保与单元测试隔离。

### 2. 测试客户端使用

指导如何创建和使用IntegrationClient，包括：
- 自动跳过机制（环境变量检查）
- Mock服务器支持
- 自动清理资源
- 上下文管理

### 3. 测试场景设计

提供核心功能的集成测试场景：
- 认证测试（静态凭证、临时凭证）
- 上传下载测试（小文件、大文件、分块上传）
- 断点续传测试
- 存储桶管理测试

### 4. 测试数据管理

测试对象的创建、清理和管理规范：
- 唯一键生成
- 测试数据准备
- 资源清理策略

### 5. 资源清理规范

确保测试后资源正确清理，避免测试环境污染：
- 注册清理函数
- 自动执行清理
- 错误处理和日志记录

## 使用方法

### 基本用法

```bash
/go-sdk-integration
```

### 带参数使用

```bash
/go-sdk-integration --modules=bucket,object,auth
/go-sdk-integration --name=MyIntegrationTest
/go-sdk-integration --output=./test/integration/
```

### 技能输出

1. **测试文件结构**：
   ```
   obs/test/integration/e2e/
   ├── [name]_integration_test.go
   └── README.md
   ```

2. **测试客户端代码**：
   - IntegrationClient使用示例
   - Mock服务器配置
   - 清理函数注册

3. **测试场景模板**：
   - 认证测试模板
   - 上传下载测试模板
   - 分块上传测试模板

4. **配置和规范**：
   - 环境变量说明
   - 测试配置管理
   - 最佳实践指南

## 输出示例

### 生成的测试文件示例

```go
//go:build integration

package e2e

import (
	"testing"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// TestAuthentication_ShouldConnectSuccessfully_GivenValidCredentials
func TestAuthentication_ShouldConnectSuccessfully_GivenValidCredentials(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	// 验证连接是否正常
	input := &obs.GetBucketLocationInput{
		Bucket: client.GetTestBucket(),
	}

	_, err := client.TestClient().GetBucketLocation(input)
	if err != nil {
		t.Fatalf("认证失败: %v", err)
	}
}

// TestUploadDownload_ShouldUploadAndDownloadSuccessfully_GivenSmallFile
func TestUploadDownload_ShouldUploadAndDownloadSuccessfully_GivenSmallFile(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()
	objectKey := client.GetTestObjectKey("small-file.txt")
	content := "This is a small test file."

	// 上传对象
	err := client.TestClient().PutObject(bucket, objectKey, content, nil)
	if err != nil {
		t.Fatalf("上传对象失败: %v", err)
	}

	// 清理函数
	client.AddCleanup(func(t *testing.T) {
		err := client.TestClient().DeleteObject(bucket, objectKey)
		if err != nil {
			t.Logf("清理对象失败: %v", err)
		}
	})

	// 下载对象
	input := &obs.GetObjectInput{
		Bucket: bucket,
		Key:    objectKey,
	}

	output, err := client.TestClient().GetObject(input)
	if err != nil {
		t.Fatalf("下载对象失败: %v", err)
	}
	defer output.Body.Close()

	// 验证内容
	// ...
}
```

### 环境变量配置示例

```bash
# 必需配置
export OBS_TEST_AK="your-access-key"
export OBS_TEST_SK="your-secret-key"
export OBS_TEST_ENDPOINT="https://obs.cn-north-4.myhuaweicloud.com"
export OBS_TEST_BUCKET="your-test-bucket"

# 可选配置
export OBS_TEST_REGION="cn-north-4"
export OBS_TEST_TOKEN="your-temporary-token"
export OBS_MOCK_ENABLED="true"
export OBS_MOCK_PORT="8080"
```

## 最佳实践

### 1. 测试命名规范

使用BDD命名格式：
```go
Test{Function}_Should{Expected}_When{Condition}_Given{Precondition}
```

示例：
- `TestUpload_ShouldSucceed_GivenValidObject`
- `TestDownload_ShouldFail_GivenNonexistentObject`
- `TestAuthentication_ShouldConnect_GivenValidCredentials`

### 2. 资源管理

1. **使用defer进行清理**：
   ```go
   func TestExample(t *testing.T) {
       client := integration.NewTestClient(t)
       defer client.Cleanup(t)
       // ...
   }
   ```

2. **注册具体的清理函数**：
   ```go
   client.AddCleanup(func(t *testing.T) {
       // 清理特定资源
   })
   ```

### 3. Mock服务器使用

```go
// 启动Mock服务器
mockServer := integration.NewMockServer()
mockServer.Start(":8080")
defer mockServer.Stop()

// 创建指向Mock服务器的客户端
mockEndpoint := "http://localhost:8080"
client := integration.NewTestClientWithEndpoint(t, mockEndpoint)
```

### 4. 并发测试

```go
func TestConcurrentOperations(t *testing.T) {
    client := integration.NewTestClient(t)
    defer client.Cleanup(t)

    var wg sync.WaitGroup
    errors := make(chan error, 10)

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            objectKey := client.GetTestObjectKey(fmt.Sprintf("concurrent-%d.txt", id))
            err := client.TestClient().PutObject(client.GetTestBucket(), objectKey, fmt.Sprintf("test %d", id), nil)
            if err != nil {
                errors <- err
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    for err := range errors {
        t.Error(err)
    }
}
```

## 注意事项

1. **环境配置**：确保设置了必需的环境变量
2. **资源清理**：必须使用`defer client.Cleanup(t)`清理资源
3. **Mock服务器**：开发环境建议使用Mock服务器
4. **错误处理**：测试失败时提供详细的错误信息
5. **测试隔离**：每个测试用例应该独立，不依赖其他测试的状态

## 常见问题

### Q: 如何跳过集成测试？
A: 设置环境变量`OBS_SKIP_INTEGRATION_TESTS=true`

### Q: 测试失败后如何查看详细日志？
A: 使用`-v`参数运行测试：`go test -tags integration ./test/integration -v`

### Q: 如何使用Mock服务器？
A: 设置`OBS_MOCK_ENABLED=true`和`OBS_MOCK_PORT=8080`

### Q: 如何配置超时时间？
A: 在代码中使用带超时的context：
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

## 技能版本

- 版本：1.1.0
- 最后更新：2024-01-09
- 兼容性：OBS SDK v3.x
- **支持技能调用**: 可被go-sdk-dev-task技能调用和协调