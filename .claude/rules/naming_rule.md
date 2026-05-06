# Naming Rules

约束 OBS SDK 中类型、函数、变量、常量和文件的命名规范，确保命名风格统一且符合 Go 惯例。

---

## NAME-01: 导出类型命名以功能后缀结尾

**条款**: 禁止对 API 输入/输出结构体使用 `Input`/`Output` 以外的后缀（如 `Request`、`Response`、`Dto`）。

### 正确示例

```go
type DeleteBucketCustomDomainInput struct {
    Bucket       string
    CustomDomain string
}

type GetBucketCustomDomainOutput struct {
    BaseModel
    Domains []Domain `xml:"Domains"`
}
```

### 错误示例

```go
type DeleteBucketRequest struct { // 错误：应使用 Input 后缀
    Bucket string
}

type GetBucketResponse struct { // 错误：应使用 Output 后缀
    Domains []Domain `xml:"Domains"`
}
```

### 验证清单
- [ ] 是否不存在输入结构体不以 `Input` 结尾的情况
- [ ] 是否不存在输出结构体不以 `Output` 结尾的情况
- [ ] 是否不存在使用 `Request`/`Response` 等非标准后缀的情况

---

## NAME-02: 类型常量值全大写，类型名以 Type 结尾

**条款**: 禁止枚举类型缺少 `Type` 后缀（如用 `Signature` 而非 `SignatureType`），禁止使用缩写常量名（如用 `SigV2` 而非 `SignatureV2`）。

### 正确示例

```go
// type.go
type SignatureType string

const (
    SignatureV2  SignatureType = "v2"
    SignatureV4  SignatureType = "v4"
    SignatureObs SignatureType = "OBS"
)

type HttpMethodType string

const (
    HttpMethodGet    HttpMethodType = HTTP_GET
    HttpMethodPut    HttpMethodType = HTTP_PUT
)
```

### 错误示例

```go
type Signature string // 错误：缺少 Type 后缀

const (
    SigV2  Signature = "v2"  // 错误：常量名应使用全称 SignatureV2
)
```

### 验证清单
- [ ] 是否不存在枚举类型缺少 `Type` 后缀的情况
- [ ] 是否不存在使用缩写常量名（如 `SigV2`）的情况
- [ ] 是否不存在常量名与变量名语义不一致的情况

---

## NAME-03: 配置选项函数以 With 开头

**条款**: 禁止使用 `Set`、`Get`、`Config` 等非 `With` 前缀命名配置选项函数。

### 正确示例

```go
// conf.go 中的配置选项
func WithRegion(region string) configurer {
    return func(conf *config) {
        conf.region = region
    }
}

func WithSignature(signature SignatureType) configurer {
    return func(conf *config) {
        conf.signature = signature
    }
}

func WithConnectTimeout(connectTimeout int) configurer {
    return func(conf *config) {
        conf.connectTimeout = connectTimeout
    }
}
```

### 错误示例

```go
func SetRegion(region string) configurer { ... }     // 错误：应使用 With 前缀
func Region(region string) configurer { ... }         // 错误：缺少 With 前缀
```

### 验证清单
- [ ] 是否不存在配置选项函数不以 `With` 开头的情况
- [ ] 是否不存在函数名不准确描述配置属性的情况（如 `WithLocation` 而非 `WithRegion`）
- [ ] 是否不存在返回类型非 `configurer` 的配置选项函数

---

## NAME-04: 扩展选项函数以 With 开头

**条款**: 禁止使用 `Add`、`Set`、`Register` 等非 `With` 前缀命名扩展选项函数。

### 正确示例

```go
// extension.go
func WithProgress(progressListener ProgressListener) extensionProgressListener { ... }
func WithReqPaymentHeader(requester PayerType) extensionHeaders { ... }
func WithTrafficLimitHeader(trafficLimit int64) extensionHeaders { ... }
func WithCallbackHeader(callback string) extensionHeaders { ... }
func WithCustomHeader(key string, value string) extensionHeaders { ... }
```

### 错误示例

```go
func AddProgress(l ProgressListener) extensionProgressListener { ... }    // 错误：应使用 With 前缀
func SetTrafficLimit(limit int64) extensionHeaders { ... }                 // 错误：应使用 With 前缀
```

### 验证清单
- [ ] 是否不存在扩展选项函数不以 `With` 开头的情况
- [ ] 是否不存在涉及 HTTP 头的扩展不以 `Header` 结尾的情况（如 `WithTrafficLimitHeader`）
- [ ] 是否不存在返回类型非 `extensionHeaders` 或 `extensionProgressListener` 的扩展函数

---

## NAME-05: 导出标识符用 PascalCase，未导出用 camelCase

**条款**: 禁止在标识符命名中使用下划线（如 `obs_client`、`do_action`），导出标识符禁止使用 camelCase，未导出标识符禁止使用 PascalCase。

### 正确示例

```go
// 导出
type ObsClient struct { ... }
func New(ak, sk, endpoint string, configurers ...configurer) (*ObsClient, error) { ... }
const SignatureV2 SignatureType = "v2"

// 未导出
type config struct { ... }
type securityHolder struct { ... }
func setHeaders(headers map[string][]string, key string, value []string, isObs bool) { ... }
```

### 错误示例

```go
type obs_client struct { ... }        // 错误：应使用 PascalCase
func doAction_withBucket() { ... }    // 错误：应使用 camelCase
```

### 验证清单
- [ ] 是否不存在导出标识符未使用 PascalCase（含下划线）的情况
- [ ] 是否不存在未导出函数未使用 camelCase 的情况
- [ ] 是否不存在缩写词未全大写的情况（如 `OBS`、`AK`、`SK`、`HTTP`）

---

## NAME-06: HTTP 头常量用 HEADER_ 前缀，子资源用 SubResource 前缀

**条款**: 禁止 HTTP 头常量缺少 `HEADER_` 前缀，禁止子资源常量缺少 `SubResource` 前缀，禁止 HTTP 方法常量缺少 `HTTP_` 前缀。

### 正确示例

```go
// const.go
const (
    HEADER_PREFIX           = "x-amz-"
    HEADER_ACL              = "acl"
    HEADER_STORAGE_CLASS    = "x-default-storage-class"
    HEADER_STORAGE_CLASS_OBS = "x-obs-storage-class"
)

// type.go
const (
    SubResourceStoragePolicy SubResourceType = "storagePolicy"
    SubResourceAcl           SubResourceType = "acl"
)
```

### 错误示例

```go
const ACL_HEADER = "acl"                    // 错误：应使用 HEADER_ 前缀
const StoragePolicyResource = "storagePolicy" // 错误：应使用 SubResource 前缀
```

### 验证清单
- [ ] 是否不存在 HTTP 头常量缺少 `HEADER_` 前缀的情况
- [ ] 是否不存在子资源常量缺少 `SubResource` 前缀的情况
- [ ] 是否不存在 HTTP 方法常量缺少 `HTTP_` 前缀的情况
