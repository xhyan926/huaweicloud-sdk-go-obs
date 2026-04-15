/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.integrated_test.objectlock;

import static org.junit.Assert.assertEquals;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.AuthTypeEnum;
import com.obs.services.model.BucketVersioningConfiguration;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ListVersionsRequest;
import com.obs.services.model.ListVersionsResult;
import com.obs.services.model.ObjectMetadata;
import com.obs.services.model.PutObjectRequest;
import com.obs.services.model.PutObjectResult;
import com.obs.services.model.VersionOrDeleteMarker;
import com.obs.services.model.VersioningStatusEnum;
import com.obs.services.model.objectlock.DefaultRetention;
import com.obs.services.model.objectlock.ObjectLockConfiguration;
import com.obs.services.model.objectlock.ObjectLockRule;
import com.obs.services.model.objectlock.ObjectRetention;
import com.obs.services.model.objectlock.SetObjectLockConfigurationRequest;
import com.obs.services.model.objectlock.SetObjectRetentionRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PropertiesTools;

import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.util.Arrays;
import java.util.Collection;
import java.util.Locale;

@RunWith(Parameterized.class)
public class ObjectRetentionIT {
    @Rule
    public TestName testName = new TestName();

    @Parameterized.Parameter()
    public String authTypeName;

    @Parameterized.Parameter(1)
    public AuthTypeEnum authType;

    @Parameterized.Parameters(name = "{0}")
    public static Collection<Object[]> authTypeData() {
        return Arrays.asList(new Object[][] {
            {"OBS", AuthTypeEnum.OBS},
            {"V2", AuthTypeEnum.V2}
        });
    }

    private ObsClient obsClient;
    private String bucketName;

    private static File getConfigFile() {
        String env = System.getProperty("test.env", "");
        String fileName = env.isEmpty()
            ? "test_data.properties"
            : "test_data_" + env + ".properties";
        return new File("./app/src/test/resource/" + fileName);
    }

    @Before
    public void setUp() throws Exception {
        if (authType == AuthTypeEnum.OBS) {
            obsClient = TestTools.getPipelineEnvironment();
        } else {
            obsClient = TestTools.getPipelineEnvironment_V2();
        }
        Assert.assertNotNull("ObsClient should not be null", obsClient);
        bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            .replace("[", "").replace("]", "");

        String location = PropertiesTools.getInstance(getConfigFile()).getProperties("environment.location");
        boolean isPosix = Boolean.parseBoolean(
            PropertiesTools.getInstance(getConfigFile()).getProperties("isPosix"));
        HeaderResponse createResponse = TestTools.createBucket(obsClient, bucketName, location, isPosix);
        Assert.assertEquals("Create bucket should succeed", 200, createResponse.getStatusCode());
    }

    @After
    public void tearDown() {
        if (obsClient != null && bucketName != null) {
            // WORM保护的对象不可删除，清理时需先列举对象再尝试删除
            // 忽略因WORM保护导致的删除失败
            try {
                ListVersionsRequest listRequest = new ListVersionsRequest(bucketName);
                listRequest.setMaxKeys(1000);
                ListVersionsResult versionsResult = obsClient.listVersions(listRequest);
                for (VersionOrDeleteMarker version : versionsResult.getVersions()) {
                    try {
                        obsClient.deleteObject(bucketName, version.getObjectKey(), version.getVersionId());
                    } catch (ObsException ignored) {
                        // WORM保护的对象删除失败，忽略
                    }
                }
            } catch (ObsException ignored) {
            }
            try {
                obsClient.deleteBucket(bucketName);
            } catch (ObsException ignored) {
                // WORM保护的对象仍存在时桶不可删除，忽略
            }
        }
    }

    /**
     * 开启桶级WORM、上传对象、设置对象级WORM保护策略并验证
     */
    @Test
    public void test_SDK_objectretention_001() {
        String objectKey = "retention-test-object";

        // 开启多版本
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // 开启桶级WORM
        DefaultRetention defaultRetention = new DefaultRetention("COMPLIANCE", 1, null);
        ObjectLockRule rule = new ObjectLockRule(defaultRetention);
        ObjectLockConfiguration lockConfig = new ObjectLockConfiguration("Enabled", rule);
        SetObjectLockConfigurationRequest setLockRequest =
            new SetObjectLockConfigurationRequest(bucketName, lockConfig);
        HeaderResponse lockResponse = obsClient.setObjectLockConfiguration(setLockRequest);
        Assert.assertEquals(200, lockResponse.getStatusCode());

        // 上传对象
        String content = "test content for object retention";
        PutObjectRequest putRequest = new PutObjectRequest();
        putRequest.setBucketName(bucketName);
        putRequest.setObjectKey(objectKey);
        putRequest.setInput(new ByteArrayInputStream(content.getBytes()));
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setContentLength((long) content.length());
        putRequest.setMetadata(metadata);
        PutObjectResult putResult = obsClient.putObject(putRequest);
        Assert.assertNotNull(putResult);

        // 设置对象级WORM保护策略
        long retainUntilDate = System.currentTimeMillis() + 30L * 24 * 60 * 60 * 1000;
        ObjectRetention retention = new ObjectRetention("COMPLIANCE", retainUntilDate);
        SetObjectRetentionRequest setRetentionRequest =
            new SetObjectRetentionRequest(bucketName, objectKey, retention);
        HeaderResponse retentionResponse = obsClient.setObjectRetention(setRetentionRequest);
        Assert.assertEquals(200, retentionResponse.getStatusCode());
    }

    /**
     * 设置对象级WORM保护策略后，验证可以延长保护期限
     */
    @Test
    public void test_SDK_objectretention_002() {
        String objectKey = "retention-extend-test-object";

        // 开启多版本
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // 开启桶级WORM（不设置默认保护期限）
        ObjectLockConfiguration lockConfig = new ObjectLockConfiguration("Enabled", null);
        SetObjectLockConfigurationRequest setLockRequest =
            new SetObjectLockConfigurationRequest(bucketName, lockConfig);
        HeaderResponse lockResponse = obsClient.setObjectLockConfiguration(setLockRequest);
        Assert.assertEquals(200, lockResponse.getStatusCode());

        // 上传对象
        String content = "test content for extend retention";
        PutObjectRequest putRequest = new PutObjectRequest();
        putRequest.setBucketName(bucketName);
        putRequest.setObjectKey(objectKey);
        putRequest.setInput(new ByteArrayInputStream(content.getBytes()));
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setContentLength((long) content.length());
        putRequest.setMetadata(metadata);
        obsClient.putObject(putRequest);

        // 首次设置对象级WORM保护策略
        long initialRetainDate = System.currentTimeMillis() + 10L * 24 * 60 * 60 * 1000;
        ObjectRetention retention = new ObjectRetention("COMPLIANCE", initialRetainDate);
        SetObjectRetentionRequest setRetentionRequest =
            new SetObjectRetentionRequest(bucketName, objectKey, retention);
        HeaderResponse response = obsClient.setObjectRetention(setRetentionRequest);
        Assert.assertEquals(200, response.getStatusCode());

        // 延长保护期限（仅允许延长，不允许缩短）
        long extendedRetainDate = System.currentTimeMillis() + 60L * 24 * 60 * 60 * 1000;
        retention = new ObjectRetention("COMPLIANCE", extendedRetainDate);
        setRetentionRequest = new SetObjectRetentionRequest(bucketName, objectKey, retention);
        response = obsClient.setObjectRetention(setRetentionRequest);
        Assert.assertEquals(200, response.getStatusCode());
    }

    /**
     * 验证未开启WORM的桶设置对象级保护策略时返回错误
     */
    @Test
    public void test_SDK_objectretention_003() {
        String objectKey = "retention-no-worm-object";

        // 开启多版本（但不开启WORM）
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // 上传对象
        String content = "test content no worm";
        PutObjectRequest putRequest = new PutObjectRequest();
        putRequest.setBucketName(bucketName);
        putRequest.setObjectKey(objectKey);
        putRequest.setInput(new ByteArrayInputStream(content.getBytes()));
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setContentLength((long) content.length());
        putRequest.setMetadata(metadata);
        obsClient.putObject(putRequest);

        // 在未开启WORM的桶上设置对象级保护策略，应返回错误
        long retainUntilDate = System.currentTimeMillis() + 30L * 24 * 60 * 60 * 1000;
        ObjectRetention retention = new ObjectRetention("COMPLIANCE", retainUntilDate);
        SetObjectRetentionRequest setRetentionRequest =
            new SetObjectRetentionRequest(bucketName, objectKey, retention);

        try {
            obsClient.setObjectRetention(setRetentionRequest);
            Assert.fail("Expected ObsException for bucket without WORM enabled");
        } catch (ObsException e) {
            Assert.assertEquals(400, e.getResponseCode());
        }
    }

    /**
     * 验证缩短保护期限时返回错误（保护期限仅允许延长，不允许缩短）
     */
    @Test
    public void test_SDK_objectretention_004() {
        String objectKey = "retention-shorten-object";

        // 开启多版本
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // 开启桶级WORM
        ObjectLockConfiguration lockConfig = new ObjectLockConfiguration("Enabled", null);
        SetObjectLockConfigurationRequest setLockRequest =
            new SetObjectLockConfigurationRequest(bucketName, lockConfig);
        HeaderResponse lockResponse = obsClient.setObjectLockConfiguration(setLockRequest);
        Assert.assertEquals(200, lockResponse.getStatusCode());

        // 上传对象
        String content = "test content shorten retention";
        PutObjectRequest putRequest = new PutObjectRequest();
        putRequest.setBucketName(bucketName);
        putRequest.setObjectKey(objectKey);
        putRequest.setInput(new ByteArrayInputStream(content.getBytes()));
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setContentLength((long) content.length());
        putRequest.setMetadata(metadata);
        obsClient.putObject(putRequest);

        // 先设置一个较长的保护期限
        long longRetainDate = System.currentTimeMillis() + 60L * 24 * 60 * 60 * 1000;
        ObjectRetention retention = new ObjectRetention("COMPLIANCE", longRetainDate);
        SetObjectRetentionRequest setRetentionRequest =
            new SetObjectRetentionRequest(bucketName, objectKey, retention);
        HeaderResponse response = obsClient.setObjectRetention(setRetentionRequest);
        Assert.assertEquals(200, response.getStatusCode());

        // 尝试缩短保护期限，应返回错误
        long shortRetainDate = System.currentTimeMillis() + 1L * 24 * 60 * 60 * 1000;
        retention = new ObjectRetention("COMPLIANCE", shortRetainDate);
        setRetentionRequest = new SetObjectRetentionRequest(bucketName, objectKey, retention);

        try {
            obsClient.setObjectRetention(setRetentionRequest);
            Assert.fail("Expected ObsException when shortening retention period");
        } catch (ObsException e) {
            Assert.assertEquals(400, e.getResponseCode());
        }
    }

    /**
     * 验证PUT新版本可以成功（多版本场景下PUT创建新版本而非覆盖），
     * 但受WORM保护的原始版本不可删除
     */
    @Test
    public void test_SDK_objectretention_005() {
        String objectKey = "retention-overwrite-object";

        // 开启多版本
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // 开启桶级WORM
        ObjectLockConfiguration lockConfig = new ObjectLockConfiguration("Enabled", null);
        SetObjectLockConfigurationRequest setLockRequest =
            new SetObjectLockConfigurationRequest(bucketName, lockConfig);
        HeaderResponse lockResponse = obsClient.setObjectLockConfiguration(setLockRequest);
        Assert.assertEquals(200, lockResponse.getStatusCode());

        // 上传对象v1
        String content = "original content";
        PutObjectRequest putRequest = new PutObjectRequest();
        putRequest.setBucketName(bucketName);
        putRequest.setObjectKey(objectKey);
        putRequest.setInput(new ByteArrayInputStream(content.getBytes()));
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setContentLength((long) content.length());
        putRequest.setMetadata(metadata);
        PutObjectResult putResult = obsClient.putObject(putRequest);
        String versionId = putResult.getVersionId();

        // 为v1设置对象级WORM保护策略
        long retainUntilDate = System.currentTimeMillis() + 30L * 24 * 60 * 60 * 1000;
        ObjectRetention retention = new ObjectRetention("COMPLIANCE", retainUntilDate);
        SetObjectRetentionRequest setRetentionRequest =
            new SetObjectRetentionRequest(bucketName, objectKey, retention, versionId);
        HeaderResponse retentionResponse = obsClient.setObjectRetention(setRetentionRequest);
        Assert.assertEquals(200, retentionResponse.getStatusCode());

        // 多版本场景下PUT同一key会创建新版本，不报错
        String newContent = "new version content";
        PutObjectRequest newPutRequest = new PutObjectRequest();
        newPutRequest.setBucketName(bucketName);
        newPutRequest.setObjectKey(objectKey);
        newPutRequest.setInput(new ByteArrayInputStream(newContent.getBytes()));
        ObjectMetadata newMetadata = new ObjectMetadata();
        newMetadata.setContentLength((long) newContent.length());
        newPutRequest.setMetadata(newMetadata);
        PutObjectResult newPutResult = obsClient.putObject(newPutRequest);
        Assert.assertNotNull(newPutResult);

        // 但受WORM保护的v1版本不可删除
        try {
            obsClient.deleteObject(bucketName, objectKey, versionId);
            Assert.fail("Expected ObsException when deleting WORM protected version");
        } catch (ObsException e) {
            Assert.assertEquals(403, e.getResponseCode());
        }
    }

    /**
     * 验证删除受WORM保护的对象时返回错误
     */
    @Test
    public void test_SDK_objectretention_006() {
        String objectKey = "retention-delete-object";

        // 开启多版本
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // 开启桶级WORM
        ObjectLockConfiguration lockConfig = new ObjectLockConfiguration("Enabled", null);
        SetObjectLockConfigurationRequest setLockRequest =
            new SetObjectLockConfigurationRequest(bucketName, lockConfig);
        HeaderResponse lockResponse = obsClient.setObjectLockConfiguration(setLockRequest);
        Assert.assertEquals(200, lockResponse.getStatusCode());

        // 上传对象
        String content = "content to be deleted";
        PutObjectRequest putRequest = new PutObjectRequest();
        putRequest.setBucketName(bucketName);
        putRequest.setObjectKey(objectKey);
        putRequest.setInput(new ByteArrayInputStream(content.getBytes()));
        ObjectMetadata metadata = new ObjectMetadata();
        metadata.setContentLength((long) content.length());
        putRequest.setMetadata(metadata);
        PutObjectResult putResult = obsClient.putObject(putRequest);
        String versionId = putResult.getVersionId();

        // 设置对象级WORM保护策略
        long retainUntilDate = System.currentTimeMillis() + 30L * 24 * 60 * 60 * 1000;
        ObjectRetention retention = new ObjectRetention("COMPLIANCE", retainUntilDate);
        SetObjectRetentionRequest setRetentionRequest =
            new SetObjectRetentionRequest(bucketName, objectKey, retention, versionId);
        HeaderResponse retentionResponse = obsClient.setObjectRetention(setRetentionRequest);
        Assert.assertEquals(200, retentionResponse.getStatusCode());

        // 尝试删除受保护的对象（指定版本号），应返回错误
        try {
            obsClient.deleteObject(bucketName, objectKey, versionId);
            Assert.fail("Expected ObsException when deleting WORM protected object");
        } catch (ObsException e) {
            Assert.assertEquals(403, e.getResponseCode());
        }
    }
}
