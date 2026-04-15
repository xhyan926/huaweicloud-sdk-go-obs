# 测试约定和示例

本文档提供 OBS SDK 单元测试和集成测试的约定与代码模板。

## 测试文件结构

```
src/test/java/com/obs/services/model/<feature>/
└── <Feature>UnitTest.java          # 单元测试

src/it/java/com/obs/integrated_test/<feature>/
└── <Feature>IT.java                # 集成测试
```

---

## 一、单元测试

基于 MockServer + JUnit 4，隔离外部依赖，验证参数校验和请求构造。

基于 `AbstractBucketClientQosUnitTest.java` 提取的完整模板：

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.services.model.<feature>;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.test.TestTools;

import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

public class <Feature>UnitTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-<feature>";

    @BeforeClass
    public static void setMockServer() {
        mockServer = ClientAndServer.startClientAndServer();
        System.setProperty(PROXY_HOST_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_PROPERTY_NAME, "" + mockServer.getLocalPort());
        System.setProperty(PROXY_HOST_S_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_S_PROPERTY_NAME, "" + mockServer.getLocalPort());
    }

    @AfterClass
    public static void clearEnv() {
        mockServer.close();
        System.clearProperty(PROXY_HOST_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_PROPERTY_NAME);
        System.clearProperty(PROXY_HOST_S_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_S_PROPERTY_NAME);
    }

    // ==================== GET 测试 ====================

    @Test
    public void should_throw_exception_when_getRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Get<Feature>Request is null");
        obsClient.get<Feature>(null);
    }

    @Test
    public void should_throw_exception_when_getRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Get<Feature>Request request = new Get<Feature>Request();
        request.setBucketName(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.get<Feature>(request);
    }

    @Test
    public void should_succeed_when_getFeature_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        // 准备 Mock XML 响应
        String validResponseXml = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
                + "<<Feature>Configuration>"
                + "  <Field1>value1</Field1>"
                + "  <Field2>value2</Field2>"
                + "</<Feature>Configuration>";

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("<special-param-name>"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/xml;charset=utf-8")
                        .withBody(validResponseXml));

        Get<Feature>Request request = new Get<Feature>Request(bucketNameForTest);
        Get<Feature>Result result = obsClient.get<Feature>(request);

        Assert.assertNotNull(result);
        Assert.assertNotNull(result.get<Feature>Configuration());
        // 添加字段断言
        Assert.assertEquals("value1", result.get<Feature>Configuration().getField1());
    }

    // ==================== PUT 测试 ====================

    @Test
    public void should_throw_exception_when_setRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Set<Feature>Request is null");
        obsClient.set<Feature>(null);
    }

    @Test
    public void should_throw_exception_when_setRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        <Feature>Configuration config = new <Feature>Configuration(/* params */);
        Set<Feature>Request request = new Set<Feature>Request(bucketNameForTest, config);
        request.setBucketName(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.set<Feature>(request);
    }

    @Test
    public void should_succeed_when_setFeature_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        <Feature>Configuration config = new <Feature>Configuration(/* params */);
        Set<Feature>Request request = new Set<Feature>Request(bucketNameForTest, config);
        HeaderResponse response = obsClient.set<Feature>(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    // ==================== DELETE 测试 ====================

    @Test
    public void should_throw_exception_when_deleteRequest_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Delete<Feature>Request is null");
        obsClient.delete<Feature>(null);
    }

    @Test
    public void should_throw_exception_when_deleteRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Delete<Feature>Request request = new Delete<Feature>Request(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.delete<Feature>(request);
    }

    @Test
    public void should_succeed_when_deleteFeature_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 204;
        mockServer.reset();
        mockServer.when(request().withMethod("DELETE").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        Delete<Feature>Request request = new Delete<Feature>Request(bucketNameForTest);
        HeaderResponse response = obsClient.delete<Feature>(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }
}
```

---

## 测试约定

### 命名规范

- **文件命名**: `<Feature>UnitTest.java`
- **方法命名**: `should_[ExpectedBehavior]_when_[Condition]`
- **常量命名**: `UPPER_SNAKE_CASE`

### MockServer 使用规则

1. `@BeforeClass` 启动 MockServer 并设置系统属性
2. `@AfterClass` 关闭 MockServer 并清理系统属性
3. 每个测试方法前调用 `mockServer.reset()` 清除之前的 mock 配置
4. GET 测试需要匹配 `queryStringParameter`（子资源参数）
5. PUT/DELETE 测试仅需匹配 HTTP 方法

### ExpectedException 规则

```java
@Rule
public ExpectedException expectedException = ExpectedException.none();
```

使用方式：
```java
expectedException.expect(IllegalArgumentException.class);
expectedException.expectMessage("预期的异常消息");
// 触发异常的调用
```

### ObsClient 获取方式

```java
ObsClient obsClient = TestTools.getPipelineEnvironment();
assert obsClient != null;
```

### GET 操作的 XML 响应 Mock

GET 测试需要提供符合 SAX Handler 解析格式的 XML：

```java
String validResponseXml = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
        + "<RootNode>"
        + "  <Element>value</Element>"
        + "</RootNode>";

mockServer.when(request()
                .withMethod("GET")
                .withPath("")
                .withQueryStringParameter("x-obs-<feature>"))
        .respond(response()
                .withStatusCode(200)
                .withHeader("Content-Type", "application/xml;charset=utf-8")
                .withBody(validResponseXml));
```

### PUT 操作的 Mock

PUT 测试仅需模拟成功响应：

```java
mockServer.when(request().withMethod("PUT").withPath(""))
        .respond(response().withStatusCode(200));
```

### DELETE 操作的 Mock

DELETE 测试模拟 204 响应：

```java
mockServer.when(request().withMethod("DELETE").withPath(""))
        .respond(response().withStatusCode(204));
```

---

## @AIGenerated 注解规范

AI 生成的测试方法必须添加 `@AIGenerated` 注解：

```java
@Test
@com.obs.aitool.AIGenerated(
    author = System.getProperty("user.name"),
    date = java.time.LocalDate.now().toString(),
    description = "测试 XXX 功能的正常和异常场景"
)
public void should_succeed_when_xxx_with_valid_parameters() {
    // 测试逻辑
}
```

**关键规则**：
- `author`: 使用 `System.getProperty("user.name")` 动态获取
- `date`: 使用 `java.time.LocalDate.now().toString()` 动态获取
- `description`: 简要描述测试目的
- 禁止硬编码日期和作者信息

---

## 测试覆盖要求

### 必须覆盖的场景

| 操作 | 测试场景 | 预期结果 |
|------|----------|----------|
| GET | request 为 null | IllegalArgumentException |
| GET | bucketName 为 null | IllegalArgumentException |
| GET | 正常请求 | 返回正确解析的 Result |
| PUT | request 为 null | IllegalArgumentException |
| PUT | bucketName 为 null | IllegalArgumentException |
| PUT | 正常请求 | 返回 HeaderResponse |
| DELETE | request 为 null | IllegalArgumentException |
| DELETE | bucketName 为 null | IllegalArgumentException |
| DELETE | 正常请求 | 返回 HeaderResponse |

### PUT 特有场景（如有 Configuration 校验）

| 场景 | 预期结果 |
|------|----------|
| Configuration 为 null | ObsException / IllegalArgumentException |
| Configuration 字段不合法 | ObsException |

### 分支覆盖率

- 目标：**100%** 分支覆盖率
- 所有 if/else 分支必须有对应测试用例
- 所有异常路径必须有测试覆盖

---

## 二、集成测试

连接真实 OBS 环境，验证端到端功能链路。基于 `ObsTrashIT.java` 提取的模式。

### 集成测试类模板

```java
/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.integrated_test.<feature>;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.<feature>.<Feature>Configuration;
import com.obs.services.model.<feature>.Delete<Feature>Request;
import com.obs.services.model.<feature>.Get<Feature>Request;
import com.obs.services.model.<feature>.Get<Feature>Result;
import com.obs.services.model.<feature>.Set<Feature>Request;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.Locale;

public class <Feature>IT {

    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    // IT-001: 设置配置 → 查询验证 → 边界值测试
    @Test
    public void test_SDK_<module>_<feature>_001() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);

        // 如需特殊桶类型（如 POSIX 桶）
        String bucketNamePosix = bucketName + "-posix";

        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        // 创建测试桶（如需特殊类型）
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNamePosix);
        createBucketRequest.setBucketType(BucketTypeEnum.PFS);

        try {
            // Step 1: 创建桶
            ObsBucket obsBucket = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, obsBucket.getStatusCode());

            // Step 2: 设置配置
            <Feature>Configuration config = new <Feature>Configuration(/* 正常参数 */);
            Set<Feature>Request setRequest = new Set<Feature>Request(bucketNamePosix, config);
            HeaderResponse headerResponse = obsClient.set<Feature>(setRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // Step 3: 查询配置并验证
            Get<Feature>Request getRequest = new Get<Feature>Request();
            getRequest.setBucketName(bucketNamePosix);
            Get<Feature>Result getResult = obsClient.get<Feature>(getRequest);
            // 断言配置值
            assertEquals(/* expected */, /* actual */);

            // Step 4: 边界值测试（验证取值范围约束）
            try {
                config.setField(/* 非法值 */);
                obsClient.set<Feature>(setRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertTrue(400 <= e.getResponseCode());
            }

        } finally {
            try {
                TestTools.delete_bucket(obsClient, bucketNamePosix);
            } catch (Throwable ignore) {
            }
        }
    }

    // IT-002: 设置 → 查询 → 删除 → 再查询(404)
    @Test
    public void test_SDK_<module>_<feature>_002() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNamePosix = bucketName + "-posix";

        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        try {
            // Step 1: 创建桶
            CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNamePosix);
            createBucketRequest.setBucketType(BucketTypeEnum.PFS);
            obsClient.createBucket(createBucketRequest);

            // Step 2: 设置配置
            <Feature>Configuration config = new <Feature>Configuration(/* params */);
            Set<Feature>Request setRequest = new Set<Feature>Request(bucketNamePosix, config);
            HeaderResponse headerResponse = obsClient.set<Feature>(setRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // Step 3: 查询验证
            Get<Feature>Request getRequest = new Get<Feature>Request();
            getRequest.setBucketName(bucketNamePosix);
            Get<Feature>Result getResult = obsClient.get<Feature>(getRequest);
            Assert.assertNotNull(getResult.get<Feature>Configuration());

            // Step 4: 删除配置
            Delete<Feature>Request deleteRequest = new Delete<Feature>Request();
            deleteRequest.setBucketName(bucketNamePosix);
            headerResponse = obsClient.delete<Feature>(deleteRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // Step 5: 再次查询应返回 404
            try {
                obsClient.get<Feature>(getRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertEquals(404, e.getResponseCode());
            }

        } finally {
            try {
                TestTools.delete_bucket(obsClient, bucketNamePosix);
            } catch (Throwable ignore) {
            }
        }
    }

    // IT-003: 不支持的桶类型（如适用）
    @Test
    public void test_SDK_<module>_<feature>_003() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);

        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        <Feature>Configuration config = new <Feature>Configuration(/* params */);
        Set<Feature>Request setRequest = new Set<Feature>Request(bucketName, config);

        try {
            obsClient.set<Feature>(setRequest);
            fail();
        } catch (ObsException e) {
            Assert.assertTrue(400 <= e.getResponseCode());
        }

        Get<Feature>Request getRequest = new Get<Feature>Request();
        getRequest.setBucketName(bucketName);
        try {
            obsClient.get<Feature>(getRequest);
            fail();
        } catch (ObsException e) {
            Assert.assertTrue(400 <= e.getResponseCode());
        }
    }

    // IT-004: S3/V2 兼容接口（如适用）
    @Test
    public void test_SDK_<module>_<feature>_004() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String bucketNamePosix = bucketName + "-posix";

        // 使用 V2 接口
        ObsClient obsClient = TestTools.getPipelineEnvironment_V2();
        assert obsClient != null;

        try {
            CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketNamePosix);
            createBucketRequest.setBucketType(BucketTypeEnum.PFS);
            obsClient.createBucket(createBucketRequest);

            // 完整 CRUD 验证（同 IT-002）
            // ...

        } finally {
            try {
                TestTools.delete_bucket(obsClient, bucketNamePosix);
            } catch (Throwable ignore) {
            }
        }
    }
}
```

### 集成测试约定

#### 命名规范

- **文件命名**: `<Feature>IT.java`（IT = Integration Test）
- **方法命名**: `test_SDK_<module>_<feature>_<序号>`
- **桶名生成**: `testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)`

#### 桶管理规则

1. **桶名生成**: 使用 `@Rule TestName` 自动生成唯一桶名
2. **桶类型**: 如需特殊桶类型（如 POSIX），使用 `CreateBucketRequest.setBucketType()`
3. **资源清理**: `finally` 块中调用 `TestTools.delete_bucket()` 清理测试桶
4. **清理容错**: 清理操作用 `try/catch(Throwable ignore)` 包裹，不阻断其他清理

#### ObsClient 获取方式

```java
// OBS 接口
ObsClient obsClient = TestTools.getPipelineEnvironment();

// S3/V2 兼容接口
ObsClient obsClient = TestTools.getPipelineEnvironment_V2();
```

#### 测试场景矩阵

| 编号 | 测试内容 | 核心验证点 |
|------|----------|------------|
| IT-001 | SET + GET + 边界值 | 正常设置、查询一致性、取值范围约束 |
| IT-002 | SET + GET + DELETE + GET(404) | 完整 CRUD 生命周期 |
| IT-003 | 不支持的桶类型 | 异常约束（如对象桶不支持回收站） |
| IT-004 | V2 接口兼容 | S3 兼容接口端到端验证（如适用） |

#### 测试隔离

- 每个测试方法独立创建和清理桶
- 不依赖测试执行顺序
- 不共享测试数据
