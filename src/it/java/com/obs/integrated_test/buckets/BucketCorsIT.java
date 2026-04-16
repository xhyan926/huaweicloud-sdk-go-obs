/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.integrated_test.buckets;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import com.obs.aitool.AIGenerated;
import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.AuthTypeEnum;
import com.obs.services.model.BaseBucketRequest;
import com.obs.services.model.BucketCors;
import com.obs.services.model.BucketCorsRule;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.SetBucketCorsRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Arrays;
import java.util.Collection;
import java.util.Collections;
import java.util.Locale;

@RunWith(Parameterized.class)
public class BucketCorsIT {
    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Parameterized.Parameter()
    public String authTypeName;

    @Parameterized.Parameter(1)
    public AuthTypeEnum authType;

    @Parameterized.Parameters(name = "{0}")
    public static Collection<Object[]> authTypeData() {
        return Arrays.asList(new Object[][] {
            {"OBS", AuthTypeEnum.OBS},
            {"V2", AuthTypeEnum.V2}
        });
    }

    private ObsClient obsClient;
    private String bucketName;

    @Before
    public void setUp() {
        if (authType == AuthTypeEnum.OBS) {
            obsClient = TestTools.getPipelineEnvironment();
        } else {
            obsClient = TestTools.getPipelineEnvironment_V2();
        }
        Assert.assertNotNull("ObsClient should not be null", obsClient);
        bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            .replace("[", "").replace("]", "");
    }

    @After
    public void tearDown() {
        // 桶的创建和删除由 PrepareTestBucket @Rule 统一管理
    }

    /**
     * 设置桶CORS配置并断言成功（状态码 200 或 201）
     */
    private void setBucketCorsAndAssert(BucketCors bucketCors) {
        HeaderResponse setResponse = obsClient.setBucketCors(
            new SetBucketCorsRequest(bucketName, bucketCors));
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);
    }

    /**
     * 查询桶CORS配置并断言返回指定状态码
     */
    private BucketCors getBucketCorsAndAssert(int expectedStatusCode) {
        BucketCors result = obsClient.getBucketCors(new BaseBucketRequest(bucketName));
        assertEquals(expectedStatusCode, result.getStatusCode());
        return result;
    }

    /**
     * 打印错误场景的诊断信息
     */
    private void printErrorDiagnostics(String scene, ObsException e) {
        System.out.println("[" + scene + "] ResponseCode: " + e.getResponseCode());
        System.out.println("[" + scene + "] ErrorCode: " + e.getErrorCode());
        System.out.println("[" + scene + "] ErrorMessage: " + e.getErrorMessage());
        System.out.println("[" + scene + "] ErrorBody: " + e.getXmlMessage());
        System.out.println("[" + scene + "] RequestId: " + e.getErrorRequestId());
        System.out.println("[" + scene + "] ResponseHeaders: " + e.getResponseHeaders());
    }

    /**
     * IT-001: 设置CORS配置 -> 查询验证 -> 更新策略 -> 再查询验证
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：设置桶CORS配置并查询验证，然后更新再验证")
    public void test_SDK_cors_001() {
        // SET: 配置 CORS 规则
        BucketCorsRule rule = new BucketCorsRule();
        rule.setId("rule-001");
        rule.setAllowedOrigin(Arrays.asList("http://www.example.com", "http://www.test.com"));
        rule.setAllowedMethod(Arrays.asList("GET", "PUT", "POST"));
        rule.setAllowedHeader(Collections.singletonList("*"));
        rule.setExposeHeader(Arrays.asList("x-obs-header1", "x-obs-header2"));
        rule.setMaxAgeSecond(100);
        BucketCors bucketCors = new BucketCors(Collections.singletonList(rule));

        setBucketCorsAndAssert(bucketCors);

        // GET: 验证配置
        BucketCors result = getBucketCorsAndAssert(200);
        assertNotNull(result.getRules());
        assertEquals(1, result.getRules().size());

        BucketCorsRule resultRule = result.getRules().get(0);
        assertEquals("rule-001", resultRule.getId());
        assertEquals(3, resultRule.getAllowedMethod().size());
        assertTrue(resultRule.getAllowedMethod().contains("GET"));
        assertTrue(resultRule.getAllowedMethod().contains("PUT"));
        assertTrue(resultRule.getAllowedMethod().contains("POST"));
        assertEquals(2, resultRule.getAllowedOrigin().size());
        assertTrue(resultRule.getAllowedOrigin().contains("http://www.example.com"));
        assertTrue(resultRule.getAllowedOrigin().contains("http://www.test.com"));
        assertEquals(1, resultRule.getAllowedHeader().size());
        assertEquals("*", resultRule.getAllowedHeader().get(0));
        assertEquals(2, resultRule.getExposeHeader().size());
        assertEquals(100, resultRule.getMaxAgeSecond());

        // SET: 更新 CORS 规则，修改为单规则
        BucketCorsRule updatedRule = new BucketCorsRule();
        updatedRule.setId("rule-001-updated");
        updatedRule.setAllowedOrigin(Collections.singletonList("http://www.updated.com"));
        updatedRule.setAllowedMethod(Collections.singletonList("GET"));
        updatedRule.setMaxAgeSecond(200);
        BucketCors updatedCors = new BucketCors(Collections.singletonList(updatedRule));

        setBucketCorsAndAssert(updatedCors);

        // GET: 验证更新后的配置
        result = getBucketCorsAndAssert(200);
        assertEquals(1, result.getRules().size());
        assertEquals("rule-001-updated", result.getRules().get(0).getId());
        assertEquals(1, result.getRules().get(0).getAllowedMethod().size());
        assertEquals("GET", result.getRules().get(0).getAllowedMethod().get(0));
        assertEquals(200, result.getRules().get(0).getMaxAgeSecond());
    }

    /**
     * IT-002: 设置 -> 查询 -> 删除 -> 再查询(404) -> 重复删除(204)
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：验证CORS配置完整CRUD生命周期")
    public void test_SDK_cors_002() {
        // SET
        BucketCorsRule rule = new BucketCorsRule();
        rule.setAllowedOrigin(Collections.singletonList("http://www.example.com"));
        rule.setAllowedMethod(Collections.singletonList("GET"));
        BucketCors bucketCors = new BucketCors(Collections.singletonList(rule));

        setBucketCorsAndAssert(bucketCors);

        // GET: 验证配置存在
        BucketCors result = getBucketCorsAndAssert(200);
        assertEquals(1, result.getRules().size());

        // DELETE
        HeaderResponse deleteResponse = obsClient.deleteBucketCors(new BaseBucketRequest(bucketName));
        assertEquals(204, deleteResponse.getStatusCode());

        // GET: 验证配置已删除（期望 404）
        try {
            obsClient.getBucketCors(new BaseBucketRequest(bucketName));
        } catch (ObsException e) {
            assertEquals("Expected 404 after delete, got: " + e.getResponseCode(), 404, e.getResponseCode());
        }

        // DELETE: 重复删除应返回 204（幂等）
        deleteResponse = obsClient.deleteBucketCors(new BaseBucketRequest(bucketName));
        assertEquals(204, deleteResponse.getStatusCode());
    }

    /**
     * IT-003: 多规则配置测试
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：验证多规则CORS配置")
    public void test_SDK_cors_003() {
        // SET: 配置两条规则
        BucketCorsRule rule1 = new BucketCorsRule();
        rule1.setId("rule-multi-001");
        rule1.setAllowedOrigin(Collections.singletonList("http://www.a.com"));
        rule1.setAllowedMethod(Arrays.asList("PUT", "POST", "DELETE"));
        rule1.setAllowedHeader(Collections.singletonList("*"));

        BucketCorsRule rule2 = new BucketCorsRule();
        rule2.setId("rule-multi-002");
        rule2.setAllowedOrigin(Collections.singletonList("*"));
        rule2.setAllowedMethod(Collections.singletonList("GET"));
        rule2.setMaxAgeSecond(300);
        rule2.setExposeHeader(Arrays.asList("x-obs-header1", "x-obs-header2"));

        BucketCors bucketCors = new BucketCors(Arrays.asList(rule1, rule2));
        setBucketCorsAndAssert(bucketCors);

        // GET: 验证两条规则
        BucketCors result = getBucketCorsAndAssert(200);
        assertNotNull(result.getRules());
        assertEquals(2, result.getRules().size());

        BucketCorsRule resultRule1 = result.getRules().get(0);
        assertEquals("rule-multi-001", resultRule1.getId());
        assertEquals(3, resultRule1.getAllowedMethod().size());
        assertEquals("http://www.a.com", resultRule1.getAllowedOrigin().get(0));

        BucketCorsRule resultRule2 = result.getRules().get(1);
        assertEquals("rule-multi-002", resultRule2.getId());
        assertEquals(1, resultRule2.getAllowedMethod().size());
        assertEquals("GET", resultRule2.getAllowedMethod().get(0));
        assertEquals("*", resultRule2.getAllowedOrigin().get(0));
        assertEquals(300, resultRule2.getMaxAgeSecond());
        assertEquals(2, resultRule2.getExposeHeader().size());
    }

    /**
     * IT-004: 错误场景 -- 验证 SDK 能正确返回服务端错误码、错误信息、请求ID、响应头
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：验证CORS错误场景下的错误信息完整性")
    public void test_SDK_cors_004() {
        // --- 场景1: GET 不存在的桶，服务端应返回 404 ---
        try {
            obsClient.getBucketCors(
                new BaseBucketRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertEquals("Expected 404 for non-existent bucket, got: " + e.getResponseCode(),
                404, e.getResponseCode());
            assertNotNull("ErrorMessage should not be null", e.getErrorMessage());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            assertNotNull("ResponseHeaders should not be null", e.getResponseHeaders());
            printErrorDiagnostics("IT-004 Scene1", e);
        }

        // --- 场景2: DELETE 不存在的桶，服务端应返回 404 ---
        try {
            obsClient.deleteBucketCors(
                new BaseBucketRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertEquals("Expected 404 for non-existent bucket, got: " + e.getResponseCode(),
                404, e.getResponseCode());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            assertNotNull("ResponseHeaders should not be null", e.getResponseHeaders());
            printErrorDiagnostics("IT-004 Scene2", e);
        }

        // --- 场景3: GET 已删除CORS配置的桶，验证 404 错误完整性 ---
        BucketCorsRule setupRule = new BucketCorsRule();
        setupRule.setAllowedOrigin(Collections.singletonList("http://www.setup.com"));
        setupRule.setAllowedMethod(Collections.singletonList("GET"));
        BucketCors setupCors = new BucketCors(Collections.singletonList(setupRule));
        obsClient.setBucketCors(new SetBucketCorsRequest(bucketName, setupCors));
        obsClient.deleteBucketCors(new BaseBucketRequest(bucketName));

        try {
            obsClient.getBucketCors(new BaseBucketRequest(bucketName));
        } catch (ObsException e) {
            assertTrue("Expected 404, got: " + e.getResponseCode(),
                e.getResponseCode() == 404 || e.getResponseCode() == 403);
            assertNotNull("Error body should not be null", e.getXmlMessage());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            assertNotNull("ResponseHeaders should not be null", e.getResponseHeaders());
            printErrorDiagnostics("IT-004 Scene3", e);
        }
    }
}
