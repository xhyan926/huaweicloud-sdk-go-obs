/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.model.lifecycle;

import com.obs.services.model.DeleteBucketLifecycleRequest;
import org.junit.Test;

import static org.junit.Assert.assertEquals;

public class DeleteBucketLifecycleRequestUnitTest {

    public static String bucketNameForTest = "test-bucket-lifecycle";

    @Test
    public void test_DeleteBucketLifecycleRequest_with_all_params() {
        // 创建测试对象
        DeleteBucketLifecycleRequest request = new DeleteBucketLifecycleRequest(bucketNameForTest, "validRuleId");

        // 验证RuleId是否正确设置
        assertEquals("参数不符", "validRuleId", request.getRuleId());
    }

    @Test
    public void test_DeleteBucketLifecycleRequest_with_bucketName() {
        // 创建测试对象
        DeleteBucketLifecycleRequest request = new DeleteBucketLifecycleRequest(bucketNameForTest);
        request.setRuleId("validRuleId");

        // 验证RuleId是否正确设置
        assertEquals("参数不符", "validRuleId", request.getRuleId());
    }
}