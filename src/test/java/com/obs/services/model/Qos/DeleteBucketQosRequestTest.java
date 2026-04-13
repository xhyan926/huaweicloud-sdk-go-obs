package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.DeleteBucketQosRequest;
import org.junit.Test;
import com.obs.services.model.HttpMethodEnum;

public class DeleteBucketQosRequestTest {

    private static final String TEST_BUCKET = "test-delete-qos-bucket";

    @Test
    public void testConstructorAndHttpMethod() {
        // 创建请求对象
        DeleteBucketQosRequest request = new DeleteBucketQosRequest(TEST_BUCKET);
        
        // 验证桶名称设置
        assertEquals("桶名称不匹配", TEST_BUCKET, request.getBucketName());
        
        // 验证HTTP方法是否为DELETE
        assertEquals("HTTP方法应为DELETE", HttpMethodEnum.DELETE, request.getHttpMethod());
    }

    @Test
    public void testConstructorWithNullBucketName() {
        // 测试传入null作为桶名
        DeleteBucketQosRequest request = new DeleteBucketQosRequest(null);
        
        // 验证null桶名的处理
        assertNull("桶名称应为null", request.getBucketName());
        
        // 验证HTTP方法仍为DELETE
        assertEquals("HTTP方法应为DELETE", HttpMethodEnum.DELETE, request.getHttpMethod());
    }

    @Test
    public void testConstructorWithEmptyBucketName() {
        // 测试传入空字符串作为桶名
        DeleteBucketQosRequest request = new DeleteBucketQosRequest("");
        
        // 验证空字符串桶名的处理
        assertEquals("桶名称应为空字符串", "", request.getBucketName());
        
        // 验证HTTP方法
        assertEquals("HTTP方法应为DELETE", HttpMethodEnum.DELETE, request.getHttpMethod());
    }
}
