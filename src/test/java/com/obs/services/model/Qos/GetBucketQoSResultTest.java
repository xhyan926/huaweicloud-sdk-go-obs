package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.*;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;

public class GetBucketQoSResultTest {

    // 测试类自身属性的setter和getter方法
    @Test
    public void testClassProperties() {
        GetBucketQoSResult result = new GetBucketQoSResult();

        // 测试qosGroup属性
        String testQosGroup = "QOS_GROUP_1";
        result.setQosGroup(testQosGroup);
        assertEquals("QoS组名称设置或获取错误", testQosGroup, result.getQosGroup());

        // 测试bucketQosRules属性
        List<QosRule> bucketRules = new ArrayList<>();
        bucketRules.add(createTestQosRule(1));
        bucketRules.add(createTestQosRule(2));
        result.setBucketQosRules(bucketRules);
        assertSame("Bucket规则列表引用不一致", bucketRules, result.getBucketQosRules());
        assertEquals("Bucket规则数量不匹配", 2, result.getBucketQosRules().size());

        // 测试groupQosRules属性
        List<QosRule> groupRules = new ArrayList<>();
        groupRules.add(createTestQosRule(3));
        result.setGroupQosRules(groupRules);
        assertSame("Group规则列表引用不一致", groupRules, result.getGroupQosRules());
        assertEquals("Group规则数量不匹配", 1, result.getGroupQosRules().size());
    }

    // 测试初始默认值
    @Test
    public void testDefaultValues() {
        GetBucketQoSResult result = new GetBucketQoSResult();

        // 自身属性默认值
        assertEquals("qosGroup默认值应为空字符串", "", result.getQosGroup());
        assertNotNull("bucketQosRules应初始化为空列表", result.getBucketQosRules());
        assertTrue("bucketQosRules初始应为空", result.getBucketQosRules().isEmpty());
        assertNotNull("groupQosRules应初始化为空列表", result.getGroupQosRules());
        assertTrue("groupQosRules初始应为空", result.getGroupQosRules().isEmpty());
    }

    // 测试toString()方法
    @Test
    public void testToString() {
        GetBucketQoSResult result = new GetBucketQoSResult();

        // 设置测试数据
        result.setQosGroup("TEST_GROUP");

        List<QosRule> bucketRules = new ArrayList<>();
        bucketRules.add(createTestQosRule(1000));
        result.setBucketQosRules(bucketRules);

        List<QosRule> groupRules = new ArrayList<>();
        groupRules.add(createTestQosRule(2000));
        result.setGroupQosRules(groupRules);

        // 执行测试
        String toString = result.toString();

        // 验证toString包含所有关键信息
        assertNotNull("toString()不应返回null", toString);
        assertTrue("toString应包含qosGroup信息", toString.contains("qosGroup='TEST_GROUP'"));
        assertTrue("toString应包含bucketQosRules", toString.contains("bucketQosRules="));
        assertTrue("toString应包含groupQosRules", toString.contains("groupQosRules="));

        // 验证继承属性在toString中的存在（即使我们不直接测试其getter/setter）
        assertTrue("toString应包含statusCode", toString.contains("statusCode="));
        assertTrue("toString应包含requestId", toString.contains("requestId='"));
    }

    // 辅助方法：创建测试用的QosRule对象
    private QosRule createTestQosRule(int seed) {
        QpsLimitConfiguration qpsConfig = new QpsLimitConfiguration(
                100 + seed,    // get
                200 + seed,    // putPostDelete
                300 + seed,    // list
                600 + seed     // total
        );

        BpsLimitConfiguration bpsConfig = new BpsLimitConfiguration(
                1024 + seed,   // get
                2048 + seed,   // putPost
                3072 + seed    // total
        );

        // 假设NetworkType是一个枚举，这里使用第一个枚举值
        NetworkType networkType = NetworkType.values()[0];

        return new QosRule(networkType, 1000 + seed, qpsConfig, bpsConfig);
    }
}
