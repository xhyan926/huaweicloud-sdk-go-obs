package com.obs.integrated_test.lifecycle;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.DeleteBucketLifecycleRequest;
import com.obs.services.model.GetBucketLifecycleRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.LifecycleConfiguration;
import com.obs.services.model.SetBucketLifecycleRequest;
import com.obs.services.model.StorageClassEnum;
import com.obs.test.TestTools;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.rules.TestName;

import java.io.IOException;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

public class BucketLifecycleIT {
    protected static String bucketName = "ztw-test11";

    @org.junit.Rule
    public TestName testName = new TestName();

    @Before
    public void setUp() throws InterruptedException {
        // 在每个测试开始前执行，用于初始化
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();

        HeaderResponse delResult = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult.getStatusCode());
        //  删除完规则要稍等后台处理
        TimeUnit.SECONDS.sleep(60);
    }

    @After
    public void tearDown() throws InterruptedException {
        // 在每个测试结束后执行，用于清理
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();

        HeaderResponse delResult = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult.getStatusCode());
        //  删除完规则要稍等后台处理
        TimeUnit.SECONDS.sleep(60);
    }

    @Test
    public void tc_alpha_java_lifecycle_001() throws IOException, InterruptedException {
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();
        // 不指定ruleId设置一条生命周期，在rule规则里设置id为rule1
        LifecycleConfiguration config = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = config.new Rule();
        rule.setEnabled(false);
        rule.setId("rule1");
        rule.setPrefix("prefix1");
        LifecycleConfiguration.Transition transition = config.new Transition();
        transition.setDays(30);
        transition.setObjectStorageClass(StorageClassEnum.WARM);
        rule.getTransitions().add(transition);
        config.addRule(rule);
        SetBucketLifecycleRequest request = new SetBucketLifecycleRequest(bucketName, config);
        HeaderResponse result = obsClient.setBucketLifecycle(request);
        assertEquals(200, result.getStatusCode());

        // 指定ruleId为rule2设置一条生命周期
        LifecycleConfiguration config2 = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule2 = config2.new Rule();
        rule2.setEnabled(false);
        rule2.setId("rule2");
        rule2.setPrefix("prefix1");
        rule2.getTransitions().add(transition);
        config2.addRule(rule2);
        SetBucketLifecycleRequest request2 = new SetBucketLifecycleRequest(bucketName, "rule2", config2);
        HeaderResponse result2 = obsClient.setBucketLifecycle(request2);
        assertEquals(200, result2.getStatusCode());

        //  获取桶A生命周期
        //  新增规则需要时间合并
        TimeUnit.SECONDS.sleep(10);
        LifecycleConfiguration config3 = obsClient.getBucketLifecycle(bucketName);
        assertEquals(2, config3.getRules().size());

        HeaderResponse delResult = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult.getStatusCode());
    }

    @Test
    public void tc_alpha_java_lifecycle_002() throws IOException, InterruptedException {
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();
        // 不指定ruleId设置一条生命周期，在rule规则里设置id为rule1
        LifecycleConfiguration config = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule = config.new Rule();
        rule.setEnabled(false);
        rule.setId("rule1");
        rule.setPrefix("prefix1");
        LifecycleConfiguration.Transition transition = config.new Transition();
        transition.setDays(30);
        transition.setObjectStorageClass(StorageClassEnum.WARM);
        rule.getTransitions().add(transition);
        config.addRule(rule);
        SetBucketLifecycleRequest request = new SetBucketLifecycleRequest(bucketName, config);
        HeaderResponse result = obsClient.setBucketLifecycle(request);
        assertEquals(200, result.getStatusCode());

        // 指定ruleId为rule2设置一条生命周期
        LifecycleConfiguration config2 = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule2 = config2.new Rule();
        rule2.setEnabled(false);
        rule2.setId("rule2");
        rule2.setPrefix("prefix1");
        rule2.getTransitions().add(transition);
        config2.addRule(rule2);
        SetBucketLifecycleRequest request2 = new SetBucketLifecycleRequest(bucketName, "rule2", config2);
        HeaderResponse result2 = obsClient.setBucketLifecycle(request2);
        assertEquals(200, result2.getStatusCode());

        //  指定ruleId获取桶A生命周期
        TimeUnit.SECONDS.sleep(10);
        GetBucketLifecycleRequest request3 = new GetBucketLifecycleRequest(bucketName, "rule2");
        LifecycleConfiguration config3 = obsClient.getBucketLifecycle(request3);
        assertEquals(1, config3.getRules().size());

        HeaderResponse delResult = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult.getStatusCode());
    }

    @Test
    public void tc_alpha_java_lifecycle_003() throws IOException, InterruptedException {
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();
        // 不指定ruleId设置1001条生命周期
        LifecycleConfiguration config = new LifecycleConfiguration();
        for (int i = 0; i < 1001; i++) {
            LifecycleConfiguration.Rule rule = config.new Rule();
            rule.setEnabled(false);
            rule.setId("rule" + String.valueOf(i));
            rule.setPrefix("prefix1");
            LifecycleConfiguration.Transition transition = config.new Transition();
            transition.setDays(30);
            transition.setObjectStorageClass(StorageClassEnum.WARM);
            rule.getTransitions().add(transition);
            config.addRule(rule);
        }
        HeaderResponse result = obsClient.setBucketLifecycle(bucketName, config);
        assertEquals(200, result.getStatusCode());
        TimeUnit.SECONDS.sleep(30);
        // 不指定ruleIdMarker获取生命周期
        LifecycleConfiguration config2 = obsClient.getBucketLifecycle(bucketName);
        assertEquals(1000, config2.getRules().size());
        // 指定ruleIdMarker为上一次列举到的倒数第二条ruleId获取桶A生命周期
        GetBucketLifecycleRequest request = new GetBucketLifecycleRequest(bucketName);
        request.setRuleIdMarker("rule997");
        LifecycleConfiguration config3 = obsClient.getBucketLifecycle(request);
        assertEquals(2, config3.getRules().size());
        // 不指定ruleId设置10000条生命周期
        LifecycleConfiguration config4 = new LifecycleConfiguration();
        for (int i = 0; i < 10000; i++) {
            LifecycleConfiguration.Rule rule = config4.new Rule();
            rule.setEnabled(false);
            rule.setId("rule" + String.valueOf(i));
            rule.setPrefix("prefix1");
            LifecycleConfiguration.Transition transition = config4.new Transition();
            transition.setDays(30);
            transition.setObjectStorageClass(StorageClassEnum.WARM);
            rule.getTransitions().add(transition);
            config4.addRule(rule);
        }
        HeaderResponse result2 = obsClient.setBucketLifecycle(bucketName, config4);
        assertEquals(200, result2.getStatusCode());
        //  等待60s，等服务端合并rule
        TimeUnit.SECONDS.sleep(60);
        // 循环指定ruleIdMarke为上一次列举到的最后一条ruleId，统计列举总数和列举次数
        long ruleCount = 0;
        long listCount = 0;
        String Marker = "rule";
        while (true) {
            GetBucketLifecycleRequest request2 = new GetBucketLifecycleRequest(bucketName);
            request2.setRuleIdMarker(Marker);
            try{
                LifecycleConfiguration config5 = obsClient.getBucketLifecycle(request2);
                ruleCount += config5.getRules().size();
                listCount += 1;
                if (config5.getRules().size() == 1000) {
                    List<LifecycleConfiguration.Rule>  rules = config5.getRules();
                    Marker = rules.get(999).getId();
                }else {
                    break;
                }
            } catch (ObsException obsException) {
                break;
            }
        }
        assertEquals(10000, ruleCount);
        assertEquals(10, listCount);

        HeaderResponse delResult2 = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult2.getStatusCode());
        TimeUnit.SECONDS.sleep(120);
    }

    @Test
    public void tc_alpha_java_lifecycle_004() throws IOException, InterruptedException {
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();
        // 不指定ruleId设置两条生命周期，在rule规则里设置id分别为rule1和rule2
        LifecycleConfiguration config = new LifecycleConfiguration();
        for (int i = 1; i < 3; i++) {
            LifecycleConfiguration.Rule rule = config.new Rule();
            rule.setEnabled(false);
            rule.setId("rule" + String.valueOf(i));
            rule.setPrefix("prefix1");
            LifecycleConfiguration.Transition transition = config.new Transition();
            // 指定满足前缀的对象创建30天后转换
            transition.setDays(30);
            // 指定对象转换后的存储类型
            transition.setObjectStorageClass(StorageClassEnum.WARM);
            // 直接指定满足前缀的对象转换日期
            rule.getTransitions().add(transition);
            config.addRule(rule);
        }
        HeaderResponse result = obsClient.setBucketLifecycle(bucketName, config);
        assertEquals(200, result.getStatusCode());
        TimeUnit.SECONDS.sleep(10);
        //获取桶A生命周期
        LifecycleConfiguration config2 = obsClient.getBucketLifecycle(bucketName);
        assertEquals(2, config2.getRules().size());
        // 指定ruleId为rule2删除一条生命周期
        DeleteBucketLifecycleRequest request = new DeleteBucketLifecycleRequest(bucketName, "rule2");
        HeaderResponse delResult = obsClient.deleteBucketLifecycle(request);
        assertEquals(204, delResult.getStatusCode());
        //获取桶A生命周期
        TimeUnit.SECONDS.sleep(10);
        LifecycleConfiguration config3 = obsClient.getBucketLifecycle(bucketName);
        assertEquals(1, config3.getRules().size());
        HeaderResponse delResult2 = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult2.getStatusCode());
    }

    @Test
    public void tc_alpha_java_lifecycle_005() throws IOException, InterruptedException {
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();
        // .不指定ruleId设置10001条生命周期
        LifecycleConfiguration config = new LifecycleConfiguration();
        for (int i = 0; i < 10001; i++) {
            LifecycleConfiguration.Rule rule = config.new Rule();
            rule.setEnabled(false);
            rule.setId("rule" + String.valueOf(i));
            rule.setPrefix("prefix1");
            LifecycleConfiguration.Transition transition = config.new Transition();
            // 指定满足前缀的对象创建30天后转换
            transition.setDays(30);
            // 指定对象转换后的存储类型
            transition.setObjectStorageClass(StorageClassEnum.WARM);
            // 直接指定满足前缀的对象转换日期
            rule.getTransitions().add(transition);
            config.addRule(rule);
        }
        try{
            HeaderResponse result = obsClient.setBucketLifecycle(bucketName, config);
            fail("Scenario should fail with 400, but got: " + result.getStatusCode());
        } catch (ObsException obsException) {
            assertEquals(400, obsException.getResponseCode());
        }
        TimeUnit.SECONDS.sleep(60);
        //不指定ruleId设置10000条生命周期
        LifecycleConfiguration config2 = new LifecycleConfiguration();
        for (int i = 0; i < 10000; i++) {
            LifecycleConfiguration.Rule rule = config2.new Rule();
            rule.setEnabled(false);
            rule.setId("rule" + String.valueOf(i));
            rule.setPrefix("prefix1");
            LifecycleConfiguration.Transition transition = config2.new Transition();
            // 指定满足前缀的对象创建30天后转换
            transition.setDays(30);
            // 指定对象转换后的存储类型
            transition.setObjectStorageClass(StorageClassEnum.WARM);
            // 直接指定满足前缀的对象转换日期
            rule.getTransitions().add(transition);
            config2.addRule(rule);
        }
        HeaderResponse result2 = obsClient.setBucketLifecycle(bucketName, config2);
        assertEquals(200, result2.getStatusCode());
        TimeUnit.SECONDS.sleep(60);

        // 指定ruleId为rule10001设置一条生命周期
        LifecycleConfiguration config3 = new LifecycleConfiguration();
        LifecycleConfiguration.Rule rule2 = config3.new Rule();
        rule2.setEnabled(false);
        rule2.setId("rule200000");
        rule2.setPrefix("prefix1");
        LifecycleConfiguration.Transition transition = config3.new Transition();
        // 指定满足前缀的对象创建30天后转换
        transition.setDays(30);
        // 指定对象转换后的存储类型
        transition.setObjectStorageClass(StorageClassEnum.WARM);
        rule2.getTransitions().add(transition);
        config3.addRule(rule2);
        SetBucketLifecycleRequest request2 = new SetBucketLifecycleRequest(bucketName, "rule200000", config3);
        try{
            HeaderResponse result3 = obsClient.setBucketLifecycle(request2);
            fail("Scenario should fail with 400, but got: " + result3.getStatusCode());
        } catch (ObsException obsException) {
            assertEquals(400, obsException.getResponseCode());
        }
        TimeUnit.SECONDS.sleep(10);
        HeaderResponse delResult = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult.getStatusCode());
        TimeUnit.SECONDS.sleep(120);
    }

    @Test
    public void tc_alpha_java_lifecycle_006() throws IOException, InterruptedException {
        ObsClient obsClient = TestTools.getPipelineForLifecycleEnvironment();
        // .不指定ruleId设置10000条生命周期
        LifecycleConfiguration config = new LifecycleConfiguration();
        for (int i = 0; i < 10000; i++) {
            LifecycleConfiguration.Rule rule = config.new Rule();
            rule.setEnabled(false);
            rule.setId("rule" + String.valueOf(i));
            rule.setPrefix("prefix1");
            LifecycleConfiguration.Transition transition = config.new Transition();
            // 指定满足前缀的对象创建30天后转换
            transition.setDays(30);
            // 指定对象转换后的存储类型
            transition.setObjectStorageClass(StorageClassEnum.WARM);
            // 直接指定满足前缀的对象转换日期
            rule.getTransitions().add(transition);
            config.addRule(rule);
        }
        HeaderResponse result = obsClient.setBucketLifecycle(bucketName, config);
        assertEquals(200, result.getStatusCode());
        TimeUnit.SECONDS.sleep(60);

        // 不指定ruleId删除桶A生命周期
        HeaderResponse delResult = obsClient.deleteBucketLifecycle(bucketName);
        assertEquals(204, delResult.getStatusCode());
        TimeUnit.SECONDS.sleep(180);

        // 获取桶A生命周期
        try{
            LifecycleConfiguration config2 = obsClient.getBucketLifecycle(bucketName);
            fail("Scenario should fail with 404, but got: " + config2.getStatusCode());
        } catch (ObsException obsException) {
            assertEquals(404, obsException.getResponseCode());
        }
    }
}
