package com.obs.services.internal;

import com.obs.services.ObsClient;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObsBucket;
import com.obs.test.TestTools;

import okhttp3.Dns;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.List;
import java.util.Locale;

public class CustomizedDnsTest {
    @Rule
    public TestName testName = new TestName();
    //  设置自定义dns后，创建、head、删除桶成功
    @Test
    public void test_CustomizedDns001() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        CustomizedDnsDemo customizedDnsDemo = new CustomizedDnsDemo();
        // 状态默认为false，说明还没有调用到自定义dns的lookup函数
        Assert.assertFalse(customizedDnsDemo.isUsedCustomizedDns());
        // 设置自定义dns
        ObsClient obsClient = TestTools.getPipelineEnvironmentWithCustomisedDns(customizedDnsDemo);
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
        // 状态变为true，说明已经调用了自定义dns的lookup函数
        Assert.assertTrue(customizedDnsDemo.isUsedCustomizedDns());
    }

    public static class CustomizedDnsDemo implements Dns {
        /**
         * @param hostname
         * @return
         * @throws UnknownHostException
         */
        @Override
        public List<InetAddress> lookup(String hostname) throws UnknownHostException {
            setUsedCustomizedDns(true);
            return Dns.SYSTEM.lookup(hostname);
        }

        public boolean isUsedCustomizedDns() {
            return usedCustomizedDns;
        }

        public void setUsedCustomizedDns(boolean usedCustomizedDns) {
            this.usedCustomizedDns = usedCustomizedDns;
        }

        private boolean usedCustomizedDns = false;
    }
}
