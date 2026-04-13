# OBS Java SDK

华为云对象存储服务（OBS）Java SDK，提供简单易用的API接口，支持Java和Android平台。

## 项目结构

```
obs-sdk-java/
├── pom.xml                          # 统一的Maven配置文件
├── src/
│   ├── main/
│   │   ├── java/                   # SDK业务代码
│   │   └── resources/              # 资源文件
│   ├── test/
│   │   ├── java/                   # 单元测试
│   │   └── resources/              # 单元测试资源
│   └── it/
│       ├── java/                   # 集成测试
│       └── resources/              # 集成测试资源
├── samples/                         # 示例代码
│   ├── java/                       # Java示例
│   └── android/                    # Android示例
├── scripts/                         # 辅助脚本
├── docs/                           # 文档
└── README.md                       # 项目说明
```

## 快速开始

### 环境要求

- JDK 1.8+
- Maven 3.6+

### 构建项目

```bash
# 构建Java版本（默认）
mvn clean package

# 构建Android版本
mvn clean package -Pandroid

# 只运行单元测试
mvn test

# 只运行集成测试
mvn verify -DskipTests

# 运行所有测试
mvn verify

# 跳过所有测试
mvn package -DskipTests

# 安装到本地仓库
mvn clean install

# 生成Javadoc
mvn javadoc:javadoc

# 清理构建
mvn clean
```

## 依赖管理

项目使用统一的依赖版本管理，主要依赖包括：

- OkHttp 4.12.0 - HTTP客户端
- Jackson 2.15.4 - JSON处理
- Log4j 2.24.3 - 日志框架
- JUnit 4.13.2 - 单元测试框架
- Mockito 3.12.4 - Mock框架

## 开发规范

### 零、平时要关注的外部网站-事务
1. 论坛上的问题：https://bbs.huaweicloud.com/forum/forum-620-1.html
2. 云问答：https://bbs.huaweicloud.com/ask/
3. 开源社区(issues)：https://github.com/huaweicloud/huaweicloud-sdk-java-obs
4. 客户声音处理（NPS）：https://ivoc.huaweicloud.com/voc_manager/#!/app/voice/overview/todoList?vStateList=&deStateList=

### 一、开发红线  
1.	禁止未经许可屏蔽静态检查工具告警（如：代码中使用“/*lint -e* */”、“//NOPMD”等方法绕过门禁）。
2.	禁止由外包同事确认处理CodeDEX和FOSSID问题。
3.	禁止搭便车修改代码库上的代码（如：代码上库无对应需求单、问题单，或夹带在别的需求单、问题单中上库）。
4.	禁止开发未公开接口（如：未在相关资料中明确说明的API接口、网络端口、人机/机机账号等）。
5.	禁止未经CCB批准更改接口（如：API接口、进程间接口、持久化数据格式、话单等）。
6.	禁止在CodeClub配置非正式任命Committer作为Approver，代码提交应至少经过2位Committer检视通过后方可上库。
7.	禁止未经评审私自调整门禁阈值（见《附件一：静态检查门禁阈值.xlsx》）。
8.	禁止违反安全编码规范（见《附件二：TOP5安全编码问题.xlsx》）。

## 二、红线应用  
因违反上述质量红线导致P3及以上事件的，或因违反上述质量红线导致被通报且经提醒不积极改正的，纳入个人关键黑事件，当期考评降等。

## 三、分支命名规范
1. 主干	master
2. 开发分支	dev_分支描述([版本号][特性名]/[局点])
3. release分支	br_版本号_[分支描述]
4. 线下分支	IT_分支描述([dev/release][版本号]/[特性名]/[局点])
5. 个人分支	工号_分支描述

## 四、代码检视规范
### 1. review CheckList：
```
- 高安全	符合公司安全编码规范；考虑权限、数据和资源保护；防攻击，能够抵御网络攻击，不泄漏敏感数据。	
- 可读性	易于理解和维护。对象/函数/类/变量“名”副其实，命名能够体现用途。	
- 代码简洁	通过提取公共函数、公共类，避免重复代码。	
- 资源释放	适时释放资源，包括内存、连接、文件句柄等资源，避免资源泄露。	
- 可扩展	适度应对可能的需求变化，OOP，适当的设计模式，拒绝过度设计。	
- 性能	加锁范围尽量小，降低对性能的影响；新增进程间通信、线程间互斥，对性能的影响是否做了评估和测试。	
- 容错性	异常路径的处理，Exception的捕获等。	
- OPS能力	适量的日志，日志要在“点子上”，拒绝无意义或意义不全的日志。出问题能快速定位。	
- 升级兼容	进程间新增消息或者已有消息增删字段，是否会导致升级时业务中断；indexLayer新增表或者增删字段，是否会导致升级时业务中断、导致回退版本后业务中断。	
```

### 2. 没有检视出问题时，MR的合入：
```
- MR执行Accept合入时，合入人检查没有检视出问题，需要至少一个Approver在MR的Discussion页面写明：已检视，对于本MR引起的P事件承担同等责任。
```
