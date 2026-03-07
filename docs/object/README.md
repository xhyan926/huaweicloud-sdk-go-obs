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

CreatePostPolicy 方法生成 POST 上传策略，包括 Policy JSON 和签名。生成的凭证可以用于前端 HTML 表单，实现文件直接从浏览器上传到 OBS，无需经过后端服务器。

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
| Expires | int64 | 否 | 过期时间（秒），默认300 |
| Acl | string | 否 | 对象 ACL，如 public-read |

### 返回值

**CreatePostPolicyOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| Policy | string | Base64 编码的 Policy JSON |
| Signature | string | Policy 的 HMAC-SHA1 签名 |
| BaseModel | - | 包含 HTTP 响应元数据 |

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
        Bucket:  "my-bucket",
        Key:     "uploads/test.jpg",
        Expires: 600,  // 10 分钟后过期
        Acl:     "public-read",
    }

    output, err := obsClient.CreatePostPolicy(input)
    if err != nil {
        fmt.Printf("创建策略失败: %v\n", err)
        return
    }

    // 3. 输出策略信息
    fmt.Printf("Policy: %s\n", output.Policy)
    fmt.Printf("Signature: %s\n", output.Signature)
    // 4. 使用 Policy 和 Signature 生成前端 HTML 表单
    // 参见 examples/post_upload/post_upload_sample.go
}
```

### 前端使用示例

生成的 Policy 和 Signature 可以用于前端 HTML 表单：

```html
<form action="https://obs.cn-north-1.myhuaweicloud.com/my-bucket/" method="post" enctype="multipart/form-data">
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
3. **高级功能**：对于需要自定义条件（如文件大小限制、内容类型限制）的场景，请使用 `CreateBrowserBasedSignature` 接口
4. **默认行为**：系统会自动添加桶和键条件，无需手动添加

---

## 内部函数（不公开给客户调用）

### buildPostPolicyExpiration

生成 Policy 的过期时间字符串（ISO 8601 格式）。此函数为内部函数，用于 CreatePostPolicy 方法内部。

**注意**: 此函数不公开给客户调用。

### validatePostPolicy

验证 POST Policy 的基本参数。此函数为内部函数，用于 CreatePostPolicy 方法内部。

**注意**: 此函数不公开给客户调用。

---

## 数据结构

### PostPolicyCondition

POST 策略条件结构（用于高级 POST 上传场景）。

```go
type PostPolicyCondition struct {
    Operator string      `json:"-"` // 操作符
    Key      string      `json:"-"` // 条件键
    Value    interface{} `json:"-"` // 条件值
}
```

**自定义 JSON 序列化**：条件会被序列化为数组格式 `["operator", "key", "value"]`

### PostPolicy

POST 策略结构（用于高级 POST 上传场景）。

```go
type PostPolicy struct {
    Expiration string                `json:"expiration"` // 过期时间
    Conditions []PostPolicyCondition `json:"-"`         // 条件列表
}
```

**自定义 JSON 序列化**：支持 AWS S3 POST Policy 格式

### CreatePostPolicyInput

CreatePostPolicy 输入参数。

```go
type CreatePostPolicyInput struct {
    Bucket  string // 存储桶名称
    Key     string // 对象键
    Expires  int64  // 过期时间（秒），默认300
    Acl     string // 可选，访问控制策略
}
```

### CreatePostPolicyOutput

CreatePostPolicy 输出结果。

```go
type CreatePostPolicyOutput struct {
    BaseModel
    Policy    string `json:"policy"`   // Base64 编码的 policy
    Signature string `json:"signature"` // 签名
}
```

---

## 常量定义

### 条件键 (PostPolicyConditionKeys)

| 常量 | 值 | 说明 |
|--------|-----|------|
| PostPolicyKeyBucket | `$bucket` | 桶名称 |
| PostPolicyKeyKey | `$key` | 对象键 |
| PostPolicyKeyContentType | `$content-type` | 内容类型 |
| PostPolicyKeyContentLength | `$content-length` | 内容长度 |

### 条件操作符 (PostPolicyConditionOperators)

| 常量 | 值 | 说明 |
|--------|-----|------|
| PostPolicyOpEquals | `eq` | 等于 |
| PostPolicyOpStartsWith | `starts-with` | 以...开头 |
| PostPolicyOpRange | `content-length-range` | 内容长度范围 |

---

## 使用场景

### 场景 1：基本 POST 上传

最简单的使用场景，只需要提供桶名和对象键。

```go
input := &obs.CreatePostPolicyInput{
    Bucket: "my-bucket",
    Key:    "uploads/file.jpg",
    Expires: 3600,
}
output, err := obsClient.CreatePostPolicy(input)
```

### 场景 2：带 ACL 的 POST 上传

设置上传后对象的访问权限。

```go
input := &obs.CreatePostPolicyInput{
    Bucket: "my-bucket",
    Key:    "uploads/file.jpg",
    Expires: 3600,
    Acl:     "public-read",
}
output, err := obsClient.CreatePostPolicy(input)
```

### 场景 3：高级 POST 上传（自定义条件）

使用 CreateBrowserBasedSignature 接口实现更复杂的 POST 上传场景，如文件大小限制、内容类型限制等。

```go
// 高级用法：自定义文件大小限制
input := &obs.CreateBrowserBasedSignatureInput{
    Bucket: "my-bucket",
    Key:    "uploads/",
    FormParams: map[string]string{
        "content-type": "image/jpeg",
    },
    RangeParams: []obs.RangeParams{
        {
            RangeName: "content-length-range",
            Lower:     1,
            Upper:     10 * 1024 * 1024, // 限制为 10MB
        },
    },
    Expires: 3600,
}
output, err := obsClient.CreateBrowserBasedSignature(input)
```

## 相关文档

- [CreateBrowserBasedSignature](../README.md#createbrowserbasedsignature) - 高级 POST 上传接口
- [示例代码](../../examples/post_upload/) - 完整的 POST 上传示例
- [OBS API 参考](https://support.huaweicloud.com/api-obs/obs_04_0088.html) - POST Object API 文档

## 版本信息

- **文档版本**: 1.2
- **SDK 版本**: 3.25.9+
- **更新日期**: 2026-03-07
- **更新内容**: 简化 CreatePostPolicy 接口，移除重复功能

---

**最后更新**: 2026-03-07 (任务组 2：POST 上传策略重构完成)
