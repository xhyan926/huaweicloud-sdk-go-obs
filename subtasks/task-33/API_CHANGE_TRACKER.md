# API 变更跟踪 - WORM 策略实现

## 修改的文件

### obs/trait_bucket.go
- **新增方法**：
  - `SetBucketWormTrait`
  - `GetBucketWormTrait`
  - `ExtendBucketWormTrait`

### obs/client_bucket.go
- **新增方法**：
  - `SetBucketWorm`
  - `GetBucketWorm`
  - `ExtendBucketWorm`

### obs/model_bucket.go
- **新增方法**：
  - `ToXml()`：序列化方法
  - `FromXml()`：反序列化方法

## 新增错误码

- `ErrInvalidWormConfiguration`：无效的 WORM 配置
- `ErrInvalidWormRetentionDays`：无效的保存期限
- `ErrInvalidWormMode`：无效的模式
- `ErrWormPolicyAlreadyExists`：WORM 策略已存在
- `ErrWormPolicyNotFound`：WORM 策略不存在
- `ErrWormPolicyCannotModify`：WORM 策略无法修改（COMPLIANCE 状态）
- `ErrInvalidWormVersion`：无效的版本

## 实现细节

### XML 序列化
```go
func (c *BucketWormConfiguration) ToXml() string {
    // 实现序列化逻辑
}

func (c *BucketWormConfiguration) FromXml(xmlData string) error {
    // 实现反序列化逻辑
}
```

### Trait 层方法
```go
func (cli *ObsClient) SetBucketWormTrait(input *SetBucketWormInput) *SetBucketWormOutput {
    // 实现设置逻辑，包含 WORM 约束检查
}

func (cli *ObsClient) GetBucketWormTrait(input *GetBucketWormInput) *GetBucketWormOutput {
    // 实现获取逻辑，包含状态处理
}

func (cli *ObsClient) ExtendBucketWormTrait(input *ExtendBucketWormInput) *ExtendBucketWormOutput {
    // 实现延长逻辑，包含版本控制
}
```

## 文档生成状态
- [x] 数据结构文档（子任务9.1）
- [x] 常量文档（子任务9.1）
- [ ] 客户端方法文档（待生成）
- [ ] 示例代码文档（待生成）
