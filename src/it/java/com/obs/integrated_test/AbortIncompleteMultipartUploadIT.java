package com.obs.integrated_test;

import static junit.framework.TestCase.assertEquals;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.LifecycleConfiguration;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.Locale;

public class AbortIncompleteMultipartUploadIT {
    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Test
    public void tc_set_bucket_lifecycle_011() {
        // 传入正确的桶名、无效的生命周期配置（AbortIncompleteMultipartUpload为负），设置桶的生命周期配置失败
        // 设置桶生命周期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String ruleId = "tc_set_bucket_lifecycle_011";
        String rulePrefix = "tc_set_bucket_lifecycle_rulePrefix_011";
        ObsClient obsClient = TestTools.getPipelineEnvironment();

        LifecycleConfiguration lifecycleConfiguration = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = lifecycleConfiguration.new Rule(ruleId, rulePrefix, true);
        LifecycleConfiguration.AbortIncompleteMultipartUpload abortIncompleteMultipartUpload =
                lifecycleConfiguration.new AbortIncompleteMultipartUpload();

        abortIncompleteMultipartUpload.setDaysAfterInitiation(-1);
        rule.setAbortIncompleteMultipartUpload(abortIncompleteMultipartUpload);
        lifecycleConfiguration.getRules().add(rule);

        try {
            assert obsClient != null;
            obsClient.setBucketLifecycle(bucketName, lifecycleConfiguration);
        } catch (ObsException exception) {
            // 检测设置结果，400为请求参数异常
            assertEquals(400, exception.getResponseCode());
        }
    }

    @Test
    public void tc_set_bucket_lifecycle_012() {
        // 传入正确的桶名、正确的生命周期配置（AbortIncompleteMultipartUpload），设置桶的生命周期配置成功
        // 设置桶生命周期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String ruleId = "tc_set_bucket_lifecycle_012";
        String rulePrefix = "tc_set_bucket_lifecycle_rulePrefix_012";
        ObsClient obsClient = TestTools.getPipelineEnvironment();

        LifecycleConfiguration lifecycleConfiguration = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = lifecycleConfiguration.new Rule(ruleId, rulePrefix, true);
        LifecycleConfiguration.AbortIncompleteMultipartUpload abortIncompleteMultipartUpload =
                lifecycleConfiguration.new AbortIncompleteMultipartUpload();

        abortIncompleteMultipartUpload.setDaysAfterInitiation(10);
        rule.setAbortIncompleteMultipartUpload(abortIncompleteMultipartUpload);
        lifecycleConfiguration.getRules().add(rule);

        assert obsClient != null;
        HeaderResponse response1 = obsClient.setBucketLifecycle(bucketName, lifecycleConfiguration);
        // 检测设置结果，200为成功
        assertEquals(200, response1.getStatusCode());
    }
}
