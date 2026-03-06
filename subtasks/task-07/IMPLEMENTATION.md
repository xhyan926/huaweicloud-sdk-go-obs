# 子任务 2.2：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件**: `obs/post_policy.go`（新建）

### 2. Policy 构建逻辑

```go
package obs

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "time"
)

// buildPostPolicyJSON builds the POST policy JSON
func buildPostPolicyJSON(policy *PostPolicy) (string, error) {
    if policy == nil {
        return "", errors.New("policy is nil")
    }

    // 验证过期时间
    if policy.Expiration == "" {
        return "", errors.New("expiration is required")
    }

    // 验证条件
    if len(policy.Conditions) == 0 {
        return "", errors.New("at least one condition is required")
    }

    // 构建 JSON
    jsonData, err := json.Marshal(policy)
    if err != nil {
        return "", fmt.Errorf("failed to marshal policy: %v", err)
    }

    return string(jsonData), nil
}

// BuildPostPolicyExpiration builds expiration time string
func BuildPostPolicyExpiration(expiresIn int64) string {
    expirationTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
    return expirationTime.UTC().Format("2006-01-02T15:04:05.000Z")
}

// ValidatePostPolicy validates the POST policy
func ValidatePostPolicy(policy *PostPolicy) error {
    if policy == nil {
        return errors.New("policy is nil")
    }

    // 验证过期时间
    if policy.Expiration == "" {
        return errors.New("expiration is required")
    }

    // 验证条件
    for i, cond := range policy.Conditions {
        if cond.Key == "" {
            return fmt.Errorf("condition at index %d has empty key", i)
        }
        if cond.Operator == "" {
            return fmt.Errorf("condition at index %d has empty operator", i)
        }
    }

    return nil
}
```

### 3. 辅助函数

```go
// CreatePostPolicyCondition creates a policy condition
func CreatePostPolicyCondition(operator, key string, value interface{}) PostPolicyCondition {
    return PostPolicyCondition{
        Operator: operator,
        Key:     key,
        Value:    value,
    }
}

// CreateBucketCondition creates a bucket condition
func CreateBucketCondition(bucket string) PostPolicyCondition {
    return CreatePostPolicyCondition(
        PostPolicyOpEquals,
        PostPolicyKeyBucket,
        bucket,
    )
}

// CreateKeyCondition creates a key condition
func CreateKeyCondition(key string) PostPolicyCondition {
    return CreatePostPolicyCondition(
        PostPolicyOpStartsWith,
        PostPolicyKeyKey,
        key,
    )
}

// CreateContentLengthRangeCondition creates a content-length-range condition
func CreateContentLengthRangeCondition(min, max int64) PostPolicyCondition {
    return CreatePostPolicyCondition(
        PostPolicyOpRange,
        PostPolicyKeyContentLength,
        []interface{}{min, max},
    )
}
```

### 4. 时间估算
- 文件创建和结构：20 分钟
- Policy JSON 生成：30 分钟
- 验证逻辑：20 分钟
- 过期时间处理：20 分钟
- 测试和调试：30 分钟
- **总计**: 约 2 小时（0.25 天）

## 技术要点

### JSON 格式要求
- 必须符合 AWS S3 POST Policy 规范
- 字段顺序不重要
- 条件数组格式正确

### 验证逻辑
- 检查必需字段
- 验证条件格式
- 提供清晰的错误信息

### 过期时间格式
- ISO 8601 格式: "2026-03-05T10:30:00.000Z"
- 使用 UTC 时区
- 精确到毫秒

### Base64 编码
- 在后续步骤（task-08）中进行
- 使用标准 Base64 编码
- URL 安全编码
