package com.obs.services.internal;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.ListBucketsRequest;
import com.obs.test.TestTools;
import okhttp3.Dns;
import org.junit.Assert;
import org.junit.Test;

import javax.net.ssl.SSLException;
import java.io.StringWriter;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.List;

public class NonRetrievableClassesTest {
    @Test
    public void testAddNonRetrievableClass() {
        // 初始化 log, 用于检测是否有重试
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        ObsClient obsClient = TestTools.getPipelineEnvironmentWithCustomisedDns(testDns);
        try {
            obsClient.listBuckets(new ListBucketsRequest());
        } catch (ObsException e) {
            System.out.println("testAddNonRetrievableClass writer\n" +
                    writer + "\ntestAddNonRetrievableClass writer end\n");
            Assert.assertTrue(
                    "should retry for UnknownHostException",
                    writer.toString().contains("Encountered 3 Internal Server error(s), will retry in"));
        }

        StringWriter writer2 = new StringWriter();
        TestTools.initLog(writer2);
        RestStorageService.addNonRetrievableClass(UnknownHostException.class);
        try {
            obsClient.listBuckets(new ListBucketsRequest());
        } catch (ObsException e) {
            System.out.println("testAddNonRetrievableClass writer2\n" +
                    writer2 + "\ntestAddNonRetrievableClass writer2 end\n");
            Assert.assertFalse(
                    "should not retry for UnknownHostException",
                    writer2.toString().contains("Encountered 1 Internal Server error(s), will retry in"));
        }
    }

    @Test
    public void testRemoveNonRetrievableClass() {
        // 初始化 log, 用于检测是否有重试
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        ObsClient obsClient = TestTools.getPipelineEnvironmentWithCustomisedDns(testDns);
        RestStorageService.addNonRetrievableClass(UnknownHostException.class);
        try {
            obsClient.listBuckets(new ListBucketsRequest());
        } catch (ObsException e) {
            Assert.assertFalse(
                    "should not retry for UnknownHostException",
                    writer.toString().contains("Encountered 1 Internal Server error(s), will retry in"));
        }

        RestStorageService.removeNonRetrievableClass(UnknownHostException.class);
        try {
            obsClient.listBuckets(new ListBucketsRequest());
        } catch (ObsException e) {
            Assert.assertTrue(
                    "should retry for UnknownHostException",
                    writer.toString().contains("Encountered 3 Internal Server error(s), will retry in"));
        }
    }

    @Test
    public void testGetNonRetrievableClasses() {
        RestStorageService.addNonRetrievableClass(UnknownHostException.class);
        Assert.assertTrue(
                "NonRetrievableClasses should contain UnknownHostException",
                RestStorageService.getNonRetrievableClasses().contains(UnknownHostException.class));
        RestStorageService.removeNonRetrievableClass(UnknownHostException.class);
        Assert.assertFalse(
                "NonRetrievableClasses shouldn't contain UnknownHostException",
                RestStorageService.getNonRetrievableClasses().contains(UnknownHostException.class));

        RestStorageService.addNonRetrievableClass(SSLException.class);
        Assert.assertTrue(
                "NonRetrievableClasses should contain SSLException",
                RestStorageService.getNonRetrievableClasses().contains(SSLException.class));
        RestStorageService.removeNonRetrievableClass(SSLException.class);
        Assert.assertFalse(
                "NonRetrievableClasses shouldn't contain SSLException",
                RestStorageService.getNonRetrievableClasses().contains(SSLException.class));
    }

    Dns testDns =
            new Dns() {
                /**
                 * @param s
                 * @return
                 * @throws UnknownHostException
                 */
                @Override
                public List<InetAddress> lookup(String s) throws UnknownHostException {
                    // simulate situation UnknownHostException
                    throw new UnknownHostException("test UnknownHostException");
                }
            };
}
