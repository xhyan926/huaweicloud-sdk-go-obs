# OBS SDK Go 测试工程规范

## 概述

本文档定义了OBS SDK Go项目的完整测试工程化规范，包括测试架构、开发规范、集成规范、性能规范、模糊测试规范以及测试报告系统规范。

## 目录

1. [测试架构](#测试架构) - 测试分层架构设计
2. [测试开发规范](#测试开发规范) - 测试代码编写规范
3. [测试集成规范](#测试集成规范) - 集成测试实施规范
4. [测试性能规范](#测试性能规范) - 性能测试实施规范
5. [测试模糊规范](#测试模糊规范) - 模糊测试实施规范
6. [测试报告规范](#测试报告规范) - 测试报告系统规范
7. [相关文档](#相关文档) - 其他相关测试文档

## 测试架构

### 分层架构

OBS SDK Go测试采用四层架构，通过build tags实现测试隔离：

```
测试分层
├── 单元测试层 (obs/*_test.go)
│   ├── 公共接口测试 (obs/client_*_test.go)
│   ├── 私有接口测试 (obs/*_internal_test.go)
│   └── build tag: unit
├── 集成测试层 (test/integration/**_test.go)
│   ├── 认证测试 (test/integration/auth/*_test.go)
│   ├── 存储桶测试 (test/integration/bucket/*_test.go)
│   ├── 对象测试 (test/integration/object/*_test.go)
│   ├── 分块上传测试 (test/integration/object/*_test.go)
│   ├── 断点续传测试 (test/integration/object/*_test.go)
│   └── build tag: integration
├── 性能测试层 (obs/*_benchmark_test.go)
│   ├── 上传性能测试 (obs/upload_benchmark_test.go)
│   ├── 下载性能测试 (obs/download_benchmark_test.go)
│   ├── 分块性能测试 (obs/multipart_benchmark_test.go)
│   ├── 并发性能测试 (obs/concurrent_benchmark_test.go)
│   ├── 性能基线 (obs/performance_baseline_test.go)
│   └── build tag: perf
└── 模糊测试层 (obs/*_fuzz_test.go)
    ├── XML解析模糊测试 (obs/xml_parser_fuzz_test.go)
    ├── URL解析模糊测试 (obs/url_parser_fuzz_test.go)
    ├── 签名验证模糊测试 (obs/signature_fuzz_test.go)
    ├── 响应解析模糊测试 (obs/response_parser_fuzz_test.go)
    └── build tag: fuzz
```

### Build Tags 规范

| 测试类型 | Build Tag | 文件位置 | 运行命令 |
|---------|----------|----------|----------|
| 单元测试 | `unit` | `obs/*_test.go` | `go test -tags unit ./obs -v` |
| 集成测试 | `integration` | `test/integration/**_test.go` | `go test -tags integration ./test/integration -v` |
| 性能测试 | `perf` | `obs/*_benchmark_test.go` | `go test -tags perf ./obs -bench=. -benchtime=1s -v` |
| 模糊测试 | `fuzz` | `obs/*_fuzz_test.go` | `go test -tags fuzz ./obs -fuzz=. -fuzztime=30s -v` |
| 测试报告 | `test_report` | `obs/test_report.go` | `go test -tags test_report ./obs -bench=. -v` |

### 测试文件组织

```
obs/
├── *_test.go                    # 单元测试文件
├── *_benchmark_test.go           # 性能测试文件
├── *_fuzz_test.go              # 模糊测试文件
├── test_report.go              # 测试报告系统
├── test_config.go              # 测试配置
├── test_helpers.go             # 测试辅助工具
└── test_fixtures.go            # 测试固件

test/integration/
├── suite.go                     # 集成测试套件
├── client.go                    # 集成测试客户端
├── fixtures/                    # 测试数据
├── bucket/                     # 存储桶测试
├── object/                     # 对象测试
└── auth/                       # 认证测试

benchmarks/
├── performance_baseline.json    # 性能基线数据
├── test_reports/              # 测试报告
└── test_archives/              # 测试报告归档
```

## 测试开发规范

### 命名规范

采用BDD（Behavior-Driven Development）命名规范：

```
Test{Function}_Should{Expected}_When{Condition}_Given{Precondition}
```

#### 命名模板说明

- `Function`: 被测试的函数或功能名称
- `Expected`: 预期的行为或结果
- `Condition`: 触发该行为的条件
- `Precondition`: 前置条件或状态

#### 命名示例

- **集成测试**: `TestAuthentication_ShouldConnectSuccessfully_GivenValidCredentials`
- **对象测试**: `TestUpload_ShouldSucceed_GivenValidObject`
- **性能测试**: `BenchmarkUpload_ShouldMeetBaseline_GivenNormalConditions`
- **模糊测试**: `FuzzXmlParser_ShouldHandleMaliciousInput_GivenXSSPattern`

### 测试文件组织

每个测试文件应该：

1. **文件头注释**
   ```go
   //go:build unit
   // Copyright 2019 Huawei Technologies Co.,Ltd.
   // 文件说明
   ```

2. **导入包**
   ```go
   import (
       "testing"
       "github.com/huaweicloud/huaweicloud-sdk-go-obs"
       "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/config"
   )
   ```

3. **测试函数结构**
   ```go
   func TestFunctionName(t *testing.T) {
       // 测试设置
   }
   ```

4. **测试数据准备**
   ```go
   func setupTestData(t *testing.T) {
       // 准备测试数据
   }
   ```

5. **测试清理**
   ```go
   func cleanupTestData(t *testing.T) {
       // 清理测试资源
   }
   ```

### 测试固件管理

测试固件存放在 `obs/test_fixtures.go` 中：

```go
// 常用固件
var (
    TestBucket = "test-bucket"
    TestObject = "test-object.txt"
    TestLargeFile = "test-large-file.bin"
)

// 获取测试固件
func GetTestFixture(name string) []byte {
    return []byte(fixtureContent)
}
```

### Mock 框架使用

在集成测试中使用 Mock 服务器：

```go
// 启动 Mock 服务器
mockServer := integration.NewMockServer()
mockServer.Start(":8080")
defer mockServer.Stop()

// 创建指向 Mock 服务器的客户端
mockEndpoint := "http://localhost:8080"
client := integration.NewTestClientWithEndpoint(t, mockEndpoint)
```

### 错误处理规范

1. **错误信息明确**
   ```go
   if err != nil {
       t.Errorf("Failed to upload object: %v", err)
   }
   ```

2. **错误类型检查**
   ```go
   if obsError, ok := err.(obs.ObsError); ok {
       // 处理OBS错误
   }
   ```

3. **日志记录**
   ```go
   t.Logf("Attempting to upload object: %s", objectKey)
   ```

## 测试集成规范

### 集成测试客户端使用

集成测试客户端提供统一的测试工具：

```go
// 创建测试客户端
client := integration.NewTestClient(t)
defer client.Cleanup(t)

// 获取测试配置
cfg := integration.LoadTestConfig()
bucket := cfg.GetTestBucket()

// 获取测试对象键
objectKey := client.GetTestObjectKey("test-object.txt")

// 注册清理函数
client.AddCleanup(func(t *testing.T) {
    // 清理特定资源
})
```

### 环境变量配置

必需环境变量：

```bash
export OBS_TEST_AK="your-access-key"
export OBS_TEST_SK="your-secret-key"
export OBS_TEST_ENDPOINT="https://obs.cn-north-4.myhuaweicloud.com"
export OBS_TEST_BUCKET="your-test-bucket"
```

可选环境变量：

```bash
export OBS_TEST_REGION="cn-north-4"
export OBS_TEST_TOKEN="your-temporary-token"
export OBS_MOCK_ENABLED="true"
export OBS_MOCK_PORT="8080"
```

### 测试资源管理

1. **资源创建**
   - 使用唯一的测试对象名称
   - 避免命名冲突
   - 及时清理测试资源

2. **资源清理**
   - 使用defer进行清理
   - 注册多个清理函数
   - 验证清理结果

3. **测试数据准备**
   - 使用固定的测试数据
   - 避免依赖测试间资源
   - 确保数据一致性

## 测试性能规范

### 性能测试类型

| 测试类型 | 文件大小 | 并发数 | 时长 | 说明 |
|---------|--------|--------|------|------|
| 轻量级测试 | 1MB | 1-10 | 1s | 日常测试，快速验证 |
| 深度测试 | 100MB-1GB | 100-1000 | 30s | 全面性能评估 |
| 并发测试 | 5MB | 10-500 | 10s | 并发性能测试 |
| 内存测试 | 10MB | - | - | 内存使用监控 |

### 性能基线建立

1. **基线数据收集**
   - 运行完整的性能测试套件
   - 收集关键操作的基准数据
   - 建立合理的性能阈值

2. **基线存储**
   - 基线数据存储在 `benchmarks/performance_baseline.json`
   - 包含吞吐量、延迟、内存等指标

3. **基线更新策略**
   - 定期更新基线数据
   - 基线老化后重新测试
   - 环境变化时重新建立基线

### 性能退化检测

退化检测规则：
- 基线值的 90% 为阈值
- 低于阈值时触发警告
- 低于阈值 70% 时标记为严重问题
- 记录退化原因和修复措施

## 测试模糊规范

### 模糊测试目标

需要模糊测试的关键函数类型：

1. **XML解析函数**
   - `TransToXml` - 结构体转XML
   - `XmlToTrans` - XML转结构体
   - `ConvertAclToXml` - ACL转换
   - `ConvertLifecycleConfigurationToXml` - 生命周期配置

2. **URL解析函数**
   - `formatUrls` - URL格式化
   - `prepareBaseURL` - 基础URL准备
   - `prepareObjectKey` - 对象键准备
   - `ParseResponseToBaseModel` - 响应解析

3. **签名验证函数**
   - `HmacSha256` - HMAC签名
   - `getSignature` - 签名生成
   - `Md5Hash` - MD5哈希
   - `HexSha256` - SHA256哈希

4. **响应解析函数**
   - `ParseResponseToObsError` - 错误响应解析

### 模糊测试执行

运行模糊测试：

```bash
# 运行所有模糊测试
go test -tags fuzz ./obs -fuzz=. -fuzztime=60s

# 运行特定函数的模糊测试
go test -tags fuzz ./obs -fuzz=FuzzXmlParser

# 指定模糊测试时长
go test -tags fuzz ./obs -fuzz=. -fuzztime=120s

# 使用多个工作线程
go test -tags fuzz ./obs -fuzz=. -parallel=4
```

### 模糊测试配置

模糊测试配置参数：
- 最大输入大小：100KB
- 测试时长：60秒（默认）
- 工作线程数：4
- 超时阈值：10秒
- 内存限制：2GB

### 崩溃分析

1. **崩溃记录**
   - 自动记录崩溃输入
   - 保存堆栈跟踪
   - 标记可重现性

2. **崩溃分析流程**
   - 识别崩溃模式
   - 重现崩溃问题
   - 定位根本原因
   - 生成修复建议

## 测试报告规范

### 报告类型

1. **JSON报告**
   - 测试结果详情
   - 机器可读格式
   - 用于自动化处理

2. **HTML报告**
   - 可视化展示
   - 美观的样式
   - 用于人工查看

3. **测试总览**
   - 汇总所有测试类型的结果
   - 覆盖率统计
   - 趋势分析

### 报告生成命令

```bash
# 生成测试报告
make test-report-generate

# 聚合测试报告
make test-report-aggregate

# 生成HTML总览
make test-html-summary

# 清除测试报告
make test-report-clean
```

### 报告存储位置

- JSON报告：`benchmarks/test_reports/`
- HTML总览：`benchmarks/test_summary.html`
- 归档报告：`benchmarks/test_archives/`

## 测试覆盖率目标

### 覆盖率要求

- 公共接口测试覆盖率 > 90%
- 核心功能集成测试覆盖率 > 85%
- 性能关键操作覆盖率 100%
- 安全关键函数模糊测试覆盖率 100%

### 覆盖率计算

使用Go标准覆盖率工具：

```bash
# 生成覆盖率报告
go test -tags unit ./obs -coverprofile=coverage.out -v

# 生成HTML覆盖率报告
go tool cover -html=coverage.html -o coverage.out
```

## 相关文档

1. [integration-testing.md](#测试集成规范) - 集成测试详细规范
2. [performance-testing.md](#性能测试规范) - 性能测试详细规范
3. [fuzz-testing.md](#测试模糊规范) - 模糊测试详细规范
4. [unit-testing.md](#单元测试规范) - 单元测试规范

## 测试流程最佳实践

### 1. 测试编写流程

1. **需求分析**：理解功能需求和测试范围
2. **测试设计**：设计测试用例和测试数据
3. **测试实现**：编写测试代码，遵循规范
4. **测试验证**：运行测试，验证覆盖率
5. **文档更新**：同步更新测试文档

### 2. 测试执行流程

1. **环境准备**：配置测试环境变量
2. **依赖安装**：安装测试依赖
3. **测试运行**：按测试类型运行测试
4. **结果分析**：分析测试结果，识别问题
5. **报告生成**：生成测试报告

### 3. 测试维护流程

1. **定期运行**：定期运行测试，保持基线更新
2. **问题修复**：及时修复测试失败的问题
3. **规范更新**：根据实际使用情况更新规范
4. **文档同步**：保持文档与代码同步

## 测试质量标准

### 代码质量

- 测试代码必须符合Go编码规范
- 测试代码必须有充分的注释
- 测试代码应该使用testify进行断言

### 测试完整性

- 每个公共函数都有对应的单元测试
- 每个核心功能都有集成测试覆盖
- 性能关键函数都有性能基准
- 安全关键函数都有模糊测试覆盖

### 测试可维护性

- 测试代码结构清晰，易于理解和修改
- 测试数据管理规范，易于扩展
- 测试文档完整，便于参考

## 测试技能使用

### go-sdk-ut

单元测试编写指南，用于快速验证代码逻辑正确性。

**使用场景**：
- 编写新的单元测试
- 代码审查时的单元测试
- 验证修复后的回归

### go-sdk-integration

集成测试编写指南，用于验证端到端功能。

**使用场景**：
- 为新功能编写集成测试
- Bug修复后编写回归测试
- 定期运行集成测试验证功能

### go-sdk-fuzz

模糊测试编写指南，用于发现安全漏洞。

**使用场景**：
- 为新的安全关键函数编写模糊测试
- Bug修复后重新运行模糊测试
- 定期运行模糊测试发现新问题

### go-sdk-perf

性能测试编写指南，用于评估和优化性能。

**使用场景**：
- 为性能优化编写性能测试
- 定期运行性能测试监控性能基线
- 建立新的性能基线

## 测试技能协调

### 技能调用关系

```
开发任务分解
    ↓
    ↓
go-sdk-dev-task (任务协调)
    ↓
    ┌─────────────┬────────────┐
    │             │             │
    │   go-sdk-integration     │
    │   go-sdk-fuzz          │
    │   go-sdk-perf          │
    └─────────────┴────────────┘
    ↓                    ↓              ↓
测试代码生成 ←─────────────→ 测试代码执行
```

### 技能协调机制

1. **去重检查**
   - 避免生成重复的测试用例
   - 统一测试命名规范
   - 共享测试数据和固件

2. **策略一致性**
   - 根据任务类型选择合适的测试策略
   - 协调测试覆盖率和效率
   - 统一测试执行时机

3. **资源管理**
   - 统一测试数据目录
   - 避免测试资源冲突
   - 协调测试报告存储

4. **进度同步**
   - 向go-sdk-dev-task报告测试生成进度
   - 汇总测试结果
   - 协调测试执行时机

## 测试环境管理

### 测试环境要求

1. **独立测试环境**
   - 避免与其他开发环境冲突
   - 配置独立的测试数据库
   - 使用专用的测试存储

2. **开发环境**
   - 集成测试：需要真实OBS环境
   - 性能测试：需要稳定网络环境
   - 模糊测试：建议在隔离环境运行

## 测试安全注意事项

1. **凭据管理**
   - 测试凭据不要提交到版本控制
   - 使用环境变量管理测试凭据
   - 定期轮换测试凭据

2. **数据安全**
   - 避免在测试中处理敏感数据
   - 测试数据及时清理
   - 测试报告避免泄露敏感信息

3. **环境隔离**
   - 测试环境与生产环境隔离
   - 使用独立的测试数据存储
   - 测试结束后清理测试环境

## 测试发布检查清单

在发布前确保：

- [ ] 所有测试类型覆盖率达标
- [ ] 测试报告功能正常工作
- [ ] 性能基线已建立
- [ ] 模糊测试通过关键函数
- [ ] 测试文档完整更新
- [ ] 无已知的测试阻塞问题
- [ ] 测试环境配置正确

## 版本历史

- v1.0.0 (2024-01-01): 初始版本
- v1.1.0 (2024-01-09): 添加测试报告系统，完成阶段7
- v1.2.0 (2024-01-10): 完成阶段8，完善文档规范