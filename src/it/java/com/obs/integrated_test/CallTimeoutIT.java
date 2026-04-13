/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.model.PutObjectResult;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.util.Locale;

import static com.obs.test.TestTools.genTestFile;

public class CallTimeoutIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder(new File("."));

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();
    private static final String TIMEOUT_MESSAGE = "timeout";
    @Test
    public void tc_callTimeout_normal() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        int testCallTimeout = 500;
        obsConfiguration.setCallTimeout(testCallTimeout);
        ObsClient obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        Assert.assertNotNull(obsClient);
        // 2 gb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 2 * 1024 * 1024 * 1024L);
        try {
            obsClient.putObject(bucketName, testFileName, testFile);
            Assert.fail();
        } catch (ObsException e){
            // 设置Java SDK的callTimeout为500, 上传2 gb文件提前终止，抛出异常timeout
            Assert.assertTrue(e.getErrorMessage().contains(TIMEOUT_MESSAGE));
        }
    }
    @Test
    public void tc_callTimeout_boundary() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        obsConfiguration.setCallTimeout(0);
        ObsClient obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        Assert.assertNotNull(obsClient);
        // 10 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 10 * 1024 * 1024L);
        PutObjectResult putObjectResult = obsClient.putObject(bucketName, testFileName, testFile);
        // 上传正常返回200，callTimeout为0，为默认值，不会触发超时
        Assert.assertEquals(200, putObjectResult.getStatusCode());
        obsConfiguration.setCallTimeout(Integer.MAX_VALUE);
        obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
        putObjectResult = obsClient.putObject(bucketName, testFileName, testFile);
        // 上传正常返回200，callTimeout为Integer.MAX_VALUE，远大于上传耗时，所以不会触发超时
        Assert.assertEquals(200, putObjectResult.getStatusCode());

        obsConfiguration.setCallTimeout(1);
        try {
            obsClient = TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
            Assert.assertNotNull(obsClient);
            obsClient.putObject(bucketName, testFileName, testFile);
            Assert.fail();
        } catch (ObsException e){
            // 设置Java SDK的callTimeout为1, 上传提前终止，抛出异常timeout
            Assert.assertTrue(e.getErrorMessage().contains(TIMEOUT_MESSAGE));
        }
    }
    @Test
    public void tc_callTimeout_illegal() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsConfiguration obsConfiguration = new ObsConfiguration();
        int testCallTimeout = -1;
        obsConfiguration.setCallTimeout(testCallTimeout);
        try {
            TestTools.getPipelineEnvironmentClientWithConfig(obsConfiguration);
            Assert.fail();
        } catch (Throwable t){
            // 设置失败，抛出异常提示callTimeou值非法
            Assert.assertTrue(t instanceof IllegalStateException);
        }
    }

    @Test
    public void tc_no_callTimeout() throws IOException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String testFileName = bucketName + "testFile";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        Assert.assertNotNull(obsClient);
        // 10 mb test file
        File testFile = genTestFile(temporaryFolder, testFileName, 10 * 1024 * 1024L);
        PutObjectResult putObjectResult = obsClient.putObject(bucketName, testFileName, testFile);
        // 不设置callTimeout, 上传正常返回200，不会触发超时
        Assert.assertEquals(200, putObjectResult.getStatusCode());
    }
}
