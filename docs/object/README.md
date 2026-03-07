# 对象上传 API 文档

本文档包含 OBS SDK 中对象上传相关的所有 API 接口说明。

## 目录

- [CreatePostPolicy](#createpostpolicy)
- [辅助函数](#辅助函数)
- [数据结构](#数据结构)
- [常量定义](#常量定义)
- [使用场景](#使用场景)

---

## CreatePostPolicy

创建 POST 上传策略，用于直接从浏览器向 OBS 上传文件。

### 方法签名

```go
func (obsClient ObsClient) CreatePostPolicy(input *CreatePostPolicyInput) (output *CreatePostPolicyOutput, err error)
```

### 功能说明

CreatePostPolicy 方法生成 POST 上传策略，包括 Policy JSON、签名和完整的 Token。生成的 Token 可以用于前端 HTML 表单，实现文件直接从浏览器上传到 OBS，无需经过后端服务器。

**优势**：
- 降低服务器负载和网络带宽
- 支持大文件上传
- 上传速度更快
- 文件不经过后端服务器

### 参数说明

**CreatePostPolicyInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名称 |
| Key | string | 是 | 对象名（上传后的文件路径） |
| Expires | int64 | 否 | 过期时间（Unix 时间戳，秒） |
| ExpiresIn | int64 | 否 | 过期时长（秒），与 Expires 二选一 |
| Acl | string | 否 | 对象 ACL，默认为私有 |
| Conditions | []PostPolicyCondition | 否 | 自定义条件列表 |

### 返回值

**CreatePostPolicyOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| Policy | string | Base64 编码的 Policy JSON |
| Signature | string | Policy 的 HMAC-SHA1 签名 |
| Token | string | 完整的上传 Token，格式为 `ak:signature:policy` |
| AccessKeyId | string | Access Key ID |

### 使用示例

```go
package main

import (
    "fmt"
    obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func main() {
    // 1. 创建 OBS 客户端
    obsClient, err := obs.New("your-ak", "your-sk", "your-endpoint")
    if err != nil {
        panic(err)
    }

    // 2. 创建 POST 上传策略
    input := &obs.CreatePostPolicyInput{
        Bucket:    "my-bucket",
        Key:       "uploads/test.jpg",
        ExpiresIn: 3600, // 1 小时后过期
        Acl:       "public-read",
    }

    output, err := obsClient.CreatePostPolicy(input)
    if err != nil {
        fmt.Printf("创建策略失败: %v\n", err)
        return
    }

    // 3. 输出策略信息
    fmt.Printf("Access Key ID: %s\n", output.AccessKeyId)
    fmt.Printf("Policy: %s\n", output.Policy)
    fmt.Printf("Signature: %s\n", output.Signature)
    fmt.Printf("Token: %s\n", output.Token)

    // 4. 使用 Token 生成前端 HTML 表单
    // 参见 main/post_upload_sample.go
}
```

### 前端使用示例

生成的 Token 可以用于前端 HTML 表单：

```html
<form action="https://obs.cn-north-1.myhuaweicloud.com/my-bucket/" method="post" enctype="multipart/form-data">
    <input type="text" name="AWSAccessKeyId" value="your-access-key-id" readonly>
    <input type="text" name="policy" value="base64-encoded-policy" readonly>
    <input type="text" name="signature" value="signature" readonly>
    <input type="text" name="key" value="uploads/test.jpg" readonly>
    <input type="file" name="file">
    <button type="submit">上传到 OBS</button>
</form>
```

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidArgument | 400 | 参数错误，如桶名或键名为空 |
| AccessDenied | 403 | 权限不足，检查 AK/SK 是否正确 |

### 注意事项

1. **过期时间**：Policy 有过期时间，过期后需要重新生成
2. **安全性**：Policy 和签名由后端生成，前端只使用生成的凭证
3. **条件限制**：可以在 Conditions 中添加文件大小、文件类型等限制
4. **默认条件**：系统会自动添加桶和键条件，无需手动添加

---

## 辅助函数

### BuildPostPolicyExpiration

生成 Policy 的过期时间字符串（ISO 8601 格式）。

#### 方法签名

```go
func BuildPostPolicyExpiration(expiresIn int64) string
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| expiresIn | int64 | 是 | 过期时长（秒） |

#### 返回值

返回格式为 `2006-01-02T15:04:05.000Z` 的 ISO 8601 时间字符串。

#### 示例

```go
expiresIn := int64(3600) // 1 小时
expiration := obs.BuildPostPolicyExpiration(expiresIn)
// 输出: "2026-03-06T14:30:00.000Z"
```

---

### ValidatePostPolicy

验证 POST Policy 的有效性。

#### 方法签名

```go
func ValidatePostPolicy(policy *PostPolicy) error
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| policy | *PostPolicy | 是 | 要验证的 Policy |

#### 返回值

验证成功返回 nil，失败返回错误信息。

#### 验证项

- Policy 不为 nil
- 过期时间不为空
- 每个条件的 Key 和 Operator 不为空

#### 示例

```go
policy := &obs.PostPolicy{
    Expiration: "2026-03-06T14:30:00.000Z",
    Conditions: []obs.PostPolicyCondition{
        obs.CreateBucketCondition("my-bucket"),
    },
}

err := obs.ValidatePostPolicy(policy)
if err != nil {
    fmt.Printf("验证失败: %v\n", err)
}
```

---

### CalculatePostPolicySignature

计算 POST Policy 的签名。

#### 方法签名

```go
func CalculatePostPolicySignature(policyJSON, secretAccessKey string) (string, error)
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| policyJSON | string | 是 | Policy JSON 字符串 |
| secretAccessKey | string | 是 | Secret Access Key |

#### 返回值

返回 Base64 编码的 HMAC-SHA1 签名。

#### 签名流程

1. Policy JSON 进行 Base64 编码
2. 使用 HMAC-SHA1 算法对编码后的 Policy 进行签名
3. 签名结果进行 Base64 编码

#### 示例

```go
policyJSON := `{"expiration":"2026-03-06T14:30:00.000Z","conditions":[["eq","$bucket","my-bucket"]]}`
secretKey := "your-secret-key"

signature, err := obs.CalculatePostPolicySignature(policyJSON, secretKey)
if err != nil {
    fmt.Printf("签名计算失败: %v\n", err)
}

fmt.Printf("签名: %s\n", signature)
```

---

### BuildPostPolicyToken

构建完整的 POST 上传 Token。

#### 方法签名

```go
func BuildPostPolicyToken(ak, signature, policy string) string
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ak | string | 是 | Access Key ID |
| signature | string | 是 | Policy 的签名 |
| policy | string | 是 | Base64 编码的 Policy |

#### 返回值

返回格式为 `ak:signature:policy` 的 Token 字符串。

#### 示例

```go
ak := "your-access-key-id"
signature := "base64-encoded-signature"
policy := "base64-encoded-policy"

token := obs.BuildPostPolicyToken(ak, signature, policy)
// 输出: "your-access-key-id:base64-encoded-signature:base64-encoded-policy"
```

---

### CreatePostPolicyCondition

创建 POST Policy 条件。

#### 方法签名

```go
func CreatePostPolicyCondition(operator, key string, value interface{}) PostPolicyCondition
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| operator | string | 是 | 条件操作符（eq、starts-with、content-length-range） |
| key | string | 是 | 条件键（$bucket、$key、$content-type 等） |
| value | interface{} | 是 | 条件值 |

#### 返回值

返回 PostPolicyCondition 结构体。

#### 示例

```go
condition := obs.CreatePostPolicyCondition(
    obs.PostPolicyOpEquals,
    obs.PostPolicyKeyContentType,
    "image/jpeg",
)
```

---

### CreateBucketCondition

创建桶条件的便捷函数。

#### 方法签名

```go
func CreateBucketCondition(bucket string) PostPolicyCondition
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucket | string | 是 | 桶名称 |

#### 返回值

返回一个 `eq` 操作符的桶条件。

#### 示例

```go
condition := obs.CreateBucketCondition("my-bucket")
// 等价于: ["eq", "$bucket", "my-bucket"]
```

---

### CreateKeyCondition

创建键条件的便捷函数。

#### 方法签名

```go
func CreateKeyCondition(key string) PostPolicyCondition
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| key | string | 是 | 对象名或前缀 |

#### 返回值

返回一个 `starts-with` 操作符的键条件。

#### 示例

```go
condition := obs.CreateKeyCondition("uploads/")
// 等价于: ["starts-with", "$key", "uploads/"]
```

---

## 数据结构

### PostPolicy

POST 上传策略结构。

```go
type PostPolicy struct {
    Expiration string                `json:"expiration"`
    Conditions []PostPolicyCondition `json:"-"`
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| Expiration | string | 过期时间（ISO 8601 格式） |
| Conditions | []PostPolicyCondition | 条件列表 |

---

### PostPolicyCondition

POST Policy 条件结构。

```go
type PostPolicyCondition struct {
    Operator string      `json:"-"`
    Key      string      `json:"-"`
    Value    interface{} `json:"-"`
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| Operator | string | 条件操作符（eq、starts-with、content-length-range） |
| Key | string | 条件键（$bucket、$key、$content-type 等） |
| Value | interface{} | 条件值，可以是字符串、整数或数组 |

**JSON 序列化**：

条件序列化为数组格式：`["operator", "key", "value"]`

---

### CreatePostPolicyInput

创建 POST Policy 的输入参数。

```go
type CreatePostPolicyInput struct {
    Bucket     string
    Key        string
    Expires    int64
    ExpiresIn  int64
    Acl        string
    Conditions []PostPolicyCondition
}
```

**字段说明**：

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|----------|------|
| Bucket | string | 是 | - | 桶名称 |
| Key | string | 是 | - | 对象名（上传后的文件路径） |
| Expires | int64 | 否 | - | 过期时间（Unix 时间戳） |
| ExpiresIn | int64 | 否 | 3600 | 过期时长（秒） |
| Acl | string | 否 | - | 对象 ACL |
| Conditions | []PostPolicyCondition | 否 | nil | 自定义条件列表 |

---

### CreatePostPolicyOutput

创建 POST Policy 的输出结果。

```go
type CreatePostPolicyOutput struct {
    BaseModel
    Policy      string `json:"policy"`
    Signature   string `json:"signature"`
    Token       string `json:"token"`
    AccessKeyId string `json:"accessKeyId"`
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| Policy | string | Base64 编码的 Policy JSON |
| Signature | string | Policy 的签名 |
| Token | string | 完整的上传 Token（ak:signature:policy） |
| AccessKeyId | string | Access Key ID |

---

## 常量定义

### Policy 条件键

```go
const (
    PostPolicyKeyBucket         = "$bucket"
    PostPolicyKeyKey           = "$key"
    PostPolicyKeyContentType   = "$content-type"
    PostPolicyKeyContentLength = "$content-length"
)
```

**说明**：

| 常量 | 值 | 说明 |
|------|-----|------|
| PostPolicyKeyBucket | $bucket | 桶名称 |
| PostPolicyKeyKey | $key | 对象名 |
| PostPolicyKeyContentType | $content-type | 文件类型 |
| PostPolicyKeyContentLength | $content-length | 文件大小 |

---

### Policy 条件操作符

```go
const (
    PostPolicyOpEquals     = "eq"
    PostPolicyOpStartsWith = "starts-with"
    PostPolicyOpRange      = "content-length-range"
)
```

**说明**：

| 常量 | 值 | 说明 |
|------|-----|------|
| PostPolicyOpEquals | eq | 完全匹配 |
| PostPolicyOpStartsWith | starts-with | 前缀匹配 |
| PostPolicyOpRange | content-length-range | 内容长度范围 |

---

## 使用场景

### 场景 1：基本 POST 上传

创建一个简单的 POST 上传策略，只包含桶和键条件。

```go
input := &obs.CreatePostPolicyInput{
    Bucket:    "my-bucket",
    Key:       "uploads/file.jpg",
    ExpiresIn: 3600,
}
```

---

### 场景 2：带文件大小限制

创建策略，限制上传文件的大小为 0-10MB。

```go
input := &obs.CreatePostPolicyInput{
    Bucket:    "my-bucket",
    Key:       "uploads/file.jpg",
    ExpiresIn: 3600,
    Conditions: []obs.PostPolicyCondition{
        obs.CreatePostPolicyCondition(
            "content-length-range",
            "$content-length",
            []interface{}{0, 10 * 1024 * 1024},
        ),
    },
}
```

---

### 场景 3：带文件类型限制

创建策略，限制上传文件的类型为图片。

```go
input := &obs.CreatePostPolicyInput{
    Bucket:    "my-bucket",
    Key:       "uploads/file.jpg",
    ExpiresIn: 3600,
    Conditions: []obs.PostPolicyCondition{
        obs.CreatePostPolicyCondition(
            obs.PostPolicyOpStartsWith,
            obs.PostPolicyKeyContentType,
            "image/",
        ),
    },
}
```

---

### 场景 4：前缀匹配

创建策略，允许上传到特定前缀的路径。

```go
input := &obs.CreatePostPolicyInput{
    Bucket:    "my-bucket",
    Key:       "uploads/", // 只允许上传到 uploads/ 前缀下
    ExpiresIn: 3600,
}
```

---

## 相关文档

- [POST 上传示例代码](../../main/post_upload_sample.go)
- [OBS API 文档](https://support.huaweicloud.com/api-obs/obs_04_0108.html)
- [桶清单功能](../bucket/README.md)
