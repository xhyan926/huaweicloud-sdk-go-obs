# 子任务 3.3：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/model_object.go`
- **目标文件 2**: `obs/trait_object.go`

### 2. 添加新字段到 ListObjectsInput

```go
// 在 obs/model_object.go 中的 ListObjectsInput 添加

type ListObjectsInput struct {
    Prefix  string `xml:"-"`
    Marker  string `xml:"-"`
    MaxKeys int    `xml:"-"`
    Delimiter string `xml:"-"`
    // 新增字段
    EncodingType string `xml:"-"` // 响应编码类型 (URL)
}
```

### 3. 更新 trans() 方法

```go
// 在 obs/trait_object.go 的 ListObjectsInput trans() 方法中添加

func (input ListObjectsInput) trans(isObs bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
    params = make(map[string]string)

    if input.MaxKeys > 0 {
        params["max-keys"] = IntToString(input.MaxKeys)
    }
    if input.Marker != "" {
        params["marker"] = input.Marker
    }
    if input.Prefix != "" {
        params["prefix"] = input.Prefix
    }
    if input.Delimiter != "" {
        params["delimiter"] = input.Delimiter
    }

    // 新增：添加 encoding-type 参数
    if input.EncodingType != "" {
        params["encoding-type"] = input.EncodingType
    }

    return
}
```

### 4. 更新响应处理

```go
// 在 obs/model_object.go 中的相关结构体添加

// ListObjectsOutput is result of ListObjects function
type ListObjectsOutput struct {
    BaseModel
    Prefixes       []Prefix `xml:"CommonPrefixes>Prefix"`
    Objects        []Object `xml:"Contents"`
    IsTruncated    bool     `xml:"IsTruncated"`
    Marker         string   `xml:"Marker"`
    NextMarker    string   `xml:"NextMarker"`
    Location       string   `xml:"Location"`
    Delimiter      string   `xml:"Delimiter"`
    EncodingType   string   `xml:"EncodingType"` // 新增：记录使用的编码类型
}
```

### 5. 时间估算
- 字段添加：10 分钟
- trans() 方法更新：10 分钟
- 响应处理更新：10 分钟
- 测试和验证：10 分钟
- **总计**: 约 0.7 小时（0.0875 天）

## 技术要点

### EncodingType 参数
- 值: "url" 或空字符串
- 影响: 响应中对象名和前缀的编码
- 当设置为 "url" 时，响应会使用 URL 编码

### 响应编码处理
- 当 EncodingType 为 "url" 时：
  - 对象名会进行 URL 解码
  - 前缀会进行 URL 解码
  - 方便处理特殊字符

### 向后兼容性
- 默认不设置 EncodingType
- 现有行为不受影响
- 只在明确设置时生效

### 查询参数构建
- 使用 params 映射
- encoding-type 会被添加到查询字符串
- 格式: ?encoding-type=url
