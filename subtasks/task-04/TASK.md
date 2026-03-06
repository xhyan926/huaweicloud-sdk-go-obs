# 子任务 1.4：桶清单客户端方法实现

## 目标
实现清单功能的客户端 API 方法。

## 范围
- 在 `obs/client_bucket.go` 中添加 4 个方法
- SetBucketInventory()
- GetBucketInventory(id string)
- ListBucketInventory()
- DeleteBucketInventory(id string)

## 依赖
- 前置子任务：task-03
- 阻塞：task-05

## 实施步骤
1. 在 `obs/client_bucket.go` 中添加 4 个方法
2. 实现输入验证
3. 调用 `doActionWithBucket()` 方法
4. 处理扩展选项
5. 确保符合现有 API 风格

## 验收标准
- [ ] 所有方法符合现有 API 风格
- [ ] 错误处理一致
- [ ] 扩展选项正确传递
- [ ] 文档注释完整

## 状态
pending
