# DIS 通知策略 API 接口文档

本文档包含华为云 OBS Go SDK 中 DIS 通知策略相关的 API 接口说明。

## 目录
- [SetDisPolicy](#setdispolicy) - 设置 DIS 通知策略
- [GetDisPolicy](#getdispolicy) - 获取 DIS 通知策略
- [DeleteDisPolicy](#deletedispolicy) - 删除 DIS 通知策略

---

## SetDisPolicy

设置桶的 DIS 通知策略。DIS（Data Ingestion Service）提供实时事件通知功能，当桶内发生对象操作时，可以将事件通知发送到指定的 DIS 通道进行实时处理。

### 方法签名
```go
func (obsClient ObsClient) SetDisPolicy(input *SetDisPolicyInput, extensions ...extensionOptions) (output *BaseModel, err error)
```

### 参数说明

**SetDisPolicyInput**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名 |
| Rules | []DisPolicyRule | 是 | DIS 策略规则列表，最多支持 10 条规则 |

**DisPolicyRule**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 规则 ID，桶内唯一标识，长度 1~256 个字符 |
| stream | string | 是 | DIS 服务通道名称，需提前在 DIS 服务创建 |
| project | string | 否 | DIS 服务通道所属的项目 ID |
| events | []string | 是 | OBS 事件列表，支持的事件类型 |
| prefix | string | 否 | 对象名前缀，用于过滤对象 |
| suffix | string | 否 | 对象名后缀，用于过滤对象 |
| agency | string | 是 | IAM 委托名称，需具备 DIS 服务访问权限 |

### 返回值

**BaseModel**
| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID，用于追踪和定位问题 |

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|------|------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足，需要桶拥有者或 Tenant Administrator 权限 |
| NoSuchBucket | 404 | 桶不存在 |

### 使用示例

```go
package main

import (
    "fmt"
    "obs"
)

func main() {
    // 创建 OBS 客户端
    // 请替换为您的实际 AK/SK 和 endpoint
    ak := "your-access-key-id"
    sk := "your-secret-access-key"
    endpoint := "https://obs.cn-north-4.myhuaweicloud.com"

    obsClient, err := obs.New(ak, sk, endpoint)
    if err != nil {
        panic(err)
    }

    // 设置 DIS 策略
    input := &obs.SetDisPolicyInput{
        Bucket: "my-bucket",
        Rules: []obs.DisPolicyRule{
            {
                ID:      "rule-01",
                Stream:  "my-dis-stream",
                Project: "my-project",
                Events:  []string{"ObjectCreated:*", "ObjectRemoved:*"},
                Prefix:  "",
                Suffix:  "",
                Agency:  "dis-agency",
            },
        },
    }

    output, err := obsClient.SetDisPolicy(input)
    if err != nil {
        fmt.Printf("设置 DIS 策略失败: %v\n", err)
        return
    }

    fmt.Printf("设置 DIS 策略成功，RequestId: %s\n", output.RequestId)
}
```

### 注意事项

1. **幂等性**：接口是幂等的，如果桶上已存在相同策略内容，则返回成功（状态码 200），否则返回创建（状态码 201）
2. **规则数量限制**：同一个桶最多支持 10 条规则
3. **规则唯一性**：规则 ID 必须在桶内唯一
4. **前缀和后缀**：prefix 和 suffix 加起来长度最大为 1024 个字符
5. **权限要求**：必须是桶拥有者或者拥有 Tenant Administrator 权限
6. **DIS 服务要求**：需要提前在 DIS 服务中创建流（stream），并配置相应的 IAM 委托
7. **事件类型**：支持的事件类型包括：
   - `ObjectCreated:*` - 所有对象创建事件
   - `ObjectRemoved:*` - 所有对象删除事件
   - 具体事件类型（如 `ObjectCreated:Put`、`ObjectCreated:Post`、`ObjectCreated:Copy` 等）

---

## GetDisPolicy

获取桶的 DIS 通知策略配置。

### 方法签名
```go
func (obsClient ObsClient) GetDisPolicy(bucketName string, extensions ...extensionOptions) (output *GetDisPolicyOutput, err error)
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

### 返回值

**GetDisPolicyOutput**
| 字段 | 类型 | 说明 |
|------|------|------|
| BaseModel | - | 基础响应信息 |
| DisPolicy | string | - DIS 策略配置，JSON 格式的字符串 |

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|------|------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |

### 使用示例

```go
package main

import (
    "fmt"
    "obs"
)

func main() {
    // 创建 OBS 客户端
    // 请替换为您的实际 AK/SK 和 endpoint
    ak := "your-access-key-id"
    sk := "your-secret-access-key"
    endpoint := "https://obs.cn-north-4.myhuaweicloud.com"

    obsClient, err := obs.New(ak, sk, endpoint)
    if err != nil {
        panic(err)
    }

    // 获取 DIS 策略
    output, err := obsClient.GetDisPolicy("my-bucket")
    if err != nil {
        fmt.Printf("获取 DIS 策略失败: %v\n", err)
        return
    }

    fmt.Printf("DIS 策略配置：\n%s\n", output.DisPolicy)
    fmt.Printf("RequestId: %s\n", output.RequestId)
}
```

### 注意事项

1. **响应格式**：响应为 JSON 格式，包含完整的 DIS 策略配置
2. **权限要求**：必须是桶拥有者或者拥有 Tenant Administrator 权限
3. **策略存在**：只有当桶配置了 DIS 策略时才能获取，否则返回相应的错误码
4. **JSON 解析**：返回的 `DisPolicy` 字段包含完整的 JSON 配置字符串，需要应用方进行 JSON 解析

---

## DeleteDisPolicy

删除桶的 DIS 通知策略配置。

### 方法签名
```go
func (obsClient ObsClient) DeleteDisPolicy(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error)
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

### 返回值

**BaseModel**
| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID，用于追踪和定位问题 |

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|------|------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足，需要桶拥有者或 Tenant Administrator 权限 |
| NoSuchBucket | 404 | 桶不存在 |

### 使用示例

```go
package main

import (
    "fmt"
    "obs"
)

func main() {
    // 创建 OBS 客户端
    // 请替换为您的实际 AK/SK 和 endpoint
    ak := "your-access-key-id"
    sk := "your-secret-access-key"
    endpoint := "https://obs.cn-north-4.myhuaweicloud.com"

    obsClient, err := obs.New(ak, sk, endpoint)
    if err != nil {
        panic(err)
    }

    // 删除 DIS 策略
    output, err := obsClient.DeleteDisPolicy("my-bucket")
    if err != nil {
        fmt.Printf("删除 DIS 策略失败: %v\n", err)
        return
    }

    fmt.Printf("删除 DIS 策略成功，RequestId: %s\n", output.RequestId)
}
```

### 注意事项

1. **删除条件**：只有当桶配置了 DIS 策略时才能删除
2. **权限要求**：必须是桶拥有者或者拥有 Tenant Administrator 权限
3. **不可恢复**：删除操作不可逆，删除后配置将无法恢复
4. **影响范围**：删除配置后，该桶的所有 DIS 事件通知将停止
5. **幂等性**：如果配置不存在或已删除，重复删除会返回成功

---

## 相关文档

- [OBS API 文档](https://support.huaweicloud.com/api-obs/obs_04_0139.html)
- [桶清单管理](./bucket/README.md#桶清单管理)
- [跨区域复制](./bucket/README.md#跨区域复制)
- [归档存储对象直读](./bucket/README.md#归档存储对象直读)

---

**文档版本**: 1.0

**更新日期**: 2026-03-07

**SDK 版本**: 3.25.9+
