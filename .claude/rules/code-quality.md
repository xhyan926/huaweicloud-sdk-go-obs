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
   - 新增代码可达分支覆盖率必须达到100%（不含经分析确认的不可达防御性分支）
   - 不得编写重复的独立测试方法，必须使用参数化测试
   - 测试方法命名必须遵循 `should_[ExpectedBehavior]_when_[Condition]` 格式

6. **禁止使用多个独立布尔标志表达互斥状态**
   - 当对象存在互斥的生命周期状态（如"运行中"、"已暂停"、"已取消"）时，不得使用多个 `AtomicBoolean` 分别管理
   - 必须使用 `AtomicReference<Enum>` 或等效的单变量状态机，保证状态转换的原子性和一致性
   - 独立布尔标志允许非法组合（如 `isPaused=true && isCancelled=true`），且无法原子地表达"cancel 覆盖 pause"等语义

7. **禁止在状态转换方法中缺少前置条件校验**
   - 状态转换方法必须校验当前状态是否允许本次转换
   - 非法转换必须抛出 `IllegalStateException`，并在异常消息中包含当前状态和期望状态
   - 终态（如 CANCELLED）不允许转换到任何其他状态
   - 幂等操作（如对已取消对象重复调用 cancel）不在此列，但必须明确为无副作用的 no-op

8. **禁止在公共类命名中使用 MVC 等架构框架的固定语义后缀**
   - 不得使用 `Controller`、`Service`、`Repository`、`Dao` 等在 Spring/MVC 等框架中有固定含义的后缀来命名非框架角色类
   - 类名应准确反映其在 SDK 中的实际职责，避免使用者在框架集成时产生混淆
   - 可使用的替代后缀：`Handle`（持有以控制的句柄）、`Manager`（管理器）、`Context`（上下文）、`Config`（配置）

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

// 正确的单变量状态机（规则 6/7，示例代码非项目实际类）
public class ResumableTransferHandle {
    private enum State { ACTIVE, PAUSED, CANCELLED }
    private final AtomicReference<State> state = new AtomicReference<>(State.ACTIVE);

    public void pause() {
        if (!state.compareAndSet(State.ACTIVE, State.PAUSED)) {
            throw new IllegalStateException(
                    "Cannot pause: current state is " + state.get() + ", expected ACTIVE");
        }
    }

    public void cancel() {
        if (state.get() == State.CANCELLED) {
            return; // 幂等：已取消时无副作用
        }
        state.set(State.CANCELLED);
    }

    public void resume() {
        if (!state.compareAndSet(State.PAUSED, State.ACTIVE)) {
            throw new IllegalStateException(
                    "Cannot resume: current state is " + state.get() + ", expected PAUSED");
        }
    }
}

// 正确的类名：准确反映 SDK 内职责（规则 8）
public class ResumableTransferHandle { ... }   // 句柄：持有以控制传输
public class ProgressManager { ... }            // 管理器：管理进度回调
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

// 错误：多个独立布尔标志，无法防止非法组合（规则 6）
public class BadExample {
    private final AtomicBoolean isPaused = new AtomicBoolean(false);
    private final AtomicBoolean isCancelled = new AtomicBoolean(false);

    public void pause() { isPaused.set(true); }              // cancel 后仍可 pause！
    public void resetForResume() {                            // cancel 后仍可 resume！
        isCancelled.set(false);
        isPaused.set(false);
    }
}
// 调用方：handle.cancel(); handle.pause(); // 不报错，语义矛盾

// 错误：状态转换无校验，静默接受非法操作（规则 7）
public void resume() { isPaused.set(false); }  // 未暂停时也能调用，没有报错

// 错误：非 MVC 类使用 Controller 后缀（规则 8）
public class ResumableTransferController { ... }  // 误导使用者以为是 Spring Controller
```

## 质量检查清单

- [ ] 代码符合SOLID原则
- [ ] 使用Java泛型，无原始类型
- [ ] 公共方法有JavaDoc注释
- [ ] 异常处理合理
- [ ] 测试覆盖率（可达分支）达到100%
- [ ] 测试使用参数化方式
- [ ] 互斥状态使用单变量状态机管理
- [ ] 状态转换方法有前置条件校验
- [ ] 类名不使用架构框架的固定语义后缀
- [ ] 通过静态检查工具
