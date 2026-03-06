# 子任务 3.1：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/model_bucket.go`
- **目标文件 2**: `obs/trait_bucket.go`
- **目标文件 3**: `obs/const.go`

### 2. 添加新字段到 CreateBucketInput

```go
// 在 obs/model_bucket.go 中的 CreateBucketInput 添加

type CreateBucketInput struct {
    BucketLocation
    Bucket                      string               `xml:"-"`
    ACL                         AclType              `xml:"-"`
    StorageClass                StorageClassType     `xml:"-"`
    GrantReadId                 string               `xml:"-"`
    GrantWriteId                string               `xml:"-"`
    GrantReadAcpId              string               `xml:"-"`
    GrantWriteAcpId             string               `xml:"-"`
    GrantFullControlId          string               `xml:"-"`
    GrantReadDeliveredId        string               `xml:"-"`
    GrantFullControlDeliveredId string               `xml:"-"`
    Epid                        string               `xml:"-"`
    AvailableZone               string               `xml:"-"`
    IsFSFileInterface           bool                 `xml:"-"`
    BucketRedundancy            BucketRedundancyType `xml:"-"`
    IsFusionAllowUpgrade        bool                 `xml:"-"`
    // 新增字段
    BucketType                 string               `xml:"-"`  // 桶类型
    SseKmsKeyId                string               `xml:"-"`  // KMS 密钥 ID
    SseKmsKeyProjectId         string               `xml:"-"`  // KMS 密钥项目 ID
    ServerSideDataEncryption    string               `xml:"-"`  // 数据加密算法
}
```

### 3. 添加常量定义

```go
// 在 obs/const.go 中添加

// Bucket encryption and data encryption headers
const (
    HEADER_BUCKET_TYPE                = "x-obs-bucket-type"
    HEADER_SSE_KMS_KEY_ID           = "x-obs-server-side-encryption-kms-key-id"
    HEADER_SSE_KMS_KEY_PROJECT_ID    = "x-obs-sse-kms-key-project-id"
    HEADER_SERVER_SIDE_DATA_ENCRYPTION = "x-obs-server-side-data-encryption"
)
```

### 4. 更新 trans() 方法

```go
// 在 obs/trait_bucket.go 的 CreateBucketInput trans() 方法中添加

func (input CreateBucketInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
    headers = make(map[string][]string)

    // ... 现有代码 ...

    // 新增：设置桶类型
    if bucketType := input.BucketType; bucketType != "" {
        setHeaders(headers, HEADER_BUCKET_TYPE, []string{bucketType}, true)
    }

    // 新增：设置数据加密算法
    if dataEncryption := input.ServerSideDataEncryption; dataEncryption != "" {
        setHeaders(headers, HEADER_SERVER_SIDE_DATA_ENCRYPTION, []string{dataEncryption}, true)
    }

    // 新增：设置 KMS 密钥 ID
    if kmsKeyId := input.SseKmsKeyId; kmsKeyId != "" {
        setHeaders(headers, HEADER_SSE_KMS_KEY_ID, []string{kmsKeyId}, true)
    }

    // 新增：设置 KMS 密钥项目 ID
    if kmsKeyProjectId := input.SseKmsKeyProjectId; kmsKeyProjectId != "" {
        setHeaders(headers, HEADER_SSE_KMS_KEY_PROJECT_ID, []string{kmsKeyProjectId}, true)
    }

    // ... 现有代码继续 ...

    return
}
```

### 5. 时间估算
- 字段添加：15 分钟
- 常量定义：10 分钟
- trans() 方法更新：15 分钟
- 测试和验证：15 分钟
- **总计**: 约 1 小时（0.125 天）

## 技术要点

### 向后兼容性
- 新字段为可选（指针类型）
- 不设置时不发送 HTTP 头
- 现有功能不受影响

### 常量命名
- 遵循现有命名规范
- 使用 HEADER_ 前缀
- 描述性名称

### HTTP 头处理
- 使用 setHeaders 函数
- 正确处理 OBS 和 AWS 头部
- 避免重复设置
