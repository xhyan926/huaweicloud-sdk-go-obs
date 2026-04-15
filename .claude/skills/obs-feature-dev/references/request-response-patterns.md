# 请求/响应模型代码模板

本文档提供从项目代码库中提取的精确代码模板，以 BucketTrash 功能为参考范例。

## 目录结构

```
src/main/java/com/obs/services/model/<feature>/
├── <Feature>Configuration.java      # 数据模型
├── Get<Feature>Request.java         # GET 请求
├── Set<Feature>Request.java         # PUT 请求
├── Delete<Feature>Request.java      # DELETE 请求
└── Get<Feature>Result.java          # GET 响应
```

---

## Configuration 模板

基于 `BucketTrashConfiguration.java`：

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.model.<feature>;

public class <Feature>Configuration {
    // 私有字段
    private int field1;
    private String field2;

    // 带参构造函数
    public <Feature>Configuration(int field1) {
        this.field1 = field1;
    }

    // Getter 和 Setter
    public int getField1() {
        return field1;
    }

    public void setField1(int field1) {
        this.field1 = field1;
    }

    public String getField2() {
        return field2;
    }

    public void setField2(String field2) {
        this.field2 = field2;
    }
}
```

**规则**：
- 纯 POJO，不继承任何基类
- 字段使用驼峰命名
- 提供完整的有参构造函数和 getter/setter

---

## GET Request 模板

基于 `GetBucketTrashRequest.java`：

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.model.<feature>;

import com.obs.services.model.BaseBucketRequest;
import com.obs.services.model.HttpMethodEnum;

public class Get<Feature>Request extends BaseBucketRequest {
    {
        httpMethod = HttpMethodEnum.GET;
    }

    // 可选：带桶名的构造函数
    public Get<Feature>Request(String bucketName) {
        super(bucketName);
    }

    // 可选：无参构造函数（后续通过 setBucketName 设置）
    public Get<Feature>Request() {
    }
}
```

**规则**：
- 继承 `BaseBucketRequest`
- 使用实例初始化块 `{}` 设置 `httpMethod`
- 无业务字段（桶名继承自父类）

---

## PUT Request 模板

基于 `SetBucketTrashRequest.java`：

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.model.<feature>;

import com.obs.services.model.BaseBucketRequest;
import com.obs.services.model.HttpMethodEnum;

public class Set<Feature>Request extends BaseBucketRequest {
    {
        httpMethod = HttpMethodEnum.PUT;
    }

    private <Feature>Configuration <feature>Configuration;

    public Set<Feature>Request(String bucketName, <Feature>Configuration <feature>Configuration) {
        super(bucketName);
        this.<feature>Configuration = <feature>Configuration;
    }

    public <Feature>Configuration get<Feature>Configuration() {
        return <feature>Configuration;
    }

    public void set<Feature>Configuration(<Feature>Configuration <feature>Configuration) {
        this.<feature>Configuration = <feature>Configuration;
    }
}
```

**规则**：
- 继承 `BaseBucketRequest`
- 持有 Configuration 对象
- `httpMethod = HttpMethodEnum.PUT`
- 提供带桶名和 Configuration 的构造函数

---

## DELETE Request 模板

基于 `DeleteBucketTrashRequest.java`：

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.model.<feature>;

import com.obs.services.model.BaseBucketRequest;
import com.obs.services.model.HttpMethodEnum;

public class Delete<Feature>Request extends BaseBucketRequest {
    {
        httpMethod = HttpMethodEnum.DELETE;
    }

    public Delete<Feature>Request(String bucketName) {
        super(bucketName);
    }

    public Delete<Feature>Request() {
    }
}
```

**规则**：
- 继承 `BaseBucketRequest`
- 与 GET Request 结构相同，仅 `httpMethod` 不同
- 无业务字段

---

## GET Result 模板

基于 `GetBucketTrashResult.java`：

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.model.<feature>;

import com.obs.services.model.HeaderResponse;

public class Get<Feature>Result extends HeaderResponse {
    private <Feature>Configuration <feature>Configuration;

    public <Feature>Configuration get<Feature>Configuration() {
        return <feature>Configuration;
    }

    public void set<Feature>Configuration(<Feature>Configuration <feature>Configuration) {
        this.<feature>Configuration = <feature>Configuration;
    }
}
```

**规则**：
- 继承 `HeaderResponse`（自动包含 HTTP 响应头和状态码）
- 持有 Configuration 对象
- 仅提供 getter/setter，无业务逻辑

---

## XML Builder 模板

基于 `BucketTrashConfigurationXMLBuilder.java`，放置于 `src/main/java/com/obs/services/internal/xml/`：

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.internal.xml;

import com.obs.log.ILogger;
import com.obs.log.LoggerBuilder;
import com.obs.services.exception.ObsException;
import com.obs.services.model.<feature>.<Feature>Configuration;

public class <Feature>ConfigurationXMLBuilder extends ObsSimpleXMLBuilder {
    private static final ILogger log = LoggerBuilder.getLogger("com.obs.services.ObsClient");
    private final static String <FEATURE>_CONFIGURATION = "<Feature>Configuration";
    public final static String FIELD1 = "Field1";
    public final static String FIELD2 = "Field2";

    public String buildXML(<Feature>Configuration config) {
        checkBucketPublicAccessBlock(config);
        startElement(<FEATURE>_CONFIGURATION);
        startElement(FIELD1);
        append(config.getField1());
        endElement(FIELD1);
        startElement(FIELD2);
        append(config.getField2());
        endElement(FIELD2);
        endElement(<FEATURE>_CONFIGURATION);
        return getXmlBuilder().toString();
    }

    @Override
    protected void checkBucketPublicAccessBlock(Object config) {
        if (config == null) {
            String errorMessage = "<feature>Configuration is null, failed to build request XML!";
            log.error(errorMessage);
            throw new ObsException(errorMessage);
        }
    }
}
```

**规则**：
- 继承 `ObsSimpleXMLBuilder`
- XML 元素名常量使用 `UPPER_SNAKE_CASE`，且声明为 `public final static`（供 SAX Handler 引用）
- 使用 `startElement()` / `append()` / `endElement()` 构建 XML
- 空值校验抛出 `ObsException`

---

## SAX Handler 模板

添加到 `XmlResponsesSaxParser.java` 的内部类：

```java
// 导入 XML Builder 的常量
import static com.obs.services.internal.xml.<Feature>ConfigurationXMLBuilder.*;

// 在 XmlResponsesSaxParser 类内部添加
public static class <Feature>ConfigurationXMLHandler extends DefaultXmlHandler {
    private String field1;
    private String field2;

    public String getField1() {
        return field1;
    }

    public String getField2() {
        return field2;
    }

    @Override
    public void endElement(String name, String elementText) {
        if (FIELD1.equals(name)) {
            field1 = elementText;
        } else if (FIELD2.equals(name)) {
            field2 = elementText;
        }
    }
}
```

**规则**：
- 继承 `DefaultXmlHandler`
- 字段类型使用 `String`（在 Service 层再做类型转换）
- 使用 XML Builder 中定义的常量进行元素名匹配
- 只重写 `endElement()` 方法

---

## SpecialParamEnum 注册模板

在 `SpecialParamEnum.java` 中添加枚举项：

```java
<FEATURE_ENUM>("x-obs-<feature>"),
```

**查找位置**: 找到类似的枚举项（如 `OBS_TRASH`），在其附近添加。

---

## IObsClient 接口声明模板

在 `IObsClient.java` 中添加方法声明（带 JavaDoc）：

```java
/**
 * Set <feature> configuration for a bucket.
 *
 * @param request
 *            Request parameters
 * @return Common response headers
 * @throws ObsException
 *             OBS SDK self-defined exception, thrown when the interface
 *             fails to be called or access to OBS fails
 */
HeaderResponse set<Feature>(Set<Feature>Request request) throws ObsException;

/**
 * Get <feature> configuration of a bucket.
 *
 * @param request
 *            Request parameters
 * @return <Feature> configuration result
 * @throws ObsException
 *             OBS SDK self-defined exception, thrown when the interface
 *             fails to be called or access to OBS fails
 */
Get<Feature>Result get<Feature>(Get<Feature>Request request) throws ObsException;

/**
 * Delete <feature> configuration of a bucket.
 *
 * @param request
 *            Request parameters
 * @return Common response headers
 * @throws ObsException
 *             OBS SDK self-defined exception, thrown when the interface
 *             fails to be called or access to OBS fails
 */
HeaderResponse delete<Feature>(Delete<Feature>Request request) throws ObsException;
```
