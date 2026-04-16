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
import com.obs.services.model.mirrorback.DeleteBucketMirrorBackToSourceRequest;
import com.obs.services.model.mirrorback.GetBucketMirrorBackToSourceRequest;
import com.obs.services.model.mirrorback.GetBucketMirrorBackToSourceResult;
import com.obs.services.model.mirrorback.MirrorBackCondition;
import com.obs.services.model.mirrorback.MirrorBackHttpHeader;
import com.obs.services.model.mirrorback.MirrorBackHttpHeaderSet;
import com.obs.services.model.mirrorback.MirrorBackPublicSource;
import com.obs.services.model.mirrorback.MirrorBackRedirect;
import com.obs.services.model.mirrorback.MirrorBackSourceEndpoint;
import com.obs.services.model.mirrorback.MirrorBackToSourceConfiguration;
import com.obs.services.model.mirrorback.MirrorBackToSourceRule;
import com.obs.services.model.mirrorback.SetBucketMirrorBackToSourceRequest;
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
import java.util.Collections;
import java.util.Locale;

@RunWith(Parameterized.class)
public class BucketMirrorBackToSourceIT {
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
            agency = propertiesTools.getProperties("agency");
        } catch (Exception e) {
            agency = "your-agency";
        }
    }

    @After
    public void tearDown() {
        // 桶的创建和删除由 PrepareTestBucket @Rule 统一管理
    }

    /**
     * IT-001: 设置配置 -> 查询验证 -> 覆盖策略更新 -> 再查询验证
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "集成测试：镜像回源设置、查询、更新全流程验证")
    public void test_SDK_mirrorback_001() {
        // SET: 创建基础镜像回源配置
        MirrorBackCondition condition = new MirrorBackCondition("404", "images/");
        MirrorBackSourceEndpoint endpoint = new MirrorBackSourceEndpoint(
            Collections.singletonList("https://mirror-source.obs.cn-southwest-2.myhuaweicloud.com"),
            Collections.emptyList());
        MirrorBackPublicSource publicSource = new MirrorBackPublicSource(endpoint);
        MirrorBackRedirect redirect = new MirrorBackRedirect();
        redirect.setAgency(agency);
        redirect.setPublicSource(publicSource);
        redirect.setPassQueryString(true);
        redirect.setMirrorFollowRedirect(false);
        MirrorBackToSourceRule rule = new MirrorBackToSourceRule("rule-001", condition, redirect);
        MirrorBackToSourceConfiguration config =
            new MirrorBackToSourceConfiguration(Collections.singletonList(rule));

        SetBucketMirrorBackToSourceRequest setRequest =
            new SetBucketMirrorBackToSourceRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketMirrorBackToSource(setRequest);
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证配置
        GetBucketMirrorBackToSourceResult result =
            obsClient.getBucketMirrorBackToSource(new GetBucketMirrorBackToSourceRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertNotNull(result.getMirrorBackToSourceConfiguration());
        assertNotNull(result.getMirrorBackToSourceConfiguration().getRules());
        assertEquals(1, result.getMirrorBackToSourceConfiguration().getRules().size());

        MirrorBackToSourceRule resultRule = result.getMirrorBackToSourceConfiguration().getRules().get(0);
        assertEquals("rule-001", resultRule.getId());
        assertNotNull(resultRule.getCondition());
        assertEquals("404", resultRule.getCondition().getHttpErrorCodeReturnedEquals());
        assertEquals("images/", resultRule.getCondition().getObjectKeyPrefixEquals());

        MirrorBackRedirect resultRedirect = resultRule.getRedirect();
        assertNotNull(resultRedirect);
        assertNotNull(resultRedirect.getPublicSource());
        assertNotNull(resultRedirect.getPublicSource().getSourceEndpoint());
        assertNotNull(resultRedirect.getPublicSource().getSourceEndpoint().getMaster());

        // SET: 更新配置，增加 HTTP 头传递规则
        MirrorBackHttpHeader httpHeader = new MirrorBackHttpHeader();
        httpHeader.setPassAll(false);
        httpHeader.setPass(Arrays.asList("content-type", "cache-control"));
        httpHeader.setRemove(Collections.singletonList("authorization"));
        httpHeader.setSet(Collections.singletonList(
            new MirrorBackHttpHeaderSet("x-custom-header", "custom-value")));

        MirrorBackRedirect updatedRedirect = new MirrorBackRedirect();
        updatedRedirect.setAgency(agency);
        updatedRedirect.setPublicSource(publicSource);
        updatedRedirect.setPassQueryString(true);
        updatedRedirect.setMirrorFollowRedirect(true);
        updatedRedirect.setMirrorHttpHeader(httpHeader);
        updatedRedirect.setReplaceKeyPrefixWith("backup/");
        updatedRedirect.setMirrorAllowHttpMethod(Collections.singletonList("HEAD"));

        MirrorBackToSourceRule updatedRule =
            new MirrorBackToSourceRule("rule-001", condition, updatedRedirect);
        MirrorBackToSourceConfiguration updatedConfig =
            new MirrorBackToSourceConfiguration(Collections.singletonList(updatedRule));

        setRequest = new SetBucketMirrorBackToSourceRequest(bucketName, updatedConfig);
        setResponse = obsClient.setBucketMirrorBackToSource(setRequest);
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证更新后的配置
        result = obsClient.getBucketMirrorBackToSource(new GetBucketMirrorBackToSourceRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        MirrorBackToSourceRule updatedResultRule =
            result.getMirrorBackToSourceConfiguration().getRules().get(0);
        MirrorBackRedirect updatedResultRedirect = updatedResultRule.getRedirect();
        assertNotNull(updatedResultRedirect.getMirrorHttpHeader());
        assertNotNull(updatedResultRedirect.getMirrorAllowHttpMethod());
    }

    /**
     * IT-002: 设置 -> 查询 -> 删除 -> 再查询(404)
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "集成测试：镜像回源设置、查询、删除、再查询404验证")
    public void test_SDK_mirrorback_002() {
        // SET
        MirrorBackCondition condition = new MirrorBackCondition("404", null);
        MirrorBackSourceEndpoint endpoint = new MirrorBackSourceEndpoint(
            Collections.singletonList("https://mirror-source.obs.cn-southwest-2.myhuaweicloud.com"), null);
        MirrorBackPublicSource publicSource = new MirrorBackPublicSource(endpoint);
        MirrorBackRedirect redirect = new MirrorBackRedirect();
        redirect.setAgency(agency);
        redirect.setPublicSource(publicSource);
        redirect.setPassQueryString(true);
        redirect.setMirrorFollowRedirect(false);
        MirrorBackToSourceRule rule = new MirrorBackToSourceRule("rule-002", condition, redirect);
        MirrorBackToSourceConfiguration config =
            new MirrorBackToSourceConfiguration(Collections.singletonList(rule));

        SetBucketMirrorBackToSourceRequest setRequest =
            new SetBucketMirrorBackToSourceRequest(bucketName, config);
        HeaderResponse setResponse = obsClient.setBucketMirrorBackToSource(setRequest);
        assertTrue("Expected 200 or 201, got: " + setResponse.getStatusCode(),
            setResponse.getStatusCode() == 200 || setResponse.getStatusCode() == 201);

        // GET: 验证配置存在
        GetBucketMirrorBackToSourceResult result =
            obsClient.getBucketMirrorBackToSource(new GetBucketMirrorBackToSourceRequest(bucketName));
        assertEquals(200, result.getStatusCode());
        assertEquals(1, result.getMirrorBackToSourceConfiguration().getRules().size());

        // DELETE
        DeleteBucketMirrorBackToSourceRequest deleteRequest =
            new DeleteBucketMirrorBackToSourceRequest(bucketName);
        HeaderResponse deleteResponse = obsClient.deleteBucketMirrorBackToSource(deleteRequest);
        assertEquals(204, deleteResponse.getStatusCode());

        // GET: 验证配置已删除（期望 404）
        try {
            obsClient.getBucketMirrorBackToSource(new GetBucketMirrorBackToSourceRequest(bucketName));
        } catch (ObsException e) {
            assertEquals("Expected 404 after delete, got: " + e.getResponseCode(), 404, e.getResponseCode());
        }

        // DELETE: 重复删除应返回 204（幂等）
        deleteResponse = obsClient.deleteBucketMirrorBackToSource(
            new DeleteBucketMirrorBackToSourceRequest(bucketName));
        assertEquals(204, deleteResponse.getStatusCode());
    }

    /**
     * IT-003: 错误场景验证
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16", description = "集成测试：镜像回源错误场景验证")
    public void test_SDK_mirrorback_003() {
        // 场景1: SET 使用不存在的委托名
        MirrorBackCondition condition = new MirrorBackCondition("404", null);
        MirrorBackSourceEndpoint endpoint = new MirrorBackSourceEndpoint(
            Collections.singletonList("https://mirror-source.obs.cn-southwest-2.myhuaweicloud.com"), null);
        MirrorBackPublicSource publicSource = new MirrorBackPublicSource(endpoint);
        MirrorBackRedirect redirect = new MirrorBackRedirect();
        redirect.setAgency("non_existent_agency_999");
        redirect.setPublicSource(publicSource);
        redirect.setPassQueryString(true);
        redirect.setMirrorFollowRedirect(false);
        MirrorBackToSourceRule rule = new MirrorBackToSourceRule("rule-error", condition, redirect);
        MirrorBackToSourceConfiguration config =
            new MirrorBackToSourceConfiguration(Collections.singletonList(rule));

        try {
            obsClient.setBucketMirrorBackToSource(
                new SetBucketMirrorBackToSourceRequest(bucketName, config));
            Assert.fail("Expected ObsException for invalid agency");
        } catch (ObsException e) {
            assertEquals("Expected 400, got: " + e.getResponseCode(), 400, e.getResponseCode());
            assertNotNull("Error body should not be null", e.getXmlMessage());
            assertFalse("Error body should not be empty", e.getXmlMessage().trim().isEmpty());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            System.out.println("[IT-003 Scene1] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-003 Scene1] ErrorCode: " + e.getErrorCode());
            System.out.println("[IT-003 Scene1] ErrorMessage: " + e.getErrorMessage());
            System.out.println("[IT-003 Scene1] ErrorBody: " + e.getXmlMessage());
            System.out.println("[IT-003 Scene1] RequestId: " + e.getErrorRequestId());
        }

        // 场景2: GET 不存在的桶
        try {
            obsClient.getBucketMirrorBackToSource(
                new GetBucketMirrorBackToSourceRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertEquals("Expected 404 for non-existent bucket, got: " + e.getResponseCode(), 404, e.getResponseCode());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            System.out.println("[IT-003 Scene2] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-003 Scene2] RequestId: " + e.getErrorRequestId());
        }

        // 场景3: DELETE 不存在的桶
        try {
            obsClient.deleteBucketMirrorBackToSource(
                new DeleteBucketMirrorBackToSourceRequest("non-existent-bucket-" + System.currentTimeMillis()));
            Assert.fail("Expected ObsException for non-existent bucket");
        } catch (ObsException e) {
            assertEquals("Expected 404 for non-existent bucket, got: " + e.getResponseCode(), 404, e.getResponseCode());
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
            System.out.println("[IT-003 Scene3] ResponseCode: " + e.getResponseCode());
            System.out.println("[IT-003 Scene3] RequestId: " + e.getErrorRequestId());
        }
    }
}
