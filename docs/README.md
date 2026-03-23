# OBS SDK Go API 接口文档

欢迎使用华为云 OBS Go SDK API 接口文档。本文档提供了 SDK 中所有 API 接口的详细说明、使用示例和最佳实践。

## 目录结构

本文档按照功能特性进行组织：

```
docs/
├── README.md                        # 文档总索引（本文件）
├── API_DEVELOPMENT_GUIDE.md         # API 开发规范指南 ⭐
├── replication/                     # 跨区域复制
│   └── README.md
├── dis_policy/                      # DIS 事件通知
│   └── README.md
├── encryption/                      # 加密
├── lifecycle/                       # 生命周期管理
└── multipart/                       # 分块上传
```

> 💡 **开发者注意**: 如果需要为 SDK 添加新功能或 API，请先阅读 [API 开发规范指南](./API_DEVELOPMENT_GUIDE.md)。该指南包含了完整的数据结构设计规则、trans() 方法实现模板和常见错误避免方法。

## 快速导航

### 桶管理

- [跨区域复制](./replication/README.md)
  - [设置跨区域复制规则](./replication/README.md#设置跨区域复制规则-setbucketreplication)
  - [获取跨区域复制配置](./replication/README.md#获取跨区域复制配置-getbucketreplication)
  - [删除跨区域复制规则](./replication/README.md#删除跨区域复制规则-deletebucketreplication)

- [DIS 事件通知](./dis_policy/README.md)
  - [设置 DIS 事件通知策略](./dis_policy/README.md#设置-dis-事件通知策略-setbucketdispolicy)
  - [获取 DIS 事件通知配置](./dis_policy/README.md#获取-dis-事件通知配置-getbucketdispolicy)
  - [删除 DIS 事件通知策略](./dis_policy/README.md#删除-dis-事件通知策略-deletebucketdispolicy)

## SDK 使用基础

### 创建客户端

```go
import obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

// 创建 OBS 客户端
config := &obs.Configuration{
    Endpoint:        "https://obs.cn-north-4.myhuaweicloud.com",
    AccessKeyID:     "your-access-key",
    SecretAccessKey: "your-secret-key",
    Signature:       obs.SignatureObs, // 或 obs.SignatureV4
}

client, err := obs.New(config)
if err != nil {
    panic(err)
}
```

### 使用扩展选项

```go
// 带扩展选项的 API 调用
output, err := client.Method(input, obs.WithOption(value))
```

### 错误处理

```go
output, err := client.Method(input)
if err != nil {
    if obsError, ok := err.(obs.ObsError); ok {
        fmt.Printf("错误码: %s\n", obsError.Code)
        fmt.Printf("错误信息: %s\n", obsError.Message)
    }
    return
}
```

## API 命名规范

| 操作类型 | 前缀 | 示例 |
|---------|------|------|
| 创建/设置 | Set/Create | SetBucketReplication, CreateBucket |
| 获取 | Get | GetBucketReplication, GetObject |
| 删除 | Delete | DeleteBucketReplication, DeleteObject |
| 列举 | List | ListObjects, ListBuckets |

## 签名协议支持

OBS SDK 支持多种签名协议：

| 签名类型 | 常量 | 说明 | 支持的功能 |
|---------|------|------|-----------|
| OBS 签名 | `obs.SignatureObs` | 华为云 OBS 专用签名 | 所有功能 |
| AWS V4 签名 | `obs.SignatureV4` | 兼容 Amazon S3 | 大部分功能 |

**注意**：部分功能仅支持 OBS 签名，如：
- 跨区域复制（Cross-Region Replication）
- DIS 事件通知（DIS Event Notification）

## 新增功能 (v3.26.0)

### 跨区域复制 (Cross-Region Replication)

实现对象数据到目标桶的异步复制，支持：

- **API 接口**
  - `SetBucketReplication` - 设置跨区域复制规则
  - `GetBucketReplication` - 获取跨区域复制配置
  - `DeleteBucketReplication` - 删除跨区域复制规则

- **功能特性**
  - 多规则配置（最多 100 条）
  - 前缀过滤
  - 历史对象复制
  - 存储类型指定
  - 仅支持 OBS 签名

**限制**：
- 配置大小限制：50KB
- 规则数量上限：100 条
- 仅支持 OBS 签名类型

[查看详细文档](./replication/README.md)

### DIS 事件通知 (DIS Event Notification)

将对象操作事件发送到 DIS 数据接入服务，支持：

- **API 接口**
  - `SetBucketDisPolicy` - 设置 DIS 事件通知策略
  - `GetBucketDisPolicy` - 获取 DIS 事件通知配置
  - `DeleteBucketDisPolicy` - 删除 DIS 事件通知策略

- **功能特性**
  - 多规则配置（最多 10 条）
  - 事件类型过滤
  - 对象前缀/后缀过滤
  - 委托授权机制
  - JSON 格式配置
  - 仅支持 OBS 签名

**数据格式**：
- 请求体：JSON（`Content-Type: application/json`）
- 事件列表：字符串数组（如 `["ObjectCreated:*", "ObjectRemoved:*"]`）

**限制**：
- 规则数量上限：10 条
- 仅支持 OBS 签名类型

[查看详细文档](./dis_policy/README.md)

## 版本信息

- 文档版本: 1.0
- SDK 版本: 3.26.0+
- 更新日期: 2026-03-21

## 相关资源

- [华为云 OBS 官方文档](https://support.huaweicloud.com/obs/)
- [OBS API 参考文档](https://support.huaweicloud.com/api-obs/)
- [GitHub 仓库](https://github.com/huaweicloud/huaweicloud-sdk-go-obs)
- [示例代码](../examples/)

## 反馈与支持

如有问题或建议，请通过以下方式反馈：

- 提交 Issue: [GitHub Issues](https://github.com/huaweicloud/huaweicloud-sdk-go-obs/issues)
- 华为云支持: [技术支持](https://support.huaweicloud.com/)
