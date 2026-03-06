# 桶相关 API 接口文档

本文档包含华为云 OBS Go SDK 中所有桶操作相关的 API 接口说明。

## 目录

- [桶清单管理](#桶清单管理)

---

## 桶清单管理

桶清单功能可以定期列举桶内对象，并将对象元数据信息以 CSV 格式存储到指定的桶中。

### SetBucketInventory

设置桶清单配置。

#### 方法签名

```go
func (obsClient ObsClient) SetBucketInventory(input *SetBucketInventoryInput, extensions ...extensionOptions) (output *BaseModel, err error)
```

#### 参数说明

**SetBucketInventoryInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名 |
| InventoryConfiguration | InventoryConfiguration | 是 | 清单配置 |

**InventoryConfiguration**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Id | string | 是 | 清单规则 ID，桶内唯一 |
| IsEnabled | bool | 是 | 是否启用清单 |
| Destination | InventoryDestination | 是 | 清单报告的目标配置 |
| Schedule | InventorySchedule | 是 | 清单生成调度 |
| Filter | *InventoryFilter | 否 | 对象筛选条件 |
| IncludedObjectVersions | string | 否 | 版本包含策略，All 或 Current |
| OptionalFields | *InventoryOptionalFields | 否 | 可选的元数据字段 |

**InventoryDestination**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Format | string | 是 | 报告格式，支持 CSV |
| Bucket | string | 是 | 存储报告的目标桶 |
| Prefix | string | 否 | 报告对象的前缀 |

**InventorySchedule**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Frequency | InventoryFrequencyType | 是 | 调度频率，Daily 或 Weekly |

**InventoryFilter**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Prefix | string | 否 | 对象名前缀 |

**InventoryOptionalFields**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Fields | []string | 否 | 可选字段列表 |

**可选字段常量**

| 常量 | 说明 |
|------|------|
| InventoryFieldSize | 对象大小 |
| InventoryFieldLastModifiedDate | 最后修改时间 |
| InventoryFieldETag | 对象 ETag |
| InventoryFieldStorageClass | 存储类型 |
| InventoryFieldIsMultipartUploaded | 是否分块上传 |
| InventoryFieldReplicationStatus | 复制状态 |
| InventoryFieldEncryptionStatus | 加密状态 |
| InventoryFieldObjectLockRetainUntilDate | 对象锁保留到期日期 |
| InventoryFieldObjectLockMode | 对象锁模式 |

#### 返回值

**BaseModel**

| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID |

#### 使用示例

```go
package main

import (
    "fmt"
    obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func main() {
    // 创建客户端
    obsClient, err := obs.New("your-ak", "your-sk", "https://obs.cn-north-4.myhuaweicloud.com")
    if err != nil {
        panic(err)
    }

    // 设置每日清单报告
    input := &obs.SetBucketInventoryInput{
        Bucket: "my-bucket",
        InventoryConfiguration: obs.InventoryConfiguration{
            Id:        "daily-inventory",
            IsEnabled: true,
            Destination: obs.InventoryDestination{
                Format: "CSV",
                Bucket: "report-bucket",
                Prefix: "daily-reports/",
            },
            Schedule: obs.InventorySchedule{
                Frequency: obs.InventoryFrequencyDaily,
            },
            IncludedObjectVersions: "All",
            OptionalFields: &obs.InventoryOptionalFields{
                Fields: []string{
                    obs.InventoryFieldSize,
                    obs.InventoryFieldLastModifiedDate,
                    obs.InventoryFieldETag,
                    obs.InventoryFieldStorageClass,
                },
            },
        },
    }

    output, err := obsClient.SetBucketInventory(input)
    if err != nil {
        fmt.Printf("设置桶清单失败: %v\n", err)
        return
    }

    fmt.Printf("设置桶清单成功，RequestId: %s\n", output.RequestId)
}
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |
| MalformedXML | 400 | XML 格式错误 |
| InvalidArgument | 400 | 参数错误 |

#### 注意事项

1. 清单规则 ID 必须在桶内唯一
2. 目标桶必须存在且具有写入权限
3. 每个桶最多支持 1000 个清单配置
4. 清单报告生成可能需要一定时间

---

### GetBucketInventory

获取指定 ID 的桶清单配置。

#### 方法签名

```go
func (obsClient ObsClient) GetBucketInventory(bucketName, inventoryId string) (output *GetBucketInventoryOutput, err error)
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |
| inventoryId | string | 是 | 清单规则 ID |

#### 返回值

**GetBucketInventoryOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| BaseModel | - | 基础响应信息 |
| InventoryConfiguration | InventoryConfiguration | 清单配置 |

#### 使用示例

```go
output, err := obsClient.GetBucketInventory("my-bucket", "daily-inventory")
if err != nil {
    fmt.Printf("获取桶清单失败: %v\n", err)
    return
}

fmt.Printf("清单 ID: %s\n", output.InventoryConfiguration.Id)
fmt.Printf("启用状态: %v\n", output.InventoryConfiguration.IsEnabled)
fmt.Printf("调度频率: %s\n", output.InventoryConfiguration.Schedule.Frequency)
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |
| NoSuchInventoryConfiguration | 404 | 清单配置不存在 |

---

### ListBucketInventory

列举桶的所有清单配置。

#### 方法签名

```go
func (obsClient ObsClient) ListBucketInventory(bucketName string) (output *ListBucketInventoryOutput, err error)
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

#### 返回值

**ListBucketInventoryOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| BaseModel | - | 基础响应信息 |
| InventoryConfigurations | []InventoryConfiguration | 清单配置列表 |

#### 使用示例

```go
output, err := obsClient.ListBucketInventory("my-bucket")
if err != nil {
    fmt.Printf("列举桶清单失败: %v\n", err)
    return
}

fmt.Printf("共有 %d 个清单配置:\n", len(output.InventoryConfigurations))
for i, config := range output.InventoryConfigurations {
    fmt.Printf("%d. ID: %s, 启用: %v\n", i+1, config.Id, config.IsEnabled)
}
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |

---

### DeleteBucketInventory

删除指定 ID 的桶清单配置。

#### 方法签名

```go
func (obsClient ObsClient) DeleteBucketInventory(bucketName, inventoryId string) (output *BaseModel, err error)
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |
| inventoryId | string | 是 | 清单规则 ID |

#### 返回值

**BaseModel**

| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID |

#### 使用示例

```go
output, err := obsClient.DeleteBucketInventory("my-bucket", "daily-inventory")
if err != nil {
    fmt.Printf("删除桶清单失败: %v\n", err)
    return
}

fmt.Printf("删除桶清单成功，RequestId: %s\n", output.RequestId)
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |
| NoSuchInventoryConfiguration | 404 | 清单配置不存在 |

#### 注意事项

1. 删除清单配置不会删除已生成的清单报告
2. 建议先禁用清单配置，等待当前报告生成完成后再删除

---

## 常量定义

### 调度频率常量

```go
type InventoryFrequencyType string

const (
    InventoryFrequencyDaily  InventoryFrequencyType = "Daily"   // 每日
    InventoryFrequencyWeekly InventoryFrequencyType = "Weekly"  // 每周
)
```

### 子资源常量

```go
const (
    SubResourceInventory = "inventory"  // 清单子资源
)
```

---

## 使用场景

### 场景 1：每日完整清单报告

适用于需要每天了解桶内所有对象变化情况的场景。

```go
input := &obs.SetBucketInventoryInput{
    Bucket: "my-bucket",
    InventoryConfiguration: obs.InventoryConfiguration{
        Id:        "daily-full-inventory",
        IsEnabled: true,
        Destination: obs.InventoryDestination{
            Format: "CSV",
            Bucket: "report-bucket",
            Prefix: "daily-reports/",
        },
        Schedule: obs.InventorySchedule{
            Frequency: obs.InventoryFrequencyDaily,
        },
        IncludedObjectVersions: "All",
    },
}
```

### 场景 2：按前缀筛选的清单报告

适用于需要根据对象前缀分类生成清单报告的场景。

```go
input := &obs.SetBucketInventoryInput{
    Bucket: "my-bucket",
    InventoryConfiguration: obs.InventoryConfiguration{
        Id:        "prefix-inventory",
        IsEnabled: true,
        Destination: obs.InventoryDestination{
            Format: "CSV",
            Bucket: "report-bucket",
            Prefix: "prefix-reports/",
        },
        Schedule: obs.InventorySchedule{
            Frequency: obs.InventoryFrequencyWeekly,
        },
        Filter: &obs.InventoryFilter{
            Prefix: "documents/",
        },
        IncludedObjectVersions: "Current",
    },
}
```

### 场景 3：禁用清单

适用于临时停止清单报告生成的场景。

```go
input := &obs.SetBucketInventoryInput{
    Bucket: "my-bucket",
    InventoryConfiguration: obs.InventoryConfiguration{
        Id:        "daily-inventory",
        IsEnabled: false,  // 设置为 false 来禁用
        Destination: obs.InventoryDestination{
            Format: "CSV",
            Bucket: "report-bucket",
            Prefix: "daily-reports/",
        },
        Schedule: obs.InventorySchedule{
            Frequency: obs.InventoryFrequencyDaily,
        },
    },
}
```

---

## 相关文档

- [OBS API 文档](https://support.huaweicloud.com/api-obs/obs_04_0086.html)
- [桶清单功能说明](https://support.huaweicloud.com/productdesc-obs/obs_04_0017.html)
- [示例代码](../../examples/bucket_inventory_README.md)

---

**文档版本**: 1.0
**更新日期**: 2026-03-06
**SDK 版本**: 3.25.9+
