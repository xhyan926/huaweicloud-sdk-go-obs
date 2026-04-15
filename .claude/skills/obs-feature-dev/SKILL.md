---
name: obs-feature-dev
description: |
  当用户提到为 OBS SDK 添加新功能、新 API 方法、扩展桶/对象操作、实现新的配置接口时激活。
  关键词：添加新功能、新API、桶操作、对象操作、扩展客户端、新配置、实现接口方法。
  示例："为OBS SDK添加桶标签管理功能"、"实现一个新的对象操作API"、"扩展桶的回收站配置"。
---

# OBS SDK 新功能开发指南

本 Skill 指导完成 OBS Java SDK 新功能开发的全流程，覆盖从接口文档分析到样例代码编写的九个阶段。

## 关键源文件索引

| 文件 | 路径 | 用途 |
|------|------|------|
| 特殊参数枚举 | `src/main/java/com/obs/services/model/SpecialParamEnum.java` | 子资源参数注册 |
| 主接口 | `src/main/java/com/obs/services/IObsClient.java` | 公共 API 声明 |
| 抽象客户端 | `src/main/java/com/obs/services/AbstractBucketAdvanceClient.java` | 参数校验 + doAction 委托 |
| 内部服务 | `src/main/java/com/obs/services/internal/service/ObsBucketAdvanceService.java` | REST 调用实现 |
| XML 解析器 | `src/main/java/com/obs/services/internal/handler/XmlResponsesSaxParser.java` | GET 响应解析 |
| XML 构建器目录 | `src/main/java/com/obs/services/internal/xml/` | PUT 请求 XML 构建 |
| 模型类目录 | `src/main/java/com/obs/services/model/` | Request/Result/Configuration 类 |

---

## Phase 1: 接口文档读取与计划生成

**目标**: 读取接口文档（在线 URL 或本地文件），提取完整的 API 规范，生成结构化开发计划。

### 1.1 文档输入

支持两种文档来源：

| 类型 | 输入方式 | 工具 |
|------|----------|------|
| 在线文档 | 提供 URL | `mcp__web-reader__webReader` 读取网页内容 |
| 本地文档 | 提供文件路径 | `Read` 工具读取文件内容 |

支持的文档格式：HTML 网页、Markdown、纯文本。

### 1.2 文档解析与信息提取

从接口文档中提取以下结构化信息，**必须完整提取，不得遗漏**：

**API 基本信息：**
- 接口名称、功能描述
- HTTP 方法（GET / PUT / POST / DELETE）
- URI 路径和查询参数（子资源标识）
- 是否为桶级别操作或对象级别操作

**请求规范：**
- 请求 URI 参数（路径参数、查询参数）
- 请求头（Content-Type、自定义头等）
- 请求体 XML/JSON 结构（完整列出每个元素、类型、是否必填、取值范围）

**响应规范：**
- 响应状态码及含义（200/204/400/404 等）
- 响应头（特殊响应头字段）
- 响应体 XML/JSON 结构（完整列出每个元素、类型）
- 错误响应体结构

**特殊处理：**
- 权限要求
- 约束条件（如桶类型限制、取值范围）
- 与其他接口的依赖关系

### 1.3 生成开发计划

基于提取的接口信息，生成结构化计划，包含以下板块：

- **接口信息**: 名称、描述、HTTP 方法、URI、操作级别
- **子资源参数**: 参数名、SpecialParamEnum 枚举名
- **请求参数**: 路径参数、查询参数、请求头、请求体结构（含字段类型/必填/取值范围）
- **响应规范**: 状态码、响应头、响应体结构、错误码
- **约束条件**: 桶类型限制、参数取值范围等
- **文件清单**: 需创建和需修改的文件列表

### 1.4 质量关卡

**计划生成后必须向用户展示，等待用户确认后再进入下一阶段。**

向用户展示的内容：
1. 提取的接口规范摘要
2. 生成的文件清单
3. 识别到的特殊约束
4. 询问："以上开发计划是否正确？是否需要调整？"

---

## Phase 2: 开发计划校验

**目标**: 使用独立 Agent 重新读取接口文档，逐项比对开发计划的准确性，确保零遗漏。

### 2.1 启动独立校验 Agent

**必须使用 `Agent` 工具启动一个独立的 `Explore` 类型 Agent**，与主流程隔离执行校验。

校验 Agent 的职责：
1. 重新从原始文档源读取接口文档（不依赖 Phase 1 的提取结果）
2. 逐项比对 Phase 1 生成的开发计划
3. 输出差异报告

### 2.2 校验检查项

校验 Agent 必须逐一核对以下维度：

**URI 与参数校验：**
- [ ] HTTP 方法是否正确（GET/PUT/DELETE）
- [ ] URI 路径是否正确
- [ ] 子资源参数名称是否与文档完全一致（区分大小写）
- [ ] 查询参数是否完整、无遗漏

**请求体校验（PUT 操作）：**
- [ ] XML/JSON 根元素名称是否与文档一致
- [ ] 每个字段名称是否与文档完全一致（区分大小写）
- [ ] 每个字段类型是否正确（String/int/long/boolean/enum 等）
- [ ] 必填/选填标记是否正确
- [ ] 字段取值范围是否正确记录
- [ ] 嵌套结构是否完整（不遗漏子元素）
- [ ] XML 命名空间是否正确（如有）

**响应体校验（GET 操作）：**
- [ ] 响应根元素名称是否与文档一致
- [ ] 每个字段名称、类型是否与文档完全一致
- [ ] 响应状态码是否完整记录
- [ ] 响应头字段是否遗漏

**请求头/响应头校验：**
- [ ] Content-Type 是否正确
- [ ] 自定义请求头是否完整记录
- [ ] 特殊响应头是否遗漏

**约束条件校验：**
- [ ] 参数约束（取值范围、长度限制）是否完整
- [ ] 桶类型限制是否记录
- [ ] 权限要求是否记录

### 2.3 差异处理

校验 Agent 输出差异报告后：

1. **无差异** → 进入 Phase 3（需求分析）继续开发
2. **有差异** → 按以下方式处理：
   - 将差异报告展示给用户
   - 根据文档修正开发计划
   - 修正后可再次校验（如差异较多）
   - 确认无误后进入 Phase 3（需求分析）

### 2.4 校验 Agent Prompt 模板

```
你是一个 OBS SDK 接口文档校验专家。请完成以下任务：

1. 读取接口文档（来源：{文档URL或路径}）
2. 从文档中提取完整的接口规范（HTTP方法、URI、请求参数、请求体、响应体、请求头、响应头、错误码）
3. 将提取结果与以下开发计划逐项比对：
   {Phase 1 生成的开发计划}
4. 输出差异报告，格式如下：
   - ✅ 一致项：[列出]
   - ❌ 不一致项：[列出具体差异，包含文档原文和计划中的值]
   - ⚠️ 遗漏项：[列出文档中存在但计划中未包含的内容]

注意：
- 字段名称区分大小写
- 必须比对每一个参数、每一个字段
- 不要忽略任何文档中提到的约束条件
```

---

## Phase 3: 需求分析

**目标**: 明确功能边界，确定需要创建/修改的文件清单。

### 3.1 确定操作类型

| 操作 | HTTP 方法 | 典型场景 |
|------|-----------|----------|
| 查询配置 | GET | 获取桶的某项配置 |
| 设置配置 | PUT | 设置/更新桶的某项配置 |
| 删除配置 | DELETE | 删除桶的某项配置 |
| 列举 | GET (列表) | 列举桶的多条规则 |

### 3.2 确定参数和返回值

- 请求参数：桶名 + 配置对象（PUT 需要）
- 响应字段：配置对象（GET 需要）或仅 HeaderResponse（PUT/DELETE）
- 特殊参数（子资源标识）：需要在 `SpecialParamEnum` 中注册

### 3.3 质量关卡

向用户展示功能摘要并确认：

```
功能名称: [如 BucketTrash]
操作集合: GET / PUT / DELETE
子资源参数: [如 x-obs-trash]
请求参数: [列出]
响应字段: [列出]
需新增的文件:
  - model 层: [Configuration/Request/Result 类]
  - xml 层: [XMLBuilder 类]
需修改的文件:
  - SpecialParamEnum.java (新增枚举项)
  - IObsClient.java (声明方法)
  - AbstractBucketAdvanceClient.java (实现委托)
  - ObsBucketAdvanceService.java (REST 实现)
  - XmlResponsesSaxParser.java (GET 响应处理器)
```

---

## Phase 4: 请求/响应模型设计

**目标**: 创建 Configuration、Request、Result 模型类。

所有模型类放置于 `src/main/java/com/obs/services/model/<feature>/` 包下。

详细代码模板参考 → `references/request-response-patterns.md`

### 4.1 类设计规则

- **Configuration 类**: 纯 POJO，包含业务字段 + 构造函数 + getter/setter
- **GET Request**: 继承 `BaseBucketRequest`，`httpMethod = HttpMethodEnum.GET`
- **PUT Request**: 继承 `BaseBucketRequest`，`httpMethod = HttpMethodEnum.PUT`，持有 Configuration 对象
- **DELETE Request**: 继承 `BaseBucketRequest`，`httpMethod = HttpMethodEnum.DELETE`
- **GET Result**: 继承 `HeaderResponse`，持有 Configuration 对象

### 4.2 命名约定

```
Configuration:  <Feature>Configuration     如 BucketTrashConfiguration
GET Request:    Get<Feature>Request        如 GetBucketTrashRequest
PUT Request:    Set<Feature>Request        如 SetBucketTrashRequest
DELETE Request: Delete<Feature>Request     如 DeleteBucketTrashRequest
GET Result:     Get<Feature>Result         如 GetBucketTrashResult
```

---

## Phase 5: 接口声明

**目标**: 在 `IObsClient.java` 中声明公共 API 方法。

### 5.1 方法签名模式

```java
// PUT 操作 - 设置配置
HeaderResponse set<Feature>(Set<Feature>Request request) throws ObsException;

// GET 操作 - 获取配置
Get<Feature>Result get<Feature>(Get<Feature>Request request) throws ObsException;

// DELETE 操作 - 删除配置
HeaderResponse delete<Feature>(Delete<Feature>Request request) throws ObsException;
```

### 5.2 JavaDoc 模板

```java
/**
 * [操作描述].
 *
 * @param request
 *            Request parameters
 * @return Common response headers / Result object
 * @throws ObsException
 *             OBS SDK self-defined exception, thrown when the interface
 *             fails to be called or access to OBS fails
 */
```

---

## Phase 6: 实现

**目标**: 在客户端层和服务层实现业务逻辑。

### 6.1 SpecialParamEnum 注册

在 `SpecialParamEnum.java` 中添加子资源参数：

```java
<FEATURE_ENUM>("x-obs-<feature>"),
```

### 6.2 AbstractBucketAdvanceClient 实现

对每个操作方法，使用以下模式：

```java
@Override
public Get<Feature>Result get<Feature>(Get<Feature>Request request) throws ObsException {
    ServiceUtils.assertParameterNotNull(request, "Get<Feature>Request is null");
    ServiceUtils.assertParameterNotNull(request.getBucketName(), "bucketName is null");
    return this.doActionWithResult("get<Feature>", request.getBucketName(),
        new ActionCallbackWithResult<Get<Feature>Result>() {
            @Override
            public Get<Feature>Result action() throws ServiceException {
                return AbstractBucketAdvanceClient.this.get<Feature>Impl(request);
            }
        });
}
```

**PUT/DELETE 操作**返回类型使用 `HeaderResponse` 而非自定义 Result。

### 6.3 ObsBucketAdvanceService 实现

三个核心 REST 方法模式：

**GET 操作** - 使用 `performRestGet()`：
1. 构建 `requestParameters` Map，放入 `SpecialParamEnum` 子资源参数
2. 调用 `performRestGet()` 获取响应
3. 调用 `verifyResponseContentType()` 验证
4. 使用 `getXmlResponseSaxParser().parse()` 解析 XML
5. 构建 Result 对象，调用 `setHeadersAndStatus()`

**PUT 操作** - 使用 `transRequest()` + `performRequest()`：
1. 构建 `requestParams` Map，放入 `SpecialParamEnum` 子资源参数
2. 创建 XML 构建器，生成 XML 字符串
3. 设置 `Content-Length`、`Content-MD5`、`Content-Type` 请求头
4. 调用 `transRequest()` 创建 `NewTransResult`
5. 设置 headers、params、body
6. 调用 `performRequest()` 发送请求
7. 返回 `build(response)`

**DELETE 操作** - 使用 `performRestDelete()`：
1. 构建 `requestParams` Map，放入 `SpecialParamEnum` 子资源参数
2. 调用 `performRestDelete()` 获取响应
3. 返回 `this.build(response)`

### 6.4 XML 构建（PUT 操作需要）

在 `src/main/java/com/obs/services/internal/xml/` 下创建 XML 构建器类：
- 继承 `ObsSimpleXMLBuilder`
- 使用 `startElement()`、`append()`、`endElement()` 构建 XML
- 空值校验：使用 `checkBucketPublicAccessBlock()` 方法（继承自父类）

### 6.5 XML 解析（GET 操作需要）

在 `XmlResponsesSaxParser.java` 中添加内部处理器类：
- 继承 `DefaultXmlHandler`
- 重写 `endElement(String name, String elementText)` 方法
- 使用 XML Builder 中定义的常量进行元素名匹配

---

## Phase 7: 测试

**目标**: 编写单元测试和集成测试，确保 100% 分支覆盖率和端到端功能正确性。

详细测试模板参考 → `references/test-patterns.md`

### 7.1 单元测试

**目标**: 使用 MockServer 隔离外部依赖，验证参数校验逻辑和请求构造的正确性。

- 使用 **JUnit 4** + **MockServer**
- 测试类命名: `<Feature>UnitTest.java`
- 放置于 `src/test/java/com/obs/services/model/<feature>/`
- 通过 `TestTools.getPipelineEnvironment()` 获取 ObsClient（流量走 MockServer 代理）
- 每个测试方法前调用 `mockServer.reset()` 清除旧配置

**单元测试必须覆盖的场景**（每个操作）：

| 场景 | 预期结果 |
|------|----------|
| request 为 null | `IllegalArgumentException` |
| bucketName 为 null | `IllegalArgumentException` |
| 正常请求（MockServer 模拟响应） | 验证返回值字段 |
| PUT 特有: Configuration 为 null | 异常 |
| PUT 特有: Configuration 字段不合法 | 异常 |
| GET 特有: XML 响应解析 | 验证 Configuration 对象字段值 |

**单元测试验证命令**：
```bash
mvn test -Dtest=<Feature>UnitTest -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
```

### 7.2 集成测试

**目标**: 连接真实 OBS 环境，验证完整的端到端功能链路（创建桶 → 设置配置 → 查询配置 → 删除配置 → 清理桶）。

- 测试类命名: `<Feature>IT.java`（IT = Integration Test）
- 放置于 `src/it/java/com/obs/integrated_test/<feature>/`
- 通过 `TestTools.getPipelineEnvironment()` 获取真实 ObsClient
- 使用 `@Rule TestName` 生成唯一桶名
- 使用 `try/finally` 确保测试桶被清理

**集成测试必须覆盖的场景**：

| 场景编号 | 测试内容 | 说明 |
|----------|----------|------|
| IT-001 | 设置配置 → 查询验证 → 边界值测试 | 验证正常设置和取值范围 |
| IT-002 | 设置 → 查询 → 删除 → 再查询(404) | 验证完整 CRUD 生命周期 |
| IT-003 | 不支持的桶类型 | 验证异常约束（如对象桶不支持回收站） |
| IT-004 | S3 兼容接口 | 验证 V2 接口方式调用（如适用） |

**集成测试验证命令**：
```bash
mvn verify -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
```

### 7.3 测试方法命名规范

- **单元测试**: `should_[ExpectedBehavior]_when_[Condition]`
- **集成测试**: `test_SDK_<module>_<feature>_<序号>`（如 `test_SDK_fs_trash_001`）

### 7.4 @AIGenerated 注解

AI 生成的测试方法必须添加注解：

```java
@com.obs.aitool.AIGenerated(
    author = System.getProperty("user.name"),
    date = java.time.LocalDate.now().toString(),
    description = "..."
)
```

---

## Phase 8: 构建验证与代码审查

**目标**: 确保代码通过编译、测试和静态分析。

详细验证命令参考 → `references/verification-checklist.md`

### 8.1 验证步骤（按顺序执行）

1. **编译验证**
2. **单元测试验证**
3. **集成测试验证**
4. **静态分析**

### 8.2 代码审查要点

- [ ] 所有公共方法有 JavaDoc
- [ ] 使用 Java 泛型，无原始类型
- [ ] 异常处理合理
- [ ] 参数校验完整（null 检查）
- [ ] SOLID 原则合规
- [ ] 向后兼容性（不影响现有 API）

---

## Phase 9: 编写样例代码

**目标**: 在 `samples/` 目录下编写用户可直接参考的样例代码，演示新功能的基本用法。

详细样例模板参考 → `references/sample-patterns.md`

### 9.1 样例文件规范

- 文件位置: `samples/src/main/java/com/example/obs/<Feature>Sample.java`
- 包名: `com.example.obs`
- 许可证头: Apache License 2.0

### 9.2 样例结构要点

- endPoint/ak/sk 使用占位字符串，不包含真实凭证
- 每个操作（set/get/delete）拆分为独立 `private static` 方法
- 添加步骤注释说明操作目的
- `main()` 中统一捕获 `ObsException`，`finally` 中关闭 ObsClient
- 使用 `System.out.println()` 打印操作结果

---

## 参考文档

根据需要加载项目文档：
- 架构设计: `.claude/docs/architecture.md`
- 构建指南: `.claude/docs/build-guide.md`
- 开发工作流: `.claude/rules/development-workflow.md`
- 代码质量: `.claude/rules/code-quality.md`
- AI 代码生成: `.claude/rules/ai-code-generation.md`
