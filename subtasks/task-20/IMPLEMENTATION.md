# 子任务 5.2：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件**: `obs/type.go`
- **追加位置**: 在现有 SubResourceType 常量之后

### 2. 常量定义

```go
// SubResourceStorageInfo subResource value: storageinfo
SubResourceStorageInfo SubResourceType = "storageinfo"
```

### 3. 时间估算
- 常量定义：5 分钟
- 代码审查：5 分钟
- **总计**: 约 0.17 小时（0.021 天）

## 技术要点

### 命名规范
- 常量使用大写字母
- 遵循现有命名模式
- 描述性名称

### 常量位置
- 添加到 SubResourceType 常量组
- 按字母顺序或逻辑顺序排列
- 添加注释说明用途

### 常量值
- 值: "storageinfo"
- 与 API 规范一致
- 用于查询字符串参数
