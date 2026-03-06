# 子任务 1.3：桶清单 Trait 层实现

## 目标
实现清单功能的请求转换逻辑，处理参数映射和序列化。

## 范围
- 在 `obs/trait_bucket.go` 中实现 `SetBucketInventoryInput.trans()` 方法
- 实现 `DeleteBucketInventoryInput` 的序列化逻辑
- 添加参数验证
- 处理 inventory 子资源参数

## 依赖
- 前置子任务：task-01, task-02
- 阻塞：task-04

## 实施步骤
1. 在 `obs/trait_bucket.go` 中实现 `SetBucketInventoryInput.trans()` 方法
2. 添加 XML 序列化逻辑
3. 实现参数验证
4. 处理 inventory 子资源参数
5. 确保符合现有 trans() 方法模式

## 验收标准
- [ ] trans() 方法返回正确的参数映射
- [ ] XML 序列化结果符合 API 规范
- [ ] 参数验证逻辑完整
- [ ] 错误处理一致

## 状态
pending
