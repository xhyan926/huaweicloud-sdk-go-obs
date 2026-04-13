# 依赖管理文档

## 依赖健康度评估

### 总体评估：🟢 良好

**优点**：
1. ✅ 没有严重的版本冲突
2. ✅ Maven正确处理了依赖传递
3. ✅ Shade插件成功打包
4. ✅ 测试和主依赖分离良好

**需要改进的地方**：
1. ⚠️ Log4j依赖范围应该改为test
2. ⚠️ 可以考虑添加依赖管理规则

---

## 版本冲突分析

### 1. OkHttp/Okio版本冲突（已解决）

**问题描述**：
- OkHttp 4.12.0 依赖 OkIo 3.6.0
- 项目直接声明 OkIo 3.8.0
- Maven自动选择了更高版本 3.8.0

**冲突状态**：✅ 已解决
- Maven自动处理了版本冲突
- 最终使用 OkIo 3.8.0

**影响**：无影响，使用更高版本是安全的

---

### 2. Jackson版本冲突（已解决）

**问题描述**：
- MockWebServer 4.12.0 依赖 Jackson 2.14.x
- 项目使用 Jackson 2.13.4
- Maven选择了项目声明的 2.13.4 版本

**冲突状态**：✅ 已解决
- 测试依赖的Jackson版本被主依赖版本覆盖
- 避免了版本不一致问题

**影响**：无影响，统一使用项目声明的版本

---

## 传递依赖冲突（正常）

### 1. Kotlin标准库重复

**问题描述**：
- OkHttp 4.12.0 引入了多个Kotlin标准库
- kotlin-stdlib, kotlin-stdlib-jdk7, kotlin-stdlib-jdk8
- 这些库之间存在依赖关系

**冲突状态**：✅ 正常现象
- 这是Kotlin生态系统的正常依赖结构
- Maven会正确处理这些依赖

**影响**：无影响

---

### 2. Jackson模块重复

**问题描述**：
- jackson-core, jackson-databind, jackson-annotations
- jackson-databind依赖其他两个模块
- 存在META-INF资源文件重叠

**冲突状态**：✅ 正常现象
- Shade插件会正确处理资源文件重叠
- 只保留一个版本的资源文件

**影响**：无影响，Shade插件已处理

---

## 未声明但使用的依赖

### 检测到的问题依赖
```
Used undeclared dependencies found:
- com.squareup.okio:okio-jvm:jar:3.8.0:compile
- org.powermock:powermock-core:jar:2.0.9:test
- org.mock-server:mockserver-core:jar:5.15.0:test
- org.mock-server:mockserver-client-java:jar:5.15.0:test
- org.powermock:powermock-api-support:jar:2.0.9:test
```

**分析结果**：
1. **okio-jvm**: 这是OkIo 3.8.0的传递依赖，不需要显式声明
2. **PowerMock相关**: PowerMock的内部依赖，通过powermock-module-junit4和powermock-api-mockito2引入
3. **MockServer相关**: MockServer的内部依赖，通过mockserver-netty引入

**建议**：这些是正常的传递依赖，不需要显式声明

---

## 已声明但未使用的依赖

### 检测到的问题依赖
```
Unused declared dependencies found:
- com.jamesmurty.utils:java-xmlbuilder:jar:1.3:compile
- com.squareup.okio:okio:jar:3.8.0:compile
- org.powermock:powermock-module-junit4:jar:2.0.9:test
```

**分析结果**：
1. **java-xmlbuilder**:
   - Maven分析显示未使用，但实际代码中可能通过反射使用
   - 建议保留，确保功能完整性

2. **okio**:
   - 虽然显示未使用，但OkHttp内部依赖它
   - 建议保留显式声明，控制版本

3. **powermock-module-junit4**:
   - 测试代码中确实使用了PowerMock
   - Maven分析可能存在误判
   - 建议保留

**建议**：保留这些依赖，确保功能完整性

---

## 测试依赖问题

### Log4j依赖范围问题
```
Non-test scoped test only dependencies found:
- org.apache.logging.log4j:log4j-api:jar:2.17.1:compile
- org.apache.logging.log4j:log4j-core:jar:2.17.1:compile
```

**问题描述**：
- Log4j依赖声明为compile范围
- 但实际只在测试中使用
- 应该改为test范围

**建议修复**：
```xml
<dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-core</artifactId>
    <version>${log4j.version}</version>
    <scope>test</scope>  <!-- 添加scope -->
</dependency>
<dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-api</artifactId>
    <version>${log4j.version}</version>
    <scope>test</scope>  <!-- 添加scope -->
</dependency>
```

---

## 打包警告分析

### Shade插件警告
```
[WARNING] jackson-core-2.13.4.jar, jackson-databind-2.13.4.jar define 1 overlapping resource:
[WARNING]   - META-INF/NOTICE
[WARNING] annotations-13.0.jar, esdk-obs-java-3.25.7.jar, ... define 1 overlapping resource:
[WARNING]   - META-INF/MANIFEST.MF
```

**分析结果**：
- META-INF资源文件重叠是正常现象
- Shade插件会自动处理，只保留一个版本
- 不影响功能

**建议**：可以忽略这些警告，或者配置Shade插件排除特定资源

---

## 优化建议

### 1. 修复Log4j依赖范围
```xml
<dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-core</artifactId>
    <version>${log4j.version}</version>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-api</artifactId>
    <version>${log4j.version}</version>
    <scope>test</scope>
</dependency>
```

---

### 2. 添加Enforcer插件规则
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-enforcer-plugin</artifactId>
    <version>3.6.2</version>
    <executions>
        <execution>
            <id>enforce-versions</id>
            <goals>
                <goal>enforce</goal>
            </goals>
            <configuration>
                <rules>
                    <requireMavenVersion>
                        <version>[3.6.0,)</version>
                    </requireMavenVersion>
                    <requireJavaVersion>
                        <version>[1.8,)</version>
                    </requireJavaVersion>
                    <dependencyConvergence/>
                </rules>
            </configuration>
        </execution>
    </executions>
</plugin>
```

---

### 3. 优化Shade插件配置
```xml
<filters>
    <filter>
        <artifact>*:*</artifact>
        <excludes>
            <exclude>META-INF/*.SF</exclude>
            <exclude>META-INF/*.DSA</exclude>
            <exclude>META-INF/*.RSA</exclude>
            <exclude>META-INF/LICENSE*</exclude>
            <exclude>META-INF/NOTICE*</exclude>
            <exclude>META-INF/MANIFEST.MF</exclude>
        </excludes>
    </filter>
</filters>
```

---

## 结论

当前项目的依赖管理整体上是健康的：
- 没有严重的依赖冲突
- 版本选择合理
- 打包过程正常

主要需要改进的是Log4j的依赖范围，以及可以考虑添加一些依赖管理规则来预防未来的依赖问题。

**风险等级**：🟢 低风险
**建议操作**：可选优化，不影响当前功能

---

## 依赖版本管理

所有依赖版本在`pom.xml`的`<properties>`部分统一管理：

```xml
<properties>
    <okhttp.version>4.12.0</okhttp.version>
    <okio.version>3.8.0</okio.version>
    <jackson.version>2.15.4</jackson.version>
    <log4j.version>2.24.3</log4j.version>
    <!-- 其他依赖版本 -->
</properties>
```

---

## 依赖重定位

为避免与应用依赖冲突，以下依赖被重定位：

- `com.jamesmurty.utils` → `shade.jamesmurty.utils`
- `okhttp3` → `shade.okhttp3`
- `okio` → `shade.okio`
- `com.fasterxml.jackson.annotation` → `shade.fasterxml.jackson.annotation`
- `com.fasterxml.jackson.databind` → `shade.fasterxml.jackson.databind`
- `com.fasterxml.jackson.core` → `shade.fasterxml.jackson.core`

---

## 常用依赖管理命令

```bash
# 查看依赖树
mvn dependency:tree

# 分析依赖使用情况
mvn dependency:analyze

# 查看依赖冲突
mvn dependency:tree -Dverbose

# 强制更新依赖
mvn clean install -U

# 查看有效POM
mvn help:effective-pom
```