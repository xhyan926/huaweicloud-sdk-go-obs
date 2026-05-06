# API Rules

约束 OBS SDK 客户端方法签名、参数校验、返回语义和参数选择，确保 API 风格统一。

---

## API-01: 客户端方法使用值接收者

**条款**: 禁止对 `ObsClient` 方法使用指针接收者 `(obsClient *ObsClient)`。

### 正确示例

```go
func (obsClient ObsClient) DeleteBucketCustomDomain(input *DeleteBucketCustomDomainInput, extensions ...extensionOptions) (output *BaseModel, err error) {
    // ...
}
```

### 错误示例

```go
func (obsClient *ObsClient) DeleteBucketCustomDomain(input *DeleteBucketCustomDomainInput) (*BaseModel, error) {
    // 错误：不应使用指针接收者
}
```

### 验证清单
- [ ] 是否不存在 `ObsClient` 方法使用指针接收者 `(obsClient *ObsClient)` 的情况
- [ ] 是否所有 `ObsClient` 方法均使用值接收者 `(obsClient ObsClient)`

---

## API-02: 方法签名统一格式

**条款**: 禁止使用匿名返回值（如 `(*BaseModel, error)`），禁止使用显式 `return x, y` 语句（必须使用 named return + bare return）。

### 正确示例

```go
func (obsClient ObsClient) SetBucketCustomDomain(input *SetBucketCustomDomainInput, extensions ...extensionOptions) (output *BaseModel, err error) {
    if input == nil {
        return nil, errors.New("SetBucketCustomDomainInput is nil")
    }
    output = &BaseModel{}
    err = obsClient.doActionWithBucket("SetBucketCustomDomain", HTTP_PUT, input.Bucket, input, output, extensions)
    if err != nil {
        output = nil
    }
    return
}
```

### 错误示例

```go
// 错误：未使用 named return，使用显式返回
func (obsClient ObsClient) SetBucketCustomDomain(input *SetBucketCustomDomainInput, extensions ...extensionOptions) (*BaseModel, error) {
    output := &BaseModel{}
    err := obsClient.doActionWithBucket(...)
    if err != nil {
        return nil, err  // 错误：应使用 named return + bare return
    }
    return output, nil
}
```

### 验证清单
- [ ] 是否不存在使用匿名返回值的方法签名
- [ ] 是否不存在方法末尾使用显式 `return x, y` 的情况
- [ ] 是否不存在第二个参数非 `extensions ...extensionOptions` 的情况

---

## API-03: input 为 nil 时立即返回错误

**条款**: 禁止在未检查 `input == nil` 的情况下直接访问 input 的字段，禁止 nil 检查的错误消息使用非 `"{TypeName} is nil"` 格式。

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
func (obsClient ObsClient) DeleteBucketCustomDomain(input *DeleteBucketCustomDomainInput, extensions ...extensionOptions) (output *BaseModel, err error) {
    // 错误：未检查 input == nil，直接使用 input.Bucket
    output = &BaseModel{}
    err = obsClient.doActionWithBucket("DeleteBucketCustomDomain", HTTP_DELETE, input.Bucket, ...)
}
```

### 验证清单
- [ ] 是否不存在接受 `*Input` 参数但未在开头检查 `input == nil` 的方法
- [ ] 是否不存在 nil 检查错误消息使用非 `"{MethodName}Input is nil"` 格式的情况
- [ ] 是否不存在 nil 检查后未立即返回的情况

---

## API-04: 错误时 output 置 nil，成功时 output 非 nil

**条款**: 禁止在 `doAction` 返回错误时未将 `output` 置为 `nil`，禁止方法末尾使用显式 `return x, y` 而非 bare return。

### 正确示例

```go
func (obsClient ObsClient) GetBucketCustomDomain(bucketName string, extensions ...extensionOptions) (output *GetBucketCustomDomainOutput, err error) {
    output = &GetBucketCustomDomainOutput{}
    err = obsClient.doActionWithBucket("GetBucketCustomDomain", HTTP_GET, bucketName, newSubResourceSerial(SubResourceCustomDomain), output, extensions)
    if err != nil {
        output = nil
    }
    return
}
```

### 错误示例

```go
func (obsClient ObsClient) GetBucketCustomDomain(bucketName string, extensions ...extensionOptions) (output *GetBucketCustomDomainOutput, err error) {
    // 错误：未初始化 output
    err = obsClient.doActionWithBucket(...)
    if err != nil {
        return nil, err  // 错误：应使用 output = nil + bare return
    }
    output = &GetBucketCustomDomainOutput{}
    return
}
```

### 验证清单
- [ ] 是否不存在 `doAction` 返回错误时未将 `output` 置为 `nil` 的情况
- [ ] 是否不存在方法末尾使用显式 `return x, y` 而非 bare return 的情况
- [ ] 是否不存在方法开头未初始化 `output = &XxxOutput{}` 的情况

---

## API-05: 只读操作可使用 string 参数，复杂操作使用 input 结构体

**条款**: 禁止为仅需单个 `string` 参数的只读操作创建 `Input` 结构体，禁止需要多个参数的写入操作使用裸 `string` 参数。

### 正确示例

```go
// 只读操作 — 直接使用 string 参数
func (obsClient ObsClient) GetBucketCustomDomain(bucketName string, extensions ...extensionOptions) (output *GetBucketCustomDomainOutput, err error) {
    output = &GetBucketCustomDomainOutput{}
    err = obsClient.doActionWithBucket("GetBucketCustomDomain", HTTP_GET, bucketName, newSubResourceSerial(SubResourceCustomDomain), output, extensions)
    if err != nil {
        output = nil
    }
    return
}

// 写入操作 — 使用 Input 结构体
func (obsClient ObsClient) DeleteBucketCustomDomain(input *DeleteBucketCustomDomainInput, extensions ...extensionOptions) (output *BaseModel, err error) {
    if input == nil {
        return nil, errors.New("DeleteBucketCustomDomainInput is nil")
    }
    output = &BaseModel{}
    err = obsClient.doActionWithBucket("DeleteBucketCustomDomain", HTTP_DELETE, input.Bucket, newSubResourceSerialV2(SubResourceCustomDomain, input.CustomDomain), output, extensions)
    if err != nil {
        output = nil
    }
    return
}
```

### 错误示例

```go
// 错误：简单的只读操作不应定义 Input 结构体
type GetBucketCustomDomainInput struct {
    Bucket string
}
func (obsClient ObsClient) GetBucketCustomDomain(input *GetBucketCustomDomainInput, ...) { ... }

// 错误：写入操作需要额外参数却用 string
func (obsClient ObsClient) DeleteBucketCustomDomain(bucketName string, customDomain string, ...) { ... }
```

### 验证清单
- [ ] 是否不存在为仅需单个 `string` 参数的只读操作创建 `Input` 结构体的情况
- [ ] 是否不存在需要多个参数的写入操作使用裸 `string` 参数的情况
- [ ] 是否不存在简单只读操作（仅需要 bucketName）未使用 `string` 参数的情况
