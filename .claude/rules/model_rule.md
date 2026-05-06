# Model Rules

约束 OBS SDK 中结构体定义、BaseModel 嵌入、序列化标签和 OBS/S3 双版本模式。

---

## MODEL-01: 所有 Output 结构体必须嵌入 BaseModel

**条款**: 禁止定义未嵌入 `BaseModel` 的 Output 结构体，禁止在 Output 中重复定义 `BaseModel` 已有的字段（如 `RequestId`、`StatusCode`）。

### 正确示例

```go
type GetBucketCustomDomainOutput struct {
    BaseModel
    Domains []Domain `xml:"Domains"`
}

type GetBucketMirrorBackToSourceOutput struct {
    BaseModel
    Rules string `json:"body"`
}
```

### 错误示例

```go
type GetBucketCustomDomainOutput struct {
    // 错误：缺少 BaseModel 嵌入
    Domains []Domain `xml:"Domains"`
    RequestId string  // 错误：不应手动定义 BaseModel 中已有的字段
}
```

### 验证清单
- [ ] 是否不存在未嵌入 `BaseModel` 的 `*Output` 结构体
- [ ] 是否不存在 `BaseModel` 非结构体第一个字段的情况
- [ ] 是否不存在在 Output 中重复定义 `BaseModel` 已有字段的情况

---

## MODEL-02: HTTP 头字段标注 xml:"-"

**条款**: 禁止 HTTP 头传递的字段缺少 `xml:"-"` 标签（会导致错误的 XML 序列化），禁止 XML 正文字段标注 `xml:"-"`。

### 正确示例

```go
type ObjectMetadata struct {
    CacheControl       string `xml:"-" json:"-"`
    ContentDisposition string `xml:"-" json:"-"`
    ContentEncoding    string `xml:"-" json:"-"`
    ContentType        string `xml:"-" json:"-"`
}
```

### 错误示例

```go
type ObjectMetadata struct {
    CacheControl string `xml:"CacheControl"` // 错误：HTTP 头字段不应参与 XML 序列化
}
```

### 验证清单
- [ ] 是否不存在 HTTP 头传递的字段缺少 `xml:"-"` 标签的情况
- [ ] 是否不存在 XML 正文字段标注 `xml:"-"` 的情况
- [ ] 是否不存在未区分头部字段和正文字段的情况

---

## MODEL-03: 序列化标签规范

**条款**: 禁止 XML 标签使用非 PascalCase（如 `xml:"name"`），禁止在必填字段上使用 `,omitempty`，禁止 JSON 标签使用非 camelCase。

### 正确示例

```go
type CustomDomainConfiguration struct {
    Name             string `xml:"Name"`
    CertificateId    string `xml:"CertificateId,omitempty"` // 可选字段
    Certificate      string `xml:"Certificate"`
    CertificateChain string `xml:"CertificateChain,omitempty"`
    PrivateKey       string `xml:"PrivateKey"`
}

type SetBucketMirrorBackToSourceInput struct {
    Bucket string
    Rules  string `json:"body"`
}
```

### 错误示例

```go
type CustomDomainConfiguration struct {
    Name          string `xml:"name"`               // 错误：XML 标签应使用 PascalCase
    CertificateId string `xml:"CertificateId"`       // 错误：可选字段缺少 omitempty
}
```

### 验证清单
- [ ] 是否不存在 XML 标签使用非 PascalCase 的情况
- [ ] 是否不存在必填字段使用 `,omitempty` 的情况
- [ ] 是否不存在 JSON 标签使用非 camelCase 的情况
- [ ] 是否不存在可选字段缺少 `,omitempty` 的情况

---

## MODEL-04: Input 的 trans 方法负责字段转换

**条款**: 禁止在 Client 层中处理 OBS/S3 字段差异（如判断 `signature == SignatureObs`），字段转换逻辑必须限制在 trait 层的 `trans` 方法中。

### 正确示例

```go
// trait 层中为 Input 定义的 trans 方法
func (input SetBucketStoragePolicyInput) trans(isObs bool) bucketStoragePolicy {
    policy := bucketStoragePolicy{}
    if isObs {
        policy.StorageClass = input.StorageClass
    } else {
        policy.StorageClass = input.StorageClass
    }
    return policy
}
```

### 错误示例

```go
// 错误：在 Client 层直接做字段转换
func (obsClient ObsClient) SetBucketStoragePolicy(input *SetBucketStoragePolicyInput, ...) {
    if obsClient.conf.signature == SignatureObs {
        // 不应在 Client 层处理 OBS/S3 差异
    }
}
```

### 验证清单
- [ ] 是否不存在字段转换逻辑位于非 trait 层 `trans` 方法中的情况
- [ ] 是否不存在 Client 层包含 OBS/S3 字段差异处理逻辑的情况
- [ ] 是否不存在 `trans` 方法签名非 `func (input XxxInput) trans(isObs bool) yyyType` 的情况

---

## MODEL-05: OBS/S3 双版本结构使用小写未导出类型

**条款**: 禁止导出 OBS/S3 双版本结构体（如 `BucketStoragePolicyS3`），禁止在 trait 层之外引用双版本结构。

### 正确示例

```go
// trait_bucket.go — 未导出的双版本结构
type bucketStoragePolicy struct {
    XMLName       xml.Name `xml:"StoragePolicy"`
    StorageClass  string   `xml:"DefaultStorageClass"`
}

type bucketStoragePolicyObs struct {
    XMLName       xml.Name `xml:"StoragePolicy"`
    StorageClass  string   `xml:"DefaultStorageClass"`
}
```

### 错误示例

```go
// 错误：导出的双版本结构，暴露内部差异
type BucketStoragePolicyS3 struct { ... }
type BucketStoragePolicyOBS struct { ... }
```

### 验证清单
- [ ] 是否不存在 OBS/S3 双版本结构为导出类型（大写开头）的情况
- [ ] 是否不存在双版本结构定义在 trait 层之外的情况
- [ ] 是否不存在导出的 OBS/S3 特定版本结构体
