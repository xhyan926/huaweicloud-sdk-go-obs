/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.integrated_test;

import com.obs.services.SecretFlexibleObsClient;
import com.obs.services.model.AvailableZoneEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.ObsBucket;
import com.obs.test.TestTools;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.Locale;

public class SecretFlexibleBucketObsClientIT {

    @Rule
    public TestName testName = new TestName();

    @Test
    public void test_CreateBucketRequestWithAkSk() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        TestTools.TestPipelineAkSk env = TestTools.getPipelineAkSk();
        assert env != null;
        SecretFlexibleObsClient secretFlexibleObsClient =
            new SecretFlexibleObsClient(env.endPoint);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketName);
        createBucketRequest.setAvailableZone(AvailableZoneEnum.MULTI_AZ);
        try {
            ObsBucket obsBucket = secretFlexibleObsClient.createBucket(createBucketRequest, env.ak, env.sk);
            Assert.assertEquals(200, obsBucket.getStatusCode());
        } finally {
            try {
                secretFlexibleObsClient.deleteBucket(bucketName, env.ak, env.sk);
            } catch (Throwable ignore){}
        }
    }
    @Test
    public void test_CreateBucketRequestWithAkSkToken() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        TestTools.TestPipelineAkSk env = TestTools.getPipelineAkSkToken();
        assert env != null;
        SecretFlexibleObsClient secretFlexibleObsClient =
            new SecretFlexibleObsClient(env.endPoint);
        CreateBucketRequest createBucketRequest = new CreateBucketRequest(bucketName);
        createBucketRequest.setAvailableZone(AvailableZoneEnum.MULTI_AZ);
        try {
            ObsBucket obsBucket = secretFlexibleObsClient.createBucket(createBucketRequest, env.ak, env.sk, env.securityToken);
            Assert.assertEquals(200, obsBucket.getStatusCode());
        } finally {
            try {
                secretFlexibleObsClient.deleteBucket(bucketName, env.ak, env.sk, env.securityToken);
            } catch (Throwable ignore){}
        }

    }

}
