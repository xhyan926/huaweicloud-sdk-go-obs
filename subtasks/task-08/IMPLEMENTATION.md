# 子任务 2.3：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/post_policy.go`
- **目标文件 2**: `obs/client_object.go`
- **追加位置**: 在现有对象方法之后

### 2. 签名计算逻辑

```go
// 在 post_policy.go 中添加

// CalculatePostPolicySignature calculates the signature for POST policy
func CalculatePostPolicySignature(policyJSON, secretAccessKey string) (string, error) {
    // Base64 编码 Policy
    encodedPolicy := base64.StdEncoding.EncodeToString([]byte(policyJSON))

    // 计算签名（复用现有逻辑）
    // 参考 obs/auth.go 中的签名计算
    signature := calculateHmacSHA1(encodedPolicy, secretAccessKey)

    return signature, nil
}

// BuildPostPolicyToken builds the complete token string
func BuildPostPolicyToken(ak, signature, policy string) string {
    return fmt.Sprintf("%s:%s:%s", ak, signature, policy)
}
```

### 3. CreatePostPolicy 客户端方法

```go
// 在 client_object.go 中添加

// CreatePostPolicy creates a POST upload policy with signature
func (obsClient ObsClient) CreatePostPolicy(input *CreatePostPolicyInput) (output *CreatePostPolicyOutput, err error) {
    if input == nil {
        return nil, errors.New("CreatePostPolicyInput is nil")
    }
    if input.Bucket == "" {
        return nil, errors.New("bucket is empty")
    }
    if input.Key == "" {
        return nil, errors.New("key is empty")
    }

    // 构建过期时间
    expiration := BuildPostPolicyExpiration(input.Expires)

    // 构建 Policy 结构
    policy := &PostPolicy{
        Expiration: expiration,
        Conditions: input.Conditions,
    }

    // 添加默认条件
    bucketCond := CreateBucketCondition(input.Bucket)
    keyCond := CreateKeyCondition(input.Key)
    policy.Conditions = append(policy.Conditions, bucketCond, keyCond)

    // 验证 Policy
    if err := ValidatePostPolicy(policy); err != nil {
        return nil, err
    }

    // 生成 JSON
    policyJSON, err := buildPostPolicyJSON(policy)
    if err != nil {
        return nil, err
    }

    // 计算签名
    signature, err := CalculatePostPolicySignature(policyJSON, obsClient.conf.securityProvider.Secret())
    if err != nil {
        return nil, err
    }

    // Base64 编码 Policy
    encodedPolicy := base64.StdEncoding.EncodeToString([]byte(policyJSON))

    // 构建 Token
    token := BuildPostPolicyToken(
        obsClient.conf.securityProvider.Access(),
        signature,
        encodedPolicy,
    )

    output = &CreatePostPolicyOutput{
        BaseModel:    BaseModel{},
        Policy:      encodedPolicy,
        Signature:    signature,
        Token:       token,
        AccessKeyId: obsClient.conf.securityProvider.Access(),
    }

    return
}
```

### 4. 时间估算
- 签名计算实现：30 分钟
- Token 生成逻辑：15 分钟
- 客户端方法实现：30 分钟
- 测试和调试：30 分钟
- **总计**: 约 1.8 小时（0.22 天）

## 技术要点

### 签名算法
- 使用 HMAC-SHA1（V2 签名）
- 或 HMAC-SHA256（V4 签名）
- 参考 obs/auth.go 的实现

### Base64 编码
- 使用标准 Base64 编码
- 不是 URL 安全编码
- Policy 需要先编码再签名

### Token 格式
- 格式: `ak:signature:policy`
- ak: Access Key ID
- signature: 计算出的签名
- policy: Base64 编码的 Policy JSON

### 安全性
- 不要在日志中输出 SK
- 验证 Policy 签名时使用正确的 SK
- 保护敏感信息
