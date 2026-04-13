# AI代码生成规则

**应用场景**: 当使用AI助手生成代码时

## 禁止事项

1. **禁止缺少@AIGenerated注解**
   - 所有AI生成的测试方法必须添加 `java.com.obs.aitool.AIGenerated` 注解
   - 不得使用硬编码日期，必须动态获取当前系统日期
   - 不得使用硬编码作者信息，必须动态获取git用户信息

2. **禁止跳过验证步骤**
   - 代码生成后必须执行编译验证
   - 不得仅生成代码而不验证功能
   - 识别到的问题必须主动解决，不能仅报告问题

3. **禁止违反测试规范**
   - 不得编写重复的独立测试方法
   - 测试方法命名必须遵循规范格式
   - 参数化测试必须按类别分离

## 正向示例

```java
// 正确的@AIGenerated注解使用
@Test
@AIGenerated(author = "zhanghaoliang", date = "2025-01-10", description = "测试用户登录成功场景")
void should_login_successfully_when_valid_credentials_provided() {
    // 测试逻辑
}

// 正确的参数化测试结构
@RunWith(PowerMockRunner.class)
@PowerMockRunnerDelegate(Parameterized.class)
@PrepareForTest({SystemProperties.class})
public class AsyncTaskServiceTest {
    @Parameterized.Parameter(0) public String testName;
    @Parameterized.Parameter(1) public InputType input;
    @Parameterized.Parameter(2) public ExpectedType expected;
    @Parameterized.Parameter(3) public String testCategory;

    @Parameterized.Parameters(name = "{0}")
    public static Collection<Object[]> testData() {
        return Arrays.asList(new Object[][] {
            {"BothEnabled", BOTH_ENABLED, expected1, "SWITCH_CONFIG"},
            {"MethodNotAccessible", 0L, null, "EXCEPTION"}
        });
    }

    @Test
    public void should_configure_switches_correctly() {
        if (!"SWITCH_CONFIG".equals(testCategory)) {
            return;
        }
        // 测试逻辑
    }
}
```

## 反向示例

```java
// 错误的@AIGenerated注解使用
@Test
@AIGenerated(author = "AI", date = "2024-01-01", description = "测试")
void test() {  // 硬编码信息，命名不规范
    // 测试逻辑
}

// 错误的重复测试方法
@Test
public void testConfig1() { /* 重复逻辑 */ }
@Test
public void testConfig2() { /* 重复逻辑 */ }
// 应该使用参数化测试
```

## 强制验证清单

代码生成完成后必须执行：
- [ ] 代码文件已创建且语法正确
- [ ] 相关依赖已添加到pom.xml
- [ ] @AIGenerated注解正确导入和使用
- [ ] 编译无错误：`mvn compile test-compile -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C`
- [ ] 测试执行通过：`mvn test -Dtest=GeneratedTestClassName -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C`
- [ ] 测试覆盖率满足要求（100%）
- [ ] 功能验证通过