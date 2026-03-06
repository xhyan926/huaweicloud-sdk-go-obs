# 子任务 5.1：存量信息数据模型

## 目标
创建桶存量信息的数据结构。

## 范围
- 在 `obs/model_bucket.go` 中添加存量信息结构体
- 定义 GetBucketStorageInfoOutput
- 定义 StorageInfo 子结构
- 添加 XML 标签映射

## 依赖
- 前置子任务：无
- 阻塞：task-20

## 实施步骤
1. 在 `obs/model_bucket.go` 中添加以下结构体：
   - GetBucketStorageInfoOutput
   - StorageInfo
2. 为所有结构体添加 XML 标签映射
3. 确保结构体正确映射 API 响应

## 验收标准
- [ ] 所有结构体通过 `go vet` 检查
- [ ] XML 标签正确映射 API 响应
- [ ] 字段类型定义正确

## 状态
pending
