# 子任务 9.2：实施计划

## 详细实施步骤

### 1. 文件位置
- **目标文件 1**: `obs/trait_bucket.go`
- **目标文件 2**: `obs/client_bucket.go`

### 2. 客户端方法实现

```go
// SetBucketObjectLock sets up bucket-level WORM policy
func (obsClient ObsClient) SetBucketObjectLock(input *SetBucketObjectLockInput, extensions ...extensionOptions) (output *BaseModel, err error) {
    if input == nil {
        return nil, errors.New("SetBucketObjectLockInput is nil")
    }
    if input.Bucket == "" {
        return nil, errors.New("bucket is empty")
    }

    output = &BaseModel{}
    err = obsClient.doActionWithBucket("SetBucketObjectLock", HTTP_PUT, input.Bucket, input, output, extensions)
    if err != nil {
        output = nil
    }
    return
}

// GetBucketObjectLock gets bucket-level WORM policy
func (obsClient ObsClient) GetBucketObjectLock(bucketName string, extensions ...extensionOptions) (output *GetBucketObjectLockOutput, err error) {
    if bucketName == "" {
        return nil, errors.New("bucketName is empty")
    }

    output = &GetBucketObjectLockOutput{}
    err = obsClient.doActionWithBucket("GetBucketObjectLock", HTTP_GET, bucketName,
        newSubResourceSerial(SubResourceObjectLock), output, extensions)
    if err != nil {
        output = nil
    }
    return
}

// DeleteBucketObjectLock deletes bucket-level WORM policy
func (obsClient ObsClient) DeleteBucketObjectLock(bucketName string, extensions ...extensionOptions) (output *BaseModel, err error) {
    if bucketName == "" {
        return nil, errors.New("bucketName is empty")
    }

    output = &BaseModel{}
    err = obsClient.doActionWithBucket("DeleteBucketObjectLock", HTTP_DELETE, bucketName,
        newSubResourceSerial(SubResourceObjectLock), output, extensions)
    if err != nil {
        output = nil
    }
    return
}
```

### 3. 时间估算
- SetBucketObjectLock 实现：20 分钟
- GetBucketObjectLock 实现：15 分钟
- DeleteBucketObjectLock 实现：15 分钟
- 测试和调试：20 分钟
- **总计**: 约 1.2 小时（0.15 天）

## 技术要点

### 方法命名
- 遵循现有命名模式
- 使用描述性名称

### 参数验证
- 验证输入不为 nil
- 验证桶名称不为空
