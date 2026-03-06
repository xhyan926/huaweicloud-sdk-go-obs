# 子任务 2.4：实施计划

## 详细实施步骤

### 1. 使用 go-sdk-ut skill

**必须调用** `/go-sdk-ut` skill 来编写测试

### 2. 测试文件位置
- **目标文件**: `obs/post_policy_test.go`
- **目标文件**: `obs/client_object_test.go`

### 3. 测试用例结构

#### 3.1 Policy 生成测试

```go
func TestBuildPostPolicyJSON_ShouldGenerateValidPolicy_GivenSimpleCondition(t *testing.T) {
    policy := &PostPolicy{
        Expiration: "2026-03-05T10:30:00.000Z",
        Conditions: []PostPolicyCondition{
            CreateBucketCondition("test-bucket"),
            CreateKeyCondition("uploads/"),
        },
    }

    jsonStr, err := buildPostPolicyJSON(policy)

    assert.NoError(t, err)
    assert.Contains(t, jsonStr, "expiration")
    assert.Contains(t, jsonStr, "conditions")
}
```

#### 3.2 签名计算测试

```go
func TestCalculatePostPolicySignature_ShouldCalculateCorrectSignature_GivenValidInput(t *testing.T) {
    policyJSON := `{"expiration":"2026-03-05T10:30:00.000Z","conditions":[...]}`
    secretKey := "testSecretKey123"

    signature, err := CalculatePostPolicySignature(policyJSON, secretKey)

    assert.NoError(t, err)
    assert.NotEmpty(t, signature)
    // 可以验证签名是否符合预期（如果已知的测试向量）
}
```

#### 3.3 集成测试

```go
func TestCreatePostPolicy_ShouldGenerateCompleteToken_GivenValidInput(t *testing.T) {
    input := &CreatePostPolicyInput{
        Bucket:  "test-bucket",
        Key:     "uploads/test.jpg",
        Expires: 3600, // 1 小时后过期
        Conditions: []PostPolicyCondition{
            CreateContentLengthRangeCondition(0, 10485760), // 0-10MB
        },
    }

    output, err := client.CreatePostPolicy(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
    assert.NotEmpty(t, output.Policy)
    assert.NotEmpty(t, output.Signature)
    assert.NotEmpty(t, output.Token)
    assert.NotEmpty(t, output.AccessKeyId)

    // 验证 Token 格式
    assert.Regexp(t, `^[^:]+:[^:]+:[^:]+$`, output.Token)
}
```

#### 3.4 边界条件测试

```go
func TestCreatePostPolicy_ShouldHandleZeroExpiration_GivenInput(t *testing.T) {
    input := &CreatePostPolicyInput{
        Bucket:  "test-bucket",
        Key:     "test.txt",
        Expires: 0, // 立即过期
    }

    output, err := client.CreatePostPolicy(input)

    assert.Error(t, err)
    assert.Nil(t, output)
}

func TestCreatePostPolicy_ShouldHandleLongExpiration_GivenInput(t *testing.T) {
    input := &CreatePostPolicyInput{
        Bucket:  "test-bucket",
        Key:     "test.txt",
        Expires: 86400 * 7, // 7 天后过期
    }

    output, err := client.CreatePostPolicy(input)

    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

### 4. 时间估算
- Policy 生成测试：30 分钟
- 签名计算测试：30 分钟
- 集成测试：30 分钟
- 边界条件测试：30 分钟
- 测试覆盖率优化：30 分钟
- **总计**: 约 2.5 小时（0.31 天）

## 技术要点

### 测试策略
- 从简单到复杂
- 先测试单元功能
- 再测试集成场景
- 最后测试边界条件

### BDD 风格命名
- Test<功能>_Should<预期结果>_When<条件>
- 描述测试意图
- 易于维护

### 签名验证
- 如果有已知测试向量，使用它们验证
- 否则验证签名过程正确
- 验证 HMAC 使用正确

### Token 格式验证
- 使用正则表达式验证格式
- 验证组件分隔符
- 验证所有组件存在
