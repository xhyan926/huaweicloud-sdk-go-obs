/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.model.lifecycle;

import com.obs.services.model.SetBucketLifecycleRequest;
import org.junit.Test;
import static org.junit.Assert.assertEquals;

public class SetBucketLifecycleRequestUnitTest {

    public static String bucketNameForTest = "test-bucket-lifecycle";

    @Test
    public void test_SetBucketLifecycleRequest_with_all_params() {
        // 创建测试对象
        SetBucketLifecycleRequest request = new SetBucketLifecycleRequest(bucketNameForTest, "validRuleId",
            null);

        // 验证RuleId是否正确设置
        assertEquals("参数不符", "validRuleId", request.getRuleId());
    }

    @Test
    public void test_SetBucketLifecycleRequest_with_bucketName_and_lifecycleConfig() {
        // 创建测试对象
        SetBucketLifecycleRequest request = new SetBucketLifecycleRequest(bucketNameForTest, null);
        request.setRuleId("validRuleId");

        // 验证RuleId是否正确设置
        assertEquals("参数不符", "validRuleId", request.getRuleId());
    }


}