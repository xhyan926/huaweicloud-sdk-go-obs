package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.model.ListBucketsRequest;
import com.obs.test.TestTools;

import org.junit.Assert;
import org.junit.Test;

import java.io.IOException;
import java.io.StringWriter;

public class ClientCloseIT {
    static String FinishCloseMsg = "client closed";

    @Test
    public void testClose() throws IOException {
        // 初始化 log, 用于检测close是否有打印日志
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        obsClient.close();
        System.out.println("ClientCloseTest writer:" + writer);
        Assert.assertTrue("log should contains close", writer.toString().contains(FinishCloseMsg));
    }

    @Test
    public void testCloseTwice() throws IOException {
        // 初始化 log, 用于检测close是否有打印日志
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        obsClient.close();
        obsClient.close();
        System.out.println("ClientCloseTest writer:" + writer);
        Assert.assertTrue("log should contains close", writer.toString().contains(FinishCloseMsg));
    }

    @Test
    public void testUseAfterClose() throws IOException {
        // 初始化 log, 用于检测close是否有打印日志
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        obsClient.close();
        try {
            obsClient.listBuckets(new ListBucketsRequest());
        } catch (Exception e) {
            e.printStackTrace();
        }
        Assert.assertTrue("log should contains close", writer.toString().contains(FinishCloseMsg));
    }

    class TestThreadOfUseAfterClose extends Thread
    {
        ObsClient obsClient;
        TestThreadOfUseAfterClose(ObsClient obsClient_p){
            this.obsClient = obsClient_p;
        }
        @Override
        public void run()
        {
            try {
                obsClient.listBuckets(new ListBucketsRequest());
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    };
    @Test
    public void testUseAfterCloseInMultiThreads() throws IOException, InterruptedException
    {
        // 初始化 log, 用于检测close是否有打印日志
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        obsClient.close();
        Thread thread1 = new TestThreadOfUseAfterClose(obsClient);
        Thread thread2 = new TestThreadOfUseAfterClose(obsClient);
        thread1.start();
        thread2.start();
        thread1.join(); // 等待子线程1执行完毕
        thread2.join(); // 等待子线程2执行完毕
        Assert.assertTrue("log should contains close", writer.toString().contains(FinishCloseMsg));
    }
    @Test
    public void testUseAndClose() throws IOException {
        // 初始化 log, 用于检测close是否有打印日志
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            obsClient.listBuckets(new ListBucketsRequest());
            Assert.assertFalse("log should not contains close", writer.toString().contains(FinishCloseMsg));
        }
    }
}
