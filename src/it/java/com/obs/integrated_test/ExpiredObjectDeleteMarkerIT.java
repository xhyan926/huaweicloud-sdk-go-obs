package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.BucketTagInfo;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.LifecycleConfiguration;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.Calendar;
import java.util.Date;
import java.util.Locale;

import static junit.framework.TestCase.assertEquals;
import static junit.framework.TestCase.assertNotNull;
import static junit.framework.TestCase.assertNull;
import static junit.framework.TestCase.assertTrue;
import static junit.framework.TestCase.fail;

public class ExpiredObjectDeleteMarkerIT
{
    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Test
    public void tc_ExpirationWithExpiredObjectDeleteMarker() {
        // 测试生命周期只设置Expiration，子元素为ExpiredObjectDeleteMarker
        // 设置桶生命周期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String ruleId = "tc_set_bucket_lifecycle_011";
        String rulePrefix = "tc_set_bucket_lifecycle_rulePrefix_011";
        ObsClient obsClient = TestTools.getPipelineEnvironment();

        LifecycleConfiguration lifecycleConfiguration = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = lifecycleConfiguration.new Rule(ruleId, rulePrefix, true);
        LifecycleConfiguration.Expiration expiration = lifecycleConfiguration.new Expiration();
        Boolean testExpiredObjectDeleteMarker = true;
        expiration.setExpiredObjectDeleteMarker(testExpiredObjectDeleteMarker);
        rule.setExpiration(expiration);
        lifecycleConfiguration.getRules().add(rule);
        try {
            assert obsClient != null;
            HeaderResponse headerResponse = obsClient.setBucketLifecycle(bucketName, lifecycleConfiguration);
            assertEquals(200, headerResponse.getStatusCode());
            LifecycleConfiguration getLifecycleConfiguration = obsClient.getBucketLifecycle(bucketName);
            assertEquals(200, getLifecycleConfiguration.getStatusCode());
            assertNotNull(getLifecycleConfiguration.getRules());
            assertNotNull(getLifecycleConfiguration.getRules().get(0));
            assertNotNull(getLifecycleConfiguration.getRules().get(0).getExpiration());
            assertNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getDays());
            assertNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getDate());
            assertNotNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getExpiredObjectDeleteMarker());
            assertEquals(testExpiredObjectDeleteMarker,
                    getLifecycleConfiguration.getRules().get(0).getExpiration().getExpiredObjectDeleteMarker());
        } catch (ObsException exception) {
            System.out.println("ErrorMessage" + exception.getErrorMessage());
            exception.printStackTrace();
            fail();
        }
    }
    @Test
    public void tc_ExpirationWithDays() {
        // 测试生命周期设置Expiration，子元素为Days
        // 设置桶生命周期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String ruleId = "tc_set_bucket_lifecycle_011";
        String rulePrefix = "tc_set_bucket_lifecycle_rulePrefix_011";
        ObsClient obsClient = TestTools.getPipelineEnvironment();

        LifecycleConfiguration lifecycleConfiguration = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = lifecycleConfiguration.new Rule(ruleId, rulePrefix, true);
        LifecycleConfiguration.Expiration expiration = lifecycleConfiguration.new Expiration();
        Integer testDays = 30;
        expiration.setDays(testDays);
        rule.setExpiration(expiration);
        lifecycleConfiguration.getRules().add(rule);
        try {
            assert obsClient != null;
            HeaderResponse headerResponse = obsClient.setBucketLifecycle(bucketName, lifecycleConfiguration);
            assertEquals(200, headerResponse.getStatusCode());
            LifecycleConfiguration getLifecycleConfiguration = obsClient.getBucketLifecycle(bucketName);
            assertEquals(200, getLifecycleConfiguration.getStatusCode());
            assertNotNull(getLifecycleConfiguration.getRules());
            assertNotNull(getLifecycleConfiguration.getRules().get(0));
            assertNotNull(getLifecycleConfiguration.getRules().get(0).getExpiration());
            assertNotNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getDays());
            assertNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getDate());
            assertNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getExpiredObjectDeleteMarker());
            assertEquals(testDays, getLifecycleConfiguration.getRules().get(0).getExpiration().getDays());
        } catch (ObsException exception) {
            System.out.println("ErrorMessage" + exception.getErrorMessage());
            exception.printStackTrace();
            fail();
        }
    }
    @Test
    public void tc_ExpirationWithDate() {
        // 测试生命周期设置Expiration，子元素为Date
        // 设置桶生命周期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String ruleId = "tc_set_bucket_lifecycle_011";
        String rulePrefix = "tc_set_bucket_lifecycle_rulePrefix_011";
        ObsClient obsClient = TestTools.getPipelineEnvironment();

        LifecycleConfiguration lifecycleConfiguration = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = lifecycleConfiguration.new Rule(ruleId, rulePrefix, true);
        LifecycleConfiguration.Expiration expiration = lifecycleConfiguration.new Expiration();
        Date now = new Date();
        expiration.setDate(now);
        rule.setExpiration(expiration);
        lifecycleConfiguration.getRules().add(rule);
        try {
            assert obsClient != null;
            HeaderResponse headerResponse = obsClient.setBucketLifecycle(bucketName, lifecycleConfiguration);
            assertEquals(200, headerResponse.getStatusCode());
            LifecycleConfiguration getLifecycleConfiguration = obsClient.getBucketLifecycle(bucketName);
            assertEquals(200, getLifecycleConfiguration.getStatusCode());
            assertNotNull(getLifecycleConfiguration.getRules());
            assertNotNull(getLifecycleConfiguration.getRules().get(0));
            assertNotNull(getLifecycleConfiguration.getRules().get(0).getExpiration());
            assertNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getDays());
            assertNotNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getDate());
            assertNull(getLifecycleConfiguration.getRules().get(0).getExpiration().getExpiredObjectDeleteMarker());
            Date getDate = getLifecycleConfiguration.getRules().get(0).getExpiration().getDate();

            Calendar cal1 = Calendar.getInstance();
            cal1.setTime(now);
            Calendar cal2 = Calendar.getInstance();
            cal2.setTime(getDate);

            // 两个date的日期相同即可
            assertTrue(cal1.get(Calendar.YEAR) == cal2.get(Calendar.YEAR) &&
                    cal1.get(Calendar.MONTH) == cal2.get(Calendar.MONTH) &&
                    cal1.get(Calendar.DAY_OF_MONTH) == cal2.get(Calendar.DAY_OF_MONTH));
        } catch (ObsException exception) {
            System.out.println("ErrorMessage" + exception.getErrorMessage());
            exception.printStackTrace();
            fail();
        }
    }
    @Test
    public void tc_EmptyExpiration() {
        // 测试生命周期设置Expiration,但不设置任何子元素
        // 设置桶生命周期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String ruleId = "tc_set_bucket_lifecycle_011";
        String rulePrefix = "tc_set_bucket_lifecycle_rulePrefix_011";
        ObsClient obsClient = TestTools.getPipelineEnvironment();

        LifecycleConfiguration lifecycleConfiguration = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = lifecycleConfiguration.new Rule(ruleId, rulePrefix, true);
        LifecycleConfiguration.Expiration expiration = lifecycleConfiguration.new Expiration();
        rule.setExpiration(expiration);
        lifecycleConfiguration.getRules().add(rule);
        try {
            assert obsClient != null;
            obsClient.setBucketLifecycle(bucketName, lifecycleConfiguration);
            fail();
        } catch (ObsException exception) {
            // 检测设置结果，400为请求参数异常
            assertEquals(400, exception.getResponseCode());
        }
    }
    @Test
    public void tc_ExpirationWithExpiredObjectDeleteMarkerAndTag() {
        // 测试生命周期只设置Expiration，子元素为ExpiredObjectDeleteMarker,同时设置tag, 抛出异常
        // 设置桶生命周期
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String ruleId = "tc_set_bucket_lifecycle_011";
        String rulePrefix = "tc_set_bucket_lifecycle_rulePrefix_011";
        ObsClient obsClient = TestTools.getPipelineEnvironment();

        LifecycleConfiguration lifecycleConfiguration = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = lifecycleConfiguration.new Rule(ruleId, rulePrefix, true);
        LifecycleConfiguration.Expiration expiration = lifecycleConfiguration.new Expiration();
        Boolean testExpiredObjectDeleteMarker = true;
        expiration.setExpiredObjectDeleteMarker(testExpiredObjectDeleteMarker);
        rule.setExpiration(expiration);
        BucketTagInfo.TagSet tagSet = new BucketTagInfo.TagSet();
        tagSet.addTag("testTagKey", "testTagValue");
        rule.setTagSet(tagSet);
        lifecycleConfiguration.getRules().add(rule);
        try {
            assert obsClient != null;
            obsClient.setBucketLifecycle(bucketName, lifecycleConfiguration);
            fail();
        } catch (ObsException exception) {
            // 检测设置结果，400为请求参数异常
            assertEquals(400, exception.getResponseCode());
        }
    }
}
