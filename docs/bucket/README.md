# 桶相关 API 接口文档

本文档包含华为云 OBS Go SDK 中所有桶操作相关的 API 接口说明。

## 目录

- [获取桶存量信息](#获取桶存量信息)
- [桶清单管理](#桶清单管理)
- [跨区域复制](#跨区域复制)
- [归档存储对象直读](#归档存储对象直读)

---

## 获取桶存量信息

获取桶中的对象个数及对象占用空间，包括标准存储、低频存储和归档存储的详细信息。

### GetBucketStorageInfo

获取指定桶的存量信息，包括对象数量和存储空间大小，并按存储类型分类统计。

#### 方法签名

```go
func (obsClient ObsClient) GetBucketStorageInfo(bucketName string, extensions ...extensionOptions) (output *GetBucketStorageInfoOutput, err error)
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

#### 返回值

**GetBucketStorageInfoOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| Location | string | 桶所在区域 |
| Bucket | string | 桶名 |
| Size | int64 | 桶中对象占用的总存储空间（字节） |
| ObjectNumber | int | 桶中的总对象个数 |
| Standard | StorageTypeDetail | 标准存储的详细信息 |
| Warm | StorageTypeDetail | 低频存储的详细信息 |
| Cold | StorageTypeDetail | 归档存储的详细信息 |

**StorageTypeDetail**

| 字段 | 类型 | 说明 |
|------|------|------|
| ObjectNumber | int64 | 该存储类型的对象数量 |
| Size | int64 | 该存储类型的占用空间（字节） |

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

    // 获取桶存量信息
    output, err := obsClient.GetBucketStorageInfo("my-bucket")
    if err != nil {
        fmt.Printf("获取桶存量信息失败: %v\n", err)
        return
    }

    fmt.Printf("桶存量信息：\n")
    fmt.Printf("  桶名: %s\n", output.Bucket)
    fmt.Printf("  区域: %s\n", output.Location)
    fmt.Printf("  总对象数: %d\n", output.ObjectNumber)
    fmt.Printf("  总存储空间: %d 字节\n", output.Size)

    // 按存储类型显示详细信息
    fmt.Printf("  标准存储: %d 对象, %d 字节\n", output.Standard.ObjectNumber, output.Standard.Size)
    fmt.Printf("  低频存储: %d 对象, %d 字节\n", output.Warm.ObjectNumber, output.Warm.Size)
    fmt.Printf("  归档存储: %d 对象, %d 字节\n", output.Cold.ObjectNumber, output.Cold.Size)
}
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |

#### 注意事项

1. 桶名必须符合 OBS 桶名命名规范
2. 需要具有桶的读取权限
3. 存量信息数据可能存在一定的延迟，不是实时数据
4. 对象数量和存储空间为统计值，可能存在一定的误差

#### 使用场景

**场景 1：监控桶存储使用情况**

适用于需要实时了解桶存储容量和对象数量的场景。

```go
output, err := obsClient.GetBucketStorageInfo("my-bucket")
if err != nil {
    log.Printf("获取存量信息失败: %v", err)
    return
}

// 计算存储使用率
totalSize := output.Size
quota := 10 * 1024 * 1024 * 1024 // 10GB
usageRate := float64(totalSize) / float64(quota) * 100

log.Printf("存储使用率: %.2f%%", usageRate)
if usageRate > 80 {
    log.Printf("警告：存储空间即将用尽")
}
```

**场景 2：按存储类型分析存储分布**

适用于需要分析不同存储类型对象分布的场景。

```go
output, err := obsClient.GetBucketStorageInfo("my-bucket")
if err != nil {
    log.Printf("获取存量信息失败: %v", err)
    return
}

// 统计各存储类型的占比
standardRatio := float64(output.Standard.Size) / float64(output.Size) * 100
warmRatio := float64(output.Warm.Size) / float64(output.Size) * 100
coldRatio := float64(output.Cold.Size) / float64(output.Size) * 100

log.Printf("存储类型分布：")
log.Printf("  标准存储: %.2f%%", standardRatio)
log.Printf("  低频存储: %.2f%%", warmRatio)
log.Printf("  归档存储: %.2f%%", coldRatio)
```

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

#### 使用场景

**场景 1：每日完整清单报告**

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

**场景 2：按前缀筛选的清单报告**

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

#### 使用场景

**场景：查询并判断清单配置**

适用于需要查询特定清单规则配置并进行条件判断的场景。

```go
output, err := obsClient.GetBucketInventory("my-bucket", "daily-inventory")
if err != nil {
    return
}

// 判断清单是否启用
if output.InventoryConfiguration.IsEnabled {
    fmt.Println("清单已启用")
    fmt.Printf("报告存储在: %s/%s\n",
        output.InventoryConfiguration.Destination.Bucket,
        output.InventoryConfiguration.Destination.Prefix)
} else {
    fmt.Println("清单已禁用")
}
```

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

#### 使用场景

**场景：批量管理清单配置**

适用于需要查看所有清单规则并进行批量管理的场景。

```go
output, err := obsClient.ListBucketInventory("my-bucket")
if err != nil {
    return
}

// 统计启用和禁用的清单数量
enabledCount := 0
disabledCount := 0

for _, config := range output.InventoryConfigurations {
    if config.IsEnabled {
        enabledCount++
    } else {
        disabledCount++
    }
}

fmt.Printf("启用的清单: %d\n", enabledCount)
fmt.Printf("禁用的清单: %d\n", disabledCount)
```

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

#### 使用场景

**场景：安全删除清单配置**

适用于需要在删除前确认清单配置存在且已禁用的场景。

```go
// 先查询清单配置
output, err := obsClient.GetBucketInventory("my-bucket", "daily-inventory")
if err != nil {
    fmt.Printf("清单配置不存在或获取失败: %v\n", err)
    return
}

// 检查是否已禁用
if output.InventoryConfiguration.IsEnabled {
    fmt.Println("请先禁用清单配置再删除")
    return
}

// 删除清单配置
_, err = obsClient.DeleteBucketInventory("my-bucket", "daily-inventory")
if err != nil {
    fmt.Printf("删除失败: %v\n", err)
    return
}

fmt.Println("清单配置已安全删除")
```

---

## 跨区域复制

跨区域复制功能允许在不同区域之间自动复制对象，实现数据的异地容灾。

### SetBucketReplication

设置桶的跨区域复制配置。

#### 方法签名

```go
func (obsClient ObsClient) SetBucketReplication(input *SetBucketReplicationInput, extensions ...extensionOptions) (output *BaseModel, err error)
```

#### 参数说明

**SetBucketReplicationInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 源桶名 |
| Agency | string | 是 | IAM 委托名称 |
| Rules | []ReplicationRule | 是 | 复制规则列表 |

**ReplicationRule**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ID | string | 否 | 规则 ID |
| Prefix | string | 是 | 对象前缀 |
| Status | string | 是 | 规则状态，Enabled 或 Disabled |
| Destination | ReplicationDestination | 是 | 目标配置 |
| HistoricalObjectReplication | string | 否 | 历史对象复制 |

**ReplicationDestination**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 目标桶名 |
| StorageClass | string | 否 | 目标存储类型 |
| DeleteData | string | 否 | 删除数据同步 |

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

    // 设置跨区域复制配置
    input := &obs.SetBucketReplicationInput{
        Bucket: "source-bucket",
        Agency: "your-agency-name",
        Rules: []obs.ReplicationRule{
            {
                ID:    "rule-1",
                Prefix: "documents/",
                Status: string(obs.ReplicationStatusEnabled),
                Destination: obs.ReplicationDestination{
                    Bucket:       "dest-bucket",
                    StorageClass: "STANDARD",
                },
            },
        },
    }

    output, err := obsClient.SetBucketReplication(input)
    if err != nil {
        fmt.Printf("设置跨区域复制失败: %v\n", err)
        return
    }

    fmt.Printf("设置跨区域复制成功，RequestId: %s\n", output.RequestId)
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

1. 目标桶必须存在且在不同区域
2. 必须先创建 IAM 委托并授权
3. 复制规则的前缀不能重叠
4. 同一桶内最多支持 100 条复制规则

#### 使用场景

**场景 1：跨区域数据容灾**

适用于需要在不同区域之间复制数据以实现容灾的场景。

```go
input := &obs.SetBucketReplicationInput{
    Bucket: "source-bucket",
    Agency: "disaster-recovery-agency",
    Rules: []obs.ReplicationRule{
        {
            ID:    "dr-rule",
            Prefix: "",
            Status: string(obs.ReplicationStatusEnabled),
            Destination: obs.ReplicationDestination{
                Bucket:       "dest-bucket",
                StorageClass: "STANDARD",
            },
        },
    },
}
```

**场景 2：按前缀选择性复制**

适用于只需要复制特定前缀对象的场景。

```go
input := &obs.SetBucketReplicationInput{
    Bucket: "source-bucket",
    Agency: "selective-replication-agency",
    Rules: []obs.ReplicationRule{
        {
            ID:    "important-docs",
            Prefix: "important/",
            Status: string(obs.ReplicationStatusEnabled),
            Destination: obs.ReplicationDestination{
                Bucket:       "backup-bucket",
                StorageClass: "STANDARD_IA",
            },
        },
    },
}
```

---

### GetBucketReplication

获取桶的跨区域复制配置。

#### 方法签名

```go
func (obsClient ObsClient) GetBucketReplication(input *GetBucketReplicationInput) (output *GetBucketReplicationOutput, err error)
```

#### 参数说明

**GetBucketReplicationInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名 |

#### 返回值

**GetBucketReplicationOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| BaseModel | - | 基础响应信息 |
| Agency | string | IAM 委托名称 |
| Rules | []ReplicationRule | 复制规则列表 |

#### 使用示例

```go
output, err := obsClient.GetBucketReplication(&obs.GetBucketReplicationInput{
    Bucket: "source-bucket",
})
if err != nil {
    fmt.Printf("获取跨区域复制配置失败: %v\n", err)
    return
}

fmt.Printf("IAM 委托: %s\n", output.Agency)
fmt.Printf("复制规则数量: %d\n", len(output.Rules))
for i, rule := range output.Rules {
    fmt.Printf("%d. 规则 ID: %s, 前缀: %s, 状态: %s\n",
        i+1, rule.ID, rule.Prefix, rule.Status)
}
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |
| ReplicationConfigurationNotFoundError | 404 | 复制配置不存在 |

#### 使用场景

**场景：分析复制规则状态**

适用于需要查看所有复制规则并分析其状态的场景。

```go
output, err := obsClient.GetBucketReplication(&obs.GetBucketReplicationInput{
    Bucket: "source-bucket",
})
if err != nil {
    return
}

// 分析复制规则状态
enabledRules := 0
disabledRules := 0

for _, rule := range output.Rules {
    if rule.Status == string(obs.ReplicationStatusEnabled) {
        enabledRules++
    } else {
        disabledRules++
    }
}

fmt.Printf("启用的规则: %d\n", enabledRules)
fmt.Printf("禁用的规则: %d\n", disabledRules)
fmt.Printf("目标桶: ")
for i, rule := range output.Rules {
    if i > 0 {
        fmt.Print(", ")
    }
    fmt.Printf("%s (%s)", rule.Destination.Bucket, rule.Prefix)
}
fmt.Println()
```

---

### DeleteBucketReplication

删除桶的跨区域复制配置。

#### 方法签名

```go
func (obsClient ObsClient) DeleteBucketReplication(input *DeleteBucketReplicationInput, extensions ...extensionOptions) (output *BaseModel, err error)
```

#### 参数说明

**DeleteBucketReplicationInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名 |

#### 返回值

**BaseModel**

| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID |

#### 使用示例

```go
output, err := obsClient.DeleteBucketReplication(&obs.DeleteBucketReplicationInput{
    Bucket: "source-bucket",
})
if err != nil {
    fmt.Printf("删除跨区域复制配置失败: %v\n", err)
    return
}

fmt.Printf("删除跨区域复制配置成功，RequestId: %s\n", output.RequestId)
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |

#### 注意事项

1. 删除复制配置会停止所有复制任务
2. 已复制到目标桶的数据不会自动删除
3. 建议先确认所有复制任务完成后再删除配置

#### 使用场景

**场景：安全删除复制配置**

适用于需要在删除前确认复制配置状态并提醒用户的场景。

```go
// 先查询复制配置
output, err := obsClient.GetBucketReplication(&obs.GetBucketReplicationInput{
    Bucket: "source-bucket",
})
if err != nil {
    fmt.Printf("复制配置不存在或获取失败: %v\n", err)
    return
}

// 显示将要删除的配置
fmt.Println("即将删除以下跨区域复制配置:")
fmt.Printf("IAM 委托: %s\n", output.Agency)
for i, rule := range output.Rules {
    fmt.Printf("%d. 源前缀: %s -> 目标桶: %s\n",
        i+1, rule.Prefix, rule.Destination.Bucket)
}

fmt.Println("\n注意: 已复制的数据将保留在目标桶中")

// 删除复制配置
_, err = obsClient.DeleteBucketReplication(&obs.DeleteBucketReplicationInput{
    Bucket: "source-bucket",
})
if err != nil {
    fmt.Printf("删除失败: %v\n", err)
    return
}

fmt.Println("跨区域复制配置已删除")
```

---

## 归档存储对象直读

归档对象直读功能允许用户直接下载归档存储的对象，而无需先进行恢复操作。这可以简化归档对象的访问流程，提高访问效率。

### SetBucketDirectColdAccess

设置桶的归档对象直读配置。

#### 方法签名

```go
func (obsClient ObsClient) SetBucketDirectColdAccess(input *SetBucketDirectColdAccessInput, extensions ...extensionOptions) (output *BaseModel, err error)
```

#### 参数说明

**SetBucketDirectColdAccessInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名 |
| Enabled | bool | 是 | 是否启用归档对象直读 |

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

    // 启用归档对象直读
    input := &obs.SetBucketDirectColdAccessInput{
        Bucket:  "my-bucket",
        Enabled: true,
    }

    output, err := obsClient.SetBucketDirectColdAccess(input)
    if err != nil {
        fmt.Printf("设置归档直读失败: %v\n", err)
        return
    }

    fmt.Printf("设置归档直读成功，RequestId: %s\n", output.RequestId)
}
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |

#### 注意事项

1. 启用归档直读后，归档对象可以直接下载，无需恢复
2. 归档直读可能会产生额外的费用，请根据实际业务需求选择是否启用
3. 只有桶的所有者或具有相应权限的用户才能设置归档直读

#### 使用场景

**场景 1：启用归档直读以提高访问效率**

适用于需要频繁访问归档对象但希望避免恢复操作的场景。

```go
input := &obs.SetBucketDirectColdAccessInput{
    Bucket:  "archive-bucket",
    Enabled: true,
}

output, err := obsClient.SetBucketDirectColdAccess(input)
if err != nil {
    log.Printf("启用归档直读失败: %v", err)
    return
}
log.Printf("归档直读已启用")
```

**场景 2：禁用归档直读以降低成本**

适用于需要控制成本或不常访问归档对象的场景。

```go
input := &obs.SetBucketDirectColdAccessInput{
    Bucket:  "archive-bucket",
    Enabled: false,
}

output, err := obsClient.SetBucketDirectColdAccess(input)
if err != nil {
    log.Printf("禁用归档直读失败: %v", err)
    return
}
log.Printf("归档直读已禁用")
```

---

### GetBucketDirectColdAccess

获取桶的归档对象直读配置。

#### 方法签名

```go
func (obsClient ObsClient) GetBucketDirectColdAccess(bucketName string, extensions ...extensionOptions) (output *GetBucketDirectColdAccessOutput, err error)
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

#### 返回值

**GetBucketDirectColdAccessOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| BaseModel | - | 基础响应信息 |
| Enabled | bool | 是否启用归档对象直读 |

#### 使用示例

```go
output, err := obsClient.GetBucketDirectColdAccess("my-bucket")
if err != nil {
    fmt.Printf("获取归档直读配置失败: %v\n", err)
    return
}

fmt.Printf("归档直读状态: %v\n", output.Enabled)
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |

#### 使用场景

**场景：查询并判断归档直读状态**

适用于需要查询归档直读配置并进行条件判断的场景。

```go
output, err := obsClient.GetBucketDirectColdAccess("my-bucket")
if err != nil {
    log.Printf("获取归档直读配置失败: %v", err)
    return
}

// 判断归档直读是否启用
if output.Enabled {
    log.Println("归档直读已启用，归档对象可直接访问")
} else {
    log.Println("归档直读未启用，归档对象需要先恢复")
}
```

---

### DeleteBucketDirectColdAccess

删除桶的归档对象直读配置。

#### 方法签名

```go
func (obsClient ObsClient) DeleteBucketDirectColdAccess(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error)
```

#### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

#### 返回值

**BaseModel**

| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID |

#### 使用示例

```go
output, err := obsClient.DeleteBucketDirectColdAccess("my-bucket")
if err != nil {
    fmt.Printf("删除归档直读配置失败: %v\n", err)
    return
}

fmt.Printf("删除归档直读配置成功，RequestId: %s\n", output.RequestId)
```

#### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| AccessDenied | 403 | 权限不足 |
| NoSuchBucket | 404 | 桶不存在 |

#### 注意事项

1. 删除归档直读配置后，归档对象需要先恢复才能下载
2. 已下载的对象不受影响
3. 建议在删除前确认所有需要的访问已完成

#### 使用场景

**场景：安全删除归档直读配置**

适用于需要在删除前确认归档直读配置状态并提醒用户的场景。

```go
// 先查询归档直读配置
output, err := obsClient.GetBucketDirectColdAccess("my-bucket")
if err != nil {
    log.Printf("归档直读配置不存在或获取失败: %v", err)
    return
}

// 检查是否已启用
if output.Enabled {
    log.Println("注意：删除归档直读配置后，归档对象将需要先恢复才能下载")
}

// 删除归档直读配置
_, err = obsClient.DeleteBucketDirectColdAccess("my-bucket")
if err != nil {
    log.Printf("删除失败: %v", err)
    return
}

log.Println("归档直读配置已删除")
```

---

## 常量定义

### 桶清单相关常量

#### 调度频率常量

```go
type InventoryFrequencyType string

const (
    InventoryFrequencyDaily  InventoryFrequencyType = "Daily"   // 每日
    InventoryFrequencyWeekly InventoryFrequencyType = "Weekly"  // 每周
)
```

#### 子资源常量

```go
const (
    SubResourceInventory = "inventory"  // 清单子资源
)
```

### 跨区域复制相关常量

#### 复制状态常量

```go
type ReplicationStatusType string

const (
    ReplicationStatusEnabled  ReplicationStatusType = "Enabled"   // 启用
    ReplicationStatusDisabled ReplicationStatusType = "Disabled"  // 禁用
)
```

#### 子资源常量

```go
const (
    SubResourceReplication = "replication"  // 跨区域复制子资源
)
```

#### 归档直读相关子资源常量
```go
const (
    SubResourceDirectcoldaccess = "directcoldaccess"  // 归档对象直读子资源
)
```

---

## 相关文档

- [OBS API 文档](https://support.huaweicloud.com/api-obs/obs_04_0086.html)
- [桶清单功能说明](https://support.huaweicloud.com/productdesc-obs/obs_04_0017.html)
- [跨区域复制功能说明](https://support.huaweicloud.com/productdesc-obs/obs_04_0033.html)
- [示例代码](../../examples/)

---

**文档版本**: 1.2
**更新日期**: 2026-03-07
**SDK 版本**: 3.25.9+
