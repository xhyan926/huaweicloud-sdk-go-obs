# 子任务 3.1：创建桶参数补充

## 目标
补充 CreateBucket 的缺失参数。

## 范围
- 在 `CreateBucketInput` 结构体中添加新字段
- 在 `trait_bucket.go` 的 `trans()` 方法中添加头部设置逻辑
- 更新相关常量定义

## 依赖
- 前置子任务：无
- 阻塞：task-14

## 实施步骤
1. 在 `CreateBucketInput` 中添加以下字段：
   - BucketType (对应 x-obs-bucket-type)
   - SseKmsKeyId (对应 x-obs-server-side-encryption-kms-key-id)
   - SseKmsKeyProjectId (对应 x-obs-sse-kms-key-project-id)
   - ServerSideDataEncryption (对应 x-obs-server-side-data-encryption)

2. 在 `trait_bucket.go` 的 `trans()` 方法中添加 HTTP 头部设置逻辑

3. 更新 `obs/const.go` 添加新的头部常量

## 验收标准
- [ ] 所有新字段正确映射到 HTTP 头部
- [ ] 不影响现有功能
- [ ] 代码通过 go vet 检查

## 状态
pending
