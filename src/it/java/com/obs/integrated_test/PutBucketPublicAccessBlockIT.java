/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.integrated_test;

import static org.junit.Assert.fail;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.AccessControlList;
import com.obs.services.model.bpa.BucketPublicAccessBlock;
import com.obs.services.model.bpa.DeleteBucketPublicAccessBlockRequest;
import com.obs.services.model.bpa.GetBucketPolicyPublicStatusRequest;
import com.obs.services.model.bpa.GetBucketPolicyPublicStatusResult;
import com.obs.services.model.bpa.GetBucketPublicAccessBlockRequest;
import com.obs.services.model.bpa.GetBucketPublicAccessBlockResult;
import com.obs.services.model.bpa.GetBucketPublicStatusRequest;
import com.obs.services.model.bpa.GetBucketPublicStatusResult;
import com.obs.services.model.bpa.PutBucketPublicAccessBlockRequest;
import com.obs.services.model.BucketPolicyResponse;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.SetBucketAclRequest;
import com.obs.services.model.SetBucketPolicyRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import com.obs.test.tools.TestCaseIgnore;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;
import org.junit.rules.TestName;

import java.io.IOException;
import java.util.Locale;

public class PutBucketPublicAccessBlockIT {
    @Rule
    public TemporaryFolder temporaryFolder = new TemporaryFolder();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Rule
    public TestName testName = new TestName();

    /***
     * 1、用户创桶A
     * 2、调用查询接口GetBucketPublicAccessBlock查看桶BPA状态
     * 3、调用设置接口PutBucketPublicAccessBlock设置BPA开关，状态分别为
     * BlockPublicAcls:false     BlockPublicPolicy:FALSE
     * IgnorePublicAcls:FaLsE      RestrictPublicBuckets:faLSE
     * 4、调用查询接口GetBucketPublicAccessBlock查看桶BPA状态
     * 5、调用设置接口PutBucketPublicAccessBlock设置BPA开关，状态分别为
     * BlockPublicAcls:true       BlockPublicPolicy:TRUE
     * IgnorePublicAcls:TrUe      RestrictPublicBuckets:trUE
     * 6、调用查询接口GetBucketPublicAccessBlock查看桶BPA状态
     * 7、调用删除接口DeleteBucketPublicAccessBlock删除BPA配置
     * 8、调用查询接口GetBucketPublicAccessBlock查看桶BPA状态
     *
     * 预期结果:
     * 1、创桶成功，200
     * 2、查询新创桶BPA状态，四个状态均为true
     * 3、设置状态false成功
     * 4、查询桶BPA状态，四个状态均为false
     * 5、设置状态true成功
     * 6、查询桶BPA状态，四个状态均为true
     * 7、删除BPA配置成功
     * 8、查询桶BPA状态，四个状态均为false
     *
     * @throws IOException
     */
    @Test
    public void test_SDK_bpa_001() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);

        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                    new PutBucketPublicAccessBlockRequest(bucketName, bucketPublicAccessBlock);
            HeaderResponse headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            GetBucketPublicAccessBlockRequest getBucketPublicAccessBlockRequest =
                    new GetBucketPublicAccessBlockRequest();
            getBucketPublicAccessBlockRequest.setBucketName(bucketName);
            GetBucketPublicAccessBlockResult getBucketPublicAccessBlockResult =
                    obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            bucketPublicAccessBlock.setBlockPublicACLs(false);
            bucketPublicAccessBlock.setIgnorePublicACLs(false);
            bucketPublicAccessBlock.setBlockPublicPolicy(false);
            bucketPublicAccessBlock.setRestrictPublicBuckets(false);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(
                    getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            DeleteBucketPublicAccessBlockRequest deleteBucketPublicAccessBlockRequest =
                    new DeleteBucketPublicAccessBlockRequest();
            deleteBucketPublicAccessBlockRequest.setBucketName(bucketName);
            headerResponse = obsClient.deleteBucketPublicAccessBlock(deleteBucketPublicAccessBlockRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(
                    getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    @Test
    public void test_SDK_bpa_002() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);

        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(false);
            bucketPublicAccessBlock.setIgnorePublicACLs(false);
            bucketPublicAccessBlock.setBlockPublicPolicy(false);
            bucketPublicAccessBlock.setRestrictPublicBuckets(false);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                    new PutBucketPublicAccessBlockRequest(bucketName, bucketPublicAccessBlock);
            HeaderResponse headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            GetBucketPublicAccessBlockRequest getBucketPublicAccessBlockRequest =
                    new GetBucketPublicAccessBlockRequest();
            getBucketPublicAccessBlockRequest.setBucketName(bucketName);
            GetBucketPublicAccessBlockResult getBucketPublicAccessBlockResult =
                    obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(
                    getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(true);
            putBucketPublicAccessBlockRequest.setBucketPublicAccessBlock(bucketPublicAccessBlock);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(
                    getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            putBucketPublicAccessBlockRequest.setBucketPublicAccessBlock(bucketPublicAccessBlock);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(
                    getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            putBucketPublicAccessBlockRequest.setBucketPublicAccessBlock(bucketPublicAccessBlock);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            putBucketPublicAccessBlockRequest.setBucketPublicAccessBlock(bucketPublicAccessBlock);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(
                    getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());

            DeleteBucketPublicAccessBlockRequest deleteBucketPublicAccessBlockRequest =
                    new DeleteBucketPublicAccessBlockRequest();
            deleteBucketPublicAccessBlockRequest.setBucketName(bucketName);
            headerResponse = obsClient.deleteBucketPublicAccessBlock(deleteBucketPublicAccessBlockRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            getBucketPublicAccessBlockResult = obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(
                    getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    @Test
    public void test_SDK_bpa_004() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                    new PutBucketPublicAccessBlockRequest(bucketName, bucketPublicAccessBlock);
            HeaderResponse headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            try {
                SetBucketAclRequest setBucketAclRequest =
                        new SetBucketAclRequest(bucketName, AccessControlList.REST_CANNED_PUBLIC_READ);
                obsClient.setBucketAcl(setBucketAclRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertEquals(403, e.getResponseCode());
            }

            try {
                String publicPolicy =
                        "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Allow\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"GetObject\",\"GetObjectVersion\"],\"Resource\":[\""
                                + bucketName
                                + "\",\""
                                + bucketName
                                + "/*\"]}]}";
                SetBucketPolicyRequest setBucketPolicyRequest = new SetBucketPolicyRequest(bucketName, publicPolicy);
                obsClient.setBucketPolicy(setBucketPolicyRequest);
                fail();
            } catch (ObsException e) {
                if (e.getResponseCode() != 403) {
                    TestTools.printObsException(e);
                }
                Assert.assertEquals(403, e.getResponseCode());
            }

            GetBucketPublicStatusRequest getBucketPublicStatusRequest = new GetBucketPublicStatusRequest();
            getBucketPublicStatusRequest.setBucketName(bucketName);
            GetBucketPublicStatusResult getBucketPublicStatusResult =
                    obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertFalse(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());

            String nonPublicPolicy =
                    "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Deny\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"PutObject\",\"GetObjectVersion\"],\"Resource\":[\""
                            + bucketName
                            + "\",\""
                            + bucketName
                            + "/*\"]}]}";
            SetBucketPolicyRequest setBucketPolicyRequest = new SetBucketPolicyRequest(bucketName, nonPublicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            getBucketPublicStatusResult = obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertFalse(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    @Test
    public void test_SDK_bpa_005() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;

            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                    new PutBucketPublicAccessBlockRequest(bucketName, bucketPublicAccessBlock);
            HeaderResponse headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            try {
                SetBucketAclRequest setBucketAclRequest =
                        new SetBucketAclRequest(bucketName, AccessControlList.REST_CANNED_PUBLIC_READ);
                obsClient.setBucketAcl(setBucketAclRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertEquals(403, e.getResponseCode());
            }

            try {
                String publicPolicy =
                        "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Allow\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"GetObject\",\"GetObjectVersion\"],\"Resource\":[\""
                                + bucketName
                                + "\",\""
                                + bucketName
                                + "/*\"]}]}";
                SetBucketPolicyRequest setBucketPolicyRequest = new SetBucketPolicyRequest(bucketName, publicPolicy);
                obsClient.setBucketPolicy(setBucketPolicyRequest);
                fail();
            } catch (ObsException e) {
                if (e.getResponseCode() != 403) {
                    TestTools.printObsException(e);
                }
                Assert.assertEquals(403, e.getResponseCode());
            }

            try {
                GetBucketPolicyPublicStatusRequest getBucketPolicyPublicStatusRequest =
                        new GetBucketPolicyPublicStatusRequest();
                getBucketPolicyPublicStatusRequest.setBucketName(bucketName);
                obsClient.getBucketPolicyPublicStatus(getBucketPolicyPublicStatusRequest);
                fail();
            } catch (ObsException e) {
                if (e.getResponseCode() != 404) {
                    TestTools.printObsException(e);
                }
                Assert.assertEquals(404, e.getResponseCode());
            }

            String nonPublicPolicy =
                    "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Deny\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"PutObject\",\"GetObjectVersion\"],\"Resource\":[\""
                            + bucketName
                            + "\",\""
                            + bucketName
                            + "/*\"]}]}";
            SetBucketPolicyRequest setBucketPolicyRequest = new SetBucketPolicyRequest(bucketName, nonPublicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            GetBucketPolicyPublicStatusRequest getBucketPolicyPublicStatusRequest =
                    new GetBucketPolicyPublicStatusRequest();
            getBucketPolicyPublicStatusRequest.setBucketName(bucketName);
            GetBucketPolicyPublicStatusResult getBucketPolicyPublicStatusResult =
                    obsClient.getBucketPolicyPublicStatus(getBucketPolicyPublicStatusRequest);
            Assert.assertFalse(getBucketPolicyPublicStatusResult.getPolicyStatus().getIsPublic());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    @Test
    public void test_SDK_bpa_006() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;

            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(false);
            bucketPublicAccessBlock.setIgnorePublicACLs(false);
            bucketPublicAccessBlock.setBlockPublicPolicy(false);
            bucketPublicAccessBlock.setRestrictPublicBuckets(false);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                    new PutBucketPublicAccessBlockRequest(bucketName, bucketPublicAccessBlock);
            HeaderResponse headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            // 2
            String publicPolicy =
                    "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Allow\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"GetObject\",\"GetObjectVersion\"],\"Resource\":[\""
                            + bucketName
                            + "\",\""
                            + bucketName
                            + "/*\"]}]}";
            SetBucketPolicyRequest setBucketPolicyRequest = new SetBucketPolicyRequest(bucketName, publicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // 3
            SetBucketAclRequest setBucketAclRequest =
                    new SetBucketAclRequest(bucketName, AccessControlList.REST_CANNED_PUBLIC_READ);
            headerResponse = obsClient.setBucketAcl(setBucketAclRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            // 4
            GetBucketPublicStatusRequest getBucketPublicStatusRequest = new GetBucketPublicStatusRequest();
            getBucketPublicStatusRequest.setBucketName(bucketName);
            GetBucketPublicStatusResult getBucketPublicStatusResult =
                    obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertTrue(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());

            // 5
            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            // 6
            try {
                obsClient.setBucketAcl(setBucketAclRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertEquals(403, e.getResponseCode());
            }

            // 7
            try {
                obsClient.setBucketPolicy(setBucketPolicyRequest);
                fail();
            } catch (ObsException e) {
                if (e.getResponseCode() != 403) {
                    TestTools.printObsException(e);
                }
                Assert.assertEquals(403, e.getResponseCode());
            }

            // 8
            getBucketPublicStatusResult = obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertTrue(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());

            // 9
            setBucketAclRequest.setAcl(AccessControlList.REST_CANNED_PRIVATE);
            headerResponse = obsClient.setBucketAcl(setBucketAclRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            // 10
            getBucketPublicStatusResult = obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertTrue(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());

            // 11
            String nonPublicPolicy =
                    "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Deny\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"PutObject\",\"GetObjectVersion\"],\"Resource\":[\""
                            + bucketName
                            + "\",\""
                            + bucketName
                            + "/*\"]}]}";
            setBucketPolicyRequest.setPolicy(nonPublicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // 12
            getBucketPublicStatusResult = obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertFalse(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());

            // 13
            BucketPolicyResponse bucketPolicyResponse = obsClient.getBucketPolicyV2(bucketName);
            Assert.assertEquals(bucketPolicyResponse.getStatusCode(), 200);

            // 14
            AccessControlList bucketAcl = obsClient.getBucketAcl(bucketName);
            Assert.assertEquals(bucketAcl.getStatusCode(), 200);
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    @Test
    public void test_SDK_bpa_007() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        try (ObsClient obsClient = TestTools.getPipelineEnvironment()) {
            assert obsClient != null;
            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(false);
            bucketPublicAccessBlock.setIgnorePublicACLs(false);
            bucketPublicAccessBlock.setBlockPublicPolicy(false);
            bucketPublicAccessBlock.setRestrictPublicBuckets(false);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                new PutBucketPublicAccessBlockRequest(bucketName, bucketPublicAccessBlock);
            HeaderResponse headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            // 2
            String publicPolicy =
                    "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Allow\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"GetObject\",\"GetObjectVersion\"],\"Resource\":[\""
                            + bucketName
                            + "\",\""
                            + bucketName
                            + "/*\"]}]}";
            SetBucketPolicyRequest setBucketPolicyRequest = new SetBucketPolicyRequest(bucketName, publicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // 2
            SetBucketAclRequest setBucketAclRequest =
                    new SetBucketAclRequest(bucketName, AccessControlList.REST_CANNED_PUBLIC_READ);
            headerResponse = obsClient.setBucketAcl(setBucketAclRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
            // 3
            GetBucketPublicStatusRequest getBucketPublicStatusRequest = new GetBucketPublicStatusRequest();
            getBucketPublicStatusRequest.setBucketName(bucketName);
            GetBucketPublicStatusResult getBucketPublicStatusResult =
                obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertTrue(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());
            // 3
            GetBucketPolicyPublicStatusRequest getBucketPolicyPublicStatusRequest = new GetBucketPolicyPublicStatusRequest();
            getBucketPolicyPublicStatusRequest.setBucketName(bucketName);
            GetBucketPolicyPublicStatusResult getBucketPolicyPublicStatusResult =
                obsClient.getBucketPolicyPublicStatus(getBucketPolicyPublicStatusRequest);
            Assert.assertTrue(getBucketPolicyPublicStatusResult.getPolicyStatus().getIsPublic());
            // 4
            String nonPublicPolicy =
                "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Deny\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"PutObject\",\"GetObjectVersion\"],\"Resource\":[\""
                    + bucketName
                    + "\",\""
                    + bucketName
                    + "/*\"]}]}";
            setBucketPolicyRequest.setPolicy(nonPublicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // 4
            headerResponse = obsClient.setBucketAcl(setBucketAclRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
            // 5
            getBucketPublicStatusResult =
                obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertTrue(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());
            // 5
            getBucketPolicyPublicStatusResult =
                obsClient.getBucketPolicyPublicStatus(getBucketPolicyPublicStatusRequest);
            Assert.assertFalse(getBucketPolicyPublicStatusResult.getPolicyStatus().getIsPublic());

            // 6
            setBucketPolicyRequest.setPolicy(nonPublicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());

            // 6
            setBucketAclRequest.setAcl(AccessControlList.REST_CANNED_PRIVATE);
            headerResponse = obsClient.setBucketAcl(setBucketAclRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());

            // 7
            getBucketPublicStatusResult =
                obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertFalse(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());
            // 7
            getBucketPolicyPublicStatusResult =
                obsClient.getBucketPolicyPublicStatus(getBucketPolicyPublicStatusRequest);
            Assert.assertFalse(getBucketPolicyPublicStatusResult.getPolicyStatus().getIsPublic());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    @Test
    public void test_SDK_bpa_008() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            + "-posix";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        try {
            assert obsClient != null;
            // 1
            CreateBucketRequest createBucketRequest = new CreateBucketRequest();
            createBucketRequest.setBucketName(bucketName);
            createBucketRequest.setBucketType(BucketTypeEnum.PFS);
            HeaderResponse headerResponse = obsClient.createBucket(createBucketRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
            // 1
            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                new PutBucketPublicAccessBlockRequest(bucketName, bucketPublicAccessBlock);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
            // 2
            GetBucketPublicAccessBlockRequest getBucketPublicAccessBlockRequest =
                new GetBucketPublicAccessBlockRequest();
            getBucketPublicAccessBlockRequest.setBucketName(bucketName);
            GetBucketPublicAccessBlockResult getBucketPublicAccessBlockResult =
                obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());
            // 3
            DeleteBucketPublicAccessBlockRequest deleteBucketPublicAccessBlockRequest =
                new DeleteBucketPublicAccessBlockRequest();
            deleteBucketPublicAccessBlockRequest.setBucketName(bucketName);
            headerResponse = obsClient.deleteBucketPublicAccessBlock(deleteBucketPublicAccessBlockRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());
            // 4
            getBucketPublicAccessBlockResult =
                obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertFalse(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());
            // 5
            bucketPublicAccessBlock.setBlockPublicACLs(true);
            bucketPublicAccessBlock.setIgnorePublicACLs(true);
            bucketPublicAccessBlock.setBlockPublicPolicy(true);
            bucketPublicAccessBlock.setRestrictPublicBuckets(true);
            putBucketPublicAccessBlockRequest.setBucketPublicAccessBlock(bucketPublicAccessBlock);
            headerResponse = obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            // 6
            getBucketPublicAccessBlockResult =
                obsClient.getBucketPublicAccessBlock(getBucketPublicAccessBlockRequest);
            Assert.assertEquals(200, getBucketPublicAccessBlockResult.getStatusCode());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicPolicy());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getBlockPublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getIgnorePublicACLs());
            Assert.assertTrue(getBucketPublicAccessBlockResult.getBucketPublicAccessBlock().getRestrictPublicBuckets());
            Assert.assertEquals(200, headerResponse.getStatusCode());
            // 7
            try {
                SetBucketAclRequest setBucketAclRequest =
                    new SetBucketAclRequest(bucketName, AccessControlList.REST_CANNED_PUBLIC_READ);
                obsClient.setBucketAcl(setBucketAclRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertEquals(403, e.getResponseCode());
            }

            // 8
            String publicPolicy =
                "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Allow\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"GetObject\",\"GetObjectVersion\"],\"Resource\":[\""
                    + bucketName
                    + "\",\""
                    + bucketName
                    + "/*\"]}]}";
            SetBucketPolicyRequest setBucketPolicyRequest = new SetBucketPolicyRequest(bucketName, publicPolicy);
            try {
                obsClient.setBucketPolicy(setBucketPolicyRequest);
                fail();
            } catch (ObsException e) {
                Assert.assertEquals(403, e.getResponseCode());
            }
            // 9
            GetBucketPublicStatusRequest getBucketPublicStatusRequest = new GetBucketPublicStatusRequest();
            getBucketPublicStatusRequest.setBucketName(bucketName);
            GetBucketPublicStatusResult getBucketPublicStatusResult =
                obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertFalse(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());
            // 10
            String nonPublicPolicy =
                "{\"Statement\":[{\"Sid\":\"44d8\",\"Effect\":\"Deny\",\"Principal\":{\"ID\":[\"*\"]},\"Action\":[\"ListBucket\",\"PutObject\",\"GetObjectVersion\"],\"Resource\":[\""
                    + bucketName
                    + "\",\""
                    + bucketName
                    + "/*\"]}]}";
            setBucketPolicyRequest.setPolicy(nonPublicPolicy);
            headerResponse = obsClient.setBucketPolicy(setBucketPolicyRequest);
            Assert.assertEquals(204, headerResponse.getStatusCode());
            // 9
            GetBucketPolicyPublicStatusRequest getBucketPolicyPublicStatusRequest = new GetBucketPolicyPublicStatusRequest();
            getBucketPolicyPublicStatusRequest.setBucketName(bucketName);
            GetBucketPolicyPublicStatusResult getBucketPolicyPublicStatusResult =
                obsClient.getBucketPolicyPublicStatus(getBucketPolicyPublicStatusRequest);
            Assert.assertFalse(getBucketPolicyPublicStatusResult.getPolicyStatus().getIsPublic());
            // 10
            SetBucketAclRequest setBucketAclRequest =
                new SetBucketAclRequest(bucketName, AccessControlList.REST_CANNED_PRIVATE);
            headerResponse = obsClient.setBucketAcl(setBucketAclRequest);
            Assert.assertEquals(200, headerResponse.getStatusCode());
            // 11
            getBucketPublicStatusResult =
                obsClient.getBucketPublicStatus(getBucketPublicStatusRequest);
            Assert.assertFalse(getBucketPublicStatusResult.getBucketPublicStatus().getIsPublic());
            // 11
            getBucketPolicyPublicStatusResult =
                obsClient.getBucketPolicyPublicStatus(getBucketPolicyPublicStatusRequest);
            Assert.assertFalse(getBucketPolicyPublicStatusResult.getPolicyStatus().getIsPublic());
        } catch (ObsException e) {
            TestTools.printObsException(e);
            throw e;
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        } finally {
            try {
                TestTools.delete_bucket(obsClient, bucketName);
            } catch (Throwable ignore) {

            }
        }
    }

    // 验证设置无效的BucketPublicAccessBlock时，抛出异常提示参数无效
    @Test
    public void test_SDK_bpa_009() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            + "-posix";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        try {
            assert obsClient != null;
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                new PutBucketPublicAccessBlockRequest(bucketName, null);
            obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            fail();
        } catch (ObsException e) {
            boolean buildXmlFailed = e.getMessage() != null && e.getMessage().contains("failed to build request XML");
            if (!buildXmlFailed) {
                TestTools.printObsException(e);
                throw e;
            } else {
                Assert.assertTrue(true);
            }
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }
    // 验证设置无效的BucketPublicAccessBlock时，抛出异常提示参数无效
    @Test
    public void test_SDK_bpa_0010() {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            + "-posix";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
        bucketPublicAccessBlock.setBlockPublicACLs(true);
        try {
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                new PutBucketPublicAccessBlockRequest(null, bucketPublicAccessBlock);
            obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            fail();
        } catch (IllegalArgumentException e) {
            boolean bucketNameInvalid = e.getMessage() != null && e.getMessage().contains("bucketName is null");
            if (!bucketNameInvalid) {
                throw e;
            } else {
                Assert.assertTrue(true);
            }
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }

        try {
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                new PutBucketPublicAccessBlockRequest("", bucketPublicAccessBlock);
            obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            fail();
        } catch (IllegalArgumentException e) {
            boolean bucketNameInvalid = e.getMessage() != null && e.getMessage().contains("bucketName is null");
            if (!bucketNameInvalid) {
                throw e;
            } else {
                Assert.assertTrue(true);
            }
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }

    // 验证设置无效的BucketPublicAccessBlock时，抛出异常提示参数无效
    @Test
    public void test_SDK_bpa_0011() throws IOException {
        if(TestCaseIgnore.needIgnore()) {
            return;
        }
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            + "-posix";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        try {
            assert obsClient != null;
            BucketPublicAccessBlock bucketPublicAccessBlock = new BucketPublicAccessBlock();
            bucketPublicAccessBlock.setBlockPublicACLs(false);
            bucketPublicAccessBlock.setRestrictPublicBuckets(false);
            bucketPublicAccessBlock.setIgnorePublicACLs(false);
            bucketPublicAccessBlock.setBlockPublicPolicy(false);
            PutBucketPublicAccessBlockRequest putBucketPublicAccessBlockRequest =
                new PutBucketPublicAccessBlockRequest(bucketName, null);
            obsClient.putBucketPublicAccessBlock(putBucketPublicAccessBlockRequest);
            fail();
        } catch (ObsException e) {
            boolean buildXmlFailed = e.getMessage() != null && e.getMessage().contains("failed to build request XML");
            if (!buildXmlFailed) {
                TestTools.printObsException(e);
                throw e;
            } else {
                Assert.assertTrue(true);
            }
        } catch (Exception e) {
            TestTools.printException(e);
            throw e;
        }
    }
}
