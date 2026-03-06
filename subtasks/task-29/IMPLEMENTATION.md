# 子任务 8.1：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/model_bucket.go`
- **目标文件 2**: `obs/type.go`

### 2. 结构体定义

```go
// SetZipPolicyInput sets up ZIP extraction policy
type SetZipPolicyInput struct {
    BaseModel
    Bucket      string `xml:"-"`
    ZipPolicy   string `xml:"-"` // ZIP 解压策略配置
}

// GetZipPolicyOutput gets ZIP extraction policy
type GetZipPolicyOutput struct {
    BaseModel
    ZipPolicy string `xml:"zipPolicy"`
}

// DeleteZipPolicyInput deletes ZIP extraction policy
type DeleteZipPolicyInput struct {
    BaseModel
    Bucket string `xml:"-"`
}
```

### 3. 常量定义

```go
SubResourceZip SubResourceType = "policy=zip"
```

### 4. 时间估算
- 结构体定义：30 分钟
- 常量定义：10 分钟
- 代码审查：10 分钟
- **总计**: 约 0.8 小时（0.1 天）

## 技术要点

### 在线解压功能
- 自动解压上传的 ZIP 文件
- 解压后的文件存放到指定位置
- 简化批量文件操作

### 常量命名
- SubResourceZip
- 与 API 子资源名称一致
