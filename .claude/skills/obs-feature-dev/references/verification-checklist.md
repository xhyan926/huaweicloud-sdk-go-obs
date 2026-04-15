# 构建验证和质量检查清单

本文档提供 OBS SDK 代码提交前的完整验证命令和质量检查清单。

---

## Maven 命令规范

**所有 Maven 命令必须包含 settings 参数和 pom 文件路径：**

```bash
mvn [command] -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
```

**禁止**使用不带 settings 参数的简化命令。

---

## 验证步骤（按顺序执行）

### 步骤 1: 编译验证

```bash
mvn clean compile test-compile -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
```

**检查点**：
- 无编译错误
- 无编译警告（如 raw types、unchecked）
- 新增文件被正确编译

### 步骤 2: 单元测试执行

```bash
# 运行指定测试类
mvn test -Dtest=<Feature>UnitTest -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C

# 运行全部测试
mvn test -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
```

**检查点**：
- 所有测试通过
- 无测试失败或错误
- 无测试跳过（除非有合理原因）

### 步骤 3: 静态分析

```bash
mvn checkstyle:check spotbugs:check pmd:check -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
```

**检查点**：
- Checkstyle 检查通过
- SpotBugs 无缺陷
- PMD 无告警

### 步骤 4: 完整构建验证

```bash
mvn clean verify -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
```

**检查点**：
- 所有阶段通过
- 构建产物正确生成

---

## 代码审查检查清单

### 代码质量

- [ ] 代码符合 SOLID 原则
- [ ] 使用 Java 泛型，无原始类型（`List` → `List<String>`）
- [ ] 公共方法有 JavaDoc 注释
- [ ] 异常处理合理，不吞异常
- [ ] 无硬编码敏感信息（密码、密钥）
- [ ] 无未使用的 import 和变量

### 架构合规

- [ ] Request 类正确继承 `BaseBucketRequest`
- [ ] Result 类正确继承 `HeaderResponse`
- [ ] `httpMethod` 在实例初始化块中设置
- [ ] SpecialParamEnum 注册了子资源参数
- [ ] IObsClient 接口方法签名正确
- [ ] AbstractBucketAdvanceClient 使用 `doActionWithResult` + `ActionCallbackWithResult`
- [ ] ObsBucketAdvanceService 使用正确的 REST 方法（`performRestGet/Put/Delete`）

### 测试质量

- [ ] 测试方法命名遵循 `should_[ExpectedBehavior]_when_[Condition]`
- [ ] 覆盖所有 null 参数场景
- [ ] 覆盖正常请求场景
- [ ] GET 操作有 XML 响应 Mock
- [ ] 使用 `ExpectedException` 规则验证异常
- [ ] AI 生成代码有 `@AIGenerated` 注解
- [ ] `@AIGenerated` 使用动态日期和作者信息

### 向后兼容

- [ ] 不修改现有公共 API 的签名
- [ ] 不删除现有公共方法
- [ ] 不修改现有枚举值
- [ ] 新增方法在接口中声明（不破坏实现类）

---

## 提交信息格式

```
<type>: <简要描述>

<详细说明>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

### type 规范

| type | 说明 |
|------|------|
| feat | 新功能 |
| fix | 修复缺陷 |
| refactor | 重构（不改变行为） |
| test | 添加或修改测试 |
| docs | 文档修改 |
| chore | 构建/工具变更 |

### 示例

```
feat: 添加桶回收站配置管理接口

- 新增 BucketTrashConfiguration 数据模型
- 实现 GET/PUT/DELETE BucketTrash 操作
- 添加 SpecialParamEnum.OBS_TRASH 枚举项
- 编写 MockServer 单元测试
```

---

## 常见问题

### 编译失败

1. 检查 import 是否完整
2. 检查泛型是否正确使用
3. 检查方法签名是否与接口声明一致

### 测试失败

1. 检查 MockServer XML 响应格式是否与 SAX Handler 匹配
2. 检查 ExpectedException 的异常消息是否精确匹配
3. 检查 `mockServer.reset()` 是否在每个测试方法前调用

### 静态分析告警

1. Checkstyle：检查缩进、空格、import 顺序
2. SpotBugs：检查空指针、资源泄漏
3. PMD：检查代码复杂度、命名规范
