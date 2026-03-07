# Huawei Cloud OBS Go SDK API 接口文档

欢迎使用华为云 OBS (Object Storage Service) Go SDK API 接口文档。本文档提供了 SDK 中所有 API 接口的详细说明、使用示例和最佳实践。

## 目录结构

本文档按照功能特性进行组织：

```
docs/
├── README.md                 # 本文档（API 文档总索引）
├── bucket/                   # 桶操作相关 API
│   └── README.md            # 桶操作接口文档
├── object/                   # 对象操作相关 API
├── multipart/                # 分块上传相关 API
├── lifecycle/                # 生命周期管理相关 API
├── encryption/               # 加密相关 API
└── ...                       # 其他特性 API
```

## 快速导航

### 桶操作 (bucket/)

桶是 OBS 中存储对象的容器。桶操作 API 包括：

- [桶清单管理](./bucket/README.md#桶清单管理)
  - [SetBucketInventory](./bucket/README.md#setbucketinventory) - 设置桶清单配置
  - [GetBucketInventory](./bucket/README.md#getbucketinventory) - 获取桶清单配置
  - [ListBucketInventory](./bucket/README.md#listbucketinventory) - 列举桶清单配置
  - [DeleteBucketInventory](./bucket/README.md#deletebucketinventory) - 删除桶清单配置
- [跨区域复制](./bucket/README.md#跨区域复制)
  - [SetBucketReplication](./bucket/README.md#setbucketreplication) - 设置桶的跨区域复制配置
  - [GetBucketReplication](./bucket/README.md#getbucketreplication) - 获取桶的跨区域复制配置
  - [DeleteBucketReplication](./bucket/README.md#deletebucketreplication) - 删除桶的跨区域复制配置
- [归档存储对象直读](./bucket/README.md#归档存储对象直读)
- [SetBucketDirectColdAccess](./bucket/README.md#setbucketdirectcoldaccess) - 设置桶的归档对象直读配置
- [GetBucketDirectColdAccess](./bucket/README.md#getbucketdirectcoldaccess) - 获取桶的归档对象直读配置
- [DeleteBucketDirectColdAccess](./bucket/README.md#deletebucketdirectcoldaccess) - 删除桶的归档对象直读配置

### 对象操作 (object/)

对象是 OBS 中存储的基本数据单元。对象操作 API 包括：

- [POST 上传策略](./object/README.md#createpostpolicy) - 创建 POST 上传策略（简化版）
- [CreateBrowserBasedSignature](../README.md#createbrowserbasedsignature) - 高级 POST 上传策略（支持自定义条件）
- [辅助函数](./object/README.md#辅助函数) - Policy 构建和验证
- [数据结构](./object/README.md#数据结构) - Policy 相关数据结构
- [常量定义](./object/README.md#常量定义) - Policy 条件键和操作符
- [使用场景](./object/README.md#使用场景) - 常见使用场景
- [示例代码](../../examples/post_upload/) - 完整的 POST 上传示例

### 分块上传 (multipart/)

分块上传适用于大文件上传场景。相关 API 包括：

- 初始化分块上传
- 上传分块
- 合并分块
- 列举分块
- 取消分块上传
- ...

### 生命周期管理 (lifecycle/)

生命周期管理可以自动管理对象的过期和转换。相关 API 包括：

- 设置生命周期规则
- 获取生命周期规则
- 删除生命周期规则
- ...

### 加密 (encryption/)

加密功能保护对象数据的安全性。相关 API 包括：

- 设置桶加密
- 获取桶加密
- 删除桶加密
- ...

## SDK 使用基础

### 创建客户端

```go
import obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

// 使用默认配置创建客户端
obsClient, err := obs.New("your-ak", "your-sk", "https://obs.cn-north-4.myhuaweicloud.com")
if err != nil {
    panic(err)
}

// 使用自定义配置创建客户端
obsClient, err := obs.New(
    "your-ak",
    "your-sk",
    "https://obs.cn-north-4.myhuaweicloud.com",
    obs.WithMaxConnections(100),
    obs.WithSecurityToken("your-token"),
    obs.WithTimeout(60),
)
if err != nil {
    panic(err)
}
```

### 使用扩展选项

SDK 支持通过扩展选项增强功能：

```go
// 使用进度监听
output, err := obsClient.PutObject(input, obs.WithProgress(listener))

// 使用自定义请求头
output, err := obsClient.PutObject(input, obs.WithCustomHeader("x-obs-meta-key", "value"))

// 使用请求者付费
output, err := obsClient.GetObject(input, obs.WithReqPaymentHeader("requester"))
```

### 错误处理

```go
output, err := obsClient.SetBucketInventory(input)
if err != nil {
    if obsError, ok := err.(obs.ObsError); ok {
        fmt.Printf("错误码: %s\n", obsError.Code)
        fmt.Printf("错误信息: %s\n", obsError.Message)
        fmt.Printf("请求 ID: %s\n", obsError.RequestId)
        fmt.Printf("HTTP 状态码: %d\n", obsError.StatusCode)
    } else {
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

## API 命名规范

SDK 中的 API 方法命名遵循以下规范：

| 操作类型 | 前缀 | 示例 |
|---------|------|------|
| 创建/设置 | Set | SetBucketInventory |
| 获取 | Get | GetBucketInventory |
| 列举 | List | ListBucketInventory |
| 删除 | Delete | DeleteBucketInventory |
| 上传 | Put | PutObject |
| 下载 | Get | GetObject |
| 复制 | Copy | CopyObject |

## 文档规范

本文档中的每个 API 接口文档包含以下部分：

1. **方法签名** - 完整的函数定义
2. **参数说明** - 所有参数的详细说明表格
3. **返回值** - 返回值结构的详细说明
4. **使用示例** - 完整的代码示例
5. **错误码** - 可能的错误码及说明
6. **注意事项** - 使用时的注意事项和最佳实践

## 版本信息

- **文档版本**: 1.1
- **SDK 版本**: 3.25.9+
- **更新日期**: 2026-03-06
- **Go 版本**: 1.16+

## 相关资源

- [华为云 OBS 官方文档](https://support.huaweicloud.com/obs/index.html)
- [OBS API 参考文档](https://support.huaweicloud.com/api-obs/obs_04_0008.html)
- [SDK GitHub 仓库](https://github.com/huaweicloud/huaweicloud-sdk-go-obs)
- [示例代码](../examples/)
- [更新日志](../README_CN.MD)

## 反馈与支持

如果您在使用过程中遇到问题或有改进建议，欢迎：

- 提交 [Issue](https://github.com/huaweicloud/huaweicloud-sdk-go-obs/issues)
- 查阅 [华为云支持中心](https://support.huaweicloud.com/)
- 加入 [华为云开发者社区](https://bbs.huaweicloud.com/)

## 贡献指南

我们欢迎社区贡献。如果您想为本文档做出贡献：

1. Fork 项目仓库
2. 创建特性分支
3. 提交您的更改
4. 发起 Pull Request

---

**最后更新**: 2026-03-07 (任务组 6：桶归档存储对象直读完成)
