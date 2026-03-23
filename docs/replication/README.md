# 跨区域复制 API 接口文档

本文档包含华为云 OBS Go SDK 中跨区域复制（Cross-Region Replication）相关的所有 API 接口说明。

## 目录

- [设置跨区域复制规则](#设置跨区域复制规则-setbucketreplication)
- [获取跨区域复制配置](#获取跨区域复制配置-getbucketreplication)
- [删除跨区域复制规则](#删除跨区域复制规则-deletebucketreplication)
- [数据结构定义](#数据结构定义)
- [常量定义](#常量定义)
- [使用场景](#使用场景)

---

## 设置跨区域复制规则

为指定桶设置跨区域复制规则，实现对象数据到目标桶的异步复制。

### 方法签名

```go
func (obsClient ObsClient) SetBucketReplication(input *SetBucketReplicationInput, extensions ...extensionOptions) (output *BaseModel, err error)
```

### 参数说明

**SetBucketReplicationInput**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 桶名 |
| ReplicationConfiguration | ReplicationConfiguration | 是 | 复制配置 |

**ReplicationConfiguration**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Rules | []ReplicationRule | 是 | 复制规则列表，最多 100 条 |

**ReplicationRule**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ID | string | 是 | 规则 ID，同一桶内唯一 |
| Status | RuleStatusType | 是 | 规则状态，Enabled 或 Disabled |
| Prefix | *ReplicationPrefix | 否 | 前缀配置，指定要复制的对象前缀 |
| Destination | ReplicationDestination | 是 | 目标桶配置 |
| HistoricalObjectReplication | string | 否 | 历史对象复制，Enabled 或 Disabled |

**ReplicationPrefix**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| PrefixSet | PrefixSet | 是 | 前缀集合 |

**PrefixSet**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Prefixes | []string | 是 | 对象前缀列表 |

**ReplicationDestination**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Bucket | string | 是 | 目标桶名 |
| StorageClass | StorageClassType | 否 | 目标存储类型 |
| Location | string | 否 | 目标桶所在区域 |

### 返回值

**BaseModel**

| 字段 | 类型 | 说明 |
|------|------|------|
| StatusCode | int | HTTP 状态码 |
| RequestId | string | 请求 ID |
| Headers | map[string][]string | 响应头 |

### 使用示例

```go
package main

import (
	"fmt"
	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func main() {
	// 创建 OBS 客户端
	ak := "your-access-key"
	sk := "your-secret-key"
	endpoint := "https://obs.cn-north-4.myhuaweicloud.com"

	config := &obs.Configuration{
		Endpoint:        endpoint,
		AccessKeyID:     ak,
		SecretAccessKey: sk,
		Signature:       obs.SignatureObs, // 跨区域复制仅支持 OBS 签名
	}

	client, err := obs.New(config)
	if err != nil {
		panic(err)
	}

	// 设置跨区域复制规则
	input := &obs.SetBucketReplicationInput{
		Bucket: "source-bucket",
		ReplicationConfiguration: obs.ReplicationConfiguration{
			Rules: []obs.ReplicationRule{
				{
					ID:     "rule-1",
					Status: obs.RuleStatusEnabled,
					Prefix: &obs.ReplicationPrefix{
						PrefixSet: obs.PrefixSet{
							Prefixes: []string{"images/", "videos/"},
						},
					},
					Destination: obs.ReplicationDestination{
						Bucket:       "dest-bucket",
						StorageClass: obs.StorageClassStandard,
						Location:     "cn-south-1",
					},
					HistoricalObjectReplication: "Enabled",
				},
			},
		},
	}

	output, err := client.SetBucketReplication(input)
	if err != nil {
		fmt.Printf("设置跨区域复制规则失败: %v\n", err)
		return
	}

	fmt.Printf("设置成功，RequestId: %s\n", output.RequestId)
}
```

### 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| InvalidBucketName | 400 | 桶名无效 |
| MalformedXML | 400 | XML 格式错误 |
| InvalidRequest | 400 | 请求参数无效 |
| NoSuchBucket | 404 | 桶不存在 |
| Unauthorized | 403 | 未授权 |

### 注意事项

1. 跨区域复制功能仅支持 OBS 签名，不支持 AWS 签名
2. 同一桶内最多可配置 100 条复制规则
3. 配置大小限制为 50KB
4. 源桶和目标桶必须都已存在且版本控制状态相同
5. 目标桶必须与源桶位于不同的区域

---

## 获取跨区域复制配置

获取指定桶的跨区域复制配置。

### 方法签名

```go
func (obsClient ObsClient) GetBucketReplication(bucketName string, extensions ...extensionOptions) (output *GetBucketReplicationOutput, err error)
```

### 参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bucketName | string | 是 | 桶名 |

### 返回值

**GetBucketReplicationOutput**

| 字段 | 类型 | 说明 |
|------|------|------|
| BaseModel | - | 基础响应模型 |
| ReplicationConfiguration | ReplicationConfiguration | 复制配置 |

### 使用示例

```go
// 获取跨区域复制配置
output, err := client.GetBucketReplication("source-bucket")
if err != nil {
	fmt.Printf("获取跨区域复制配置失败: %v\n", err)
	return
}

// 遍历复制规则
for _, rule := range output.ReplicationConfiguration.Rules {
	fmt.Printf("规则 ID: %s\n", rule.ID)
	fmt.Printf("规则状态: %s\n", rule.Status)
	fmt.Printf("目标桶: %s\n", rule.Destination.Bucket)
}
```

---

## 删除跨区域复制规则

删除指定桶的跨区域复制配置。

### 方法签名

```go
func (obsClient ObsClient) DeleteBucketReplication(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error)
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
| RequestId | string | 请求 ID |
| Headers | map[string][]string | 响应头 |

### 使用示例

```go
// 删除跨区域复制配置
output, err := client.DeleteBucketReplication("source-bucket")
if err != nil {
	fmt.Printf("删除跨区域复制配置失败: %v\n", err)
	return
}

fmt.Printf("删除成功，RequestId: %s\n", output.RequestId)
```

---

## 数据结构定义

```go
// 复制配置
type ReplicationConfiguration struct {
	XMLName xml.Name          `xml:"ReplicationConfiguration"`
	Rules   []ReplicationRule `xml:"Rule"`
}

// 复制规则
type ReplicationRule struct {
	XMLName                     xml.Name               `xml:"Rule"`
	ID                          string                 `xml:"Id,omitempty"`
	Status                      RuleStatusType         `xml:"Status"`
	Prefix                      *ReplicationPrefix     `xml:"Prefix,omitempty"`
	Destination                 ReplicationDestination `xml:"Destination"`
	HistoricalObjectReplication string                 `xml:"HistoricalObjectReplication,omitempty"`
}

// 复制前缀配置
type ReplicationPrefix struct {
	PrefixSet PrefixSet `xml:"PrefixSet"`
}

// 前缀集合
type PrefixSet struct {
	Prefixes []string `xml:"Prefix"`
}

// 目标桶配置
type ReplicationDestination struct {
	XMLName      xml.Name         `xml:"Destination"`
	Bucket       string           `xml:"Bucket"`
	StorageClass StorageClassType `xml:"StorageClass,omitempty"`
	Location     string           `xml:"Location,omitempty"`
}
```

---

## 常量定义

```go
// 规则状态
type RuleStatusType string

const (
	RuleStatusEnabled  RuleStatusType = "Enabled"  // 启用
	RuleStatusDisabled RuleStatusType = "Disabled" // 禁用
)

// 子资源类型
const (
	SubResourceReplication SubResourceType = "replication"
)
```

---

## 使用场景

### 场景 1：按前缀复制对象

复制指定前缀的所有对象到目标桶。

```go
input := &obs.SetBucketReplicationInput{
	Bucket: "source-bucket",
	ReplicationConfiguration: obs.ReplicationConfiguration{
		Rules: []obs.ReplicationRule{
			{
				ID:     "prefix-rule",
				Status: obs.RuleStatusEnabled,
				Prefix: &obs.ReplicationPrefix{
					PrefixSet: obs.PrefixSet{
						Prefixes: []string{"logs/", "backup/"},
					},
				},
				Destination: obs.ReplicationDestination{
					Bucket: "dest-bucket",
				},
			},
		},
	},
}
```

### 场景 2：指定存储类型复制

复制对象时指定目标存储类型。

```go
input := &obs.SetBucketReplicationInput{
	Bucket: "source-bucket",
	ReplicationConfiguration: obs.ReplicationConfiguration{
		Rules: []obs.ReplicationRule{
			{
				ID:     "storage-class-rule",
				Status: obs.RuleStatusEnabled,
				Destination: obs.ReplicationDestination{
					Bucket:       "dest-bucket",
					StorageClass: obs.StorageClassWarm, // 温存储
				},
			},
		},
	},
}
```

### 场景 3：多规则配置

配置多条复制规则实现复杂的复制策略。

```go
input := &obs.SetBucketReplicationInput{
	Bucket: "source-bucket",
	ReplicationConfiguration: obs.ReplicationConfiguration{
		Rules: []obs.ReplicationRule{
			{
				ID:     "rule-1",
				Status: obs.RuleStatusEnabled,
				Destination: obs.ReplicationDestination{
					Bucket: "dest-bucket-1",
				},
			},
			{
				ID:     "rule-2",
				Status: obs.RuleStatusEnabled,
				Destination: obs.ReplicationDestination{
					Bucket: "dest-bucket-2",
				},
			},
		},
	},
}
```

---

## 相关文档

- [华为云 OBS 跨区域复制文档](https://support.huaweicloud.com/obs/productdesc-obs/obs_04_0063.html)
- [API 参考文档](https://support.huaweicloud.com/api-obs/obs_04_0084.html)
- [示例代码](../../examples/)
