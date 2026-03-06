# 子任务 4.1：跨区域复制数据模型定义

## 目标
创建跨区域复制的数据结构。

## 范围
- 在 `obs/model_bucket.go` 中添加复制配置结构体
- 定义 Input/Output 结构体
- 定义 ReplicationRule、Destination 等子结构
- 添加 XML 标签映射

## 依赖
- 前置子任务：无
- 阻塞：task-16

## 实施步骤
1. 在 `obs/model_bucket.go` 中添加以下结构体：
   - SetBucketReplicationInput / Output
   - GetBucketReplicationOutput
   - ReplicationRule（复制规则）
   - ReplicationDestination（目标配置）
   - ReplicationFilter（筛选配置）

2. 为所有结构体添加 XML 标签映射
3. 确保结构体符合 API 规范

## 验收标准
- [ ] 所有结构体通过 `go vet` 检查
- [ ] XML 标签正确映射 API 规范
- [ ] 字段类型定义正确

## 状态
pending
