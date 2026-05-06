# Error Rules

约束 OBS SDK 中的错误类型选择、错误消息格式和日志记录规范。

---

## ERR-01: 服务端错误使用 ObsError 结构体

**条款**: 禁止使用 `fmt.Errorf` 或自定义错误类型表示 OBS 服务端错误，服务端错误必须统一解析为 `ObsError`。

### 正确示例

```go
// error.go — ObsError 定义
type ObsError struct {
    BaseModel
    Status    string
    XMLName   xml.Name `xml:"Error"`
    Code      string   `xml:"Code" json:"code"`
    Message   string   `xml:"Message" json:"message"`
    Resource  string   `xml:"Resource"`
    HostId    string   `xml:"HostId"`
    Indicator string
}

// 调用方使用
output, err := obsClient.GetBucketCustomDomain("my-bucket")
if err != nil {
    if obsError, ok := err.(obs.ObsError); ok {
        fmt.Printf("Status=%s, Code=%s, Message=%s\n", obsError.Status, obsError.Code, obsError.Message)
    }
}
```

### 错误示例

```go
// 错误：手动构造服务端错误
if resp.StatusCode >= 400 {
    return fmt.Errorf("request failed with status %d", resp.StatusCode)  // 应使用 ObsError
}
```

### 验证清单
- [ ] 是否不存在使用 `fmt.Errorf` 或自定义错误类型表示服务端错误的情况
- [ ] 是否不存在 `ObsError` 缺少 `Status`、`Code`、`Message`、`RequestId` 字段的情况
- [ ] 是否不存在未通过类型断言访问 `ObsError` 详细字段的情况

---

## ERR-02: 客户端校验错误使用标准 errors.New()

**条款**: 禁止对客户端参数校验错误使用 `ObsError`，客户端校验错误必须使用标准库 `errors.New()`。

### 正确示例

```go
func (obsClient ObsClient) DeleteBucketCustomDomain(input *DeleteBucketCustomDomainInput, extensions ...extensionOptions) (output *BaseModel, err error) {
    if input == nil {
        return nil, errors.New("DeleteBucketCustomDomainInput is nil")
    }
    // ...
}
```

### 错误示例

```go
// 错误：客户端校验使用 ObsError
if input == nil {
    return nil, ObsError{Code: "InvalidInput", Message: "input is nil"}
}
```

### 验证清单
- [ ] 是否不存在客户端参数校验错误使用 `ObsError` 的情况
- [ ] 是否不存在参数校验错误未使用标准库 `errors.New()` 的情况
- [ ] 是否不存在 `errors` 包未从标准库导入的情况

---

## ERR-03: 日志记录使用 doLog 函数

**条款**: 日志记录使用 `doLog(LEVEL_ERROR/WARN, ...)` 函数，禁止使用 `fmt.Println`、`log.Println` 等标准输出方式进行日志输出。

### 正确示例

```go
doLog(LEVEL_WARN, "Try to get security from ecs failed, cost %d ms, err %s", cost, err.Error())
doLog(LEVEL_ERROR, "Failed to parse response: %s", err.Error())
doLog(LEVEL_INFO, "Get security from ecs succeed, AK:xxxx, SK:xxxx")
doLog(LEVEL_DEBUG, "Get the json data from ecs succeed")
```

### 错误示例

```go
fmt.Println("Error:", err)           // 错误：禁止使用 fmt.Println
log.Printf("Failed: %s", err)        // 错误：禁止使用标准 log
fmt.Fprintf(os.Stderr, "err: %v", e) // 错误：禁止直接输出到 stderr
```

### 验证清单
- [ ] 所有日志是否通过 `doLog` 函数记录
- [ ] 是否不存在 `fmt.Println`、`log.Println` 等直接输出
- [ ] 日志级别是否选择正确（DEBUG/INFO/WARN/ERROR）

---

## ERR-04: 参数校验错误消息保持一致

**条款**: 禁止使用与现有代码不一致的错误消息格式（如 `"input cannot be null"`、`"bucket name is required"`），必须使用：
- `"XxxInput is nil"` — input 为 nil
- `"Bucket is empty"` — bucket 名为空
- `"Key is empty"` — object key 为空

### 正确示例

```go
return nil, errors.New("DeleteBucketCustomDomainInput is nil")
return nil, errors.New("Bucket is empty")
return nil, errors.New("Key is empty")
```

### 错误示例

```go
return nil, errors.New("input cannot be null")       // 错误：格式应为 "XxxInput is nil"
return nil, errors.New("bucket name is required")    // 错误：格式应为 "Bucket is empty"
return nil, errors.New("object key not provided")    // 错误：格式应为 "Key is empty"
```

### 验证清单
- [ ] 是否不存在 nil 检查的错误消息使用非 `"{TypeName} is nil"` 格式的情况
- [ ] 是否不存在必填字段为空的错误消息使用非 `"{FieldName} is empty"` 格式的情况
- [ ] 是否不存在新增错误消息与现有代码风格不一致的情况
