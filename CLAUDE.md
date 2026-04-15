# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

华为云对象存储服务（OBS）Java SDK，提供简单易用的API接口，支持Java和Android平台。

## 快速开始指南

### 环境要求
- JDK 1.8+
- Maven 3.6+

### 基本构建命令
```bash
# 构建Java版本（默认）
mvn clean package

# 构建Android版本
mvn clean package -Pandroid

# 运行测试
mvn test

# 运行集成测试
mvn verify
```

## 渐进式披露文档导航

### 📋 规则文件 (.claude/rules/)

根据当前任务类型，加载相应的规则文件：

#### 1. 开发工作流规则
**加载场景**: 当进行代码的开发活动时加载
**文件**: `.claude/rules/core-development-principles.md`
**主要内容**:
- 核心的开发原则

#### 2. 开发工作流规则
**加载场景**: 当进行代码开发、构建、测试、提交等开发活动时
**文件**: `.claude/rules/development-workflow.md`
**主要内容**:
- Maven命令规范（必须包含settings参数）
- 静态检查要求
- 开发红线
- 验证清单

#### 3. 代码质量规则
**加载场景**: 当编写、修改、审查代码时
**文件**: `.claude/rules/code-quality.md`
**主要内容**:
- SOLID原则要求
- 代码风格规范
- 测试质量要求
- 质量检查清单

#### 4. AI代码生成规则
**加载场景**: 当使用AI助手生成代码时
**文件**: `.claude/rules/ai-code-generation.md`
**主要内容**:
- @AIGenerated注解规范
- 动态日期和git用户信息获取
- 参数化测试要求
- 强制验证清单

#### 5. 集成测试规则
**加载场景**: 当编写、修改、调试集成测试（`*IT.java`）或新增 `SpecialParamEnum` 子资源参数时
**文件**: `.claude/rules/integration-test.md`
**主要内容**:
- 签名白名单注册要求
- 多环境配置一致性
- 测试资源生命周期管理
- 桶名转换一致性
- 断言与服务端行为对齐
- 集成测试调试清单

### 📚 技术文档 (.claude/docs/)

根据需要了解的深度，加载相应的技术文档：

#### 1. 架构文档
**加载场景**: 当需要了解项目整体架构、设计模式、核心组件时
**文件**: `.claude/docs/architecture.md`
**主要内容**:
- 分层架构设计
- 核心组件说明
- 包结构组织
- 设计模式应用
- 关键特性介绍

#### 2. 构建指南
**加载场景**: 当需要构建项目、运行测试、生成文档时
**文件**: `.claude/docs/build-guide.md`
**主要内容**:
- 环境要求
- 构建命令详解
- Profile配置
- 构建产物说明
- 常见问题解决

#### 3. Maven插件详解
**加载场景**: 当需要了解Maven插件配置、自定义构建流程时
**文件**: `.claude/docs/maven-plugins.md`
**主要内容**:
- 核心插件配置
- 插件执行顺序
- 依赖关系图
- 最佳实践

#### 4. 依赖管理
**加载场景**: 当遇到依赖问题、需要添加新依赖时
**文件**: `.claude/docs/dependency-management.md`
**主要内容**:
- 依赖健康度评估
- 版本冲突分析
- 优化建议
- 常用依赖管理命令

#### 5. 开发规范
**加载场景**: 当需要了解团队开发规范、代码审查标准时
**文件**: `.claude/docs/development-standards.md`
**主要内容**:
- 开发红线
- 分支命名规范
- 代码检视规范
- 提交信息规范
- 外部资源链接

## 项目结构

```
obs-sdk-java/
├── src/
│   ├── main/java/              # SDK业务代码
│   │   └── com/obs/
│   │       ├── services/       # 服务层（客户端、接口）
│   │       ├── log/            # 日志模块
│   │       └── ...
│   ├── test/java/              # 单元测试（*Test.java）
│   └── it/java/                # 集成测试（*IT.java）
├── samples/                    # 示例代码
│   ├── java/                   # Java示例
│   └── android/                # Android示例
├── .claude/                    # Claude知识库
│   ├── rules/                  # 规则文件
│   └── docs/                   # 技术文档
├── pom.xml                     # Maven配置
└── README.md                   # 项目说明
```

## 核心架构概念

### 分层架构
- **API接口层**: `IObsClient`, `IObsBucketExtendClient`
- **客户端实现层**: `ObsClient`, `ObsClientAsync`
- **服务抽象层**: `AbstractClient`, `AbstractBucketClient`
- **内部服务层**: `RestStorageService`, `ObsService`
- **工具和工具层**: 认证、IO处理、XML解析等

### 关键设计模式
- **工厂模式**: `LoggerBuilder`, `ObsConfiguration`
- **策略模式**: `IAuthentication` 及其实现
- **责任链模式**: `OBSCredentialsProviderChain`
- **模板方法模式**: `AbstractClient` 及其子类
- **适配器模式**: `ObsConvertor`, `V2Convertor`

## 重要约束

### Maven命令约束
所有Maven命令必须包含settings参数：
```bash
mvn [command] -s ./esdk_obs_java_android_en/CI/settings.xml -T 1C
```

### 测试约束
- 新增代码分支覆盖率必须达到100%
- 测试方法命名必须遵循 `should_[ExpectedBehavior]_when_[Condition]` 格式
- 必须使用参数化测试避免重复代码

### 代码质量约束
- 必须符合SOLID原则
- 不得使用原始类型，必须使用Java泛型
- 公共方法必须有JavaDoc注释

### AI生成代码约束
- 必须添加 `@AIGenerated` 注解
- 必须使用动态日期和git用户信息
- 必须通过编译和测试验证

## 常见任务场景

### 场景1: 添加新功能
1. 加载 `.claude/rules/development-workflow.md` - 了解开发流程
2. 加载 `.claude/rules/code-quality.md` - 了解代码质量要求
3. 加载 `.claude/docs/architecture.md` - 了解架构设计
4. 按照SOLID原则设计接口和实现
5. 编写单元测试（覆盖率100%）
6. 执行验证清单

### 场景2: 修复缺陷
1. 加载 `.claude/rules/development-workflow.md` - 了解开发流程
2. 加载 `.claude/docs/build-guide.md` - 了解构建和测试
3. 定位问题根因
4. 编写修复代码
5. 编写或更新测试
6. 执行验证清单

### 场景3: 使用AI生成代码
1. 加载 `.claude/rules/ai-code-generation.md` - 了解AI生成规范
2. 加载 `.claude/rules/code-quality.md` - 了解代码质量要求
3. 生成代码并添加@AIGenerated注解
4. 执行强制验证清单

### 场景4: 解决依赖问题
1. 加载 `.claude/docs/dependency-management.md` - 了解依赖管理
2. 运行 `mvn dependency:tree` 查看依赖树
3. 分析冲突原因
4. 按照优化建议解决

### 场景5: 代码审查
1. 加载 `.claude/rules/code-quality.md` - 了解代码质量要求
2. 加载 `.claude/docs/development-standards.md` - 了解审查标准
3. 按照Review CheckList进行检查
4. 提供改进建议

## 外部资源

- **论坛**: https://bbs.huaweicloud.com/forum/forum-620-1.html
- **云问答**: https://bbs.huaweicloud.com/ask/
- **开源社区**: https://github.com/huaweicloud/huaweicloud-sdk-java-obs
- **客户声音**: https://ivoc.huaweicloud.com/voc_manager/#!/app/voice/overview/todoList

## 注意事项

1. **渐进式加载**: 根据任务需要加载相应文档，避免一次性加载过多信息
2. **规则优先**: 规则文件中的禁止事项必须严格遵守
3. **验证必做**: 任何代码修改都必须完成相应的验证步骤
4. **质量第一**: 代码质量优先于开发速度
5. **向后兼容**: 修改公共API时需考虑向后兼容性