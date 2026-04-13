# 构建指南

## 环境要求

- JDK 1.8+
- Maven 3.6+

## 基本构建命令

### 清理和编译
```bash
mvn clean compile
```

### 打包
```bash
# 打包（不运行测试）
mvn clean package -DskipTests

# 打包并运行测试
mvn clean package
```

### 网络环境配置
如果遇到SSL证书问题，可以使用以下参数：
```bash
mvn clean package -DskipTests -Daether.connector.https.securityMode=insecure
```

## 测试相关命令

### 运行测试
```bash
# 只运行单元测试
mvn test

# 只运行集成测试
mvn verify -DskipTests

# 运行所有测试（单元测试 + 集成测试）
mvn verify

# 跳过所有测试
mvn package -DskipTests

# 运行指定测试类
mvn test -Dtest=ClassName
```

## 安装和部署

```bash
# 安装到本地Maven仓库
mvn clean install

# 部署到远程仓库
mvn clean deploy
```

## 文档生成

```bash
# 生成Javadoc
mvn javadoc:javadoc

# 生成源码包
mvn source:jar

# 生成所有文档
mvn javadoc:javadoc source:jar
```

## Profile配置

### Java版本构建（默认）
```bash
mvn clean package

# 如果遇到网络问题
mvn clean package -Daether.connector.https.securityMode=insecure
```

### Android版本构建
```bash
mvn clean package -Pandroid

# 如果遇到网络问题
mvn clean package -Pandroid -Daether.connector.https.securityMode=insecure
```

## 构建产物

构建完成后，在`target`目录下会生成以下文件：

- `esdk-obs-java-3.25.7.jar` - 主JAR包（包含依赖的Shade包，约6.2MB）
- `esdk-obs-java-3.25.7-sources.jar` - 源码包
- `esdk-obs-java-3.25.7-javadoc.jar` - Javadoc包
- `esdk-obs-java-3.25.7-android.jar` - Android版本JAR包（使用-Pandroid时，约5.9MB）

## Samples模块构建

Samples模块是独立的Maven模块，用于存放示例代码，包名为`com.example.obs`。

### 构建Samples模块
```bash
# 进入samples目录
cd samples

# 清理并打包
mvn clean package -DskipTests

# 如果遇到网络问题
mvn clean package -DskipTests -Daether.connector.https.securityMode=insecure
```

### Samples模块产物
构建完成后，在`samples/target`目录下会生成：
- `obs-samples-3.25.7.jar` - Samples模块JAR包

## 依赖说明

### 核心依赖
- `com.squareup.okhttp3:okhttp:4.12.0` - HTTP客户端
- `com.squareup.okio:okio:3.8.0` - OkHttp的IO库
- `com.fasterxml.jackson.core:jackson-core:2.13.4` - Jackson核心
- `com.fasterxml.jackson.core:jackson-databind:2.13.4` - Jackson数据绑定
- `com.fasterxml.jackson.core:jackson-annotations:2.13.4` - Jackson注解
- `org.apache.logging.log4j:log4j-core:2.17.1` - Log4j核心（测试范围）
- `org.apache.logging.log4j:log4j-api:2.17.1` - Log4j API（测试范围）

### 测试依赖
- `junit:junit:4.13.2` - 单元测试框架
- `org.mockito:mockito-core:3.12.4` - Mock框架
- `org.powermock:powermock-module-junit4:2.0.9` - PowerMock模块
- `org.powermock:powermock-api-mockito2:2.0.9` - PowerMock API
- `com.squareup.okhttp3:mockwebserver:4.12.0` - MockWebServer
- `org.mock-server:mockserver-netty:5.15.0` - MockServer

## 常见问题

### SSL证书问题
如果遇到Maven仓库SSL证书问题，可以使用项目提供的Maven配置文件：
```bash
mvn clean package
```

### 内存不足
如果构建时遇到内存不足问题，可以增加Maven内存：
```bash
export MAVEN_OPTS="-Xmx2048m -XX:MaxPermSize=512m"
mvn clean package
```

### 依赖冲突
如果遇到依赖冲突，可以查看依赖树：
```bash
mvn dependency:tree
```

## CI/CD集成

### Jenkins Pipeline示例
```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'mvn clean package -DskipTests'
            }
        }
        stage('Test') {
            steps {
                sh 'mvn test'
            }
        }
        stage('Integration Test') {
            steps {
                sh 'mvn verify -DskipTests'
            }
        }
        stage('Deploy') {
            steps {
                sh 'mvn clean deploy'
            }
        }
    }
}
```

### GitHub Actions示例
```yaml
name: Build and Test

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up JDK 1.8
        uses: actions/setup-java@v2
        with:
          java-version: '1.8'
          distribution: 'adopt'
      - name: Build with Maven
        run: mvn clean package
      - name: Run tests
        run: mvn test
```