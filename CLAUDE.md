# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 概述

这是华为云 OBS (Object Storage Service) Go SDK。它提供了一个 Go 客户端库，用于与华为云 OBS 服务交互，该服务兼容 S3 API。SDK 支持 OBS 专用签名 (SignatureObs) 和 AWS S3 签名 (SignatureV2, SignatureV4)。

## 运行示例

`examples/` 目录包含展示各种 SDK 操作的示例代码。运行示例：

```bash
cd main
# 编辑 obs_go_sample.go，设置您的 endpoint、ak、sk、bucketName 和 objectKey
go run obs_go_sample.go
```

`examples/` 中的单独示例文件可以直接运行，但需要先设置凭据和配置。

## 包结构

SDK 组织在 `obs/` 包中，包含以下关键架构组件：

### 核心客户端 (`client_base.go`, `conf.go`)
- **`ObsClient`**: 主客户端结构，通过 `obs.New(ak, sk, endpoint, configurers...)` 创建
- **`config`**: 通过函数式配置器管理的配置结构
- 配置使用 `WithXXX` 函数的"函数式选项"模式

### API 层次
1. **客户端方法** (`client_bucket.go`, `client_object.go`, `client_part.go`):
   - 公共 API 方法，如 `PutObject`、`GetObject`、`CreateBucket` 等
   - 方法通常返回结果结构和错误

2. **特性层** (`trait_bucket.go`, `trait_object.go`, `trait_part.go`):
   - 存储桶、对象和分块操作的内部实现
   - 由客户端方法调用

3. **HTTP 层** (`http.go`, `auth.go`, `authV2.go`, `authV4.go`):
   - `doAction()` - 中央 HTTP 请求处理器
   - 请求签名 (v2 和 v4 签名)
   - 重试逻辑和连接管理

4. **模型层** (`model_bucket.go`, `model_object.go`, `model_part.go` 等):
   - API 调用的输入/输出结构
   - 请求/响应的 XML 序列化

### 扩展系统 (`extension.go`)
API 调用通过可变参数支持扩展选项：
- `WithProgress(listener)` - 进度回调
- `WithReqPaymentHeader(requester)` - 请求者付费
- `WithTrafficLimitHeader(limit)` - 流量限制
- `WithCallbackHeader(callback)` - 上传回调
- `WithCustomHeader(key, value)` - 自定义头

### 文件操作 (`transfer.go`, `client_resume.go`)
- **断点续传上传/下载**: `UploadFile()` 和 `DownloadFile()`
- 支持并发分块上传的分块操作

### 认证
- `BasicSecurityProvider` (provider.go): 使用 AK/SK 的默认提供者
- 支持带安全令牌的临时凭证
- 三种签名类型：SignatureV2、SignatureV4、SignatureObs

### 常量 (`const.go`)
- 集中定义所有 HTTP 头、参数名称和常量
- 两个头前缀：`x-amz-` (AWS S3) 和 `x-obs-` (OBS 专用)

### 错误处理 (`error.go`)
- `ObsError`: 结构化错误，包含 Status、Code、Message、RequestId
- 类型断言模式：`if obsError, ok := err.(obs.ObsError); ok`

### 日志 (`log.go`)
- `InitLog(path, maxSize, backupCount, level, isCompress)`
- 日志级别：LEVEL_DEBUG、LEVEL_INFO、LEVEL_WARN、LEVEL_ERROR
- `CloseLog()` 应使用 defer 调用以刷新日志

### 单元测试原则
**合并重复测试用例**：
- 当多个测试用例针对相同场景时，应合并为一个具有充分覆盖的测试用例
- 保留最有代表性的测试场景，删除冗余和重复的测试
- 通过参数化测试或组合多种场景来提高测试覆盖率

**采用BDD风格命名规范**：
- 测试命名应采用 Should_xxx_When_xxx_Given_xxx 格式
- 命名应清晰表达测试目的和预期行为

**提升测试质量**：
- 优先关注功能逻辑而非代码覆盖率
- 确保每个测试用例都有明确的业务价值

**测试工具**：
- 使用testify进行断言
- 使用httptest模拟http server
- 使用gomonkey进行mock
