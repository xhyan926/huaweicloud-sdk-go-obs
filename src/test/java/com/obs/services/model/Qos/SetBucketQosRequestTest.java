package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.QosConfiguration;
import com.obs.services.model.Qos.SetBucketQosRequest;
import org.junit.Test;
import com.obs.services.model.HttpMethodEnum;

public class SetBucketQosRequestTest {

    private static final String TEST_BUCKET = "test-bucket-123";

    @Test
    public void testConstructorAndGetters() {
        // 准备测试数据
        QosConfiguration testConfig = new QosConfiguration();
        
        // 创建请求对象
        SetBucketQosRequest request = new SetBucketQosRequest(TEST_BUCKET, testConfig);
        
        // 验证基本属性
        assertEquals("桶名称不匹配", TEST_BUCKET, request.getBucketName());
        assertSame("QoS配置对象不匹配", testConfig, request.getQosConfig());
        assertEquals("HTTP方法应为PUT", HttpMethodEnum.PUT, request.getHttpMethod());
    }

    @Test
    public void testSetQosConfiguration() {
        // 创建初始请求（带空配置）
        SetBucketQosRequest request = new SetBucketQosRequest(TEST_BUCKET, null);
        assertNull("初始QoS配置应为null", request.getQosConfig());
        
        // 创建新的配置对象并设置
        QosConfiguration newConfig = new QosConfiguration();
        request.setQosConfiguration(newConfig);
        
        // 验证配置已更新
        assertSame("QoS配置未正确更新", newConfig, request.getQosConfig());
    }

    @Test
    public void testNullBucketName() {
        // 测试空桶名的情况
        QosConfiguration config = new QosConfiguration();
        SetBucketQosRequest request = new SetBucketQosRequest(null, config);
        
        assertNull("桶名称应为null", request.getBucketName());
        assertSame("QoS配置对象不匹配", config, request.getQosConfig());
        assertEquals("HTTP方法设置不正确", HttpMethodEnum.PUT, request.getHttpMethod());
    }

    @Test
    public void testNullQosConfigInConstructor() {
        // 测试构造函数传入null配置
        SetBucketQosRequest request = new SetBucketQosRequest(TEST_BUCKET, null);
        
        assertEquals("桶名称不匹配", TEST_BUCKET, request.getBucketName());
        assertNull("QoS配置应为null", request.getQosConfig());
    }

    @Test
    public void testSetQosConfigurationToNull() {
        // 测试将配置设置为null
        QosConfiguration initialConfig = new QosConfiguration();
        SetBucketQosRequest request = new SetBucketQosRequest(TEST_BUCKET, initialConfig);
        
        request.setQosConfiguration(null);
        assertNull("QoS配置应被设置为null", request.getQosConfig());
    }
}
