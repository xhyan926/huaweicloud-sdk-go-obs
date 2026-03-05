---
name: go-sdk-ut
description: OBS SDK Go单元测试编写指南。当用户需要编写OBS SDK或类似Go SDK的单元测试、需要测试代码审查、想要改进现有测试质量、询问测试命名规范或测试最佳实践时，请使用此技能。适用于Go语言单元测试编写、测试用例优化、测试代码规范检查等场景。
---

# OBS SDK Go 单元测试编写指南

## 概述

本技能为华为云 OBS SDK 及类似 Go SDK 项目提供单元测试编写的完整指导。它包含了编写规范、最佳实践、常见反例纠正、以及工具使用指南，帮助开发者编写高质量、可维护的 Go 单元测试。

## 适用场景

- 编写 OBS SDK 或类似 Go 项目的单元测试
- 审查和改进现有测试代码
- 寻找测试命名规范或最佳实践
- 编写 HTTP Mock、并发测试等复杂测试场景
- 需要了解 Go 测试工具（testify、gomonkey、httptest）的使用方法

## 编写规范

### 1. 文件命名规范

| 测试类型 | 命名格式 | 示例 |
|---------|-----------|------|
| 内部函数测试 | `{源文件名}_internal_test.go` | `util_internal_test.go` |
| 公共函数测试 | `{源文件名}_test.go` | `client_object_test.go` |

### 2. 测试函数命名规范（BDD 风格）

采用以下格式：
```
Test{函数名}_Should{预期结果}_When{条件}_Given{前置条件}
```

**示例**：
```go
func TestStringToInt64_ShouldReturnDefaultValue_WhenValueIsEmptyString(t *testing.T)
func TestHmacSha256_ShouldReturnHash_WhenValidKeyAndValue(t *testing.T)
func TestConnDelegate_Read_ShouldReadData_WhenValidDelegate(t *testing.T)
```

### 3. 测试函数结构

每个测试函数应包含三个阶段：

```go
func TestFunctionName_ShouldResult_WhenCondition_GivenPrecondition(t *testing.T) {
    // 1. 准备阶段 - 设置测试数据
    input := "test input"
    expected := "expected output"

    // 2. 执行阶段 - 调用被测函数
    result := FunctionName(input)

    // 3. 断言阶段 - 验证结果
    assert.Equal(t, expected, result)
}
```

## 编写原则

### 1. 合并重复测试用例

避免多个测试用例测试相同场景，应使用参数化测试合并：

**正确做法**：
```go
func TestStringToInt64_ShouldReturnCorrectValue_WhenValidNumericString(t *testing.T) {
    cases := []struct {
        name     string
        input    string
        def      int64
        expected int64
    }{
        {"Positive number", "123", 0, 123},
        {"Negative number", "-456", 0, -456},
        {"Zero", "0", 0, 0},
        {"Large number", "9223372036854775807", 0, 9223372036854775807},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            result := StringToInt64(tc.input, tc.def)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### 2. 关注功能逻辑而非代码覆盖率

每个测试用例应验证明确的业务行为，而非单纯追求覆盖代码分支。

**正确做法**：
```go
func TestBase64Md5_ShouldReturnBase64EncodedMD5_WhenValidBytes(t *testing.T) {
    input := []byte("hello world")
    result := Base64Md5(input)

    // 验证输出格式（Base64）和内容正确性
    _, err := Base64Decode(result)
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

### 3. 确保测试独立性

每个测试应该独立运行，不依赖其他测试的状态或执行顺序。

**正确做法**：
```go
func TestCounter_Increment_ShouldIncreaseValue(t *testing.T) {
    counter := NewCounter()  // 每个测试创建新实例
    counter.Increment()
    assert.Equal(t, 1, counter.Value())
}
```

## 经验教训

### 1. 结构体嵌入字段的处理

Go 的结构体嵌入可能导致字段访问混淆，需要明确指定嵌入字段。

**正确写法**：
```go
output := &ListVersionsOutput{
    Versions: []Version{
        {DeleteMarker: DeleteMarker{Key: "object%2Fkey"}},  // 明确指定嵌入字段
    },
}
```

### 2. HTTP Mock 的正确使用

为 HTTP 相关测试创建自定义 mock transport，而不是依赖真实网络连接。

```go
type jsonTransport struct {
    jsonStr string
}

func (jt *jsonTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    return &http.Response{
        StatusCode: http.StatusOK,
        Body:       ioutil.NopCloser(bytes.NewReader([]byte(jt.jsonStr))),
        Header:     make(http.Header),
    }, nil
}
```

### 3. 资源清理模式

始终使用 `defer` 确保资源被清理，即使测试失败。

**正确做法**：
```go
func TestConnOperations(t *testing.T) {
    server, client := net.Pipe()
    defer server.Close()  // 确保总是执行
    defer client.Close()

    // 测试代码...
}
```

### 4. 并发测试的处理

使用 goroutine 时确保主测试等待异步操作完成。

**正确做法**：
```go
func TestConn_Write_ShouldWriteData_WhenValidDelegate(t *testing.T) {
    server, client := net.Pipe()
    defer server.Close()
    defer client.Close()

    delegate := getConnDelegate(client, 10, 100)
    data := []byte("test")

    // 启动 goroutine 读取，防止写入阻塞
    done := make(chan bool)
    go func() {
        buffer := make([]byte, 10)
        _, _ = server.Read(buffer)
        close(done)
    }()

    n, err := delegate.Write(data)
    assert.NoError(t, err)
    assert.Equal(t, len(data), n)

    <-done  // 等待读取完成
}
```

### 5. 类型匹配问题

时间相关的类型容易混淆 `time.Time` 和 `time.Duration`。

**正确写法**：
```go
err := delegate.SetDeadline(time.Now().Add(delegate.socketTimeout))  // 正确转换为 time.Time
```

### 6. XML 输出断言

XML 输出包含属性时，断言应该灵活处理。

**正确写法**：
```go
assert.Contains(t, result, "<Grant><Grantee")  // 不包含闭合 >，允许属性存在
```

## 常见正例与反例

### 断言的使用

**反例**：
```go
if result != 123 {  // 原生 if，缺少详细错误信息
    t.Error("failed")
}
```

**正例**：
```go
assert.Equal(t, int64(123), result)  // 使用 assert，自动输出详细差异
```

### 错误处理测试

**反例**：
```go
assert.Error(t, err)  // 只验证有错误，不验证错误类型
```

**正例**：
```go
assert.Error(t, err)
assert.Contains(t, err.Error(), "illegal base64 data")  // 验证错误消息
```

### 边界条件测试

**反例**：
```go
func TestStringToInt64(t *testing.T) {
    assert.Equal(t, int64(100), StringToInt64("100", 0))
    // 只测试正常情况，缺少边界
}
```

**正例**：
```go
func TestStringToInt64_ShouldHandleBoundaryValues(t *testing.T) {
    // 空字符串使用默认值
    assert.Equal(t, int64(0), StringToInt64("", 0))
    assert.Equal(t, int64(999), StringToInt64("", 999))

    // 无效字符串使用默认值
    assert.Equal(t, int64(-1), StringToInt64("invalid", -1))

    // 边界数值
    assert.Equal(t, int64(9223372036854775807), StringToInt64("9223372036854775807", 0))
    assert.Equal(t, int64(-9223372036854775808), StringToInt64("-9223372036854775808", 0))
}
```

## 测试工具使用

### testify 库

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// assert：失败继续执行
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.Contains(t, str, substr)
assert.NotNil(t, obj)
assert.Empty(t, value)

// require：失败立即停止
require.NoError(t, err)  // 依赖 err 成功才能继续的测试
require.NotNil(t, obj)   // 依赖 obj 非空才能继续的测试
```

### httptest 使用

```go
import (
    "net/http/httptest"
)

server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}))
defer server.Close()

// 使用 server.URL 进行请求
```

### gomonkey 使用（Mock 函数）

```go
import "github.com/agiledragon/gomonkey/v2"

// Mock 函数
patches := gomonkey.ApplyFunc(ExternalFunction, func() string {
    return "mocked value"
})
defer patches.Reset()

// Mock 方法
patches.ApplyMethod(&obj, "MethodName", func(_ *MyType) string {
    return "mocked method result"
})
```

## 检查清单

在提交测试代码前，检查以下内容：

- [ ] 测试命名遵循 BDD 风格
- [ ] 每个测试用例有明确的业务意义
- [ ] 测试之间相互独立
- [ ] 使用 defer 清理资源
- [ ] 边界条件已测试
- [ ] 错误情况已测试
- [ ] 断言消息清晰明确
- [ ] 避免重复的测试用例
- [ ] 使用参数化测试合并相似场景
- [ ] 所有测试通过

## 工作流程

当使用此技能时：

1. **理解需求** - 确认用户需要编写什么类型的测试
2. **应用规范** - 根据编写规范提供测试代码建议
3. **检查质量** - 对照检查清单验证测试质量
4. **提供改进** - 指出反例并提供正例
5. **确保通过** - 确保测试用例能够正确执行
