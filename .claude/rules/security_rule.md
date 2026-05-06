# Security Rules

约束 OBS SDK 中的凭证提供、签名类型选择、敏感信息处理和临时凭证使用。

---

## SEC-01: 凭证通过 securityProvider 接口提供

**条款**: 禁止在代码中硬编码 AK/SK 凭证，禁止绕过 `securityProvider` 接口直接获取凭证。

### 正确示例

```go
// provider.go
type securityProvider interface {
    getSecurity() securityHolder
}

// 默认实现
bsp := NewBasicSecurityProvider("ak", "sk", "")
client, _ := New("ak", "sk", endpoint, WithSecurityProviders(bsp))

// 环境变量实现
esp := NewEnvSecurityProvider("")
client, _ := New("", "", endpoint, WithSecurityProviders(esp))
```

### 错误示例

```go
// 错误：绕过 securityProvider 直接在代码中硬编码
func getAK() string {
    return "my-hardcoded-access-key"  // 禁止硬编码
}
```

### 验证清单
- [ ] 是否不存在硬编码 AK/SK 凭证的情况
- [ ] 是否不存在绕过 `securityProvider` 接口直接获取凭证的情况
- [ ] 是否不存在不支持通过配置切换 Provider 的情况

---

## SEC-02: 签名类型通过 config.signature 选择

**条款**: 禁止通过 endpoint 格式或其他非标准方式推断签名类型，禁止手动拼接 `x-amz-`/`x-obs-` 头前缀。

### 正确示例

```go
// 使用 OBS 签名
client, _ := obs.New(ak, sk, endpoint, obs.WithSignature(obs.SignatureObs))

// 使用 V4 签名
client, _ := obs.New(ak, sk, endpoint, obs.WithSignature(obs.SignatureV4))

// isObs 在 auth 层判断头前缀
if isObs {
    // 使用 x-obs- 前缀
} else {
    // 使用 x-amz- 前缀
}
```

### 错误示例

```go
// 错误：手动判断签名类型而非通过 config
if strings.HasPrefix(endpoint, "obs") {
    useOBSSignature()
} else {
    useV2Signature()
}
```

### 验证清单
- [ ] 是否不存在通过非标准方式推断签名类型的情况
- [ ] 是否不存在手动拼接 `x-amz-`/`x-obs-` 头前缀的情况
- [ ] 是否不存在未通过 `WithSignature` 配置签名类型的情况

---

## SEC-03: 日志禁止输出 AK/SK/SecurityToken 实际值

**条款**: 日志中涉及 AK/SK/SecurityToken 的输出必须脱敏处理，使用 `"xxxx"` 替代实际值。参考 `EcsSecurityProvider` 中的日志格式。

### 正确示例

```go
// provider.go — 正确的脱敏日志
doLog(LEVEL_INFO, "Get security from ecs succeed, AK:xxxx, SK:xxxx, SecurityToken:xxxx, ExpireDate %s", _sh.expireDate)
```

### 错误示例

```go
// 错误：日志中输出实际凭证值
doLog(LEVEL_INFO, "AK: %s, SK: %s", holder.ak, holder.sk)
doLog(LEVEL_DEBUG, "Using token: %s", holder.securityToken)
```

### 验证清单
- [ ] 日志中是否不存在 AK/SK/SecurityToken 的实际值
- [ ] 凭证相关日志是否使用 `"xxxx"` 替代
- [ ] `config.String()` 方法是否已脱敏（不输出 credential 信息）

---

## SEC-04: 临时凭证通过 securityToken 字段传递

**条款**: 禁止通过 URL 参数或自定义 HTTP 头传递临时安全令牌（必须通过 `securityToken` 字段和标准头）。

### 正确示例

```go
// provider.go — 临时凭证传递
func (bsp *BasicSecurityProvider) refresh(ak, sk, securityToken string) {
    bsp.val.Store(securityHolder{
        ak:            strings.TrimSpace(ak),
        sk:            strings.TrimSpace(sk),
        securityToken: strings.TrimSpace(securityToken),
    })
}

// 创建带临时凭证的客户端
bsp := NewBasicSecurityProvider(ak, sk, securityToken)
client, _ := New(ak, sk, endpoint, WithSecurityProviders(bsp))
```

### 错误示例

```go
// 错误：将 securityToken 放入 URL 参数
url := fmt.Sprintf("%s?token=%s", endpoint, securityToken)

// 错误：自定义头传递 token
headers["X-Custom-Token"] = []string{securityToken}
```

### 验证清单
- [ ] 是否不存在通过 URL 参数传递临时令牌的情况
- [ ] 是否不存在通过自定义 HTTP 头传递临时令牌的情况
- [ ] 是否不存在 HTTP 头未使用标准字段（`x-amz-security-token` / `x-obs-security-token`）的情况
