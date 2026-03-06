# 子任务 1.3：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件**: `obs/trait_bucket.go`
- **追加位置**: 在现有 trans() 方法之后

### 2. SetBucketInventoryInput trans() 实现

```go
func (input SetBucketInventoryInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
    params = make(map[string]string)
    params[string(SubResourceInventory)] = ""

    data, err = ConvertRequestToIoReader(input)
    if err != nil {
        return
    }

    return
}
```

### 3. DeleteBucketInventoryInput trans() 实现

```go
func (input DeleteBucketInventoryInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
    params = make(map[string]string)
    if input.Id != "" {
        params[string(SubResourceInventory)] = input.Id
    }

    return
}
```

### 4. 参数验证逻辑

在 `trans()` 方法中添加验证：
- Bucket 名称不为空
- Inventory ID 不为空（对于删除操作）
- IsEnabled 为 true 时必须包含 Destination

### 5. 时间估算
- trans() 方法实现：30 分钟
- 参数验证：20 分钟
- 测试和调试：30 分钟
- **总计**: 约 1.3 小时（0.17 天）

## 技术要点

### 子资源参数处理
- inventory 子资源可以带 ID 参数
- 使用 SubResourceInventory 常量
- 正确构建 params 映射

### XML 序列化
- 使用 ConvertRequestToIoReader 函数
- 确保结构体正确转换为 XML
- 处理嵌套结构的序列化

### 参数验证
- 参考现有方法的验证逻辑
- 提供清晰的错误信息
- 验证必选字段
