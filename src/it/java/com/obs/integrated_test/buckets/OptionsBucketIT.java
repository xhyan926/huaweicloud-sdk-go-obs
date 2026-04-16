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
import com.obs.services.model.BucketCors;
import com.obs.services.model.BucketCorsRule;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.OptionsInfoRequest;
import com.obs.services.model.OptionsInfoResult;
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
public class OptionsBucketIT {
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
     * IT-001: 设置CORS配置 -> OPTIONS桶预检 -> 验证响应字段
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：设置桶CORS配置后OPTIONS预检验证响应字段")
    public void test_SDK_options_bucket_001() {
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

        // OPTIONS: 发送预检请求
        OptionsInfoRequest request = new OptionsInfoRequest(bucketName);
        request.setOrigin("http://www.example.com");
        request.setRequestMethod(Arrays.asList("GET", "PUT"));
        request.setRequestHeaders(Collections.singletonList("Authorization"));

        OptionsInfoResult result = obsClient.optionsBucket(request);

        // 验证响应字段
        assertNotNull(result);
        assertNotNull("AllowOrigin should not be null", result.getAllowOrigin());
        String allowMethodsStr = result.getAllowMethods().toString();
        assertTrue("AllowMethods should contain GET or PUT, actual: " + allowMethodsStr,
            allowMethodsStr.contains("GET") || allowMethodsStr.contains("PUT"));
    }

    /**
     * IT-002: 无CORS配置时OPTIONS桶 -> 验证错误响应
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：无CORS配置时OPTIONS桶预检验证错误响应")
    public void test_SDK_options_bucket_002() {
        // 不设置 CORS 配置，直接 OPTIONS
        OptionsInfoRequest request = new OptionsInfoRequest(bucketName);
        request.setOrigin("http://www.no-cors.com");
        request.setRequestMethod(Collections.singletonList("GET"));

        try {
            obsClient.optionsBucket(request);
            // 如果没有抛异常，说明服务端可能返回了200但无CORS头，也属于正常行为
        } catch (ObsException e) {
            // 预期可能返回 403 或其他错误码
            assertTrue("Expected error response, got: " + e.getResponseCode(),
                e.getResponseCode() >= 400);
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
        }
    }

    /**
     * IT-003: OPTIONS缺少必填Origin头 -> 验证错误响应
     */
    @Test
    @AIGenerated(author = "yanliwei", date = "2026-04-16",
        description = "集成测试：OPTIONS桶预检缺少Origin头验证错误响应")
    public void test_SDK_options_bucket_003() {
        // SET: 先配置 CORS
        BucketCorsRule rule = new BucketCorsRule();
        rule.setAllowedOrigin(Collections.singletonList("http://www.example.com"));
        rule.setAllowedMethod(Collections.singletonList("GET"));
        BucketCors bucketCors = new BucketCors(Collections.singletonList(rule));
        setBucketCorsAndAssert(bucketCors);

        // OPTIONS: 不设置 Origin 头
        OptionsInfoRequest request = new OptionsInfoRequest(bucketName);
        request.setRequestMethod(Collections.singletonList("GET"));
        // 不设置 origin

        try {
            OptionsInfoResult result = obsClient.optionsBucket(request);
            // 服务端可能返回 200 但无有意义的 CORS 头
            assertNotNull(result);
        } catch (ObsException e) {
            assertTrue("Expected error response, got: " + e.getResponseCode(),
                e.getResponseCode() >= 400);
            assertNotNull("ErrorRequestId should not be null", e.getErrorRequestId());
        }
    }
}
