# 桶清单功能示例

本 README 说明如何使用华为云 OBS Go SDK 的桶清单功能。

## 功能说明

桶清单功能可以定期列举桶内对象，并将对象元数据信息以 CSV 格式存储到指定的桶中。

## API 方法

### SetBucketInventory
设置桶清单配置。

```go
input := &obs.SetBucketInventoryInput{
    Bucket: "your-bucket-name",
    InventoryConfiguration: obs.InventoryConfiguration{
        Id:        "inventory-id",
        IsEnabled: true,
        Destination: obs.InventoryDestination{
            Format: "CSV",
            Bucket: "destination-bucket",
            Prefix: "inventory/",
        },
        Schedule: obs.InventorySchedule{
            Frequency: obs.InventoryFrequencyDaily,
        },
    },
}

output, err := client.SetBucketInventory(input)
```

### GetBucketInventory
获取指定 ID 的桶清单配置。

```go
output, err := client.GetBucketInventory("your-bucket-name", "inventory-id")
```

### ListBucketInventory
列举桶的所有清单配置。

```go
output, err := client.ListBucketInventory("your-bucket-name")
```

### DeleteBucketInventory
删除指定 ID 的桶清单配置。

```go
output, err := client.DeleteBucketInventory("your-bucket-name", "inventory-id")
```

## 数据结构

### InventoryConfiguration
清单配置主结构，包含以下字段：

- **Id**: 清单规则 ID（必填）
- **IsEnabled**: 是否启用清单（必填）
- **Destination**: 清单报告的目标配置
  - **Format**: 报告格式，支持 CSV
  - **Bucket**: 存储报告的目标桶（必填）
  - **Prefix**: 报告对象的前缀
- **Schedule**: 清单生成调度
    - **Frequency**: 调度频率，支持 Daily 或 Weekly
- **Filter**: 对象筛选条件（可选）
  - **Prefix**: 对象名前缀
- **IncludedObjectVersions**: 版本包含策略，支持 All 或 Current
- **OptionalFields**: 可选的元数据字段列表

### InventoryFrequency
清单调度频率：
- **Daily**: 每日生成
- **Weekly**: 每周生成

### InventoryOptionalFields
可选的元数据字段：
- **Size**: 对象大小
- **LastModifiedDate**: 最后修改时间
- **ETag**: 对象 ETag
- **StorageClass**: 存储类型
- **ReplicationStatus**: 复制状态
- **EncryptionStatus**: 加密状态
- **ObjectLockRetainUntilDate**: 对象锁保留到期日期
- **ObjectLockMode**: 对象锁模式

## 使用场景

### 场景 1：每日清单报告

适用于需要每天了解桶内对象变化情况的场景。

```go
input := &obs.SetBucketInventoryInput{
    Bucket: "my-bucket",
    InventoryConfiguration: obs.InventoryConfiguration{
        Id: "daily-inventory",
        IsEnabled: true,
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

output, err := client.SetBucketInventory(input)
```

### 场景 2：按对象类型筛选

适用于需要根据对象类型分别生成清单报告的场景。

```go
input := &obs.SetBucketInventoryInput{
    Bucket: "my-bucket",
    InventoryConfiguration: obs.InventoryConfiguration{
        Id: "type-inventory",
        IsEnabled: true,
        Destination: obs.InventoryDestination{
            Format: "CSV",
            Bucket: "report-bucket",
            Prefix: "type-reports/",
        },
        Filter: &obs.InventoryFilter{
            Prefix: "documents/",
        },
    },
}

output, err := client.SetBucketInventory(input)
```

### 场景 3：获取特定清单

适用于需要查看或修改特定清单配置的场景。

```go
output, err := client.GetBucketInventory("my-bucket", "daily-inventory")
```

### 场景 4：列举所有清单

适用于需要查看和管理所有清单配置的场景。

```go
output, err := client.ListBucketInventory("my-bucket")
for i, config := range output.InventoryConfigurations {
    fmt.Printf("清单 %d: %s (启用: %v)\n", i+1, config.Id, config.IsEnabled)
}
```

### 场景 5：禁用清单

适用于临时停止清单报告生成的场景。

```go
input := &obs.SetBucketInventoryInput{
    Bucket: "my-bucket",
    InventoryConfiguration: obs.InventoryConfiguration{
        Id: "daily-inventory",
        IsEnabled: false,  // 设置为 false 来禁用
    },
}

output, err := client.SetBucketInventory(input)
```

### 场景 6：删除清单

适用于不再需要的清单配置。

```go
output, err := client.DeleteBucketInventory("my-bucket", "daily-inventory")
```

## 注意事项

1. 清单规则 ID 必须唯一，不能重复
2. 清单报告将存储在指定的目标桶中，请确保该桶存在
3. 清单报告可能需要一些时间才能生成，请耐心等待
4. 每个桶最多支持 1000 个清单配置
5. 如果需要删除清单，请先禁用它，然后再删除

## 环境变量

配置示例时需要设置以下环境变量：

```bash
export AccessKeyID=your_access_key_id
export SecretAccessKey=your_secret_access_key
export Endpoint=https://obs.cn-north-4.myhuaweicloud.com
export BucketName=your_source_bucket
export DestinationBucket=your_report_bucket
```

## 完整示例

请参考 `bucket_operations_sample.go` 和 `object_operations_sample.go` 了解完整的 SDK 使用模式。
