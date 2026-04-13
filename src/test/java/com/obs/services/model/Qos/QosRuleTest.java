package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.BpsLimitConfiguration;
import com.obs.services.model.Qos.QosRule;
import com.obs.services.model.Qos.QpsLimitConfiguration;
import org.junit.Rule;
import org.junit.Test;
import com.obs.services.model.Qos.NetworkType;
import org.junit.rules.ExpectedException;

public class QosRuleTest {
    @Rule
    public ExpectedException expectedException = ExpectedException.none();

    // 测试构造函数和getter方法
    @Test
    public void testConstructorAndGetters() {
        // 准备测试数据
        NetworkType networkType = NetworkType.INTRANET;
        long concurrentLimit = 5000;
        QpsLimitConfiguration qpsConfig = new QpsLimitConfiguration(1000, 2000, 3000, 4000);
        BpsLimitConfiguration bpsConfig = new BpsLimitConfiguration(5000, 6000, 7000);

        // 创建测试对象
        QosRule qosRule = new QosRule(networkType, concurrentLimit, qpsConfig, bpsConfig);

        // 验证属性值
        assertEquals("网络类型不匹配", networkType, qosRule.getNetworkType());
        assertEquals("并发请求限制不匹配", concurrentLimit, qosRule.getConcurrentRequestLimit());
        assertSame("QPS配置对象不匹配", qpsConfig, qosRule.getQpsLimit());
        assertSame("BPS配置对象不匹配", bpsConfig, qosRule.getBpsLimit());
    }

    // 测试setter方法
    @Test
    public void testSetters() {
        // 创建初始对象
        QosRule qosRule = new QosRule(null, 0, null, null);

        // 准备新的测试数据
        NetworkType newNetworkType = NetworkType.EXTRANET;
        long newConcurrentLimit = 10000;
        QpsLimitConfiguration newQpsConfig = new QpsLimitConfiguration(0, 0, 0, 0);
        BpsLimitConfiguration newBpsConfig = new BpsLimitConfiguration(0, 0, 0);

        // 设置新属性
        qosRule.setNetworkType(newNetworkType);
        qosRule.setConcurrentRequestLimit(newConcurrentLimit);
        qosRule.setQpsLimit(newQpsConfig);
        qosRule.setBpsLimit(newBpsConfig);

        // 验证setter是否生效
        assertEquals("网络类型设置不正确", newNetworkType, qosRule.getNetworkType());
        assertEquals("并发请求限制设置不正确", newConcurrentLimit, qosRule.getConcurrentRequestLimit());
        assertSame("QPS配置设置不正确", newQpsConfig, qosRule.getQpsLimit());
        assertSame("BPS配置设置不正确", newBpsConfig, qosRule.getBpsLimit());
    }

    // 测试构造函数传入负数并发请求限制时抛出异常
    @Test
    public void testConstructorWithNegativeConcurrentLimit() {
        // 准备测试数据
        NetworkType networkType = NetworkType.INTRANET;
        long negativeConcurrentLimit = -100; // 负数
        QpsLimitConfiguration qpsConfig = new QpsLimitConfiguration(1000, 2000, 3000, 4000);
        BpsLimitConfiguration bpsConfig = new BpsLimitConfiguration(5000, 6000, 7000);

        // 期望抛出IllegalArgumentException
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("concurrentRequestLimit value cannot be negative: " + negativeConcurrentLimit);

        // 尝试创建对象（应抛出异常）
        new QosRule(networkType, negativeConcurrentLimit, qpsConfig, bpsConfig);
    }

    // 测试setter方法设置负数并发请求限制时抛出异常
    @Test
    public void testSetConcurrentRequestLimitWithNegative() {
        // 创建初始对象
        QosRule qosRule = new QosRule(NetworkType.EXTRANET, 1000, null, null);
        long originalLimit = qosRule.getConcurrentRequestLimit();
        long negativeLimit = -500; // 负数

        // 期望抛出IllegalArgumentException
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("concurrentRequestLimit value cannot be negative: " + negativeLimit);

        try {
            // 尝试设置负数（应抛出异常）
            qosRule.setConcurrentRequestLimit(negativeLimit);
        } finally {
            // 验证原始值未被修改
            assertEquals("设置负数后原始值不应改变", originalLimit, qosRule.getConcurrentRequestLimit());
        }
    }
}
    