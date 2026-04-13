# OBS Java SDK 架构文档

## 项目概述

华为云对象存储服务（OBS）Java SDK，提供简单易用的API接口，支持Java和Android平台。

## 核心架构

### 分层架构

```
┌─────────────────────────────────────┐
│         API接口层                    │
│  (IObsClient, IObsBucketExtendClient)│
├─────────────────────────────────────┤
│         客户端实现层                  │
│  (ObsClient, ObsClientAsync)         │
├─────────────────────────────────────┤
│         服务抽象层                    │
│  (AbstractClient, AbstractBucketClient)│
├─────────────────────────────────────┤
│         内部服务层                    │
│  (RestStorageService, ObsService)    │
├─────────────────────────────────────┤
│         工具和工具层                  │
│  (认证、IO处理、XML解析等)             │
└─────────────────────────────────────┘
```

### 核心组件

#### 1. 客户端接口
- `IObsClient`: 同步客户端接口
- `IObsClientAsync`: 异步客户端接口
- `IObsBucketExtendClient`: 桶扩展操作接口
- `IFSClient`: 文件系统接口

#### 2. 客户端实现
- `ObsClient`: 同步客户端主实现
- `ObsClientAsync`: 异步客户端实现
- `SecretFlexibleObsClient`: 密钥灵活客户端
- `CryptoObsClient`: 加密客户端

#### 3. 抽象基类
- `AbstractClient`: 客户端抽象基类
- `AbstractBucketClient`: 桶操作抽象基类
- `AbstractObjectClient`: 对象操作抽象基类
- `AbstractBatchClient`: 批量操作抽象基类

#### 4. 内部服务
- `RestStorageService`: REST存储服务
- `RestConnectionService`: REST连接服务
- `ObsService`: OBS服务基类
- `UploadResumableClient`: 断点续传上传
- `DownloadResumableClient`: 断点续传下载

#### 5. 认证模块
- `IObsCredentialsProvider`: 认证提供者接口
- `BasicObsCredentialsProvider`: 基本认证提供者
- `EcsObsCredentialsProvider`: ECS认证提供者
- `EnvironmentVariableObsCredentialsProvider`: 环境变量认证提供者
- `OBSCredentialsProviderChain`: 认证链

#### 6. 加密模块
- `CryptoObsClient`: 加密客户端
- `CTRCipherGenerator`: CTR密码生成器
- `CtrRSACipherGenerator`: CTR RSA密码生成器

## 包结构

```
com.obs
├── services/                    # 服务层
│   ├── internal/               # 内部实现
│   │   ├── service/           # 内部服务
│   │   ├── utils/             # 工具类
│   │   ├── io/                # IO处理
│   │   ├── security/          # 安全相关
│   │   ├── handler/           # 处理器
│   │   └── xml/               # XML处理
│   ├── model/                 # 数据模型
│   ├── crypto/                # 加密模块
│   └── exception/             # 异常定义
├── log/                        # 日志模块
└── integrated_test/            # 集成测试
```

## 设计模式

### 1. 工厂模式
- `LoggerBuilder`: 日志构建器
- `ObsConfiguration`: 配置构建

### 2. 策略模式
- `IAuthentication`: 认证策略接口
- `DefaultAuthentication`: 默认认证策略
- `ObsAuthentication`: OBS认证策略
- `V2Authentication`: V2认证策略
- `V4Authentication`: V4认证策略

### 3. 责任链模式
- `OBSCredentialsProviderChain`: 认证提供者链
- `DefaultCredentialsProviderChain`: 默认认证链

### 4. 模板方法模式
- `AbstractClient`: 定义客户端操作模板
- `AbstractBucketClient`: 定义桶操作模板
- `AbstractObjectClient`: 定义对象操作模板

### 5. 适配器模式
- `ObsConvertor`: 对象转换适配器
- `V2Convertor`: V2协议适配器
- `V2BucketConvertor`: V2桶操作适配器

## 关键特性

### 1. 多协议支持
- 支持V2和V4认证协议
- 支持S3兼容协议
- 支持路径样式和虚拟主机样式访问

### 2. 断点续传
- `UploadResumableClient`: 支持断点续传上传
- `DownloadResumableClient`: 支持断点续传下载
- 自动记录和恢复传输状态

### 3. 加密存储
- 客户端加密支持
- 支持CTR和CTR RSA加密模式
- 密钥管理和轮换

### 4. 异步操作
- `ObsClientAsync`: 异步客户端
- 回调机制支持
- 并发控制

### 5. 进度监控
- `ProgressListener`: 进度监听接口
- `ProgressManager`: 进度管理器
- `ConcurrentProgressManager`: 并发进度管理

## 依赖管理

### 核心依赖
- OkHttp 4.12.0: HTTP客户端
- Jackson 2.15.4: JSON处理
- Log4j 2.24.3: 日志框架（仅测试）

### 测试依赖
- JUnit 4.13.2: 单元测试框架
- Mockito 3.12.4: Mock框架
- PowerMock 2.0.9: 静态方法Mock
- MockWebServer 4.12.0: HTTP服务Mock

## 构建配置

### Maven插件
- `maven-compiler-plugin`: 代码编译
- `maven-shade-plugin`: 依赖打包和重定位
- `maven-surefire-plugin`: 单元测试
- `maven-failsafe-plugin`: 集成测试
- `maven-javadoc-plugin`: 文档生成
- `maven-source-plugin`: 源码打包

### Profile配置
- `java`: Java版本构建（默认）
- `android`: Android版本构建

### 依赖重定位
为避免与应用依赖冲突，以下依赖被重定位：
- `com.jamesmurty.utils` → `shade.jamesmurty.utils`
- `okhttp3` → `shade.okhttp3`
- `okio` → `shade.okio`
- `com.fasterxml.jackson.*` → `shade.fasterxml.jackson.*`

## 扩展点

### 1. 自定义认证
实现 `IObsCredentialsProvider` 接口

### 2. 自定义日志
实现 `ILogger` 接口

### 3. 自定义进度监听
实现 `ProgressListener` 接口

### 4. 自定义加密
实现加密接口并扩展 `CryptoObsClient`

## 性能优化

### 1. 连接池
- OkHttp连接池管理
- 可配置连接池大小

### 2. 并发控制
- 多线程上传下载
- 可配置并发数

### 3. 缓存机制
- 元数据缓存
- DNS缓存

### 4. 压缩传输
- 支持GZIP压缩
- 可配置压缩策略

## 安全特性

### 1. 认证安全
- 多种认证方式支持
- AK/SK安全管理
- STS Token支持

### 2. 传输安全
- HTTPS支持
- 证书验证
- 加密传输

### 3. 数据安全
- 客户端加密
- 服务端加密支持
- 密钥轮换

## 兼容性

### Java版本
- 最低要求: Java 8
- 推荐版本: Java 8+

### Android版本
- 最低要求: Android 4.0+
- 推荐版本: Android 5.0+

### 平台支持
- Windows
- Linux
- macOS
- Android