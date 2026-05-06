# Architecture Rules

约束 OBS SDK 的分层架构、文件职责和依赖方向，确保代码组织清晰、调用链单向。

---

## ARCH-01: 四层调用方向单向向下

**条款**: 调用链必须遵循 Client → Trait → HTTP → Model 的单向方向，严禁反向调用或跨层调用。客户端方法调用 trait 层，trait 层调用 HTTP 层（`doAction` 系列），HTTP 层操作 Model 进行序列化/反序列化。

### 正确示例

```go
// client_bucket.go — Client 层
func (obsClient ObsClient) SetBucketStoragePolicy(input *SetBucketStoragePolicyInput, extensions ...extensionOptions) (output *BaseModel, err error) {
    // Client 层调用 Trait 层实现
    output = &BaseModel{}
    err = obsClient.doActionWithBucket("SetBucketStoragePolicy", HTTP_PUT, input.Bucket, input, output, extensions)
    if err != nil {
        output = nil
    }
    return
}
```

### 错误示例

```go
// 错误：Trait 层直接调用 Client 层方法
func (obsClient ObsClient) setBucketStoragePolicyBase(input *SetBucketStoragePolicyInput) {
    // 不能在 trait 层调用其他 client 方法
    obsClient.SetBucketStoragePolicy(input)
}
```

### 验证清单
- [ ] Client 层方法是否仅通过 `doAction` 系列方法委托给 HTTP 层
- [ ] HTTP 层是否未直接引用 Client 层的公共方法
- [ ] Model 层是否仅被 HTTP/Trait 层引用，未持有上层逻辑

---

## ARCH-02: 客户端方法与特性实现分离

**条款**: 公共 API 方法定义在 `client_*.go` 文件中，内部特性实现定义在 `trait_*.go` 文件中。`client_*.go` 文件禁止引用 `trait_*.go` 中的未导出函数，`trait_*.go` 禁止导出函数。

### 正确示例

```go
// client_bucket.go — 公共 API 方法
func (obsClient ObsClient) GetBucketStoragePolicy(input *GetBucketStoragePolicyInput, extensions ...extensionOptions) (output *GetBucketStoragePolicyOutput, err error) {
    output = &GetBucketStoragePolicyOutput{}
    err = obsClient.doActionWithBucket("GetBucketStoragePolicy", HTTP_GET, input.Bucket, newSubResourceSerial(SubResourceStoragePolicy), output, extensions)
    if err != nil {
        output = nil
    }
    return
}

// trait_bucket.go — 内部特性实现（未导出）
func (obsClient ObsClient) getBucketStoragePolicyBase(input *GetBucketStoragePolicyInput, output IGetBucketStoragePolicyOutput) error {
    return obsClient.doActionWithBucket("GetBucketStoragePolicy", HTTP_GET, input.Bucket, newSubResourceSerial(SubResourceStoragePolicy), output)
}
```

### 错误示例

```go
// 错误：在 client_bucket.go 中定义未导出的特性实现
func (obsClient ObsClient) getInternalPolicy() { ... }
```

### 验证清单
- [ ] `client_*.go` 文件中是否仅包含导出的公共方法
- [ ] `trait_*.go` 文件中的函数是否均为未导出
- [ ] 是否存在反向引用（trait 调用 client 的导出方法）

---

## ARCH-03: 公共类型和常量集中定义

**条款**: 公共类型集中定义在 `type.go`，公共常量集中定义在 `const.go`。禁止在 `client_*.go`、`trait_*.go`、`model_*.go` 中定义导出类型或导出常量。

### 正确示例

```go
// type.go — 集中定义导出类型
type SignatureType string

const (
    SignatureV2  SignatureType = "v2"
    SignatureV4  SignatureType = "v4"
    SignatureObs SignatureType = "OBS"
)
```

### 错误示例

```go
// 错误：在 client_bucket.go 中定义导出类型
type BucketOperation string
```

### 验证清单
- [ ] 所有导出类型是否定义在 `type.go` 中
- [ ] 所有导出常量是否定义在 `const.go` 中
- [ ] 业务文件中是否只包含未导出的类型和常量

---

## ARCH-04: HTTP 请求统一通过 doAction 入口

**条款**: 所有 HTTP 请求必须通过 `doAction` 系列方法（`doActionWithBucket`、`doActionWithBucketV2`、`doActionForObject` 等）统一入口，禁止直接使用 `http.Client` 发送请求。

### 正确示例

```go
// 通过 doActionWithBucket 统一入口
err = obsClient.doActionWithBucket("SetBucketStoragePolicy", HTTP_PUT, input.Bucket, input, output, extensions)
```

### 错误示例

```go
// 错误：绕过 doAction 直接发送 HTTP 请求
req, _ := http.NewRequest("PUT", url, body)
resp, _ := http.DefaultClient.Do(req)
```

### 验证清单
- [ ] 是否所有 HTTP 请求都通过 `doAction` 系列方法
- [ ] 是否存在直接使用 `http.Client` 或 `http.NewRequest` 的业务代码
- [ ] `doAction` 是否统一处理签名、重试、错误解析

---

## ARCH-05: 文件按领域划分

**条款**: 禁止创建未按领域命名的源文件（如 `utils.go`、`helpers.go`、`api.go`），禁止将公共模型定义在非 `model_*.go` 文件中。

### 正确示例

```
obs/
  client_bucket.go    # 存储桶相关公共 API
  client_object.go    # 对象相关公共 API
  client_part.go      # 分块相关公共 API
  trait_bucket.go     # 存储桶内部实现
  trait_object.go     # 对象内部实现
  model_bucket.go     # 存储桶模型定义
  model_object.go     # 对象模型定义
```

### 错误示例

```
obs/
  api.go              # 错误：未按领域命名
  bucket_utils.go     # 错误：应使用 trait_ 前缀
  helpers.go          # 错误：未按领域划分
```

### 验证清单
- [ ] 是否不存在未按 `*_bucket.go` / `*_object.go` / `*_part.go` 命名的源文件
- [ ] 是否不存在未按领域划分的通用文件（如 `utils.go`、`helpers.go`）
- [ ] 是否不存在公共模型定义在非 `model_*.go` 文件中的情况
