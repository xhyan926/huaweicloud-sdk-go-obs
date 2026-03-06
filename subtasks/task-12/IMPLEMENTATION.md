# 子任务 3.2：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/model_object.go`
- **目标文件 2**: `obs/trait_object.go`
- **目标文件 3**: `obs/const.go`

### 2. 添加新字段到 PutObjectInput

```go
// 在 obs/model_object.go 中的 PutObjectInput 添加

type PutObjectInput struct {
    // ... 现有字段 ...

    // 新增字段
    Expires                      int    `xml:"-"`  // 对象过期时间（天）
    ObjectLockMode              string `xml:"-"`  // 对象 WORM 模式
    ObjectLockRetainUntilDate string `xml:"-"`  // WORM 保留截止时间
    ServerSideDataEncryption   string `xml:"-"`  // 数据加密算法
    SseKmsKeyId                 string `xml:"-"`  // KMS 密钥 ID
}
```

### 3. 添加常量定义

```go
// 在 obs/const.go 中添加

// Object lock and expiration headers
const (
    HEADER_EXPIRES                         = "x-obs-expires"
    HEADER_OBJECT_LOCK_MODE                = "x-obs-object-lock-mode"
    HEADER_OBJECT_LOCK_RETAIN_UNTIL_DATE = "x-obs-object-lock-retain-until-date"
)
```

### 4. 更新 trans() 方法

```go
// 在 obs/trait_object.go 的 PutObjectInput trans() 方法中添加

func (input PutObjectInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
    headers = make(map[string][]string)

    // ... 现有代码 ...

    // 新增：设置对象过期时间
    if input.Expires > 0 {
        setHeaders(headers, HEADER_EXPIRES, []string{IntToString(input.Expires)}, true)
    }

    // 新增：设置对象 WORM 模式
    if objectLockMode := input.ObjectLockMode; objectLockMode != "" {
        setHeaders(headers, HEADER_OBJECT_LOCK_MODE, []string{objectLockMode}, true)
    }

    // 新增：设置 WORM 保留时间
    if retainUntilDate := input.ObjectLockRetainUntilDate; retainUntilDate != "" {
        setHeaders(headers, HEADER_OBJECT_LOCK_RETAIN_UNTIL_DATE, []string{retainUntilDate}, true)
    }

    // 新增：设置数据加密算法（如果未通过 SseHeader 设置）
    if dataEncryption := input.ServerSideDataEncryption; dataEncryption != "" {
        setHeadersNext(headers, HEADER_SERVER_SIDE_DATA_ENCRYPTION, HEADER_SERVER_SIDE_DATA_ENCRYPTION, []string{dataEncryption}, true)
    }

    // 新增：设置 KMS 密钥 ID
    if kmsKeyId := input.SseKmsKeyId; kmsKeyId != "" {
        setHeaders(headers, HEADER_SSE_KMS_KEY_ID, []string{kmsKeyId}, true)
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

### WORM 参数
- ObjectLockMode: WORM 模式（如 COMPLIANCE）
- ObjectLockRetainUntilDate: 保留截止日期
- 必须一起设置才有意义

### 过期时间
- 单位：天
- 值为正整数
- 设置后对象会自动过期

### 加密参数
- ServerSideDataEncryption: 加密算法（AES256/SM4）
- SseKmsKeyId: KMS 密钥 ID
- 可以与 SSE-OBS 或 SSE-KMS 结合使用

### 向后兼容性
- 所有新字段为可选
- 不设置时不影响现有功能
- 现有测试应继续通过
