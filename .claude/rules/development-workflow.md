# 开发工作流规则

**应用场景**: 当进行代码开发、构建、测试、提交等开发活动时

## 禁止事项

1. **禁止跳过Maven settings参数**
   - 所有Maven命令必须包含 `-s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml` 参数
   - 不得使用简化的Maven命令
   - 不得省略settings.xml配置

2. **禁止未经许可屏蔽静态检查告警**
   - 不得使用 `/*lint -e* */`、`//NOPMD` 等方法绕过门禁
   - 不得未经评审私自调整门禁阈值
   - 静态检查问题必须修复，不能通过注释绕过

3. **禁止违反开发红线**
   - 不得搭便车修改代码（无对应需求单、问题单）
   - 不得开发未公开接口
   - 不得未经CCB批准更改接口
   - 不得违反安全编码规范

4. **禁止跳过验证步骤**
   - 代码生成后必须执行编译验证
   - 测试创建后必须运行测试验证
   - 不得跳过依赖检查、编译、测试、构建、静态检查任一步骤

5. **禁止修改源码后直接运行 failsafe:integration-test**
   - 修改 `src/main/java` 下的源码后，`mvn failsafe:integration-test` 可能使用 `target/classes` 中的旧编译产物
   - 必须先执行 `mvn clean test-compile` 清除旧产物并重新编译
   - 再执行 `mvn failsafe:integration-test` 确保运行的是最新代码
   - 命令组合：`mvn clean test-compile failsafe:integration-test -Dit.test=XXX -Dtest.env=prod`

## 正向示例

```bash
# 正确的Maven命令
mvn clean test -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
mvn test -Dtest=ClassName -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C
mvn clean verify -s ./esdk_obs_java_android_en/CI/settings.xml -f ./esdk_obs_java_android_en/source/pom-java.xml -T 1C

# 正确的集成测试命令（修改源码后先 clean 再跑 IT）
mvn clean test-compile failsafe:integration-test -Dit.test=BucketCompressPolicyIT -Dtest.env=prod -s /usr/local/Maven/conf/settings.xml -T 1C
```

## 反向示例

```bash
# 错误的Maven命令
mvn clean test  # 缺少settings参数
mvn test -Dtest=ClassName  # 缺少settings参数
mvn clean verify  # 缺少settings参数

# 错误：修改源码后直接跑 IT（使用旧 class 文件）
# 编辑 ObsBucketAdvanceService.java 后直接执行：
mvn failsafe:integration-test -Dit.test=XXXIT -Dtest.env=prod  # 可能运行的是修改前的代码
# 应该先 mvn clean test-compile 再 failsafe:integration-test
```

## 验证清单

提交前必须按顺序执行：
1. 依赖验证：`mvn dependency:analyze -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C`
2. 编译验证：`mvn clean compile test-compile -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C`
3. 测试验证：`mvn test -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C`
4. 构建验证：`mvn clean verify -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C`
5. 静态检查：`mvn checkstyle:check spotbugs:check pmd:check -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C`