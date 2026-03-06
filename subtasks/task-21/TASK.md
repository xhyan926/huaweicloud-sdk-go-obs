# 子任务 5.3：存量信息实现

## 目标
实现 GetBucketStorageInfo 方法。

## 范围
- 客户端方法实现
- 参数处理

## 依赖
- 前置子任务：task-19, task-20
- 阻塞：task-22

## 实施步骤
1. 在 `obs/client_bucket.go` 中添加 GetBucketStorageInfo 方法
2. 实现请求逻辑
3. 处理参数验证

## 验收标准
- [ ] 方法符合现有 API 风格
- [ ] 错误处理一致

## 状态
pending
