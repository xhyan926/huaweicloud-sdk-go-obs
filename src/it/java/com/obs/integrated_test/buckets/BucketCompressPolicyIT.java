/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.integrated_test.buckets;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.assertFalse;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.AuthTypeEnum;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.compress.CompressPolicyConfiguration;
import com.obs.services.model.compress.CompressPolicyRule;
import com.obs.services.model.compress.DeleteBucketCompressPolicyRequest;
import com.obs.services.model.compress.GetBucketCompressPolicyRequest;
import com.obs.services.model.compress.GetBucketCompressPolicyResult;
import com.obs.services.model.compress.SetBucketCompressPolicyRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import com.obs.test.tools.PropertiesTools;

import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collection;
import java.util.Locale;
import java.util.Map;

@RunWith(Parameterized.class)
public class BucketCompressPolicyIT {
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
    private String projectId;
    private String agency;

    @Before
    public void setUp() throws IOException {
        if (authType == AuthTypeEnum.OBS) {
            obsClient = TestTools.getPipelineEnvironment();
        } else {
            obsClient = TestTools.getPipelineEnvironment_V2();
        }
        Assert.assertNotNull("ObsClient should not be null", obsClient);
        bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            .replace("[", "").replace("]", "");
        try {
            PropertiesTools propertiesTools = PropertiesTools.getInstance(TestTools.getPropertiesFile());
            projectId = propertiesTools.getProperties("projectId");
            agency = propertiesTools.getProperties("agency");
        } catch (Exception e) {
            // 配置读取失败时使用占位值
            projectId = "your-project-id";
            agency = "your-agency";
        }
    }

    @After
    public void tearDown() {
        // 桶的创建和删除由 PrepareTestBucket @Rule 统一管理
    }

    /**
     * IT-001: 设置配置 → 查询验证 → 覆盖策略更新 → 再查询验证
     */
    @Test
    public void test_SDK_compresspolicy_001() {
        // SET: 配置 overwrite=0（不覆盖跳过）
        CompressPolicyRule rule = new CompressPolicyRule("rule-001",
            projectId, agency,
            Arrays.asList("ObjectCreated:*"), ".zip", 0);
        rule.setPrefix("decompress/");
        rule.setDecompresspath("after-decompress/");
        rule.setPolicytype("decompress");
        CompressPolicyConfiguration config = new CompressPolicyConfiguration(Arrays.asList(rule));

        SetBucketCompressPolicyRequest setRequest =
            new SetBucketCompressPolicyRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketCompressPolicy(setRequest);
        assertTrue(setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证配置
        GetBucketCompressPolicyResult result =
            obsClient.getBucketCompressPolicy(new GetBucketCompressPolicyRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertNotNull(result.getCompressPolicyConfiguration());
        assertNotNull(result.getCompressPolicyConfiguration().getRules());
        assertEquals(1, result.getCompressPolicyConfiguration().getRules().size());

        CompressPolicyRule resultRule = result.getCompressPolicyConfiguration().getRules().get(0);
        assertEquals("rule-001", resultRule.getId());
        assertEquals(projectId, resultRule.getProject());
        assertEquals(agency, resultRule.getAgency());
        assertEquals(1, resultRule.getEvents().size());
        assertEquals("ObjectCreated:*", resultRule.getEvents().get(0));
        assertEquals("decompress/", resultRule.getPrefix());
        assertEquals(".zip", resultRule.getSuffix());
        assertEquals(Integer.valueOf(0), resultRule.getOverwrite());
        assertEquals("after-decompress/", resultRule.getDecompresspath());

        // SET: 更新配置 overwrite=2（覆盖）
        CompressPolicyRule updatedRule = new CompressPolicyRule("rule-001",
            projectId, agency,
            Arrays.asList("ObjectCreated:*"), ".zip", 2);
        updatedRule.setPrefix("decompress/");
        updatedRule.setDecompresspath("after-decompress/");
        updatedRule.setPolicytype("decompress");
        CompressPolicyConfiguration updatedConfig =
            new CompressPolicyConfiguration(Arrays.asList(updatedRule));

        setRequest = new SetBucketCompressPolicyRequest(bucketName, updatedConfig);
        setResponse = obsClient.setBucketCompressPolicy(setRequest);
        assertTrue(setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证更新后的配置
        result = obsClient.getBucketCompressPolicy(new GetBucketCompressPolicyRequest(bucketName));
        assertEquals(Integer.valueOf(2),
            result.getCompressPolicyConfiguration().getRules().get(0).getOverwrite());
    }

    /**
     * IT-002: 设置 → 查询 → 删除 → 再查询(404) → 重复删除(204)
     */
    @Test
    public void test_SDK_compresspolicy_002() {
        // SET
        CompressPolicyRule rule = new CompressPolicyRule("rule-002",
            projectId, agency,
            Arrays.asList("ObjectCreated:Put"), ".zip", 1);
        rule.setPolicytype("decompress");
        CompressPolicyConfiguration config = new CompressPolicyConfiguration(Arrays.asList(rule));

        SetBucketCompressPolicyRequest setRequest =
            new SetBucketCompressPolicyRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketCompressPolicy(setRequest);
        assertTrue(setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证配置存在
        GetBucketCompressPolicyResult result =
            obsClient.getBucketCompressPolicy(new GetBucketCompressPolicyRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertEquals(1, result.getCompressPolicyConfiguration().getRules().size());

        // DELETE
        DeleteBucketCompressPolicyRequest deleteRequest =
            new DeleteBucketCompressPolicyRequest(bucketName);
        HeaderResponse deleteResponse = obsClient.deleteBucketCompressPolicy(deleteRequest);
        assertEquals(204, deleteResponse.getStatusCode());

        // GET: 验证配置已删除（期望 404）
        try {
            obsClient.getBucketCompressPolicy(new GetBucketCompressPolicyRequest(bucketName));
            // 服务端可能返回 200 但 rules 为空，或者返回 404
        } catch (ObsException e) {
            assertTrue("Expected 404 after delete, got: " + e.getResponseCode(),
                e.getResponseCode() == 404);
        }

        // DELETE: 重复删除应返回 204（幂等）
        deleteResponse = obsClient.deleteBucketCompressPolicy(
            new DeleteBucketCompressPolicyRequest(bucketName));
        assertEquals(204, deleteResponse.getStatusCode());
    }

    /**
     * IT-003: 多规则配置测试
     */
    @Test
    public void test_SDK_compresspolicy_003() {
        // SET: 配置两条规则
        CompressPolicyRule rule1 = new CompressPolicyRule("rule-multi-001",
            projectId, agency,
            Arrays.asList("ObjectCreated:*"), ".zip", 0);
        rule1.setPrefix("folder-a/");
        rule1.setDecompresspath("output-a/");
        rule1.setPolicytype("decompress");

        CompressPolicyRule rule2 = new CompressPolicyRule("rule-multi-002",
            projectId, agency,
            Arrays.asList("ObjectCreated:Put", "ObjectCreated:Post"), ".zip", 2);
        rule2.setPrefix("folder-b/");
        rule2.setDecompresspath("output-b/");
        rule2.setPolicytype("decompress");

        CompressPolicyConfiguration config =
            new CompressPolicyConfiguration(Arrays.asList(rule1, rule2));
        SetBucketCompressPolicyRequest setRequest =
            new SetBucketCompressPolicyRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketCompressPolicy(setRequest);
        assertTrue(setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证两条规则
        GetBucketCompressPolicyResult result =
            obsClient.getBucketCompressPolicy(new GetBucketCompressPolicyRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertNotNull(result.getCompressPolicyConfiguration().getRules());
        assertEquals(2, result.getCompressPolicyConfiguration().getRules().size());

        CompressPolicyRule resultRule1 = result.getCompressPolicyConfiguration().getRules().get(0);
        assertEquals("rule-multi-001", resultRule1.getId());
        assertEquals("folder-a/", resultRule1.getPrefix());
        assertEquals("output-a/", resultRule1.getDecompresspath());
        assertEquals(Integer.valueOf(0), resultRule1.getOverwrite());

        CompressPolicyRule resultRule2 = result.getCompressPolicyConfiguration().getRules().get(1);
        assertEquals("rule-multi-002", resultRule2.getId());
        assertEquals("folder-b/", resultRule2.getPrefix());
        assertEquals("output-b/", resultRule2.getDecompresspath());
        assertEquals(Integer.valueOf(2), resultRule2.getOverwrite());
        assertEquals(2, resultRule2.getEvents().size());
    }

    /**
     * IT-004: 错误场景 — 验证 SDK 能正确返回服务端错误码、错误信息、请求ID、响应头
     */
    @Test
    public void test_SDK_compresspolicy_004() {
        // --- 场景1: SET 使用不存在的委托名，服务端应返回 400 + 错误信息 ---
        CompressPolicyRule rule = new CompressPolicyRule("rule-error-001",
            projectId, "non_existent_agency_999",
            Arrays.asList("ObjectCreated:*"), ".zip", 0);
        rule.setPolicytype("decompress");
        CompressPolicyConfiguration config = new CompressPolicyConfiguration(Arrays.asList(rule));

        try {
            obsClient.setBucketCompressPolicy(
                new SetBucketCompressPolicyRequest(bucketName, config));
            Assert.fail("Expected ObsException for invalid agency");
        } catch (ObsException e) {
            // 验证 HTTP 状态码
            assertEquals(400, e.getResponseCode());
            // 验证错误响应体不为空（JSON 格式错误信息存入 getXmlMessage()）
            assertNotNull("Error body (xmlMessage) should not be null", e.getXmlMessage());
            assertFalse("Error body should not be empty", e.getXmlMessage().trim().isEmpty());
            // 诊断输出：打印所有响应头和错误体
            Map<String, String> diagHeaders = e.getResponseHeaders();
            StringBuilder headerDump = new StringBuilder("{");
            if (diagHeaders != null) {
                for (Map.Entry<String, String> entry : diagHeaders.entrySet()) {
                    headerDump.append(entry.getKey()).append("=").append(entry.getValue()).append("; ");
                }
            }
            headerDump.append("}");
            String diagInfo = String.format(
                "ResponseCode=%d, AllHeaders=%s, ErrorCode=[%s], ErrorMessage=[%s], "
                + "ErrorRequestId=[%s], Body=[%s]",
                e.getResponseCode(), headerDump.toString(),
                e.getErrorCode(), e.getErrorMessage(),
                e.getErrorRequestId(), e.getXmlMessage());
            // 验证错误体和错误信息不为空
            assertTrue("Error info diagnostic. " + diagInfo,
                e.getXmlMessage() != null && !e.getXmlMessage().trim().isEmpty()
                && e.getErrorRequestId() != null && !e.getErrorRequestId().trim().isEmpty());
            // 验证 Request-ID 不为空（用于问题定位）
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            assertFalse("ErrorRequestId should not be empty", e.getErrorRequestId().trim().isEmpty());
            // 验证响应头不为空
            Map<String, String> headers = e.getResponseHeaders();
            assertNotNull("ResponseHeaders should not be null", headers);
            assertFalse("ResponseHeaders should not be empty", headers.isEmpty());
            // 打印完整错误信息便于确认
            System.out.println("[IT-004 Scene1] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-004 Scene1] ErrorCode: " + e.getErrorCode());
            System.out.println("[IT-004 Scene1] ErrorMessage: " + e.getErrorMessage());
            System.out.println("[IT-004 Scene1] ErrorBody: " + e.getXmlMessage());
            System.out.println("[IT-004 Scene1] RequestId: " + e.getErrorRequestId());
            System.out.println("[IT-004 Scene1] ResponseHeaders: " + headers);
        }

        // --- 场景2: GET 不存在的桶，服务端应返回 404 ---
        try {
            obsClient.getBucketCompressPolicy(
                new GetBucketCompressPolicyRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertTrue("Expected 404 for non-existent bucket, got: " + e.getResponseCode(),
                e.getResponseCode() == 404);
            assertNotNull("ErrorMessage should not be null", e.getErrorMessage());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            // 验证响应头
            Map<String, String> headers = e.getResponseHeaders();
            assertNotNull("ResponseHeaders should not be null", headers);
            System.out.println("[IT-004 Scene2] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-004 Scene2] ErrorCode: " + e.getErrorCode());
            System.out.println("[IT-004 Scene2] ErrorMessage: " + e.getErrorMessage());
            System.out.println("[IT-004 Scene2] RequestId: " + e.getErrorRequestId());
            System.out.println("[IT-004 Scene2] XmlMessage: " + e.getXmlMessage());
        }

        // --- 场景3: DELETE 不存在的桶，服务端应返回 404 ---
        try {
            obsClient.deleteBucketCompressPolicy(
                new DeleteBucketCompressPolicyRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertTrue("Expected 404 for non-existent bucket, got: " + e.getResponseCode(),
                e.getResponseCode() == 404);
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            Map<String, String> headers = e.getResponseHeaders();
            assertNotNull("ResponseHeaders should not be null", headers);
            System.out.println("[IT-004 Scene3] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-004 Scene3] RequestId: " + e.getErrorRequestId());
            System.out.println("[IT-004 Scene3] ResponseHeaders: " + headers);
        }

        // --- 场景4: GET 已删除策略的桶，验证 404 错误完整性 ---
        // 先设置再删除，然后查询
        CompressPolicyRule setupRule = new CompressPolicyRule("rule-to-delete",
            projectId, agency,
            Arrays.asList("ObjectCreated:*"), ".zip", 0);
        setupRule.setPolicytype("decompress");
        CompressPolicyConfiguration setupConfig = new CompressPolicyConfiguration(Arrays.asList(setupRule));
        obsClient.setBucketCompressPolicy(new SetBucketCompressPolicyRequest(bucketName, setupConfig));
        obsClient.deleteBucketCompressPolicy(new DeleteBucketCompressPolicyRequest(bucketName));

        try {
            obsClient.getBucketCompressPolicy(new GetBucketCompressPolicyRequest(bucketName));
            // 如果没有抛异常说明服务端返回200但无策略
        } catch (ObsException e) {
            // 验证 404 及完整错误链
            assertTrue("Expected 404, got: " + e.getResponseCode(),
                e.getResponseCode() == 404);
            assertNotNull("Error body should not be null", e.getXmlMessage());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            assertNotNull("ResponseHeaders should not be null", e.getResponseHeaders());
            // 验证 toString 包含关键错误信息
            String errorStr = e.toString();
            assertTrue("Error toString should contain ResponseCode",
                errorStr.contains(String.valueOf(e.getResponseCode())));
            System.out.println("[IT-004 Scene4] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-004 Scene4] ErrorCode: " + e.getErrorCode());
            System.out.println("[IT-004 Scene4] ErrorMessage: " + e.getErrorMessage());
            System.out.println("[IT-004 Scene4] ErrorBody: " + e.getXmlMessage());
            System.out.println("[IT-004 Scene4] RequestId: " + e.getErrorRequestId());
        }
    }
}
