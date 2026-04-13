# Maven插件配置详解

## 核心插件

### 1. maven-compiler-plugin (编译插件)

**版本**: 3.8.1

**作用**: 编译Java源代码

**配置详情**:
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-compiler-plugin</artifactId>
    <version>3.8.1</version>
    <configuration>
        <source>1.8</source>      <!-- 源代码兼容性：Java 8 -->
        <target>1.8</target>      <!-- 目标字节码版本：Java 8 -->
        <encoding>UTF-8</encoding> <!-- 源文件编码：UTF-8 -->
    </configuration>
</plugin>
```

**功能说明**:
- 将`src/main/java`下的Java源文件编译成字节码
- 支持Java 8语法特性
- 确保字符编码正确处理中文等字符
- 生成class文件到`target/classes`目录

**重要性**: ⭐⭐⭐⭐⭐ (核心插件，必须配置)

---

### 2. maven-surefire-plugin (单元测试插件)

**版本**: 2.22.2

**作用**: 运行单元测试

**配置详情**:
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-surefire-plugin</artifactId>
    <version>2.22.2</version>
    <configuration>
        <includes>
            <include>**/*Test.java</include>        <!-- 包含Test结尾的类 -->
            <include>**/*Tests.java</include>       <!-- 包含Tests结尾的类 -->
        </includes>
        <excludes>
            <exclude>**/*IT.java</exclude>          <!-- 排除集成测试 -->
            <exclude>**/*IntegrationTest.java</exclude> <!-- 排除集成测试 -->
        </excludes>
    </configuration>
</plugin>
```

**功能说明**:
- 运行`src/test/java`下的单元测试
- 通过文件名模式区分单元测试和集成测试
- 支持JUnit测试框架
- 生成测试报告到`target/surefire-reports`

**执行时机**: `test`阶段

**重要性**: ⭐⭐⭐⭐⭐ (质量保证核心插件)

---

### 3. maven-failsafe-plugin (集成测试插件)

**版本**: 2.22.2

**作用**: 运行集成测试

**配置详情**:
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-failsafe-plugin</artifactId>
    <version>2.22.2</version>
    <executions>
        <execution>
            <goals>
                <goal>integration-test</goal>  <!-- 运行集成测试 -->
                <goal>verify</goal>              <!-- 验证测试结果 -->
            </goals>
        </execution>
    </executions>
    <configuration>
        <includes>
            <include>**/*IT.java</include>              <!-- 包含IT结尾的类 -->
            <include>**/*IntegrationTest.java</include> <!-- 包含IntegrationTest结尾的类 -->
        </includes>
    </configuration>
</plugin>
```

**功能说明**:
- 运行`src/it/java`下的集成测试
- 与单元测试分离，确保测试分类清晰
- 支持需要外部环境的测试
- 生成测试报告到`target/failsafe-reports`

**执行时机**: `integration-test`和`verify`阶段

**重要性**: ⭐⭐⭐⭐ (集成测试保证)

---

### 4. maven-shade-plugin (打包插件) - 核心插件

**版本**: 3.2.4

**作用**: 创建包含所有依赖的Uber JAR包

**配置详情**:
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-shade-plugin</artifactId>
    <version>3.2.4</version>
    <executions>
        <execution>
            <phase>package</phase>
            <goals>
                <goal>shade</goal>
            </goals>
            <configuration>
                <createDependencyReducedPom>false</createDependencyReducedPom>
                <artifactSet>
                    <excludes>
                        <!-- 排除Log4j，由应用提供 -->
                        <exclude>org.apache.logging.log4j:log4j-core:jar:</exclude>
                        <exclude>org.apache.logging.log4j:log4j-api:jar:</exclude>
                    </excludes>
                </artifactSet>
                <relocations>
                    <!-- 依赖重定位，避免冲突 -->
                    <relocation>
                        <pattern>com.jamesmurty.utils</pattern>
                        <shadedPattern>shade.jamesmurty.utils</shadedPattern>
                    </relocation>
                    <relocation>
                        <pattern>okhttp3</pattern>
                        <shadedPattern>shade.okhttp3</shadedPattern>
                    </relocation>
                    <relocation>
                        <pattern>okio</pattern>
                        <shadedPattern>shade.okio</shadedPattern>
                    </relocation>
                    <relocation>
                        <pattern>com.fasterxml.jackson.annotation</pattern>
                        <shadedPattern>shade.fasterxml.jackson.annotation</shadedPattern>
                    </relocation>
                    <relocation>
                        <pattern>com.fasterxml.jackson.databind</pattern>
                        <shadedPattern>shade.fasterxml.jackson.databind</shadedPattern>
                    </relocation>
                    <relocation>
                        <pattern>com.fasterxml.jackson.core</pattern>
                        <shadedPattern>shade.fasterxml.jackson.core</shadedPattern>
                    </relocation>
                </relocations>
                <filters>
                    <filter>
                        <artifact>*:*</artifact>
                        <excludes>
                            <!-- 排除签名文件 -->
                            <exclude>META-INF/*.SF</exclude>
                            <exclude>META-INF/*.DSA</exclude>
                            <exclude>META-INF/*.RSA</exclude>
                        </excludes>
                    </filter>
                </filters>
            </configuration>
        </execution>
    </executions>
</plugin>
```

**功能说明**:
- 创建包含所有依赖的单一JAR包（Uber JAR）
- **依赖重定位**: 将第三方依赖的包名重命名，避免与应用依赖冲突
- **排除Log4j**: 日志框架由应用环境提供，避免版本冲突
- **排除签名文件**: 避免打包后的签名验证问题
- 自动处理依赖冲突和资源文件重叠

**执行时机**: `package`阶段

**重要性**: ⭐⭐⭐⭐⭐ (SDK打包的核心插件)

**生成的文件**:
- `esdk-obs-java-3.25.7.jar` (约6.2MB)

---

### 5. maven-javadoc-plugin (文档插件)

**版本**: 3.3.0

**作用**: 生成Java API文档

**配置详情**:
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-javadoc-plugin</artifactId>
    <version>3.3.0</version>
    <configuration>
        <source>1.8</source>
        <encoding>UTF-8</encoding>
        <docencoding>UTF-8</docencoding>
        <charset>UTF-8</charset>
        <excludePackageNames>*.internal.*:*.log:*.proxy:okio:okhttp3</excludePackageNames>
        <destDir>${project.build.directory}/apidocs</destDir>
    </configuration>
</plugin>
```

**功能说明**:
- 从源代码注释生成HTML格式的API文档
- 排除内部实现类和第三方依赖包
- 支持UTF-8编码，正确处理中文注释
- 生成文档到`target/apidocs`目录

**使用方式**:
```bash
mvn javadoc:javadoc
```

**重要性**: ⭐⭐⭐ (文档生成，可选但推荐)

---

### 6. maven-source-plugin (源码插件)

**版本**: 3.2.1

**作用**: 打包源代码

**配置详情**:
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-source-plugin</artifactId>
    <version>3.2.1</version>
</plugin>
```

**功能说明**:
- 将源代码打包成JAR文件
- 方便IDE调试和代码查看
- 生成`-sources.jar`文件

**使用方式**:
```bash
mvn source:jar
```

**重要性**: ⭐⭐⭐ (开发调试支持)

---

### 7. maven-clean-plugin (清理插件)

**版本**: 3.1.0

**作用**: 清理构建输出

**配置详情**:
```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-clean-plugin</artifactId>
    <version>3.1.0</version>
    <configuration>
        <filesets>
            <fileset>
                <directory>logs</directory>  <!-- 清理日志目录 -->
            </fileset>
        </filesets>
    </configuration>
</plugin>
```

**功能说明**:
- 删除`target`目录及其他构建输出
- 清理运行时生成的日志文件
- 确保每次构建都是干净的

**执行时机**: `clean`阶段

**重要性**: ⭐⭐⭐⭐ (构建环境清理)

---

## Profile中的插件配置

### Android Profile - maven-shade-plugin

**版本**: 3.2.4

**特殊配置**:
```xml
<finalName>${project.artifactId}-${project.version}-android</finalName>
<filters>
    <filter>
        <artifact>*:*</artifact>
        <excludes>
            <exclude>META-INF/*.SF</exclude>
            <exclude>META-INF/*.DSA</exclude>
            <exclude>META-INF/*.RSA</exclude>
            <exclude>okio/*</exclude>  <!-- Android特定：排除okio -->
        </excludes>
    </filter>
</filters>
```

**功能说明**:
- 生成Android专用的JAR包
- 文件名包含`-android`后缀
- 排除okio相关内容（Android环境可能有自己的okio实现）

**生成的文件**:
- `esdk-obs-java-3.25.7-android.jar` (约5.9MB)

**重要性**: ⭐⭐⭐⭐⭐ (Android版本打包核心)

---

## Maven生命周期与插件执行顺序

### 完整构建流程
```bash
mvn clean install
```

**执行顺序**:
1. **clean阶段**:
   - `maven-clean-plugin`: 清理构建输出

2. **compile阶段**:
   - `maven-compiler-plugin`: 编译主代码

3. **test阶段**:
   - `maven-compiler-plugin`: 编译测试代码
   - `maven-surefire-plugin`: 运行单元测试

4. **package阶段**:
   - `maven-shade-plugin`: 创建Uber JAR包

5. **install阶段**:
   - 安装JAR包到本地Maven仓库

### 集成测试流程
```bash
mvn verify
```

**额外执行**:
- **integration-test阶段**: `maven-failsafe-plugin`运行集成测试
- **verify阶段**: 验证集成测试结果

---

## 插件依赖关系图

```
maven-clean-plugin
    ↓
maven-compiler-plugin (主代码编译)
    ↓
maven-compiler-plugin (测试代码编译)
    ↓
maven-surefire-plugin (单元测试)
    ↓
maven-shade-plugin (打包)
    ↓
maven-install-plugin (安装)
    ↓
maven-failsafe-plugin (集成测试)
    ↓
maven-verify-plugin (验证)
```

---

## 插件配置最佳实践

### 1. 版本管理
- 所有插件都指定了具体版本，避免使用默认版本
- 版本选择经过测试，确保兼容性

### 2. 测试分离
- 使用`maven-surefire-plugin`运行单元测试
- 使用`maven-failsafe-plugin`运行集成测试
- 通过文件名模式清晰区分测试类型

### 3. 依赖冲突处理
- `maven-shade-plugin`使用重定位避免包名冲突
- 排除Log4j避免与应用环境冲突
- 排除签名文件避免打包问题

### 4. 多环境支持
- 通过Profile支持Java和Android两种打包方式
- Android版本有特定的排除规则
- 默认使用Java版本构建

---

## 插件使用示例

### 常用构建命令
```bash
# 完整构建（跳过测试）
mvn clean package -DskipTests

# 完整构建（运行测试）
mvn clean install

# 只运行单元测试
mvn test

# 只运行集成测试
mvn verify -DskipTests

# 生成文档
mvn javadoc:javadoc

# 生成源码包
mvn source:jar

# Android版本构建
mvn clean package -Pandroid
```

### 网络环境支持
```bash
# 解决SSL证书问题
mvn clean install -DskipTests -Daether.connector.https.securityMode=insecure
```

---

## 插件重要性总结

| 插件 | 重要性 | 必需性 | 说明 |
|------|--------|--------|------|
| maven-compiler-plugin | ⭐⭐⭐⭐⭐ | 必需 | 编译代码 |
| maven-shade-plugin | ⭐⭐⭐⭐⭐ | 必需 | 打包Uber JAR |
| maven-surefire-plugin | ⭐⭐⭐⭐⭐ | 必需 | 单元测试 |
| maven-clean-plugin | ⭐⭐⭐⭐ | 必需 | 清理构建 |
| maven-failsafe-plugin | ⭐⭐⭐⭐ | 推荐 | 集成测试 |
| maven-javadoc-plugin | ⭐⭐⭐ | 推荐 | 生成文档 |
| maven-source-plugin | ⭐⭐⭐ | 推荐 | 源码打包 |

---

## 总结

当前项目的Maven插件配置：
- ✅ **功能完整**: 覆盖了编译、测试、打包、文档生成等全流程
- ✅ **配置合理**: 插件版本选择经过验证，配置参数优化
- ✅ **多环境支持**: 通过Profile支持Java和Android两种构建
- ✅ **依赖管理完善**: 使用Shade插件处理依赖冲突
- ✅ **测试分离**: 单元测试和集成测试清晰分离

插件配置整体上达到了生产级别的标准，能够满足SDK开发和发布的需求。