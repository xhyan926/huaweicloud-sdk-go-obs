# 子任务 2.1：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件**: `obs/model_object.go`
- **追加位置**: 在现有对象操作结构体之后

### 2. POST 策略结构体定义

```go
// PostPolicyCondition defines a condition for POST policy
type PostPolicyCondition struct {
    Operator string `json:"-"` // equals, starts-with, etc.
    Key     string `json:"-"`
    Value    interface{} `json:"-"`
}

// PostPolicy defines a POST upload policy
type PostPolicy struct {
    Expiration string                 `json:"expiration"`       // 过期时间
    Conditions []PostPolicyCondition `json:"conditions"`    // 条件列表
}

// CreatePostPolicyInput is input for creating POST policy
type CreatePostPolicyInput struct {
    Bucket      string
    Key         string
    Expires     int64        // 过期时间（秒）
    Conditions []PostPolicyCondition
    Acl         string       // optional
}

// CreatePostPolicyOutput is result of creating POST policy
type CreatePostPolicyOutput struct {
    BaseModel
    Policy string `json:"policy"`      // Base64 编码的 policy
    Signature string `json:"signature"`  // 签名
    Token string `json:"token"`        // 完整 token (ak:signature:policy)
    AccessKeyId string `json:"accessKeyId"` // Access Key ID
}
```

### 3. 支持的条件类型

```go
// PostPolicyCondition operators
const (
    PostPolicyOpEquals     = "eq"
    PostPolicyOpStartsWith = "starts-with"
    PostPolicyOpRange      = "content-length-range"
)

// PostPolicyCondition keys
const (
    PostPolicyKeyBucket           = "$bucket"
    PostPolicyKeyKey             = "$key"
    PostPolicyKeyContentType     = "$content-type"
    PostPolicyKeyContentLength   = "$content-length"
)
```

### 4. 时间估算
- 结构体定义：30 分钟
- JSON 序列化验证：15 分钟
- 代码审查：15 分钟
- **总计**: 约 1 小时（0.125 天）

## 技术要点

### AWS S3 POST Policy 格式
- 参考 AWS S3 POST 上传规范
- Policy 是 JSON 格式
- 需要 Base64 编码

### 条件类型
- eq: 完全匹配
- starts-with: 前缀匹配
- content-length-range: 内容长度范围

### 过期时间
- ISO 8601 格式
- 或 Unix 时间戳（秒）
- 建议：使用过期秒数

### JSON 序列化
- 使用 encoding/json
- 注意 interface{} 类型的序列化
- 确保输出格式正确
