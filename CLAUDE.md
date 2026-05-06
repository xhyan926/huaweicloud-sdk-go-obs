# CLAUDE.md

华为云 OBS Go SDK — 与华为云 OBS 服务交互的 Go 客户端库，兼容 S3 API。

## 快速参考

- **包路径**: `obs/`
- **入口**: `obs.New(ak, sk, endpoint, configurers...)`
- **架构**: Client → Trait → HTTP → Model（四层单向调用）
- **签名**: SignatureObs / SignatureV2 / SignatureV4
- **配置模式**: 函数式选项 WithXXX

## 运行示例

```bash
cd main && go run obs_go_sample.go
```

## 测试命令

```bash
go test -tags unit ./obs -v              # 单元测试
go test -tags integration ./obs/test/integration -v  # 集成测试
go test -tags perf ./obs -bench=. -benchtime=1s       # 性能测试
go test -tags fuzz ./obs -fuzz=.                       # 模糊测试
```

## 测试技能

- `/go-sdk-ut`: 单元测试编写指南
- `/go-sdk-integration`: 集成测试编写指南
- `/go-sdk-perf`: 性能测试编写指南
- `/go-sdk-fuzz`: 模糊测试编写指南

## 规则索引

所有开发规则位于 `.claude/rules/`，无优先级区分，全部同等重要。

| 规则文件 | 规则数 | 约束范围 |
|---------|--------|---------|
| `architecture_rule.md` | 5 | 分层调用、文件职责、依赖方向 |
| `naming_rule.md` | 6 | 类型、函数、变量、常量、文件命名 |
| `api_rule.md` | 5 | 方法签名、nil 检查、返回语义、参数选择 |
| `model_rule.md` | 5 | 结构体定义、BaseModel 嵌入、序列化标签 |
| `error_rule.md` | 4 | ObsError 使用、客户端校验错误、日志错误 |
| `testing_rule.md` | 6 | BDD 命名、build tag、断言工具、用例合并 |
| `security_rule.md` | 4 | 凭证提供者、签名类型、敏感信息处理 |
| `extension_rule.md` | 4 | 函数式选项、header 扩展、进度回调 |
