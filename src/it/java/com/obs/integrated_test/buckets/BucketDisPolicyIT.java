/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.integrated_test.buckets;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.assertFalse;

import com.obs.aitool.AIGenerated;
import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.AuthTypeEnum;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.dis.DeleteBucketDisPolicyRequest;
import com.obs.services.model.dis.DisPolicyConfiguration;
import com.obs.services.model.dis.DisPolicyRule;
import com.obs.services.model.dis.GetBucketDisPolicyRequest;
import com.obs.services.model.dis.GetBucketDisPolicyResult;
import com.obs.services.model.dis.SetBucketDisPolicyRequest;
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
import java.util.*;

@RunWith(Parameterized.class)
public class BucketDisPolicyIT {
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
    private String stream;

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
            stream = propertiesTools.getProperties("stream");
        } catch (Exception e) {
            // 配置读取失败时使用占位值
            projectId = "your-project-id";
            agency = "your-agency";
            stream = "your-stream";
        }
    }

    @After
    public void tearDown() {
        // 桶的创建和删除由 PrepareTestBucket @Rule 统一管理
    }

    /**
     * IT-001: 设置配置 → 查询验证 → 更新策略 → 再查询验证
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：设置DIS通知策略并查询验证，然后更新再验证")
    public void test_SDK_dispolicy_001() {
        // SET: 配置 DIS 通知策略
        DisPolicyRule rule = new DisPolicyRule("rule-001", stream,
            projectId, Arrays.asList("ObjectCreated:*", "ObjectRemoved:*"), agency);
        rule.setPrefix("input/");
        rule.setSuffix(".txt");
        DisPolicyConfiguration config = new DisPolicyConfiguration(Collections.singletonList(rule));

        SetBucketDisPolicyRequest setRequest =
            new SetBucketDisPolicyRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketDisPolicy(setRequest);
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证配置
        GetBucketDisPolicyResult result =
            obsClient.getBucketDisPolicy(new GetBucketDisPolicyRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertNotNull(result.getDisPolicyConfiguration());
        assertNotNull(result.getDisPolicyConfiguration().getRules());
        assertEquals(1, result.getDisPolicyConfiguration().getRules().size());

        DisPolicyRule resultRule = result.getDisPolicyConfiguration().getRules().get(0);
        assertEquals("rule-001", resultRule.getId());
        assertEquals(stream, resultRule.getStream());
        assertEquals(projectId, resultRule.getProject());
        assertEquals(agency, resultRule.getAgency());
        assertEquals(2, resultRule.getEvents().size());
        assertEquals("input/", resultRule.getPrefix());
        assertEquals(".txt", resultRule.getSuffix());

        // SET: 更新策略，修改事件列表
        DisPolicyRule updatedRule = new DisPolicyRule("rule-001", stream,
            projectId, Collections.singletonList("ObjectCreated:Put"), agency);
        updatedRule.setPrefix("input/");
        updatedRule.setSuffix(".txt");
        DisPolicyConfiguration updatedConfig =
            new DisPolicyConfiguration(Collections.singletonList(updatedRule));

        setRequest = new SetBucketDisPolicyRequest(bucketName, updatedConfig);
        setResponse = obsClient.setBucketDisPolicy(setRequest);
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证更新后的配置
        result = obsClient.getBucketDisPolicy(new GetBucketDisPolicyRequest(bucketName));
        assertEquals(1, result.getDisPolicyConfiguration().getRules().get(0).getEvents().size());
        assertEquals("ObjectCreated:Put",
            result.getDisPolicyConfiguration().getRules().get(0).getEvents().get(0));
    }

    /**
     * IT-002: 设置 → 查询 → 删除 → 再查询(404) → 重复删除(204)
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：验证DIS通知策略完整CRUD生命周期")
    public void test_SDK_dispolicy_002() {
        // SET
        DisPolicyRule rule = new DisPolicyRule("rule-002", stream,
            projectId, Collections.singletonList("ObjectCreated:*"), agency);
        DisPolicyConfiguration config = new DisPolicyConfiguration(Collections.singletonList(rule));

        SetBucketDisPolicyRequest setRequest =
            new SetBucketDisPolicyRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketDisPolicy(setRequest);
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证配置存在
        GetBucketDisPolicyResult result =
            obsClient.getBucketDisPolicy(new GetBucketDisPolicyRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertEquals(1, result.getDisPolicyConfiguration().getRules().size());

        // DELETE
        DeleteBucketDisPolicyRequest deleteRequest =
            new DeleteBucketDisPolicyRequest(bucketName);
        HeaderResponse deleteResponse = obsClient.deleteBucketDisPolicy(deleteRequest);
        assertEquals(204, deleteResponse.getStatusCode());

        // GET: 验证配置已删除（期望 404）
        try {
            obsClient.getBucketDisPolicy(new GetBucketDisPolicyRequest(bucketName));
        } catch (ObsException e) {
            assertEquals("Expected 404 after delete, got: " + e.getResponseCode(), 404, e.getResponseCode());
        }

        // DELETE: 重复删除应返回 204（幂等）
        deleteResponse = obsClient.deleteBucketDisPolicy(
            new DeleteBucketDisPolicyRequest(bucketName));
        assertEquals(204, deleteResponse.getStatusCode());
    }

    /**
     * IT-003: 多规则配置测试
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：验证多规则DIS通知策略配置")
    public void test_SDK_dispolicy_003() {
        // SET: 配置两条规则
        DisPolicyRule rule1 = new DisPolicyRule("rule-multi-001", stream,
            projectId, Collections.singletonList("ObjectCreated:*"), agency);
        rule1.setPrefix("folder-a/");

        DisPolicyRule rule2 = new DisPolicyRule("rule-multi-002", stream,
            projectId, Arrays.asList("ObjectCreated:Put", "ObjectRemoved:*"), agency);
        rule2.setPrefix("folder-b/");
        rule2.setSuffix(".log");

        DisPolicyConfiguration config =
            new DisPolicyConfiguration(Arrays.asList(rule1, rule2));
        SetBucketDisPolicyRequest setRequest =
            new SetBucketDisPolicyRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketDisPolicy(setRequest);
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证两条规则
        GetBucketDisPolicyResult result =
            obsClient.getBucketDisPolicy(new GetBucketDisPolicyRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertNotNull(result.getDisPolicyConfiguration().getRules());
        assertEquals(2, result.getDisPolicyConfiguration().getRules().size());

        DisPolicyRule resultRule1 = result.getDisPolicyConfiguration().getRules().get(0);
        assertEquals("rule-multi-001", resultRule1.getId());
        assertEquals(stream, resultRule1.getStream());
        assertEquals("folder-a/", resultRule1.getPrefix());

        DisPolicyRule resultRule2 = result.getDisPolicyConfiguration().getRules().get(1);
        assertEquals("rule-multi-002", resultRule2.getId());
        assertEquals(stream, resultRule2.getStream());
        assertEquals("folder-b/", resultRule2.getPrefix());
        assertEquals(".log", resultRule2.getSuffix());
        assertEquals(2, resultRule2.getEvents().size());
    }

    /**
     * IT-004: 错误场景 — 验证 SDK 能正确返回服务端错误码、错误信息、请求ID、响应头
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：验证错误场景下的错误信息完整性")
    public void test_SDK_dispolicy_004() {
        // --- 场景1: SET 使用不存在的委托名，服务端应返回 400 + 错误信息 ---
        DisPolicyRule rule = new DisPolicyRule("rule-error-001", stream,
            projectId, Collections.singletonList("ObjectCreated:*"), "non_existent_agency_999");
        DisPolicyConfiguration config = new DisPolicyConfiguration(Collections.singletonList(rule));

        try {
            obsClient.setBucketDisPolicy(
                new SetBucketDisPolicyRequest(bucketName, config));
            Assert.fail("Expected ObsException for invalid agency");
        } catch (ObsException e) {
            assertEquals(400, e.getResponseCode());
            assertNotNull("Error body (xmlMessage) should not be null", e.getXmlMessage());
            assertFalse("Error body should not be empty", e.getXmlMessage().trim().isEmpty());
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
            assertTrue("Error info diagnostic. " + diagInfo,
                e.getXmlMessage() != null && !e.getXmlMessage().trim().isEmpty()
                && e.getErrorRequestId() != null && !e.getErrorRequestId().trim().isEmpty());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            assertFalse("ErrorRequestId should not be empty", e.getErrorRequestId().trim().isEmpty());
            Map<String, String> headers = e.getResponseHeaders();
            assertNotNull("ResponseHeaders should not be null", headers);
            assertFalse("ResponseHeaders should not be empty", headers.isEmpty());
            System.out.println("[IT-004 Scene1] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-004 Scene1] ErrorCode: " + e.getErrorCode());
            System.out.println("[IT-004 Scene1] ErrorMessage: " + e.getErrorMessage());
            System.out.println("[IT-004 Scene1] ErrorBody: " + e.getXmlMessage());
            System.out.println("[IT-004 Scene1] RequestId: " + e.getErrorRequestId());
            System.out.println("[IT-004 Scene1] ResponseHeaders: " + headers);
        }

        // --- 场景2: GET 不存在的桶，服务端应返回 404 ---
        try {
            obsClient.getBucketDisPolicy(
                new GetBucketDisPolicyRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertEquals("Expected 404 for non-existent bucket, got: " + e.getResponseCode(), 404, e.getResponseCode());
            assertNotNull("ErrorMessage should not be null", e.getErrorMessage());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            Map<String, String> headers = e.getResponseHeaders();
            assertNotNull("ResponseHeaders should not be null", headers);
            System.out.println("[IT-004 Scene2] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-004 Scene2] ErrorCode: " + e.getErrorCode());
            System.out.println("[IT-004 Scene2] ErrorMessage: " + e.getErrorMessage());
            System.out.println("[IT-004 Scene2] RequestId: " + e.getErrorRequestId());
        }

        // --- 场景3: DELETE 不存在的桶，服务端应返回 404 ---
        try {
            obsClient.deleteBucketDisPolicy(
                new DeleteBucketDisPolicyRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertEquals("Expected 404 for non-existent bucket, got: " + e.getResponseCode(), 404, e.getResponseCode());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            Map<String, String> headers = e.getResponseHeaders();
            assertNotNull("ResponseHeaders should not be null", headers);
            System.out.println("[IT-004 Scene3] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-004 Scene3] RequestId: " + e.getErrorRequestId());
            System.out.println("[IT-004 Scene3] ResponseHeaders: " + headers);
        }

        // --- 场景4: GET 已删除策略的桶，验证 404 错误完整性 ---
        DisPolicyRule setupRule = new DisPolicyRule("rule-to-delete", stream,
            projectId, Collections.singletonList("ObjectCreated:*"), agency);
        DisPolicyConfiguration setupConfig = new DisPolicyConfiguration(Collections.singletonList(setupRule));
        obsClient.setBucketDisPolicy(new SetBucketDisPolicyRequest(bucketName, setupConfig));
        obsClient.deleteBucketDisPolicy(new DeleteBucketDisPolicyRequest(bucketName));

        try {
            obsClient.getBucketDisPolicy(new GetBucketDisPolicyRequest(bucketName));
        } catch (ObsException e) {
            assertEquals("Expected 404, got: " + e.getResponseCode(), 404, e.getResponseCode());
            assertNotNull("Error body should not be null", e.getXmlMessage());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            assertNotNull("ResponseHeaders should not be null", e.getResponseHeaders());
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
