/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.model.lifecycle;

import com.obs.services.model.GetBucketLifecycleRequest;
import org.junit.Test;
import static org.junit.Assert.assertEquals;

public class GetBucketLifecycleRequestUnitTest {

    public static String bucketNameForTest = "test-bucket-lifecycle";

    @Test
    public void test_GetBucketLifecycleRequest_with_all_params() {
        // 创建测试对象
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest, "validRuleId",
            "validRuleIdMarker");

        // 验证RuleId是否正确设置
        assertEquals("参数不符", "validRuleId", request.getRuleId());

        // 验证RuleIdMarker是否正确设置
        assertEquals("参数不符", "validRuleIdMarker", request.getRuleIdMarker());
    }

    @Test
    public void test_GetBucketLifecycleRequest_with_bucketName_and_ruleIdMarker() {
        // 测试空桶名的情况
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest);
        request.setRuleIdMarker("validRuleIdMarker");

        // 即使桶名为null，HTTP方法仍应正确设置
        assertEquals("参数不符", "validRuleIdMarker", request.getRuleIdMarker());
    }

    @Test
    public void test_GetBucketLifecycleRequest_with_bucketName_and_ruleId() {
        // 创建测试对象
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketNameForTest, "validRuleId");

        // 验证RuleId是否正确设置
        assertEquals("参数不符", "validRuleId", request.getRuleId());

        GetBucketLifecycleRequest request2 = new GetBucketLifecycleRequest(bucketNameForTest);
        request2.setRuleId("validRuleId2");

        // 验证RuleId是否正确设置
        assertEquals("参数不符", "validRuleId2", request2.getRuleId());
    }
}