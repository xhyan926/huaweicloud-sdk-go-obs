# 子任务 3.2：PUT 上传参数补充

## 目标
补充 PutObject 的缺失参数。

## 范围
- 在 `PutObjectInput` 结构体中添加新字段
- 在 `trait_object.go` 的 `trans()` 方法中添加头部设置逻辑
- 更新相关常量定义

## 依赖
- 前置子任务：无
- 阻塞：task-14

## 实施步骤
1. 在 `PutObjectInput` 中添加以下字段：
   - Expires (对应 x-obs-expires)
   - ObjectLockMode (对应 x-obs-object-lock-mode)
   - ObjectLockRetainUntilDate (对应 x-obs-object-lock-retain-until-date)
   - ServerSideDataEncryption (对应 x-obs-server-side-data-encryption)
   - SseKmsKeyId (对应 x-obs-server-side-encryption-kms-key-id)

2. 在 `trait_object.go` 的 `trans()` 方法中添加 HTTP 头部设置逻辑

3. 更新 `obs/const.go` 添加新的头部常量

## 验收标准
- [ ] 所有新字段正确映射到 HTTP 头部
- [ ] 向后兼容
- [ ] 代码通过 go vet 检查

## 状态
pending
