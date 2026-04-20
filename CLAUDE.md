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

### 规则文件 (.claude/rules/)

根据当前任务类型，加载相应的规则文件：

| 规则 | 文件 | 加载场景 |
|------|------|----------|
| 核心开发原则 | `core-development-principles.md` | 代码开发活动 |
| 开发工作流 | `development-workflow.md` | 开发、构建、测试、提交 |
| 代码质量 | `code-quality.md` | 编写、修改、审查代码 |
| AI代码生成 | `ai-code-generation.md` | AI助手生成代码 |
| 集成测试 | `integration-test.md` | 编写/调试集成测试或新增 SpecialParamEnum |

### 技术文档 (.claude/docs/)

根据需要了解的深度，加载相应的技术文档：

| 文档 | 文件 | 加载场景 |
|------|------|----------|
| 架构文档 | `architecture.md` | 了解整体架构、设计模式、核心组件 |
| 构建指南 | `build-guide.md` | 构建项目、运行测试、生成文档 |
| Maven插件 | `maven-plugins.md` | Maven插件配置、自定义构建流程 |
| 依赖管理 | `dependency-management.md` | 依赖问题、添加新依赖 |
| 开发规范 | `development-standards.md` | 团队规范、代码审查标准 |

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

> 详细约束见 `.claude/rules/` 下对应规则文件，以下为核心约束摘要：

- Maven命令必须包含 `-s` settings参数
- 测试覆盖率100%，命名 `should_X_when_Y`，使用参数化测试
- 符合SOLID原则，使用Java泛型，公共方法有JavaDoc
- AI生成代码须加 `@AIGenerated` 注解，动态日期和git用户信息

## 常见任务场景

| 场景 | 需加载的规则/文档 | 关键步骤 |
|------|-------------------|----------|
| 添加新功能 | `development-workflow` + `code-quality` + `architecture` | SOLID设计 → 单元测试(100%覆盖) → 验证清单 |
| 修复缺陷 | `development-workflow` + `build-guide` | 定位根因 → 修复代码 → 更新测试 → 验证清单 |
| AI生成代码 | `ai-code-generation` + `code-quality` | @AIGenerated注解 → 强制验证清单 |
| 解决依赖问题 | `dependency-management` | `mvn dependency:tree` → 分析冲突 → 解决 |
| 代码审查 | `code-quality` + `development-standards` | Review CheckList → 改进建议 |

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