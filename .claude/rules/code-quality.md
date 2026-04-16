# 代码质量规则

**应用场景**: 当编写、修改、审查代码时

## 禁止事项

1. **禁止违反SOLID原则**
   - 不得创建承担多个职责的类
   - 不得直接修改已发布的类，应该通过扩展
   - 不得破坏里氏替换原则
   - 不得创建臃肿的接口
   - 不得依赖具体实现，应该依赖抽象

2. **禁止低质量代码**
   - 不得使用原始类型，必须使用Java泛型
   - 不得直接访问私有字段，优先使用公开API
   - 不得使用未经检查的异常处理
   - 不得硬编码敏感信息（密码、密钥等）

3. **禁止新增构造函数引起 null 参数歧义**
   - 新增构造函数时，必须检查与已有构造函数的参数类型是否冲突
   - 当已有 `(String, String, Throwable)` 时，不得新增 `(String, String, String)` 构造函数而不修改现有的委托调用
   - 委托构造函数传递 `null` 时必须显式类型转换：`this(message, xmlMessage, (Throwable) null)`
   - 编译错误 `对XXX的引用不明确` 即为违反本规则

4. **禁止从 ObsException.getResponseHeaders() 读取 HTTP 原始头信息**
   - `ObsException.getResponseHeaders()` 返回的 Map 是经过 SDK 内部处理的头信息，可能不包含所有原始 HTTP 头（如 Content-Type）
   - 需要获取 Content-Type 等原始 HTTP 响应头时，必须在 `RestStorageService` 层从 OkHttp `Response` 对象中直接读取
   - 不得假设 `getResponseHeaders().get("Content-Type")` 一定返回有效值

5. **禁止不完整的测试**
   - 新增代码分支覆盖率必须达到100%
   - 不得编写重复的独立测试方法，必须使用参数化测试
   - 测试方法命名必须遵循 `should_[ExpectedBehavior]_when_[Condition]` 格式

## 正向示例

```java
// 正确的接口设计
public interface IObsClient {
    HeaderResponse getBucketMetadata(GetBucketMetadataRequest request);
}

// 正确的测试方法命名
@Test
public void should_return_success_when_valid_credentials_provided() {
    // 测试逻辑
}

// 正确的参数化测试
@Parameterized.Parameters(name = "{0}")
public static Collection<Object[]> testData() {
    return Arrays.asList(new Object[][] {
        {"TestCase1", input1, expected1},
        {"TestCase2", input2, expected2}
    });
}

// 正确的构造函数 null 参数消歧义
public ServiceException(String message, String xmlMessage) {
    this(message, xmlMessage, (Throwable) null);  // 显式类型转换消除歧义
}

// 正确的 Content-Type 获取方式（在 RestStorageService 层）
String contentType = response.header(CommonHeaders.CONTENT_TYPE);
return new ServiceException(message, xmlMessage, contentType);
```

## 反向示例

```java
// 错误的原始类型使用
List list = new ArrayList();  // 应该使用 List<String>

// 错误的测试方法命名
@Test
public void testLogin() {  // 应该使用 should_login_successfully_when_valid_credentials_provided

// 错误的重复测试方法
@Test
public void testConfig1() { /* 重复逻辑 */ }
@Test
public void testConfig2() { /* 重复逻辑 */ }
// 应该使用参数化测试合并

// 错误的构造函数 null 委托（编译错误：对ServiceException的引用不明确）
public ServiceException(String message, String xmlMessage) {
    this(message, xmlMessage, null);  // null 既匹配 String 也匹配 Throwable
}
// 已有: ServiceException(String, String, Throwable)
// 新增: ServiceException(String, String, String) → null 歧义

// 错误的 Content-Type 获取方式
catch (ObsException e) {
    String ct = e.getResponseHeaders().get("Content-Type");  // 可能返回 null
    if ("application/json".equals(ct)) { ... }  // 永远不会进入
}
```

## 质量检查清单

- [ ] 代码符合SOLID原则
- [ ] 使用Java泛型，无原始类型
- [ ] 公共方法有JavaDoc注释
- [ ] 异常处理合理
- [ ] 测试覆盖率达到100%
- [ ] 测试使用参数化方式
- [ ] 通过静态检查工具