# API 变更跟踪 - WORM 策略

## 新增数据结构

### BucketWormConfiguration
- **文件路径**：obs/model_bucket.go
- **结构定义**：
  ```go
  type BucketWormConfiguration struct {
      XMLName     xml.Name `xml:"WormConfiguration"`
      Version     string   `xml:"Version"`
      Retention   struct {
          Days     int    `xml:"Days"`
          Mode     string `xml:"Mode"` // COMPLIANCE, SUSPENDED
          Date     string `xml:"Date"`
      } `xml:"Retention"`
      CreateDate  string `xml:"CreateDate"`
      ModifyDate  string `xml:"ModifyDate"`
  }
  ```

### 输入/输出结构
- **SetBucketWormInput**：设置桶 WORM 策略
- **GetBucketWormInput**：获取桶 WORM 策略
- **GetBucketWormOutput**：获取桶 WORM 策略响应
- **ExtendBucketWormInput**：延长桶 WORM 策略期限

## 新增常量

### HTTP 头
- `HeaderWorm = "x-obs-worm"`
- `HeaderWormConfiguration = "x-obs-worm-configuration"`

### 参数
- `ParamWorm = "worm"`
- `ParamWormRetentionDays = "wormRetentionDays"`
- `ParamWormMode = "wormMode"`
- `ParamWormDate = "wormDate"`

### 状态常量
- `StatusCompliance = "COMPLIANCE"`
- `StatusSuspended = "SUSPENDED"`
- `StatusVersion = "1.0"`

## 新增客户端方法

### SetBucketWorm
- **方法签名**：`func (cli *ObsClient) SetBucketWorm(input *SetBucketWormInput, options ...OptionType) (*SetBucketWormOutput, error)`
- **功能描述**：设置桶的 WORM 策略
- **所属特性**：bucket
- **文档状态**：pending

### GetBucketWorm
- **方法签名**：`func (cli *ObsClient) GetBucketWorm(input *GetBucketWormInput, options ...OptionType) (*GetBucketWormOutput, error)`
- **功能描述**：获取桶的 WORM 策略
- **所属特性**：bucket
- **文档状态**：pending

### ExtendBucketWorm
- **方法签名**：`func (cli *ObsClient) ExtendBucketWorm(input *ExtendBucketWormInput, options ...OptionType) (*ExtendBucketWormOutput, error)`
- **功能描述**：延长桶的 WORM 策略期限
- **所属特性**：bucket
- **文档状态**：pending

## 文档生成状态
- [ ] 数据结构文档
- [ ] 常量文档
- [ ] 客户端方法文档
- [ ] 示例代码文档
