# {{ SDK_NAME }} API 接口文档

欢迎使用 {{ SDK_NAME }} API 接口文档。本文档提供了 SDK 中所有 API 接口的详细说明、使用示例和最佳实践。

## 目录结构

本文档按照功能特性进行组织：

```
docs/
├── README.md                 # 本文档（API 文档总索引）
├── bucket/                   # 桶操作相关 API
├── object/                   # 对象操作相关 API
├── multipart/                # 分块上传相关 API
└── ...                       # 其他特性 API
```

## 快速导航

### {{ FEATURE_1_NAME }}

- [API 1](./{{ FEATURE_1_DIR }}/#api1)
- [API 2](./{{ FEATURE_1_DIR }}/#api2)

## SDK 使用基础

### 创建客户端

```go
import {{ IMPORT_PATH }}

client, err := {{ IMPORT_PATH }}.New("your-key", "your-secret", "endpoint")
if err != nil {
    panic(err)
}
```

### 使用扩展选项

```go
output, err := client.Method(input, {{ IMPORT_PATH }}.WithOption(value))
```

### 错误处理

```go
output, err := client.Method(input)
if err != nil {
    if sdkError, ok := err.({{ IMPORT_PATH }}.SdkError); ok {
        fmt.Printf("错误码: %s\n", sdkError.Code)
        fmt.Printf("错误信息: %s\n", sdkError.Message)
    }
}
```

## API 命名规范

| 操作类型 | 前缀 | 示例 |
|---------|------|------|
| 创建/设置 | Set/Create | SetBucket, CreateObject |
| 获取 | Get | GetBucket, GetObject |
| 列举 | List | ListObjects |
| 删除 | Delete | DeleteBucket |

## 文档规范

本文档中的每个 API 接口文档包含以下部分：

1. **方法签名** - 完整的函数定义
2. **参数说明** - 所有参数的详细说明表格
3. **返回值** - 返回值结构的详细说明
4. **使用示例** - 完整的代码示例
5. **错误码** - 可能的错误码及说明
6. **注意事项** - 使用时的注意事项和最佳实践

## 版本信息

- **文档版本**: 1.0
- **SDK 版本**: {{ SDK_VERSION }}+
- **更新日期**: {{ UPDATE_DATE }}
- **Go 版本**: {{ GO_VERSION }}+

## 相关资源

- [官方文档]({{ OFFICIAL_DOCS_URL }})
- [API 参考文档]({{ API_REFERENCE_URL }})
- [SDK GitHub 仓库]({{ GITHUB_REPO }})
- [示例代码](../examples/)

## 反馈与支持

如果您在使用过程中遇到问题或有改进建议，欢迎：

- 提交 [Issue]({{ ISSUE_URL }})
- 查阅 [支持中心]({{ SUPPORT_URL }})

---

**最后更新**: {{ UPDATE_DATE }}
