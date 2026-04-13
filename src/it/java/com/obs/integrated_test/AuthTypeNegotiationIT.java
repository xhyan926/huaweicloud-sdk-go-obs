package com.obs.integrated_test;
import com.obs.test.TestTools;

import com.obs.services.ObsClient;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.test.tools.PropertiesTools;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Test;

import java.io.File;
import java.io.IOException;
import java.io.StringWriter;
import java.util.ArrayList;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

public class AuthTypeNegotiationIT {
    private static final File file = new File("./app/src/test/resource/test_data.properties");
    private static final ArrayList<String> createdBuckets = new ArrayList<>();

    @BeforeClass
    public static void create_demo_bucket() throws IOException {
        String beforeBucket = PropertiesTools.getInstance(file).getProperties("beforeBucket");
        String location = PropertiesTools.getInstance(file).getProperties("environment.location");
        CreateBucketRequest request = new CreateBucketRequest();
        request.setBucketName(beforeBucket);
        request.setBucketType(BucketTypeEnum.OBJECT);
        request.setLocation(location);
        ObsClient obsClient = TestTools.getCustomPipelineEnvironment();

        request.setBucketName(beforeBucket + "-bucket001");
        HeaderResponse response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket001");

        request.setBucketName(beforeBucket + "-bucket002");
        response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket002");

        request.setBucketName(beforeBucket + "-bucket003");
        response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket003");

        request.setBucketName(beforeBucket + "-bucket004");
        response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(beforeBucket + "-bucket004");
    }

    @AfterClass
    public static void delete_created_buckets() {
        ObsClient obsClient = TestTools.getCustomPipelineEnvironment();
        for (String bucket : createdBuckets) {
            obsClient.deleteBucket(bucket);
        }
    }

    @Test
    public void test_auth_type_obs() throws IOException {
        // 初始化 log
        StringWriter writer = new StringWriter();
        TestTools.initLog(writer);
        String beforeBucket = PropertiesTools.getInstance(file).getProperties("beforeBucket");
        ObsClient obsClient = TestTools.getCustomPipelineEnvironment();
        obsClient.headBucket(beforeBucket + "-bucket001");
        assertTrue(writer.toString().contains("apiversion"));
        writer.getBuffer().setLength(0);
        obsClient.headBucket(beforeBucket + "-bucket002");
        assertTrue(writer.toString().contains("apiversion"));
        writer.getBuffer().setLength(0);
        obsClient.headBucket(beforeBucket + "-bucket003");
        assertTrue(writer.toString().contains("apiversion"));
        writer.getBuffer().setLength(0);
        obsClient.headBucket(beforeBucket + "-bucket004");
        assertTrue(writer.toString().contains("apiversion"));
        writer.getBuffer().setLength(0);
        obsClient.headBucket(beforeBucket + "-bucket001");
        assertTrue(writer.toString().contains("apiversion"));
        writer.getBuffer().setLength(0);
        obsClient.headBucket(beforeBucket + "-bucket003");
        assertTrue(!writer.toString().contains("apiversion"));
        writer.getBuffer().setLength(0);
    }
}
