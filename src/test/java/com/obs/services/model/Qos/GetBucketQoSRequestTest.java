package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.GetBucketQoSRequest;
import org.junit.Test;
import com.obs.services.model.HttpMethodEnum;

public class GetBucketQoSRequestTest {

    private static final String TEST_BUCKET_NAME = "test-bucket-qos";

    @Test
    public void testConstructorAndHttpMethod() {
        // 创建测试对象
        GetBucketQoSRequest request = new GetBucketQoSRequest(TEST_BUCKET_NAME);
        
        // 验证桶名称是否正确设置
        assertEquals("桶名称不匹配", TEST_BUCKET_NAME, request.getBucketName());
        
        // 验证HTTP方法是否为GET
        assertEquals("HTTP方法应为GET", HttpMethodEnum.GET, request.getHttpMethod());
    }

    @Test
    public void testConstructorWithNullBucketName() {
        // 测试空桶名的情况
        GetBucketQoSRequest request = new GetBucketQoSRequest(null);
        
        // 验证空桶名是否被正确处理
        assertNull("桶名称应为null", request.getBucketName());
        
        // 即使桶名为null，HTTP方法仍应正确设置
        assertEquals("HTTP方法应为GET", HttpMethodEnum.GET, request.getHttpMethod());
    }

    @Test
    public void testConstructorWithEmptyBucketName() {
        // 测试空字符串桶名的情况
        GetBucketQoSRequest request = new GetBucketQoSRequest("");
        
        // 验证空字符串桶名是否被正确处理
        assertEquals("桶名称应为空字符串", "", request.getBucketName());
        
        // 验证HTTP方法
        assertEquals("HTTP方法应为GET", HttpMethodEnum.GET, request.getHttpMethod());
    }
}
