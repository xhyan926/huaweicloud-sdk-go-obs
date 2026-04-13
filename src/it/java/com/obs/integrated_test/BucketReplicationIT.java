/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.DeleteDataEnum;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.HistoricalObjectReplicationEnum;
import com.obs.services.model.ReplicationConfiguration;
import com.obs.services.model.RuleStatusEnum;
import com.obs.services.model.SetBucketReplicationRequest;
import com.obs.services.model.StorageClassEnum;
import com.obs.test.tools.PropertiesTools;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Test;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;

import static org.junit.Assert.assertEquals;

public class BucketReplicationIT {

    private static ObsClient obsClient1;
    private static ObsClient obsClient2;
    private static final File file = new File("./app/src/test/resource/test_data.properties");

    @BeforeClass
    public static void initializeClients() throws IOException, ObsException {
        // Initialize OBS clients for two different regions
        String ak = PropertiesTools.getInstance(file).getProperties("environment.crr.ak");
        String sk = PropertiesTools.getInstance(file).getProperties("environment.crr.sk");
        String securityToken = "";

        ObsConfiguration obsConfiguration1 = new ObsConfiguration();
        obsConfiguration1.setEndPoint(PropertiesTools.getInstance(file).getProperties("environment.crr.endpoint"));
        obsClient1 = new ObsClient(ak, sk, securityToken, obsConfiguration1);

        ObsConfiguration obsConfiguration2 = new ObsConfiguration();
        obsConfiguration2.setEndPoint(PropertiesTools.getInstance(file).getProperties("environment.crr.endpoint2"));
        obsClient2 = new ObsClient(ak, sk, securityToken, obsConfiguration2);

        // Create test buckets in different regions
        createBucket("example-bucket-test-1", "R1", obsClient1);
        createBucket("example-bucket-test-2", "R2", obsClient2);
    }

    private static void createBucket(String bucketName, String location, ObsClient obsClient) throws ObsException {
        CreateBucketRequest request = new CreateBucketRequest(bucketName);
        request.setLocation(location);
        obsClient.createBucket(request);
    }

    @AfterClass
    public static void deleteTestBuckets() throws ObsException {
        // Cleanup: Delete test buckets after tests
        deleteBucket("example-bucket-test-1", obsClient1);
        deleteBucket("example-bucket-test-2", obsClient2);
    }

    private static void deleteBucket(String bucketName, ObsClient obsClient) throws ObsException {
        obsClient.deleteBucket(bucketName);
    }

    @Test
    public void tc_set_BucketReplication_DeleteDataEnum() throws ObsException {
        // Step 1: Define the source and destination buckets
        String sourceBucket = "example-bucket-test-1";
        String destinationBucket = "arn:aws:s3:::example-bucket-test-2";

        // Step 2: Create and configure replication rule
        ReplicationConfiguration replicationConfiguration = new ReplicationConfiguration();
        replicationConfiguration.setAgency("testAgency");
        ArrayList<ReplicationConfiguration.Rule> rules = new ArrayList<>();

        ReplicationConfiguration.Rule rule = new ReplicationConfiguration.Rule();
        rule.setId("rule-delete-data");
        rule.setStatus(RuleStatusEnum.ENABLED);
        rule.setPrefix("key-prefix");
        rule.setHistoricalObjectReplication(HistoricalObjectReplicationEnum.ENABLED);

        ReplicationConfiguration.Destination destination = new ReplicationConfiguration.Destination();
        destination.setBucket(destinationBucket);
        destination.setObjectStorageClass(StorageClassEnum.STANDARD);
        destination.setDeleteData(DeleteDataEnum.ENABLED);

        rule.setDestination(destination);
        rules.add(rule);
        replicationConfiguration.setRules(rules);

        SetBucketReplicationRequest request = new SetBucketReplicationRequest(sourceBucket, replicationConfiguration);

        // Step 3: Set replication with DeleteDataEnum ENABLED
        HeaderResponse response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 4: Retrieve and verify DeleteDataEnum is ENABLED
        ReplicationConfiguration replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(DeleteDataEnum.ENABLED, replicationConfig.getRules().get(0).getDestination().getDeleteData());

        // Step 5: Update replication rule to set DeleteDataEnum to DISABLED
        request.getReplicationConfiguration().getRules().get(0).getDestination().setDeleteData(DeleteDataEnum.DISABLED);
        response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 6: Retrieve and verify DeleteDataEnum is DISABLED
        replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(DeleteDataEnum.DISABLED, replicationConfig.getRules().get(0).getDestination().getDeleteData());
    }

    @Test
    public void tc_set_BucketReplication_HistoricalObjectReplicationEnum() throws ObsException {
        // Step 1: Define source and destination buckets
        String sourceBucket = "example-bucket-test-1";
        String destinationBucket = "arn:aws:s3:::example-bucket-test-2";

        // Step 2: Create and configure replication rule with HistoricalObjectReplicationEnum ENABLED
        ReplicationConfiguration replicationConfiguration = new ReplicationConfiguration();
        replicationConfiguration.setAgency("testAgency");
        ArrayList<ReplicationConfiguration.Rule> rules = new ArrayList<>();

        ReplicationConfiguration.Rule rule = new ReplicationConfiguration.Rule();
        rule.setId("rule-historical-replication");
        rule.setStatus(RuleStatusEnum.ENABLED);
        rule.setPrefix("key-prefix");
        rule.setHistoricalObjectReplication(HistoricalObjectReplicationEnum.ENABLED);

        ReplicationConfiguration.Destination destination = new ReplicationConfiguration.Destination();
        destination.setBucket(destinationBucket);
        destination.setObjectStorageClass(StorageClassEnum.STANDARD);
        rule.setDestination(destination);
        rules.add(rule);
        replicationConfiguration.setRules(rules);

        SetBucketReplicationRequest request = new SetBucketReplicationRequest(sourceBucket, replicationConfiguration);

        // Step 3: Set replication with HistoricalObjectReplicationEnum ENABLED
        HeaderResponse response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 4: Retrieve and verify HistoricalObjectReplicationEnum is ENABLED
        ReplicationConfiguration replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(HistoricalObjectReplicationEnum.ENABLED, replicationConfig.getRules().get(0).getHistoricalObjectReplication());

        // Step 5: Update replication rule to set HistoricalObjectReplicationEnum to DISABLED
        request.getReplicationConfiguration().getRules().get(0).setHistoricalObjectReplication(HistoricalObjectReplicationEnum.DISABLED);
        response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 6: Retrieve and verify HistoricalObjectReplicationEnum is DISABLED
        replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(HistoricalObjectReplicationEnum.DISABLED, replicationConfig.getRules().get(0).getHistoricalObjectReplication());
    }

    @Test
    public void tc_set_BucketReplication_DeleteDataEnum_HistoricalObjectReplicationEnum() throws ObsException {
        // Step 1: Define the source and destination buckets
        String sourceBucket = "example-bucket-test-1";
        String destinationBucket = "arn:aws:s3:::example-bucket-test-2";

        // Step 2: Create and configure replication rule with DeleteDataEnum ENABLED and HistoricalObjectReplicationEnum DISABLED
        ReplicationConfiguration replicationConfiguration = new ReplicationConfiguration();
        replicationConfiguration.setAgency("testAgency");
        ArrayList<ReplicationConfiguration.Rule> rules = new ArrayList<>();

        ReplicationConfiguration.Rule rule = new ReplicationConfiguration.Rule();
        rule.setId("rule-delete-historical");
        rule.setStatus(RuleStatusEnum.ENABLED);
        rule.setPrefix("key-prefix");
        rule.setHistoricalObjectReplication(HistoricalObjectReplicationEnum.DISABLED);

        ReplicationConfiguration.Destination destination = new ReplicationConfiguration.Destination();
        destination.setBucket(destinationBucket);
        destination.setObjectStorageClass(StorageClassEnum.STANDARD);
        destination.setDeleteData(DeleteDataEnum.ENABLED);

        rule.setDestination(destination);
        rules.add(rule);
        replicationConfiguration.setRules(rules);

        SetBucketReplicationRequest request = new SetBucketReplicationRequest(sourceBucket, replicationConfiguration);

        // Step 3: Set replication rule
        HeaderResponse response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 4: Verify the settings
        ReplicationConfiguration replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(DeleteDataEnum.ENABLED, replicationConfig.getRules().get(0).getDestination().getDeleteData());
        assertEquals(HistoricalObjectReplicationEnum.DISABLED, replicationConfig.getRules().get(0).getHistoricalObjectReplication());

        // Step 5: Modify the replication rule to Disable DeleteDataEnum and Enable HistoricalObjectReplicationEnum
        rule.getDestination().setDeleteData(DeleteDataEnum.DISABLED);
        rule.setHistoricalObjectReplication(HistoricalObjectReplicationEnum.ENABLED);
        response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 6: Verify the settings
        replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(DeleteDataEnum.DISABLED, replicationConfig.getRules().get(0).getDestination().getDeleteData());
        assertEquals(HistoricalObjectReplicationEnum.ENABLED, replicationConfig.getRules().get(0).getHistoricalObjectReplication());

        // Step 7: Modify the replication rule to Enable both DeleteDataEnum and HistoricalObjectReplicationEnum
        rule.getDestination().setDeleteData(DeleteDataEnum.ENABLED);
        rule.setHistoricalObjectReplication(HistoricalObjectReplicationEnum.ENABLED);
        response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 8: Verify the settings
        replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(DeleteDataEnum.ENABLED, replicationConfig.getRules().get(0).getDestination().getDeleteData());
        assertEquals(HistoricalObjectReplicationEnum.ENABLED, replicationConfig.getRules().get(0).getHistoricalObjectReplication());

        // Step 9: Modify the replication rule to Disable both DeleteDataEnum and HistoricalObjectReplicationEnum
        rule.getDestination().setDeleteData(DeleteDataEnum.DISABLED);
        rule.setHistoricalObjectReplication(HistoricalObjectReplicationEnum.DISABLED);
        response = obsClient1.setBucketReplication(request.getBucketName(), request.getReplicationConfiguration());
        assertEquals(200, response.getStatusCode());

        // Step 10: Verify the settings
        replicationConfig = obsClient1.getBucketReplication(request.getBucketName());
        assertEquals(DeleteDataEnum.DISABLED, replicationConfig.getRules().get(0).getDestination().getDeleteData());
        assertEquals(HistoricalObjectReplicationEnum.DISABLED, replicationConfig.getRules().get(0).getHistoricalObjectReplication());
    }
}
