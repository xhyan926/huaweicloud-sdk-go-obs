# 兼容性规则（最高优先级）

**优先级**: 本规则优先级高于所有其他规则，发生冲突时以本规则为准。

**应用场景**: 所有涉及代码修改的活动，无例外。

## 禁止事项

1. **禁止修改已发布公共 API 的签名**
   - 不得更改 `IObsClient`、`IObsBucketExtendClient`、`IObsObjectExtendClient` 等公共接口中已有方法的名称、参数列表、返回类型、异常声明
   - 不得向已有接口方法添加 `default` 实现以外的任何变更
   - 新增接口方法不受此限制，但必须提供所有实现类中的具体实现

2. **禁止删除或重命名已发布的公共类、方法、字段、常量**
   - 不得删除任何 `public` 或 `protected` 的类、接口、枚举、方法、字段、常量
   - 不得重命名上述任何元素（包括包路径）
   - 如需废弃，必须使用 `@Deprecated` 注解标注并保留原实现

3. **禁止改变已发布方法的行为语义**
   - 不得修改已有方法的内部逻辑导致输入输出关系发生变化
   - 不得修改已有方法的异常抛出行为（新增异常类型、改变异常条件）
   - 不得修改已有的默认值、校验规则、错误消息格式
   - 修复缺陷时不在此列，但必须满足以下条件：明确标注修复的问题单号、修复后行为与协议文档或设计意图一致、不引入新的异常类型或异常条件

4. **禁止修改已发布的序列化/反序列化格式**
   - 不得更改 XML/JSON 的元素名称、命名空间、层级结构
   - 不得更改 HTTP 请求/响应头的名称和格式约定
   - 新增可选元素不受此限制

5. **禁止修改依赖项的版本范围导致现有用户构建失败**
   - 不得提升最低 JDK 版本要求
   - 不得提升 Maven 插件的最低版本要求
   - 不得移除或替换已有的依赖项

## 正向示例

```java
// 正确：新增接口方法（不修改已有方法）
public interface IObsClient {
    // 已有方法保持不变
    HeaderResponse deleteBucketCors(String bucketName);
    // 新增方法
    OptionsInfoResult optionsBucket(OptionsInfoRequest request);  // 安全
}

// 正确：废弃旧方法而非删除
@Deprecated
public OldResult oldMethod(String param) {
    // 保留原实现，委托到新方法
    return newMethod(new Request(param)).toOldResult();
}
```

## 反向示例

```java
// 错误：修改已有方法签名
HeaderResponse deleteBucketCors(String bucketName, String extraParam);  // 破坏所有调用方

// 错误：修改返回类型
BucketCors getBucketCors(String bucketName) → BucketCorsV2 getBucketCors(...);  // 破坏编译

// 错误：删除公共方法
// 从 IObsClient 中移除 createBucket()  // 破坏所有实现类和调用方
```

## 检查清单

任何代码修改提交前必须确认：
- [ ] 未修改任何已有公共接口的方法签名
- [ ] 未删除或重命名任何 public/protected 元素
- [ ] 未改变已有方法的输入输出行为
- [ ] 未修改 XML/JSON 序列化格式中的已有元素
- [ ] 未提升依赖项或运行环境的最低版本要求
