# API 变更跟踪 - 在线解压策略

## 新增数据结构

### DecompressionConfiguration
- **文件路径**：obs/model_bucket.go
- **结构定义**：
  ```go
  type DecompressionConfiguration struct {
      XMLName xml.Name `xml:"DecompressionConfiguration"`
      Rules   []DecompressionRule `xml:"Rules>Rule"`
  }
  
  type DecompressionRule struct {
      XMLName     xml.Name `xml:"Rule"`
      Prefix      string   `xml:"Prefix"`
      Status      string   `xml:"Status"` // Enabled, Disabled
      CreatedAt   string   `xml:"CreatedAt"`
      ModifiedAt  string   `xml:"ModifiedAt"`
  }
  ```

### 输入/输出结构
- **SetBucketDecompressionInput**：设置桶解压策略
- **GetBucketDecompressionInput**：获取桶解压策略
- **GetBucketDecompressionOutput**：获取桶解压策略响应
- **DeleteBucketDecompressionInput**：删除桶解压策略

## 新增常量

### HTTP 头
- `HeaderDecompression = "x-obs-decompression"`
- `HeaderDecompressionConfiguration = "x-obs-decompression-configuration"`

### 参数
- `ParamDecompression = "decompression"`
- `ParamDecompressionRuleId = "decompressionRuleId"`

## 新增客户端方法

### SetBucketDecompression
- **方法签名**：`func (cli *ObsClient) SetBucketDecompression(input *SetBucketDecompressionInput, options ...OptionType) (*SetBucketDecompressionOutput, error)`
- **功能描述**：设置桶的在线解压策略
- **所属特性**：bucket
- **文档状态**：pending

### GetBucketDecompression
- **方法签名**：`func (cli *ObsClient) GetBucketDecompression(input *GetBucketDecompressionInput, options ...OptionType) (*GetBucketDecompressionOutput, error)`
- **功能描述**：获取桶的在线解压策略
- **所属特性**：bucket
- **文档状态**：pending

### DeleteBucketDecompression
- **方法签名**：`func (cli *ObsClient) DeleteBucketDecompression(input *DeleteBucketDecompressionInput, options ...OptionType) (*DeleteBucketDecompressionOutput, error)`
- **功能描述**：删除桶的在线解压策略
- **所属特性**：bucket
- **文档状态**：pending

## 文档生成状态
- [ ] 数据结构文档
- [ ] 常量文档
- [ ] 客户端方法文档
- [ ] 示例代码文档
