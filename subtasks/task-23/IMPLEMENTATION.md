# 子任务 6.1：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/model_bucket.go`
- **目标文件 2**: `obs/type.go`

### 2. 结构体定义

```go
// SetDirectColdAccessInput sets up cold access configuration
type SetDirectColdAccessInput struct {
    BaseModel
    Bucket string `xml:"-"`
}

// GetDirectColdAccessOutput gets cold access configuration
type GetDirectColdAccessOutput struct {
    BaseModel
    Enabled bool `xml:"Enabled"`
}

// DeleteDirectColdAccessInput deletes cold access configuration
type DeleteDirectColdAccessInput struct {
    BaseModel
    Bucket string `xml:"-"`
}
```

### 3. 常量定义

```go
SubResourceDirectcoldaccess SubResourceType = "directcoldaccess"
```

### 4. 时间估算
- 结构体定义：30 分钟
- 常量定义：10 分钟
- 代码审查：10 分钟
- **总计**: 约 0.8 小时（0.1 天）

## 技术要点

### 归档直读功能
- 开启后归档对象不需要恢复便可下载
- 提升访问效率
- 简化操作流程

### 常量命名
- SubResourceDirectcoldaccess
- 与 API 子资源名称一致
