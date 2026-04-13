package com.obs.services.internal;

import com.obs.services.ObsClient;
import com.obs.services.internal.utils.OkhttpCallProfiler;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObsBucket;
import com.obs.test.TestTools;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.Locale;

public class EventListenerFactoryTest {
    @Rule
    public TestName testName = new TestName();
    //  设置自定义dns后，创建、head、删除桶成功
    @Test
    public void test_EventListenerFactoryTest01() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        StringBuilder stringBuilder = new StringBuilder();
        // 设置自定义EventListenerFactory
        ObsClient obsClient = TestTools.getPipelineEnvironmentWithEventListenerFactory(
                call -> new OkhttpCallProfiler(stringBuilder::append));
        assert obsClient != null;
        // 创建桶成功
        ObsBucket obsBucket = obsClient.createBucket(bucketName);
        Assert.assertEquals(200, obsBucket.getStatusCode());
        // head桶成功
        boolean bucketExists = obsClient.headBucket(bucketName);
        Assert.assertTrue(bucketExists);
        // 删除桶成功
        HeaderResponse response = obsClient.deleteBucket(bucketName);
        Assert.assertEquals(204, response.getStatusCode());
        // 成功输出统计信息
        Assert.assertTrue(stringBuilder.length() > 0);
        Assert.assertTrue(stringBuilder.toString().contains("call cost time"));
    }
}
