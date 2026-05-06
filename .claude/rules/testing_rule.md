# Testing Rules

约束 OBS SDK 的测试分层、命名规范、断言工具和测试组织方式。

---

## TEST-01: 按测试类型使用对应 build tag

**条款**: 禁止测试文件缺少 build tag，禁止使用非标准 tag 名（如 `//go:build test`）。各类型对应的正确 tag：
- 单元测试：`//go:build unit`
- 集成测试：`//go:build integration`
- 性能测试：`//go:build perf`
- 模糊测试：`//go:build fuzz`

### 正确示例

```go
//go:build unit

package obs

import "testing"

func TestDeleteBucket_ShouldReturnError_WhenInputIsNil(t *testing.T) {
    // ...
}
```

```go
//go:build integration

package obs

import "testing"

func TestBucketLifecycle_ShouldCreateAndDelete(t *testing.T) {
    // ...
}
```

### 错误示例

```go
// 缺少 build tag
package obs

import "testing"

func TestSomething(t *testing.T) { ... }
```

```go
// 错误：tag 名不匹配
//go:build test

package obs
```

### 验证清单
- [ ] 是否不存在测试文件缺少 build tag 的情况
- [ ] 是否不存在 build tag 未放在文件最开头（package 声明之前）的情况
- [ ] 是否不存在使用已废弃的 `// +build` 格式的情况

---

## TEST-02: 测试函数使用 BDD 风格命名

**条款**: 禁止使用非 BDD 格式的测试命名（如 `TestDeleteBucket`、`TestDeleteBucketNilInput`、`TestDeleteBucket_Error`）。

### 正确示例

```go
func TestDeleteBucket_ShouldReturnError_WhenInputIsNil(t *testing.T) { ... }
func TestGetObject_ShouldReturnContent_WhenObjectExists(t *testing.T) { ... }
func TestCreateBucket_ShouldSucceed_GivenValidInput(t *testing.T) { ... }
```

### 错误示例

```go
func TestDeleteBucket(t *testing.T) { ... }           // 错误：缺少行为描述
func TestDeleteBucketNilInput(t *testing.T) { ... }    // 错误：未使用 BDD 格式
func TestDeleteBucket_Error(t *testing.T) { ... }      // 错误：格式不规范
```

### 验证清单
- [ ] 是否不存在测试名不以 `Test` 开头的情况
- [ ] 是否不存在缺少 `Should` 描述预期行为的情况
- [ ] 是否不存在缺少 `When` 描述触发条件的情况
- [ ] 是否不存在名称不可读或有歧义的情况

---

## TEST-03: 使用 testify 断言和 httptest 模拟

**条款**: 禁止使用 `t.Error`/`t.Fatal` 等原始断言方式（必须使用 testify），禁止在单元测试中调用真实服务端（必须使用 httptest）。

### 正确示例

```go
import (
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGetObject_ShouldReturnContent_WhenObjectExists(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("content"))
    }))
    defer server.Close()

    client, err := obs.New("ak", "sk", server.URL)
    require.NoError(t, err)

    output, err := client.GetObject("bucket", "key")
    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

### 错误示例

```go
func TestGetObject(t *testing.T) {
    // 错误：直接使用真实服务端，未使用 httptest
    client, _ := obs.New("ak", "sk", "https://real-server.com")

    output, err := client.GetObject("bucket", "key")
    if err != nil {
        t.Error("expected no error")  // 错误：应使用 testify 断言
    }
}
```

### 验证清单
- [ ] 是否不存在使用 `t.Error`/`t.Fatal` 等原始断言的情况
- [ ] 是否不存在 HTTP 模拟未使用 `httptest.NewServer` 的情况
- [ ] 是否不存在未使用 `testify/assert` 和 `testify/require` 进行断言的情况
- [ ] 是否不存在关键前置条件未使用 `require`（失败即停止）的情况

---

## TEST-04: 测试数据集中定义

**条款**: 禁止在测试函数内部硬编码散落的测试数据（凭据、endpoint、bucket 名等），测试数据必须集中定义为包级变量或辅助结构体。

### 正确示例

```go
// testdata.go 或测试文件顶部
const (
    testBucketName = "test-bucket"
    testObjectKey  = "test-object"
    testEndpoint   = "https://obs.example.com"
)

func createTestClient(serverURL string) (*obs.ObsClient, error) {
    return obs.New("test-ak", "test-sk", serverURL)
}
```

### 错误示例

```go
func TestCreateBucket(t *testing.T) {
    client, _ := obs.New("ak", "sk", "https://obs.example.com")  // 错误：硬编码散落
    client.CreateBucket(&obs.CreateBucketInput{Bucket: "my-test-bucket-123"})  // 错误：硬编码
}
```

### 验证清单
- [ ] 是否不存在在测试函数内部硬编码散落测试数据的情况
- [ ] 是否不存在辅助函数未独立于测试用例的情况
- [ ] 是否不存在在测试函数中硬编码凭据和配置的情况

---

## TEST-05: 合并重复测试用例

**条款**: 禁止为相同场景创建多个独立的测试函数，应通过 `t.Run` 子测试合并为单一测试函数。

### 正确示例

```go
func TestDeleteObject_ShouldHandleVariousStates(t *testing.T) {
    tests := []struct {
        name      string
        setup     func()
        expectErr bool
    }{
        {"existing_object", func() { /* create object */ }, false},
        {"nonexistent_object", func() { /* no setup */ }, false},
        {"already_deleted", func() { /* create then delete */ }, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            err := doDelete()
            if tt.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 错误示例

```go
// 错误：三个测试用例测试相同逻辑，应合并
func TestDeleteObject_Existing(t *testing.T) { ... }
func TestDeleteObject_NonExistent(t *testing.T) { ... }
func TestDeleteObject_AlreadyDeleted(t *testing.T) { ... }
```

### 验证清单
- [ ] 是否不存在测试相同逻辑的多个独立测试函数
- [ ] 是否不存在未通过 `t.Run` 子测试组织相关场景的情况
- [ ] 是否不存在测试函数未覆盖完整场景的情况

---

## TEST-06: 集成测试按领域分子目录

**条款**: 禁止在 `obs/test/integration/` 根目录下直接放置测试文件（必须按领域分子目录），禁止将不同领域的测试放在同一文件中。

### 正确示例

```
obs/test/integration/
  bucket/
    bucket_basic_integration_test.go
    bucket_policy_integration_test.go
  object/
    upload_download_integration_test.go
    multipart_integration_test.go
  extension/
    callback_integration_test.go
```

### 错误示例

```
obs/test/integration/
  all_tests.go              # 错误：未按领域分目录
  test1_test.go
  test2_test.go
```

### 验证清单
- [ ] 是否不存在在 `obs/test/integration/` 根目录下直接放置测试文件的情况
- [ ] 是否不存在未按 `bucket/`、`object/`、`extension/` 等领域分目录的情况
- [ ] 是否不存在文件名缺少 `integration_test` 后缀的情况
