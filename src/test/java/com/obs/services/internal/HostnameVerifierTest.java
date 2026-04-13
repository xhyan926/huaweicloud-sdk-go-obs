package com.obs.services.internal;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import okhttp3.internal.tls.OkHostnameVerifier;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import javax.net.ssl.HttpsURLConnection;
import java.io.IOException;
import java.util.Locale;

public class HostnameVerifierTest
{
    @Rule
    public TestName testName = new TestName();
    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();
    @Test
    public void test_HostnameWithDot_Verify_Test() {
        String bucketName = "test.hostname.with.dot";
        ObsException obsException = null;
        try (ObsClient obsClient = TestTools.getPipelineEnvironmentWithHostnameVerifier(OkHostnameVerifier.INSTANCE)) {
            assert obsClient != null;
            obsClient.headBucket(bucketName);
        }
        catch (ObsException e){
            obsException = e;
        }
        catch (IOException e)
        {
            e.printStackTrace();
        }
        assert obsException != null;
        Assert.assertTrue("hostname which contains dot can't be verified", obsException.getCause().getMessage().contains("not verified"));
    }
    @Test
    public void test_HostnameWithoutDot_Verify_Test() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        try (ObsClient obsClient = TestTools.getPipelineEnvironmentWithHostnameVerifier(OkHostnameVerifier.INSTANCE)) {
            assert obsClient != null;
            obsClient.headBucket(bucketName);
        }
        catch (ObsException | IOException e){
            e.printStackTrace();
        }
    }
    @Test
    public void test_HostnameWithDot_Verify_Old_Test() {
        String bucketName = "test.hostname.with.dot";
        HttpsURLConnection.setDefaultHostnameVerifier(OkHostnameVerifier.INSTANCE);
        try (ObsClient obsClient = TestTools.getPipelineEnvironmentWithHostnameVerifier(null)) {
            assert obsClient != null;
            obsClient.headBucket(bucketName);
        }
        catch (ObsException | IOException e){
            e.printStackTrace();
        }
    }
    @Test
    public void test_HostnameWithoutDot_Verify_Old_Test() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        HttpsURLConnection.setDefaultHostnameVerifier(OkHostnameVerifier.INSTANCE);
        try (ObsClient obsClient = TestTools.getPipelineEnvironmentWithHostnameVerifier(null)) {
            assert obsClient != null;
            obsClient.headBucket(bucketName);
        }
        catch (ObsException | IOException e){
            e.printStackTrace();
        }
    }
}
