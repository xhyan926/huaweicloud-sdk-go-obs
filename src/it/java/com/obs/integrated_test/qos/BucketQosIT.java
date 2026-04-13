package com.obs.integrated_test.qos;

import com.obs.services.ObsClient;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.Qos.*;
import com.obs.test.TestTools;
import org.junit.Assert;

import org.junit.Test;
import org.junit.rules.TestName;


import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

public class BucketQosIT {
    @org.junit.Rule
    public TestName testName = new TestName();

    @org.junit.Rule
    public com.obs.test.tools.PrepareQosTestBucket prepareTestBucket = new com.obs.test.tools.PrepareQosTestBucket();

    @Test
    public void tc_alpha_java_js_sdk_QoS_001(){
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        assert obsClient != null;

        // 2. 调用setBucketQoS配置桶total为非零Qos，1条rule
        // 参数: networkType=total, qpsLimit和bpsLimit所有属性设置为1000
        // 修改BPS值：从1000改为512001，确保大于512000
        BpsLimitConfiguration bpsLimitConfiguration2 = new BpsLimitConfiguration(512001, 512001, 512001);
        QpsLimitConfiguration qpsLimitConfiguration2 = new QpsLimitConfiguration(1000, 1000, 1000, 1000);
        QosRule qosRule2 = new QosRule(NetworkType.TOTAL, 1000, qpsLimitConfiguration2, bpsLimitConfiguration2);
        QosConfiguration qosConfiguration2 = new QosConfiguration(qosRule2);
        SetBucketQosRequest setBucketQosRequest2 = new SetBucketQosRequest(bucketName, qosConfiguration2);
        HeaderResponse response2 = obsClient.setBucketQos(setBucketQosRequest2);
        Assert.assertEquals(200, response2.getStatusCode());

        // 3. 调用getBucketQoS，验证配置
        GetBucketQoSRequest getBucketQoSRequest3 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult3 = obsClient.getBucketQoS(getBucketQoSRequest3);
        Assert.assertEquals(200, getBucketQoSResult3.getStatusCode());

        List<QosRule> retrievedRules3 = getBucketQoSResult3.getBucketQosRules();
        Assert.assertEquals(1,retrievedRules3.size());
        QosRule retrievedQosRule3 = retrievedRules3.get(0);

        // 验证 networkType
        Assert.assertEquals(NetworkType.TOTAL, retrievedQosRule3.getNetworkType());
        // 验证 concurrentRequestLimit
        Assert.assertEquals(1000, retrievedQosRule3.getConcurrentRequestLimit());

        // 验证 QPS
        QpsLimitConfiguration retrievedQpsLimit3 = retrievedQosRule3.getQpsLimit();
        Assert.assertEquals(1000, retrievedQpsLimit3.getQpsGetLimit());
        Assert.assertEquals(1000, retrievedQpsLimit3.getQpsPutPostDeleteLimit());
        Assert.assertEquals(1000, retrievedQpsLimit3.getQpsTotalLimit());
        Assert.assertEquals(1000, retrievedQpsLimit3.getQpsListLimit());

        // 验证 BPS
        BpsLimitConfiguration retrievedBpsLimit3 = retrievedQosRule3.getBpsLimit();
        Assert.assertEquals(512001, retrievedBpsLimit3.getBpsGetLimit());
        Assert.assertEquals(512001, retrievedBpsLimit3.getBpsPutPostLimit());
        Assert.assertEquals(512001, retrievedBpsLimit3.getBpsTotalLimit());


        // 4. 调用setBucketQoS配置桶total全零Qos，1条rule
        // 参数: networkType=total, qpsLimit和bpsLimit所有属性设置为0
        // BPS保持为0（全零配置不需要修改）
        BpsLimitConfiguration bpsLimitConfiguration4 = new BpsLimitConfiguration(0, 0, 0);
        QpsLimitConfiguration qpsLimitConfiguration4 = new QpsLimitConfiguration(0, 0, 0, 0);
        QosRule qosRule4 = new QosRule(NetworkType.TOTAL, 0, qpsLimitConfiguration4, bpsLimitConfiguration4);
        QosConfiguration qosConfiguration4 = new QosConfiguration(qosRule4);
        SetBucketQosRequest setBucketQosRequest4 = new SetBucketQosRequest(bucketName, qosConfiguration4);
        HeaderResponse response4 = obsClient.setBucketQos(setBucketQosRequest4);
        Assert.assertEquals(200, response4.getStatusCode());

        // 5. 调用getBucketQoS，验证全零配置
        GetBucketQoSRequest getBucketQoSRequest5 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult5 = obsClient.getBucketQoS(getBucketQoSRequest5);
        Assert.assertEquals(200, getBucketQoSResult5.getStatusCode());

        List<QosRule> retrievedRules5 = getBucketQoSResult5.getBucketQosRules();
        QosRule retrievedQosRule5 = retrievedRules5.get(0);

        Assert.assertEquals(NetworkType.TOTAL, retrievedQosRule5.getNetworkType());
        Assert.assertEquals(0, retrievedQosRule5.getConcurrentRequestLimit());

        QpsLimitConfiguration retrievedQpsLimit5 = retrievedQosRule5.getQpsLimit();
        Assert.assertEquals(0, retrievedQpsLimit5.getQpsGetLimit());
        Assert.assertEquals(0, retrievedQpsLimit5.getQpsPutPostDeleteLimit());
        Assert.assertEquals(0, retrievedQpsLimit5.getQpsTotalLimit());
        Assert.assertEquals(0, retrievedQpsLimit5.getQpsListLimit());

        BpsLimitConfiguration retrievedBpsLimit5 = retrievedQosRule5.getBpsLimit();
        Assert.assertEquals(0, retrievedBpsLimit5.getBpsGetLimit());
        Assert.assertEquals(0, retrievedBpsLimit5.getBpsPutPostLimit());
        Assert.assertEquals(0, retrievedBpsLimit5.getBpsTotalLimit());


        // 6. 调用setBucketQoS配置桶公网+内网非零Qos，2条rule
        // rule1: intranet, 所有属性=2000
        // rule2: extranet, 所有属性=1000
        // 修改BPS值：内网从2000改为512100，外网从1000改为512001，确保大于512000
        BpsLimitConfiguration bpsLimitIntranet = new BpsLimitConfiguration(512100, 512100, 512100);
        QpsLimitConfiguration qpsLimitIntranet = new QpsLimitConfiguration(2000, 2000, 2000, 2000);
        QosRule ruleIntranet = new QosRule(NetworkType.INTRANET, 2000, qpsLimitIntranet, bpsLimitIntranet);

        BpsLimitConfiguration bpsLimitExtranet = new BpsLimitConfiguration(512001, 512001, 512001);
        QpsLimitConfiguration qpsLimitExtranet = new QpsLimitConfiguration(1000, 1000, 1000, 1000);
        QosRule ruleExtranet = new QosRule(NetworkType.EXTRANET, 1000, qpsLimitExtranet, bpsLimitExtranet);

        QosConfiguration qosConfiguration6 = new QosConfiguration(ruleIntranet, ruleExtranet);
        SetBucketQosRequest setBucketQosRequest6 = new SetBucketQosRequest(bucketName, qosConfiguration6);
        HeaderResponse response6 = obsClient.setBucketQos(setBucketQosRequest6);
        Assert.assertEquals(200, response6.getStatusCode());

        // 7. 调用getBucketQoS，验证双规则配置
        GetBucketQoSRequest getBucketQoSRequest7 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult7 = obsClient.getBucketQoS(getBucketQoSRequest7);
        Assert.assertEquals(200, getBucketQoSResult7.getStatusCode());

        List<QosRule> retrievedRules7 = getBucketQoSResult7.getBucketQosRules();
        Assert.assertEquals(2, retrievedRules7.size());

        // 排序不确定，需根据 networkType 判断
        QosRule intranetRule = null;
        QosRule extranetRule = null;
        for (QosRule rule : retrievedRules7) {
            if (rule.getNetworkType() == NetworkType.INTRANET) {
                intranetRule = rule;
            } else if (rule.getNetworkType() == NetworkType.EXTRANET) {
                extranetRule = rule;
            }
        }

        Assert.assertNotNull("Intranet rule should exist", intranetRule);
        Assert.assertNotNull("Extranet rule should exist", extranetRule);

        // 验证内网规则
        Assert.assertEquals(2000, intranetRule.getConcurrentRequestLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsTotalLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsListLimit());
        Assert.assertEquals(512100, intranetRule.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals(512100, intranetRule.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals(512100, intranetRule.getBpsLimit().getBpsTotalLimit());

        // 验证外网规则
        Assert.assertEquals(1000, extranetRule.getConcurrentRequestLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsTotalLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsListLimit());
        Assert.assertEquals(512001, extranetRule.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals(512001, extranetRule.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals(512001, extranetRule.getBpsLimit().getBpsTotalLimit());


        // 8. 调用setBucketQoS配置桶公网+内网全零Qos，2条rule
        // rule1: intranet, 所有属性=0
        // rule2: extranet, 所有属性=0
        // BPS保持为0（全零配置不需要修改）
        BpsLimitConfiguration bpsLimitIntranetZero = new BpsLimitConfiguration(0, 0, 0);
        QpsLimitConfiguration qpsLimitIntranetZero = new QpsLimitConfiguration(0, 0, 0, 0);
        QosRule ruleIntranetZero = new QosRule(NetworkType.INTRANET, 0, qpsLimitIntranetZero, bpsLimitIntranetZero);

        BpsLimitConfiguration bpsLimitExtranetZero = new BpsLimitConfiguration(0, 0, 0);
        QpsLimitConfiguration qpsLimitExtranetZero = new QpsLimitConfiguration(0, 0, 0, 0);
        QosRule ruleExtranetZero = new QosRule(NetworkType.EXTRANET, 0, qpsLimitExtranetZero, bpsLimitExtranetZero);

        QosConfiguration qosConfiguration8 = new QosConfiguration(ruleIntranetZero, ruleExtranetZero);
        SetBucketQosRequest setBucketQosRequest8 = new SetBucketQosRequest(bucketName, qosConfiguration8);
        HeaderResponse response8 = obsClient.setBucketQos(setBucketQosRequest8);
        Assert.assertEquals(200, response8.getStatusCode());

        // 9. 调用getBucketQoS，验证全零双规则配置
        GetBucketQoSRequest getBucketQoSRequest9 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult9 = obsClient.getBucketQoS(getBucketQoSRequest9);
        Assert.assertEquals(200, getBucketQoSResult9.getStatusCode());

        List<QosRule> retrievedRules9 = getBucketQoSResult9.getBucketQosRules();
        Assert.assertEquals(2, retrievedRules9.size());

        QosRule intranetRule9 = null;
        QosRule extranetRule9 = null;
        for (QosRule rule : retrievedRules9) {
            if (rule.getNetworkType() == NetworkType.INTRANET) {
                intranetRule9 = rule;
            } else if (rule.getNetworkType() == NetworkType.EXTRANET) {
                extranetRule9 = rule;
            }
        }

        Assert.assertNotNull("Intranet rule should exist", intranetRule9);
        Assert.assertNotNull("Extranet rule should exist", extranetRule9);

        // 验证内网全零
        Assert.assertEquals(0, intranetRule9.getConcurrentRequestLimit());
        Assert.assertEquals(0, intranetRule9.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals(0, intranetRule9.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals(0, intranetRule9.getQpsLimit().getQpsTotalLimit());
        Assert.assertEquals(0, intranetRule9.getQpsLimit().getQpsListLimit());
        Assert.assertEquals(0, intranetRule9.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals(0, intranetRule9.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals(0, intranetRule9.getBpsLimit().getBpsTotalLimit());

        // 验证外网全零
        Assert.assertEquals(0, extranetRule9.getConcurrentRequestLimit());
        Assert.assertEquals(0, extranetRule9.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals(0, extranetRule9.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals(0, extranetRule9.getQpsLimit().getQpsTotalLimit());
        Assert.assertEquals(0, extranetRule9.getQpsLimit().getQpsListLimit());
        Assert.assertEquals(0, extranetRule9.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals(0, extranetRule9.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals(0, extranetRule9.getBpsLimit().getBpsTotalLimit());
    }

    @Test
    public void tc_alpha_java_js_sdk_QoS_002(){
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        assert obsClient != null;

// 2. 调用setBucketQoS配置桶公网+内网非零Qos，2条rule
// rule1: intranet, 所有属性=2000
// rule2: extranet, 所有属性=1000
// 修改BPS值：内网从2000改为512100，外网从1000改为512001
        BpsLimitConfiguration bpsLimitIntranet = new BpsLimitConfiguration(512100, 512100, 512100);
        QpsLimitConfiguration qpsLimitIntranet = new QpsLimitConfiguration(2000, 2000, 2000, 2000);
        QosRule ruleIntranet = new QosRule(NetworkType.INTRANET, 2000, qpsLimitIntranet, bpsLimitIntranet);

        BpsLimitConfiguration bpsLimitExtranet = new BpsLimitConfiguration(512001, 512001, 512001);
        QpsLimitConfiguration qpsLimitExtranet = new QpsLimitConfiguration(1000, 1000, 1000, 1000);
        QosRule ruleExtranet = new QosRule(NetworkType.EXTRANET, 1000, qpsLimitExtranet, bpsLimitExtranet);

        QosConfiguration qosConfiguration2 = new QosConfiguration(ruleIntranet, ruleExtranet);
        SetBucketQosRequest setBucketQosRequest2 = new SetBucketQosRequest(bucketName, qosConfiguration2);
        HeaderResponse response2 = obsClient.setBucketQos(setBucketQosRequest2);

        Assert.assertEquals(200, response2.getStatusCode());

// ========================================================
// 3. 调用getBucketQoS，验证双规则配置
// ========================================================
        GetBucketQoSRequest getBucketQoSRequest3 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult3 = obsClient.getBucketQoS(getBucketQoSRequest3);
        Assert.assertEquals(200, getBucketQoSResult3.getStatusCode());

        List<QosRule> rules3 = getBucketQoSResult3.getBucketQosRules();
        Assert.assertEquals(2, rules3.size());

        QosRule intranetRule = null;
        QosRule extranetRule = null;
        for (QosRule rule : rules3) {
            if (rule.getNetworkType() == NetworkType.INTRANET) {
                intranetRule = rule;
            } else if (rule.getNetworkType() == NetworkType.EXTRANET) {
                extranetRule = rule;
            }
        }

// 验证内网规则数值
        Assert.assertEquals(2000, intranetRule.getConcurrentRequestLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsTotalLimit());
        Assert.assertEquals(2000, intranetRule.getQpsLimit().getQpsListLimit());
        Assert.assertEquals(512100, intranetRule.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals(512100, intranetRule.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals(512100, intranetRule.getBpsLimit().getBpsTotalLimit());

// 验证外网规则数值
        Assert.assertEquals(1000, extranetRule.getConcurrentRequestLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsGetLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsPutPostDeleteLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsTotalLimit());
        Assert.assertEquals(1000, extranetRule.getQpsLimit().getQpsListLimit());
        Assert.assertEquals(512001, extranetRule.getBpsLimit().getBpsGetLimit());
        Assert.assertEquals(512001, extranetRule.getBpsLimit().getBpsPutPostLimit());
        Assert.assertEquals(512001, extranetRule.getBpsLimit().getBpsTotalLimit());

// ========================================================
// 4. 调用setBucketQoS配置桶total为非零Qos，1条rule
// 参数: networkType=total, 所有属性=1000
// 注意：此操作会覆盖之前的双规则
// 修改BPS值：从1000改为512001
// ========================================================
        BpsLimitConfiguration bpsLimitTotal = new BpsLimitConfiguration(512001, 512001, 512001);
        QpsLimitConfiguration qpsLimitTotal = new QpsLimitConfiguration(1000, 1000, 1000, 1000);
        QosRule totalRule = new QosRule(NetworkType.TOTAL, 1000, qpsLimitTotal, bpsLimitTotal);

        QosConfiguration qosConfiguration4 = new QosConfiguration(totalRule);
        SetBucketQosRequest setBucketQosRequest4 = new SetBucketQosRequest(bucketName, qosConfiguration4);
        HeaderResponse response4 = obsClient.setBucketQos(setBucketQosRequest4);

        Assert.assertEquals(200, response4.getStatusCode());

// ========================================================
// 5. 调用getBucketQoS，验证已被覆盖为 total 单规则
// ========================================================
        GetBucketQoSRequest getBucketQoSRequest5 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult5 = obsClient.getBucketQoS(getBucketQoSRequest5);
        Assert.assertEquals(200, getBucketQoSResult5.getStatusCode());

        List<QosRule> rules5 = getBucketQoSResult5.getBucketQosRules();
        Assert.assertEquals(1, rules5.size());

        QosRule totalRule5 = rules5.get(0);
        Assert.assertEquals(NetworkType.TOTAL, totalRule5.getNetworkType());
        Assert.assertEquals(1000, totalRule5.getConcurrentRequestLimit());

// 验证 QPS
        QpsLimitConfiguration qpsLimit5 = totalRule5.getQpsLimit();
        Assert.assertEquals(1000, qpsLimit5.getQpsGetLimit());
        Assert.assertEquals(1000, qpsLimit5.getQpsPutPostDeleteLimit());
        Assert.assertEquals(1000, qpsLimit5.getQpsTotalLimit());
        Assert.assertEquals(1000, qpsLimit5.getQpsListLimit());

// 验证 BPS
        BpsLimitConfiguration bpsLimit5 = totalRule5.getBpsLimit();
        Assert.assertEquals(512001, bpsLimit5.getBpsGetLimit());
        Assert.assertEquals(512001, bpsLimit5.getBpsPutPostLimit());
        Assert.assertEquals(512001, bpsLimit5.getBpsTotalLimit());
    }

    @Test
    public void tc_alpha_java_js_sdk_QoS_003(){
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        assert obsClient != null;

// ========== 2. 调用setBucketQoS仅配置桶公网Qos ==========
        BpsLimitConfiguration bpsLimitExtranet = new BpsLimitConfiguration(513000, 513000, 513000);
        QpsLimitConfiguration qpsLimitExtranet = new QpsLimitConfiguration(2000, 2000, 2000, 2000);
        QosRule extranetRule = new QosRule(NetworkType.EXTRANET, 2000, qpsLimitExtranet, bpsLimitExtranet);
        QosConfiguration qosConfig2 = new QosConfiguration(extranetRule);
        SetBucketQosRequest request2 = new SetBucketQosRequest(bucketName, qosConfig2);
        try {
            obsClient.setBucketQos(request2);
            Assert.fail("步骤2应失败");
        } catch (Exception e) {
            Assert.assertTrue(e.toString().contains("InvalidArgument") || e.toString().contains("400"));
        }

// ========== 3. 调用setBucketQoS仅配置桶内网Qos ==========
        BpsLimitConfiguration bpsLimitIntranet = new BpsLimitConfiguration(513000, 513000, 513000);
        QpsLimitConfiguration qpsLimitIntranet = new QpsLimitConfiguration(2000, 2000, 2000, 2000);
        QosRule intranetRule = new QosRule(NetworkType.INTRANET, 2000, qpsLimitIntranet, bpsLimitIntranet);
        QosConfiguration qosConfig3 = new QosConfiguration(intranetRule);
        SetBucketQosRequest request3 = new SetBucketQosRequest(bucketName, qosConfig3);
        try {
            obsClient.setBucketQos(request3);
            Assert.fail("步骤3应失败");
        } catch (Exception e) {
            Assert.assertTrue(e.toString().contains("InvalidArgument") || e.toString().contains("400"));
        }

// ========== 4. 配置公网+total Qos（非法） ==========
        BpsLimitConfiguration bpsLimit = new BpsLimitConfiguration(513000, 513000, 513000);
        QpsLimitConfiguration qpsLimit = new QpsLimitConfiguration(2000, 2000, 2000, 2000);
        QosRule ruleExtranet = new QosRule(NetworkType.EXTRANET, 2000, qpsLimit, bpsLimit);
        QosRule ruleTotal = new QosRule(NetworkType.TOTAL, 2000, qpsLimit, bpsLimit);
        QosConfiguration qosConfig4 = new QosConfiguration(ruleExtranet, ruleTotal);
        SetBucketQosRequest request4 = new SetBucketQosRequest(bucketName, qosConfig4);
        try {
            obsClient.setBucketQos(request4);
            Assert.fail("步骤4应失败");
        } catch (Exception e) {
            Assert.assertTrue(e.toString().contains("InvalidArgument") || e.toString().contains("400"));
        }

// ========== 5. 配置内网+total Qos（非法） ==========
        QosConfiguration qosConfig5 = new QosConfiguration(intranetRule, ruleTotal);
        SetBucketQosRequest request5 = new SetBucketQosRequest(bucketName, qosConfig5);
        try {
            obsClient.setBucketQos(request5);
            Assert.fail("步骤5应失败");
        } catch (Exception e) {
            Assert.assertTrue(e.toString().contains("InvalidArgument") || e.toString().contains("400"));
        }

// ========== 6. 配置内网+公网+total Qos（非法） ==========
        List<QosRule> list = new ArrayList<>();
        list.add(intranetRule);
        list.add(extranetRule);
        list.add(ruleTotal);
        QosConfiguration qosConfig6 = new QosConfiguration();
        qosConfig6.setRules(list);
        SetBucketQosRequest request6 = new SetBucketQosRequest(bucketName, qosConfig6);
        try {
            obsClient.setBucketQos(request6);
            Assert.fail("步骤6应失败");
        } catch (Exception e) {
            Assert.assertTrue(e.toString().contains("InvalidArgument") || e.toString().contains("400"));
        }

// ========== 7. 不传QosRule（非法） ==========
        try {
            QosConfiguration qosConfig7 = new QosConfiguration();
            SetBucketQosRequest request7 = new SetBucketQosRequest(bucketName, qosConfig7);
            obsClient.setBucketQos(request7);
            Assert.fail("步骤7应失败");
        } catch (Exception e) {
            Assert.assertTrue( e.toString().contains("rules is null"));
        }
    }

    @Test
    public void tc_alpha_java_js_sdk_QoS_004() {
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        assert obsClient != null;

// 2. 调用getBucketQoS查询该集群租户Qos（已预配置的默认租户QoS）
        GetBucketQoSRequest getBucketQoSRequest2 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult2 = obsClient.getBucketQoS(getBucketQoSRequest2);
        Assert.assertEquals(200, getBucketQoSResult2.getStatusCode());


// 验证返回规则数量为1，表示租户级QoS已正确配置
        List<QosRule> tenantRules = getBucketQoSResult2.getGroupQosRules();
        Assert.assertFalse(tenantRules.isEmpty());
        QosRule tenantRule = tenantRules.get(0);

        long tenantBpsLimit = tenantRule.getBpsLimit().getBpsGetLimit();
        long tenantQpsLimit = tenantRule.getQpsLimit().getQpsGetLimit();

        long exceedBps = tenantBpsLimit + 1;
        long exceedQps = tenantQpsLimit + 1;

// 构造超限的bps和qps配置
        BpsLimitConfiguration bpsExceed = new BpsLimitConfiguration(exceedBps, exceedBps, exceedBps);
        QpsLimitConfiguration qpsExceed = new QpsLimitConfiguration(exceedQps, exceedQps, exceedQps, exceedQps);

// 3. 调用setBucketQoS配置值超过租户Qos（intranet + extranet），2条rule
// rule1: networkType=intranet, 所有属性=租户Qos+1
// rule2: networkType=extranet, 所有属性=租户Qos+1
        QosRule ruleIntranet = new QosRule(NetworkType.INTRANET, exceedBps, qpsExceed, bpsExceed);
        QosRule ruleExtranet = new QosRule(NetworkType.EXTRANET, exceedBps, qpsExceed, bpsExceed);
        QosConfiguration qosConfiguration3 = new QosConfiguration(ruleIntranet, ruleExtranet);
        SetBucketQosRequest setBucketQosRequest3 = new SetBucketQosRequest(bucketName, qosConfiguration3);

        try {
            HeaderResponse response3 = obsClient.setBucketQos(setBucketQosRequest3);
            Assert.fail("步骤3应失败：桶级配置不能超过租户QoS限制");
        } catch (Exception e) {
            String errorMsg = e.toString();
            Assert.assertTrue("Expected InvalidArgument or 400",
                    errorMsg.contains("InvalidArgument") || errorMsg.contains("400"));
        }

// 4. 调用setBucketQoS配置值超过租户Qos（total），1条rule
// networkType=total, 所有属性=租户Qos+1
        QosRule ruleTotal = new QosRule(NetworkType.TOTAL, exceedBps, qpsExceed, bpsExceed);
        QosConfiguration qosConfiguration4 = new QosConfiguration(ruleTotal);
        SetBucketQosRequest setBucketQosRequest4 = new SetBucketQosRequest(bucketName, qosConfiguration4);

        try {
            HeaderResponse response4 = obsClient.setBucketQos(setBucketQosRequest4);
            Assert.fail("步骤4应失败：total配置超出租户QoS限制");
        } catch (Exception e) {
            String errorMsg = e.toString();
            Assert.assertTrue("Expected InvalidArgument or 400",
                    errorMsg.contains("InvalidArgument") || errorMsg.contains("400"));
        }
    }

    @Test
    public void tc_alpha_java_js_sdk_QoS_005(){
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        assert obsClient != null;

// 2. 调用getBucketQoS，桶名为桶A的名称
        GetBucketQoSRequest getBucketQoSRequest2 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult2 = obsClient.getBucketQoS(getBucketQoSRequest2);
        Assert.assertEquals(200, getBucketQoSResult2.getStatusCode());

        List<QosRule> rules2 = getBucketQoSResult2.getBucketQosRules();
        QosRule tenantQosRule = rules2.get(0);

// 3. 调用deleteBucketQoS，桶名为桶A的名称
        DeleteBucketQosRequest deleteBucketQosRequest3 = new DeleteBucketQosRequest(bucketName);
        HeaderResponse response3 = obsClient.deleteBucketQoS(deleteBucketQosRequest3);
        Assert.assertEquals(204, response3.getStatusCode());

// 4. 调用getBucketQoS，桶名为桶A的名称
        GetBucketQoSRequest getBucketQoSRequest4 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult4 = obsClient.getBucketQoS(getBucketQoSRequest4);
        Assert.assertEquals(200, getBucketQoSResult4.getStatusCode());

        List<QosRule> rules4 = getBucketQoSResult4.getBucketQosRules();
        Assert.assertEquals(1, rules4.size());
        QosRule rule4 = rules4.get(0);

// 验证 networkType 为 TOTAL
        Assert.assertEquals(NetworkType.TOTAL, rule4.getNetworkType());

// 验证所有限流值为 0
        Assert.assertEquals(0, rule4.getConcurrentRequestLimit());

        QpsLimitConfiguration qps4 = rule4.getQpsLimit();
        Assert.assertEquals(0, qps4.getQpsGetLimit());
        Assert.assertEquals(0, qps4.getQpsPutPostDeleteLimit());
        Assert.assertEquals(0, qps4.getQpsTotalLimit());
        Assert.assertEquals(0, qps4.getQpsListLimit());

        BpsLimitConfiguration bps4 = rule4.getBpsLimit();
        Assert.assertEquals(0, bps4.getBpsGetLimit());
        Assert.assertEquals(0, bps4.getBpsPutPostLimit());
        Assert.assertEquals(0, bps4.getBpsTotalLimit());

// 5. 调用deleteBucketQoS，重复删除
        DeleteBucketQosRequest deleteBucketQosRequest5 = new DeleteBucketQosRequest(bucketName);
        HeaderResponse response5 = obsClient.deleteBucketQoS(deleteBucketQosRequest5);
        Assert.assertEquals(204, response5.getStatusCode());

// 6. 调用getBucketQoS，桶名为桶A的名称
        GetBucketQoSRequest getBucketQoSRequest6 = new GetBucketQoSRequest(bucketName);
        GetBucketQoSResult getBucketQoSResult6 = obsClient.getBucketQoS(getBucketQoSRequest6);
        Assert.assertEquals(200, getBucketQoSResult6.getStatusCode());

        List<QosRule> rules6 = getBucketQoSResult6.getBucketQosRules();
        QosRule rule6 = rules6.get(0);

// 再次验证所有值仍为 0
        Assert.assertEquals(NetworkType.TOTAL, rule6.getNetworkType());
        Assert.assertEquals(0, rule6.getConcurrentRequestLimit());

        QpsLimitConfiguration qps6 = rule6.getQpsLimit();
        Assert.assertEquals(0, qps6.getQpsGetLimit());
        Assert.assertEquals(0, qps6.getQpsPutPostDeleteLimit());
        Assert.assertEquals(0, qps6.getQpsTotalLimit());
        Assert.assertEquals(0, qps6.getQpsListLimit());

        BpsLimitConfiguration bps6 = rule6.getBpsLimit();
        Assert.assertEquals(0, bps6.getBpsGetLimit());
        Assert.assertEquals(0, bps6.getBpsPutPostLimit());
        Assert.assertEquals(0, bps6.getBpsTotalLimit());
    }
}
