# 子任务 4.3：跨区域复制实现

## 目标
实现跨区域复制的 Set/Get/Delete 方法。

## 范围
- Trait 层的 trans() 实现
- 客户端方法实现
- 参数验证

## 依赖
- 前置子任务：task-15, task-16
- 阻塞：task-18

## 实施步骤
1. 在 `obs/trait_bucket.go` 中实现 trans() 方法
2. 在 `obs/client_bucket.go` 中添加 3 个方法：
   - SetBucketReplication()
   - GetBucketReplication()
   - DeleteBucketReplication()
3. 处理参数验证
4. 确保符合现有 API 风格

## 验收标准
- [ ] 方法符合现有 API 风格
- [ ] 错误处理一致
- [ ] 扩展选项正确传递

## 状态
pending
