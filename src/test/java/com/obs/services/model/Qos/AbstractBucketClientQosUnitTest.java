/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.model.Qos;

import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.Qos.BpsLimitConfiguration;
import com.obs.services.model.Qos.DeleteBucketQosRequest;
import com.obs.services.model.Qos.GetBucketQoSRequest;
import com.obs.services.model.Qos.GetBucketQoSResult;
import com.obs.services.model.Qos.QosConfiguration;
import com.obs.services.model.Qos.QosRule;
import com.obs.services.model.Qos.QpsLimitConfiguration;
import com.obs.services.model.Qos.SetBucketQosRequest;
import com.obs.services.model.Qos.NetworkType;
import com.obs.test.TestTools;

import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;

import java.util.ArrayList;
import java.util.List;

public class AbstractBucketClientQosUnitTest {

    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    private static ClientAndServer mockServer;
    public static String bucketNameForTest = "test-bucket-qos";

    // 测试用QoS配置：Bps值已确认符合>512000的限制（此处沿用之前修正的512100，若需匹配真实响应的21474836480可直接修改）
    private static QosConfiguration createValidQosConfig() {
        QpsLimitConfiguration qpsConfig = new QpsLimitConfiguration(1000, 1000, 1000, 3000);
        BpsLimitConfiguration bpsConfig = new BpsLimitConfiguration(512100, 512100, 512100);
        QosRule rule = new QosRule(NetworkType.INTRANET, 1000, qpsConfig, bpsConfig);
        List<QosRule> rules = new ArrayList<>();
        rules.add(rule);

        QosConfiguration config = new QosConfiguration();
        config.setRules(rules);
        return config;
    }

    private static SetBucketQosRequest createValidSetRequest() {
        return new SetBucketQosRequest(bucketNameForTest, createValidQosConfig());
    }

    @BeforeClass
    public static void setMockServer() {
        // 启动 MockServer
        mockServer = ClientAndServer.startClientAndServer();
        System.setProperty(PROXY_HOST_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_PROPERTY_NAME, "" + mockServer.getLocalPort());
        System.setProperty(PROXY_HOST_S_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_S_PROPERTY_NAME, "" + mockServer.getLocalPort());
    }

    @AfterClass
    public static void clearEnv() {
        // 关闭 MockServer
        mockServer.close();
        System.clearProperty(PROXY_HOST_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_PROPERTY_NAME);
        System.clearProperty(PROXY_HOST_S_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_S_PROPERTY_NAME);
    }

    // ------------------------------ setBucketQos 测试 ------------------------------
    @Test
    public void should_throw_exception_when_setBucketQosRequest_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("setBucketQosRequest is null");
        obsClient.setBucketQos(null);
    }

    @Test
    public void should_throw_exception_when_setBucketQosRequest_bucketName_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        SetBucketQosRequest request = createValidSetRequest();
        request.setBucketName(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.setBucketQos(request);
    }

    @Test
    public void should_throw_exception_when_setBucketQosRequest_qosConfig_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        SetBucketQosRequest request = new SetBucketQosRequest(bucketNameForTest, null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("QosConfig is null");
        obsClient.setBucketQos(request);
    }

    @Test
    public void should_throw_exception_when_setBucketQosRequest_rules_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        QosConfiguration config = new QosConfiguration();
        config.setRules(null);
        SetBucketQosRequest request = new SetBucketQosRequest(bucketNameForTest, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rules is null");
        obsClient.setBucketQos(request);
    }

    @Test
    public void should_throw_exception_when_setBucketQosRequest_rules_is_empty() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        QosConfiguration config = new QosConfiguration();
        config.setRules(new ArrayList<>());
        SetBucketQosRequest request = new SetBucketQosRequest(bucketNameForTest, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("rules is empty");
        obsClient.setBucketQos(request);
    }

    @Test
    public void should_throw_exception_when_qosRule_networkType_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        QosRule rule = new QosRule(null, 1000,
                new QpsLimitConfiguration(1000, 1000, 1000, 3000),
                new BpsLimitConfiguration(512100, 512100, 512100));
        List<QosRule> rules = new ArrayList<>();
        rules.add(rule);

        QosConfiguration config = new QosConfiguration();
        config.setRules(rules);
        SetBucketQosRequest request = new SetBucketQosRequest(bucketNameForTest, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("networkType is null");
        obsClient.setBucketQos(request);
    }

    @Test
    public void should_throw_exception_when_qosRule_qpsLimit_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        QosRule rule = new QosRule(NetworkType.INTRANET, 1000, null,
                new BpsLimitConfiguration(512100, 512100, 512100));
        List<QosRule> rules = new ArrayList<>();
        rules.add(rule);

        QosConfiguration config = new QosConfiguration();
        config.setRules(rules);
        SetBucketQosRequest request = new SetBucketQosRequest(bucketNameForTest, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("qpsLimit is null");
        obsClient.setBucketQos(request);
    }

    @Test
    public void should_throw_exception_when_qosRule_bpsLimit_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        QosRule rule = new QosRule(NetworkType.INTRANET, 1000,
                new QpsLimitConfiguration(1000, 1000, 1000, 3000), null);
        List<QosRule> rules = new ArrayList<>();
        rules.add(rule);

        QosConfiguration config = new QosConfiguration();
        config.setRules(rules);
        SetBucketQosRequest request = new SetBucketQosRequest(bucketNameForTest, config);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bpsLimit is null");
        obsClient.setBucketQos(request);
    }

    @Test
    public void should_succeed_when_setBucketQos_with_valid_parameters() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request().withMethod("PUT").withPath("" ))
                .respond(response().withStatusCode(responseCodeForTest));

        SetBucketQosRequest request = createValidSetRequest();
        HeaderResponse response = obsClient.setBucketQos(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }

    // ------------------------------ getBucketQoS 测试 ------------------------------
    @Test
    public void should_throw_exception_when_getBucketQoSRequest_bucketName_is_null() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        GetBucketQoSRequest request = new GetBucketQoSRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.getBucketQoS(request);
    }

    @Test
    public void should_succeed_when_getBucketQoS_with_valid_parameters() {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        // 1. 准备符合真实格式的XML响应体：
        // - 新增命名空间 xmlns="http://obs.myhwclouds.com/doc/2015-06-30/"
        // - 新增 QoSGroup 节点（值为LZ1）
        // - 新增 QoSGroupConfiguration 节点（含独立QoSRule）
        // - 所有数值与真实response保持一致（NetworkType=total、Bps=21474836480等）
        String validQosResponseXml = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
                + "<GetBucketQoSResponse xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">"
                + "  <QoSConfiguration>"
                + "    <QoSRule>"
                + "      <NetworkType>total</NetworkType>"
                + "      <ConcurrentRequestLimit>2000</ConcurrentRequestLimit>"
                + "      <QpsLimit>"
                + "        <Get>15000</Get>"
                + "        <PutPostDelete>15000</PutPostDelete>"
                + "        <List>15000</List>"
                + "        <Total>0</Total>"
                + "      </QpsLimit>"
                + "      <BpsLimit>"
                + "        <Get>21474836480</Get>"
                + "        <PutPost>21474836480</PutPost>"
                + "        <Total>0</Total>"
                + "      </BpsLimit>"
                + "    </QoSRule>"
                + "  </QoSConfiguration>"
                + "  <QoSGroup>LZ1</QoSGroup>"
                + "  <QoSGroupConfiguration>"
                + "    <QoSRule>"
                + "      <NetworkType>total</NetworkType>"
                + "      <ConcurrentRequestLimit>2000</ConcurrentRequestLimit>"
                + "      <QpsLimit>"
                + "        <Get>15000</Get>"
                + "        <PutPostDelete>15000</PutPostDelete>"
                + "        <List>1000</List>"
                + "        <Total>0</Total>"
                + "      </QpsLimit>"
                + "      <BpsLimit>"
                + "        <Get>21474836480</Get>"
                + "        <PutPost>21474836480</PutPost>"
                + "        <Total>0</Total>"
                + "      </BpsLimit>"
                + "    </QoSRule>"
                + "  </QoSGroupConfiguration>"
                + "</GetBucketQoSResponse>";

        // 2. 配置MockServer：匹配真实请求的Method、Path、Query参数，返回上述XML
        Integer responseCodeForTest = 200;
        mockServer.reset();
        mockServer.when(request()
                        .withMethod("GET")
                        .withPath("")
                        .withQueryStringParameter("x-obs-qosInfo"))
                .respond(response()
                        .withStatusCode(responseCodeForTest)
                        .withHeader("Content-Type", "application/xml;charset=utf-8")
                        .withHeader("Content-Length", String.valueOf(validQosResponseXml.getBytes().length))  
                        .withBody(validQosResponseXml));

        // 3. 执行测试：调用getBucketQoS获取结果
        GetBucketQoSRequest request = new GetBucketQoSRequest(bucketNameForTest);
        GetBucketQoSResult result = obsClient.getBucketQoS(request);

        // 4. 增强断言：全面校验返回结果与真实响应的一致性
        // 4.1 基础非空校验
        Assert.assertNotNull("GetBucketQoSResult 不应为null", result);
        Assert.assertNotNull("QoSGroup 不应为null（真实响应含该节点）", result.getQosGroup());
        Assert.assertNotNull("QoSGroupConfiguration规则 不应为null（真实响应含该节点）", result.getGroupQosRules());
        Assert.assertNotNull("基础QoSConfiguration规则 不应为null", result.getBucketQosRules());

        // 4.2 校验QoSGroup（新增节点，值为LZ1）
        Assert.assertEquals("QoSGroup值与真实响应不一致", "LZ1", result.getQosGroup());

        // 4.3 校验【基础QoSConfiguration】中的QoSRule（对应<QoSConfiguration>节点下的规则）
        QosRule basicQosRule = result.getBucketQosRules().get(0);
        Assert.assertEquals("基础规则的NetworkType与真实响应不一致", NetworkType.TOTAL, basicQosRule.getNetworkType());  // 注意：需确保NetworkType枚举有"TOTAL"值（若枚举值为大写，需匹配）
        Assert.assertEquals("基础规则的ConcurrentRequestLimit与真实响应不一致", 2000L, basicQosRule.getConcurrentRequestLimit());
        // 校验基础规则的QPS限制
        Assert.assertEquals("基础规则QPS-Get与真实响应不一致", 15000L, basicQosRule.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals("基础规则QPS-PutPostDelete与真实响应不一致", 15000L, basicQosRule.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals("基础规则QPS-List与真实响应不一致", 15000L, basicQosRule.getQpsLimit().getQpsListLimit());
        Assert.assertEquals("基础规则QPS-Total与真实响应不一致", 0L, basicQosRule.getQpsLimit().getQpsTotalLimit());
        // 校验基础规则的BPS限制（符合>512000的要求，真实值为21474836480）
        Assert.assertEquals("基础规则BPS-Get与真实响应不一致", 21474836480L, basicQosRule.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals("基础规则BPS-PutPost与真实响应不一致", 21474836480L, basicQosRule.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals("基础规则BPS-Total与真实响应不一致", 0L, basicQosRule.getBpsLimit().getBpsTotalLimit());

        // 4.4 校验【QoSGroupConfiguration】中的QoSRule（对应<QoSGroupConfiguration>节点下的规则）
        QosRule groupQosRule = result.getGroupQosRules().get(0);
        Assert.assertEquals("Group规则的NetworkType与真实响应不一致", NetworkType.TOTAL, groupQosRule.getNetworkType());
        Assert.assertEquals("Group规则的ConcurrentRequestLimit与真实响应不一致", 2000L, groupQosRule.getConcurrentRequestLimit());
        // 校验Group规则的QPS限制（注意：List值为1000，与基础规则的80不同）
        Assert.assertEquals("Group规则QPS-Get与真实响应不一致", 15000L, groupQosRule.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals("Group规则QPS-PutPostDelete与真实响应不一致", 15000L, groupQosRule.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals("Group规则QPS-List与真实响应不一致", 1000L, groupQosRule.getQpsLimit().getQpsListLimit());
        Assert.assertEquals("Group规则QPS-Total与真实响应不一致", 0L, groupQosRule.getQpsLimit().getQpsTotalLimit());
        // 校验Group规则的BPS限制（与基础规则一致，均为21474836480）
        Assert.assertEquals("Group规则BPS-Get与真实响应不一致", 21474836480L, groupQosRule.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals("Group规则BPS-PutPost与真实响应不一致", 21474836480L, groupQosRule.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals("Group规则BPS-Total与真实响应不一致", 0L, groupQosRule.getBpsLimit().getBpsTotalLimit());
    }


    // ------------------------------ deleteBucketQoS 测试 ------------------------------
    @Test
    public void should_throw_exception_when_deleteBucketQoSRequest_bucketName_is_null() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        DeleteBucketQosRequest request = new DeleteBucketQosRequest(null);

        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        obsClient.deleteBucketQoS(request);
    }

    @Test
    public void should_succeed_when_deleteBucketQoS_with_valid_parameters() throws ObsException {
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;

        Integer responseCodeForTest = 204;
        mockServer.reset();
        mockServer.when(request().withMethod("DELETE").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        DeleteBucketQosRequest request = new DeleteBucketQosRequest(bucketNameForTest);
        HeaderResponse response = obsClient.deleteBucketQoS(request);

        Assert.assertEquals(responseCodeForTest.intValue(), response.getStatusCode());
    }
}