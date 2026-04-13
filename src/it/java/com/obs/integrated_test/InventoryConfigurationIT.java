package com.obs.integrated_test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObsBucket;
import com.obs.services.model.inventory.DeleteInventoryConfigurationRequest;
import com.obs.services.model.inventory.GetInventoryConfigurationRequest;
import com.obs.services.model.inventory.GetInventoryConfigurationResult;
import com.obs.services.model.inventory.InventoryConfiguration;
import com.obs.services.model.inventory.ListInventoryConfigurationRequest;
import com.obs.services.model.inventory.ListInventoryConfigurationResult;
import com.obs.services.model.inventory.SetInventoryConfigurationRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

public class InventoryConfigurationIT {
    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Test
    public void tc_SetInventoryConfiguration_001() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String targetBucketName = bucketName + "-target";
        // 检查client
        assert obsClient != null;
        // 创建目标桶，并检查是否成功
        ObsBucket targetBucket = obsClient.createBucket(targetBucketName);
        assertEquals(200, targetBucket.getStatusCode());
        // 设置桶清单配置，并检查是否成功
        InventoryConfiguration configuration = new InventoryConfiguration();
        configuration.setDestinationBucket(targetBucketName);
        configuration.setConfigurationId("config001");
        configuration.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration.setFrequency(InventoryConfiguration.FrequencyOptions.DAILY);
        configuration.setEnabled(true);
        configuration.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.CURRENT);
        configuration.setInventoryPrefix("inventoryPrefix");
        configuration.setObjectPrefix("objectPrefix");
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.IS_MULTIPART_UPLOADED);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ETAG);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.REPLICATION_STATUS);
        SetInventoryConfigurationRequest request = new SetInventoryConfigurationRequest(bucketName, configuration);
        HeaderResponse response = obsClient.setInventoryConfiguration(request);
        assertEquals(200, response.getStatusCode());
        // 清除用例数据
        HeaderResponse response1 = obsClient.deleteBucket(targetBucketName);
        assertEquals(204, response1.getStatusCode());
    }

    @Test
    public void tc_GetInventoryConfiguration_002() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String targetBucketName = bucketName + "-target";
        // 检查client
        assert obsClient != null;
        // 创建目标桶，并检查是否成功
        ObsBucket targetBucket = obsClient.createBucket(targetBucketName);
        assertEquals(200, targetBucket.getStatusCode());
        // 设置桶清单配置，并检查是否成功
        String configId = "config001";
        InventoryConfiguration configuration = new InventoryConfiguration();
        configuration.setDestinationBucket(targetBucketName);
        configuration.setConfigurationId(configId);
        configuration.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration.setFrequency(InventoryConfiguration.FrequencyOptions.DAILY);
        configuration.setEnabled(true);
        configuration.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.CURRENT);
        configuration.setInventoryPrefix("inventoryPrefix");
        configuration.setObjectPrefix("objectPrefix");
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.IS_MULTIPART_UPLOADED);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.REPLICATION_STATUS);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ETAG);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.SIZE);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.LAST_MODIFIED_DATE);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.STORAGE_CLASS);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ENCRYPTION_STATUS);
        SetInventoryConfigurationRequest request = new SetInventoryConfigurationRequest(bucketName, configuration);
        HeaderResponse response = obsClient.setInventoryConfiguration(request);
        assertEquals(200, response.getStatusCode());
        // 获取桶清单
        GetInventoryConfigurationRequest getInventoryConfigurationRequest =
                new GetInventoryConfigurationRequest(bucketName, configId);
        GetInventoryConfigurationResult getInventoryConfigurationResult =
                obsClient.getInventoryConfiguration(getInventoryConfigurationRequest);
        assertEquals(200, getInventoryConfigurationResult.getStatusCode());
        InventoryConfiguration getConfiguration = getInventoryConfigurationResult.getInventoryConfiguration();
        assertEquals(configuration, getConfiguration);
        // 清除用例数据
        HeaderResponse response1 = obsClient.deleteBucket(targetBucketName);
        assertEquals(204, response1.getStatusCode());
    }

    @Test
    public void tc_ListInventoryConfiguration_003() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String targetBucketName = bucketName + "-target";
        // 检查client
        assert obsClient != null;
        // 创建目标桶，并检查是否成功
        ObsBucket targetBucket = obsClient.createBucket(targetBucketName);
        assertEquals(200, targetBucket.getStatusCode());
        // 设置桶清单配置1，并检查是否成功
        InventoryConfiguration configuration1 = new InventoryConfiguration();
        configuration1.setDestinationBucket(targetBucketName);
        configuration1.setConfigurationId("config001");
        configuration1.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration1.setFrequency(InventoryConfiguration.FrequencyOptions.DAILY);
        configuration1.setEnabled(true);
        configuration1.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.CURRENT);
        configuration1.setInventoryPrefix("inventoryPrefix001");
        configuration1.setObjectPrefix("objectPrefix001");
        configuration1.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.IS_MULTIPART_UPLOADED);
        configuration1.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.REPLICATION_STATUS);
        configuration1.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ETAG);

        SetInventoryConfigurationRequest request = new SetInventoryConfigurationRequest(bucketName, configuration1);
        HeaderResponse response = obsClient.setInventoryConfiguration(request);
        assertEquals(200, response.getStatusCode());
        // 设置桶清单配置2，并检查是否成功
        InventoryConfiguration configuration2 = new InventoryConfiguration();
        configuration2.setDestinationBucket(targetBucketName);
        configuration2.setConfigurationId("config002");
        configuration2.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration2.setFrequency(InventoryConfiguration.FrequencyOptions.WEEKLY);
        configuration2.setEnabled(false);
        configuration2.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.ALL);
        configuration2.setInventoryPrefix("inventoryPrefix002");
        configuration2.setObjectPrefix("objectPrefix002");
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.LAST_MODIFIED_DATE);
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ENCRYPTION_STATUS);
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.STORAGE_CLASS);
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.SIZE);

        request.setInventoryConfiguration(configuration2);
        response = obsClient.setInventoryConfiguration(request);
        assertEquals(200, response.getStatusCode());

        // 列举桶清单
        ListInventoryConfigurationRequest listInventoryConfigurationRequest =
                new ListInventoryConfigurationRequest(bucketName);
        ListInventoryConfigurationResult listInventoryConfigurationResult =
                obsClient.listInventoryConfiguration(listInventoryConfigurationRequest);
        assertEquals(200, listInventoryConfigurationResult.getStatusCode());
        List<InventoryConfiguration> getConfigurations =
                listInventoryConfigurationResult.getInventoryConfigurations();
        // 检查是否包含桶清单配置1、2
        assertTrue(getConfigurations.contains(configuration1));
        assertTrue(getConfigurations.contains(configuration2));
        // 清除用例数据
        HeaderResponse response1 = obsClient.deleteBucket(targetBucketName);
        assertEquals(204, response1.getStatusCode());
    }

    @Test
    public void tc_DeleteInventoryConfiguration_004() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String targetBucketName = bucketName + "-target";
        // 检查client
        assert obsClient != null;
        // 创建目标桶，并检查是否成功
        ObsBucket targetBucket = obsClient.createBucket(targetBucketName);
        assertEquals(200, targetBucket.getStatusCode());
        // 设置桶清单配置，并检查是否成功
        String configId = "config001";
        InventoryConfiguration configuration = new InventoryConfiguration();
        configuration.setDestinationBucket(targetBucketName);
        configuration.setConfigurationId(configId);
        configuration.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration.setFrequency(InventoryConfiguration.FrequencyOptions.DAILY);
        configuration.setEnabled(true);
        configuration.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.CURRENT);
        configuration.setInventoryPrefix("inventoryPrefix");
        configuration.setObjectPrefix("objectPrefix");
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.IS_MULTIPART_UPLOADED);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ETAG);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.REPLICATION_STATUS);
        SetInventoryConfigurationRequest request = new SetInventoryConfigurationRequest(bucketName, configuration);
        HeaderResponse response = obsClient.setInventoryConfiguration(request);
        assertEquals(200, response.getStatusCode());
        // 删除桶清单配置，并检查状态码
        DeleteInventoryConfigurationRequest deleteInventoryConfigurationRequest =
                new DeleteInventoryConfigurationRequest(bucketName, configId);
        HeaderResponse response0 = obsClient.deleteInventoryConfiguration(deleteInventoryConfigurationRequest);
        assertEquals(204, response0.getStatusCode());
        // 清除用例数据
        HeaderResponse response1 = obsClient.deleteBucket(targetBucketName);
        assertEquals(204, response1.getStatusCode());
    }

    @Test
    public void tc_SetInventoryConfiguration_005() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String targetBucketName = bucketName + "-target";
        // 检查client
        assert obsClient != null;
        // 创建目标桶，并检查是否成功
        ObsBucket targetBucket = obsClient.createBucket(targetBucketName);
        assertEquals(200, targetBucket.getStatusCode());
        // 设置桶清单配置1，并检查是否成功
        InventoryConfiguration configuration1 = new InventoryConfiguration();
        configuration1.setDestinationBucket(targetBucketName);
        configuration1.setConfigurationId("config001");
        configuration1.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration1.setFrequency(InventoryConfiguration.FrequencyOptions.DAILY);
        configuration1.setEnabled(true);
        configuration1.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.CURRENT);
        configuration1.setInventoryPrefix("inventoryPrefix001");
        configuration1.setObjectPrefix("objectPrefix");
        configuration1.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.IS_MULTIPART_UPLOADED);
        configuration1.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.REPLICATION_STATUS);
        configuration1.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ETAG);

        SetInventoryConfigurationRequest request = new SetInventoryConfigurationRequest(bucketName, configuration1);
        HeaderResponse response = obsClient.setInventoryConfiguration(request);
        assertEquals(200, response.getStatusCode());
        // 设置桶清单配置2，前缀包含清单配置1的前缀
        InventoryConfiguration configuration2 = new InventoryConfiguration();
        configuration2.setDestinationBucket(targetBucketName);
        configuration2.setConfigurationId("config002");
        configuration2.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration2.setFrequency(InventoryConfiguration.FrequencyOptions.WEEKLY);
        configuration2.setEnabled(false);
        configuration2.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.ALL);
        configuration2.setInventoryPrefix("inventoryPrefix002");
        configuration2.setObjectPrefix("objectPrefix002");
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.LAST_MODIFIED_DATE);
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ENCRYPTION_STATUS);
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.STORAGE_CLASS);
        configuration2.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.SIZE);

        try {
            request.setInventoryConfiguration(configuration2);
            obsClient.setInventoryConfiguration(request);
        } catch (ObsException e) {
            assertEquals(400, e.getResponseCode());
            assertEquals("PrefixExistInclusionRelationship", e.getErrorCode());
        }

        // 清除用例数据
        HeaderResponse response1 = obsClient.deleteBucket(targetBucketName);
        assertEquals(204, response1.getStatusCode());
    }

    @Test
    public void tc_SetInventoryConfiguration_006() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String targetBucketName = bucketName + "-target";
        // 检查client
        assert obsClient != null;
        // 创建目标桶，并检查是否成功
        ObsBucket targetBucket = obsClient.createBucket(targetBucketName);
        assertEquals(200, targetBucket.getStatusCode());
        // 设置桶清单配置，并检查是否成功
        InventoryConfiguration configuration = new InventoryConfiguration();
        configuration.setDestinationBucket(targetBucketName);
        configuration.setConfigurationId("config@!#$%^&*()+=");
        configuration.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration.setFrequency(InventoryConfiguration.FrequencyOptions.DAILY);
        configuration.setEnabled(true);
        configuration.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.CURRENT);
        configuration.setInventoryPrefix("inventoryPrefix");
        configuration.setObjectPrefix("objectPrefix");
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.IS_MULTIPART_UPLOADED);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ETAG);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.REPLICATION_STATUS);
        SetInventoryConfigurationRequest request = new SetInventoryConfigurationRequest(bucketName, configuration);
        try {
            obsClient.setInventoryConfiguration(request);
        } catch (ObsException e) {
            assertEquals(400, e.getResponseCode());
        }
        // 清除用例数据
        HeaderResponse response1 = obsClient.deleteBucket(targetBucketName);
        assertEquals(204, response1.getStatusCode());
    }

    @Test
    public void tc_GetInventoryConfiguration_007() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        // 检查client
        assert obsClient != null;
        // 获取桶清单
        String configId = "configId001";
        GetInventoryConfigurationRequest getInventoryConfigurationRequest =
                new GetInventoryConfigurationRequest(bucketName, configId);
        try {
            obsClient.getInventoryConfiguration(getInventoryConfigurationRequest);
        } catch (ObsException e) {
            assertEquals(404, e.getResponseCode());
        }
    }

    @Test
    public void tc_SetInventoryConfiguration_008() {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        String targetBucketName = bucketName + "-target";
        // 检查client
        assert obsClient != null;
        // 创建目标桶，并检查是否成功
        ObsBucket targetBucket = obsClient.createBucket(targetBucketName);
        assertEquals(200, targetBucket.getStatusCode());
        // 设置桶清单配置，并检查是否成功
        InventoryConfiguration configuration = new InventoryConfiguration();
        configuration.setDestinationBucket(targetBucketName);
        StringBuilder stringBuilder = new StringBuilder();
        int maxConfigIdLength = 64;
        for (int i = 0; i < maxConfigIdLength + 1; ++i) {
            stringBuilder.append("T");
        }
        configuration.setConfigurationId(stringBuilder.toString());
        configuration.setInventoryFormat(InventoryConfiguration.InventoryFormatOptions.CSV);
        configuration.setFrequency(InventoryConfiguration.FrequencyOptions.DAILY);
        configuration.setEnabled(true);
        configuration.setIncludedObjectVersions(InventoryConfiguration.IncludedObjectVersionsOptions.CURRENT);
        configuration.setInventoryPrefix("inventoryPrefix");
        configuration.setObjectPrefix("objectPrefix");
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.IS_MULTIPART_UPLOADED);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.ETAG);
        configuration.getOptionalFields().add(InventoryConfiguration.OptionalFieldOptions.REPLICATION_STATUS);
        SetInventoryConfigurationRequest request = new SetInventoryConfigurationRequest(bucketName, configuration);
        try {
            obsClient.setInventoryConfiguration(request);
        } catch (ObsException exception) {
            assertEquals(400, exception.getResponseCode());
        }
        // 清除用例数据
        HeaderResponse response1 = obsClient.deleteBucket(targetBucketName);
        assertEquals(204, response1.getStatusCode());
    }
}
