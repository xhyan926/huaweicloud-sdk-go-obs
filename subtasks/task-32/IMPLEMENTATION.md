# 子任务 9.1：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/model_bucket.go`
- **目标文件 2**: `obs/type.go`

### 2. 结构体定义

```go
// SetBucketObjectLockInput sets up bucket-level WORM policy
type SetBucketObjectLockInput struct {
    BaseModel
    Bucket               string           `xml:"-"`
    ObjectLockEnabled   string           `xml:"ObjectLockEnabled"`
    Retention           RetentionConfig   `xml:"-"`   // 可选，保留策略
}

// GetBucketObjectLockOutput gets bucket-level WORM policy
type GetBucketObjectLockOutput struct {
    BaseModel
    ObjectLockEnabled string         `xml:"ObjectLockEnabled"`
    Retention         RetentionConfig `xml:"-"`
}

// DeleteBucketObjectLockInput deletes bucket-level WORM policy
type DeleteBucketObjectLockInput struct {
    BaseModel
    Bucket string `xml:"-"`
}

// RetentionConfig defines WORM retention configuration
type RetentionConfig struct {
    Mode   string `xml:"-"`   // COMPLIANCE
    Days   *int64 `xml:"-"`  // 可选，保留天数
    Years  *int64 `xml:"-"`  // 可选，保留年数
}
```

### 3. 常量定义

```go
SubResourceObjectLock SubResourceType = "object-lock"
```

### 4. 时间估算
- 结构体定义：30 分钟
- 常量定义：10 分钟
- 代码审查：10 分钟
- **总计**: 约 0.8 小时（0.1 天）

## 技术要点

### WORM 策略功能
- 防止数据被意外删除或修改
- 用于数据合规要求
- 支持合规性和治理模式

### 常量命名
- SubResourceObjectLock
- 与 API 子资源名称一致
