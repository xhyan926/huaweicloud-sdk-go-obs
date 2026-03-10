# OBS SDK Go 测试工程规范

## 测试工程总览

OBS SDK Go测试工程采用了分层的测试架构，包括单元测试、集成测试、性能测试和模糊测试四种类型。通过build tags实现测试隔离，确保每种测试类型可以独立运行。

## 测试架构

### 测试分层

```
测试类型        | 文件位置           | Build Tag | 运行方式                  | 目的
---------------|-------------------|-----------|-------------------------|------
单元测试       | obs/*_test.go     | unit      | go test -tags unit      | 验证代码逻辑正确性
集成测试       | obs/test/...      | integration| go test -tags integration| 验证端到端流程
性能测试       | obs/*_benchmark.go| perf      | go test -tags perf      | 评估性能表现
模糊测试       | obs/*_fuzz_test.go| fuzz      | go test -tags fuzz      | 发现潜在漏洞
```

### 目录结构

```
obs/
├── *_test.go                    # 单元测试（build tag: unit）
├── *_internal_test.go          # 私有接口测试（build tag: unit）
├── *_fuzz_test.go              # 模糊测试（build tag: fuzz）
├── *_benchmark.go              # 性能测试（build tag: perf）
├── test/
│   ├── config/
│   │   ├── test_config.go      # 统一测试配置管理
│   │   └── integration_env.go  # 集成测试环境配置
│   ├── integration/
│   │   ├── client.go          # 集成测试客户端
│   │   ├── e2e/               # 端到端测试
│   │   └── fixtures/          # 测试数据
│   └── mock_server/          # Mock服务器
│       ├── server.go         # Mock服务器实现
│       ├── handler.go        # Mock HTTP处理器
│       └── responses/        # Mock响应数据
```

## 测试类型说明

### 单元测试

**特点**：
- 优先使用接口Mock和testable设计
- 仅在无法通过其他方式打桩时使用gomonkey
- 代码覆盖率 > 90%
- 使用BDD命名规范

**运行方式**：
```bash
# 运行所有单元测试
go test -tags unit ./obs -v

# 运行特定测试
go test -tags unit ./obs -run TestBucket -v

# 生成覆盖率报告
go test -tags unit ./obs -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 集成测试

**特点**：
- 连接真实OBS服务或Mock服务器
- 从客户视角测试核心功能
- 自动跳过机制（环境变量配置）
- 资源自动清理

**运行方式**：
```bash
# 运行集成测试（需要配置环境变量）
go test -tags integration ./obs/test/integration -v

# 使用Mock服务器
go test -tags integration ./obs/test/integration -v -run TestWithMock

# 跳过集成测试
OBS_SKIP_INTEGRATION_TESTS=true go test -tags integration ./obs/test/integration -v
```

### 性能测试

**特点**：
- 轻量级测试：小文件（1MB）、低并发（10）、短时长（1s）
- 深度测试：大文件（100MB-1GB）、高并发（100-1000）、长时长（30s）
- 建立性能基线，检测性能退化
- 资源监控（CPU、内存、带宽）

**运行方式**：
```bash
# 轻量级性能测试
go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem

# 深度性能测试
go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem

# 生成性能报告
go test -tags perf ./obs -bench=. > bench.out
go tool benchmark -html bench.out > report.html
```

### 模糊测试

**特点**：
- 针对关键输入解析函数
- 测试XML解析、URL解析、签名验证等
- 发现潜在漏洞和边界条件
- 使用Go内置的模糊测试框架

**运行方式**：
```bash
# 运行模糊测试
go test -tags fuzz ./obs -fuzz=.

# 运行特定函数
go test -tags fuzz ./obs -fuzz=FuzzXmlParsing

# 设置测试时长
go test -tags fuzz ./obs -fuzz=. -fuzztime=60s
```

## 测试技能

### 技能列表

| 技能名称 | 描述 | 使用场景 |
|---------|------|---------|
| go-sdk-ut | 单元测试编写指南 | 编写单元测试、测试代码审查 |
| go-sdk-integration | 集成测试编写指南 | 编写集成测试、端到端测试 |
| go-sdk-perf | 性能测试编写指南 | 编写性能测试、性能分析 |
| go-sdk-fuzz | 模糊测试编写指南 | 编写模糊测试、安全测试 |

### 技能使用

```bash
# 单元测试
/go-sdk-ut

# 集成测试
/go-sdk-integration

# 性能测试
/go-sdk-perf

# 模糊测试
/go-sdk-fuzz
```

## 测试配置

### 环境变量

| 变量名 | 必需 | 说明 |
|--------|------|------|
| OBS_TEST_AK | 是 | 访问密钥 Access Key |
| OBS_TEST_SK | 是 | 访问密钥 Secret Key |
| OBS_TEST_ENDPOINT | 是 | OBS服务端点 |
| OBS_TEST_BUCKET | 是 | 测试桶名称 |
| OBS_TEST_REGION | 否 | 区域 |
| OBS_TEST_TOKEN | 否 | 临时安全令牌 |
| OBS_MOCK_ENABLED | 否 | 是否启用Mock服务器 |
| OBS_MOCK_PORT | 否 | Mock服务器端口 |
| OBS_SKIP_INTEGRATION_TESTS | 否 | 是否跳过集成测试 |
| OBS_SKIP_FUZZ_TESTS | 否 | 是否跳过模糊测试 |
| OBS_SKIP_PERF_TESTS | 否 | 是否跳过性能测试 |

### 测试配置文件

```go
// obs/test_config.go
type TestConfig struct {
    // 基础配置
    AccessKey       string
    SecretKey       string
    SecurityToken   string
    Endpoint        string
    Region          string

    // 测试资源
    TestBucket      string
    TestObject      string

    // 性能测试配置
    PerfLargeFileSize   int64
    PerfConcurrency     int
    PerfTestDuration    int

    // 跳过测试标记
    SkipIntegrationTests bool
    SkipFuzzTests       bool
    SkipPerfTests       bool
}
```

## 测试流程

### 开发流程

```
┌─────────────────────────────────────────────────────────────┐
│ 1. 功能分析                                                  │
│    - 识别缺失功能                                           │
│    - 确定功能优先级                                         │
│    └─→ 2. 测试策略制定                                      │
│          ├── 确定测试类型                                   │
│          ├── 定义测试覆盖范围                               │
│          └─→ 3. 计划制定                                    │
│                ├── 创建实施计划                              │
│                └─→ 4. 子任务拆分                            │
│                      ├── 子任务1: 数据模型                   │
│                      ├── 子任务2: Trait层实现                │
│                      ├── 子任务3: 客户端方法                 │
│                      ├── 子任务4: 单元测试 (/go-sdk-ut)      │
│                      ├── 子任务5: 集成测试 (/go-sdk-integration) │
│                      ├── 子任务6: 性能测试 (/go-sdk-perf)      │
│                      └── 子任务7: 模糊测试 (/go-sdk-fuzz)     │
│                      └─→ 5. 多层测试验证                     │
│                            ├── 单元测试验证                  │
│                            ├── 集成测试验证                  │
│                            ├── 性能测试验证                  │
│                            └── 模糊测试验证                  │
│                      └─→ 6. 测试报告生成                   │
└─────────────────────────────────────────────────────────────┘
```

### 测试验证清单

#### 单元测试
- [ ] 代码覆盖率 > 90%
- [ ] 所有测试通过
- [ ] 使用BDD命名规范
- [ ] 优先使用接口Mock和testable设计
- [ ] 仅在无法通过其他方式打桩时使用gomonkey

#### 集成测试
- [ ] 环境变量配置正确
- [ ] 测试资源清理
- [ ] 端到端场景覆盖
- [ ] 错误场景测试
- [ ] Mock服务器正常工作

#### 性能测试
- [ ] 性能基线对比
- [ ] 无性能退化
- [ ] 轻量级测试通过
- [ ] 资源监控正常

#### 模糊测试
- [ ] 关键输入覆盖
- [ ] 无崩溃和panic
- [ ] 超时设置合理

## 测试命令

### Makefile命令

```makefile
# 测试相关命令
.PHONY: test test-unit test-integration test-perf test-fuzz test-all test-report

# 运行所有单元测试
test-unit:
	go test -tags unit ./obs -run=TestUnit -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# 运行集成测试
test-integration:
	@if [ -z "$(OBS_TEST_AK)" ]; then \
		echo "Error: OBS_TEST_AK not set"; \
		exit 1; \
	fi
	go test -tags integration ./test/integration -v

# 运行轻量级性能测试
test-perf-light:
	go test -tags perf ./obs -bench=BenchmarkLight -benchtime=1s -benchmem

# 运行深度性能测试
test-perf-deep:
	@if [ -z "$(OBS_TEST_AK)" ]; then \
		echo "Error: OBS_TEST_AK not set"; \
		exit 1; \
	fi
	go test -tags perf ./obs -bench=BenchmarkDeep -benchtime=30s -benchmem

# 运行模糊测试
test-fuzz:
	go test -tags fuzz ./obs -fuzz=. -fuzztime=30s

# 生成测试报告
test-report:
	@echo "Generating test report..."
	@go run ./cmd/test-report/main.go

# 运行所有测试
test-all: test-unit
	@echo "All tests completed"
```

### 技能调用命令

```bash
# 单元测试
/go-sdk-ut

# 集成测试
/go-sdk-integration --modules=bucket,object,auth

# 性能测试
/go-sdk-perf --type=deep --concurrency=100

# 模糊测试
/go-sdk-fuzz --targets=xml,url,auth
```

## 测试报告

### 报告生成

```go
// 测试报告结构
type TestReport struct {
    Timestamp   time.Time              `json:"timestamp"`
    Commit      string                 `json:"commit"`
    Branch      string                 `json:"branch"`
    UnitTests   UnitTestReport         `json:"unit_tests"`
    Integration IntegrationTestReport  `json:"integration_tests"`
    Performance PerformanceTestReport  `json:"performance_tests"`
    Fuzz        FuzzTestReport         `json:"fuzz_tests"`
}
```

### 报告展示

- 测试覆盖率趋势
- 性能基准对比
- 错误统计分析
- 资源使用情况

## 最佳实践

### 1. 测试命名规范

使用BDD格式：`Test{Function}_Should{Expected}_When{Condition}_Given{Precondition}`

示例：
- `TestUpload_ShouldSucceed_GivenValidObject`
- `TestDownload_ShouldFail_GivenNonexistentObject`

### 2. 测试隔离

- 每个测试用例独立运行
- 使用`defer client.Cleanup(t)`清理资源
- 避免测试间相互依赖

### 3. Mock使用

```go
// 接口Mock（优先级1）
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

// testable设计（优先级2）
func processRequest(req *Request) (*Response, error) {
    parsed := parseRequest(req)
    validated := validateRequest(parsed)
    return buildResponse(validated), nil
}

// httptest（优先级3）
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}))

// gomonkey（优先级4，谨慎使用）
patches := gomonkey.ApplyFunc(time.Now, func() time.Time {
    return time.Date(2023, 1, 1, 0, 0, 0, time.UTC)
})
defer patches.Reset()
```

### 4. 性能测试原则

- 轻量级测试定期运行
- 深度测试单独运行
- 建立性能基线
- 设置性能退化告警

### 5. 模糊测试要点

- 测试关键输入解析函数
- 设置合理的超时和资源限制
- 关注崩溃和panic
- 生成语料库用于持续测试

## 故障排除

### 常见问题

1. **测试运行失败**
   - 检查build tag
   - 确认环境变量
   - 验证依赖配置

2. **性能测试不稳定**
   - 增加测试时长
   - 使用固定数据
   - 检查网络环境

3. **Mock服务器问题**
   - 检查端口占用
   - 验证路由配置
   - 查看请求日志

### 调试技巧

1. **启用详细日志**
   ```bash
   go test -tags unit ./obs -v
   ```

2. **使用特定测试**
   ```bash
   go test -tags unit ./obs -run TestSpecificFunction
   ```

3. **分析测试结果**
   ```bash
   go tool cover -html=coverage.out
   ```

## 相关文档

- [单元测试规范](unit-testing.md)
- [集成测试规范](integration-testing.md)
- [性能测试规范](performance-testing.md)
- [模糊测试规范](fuzz-testing.md)
- [测试迁移指南](migration-guide.md)
- [测试技能使用指南](skills-usage.md)

## 版本信息

- 版本：1.0.0
- 最后更新：2024-01-01
- 兼容性：OBS SDK v3.x
- Go版本：1.18+