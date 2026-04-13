/**
 * Copyright 2019 Huawei Technologies Co.,Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package com.obs.integrated_test;
import com.obs.services.model.AccessControlList;
import com.obs.services.model.DeleteSnapshotRequest;
import com.obs.services.model.GetSnapshotListRequest;
import com.obs.services.model.GetSnapshotListResponse;
import com.obs.services.model.ListObjectsRequest;
import com.obs.services.model.ObsObject;
import com.obs.services.model.RenameSnapshotRequest;
import com.obs.services.model.RenameSnapshotResponse;
import com.obs.services.model.Snapshot;
import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.model.BucketTypeEnum;
import com.obs.services.model.CreateSnapshotRequest;
import com.obs.services.model.CreateSnapshotResponse;
import com.obs.services.model.CreateBucketRequest;
import com.obs.services.model.GetSnapshottableDirListRequest;
import com.obs.services.model.GetSnapshottableDirListResult;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.ObjectListing;
import com.obs.services.model.SetDisallowSnapshotRequest;
import com.obs.services.model.SetSnapshotAllowRequest;
import com.obs.services.model.SnapshottableDir;
import com.obs.services.model.fs.DropFolderRequest;
import com.obs.services.model.fs.NewFolderRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import com.obs.test.tools.PropertiesTools;
import org.junit.After;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import static com.obs.services.internal.ObsConstraint.SNAPSHOT_MAX_KEYS;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.fail;

public class BucketSnapshotIT {
    private static final File file = new File("./app/src/test/resource/test_data.properties");
    private static final String NON_EXISTENT_SNAPSHOT_C = "snapshot-c";
    private static final String NEW_SNAPSHOT_NAME = "snapshot-d";
    private static final String DUMMY_OBJECT_KEY = "test-object";
    private static final String NON_EXISTENT_SNAPSHOT_D = "snapshot-e";
    protected static String DIRECTORY_B = "directory-b/";
    protected static String SNAPSHOT_C = "snapshot-c";
    protected static String SNAPSHOT_D = "snapshot-d";
    protected ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();

    protected static String snapshotBucket;
    protected static String objectBucket;
    protected static String location;
    protected static String testDir = "my-dir-a/";

    private static final ArrayList<String> createdBuckets = new ArrayList<>();

    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();
    @BeforeClass
    public static void create_demo_bucket() throws IOException {
        snapshotBucket = "posix-snapshot-2";
        location = PropertiesTools.getInstance(file).getProperties("environment.snapshot.location");
        CreateBucketRequest request = new CreateBucketRequest();
        request.setBucketName(snapshotBucket);
        request.setBucketType(BucketTypeEnum.PFS);
        request.setLocation(location);
        request.setAcl(AccessControlList.REST_CANNED_PUBLIC_READ_WRITE);
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();
        assertNotNull(obsClient);
        obsClient.listBuckets();
        HeaderResponse response = obsClient.createBucket(request);
        assertEquals(200, response.getStatusCode());
        createdBuckets.add(snapshotBucket);

        CreateBucketRequest request2 = new CreateBucketRequest();
        objectBucket = "object-bucket-snapshot";
        request2.setBucketName(objectBucket);
        request2.setBucketType(BucketTypeEnum.OBJECT);
        request2.setLocation(location);
        HeaderResponse response2 = obsClient.createBucket(request2);
        assertEquals(200, response2.getStatusCode());
        createdBuckets.add(objectBucket);
    }

    @AfterClass
    public static void delete_created_buckets() throws IOException {
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();

        for (String bucket : createdBuckets) {
            String objectMarker = null;
            do{
                ListObjectsRequest listObjectsRequest = new ListObjectsRequest(bucket, null, objectMarker, null, SNAPSHOT_MAX_KEYS);
                ObjectListing objectListing = obsClient.listObjects(listObjectsRequest);
                objectMarker = objectListing.getNextMarker();
                for (ObsObject object : objectListing.getObjects()){
                    DropFolderRequest dropFolderRequest = new DropFolderRequest(bucket,object.getObjectKey());
                    obsClient.dropFolder(dropFolderRequest);
                }
            } while (objectMarker != null);

            obsClient.deleteBucket(bucket);
        }
    }

    @After
    public void delete_snapshots_folders() throws IOException {
        ObsClient obsClient = TestTools.getPipelineForSnapshotEnvironment();
        String bucket = snapshotBucket;
        String objectMarker = null;
        do {
            ListObjectsRequest listObjectsRequest = new ListObjectsRequest(bucket, null, objectMarker, null, SNAPSHOT_MAX_KEYS);
            ObjectListing objectListing = obsClient.listObjects(listObjectsRequest);
            objectMarker = objectListing.getNextMarker();
            for (ObsObject object : objectListing.getObjects()) {
                if (object.getObjectKey().endsWith("/")) {
                    try {
                        GetSnapshotListResponse getSnapshotListResponse = obsClient.getSnapshotList(new GetSnapshotListRequest(bucket, object.getObjectKey()));
                        for (Snapshot snapshot : getSnapshotListResponse.getSnapshotList()) {
                            obsClient.deleteSnapshot(new DeleteSnapshotRequest(bucket, object.getObjectKey(), snapshot.getSnapshotName()));
                        }
                    } catch (ObsException exception) {
                        System.out.println(exception.getMessage() + " on postTest.");
                    }
                    String marker = null;
                    do{
                        GetSnapshottableDirListResult getSnapshottableDirListResult = obsClient.getSnapshottableDirList(new GetSnapshottableDirListRequest(bucket, marker));
                        marker = getSnapshottableDirListResult.getNextMarker();
                        for (SnapshottableDir snapshottableDir : getSnapshottableDirListResult.getSnapshottableDir()) {
                            obsClient.setDisallowSnapshot(new SetDisallowSnapshotRequest(bucket, snapshottableDir.getParentFullPath()));
                        }
                    } while (marker != null);

                    DropFolderRequest dropFolderRequest = new DropFolderRequest(bucket, object.getObjectKey());
                    obsClient.dropFolder(dropFolderRequest);
                }
            }
        } while (objectMarker != null);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_001() throws IOException {
        String testDir = "my-dir-1/";
        NewFolderRequest newFolderRequest = new NewFolderRequest(snapshotBucket, testDir);
        HeaderResponse createDirectoryResponse = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse.getStatusCode());

        SetSnapshotAllowRequest setSnapshotAllowRequest = new SetSnapshotAllowRequest(snapshotBucket, testDir);
        HeaderResponse setSnapshotAllowResponse = obsClient.setSnapshotAllow(setSnapshotAllowRequest);
        assertEquals(200, setSnapshotAllowResponse.getStatusCode());

        SetSnapshotAllowRequest setSnapshotAllowRequest2 = new SetSnapshotAllowRequest(snapshotBucket, testDir);
        HeaderResponse setSnapshotAllowResponse2 = obsClient.setSnapshotAllow(setSnapshotAllowRequest2);
        assertEquals(200, setSnapshotAllowResponse2.getStatusCode());

        CreateSnapshotRequest createSnapshotRequest = new CreateSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        CreateSnapshotResponse createSnapshotResponse = obsClient.createSnapshot(createSnapshotRequest);
        assertEquals(200, createSnapshotResponse.getStatusCode());

        DeleteSnapshotRequest deleteSnapshotRequest = new DeleteSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        HeaderResponse deleteSnapshotResponse = obsClient.deleteSnapshot(deleteSnapshotRequest);
        assertEquals(200, deleteSnapshotResponse.getStatusCode());

        SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, testDir);
        HeaderResponse setDisallowSnapshotResponse = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        assertEquals(200, setDisallowSnapshotResponse.getStatusCode());

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_002() throws IOException {
        String testDir = "my-dir-2/test/";
        String testDir1 = testDir + ".Trash/c/";
        String testDir2 = testDir + ".snapshot/c/";
        String testDir3 = "user/zhangshan/c/";
        String testDir4 = "user/hadoop/data/d/f/e/";

        SetSnapshotAllowRequest setSnapshotAllowRequest1 = new SetSnapshotAllowRequest(snapshotBucket, testDir);
        try{
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest1);
            fail("Scenario should fail with 404, but got: " + response.getStatusCode());
        } catch (ObsException obsException) {
            assertEquals(404, obsException.getResponseCode());
        }
        NewFolderRequest newFolderRequest = new NewFolderRequest(snapshotBucket, "");

        newFolderRequest.setObjectKey(testDir1);
        HeaderResponse createDirectoryResponse1 = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse1.getStatusCode());

        newFolderRequest.setObjectKey(testDir2);
        HeaderResponse createDirectoryResponse2 = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse2.getStatusCode());

        newFolderRequest.setObjectKey(testDir3);
        HeaderResponse createDirectoryResponse3 = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse3.getStatusCode());

        newFolderRequest.setObjectKey(testDir4);
        HeaderResponse createDirectoryResponse4 = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse4.getStatusCode());

        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest2 = new SetSnapshotAllowRequest(snapshotBucket, testDir+".snapshot/");
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest2);
            fail("Scenario 1 should fail with 405, but got: " + response.getStatusCode());

        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: The snapshot absolute path or subdirectory contains .snapshot or .Trash.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest3 = new SetSnapshotAllowRequest(snapshotBucket, testDir+".Trash/");
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest3);
            fail("Scenario should fail with 405, but got: " + response.getStatusCode());

        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: The snapshot absolute path or subdirectory contains .snapshot or .Trash.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest4 = new SetSnapshotAllowRequest(snapshotBucket, testDir1);
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest4);
            fail("Scenario should fail with 405, but got: " + response.getStatusCode());

        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: The snapshot absolute path or subdirectory contains .snapshot or .Trash.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest5 = new SetSnapshotAllowRequest(snapshotBucket, testDir2);
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest5);
            fail("Scenario should fail with 405, but got: " + response.getStatusCode());

        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: The snapshot absolute path or subdirectory contains .snapshot or .Trash.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest6 = new SetSnapshotAllowRequest(snapshotBucket, "user/");
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest6);
            fail("Scenario should fail with 405, but got: " + response.getStatusCode());

        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: The snapshot absolute path or subdirectory contains .snapshot or .Trash.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest7 = new SetSnapshotAllowRequest(snapshotBucket, "user/zhangshan/");
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest7);
            fail("Scenario should fail with 405, but got: " + response.getStatusCode());

        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: The snapshot absolute path or subdirectory contains .snapshot or .Trash.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        SetSnapshotAllowRequest setSnapshotAllowRequest8 = new SetSnapshotAllowRequest(snapshotBucket, testDir3);
        HeaderResponse setSnapshotResponse8 = obsClient.setSnapshotAllow(setSnapshotAllowRequest8);
        assertEquals(200, setSnapshotResponse8.getStatusCode());

        SetSnapshotAllowRequest setSnapshotAllowRequest9 = new SetSnapshotAllowRequest(snapshotBucket, "user/hadoop/data/d/f/");
        HeaderResponse setSnapshotResponse9 = obsClient.setSnapshotAllow(setSnapshotAllowRequest9);
        assertEquals(200, setSnapshotResponse9.getStatusCode());

        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest10 = new SetSnapshotAllowRequest(snapshotBucket, "user/hadoop/data/d/"); // Nested but doesnt fail
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest10);
            fail("Scenario should fail with 405, but got: " + response.getStatusCode());
        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: Nested snapshottable directories not allowed.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest11 = new SetSnapshotAllowRequest(snapshotBucket, testDir4); // Nested but doesnt fail
            HeaderResponse response = obsClient.setSnapshotAllow(setSnapshotAllowRequest11);
            fail("Scenario should fail with 405, but got: " + response.getStatusCode());
        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: Nested snapshottable directories not allowed.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }

        SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, testDir3);
        HeaderResponse setDisallowSnapshotResponse = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        assertEquals(200, setDisallowSnapshotResponse.getStatusCode());

        SetDisallowSnapshotRequest setDisallowSnapshotRequest2 = new SetDisallowSnapshotRequest(snapshotBucket, "user/hadoop/data/d/f/");
        HeaderResponse setDisallowSnapshotResponse2 = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest2);
        assertEquals(200, setDisallowSnapshotResponse2.getStatusCode());

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);

        DropFolderRequest dropFolderRequest2 = new DropFolderRequest(snapshotBucket,"user/");
        obsClient.dropFolder(dropFolderRequest2);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_003(){
        String testDir = "my-dir-3/";

        NewFolderRequest newFolderRequest = new NewFolderRequest(objectBucket, testDir);
        HeaderResponse createDirectoryResponse = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse.getStatusCode()); // 404 in design
        try{
            SetSnapshotAllowRequest setSnapshotAllowRequest = new SetSnapshotAllowRequest(objectBucket, testDir);
            obsClient.setSnapshotAllow(setSnapshotAllowRequest);
        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("MethodNotAllowed".equals(obsException.getErrorCode())&&
                    ("The specified method is not allowed against this resource.\n" +
                            "Request Error.").equals(obsException.getErrorMessage()));
        }
        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_004(){
        String testDir = "my-dir-4/";

        NewFolderRequest newFolderRequest = new NewFolderRequest(snapshotBucket, testDir);
        HeaderResponse createDirectoryResponse = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse.getStatusCode());

        SetSnapshotAllowRequest setSnapshotAllowRequest = new SetSnapshotAllowRequest(snapshotBucket, testDir);
        HeaderResponse setSnapshotAllowResponse = obsClient.setSnapshotAllow(setSnapshotAllowRequest);
        assertEquals(200, setSnapshotAllowResponse.getStatusCode());

        CreateSnapshotRequest createSnapshotRequest = new CreateSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        CreateSnapshotResponse createSnapshotResponse = obsClient.createSnapshot(createSnapshotRequest);
        assertEquals(200, createSnapshotResponse.getStatusCode());

        SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, testDir);
        try {
            obsClient.setDisallowSnapshot(setDisallowSnapshotRequest); // 200 in test design
        } catch (ObsException obsException) {
            assertEquals(obsException.getResponseCode(), 405);
            assertTrue("FsNotSupport".equals(obsException.getErrorCode()) && ("file system not support this request: " +
                    "The directory has snapshot(s). Please redo the operation after removing all the snapshots.\n" +
                    "Request Error.").equals(obsException.getErrorMessage()));
        }

        DeleteSnapshotRequest deleteSnapshotRequest = new DeleteSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        HeaderResponse deleteSnapshotResponse = obsClient.deleteSnapshot(deleteSnapshotRequest);
        assertEquals(200, deleteSnapshotResponse.getStatusCode()); //204 In test design

        HeaderResponse setDisallowSnapshotResponse2 = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        assertEquals(200, setDisallowSnapshotResponse2.getStatusCode());

        CreateSnapshotRequest createSnapshotRequest2 = new CreateSnapshotRequest(snapshotBucket, testDir, "snapshot_d");
        try{
            obsClient.createSnapshot(createSnapshotRequest2);
        } catch (ObsException obsException) {
            assertEquals(obsException.getResponseCode(), 404); // 400 in desing
        }

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_005(){
        String testDir = "my-dir-5/";

        NewFolderRequest newFolderRequest = new NewFolderRequest(snapshotBucket, testDir);
        HeaderResponse createDirectoryResponse = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse.getStatusCode());

        SetSnapshotAllowRequest setSnapshotAllowRequest = new SetSnapshotAllowRequest(snapshotBucket, testDir);
        HeaderResponse setSnapshotAllowResponse = obsClient.setSnapshotAllow(setSnapshotAllowRequest);
        assertEquals(200, setSnapshotAllowResponse.getStatusCode());

        CreateSnapshotRequest createSnapshotRequest = new CreateSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        CreateSnapshotResponse createSnapshotResponse = obsClient.createSnapshot(createSnapshotRequest);
        assertEquals(200, createSnapshotResponse.getStatusCode());

        try{
            SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, "non_existing_dir/");
            obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        } catch (ObsException obsException) {
            assertEquals(404, obsException.getResponseCode());
            assertTrue("NoSuchKey".equals(obsException.getErrorCode())&&
                    "The specified key does not exist.\nRequest Error.".equals(obsException.getErrorMessage()));
        }
        try{
            SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, testDir);
            obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        } catch (ObsException obsException) {
            assertEquals(405, obsException.getResponseCode());
            assertTrue("FsNotSupport".equals(obsException.getErrorCode())&&
                    ("file system not support this request: The directory has snapshot(s). " +
                            "Please redo the operation after removing all the snapshots." +
                            "\nRequest Error.").equals(obsException.getErrorMessage()));
        }

        DeleteSnapshotRequest deleteSnapshotRequest = new DeleteSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        HeaderResponse deleteSnapshotResponse = obsClient.deleteSnapshot(deleteSnapshotRequest);
        assertEquals(200, deleteSnapshotResponse.getStatusCode());

        SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, testDir);
        HeaderResponse setDisallowSnapshotResponse = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        assertEquals(200, setDisallowSnapshotResponse.getStatusCode());

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_006(){
        String testDir = "my-dir-11/";

        String dir_b = "my-dir-12/";
        NewFolderRequest newFolderRequest = new NewFolderRequest(snapshotBucket, testDir);
        HeaderResponse createDirectoryResponse = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse.getStatusCode());

        NewFolderRequest newFolderRequest2 = new NewFolderRequest(snapshotBucket, dir_b);
        HeaderResponse createDirectoryResponse2 = obsClient.newFolder(newFolderRequest2);
        assertEquals(200, createDirectoryResponse2.getStatusCode());

        GetSnapshottableDirListRequest getSnapshottableDirListRequest = new GetSnapshottableDirListRequest(snapshotBucket);
        GetSnapshottableDirListResult getSnapshottableDirListResult = obsClient.getSnapshottableDirList(getSnapshottableDirListRequest);
        assertEquals(200, getSnapshottableDirListResult.getStatusCode());
        assertEquals(0, getSnapshottableDirListResult.getSnapshottableDirCount());

        SetSnapshotAllowRequest setSnapshotAllowRequest = new SetSnapshotAllowRequest(snapshotBucket, testDir);
        HeaderResponse setSnapshotAllowResponse = obsClient.setSnapshotAllow(setSnapshotAllowRequest);
        assertEquals(200, setSnapshotAllowResponse.getStatusCode());

        SetSnapshotAllowRequest setSnapshotAllowRequest2 = new SetSnapshotAllowRequest(snapshotBucket, dir_b);
        HeaderResponse setSnapshotAllowResponse2 = obsClient.setSnapshotAllow(setSnapshotAllowRequest2);
        assertEquals(200, setSnapshotAllowResponse2.getStatusCode());

        GetSnapshottableDirListRequest getSnapshottableDirListRequest2 = new GetSnapshottableDirListRequest(snapshotBucket, dir_b); //Marker does not work like objectKey
        GetSnapshottableDirListResult getSnapshottableDirListResult2 = obsClient.getSnapshottableDirList(getSnapshottableDirListRequest2);
        assertEquals(200, getSnapshottableDirListResult2.getStatusCode());
        assertEquals(1, getSnapshottableDirListResult2.getSnapshottableDirCount());
        assertEquals(dir_b, getSnapshottableDirListResult2.getSnapshottableDir().get(0).getParentFullPath());

        GetSnapshottableDirListRequest getSnapshottableDirListRequest3 = new GetSnapshottableDirListRequest(snapshotBucket, 1);
        GetSnapshottableDirListResult getSnapshottableDirListResult3 = obsClient.getSnapshottableDirList(getSnapshottableDirListRequest3);
        assertEquals(200, getSnapshottableDirListResult3.getStatusCode());
        assertEquals(1, getSnapshottableDirListResult3.getSnapshottableDirCount());
        assertEquals(testDir, getSnapshottableDirListResult3.getSnapshottableDir().get(0).getParentFullPath());

        GetSnapshottableDirListRequest getSnapshottableDirListRequest4 = new GetSnapshottableDirListRequest(snapshotBucket);
        GetSnapshottableDirListResult getSnapshottableDirListResult4 = obsClient.getSnapshottableDirList(getSnapshottableDirListRequest4);
        assertEquals(200, getSnapshottableDirListResult4.getStatusCode());
        assertEquals(2, getSnapshottableDirListResult4.getSnapshottableDirCount());
        assertEquals(testDir, getSnapshottableDirListResult4.getSnapshottableDir().get(0).getParentFullPath());
        assertEquals(dir_b, getSnapshottableDirListResult4.getSnapshottableDir().get(1).getParentFullPath());

        SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, testDir);
        HeaderResponse setDisallowSnapshotResponse = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        assertEquals(200, setDisallowSnapshotResponse.getStatusCode());

        SetDisallowSnapshotRequest setDisallowSnapshotRequest2 = new SetDisallowSnapshotRequest(snapshotBucket, dir_b);
        HeaderResponse setDisallowSnapshotResponse2 = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest2);
        assertEquals(200, setDisallowSnapshotResponse2.getStatusCode());

        GetSnapshottableDirListRequest getSnapshottableDirListRequest5 = new GetSnapshottableDirListRequest(snapshotBucket);
        GetSnapshottableDirListResult getSnapshottableDirListResult5 = obsClient.getSnapshottableDirList(getSnapshottableDirListRequest5);
        assertEquals(200, getSnapshottableDirListResult5.getStatusCode());
        assertEquals(0, getSnapshottableDirListResult5.getSnapshottableDirCount());

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);

        DropFolderRequest dropFolderRequest2 = new DropFolderRequest(snapshotBucket,dir_b);
        obsClient.dropFolder(dropFolderRequest2);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_007() throws IOException {
        String testDir = "my-dir-7/";

        NewFolderRequest newFolderRequest = new NewFolderRequest(snapshotBucket, testDir);
        HeaderResponse createDirectoryResponse = obsClient.newFolder(newFolderRequest);
        assertEquals(200, createDirectoryResponse.getStatusCode());

        SetSnapshotAllowRequest setSnapshotAllowRequest = new SetSnapshotAllowRequest(snapshotBucket, testDir);
        HeaderResponse setSnapshotAllowResponse = obsClient.setSnapshotAllow(setSnapshotAllowRequest);
        assertEquals(200, setSnapshotAllowResponse.getStatusCode());

        CreateSnapshotRequest createSnapshotRequest = new CreateSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        CreateSnapshotResponse createSnapshotResponse = obsClient.createSnapshot(createSnapshotRequest);
        assertEquals(200, createSnapshotResponse.getStatusCode());

        CreateSnapshotRequest createSnapshotRequest2 = new CreateSnapshotRequest(snapshotBucket, testDir, "test_snapshot1");
        CreateSnapshotResponse createSnapshotResponse2 = obsClient.createSnapshot(createSnapshotRequest2);
        assertEquals(200, createSnapshotResponse2.getStatusCode());

        DeleteSnapshotRequest deleteSnapshotRequest = new DeleteSnapshotRequest(snapshotBucket, testDir, "test_snapshot");
        HeaderResponse deleteSnapshotResponse = obsClient.deleteSnapshot(deleteSnapshotRequest);
        assertEquals(200, deleteSnapshotResponse.getStatusCode());

        DeleteSnapshotRequest deleteSnapshotRequest2 = new DeleteSnapshotRequest(snapshotBucket, testDir, "test_snapshot1");
        HeaderResponse deleteSnapshotResponse2 = obsClient.deleteSnapshot(deleteSnapshotRequest2);
        assertEquals(200, deleteSnapshotResponse2.getStatusCode());

        SetDisallowSnapshotRequest setDisallowSnapshotRequest = new SetDisallowSnapshotRequest(snapshotBucket, testDir);
        HeaderResponse setDisallowSnapshotResponse = obsClient.setDisallowSnapshot(setDisallowSnapshotRequest);
        assertEquals(200, setDisallowSnapshotResponse.getStatusCode());

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_008() throws IOException{
        try {
            CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, SNAPSHOT_C);
            CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);
            fail("Scenario 2 should fail with 404, but got: " + createResponse.getStatusCode());

        } catch (ObsException e) {
            assertEquals("Scenario 2 should return 404", 404, e.getResponseCode());
            String errorMessage = e.getMessage().toLowerCase();
            assertTrue("Error should indicate snapshot root does not exist",
                    errorMessage.contains("not found") ||
                            errorMessage.contains("no such") ||
                            errorMessage.contains("does not exist"));
        }

        try {
            NewFolderRequest newFolderRequest2 = new NewFolderRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse createDirectoryResponse2 = obsClient.newFolder(newFolderRequest2);
            assertEquals(200, createDirectoryResponse2.getStatusCode());

        } catch (ObsException e) {
            fail("Scenario 3 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            SetSnapshotAllowRequest allowRequest = new SetSnapshotAllowRequest(snapshotBucket,DIRECTORY_B);

            HeaderResponse allowResponse = obsClient.setSnapshotAllow(allowRequest);
            assertEquals("Scenario 4 should return 200", 200, allowResponse.getStatusCode());

        } catch (ObsException e) {
            fail("Scenario 4 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, SNAPSHOT_C);
            CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);
            assertEquals("Scenario 5 should return 200", 200, createResponse.getStatusCode());

        } catch (ObsException e) {
            fail("Scenario 5 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, SNAPSHOT_C);
            CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);
            fail("Scenario 6 should fail with 409, but got: " + createResponse.getStatusCode());

        } catch (ObsException e) {
            assertEquals("Scenario 6 should return 409", 409, e.getResponseCode());
            String errorMessage = e.getMessage().toLowerCase();
            assertTrue("Error should indicate conflict/already exists",
                    errorMessage.contains("already exists") ||
                            errorMessage.contains("conflict") ||
                            errorMessage.contains("duplicate"));
        }

        int successfulSnapshots = 1;

        for (int i = 1; i <= 49; i++) {
            try {
                String snapshotName = String.format("snapshot-%02d", i);
                CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, snapshotName);
                CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);

                if (createResponse.getStatusCode() == 200) {
                    successfulSnapshots++;
                } else {
                    fail("Scenario 7 snapshot " + i + " should return 200, but got: " + createResponse.getStatusCode());
                }

            } catch (ObsException e) {
                fail("Scenario 7 snapshot " + i + " failed: " + e.getResponseCode() + " - " + e.getMessage());
            }
        }

        assertEquals("Should have created exactly 50 snapshots", 50, successfulSnapshots);

        try {
            CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, SNAPSHOT_D);
            CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);
            fail("Scenario 8 should fail with 400, but got: " + createResponse.getStatusCode());

        } catch (ObsException e) {
            assertEquals("Scenario 8 should return 400", 400, e.getResponseCode());
            String errorMessage = e.getMessage().toLowerCase();
            assertTrue("Error should indicate maximum limit exceeded",
                    errorMessage.contains("maximum") ||
                            errorMessage.contains("limit") ||
                            errorMessage.contains("exceeded") ||
                            errorMessage.contains("too many"));
        }

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);

    }

    @Test
    public void tc_alpha_java_sdk_snapshot_009() throws IOException{
        try {
            NewFolderRequest newFolderRequest2 = new NewFolderRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse createDirectoryResponse2 = obsClient.newFolder(newFolderRequest2);
            assertEquals(200, createDirectoryResponse2.getStatusCode());
            assertNotNull("Request ID should not be null", createDirectoryResponse2.getRequestId());

        } catch (ObsException e) {
            fail("Scenario 2 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            SetSnapshotAllowRequest allowRequest = new SetSnapshotAllowRequest(snapshotBucket,DIRECTORY_B);

            HeaderResponse allowResponse = obsClient.setSnapshotAllow(allowRequest);
            assertEquals("Scenario 3 should return 200", 200, allowResponse.getStatusCode());
            assertNotNull("Request ID should not be null", allowResponse.getRequestId());

        } catch (ObsException e) {
            fail("Scenario 3 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, SNAPSHOT_C);
            CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);
            assertEquals("Scenario 4 should return 200", 200, createResponse.getStatusCode());
            assertNotNull("Request ID should not be null", createResponse.getRequestId());

        } catch (ObsException e) {
            fail("Scenario 4 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            RenameSnapshotRequest renameRequest = new RenameSnapshotRequest(
                    snapshotBucket,
                    DIRECTORY_B,
                    SNAPSHOT_C,
                    SNAPSHOT_D
            );
            RenameSnapshotResponse renameResponse = obsClient.renameSnapshot(renameRequest);
            assertEquals("Scenario 5 should return 200", 200, renameResponse.getStatusCode());
            assertNotNull("Request ID should not be null", renameResponse.getRequestId());

        } catch (ObsException e) {
            fail("Scenario 5 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            GetSnapshotListRequest listRequest = new GetSnapshotListRequest(snapshotBucket, DIRECTORY_B);
            GetSnapshotListResponse listResponse = obsClient.getSnapshotList(listRequest);

            assertEquals("Verification should return 200", 200, listResponse.getStatusCode());
            assertNotNull("Snapshot list should not be null", listResponse.getSnapshotList());
            assertEquals("Should have exactly 1 snapshot", 1, listResponse.getSnapshotCount());

            boolean snapshotDExists = false;
            boolean snapshotCExists = false;

            for (Snapshot snapshot : listResponse.getSnapshotList()) {
                if (SNAPSHOT_D.equals(snapshot.getSnapshotName())) {
                    snapshotDExists = true;
                }
                if (SNAPSHOT_C.equals(snapshot.getSnapshotName())) {
                    snapshotCExists = true;
                }
            }

            assertTrue("Snapshot d should exist after rename", snapshotDExists);
            assertFalse("Snapshot c should not exist after rename", snapshotCExists);

        } catch (ObsException e) {
            fail("Verification failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_010() throws IOException {

        System.out.println(" SCENARIO 2: Create directory b");
        try {
            NewFolderRequest newFolderRequest = new NewFolderRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse createDirectoryResponse = obsClient.newFolder(newFolderRequest);

            assertEquals("Scenario 2 should return HTTP 200", 200, createDirectoryResponse.getStatusCode());
            assertNotNull("Request ID should not be null", createDirectoryResponse.getRequestId());

            System.out.println(" SUCCESS: Directory created");
            System.out.println("   Status Code: " + createDirectoryResponse.getStatusCode());
            System.out.println("   Request ID: " + createDirectoryResponse.getRequestId());
            System.out.println("   Directory Path: " + DIRECTORY_B);

        } catch (ObsException e) {
            fail(" SCENARIO 2 FAILED: " + e.getResponseCode() + " - " + e.getMessage());
        }

        System.out.println(" SCENARIO 3: Set snapshot allow for directory b");
        try {
            SetSnapshotAllowRequest allowRequest = new SetSnapshotAllowRequest(snapshotBucket, DIRECTORY_B);

            HeaderResponse allowResponse = obsClient.setSnapshotAllow(allowRequest);

            assertEquals("Scenario 3 should return HTTP 200", 200, allowResponse.getStatusCode());
            assertNotNull("Request ID should not be null", allowResponse.getRequestId());

            System.out.println(" SUCCESS: Snapshot allow enabled");
            System.out.println("   Status Code: " + allowResponse.getStatusCode());
            System.out.println("   Request ID: " + allowResponse.getRequestId());
            System.out.println("   Target Directory: " + DIRECTORY_B);
            System.out.println("   Snapshot Allow: ENABLED");

        } catch (ObsException e) {
            fail(" SCENARIO 3 FAILED: " + e.getResponseCode() + " - " + e.getMessage());
        }

        System.out.println(" SCENARIO 4: Rename non-existent snapshot c (Expected 404)");
        try {
            RenameSnapshotRequest renameRequest = new RenameSnapshotRequest(
                    snapshotBucket,
                    DIRECTORY_B,
                    "FODPGKFGOPJKDFOPGJFD",
                    "FDOIGJFDJGFDJGFDJGFD"
            );
            RenameSnapshotResponse renameResponse = obsClient.renameSnapshot(renameRequest);

            fail(" SCENARIO 4 UNEXPECTED: Expected 404 error but got " + renameResponse.getStatusCode());

        } catch (ObsException e) {
            assertEquals("Scenario 4 should return HTTP 404", 404, e.getResponseCode());

            String errorMessage = e.getMessage().toLowerCase();
            assertTrue("Error message should indicate not found",
                    errorMessage.contains("not found") ||
                            errorMessage.contains("no such") ||
                            errorMessage.contains("does not exist"));

            System.out.println(" EXPECTED ERROR: Non-existent snapshot rename failed as expected");
            System.out.println("   Status Code: " + e.getResponseCode());
            System.out.println("   Error Code: " + e.getErrorCode());
            System.out.println("   Error Message: " + e.getMessage());
            System.out.println("   Non-existent Snapshot: " + NON_EXISTENT_SNAPSHOT_C);
            System.out.println("   Target Name: " + NEW_SNAPSHOT_NAME);
        }

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_011() throws IOException {
        System.out.println(" SCENARIO 2: Create directory b");
        try {
            NewFolderRequest newFolderRequest2 = new NewFolderRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse createDirectoryResponse2 = obsClient.newFolder(newFolderRequest2);
            assertEquals(200, createDirectoryResponse2.getStatusCode());

            assertEquals("Scenario 2 should return HTTP 200", 200, createDirectoryResponse2.getStatusCode());
            assertNotNull("Request ID should not be null", createDirectoryResponse2.getRequestId());

            System.out.println(" SUCCESS: Directory created");
            System.out.println("   Status Code: " + createDirectoryResponse2.getStatusCode());
            System.out.println("   Request ID: " + createDirectoryResponse2.getRequestId());
            System.out.println("   Directory Path: " + DIRECTORY_B);
            System.out.println("   Content Type: application/x-directory");

        } catch (ObsException e) {
            fail(" SCENARIO 2 FAILED: " + e.getResponseCode() + " - " + e.getMessage());
        }

        System.out.println(" SCENARIO 3: Set snapshot allow for directory b");
        try {
            SetSnapshotAllowRequest allowRequest = new SetSnapshotAllowRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse allowResponse = obsClient.setSnapshotAllow(allowRequest);

            assertEquals("Scenario 3 should return HTTP 200", 200, allowResponse.getStatusCode());
            assertNotNull("Request ID should not be null", allowResponse.getRequestId());

            System.out.println(" SUCCESS: Snapshot allow enabled");
            System.out.println("   Status Code: " + allowResponse.getStatusCode());
            System.out.println("   Request ID: " + allowResponse.getRequestId());
            System.out.println("   Target Directory: " + DIRECTORY_B);
            System.out.println("   Snapshot Allow: ENABLED");

        } catch (ObsException e) {
            fail(" SCENARIO 3 FAILED: " + e.getResponseCode() + " - " + e.getMessage());
        }

        System.out.println(" SCENARIO 4: Create snapshot c under directory b");
        try {
            CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, SNAPSHOT_C);
            CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);

            assertEquals("Scenario 4 should return HTTP 200", 200, createResponse.getStatusCode());
            assertNotNull("Request ID should not be null", createResponse.getRequestId());

            System.out.println(" SUCCESS: Snapshot created");
            System.out.println("   Status Code: " + createResponse.getStatusCode());
            System.out.println("   Request ID: " + createResponse.getRequestId());
            System.out.println("   Bucket: " + snapshotBucket);
            System.out.println("   Directory: " + DIRECTORY_B);
            System.out.println("   Snapshot Name: " + SNAPSHOT_C);

        } catch (ObsException e) {
            fail(" SCENARIO 4 FAILED: " + e.getResponseCode() + " - " + e.getMessage());
        }

        System.out.println("🗑  SCENARIO 5: Delete snapshot c (Expected 200)");
        try {
            DeleteSnapshotRequest deleteRequest = new DeleteSnapshotRequest(snapshotBucket, DIRECTORY_B, SNAPSHOT_C);
            HeaderResponse deleteResponse = obsClient.deleteSnapshot(deleteRequest);

            assertEquals("Scenario 5 should return HTTP 200", 200, deleteResponse.getStatusCode());
            assertNotNull("Request ID should not be null", deleteResponse.getRequestId());

            System.out.println(" SUCCESS: Snapshot deleted");
            System.out.println("   Status Code: " + deleteResponse.getStatusCode());
            System.out.println("   Request ID: " + deleteResponse.getRequestId());
            System.out.println("   Bucket: " + snapshotBucket);
            System.out.println("   Directory: " + DIRECTORY_B);
            System.out.println("   Deleted Snapshot: " + SNAPSHOT_C);

        } catch (ObsException e) {
            fail(" SCENARIO 5 FAILED: " + e.getResponseCode() + " - " + e.getMessage());
        }
        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_012() throws IOException {
        try {
            DeleteSnapshotRequest deleteRequest = new DeleteSnapshotRequest(
                    snapshotBucket,
                    DUMMY_OBJECT_KEY,
                    NON_EXISTENT_SNAPSHOT_D
            );
            HeaderResponse deleteResponse = obsClient.deleteSnapshot(deleteRequest);

            fail(" SCENARIO 2 UNEXPECTED: Expected 404 error but got " + deleteResponse.getStatusCode());

        } catch (ObsException e) {
            assertEquals("Scenario 2 should return HTTP 404", 404, e.getResponseCode());


        }
        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_013() throws IOException {

        try {
            NewFolderRequest newFolderRequest2 = new NewFolderRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse createDirectoryResponse2 = obsClient.newFolder(newFolderRequest2);
            assertEquals("Scenario 2 should return 200", 200, createDirectoryResponse2.getStatusCode());

        } catch (ObsException e) {
            fail("Scenario 2 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            SetSnapshotAllowRequest allowRequest = new SetSnapshotAllowRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse allowResponse = obsClient.setSnapshotAllow(allowRequest);
            assertEquals("Scenario 3 should return 200", 200, allowResponse.getStatusCode());

        } catch (ObsException e) {
            fail("Scenario 3 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        for (int i = 1; i <= 10; i++) {
            try {
                String snapshotName = String.format("snapshot-%02d", i);
                CreateSnapshotRequest createRequest = new CreateSnapshotRequest(snapshotBucket, DIRECTORY_B, snapshotName);
                CreateSnapshotResponse createResponse = obsClient.createSnapshot(createRequest);
                assertEquals("Scenario 4 snapshot " + i + " should return 200", 200, createResponse.getStatusCode());

            } catch (ObsException e) {
                fail("Scenario 4 snapshot " + i + " failed: " + e.getResponseCode() + " - " + e.getMessage());
            }
        }

        try {
            GetSnapshotListRequest listRequest = new GetSnapshotListRequest(snapshotBucket, DIRECTORY_B);
            GetSnapshotListResponse listResponse = obsClient.getSnapshotList(listRequest);

            assertEquals("Scenario 5 should return 200", 200, listResponse.getStatusCode());
            assertEquals("Should retrieve all 10 snapshots", 10, listResponse.getSnapshotCount());
            assertNotNull("Snapshot list should not be null", listResponse.getSnapshotList());
            assertEquals("Snapshot list size should be 10", 10, listResponse.getSnapshotList().size());

        } catch (ObsException e) {
            fail("Scenario 5 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            GetSnapshotListRequest listRequest = new GetSnapshotListRequest(snapshotBucket, DIRECTORY_B);
            listRequest.setMaxKeys(2);
            GetSnapshotListResponse listResponse = obsClient.getSnapshotList(listRequest);

            assertEquals("Scenario 6 should return 200", 200, listResponse.getStatusCode());
            assertEquals("Should retrieve exactly 2 snapshots", 2, listResponse.getSnapshotCount());
            assertNotNull("Snapshot list should not be null", listResponse.getSnapshotList());
            assertEquals("Snapshot list size should be 2", 2, listResponse.getSnapshotList().size());
            assertTrue("Should be truncated", listResponse.isTruncated());
            assertNotNull("Should have next marker", listResponse.getNextMarker());

            List<Snapshot> snapshots = listResponse.getSnapshotList();
            assertEquals("First snapshot should be snapshot-01", "snapshot-01", snapshots.get(0).getSnapshotName());
            assertEquals("Second snapshot should be snapshot-02", "snapshot-02", snapshots.get(1).getSnapshotName());

        } catch (ObsException e) {
            fail("Scenario 6 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            GetSnapshotListRequest listRequest = new GetSnapshotListRequest(snapshotBucket, DIRECTORY_B);
            listRequest.setMarker("snapshot-05");
            GetSnapshotListResponse listResponse = obsClient.getSnapshotList(listRequest);

            assertEquals("Scenario 7 should return 200", 200, listResponse.getStatusCode());
            assertEquals("Should retrieve 6 snapshots including and after marker", 6, listResponse.getSnapshotCount());
            assertNotNull("Snapshot list should not be null", listResponse.getSnapshotList());
            assertEquals("Snapshot list size should be 6", 6, listResponse.getSnapshotList().size());

            List<Snapshot> snapshots = listResponse.getSnapshotList();
            assertEquals("First snapshot should be snapshot-05 (marker)", "snapshot-05", snapshots.get(0).getSnapshotName());
            assertEquals("Second snapshot should be snapshot-06", "snapshot-06", snapshots.get(1).getSnapshotName());
            assertEquals("Last snapshot should be snapshot-10", "snapshot-10", snapshots.get(5).getSnapshotName());

            for (int i = 0; i < snapshots.size(); i++) {
                String expectedName = String.format("snapshot-%02d", i + 5);
                assertEquals("Snapshot should be in dictionary order", expectedName, snapshots.get(i).getSnapshotName());
            }

        } catch (ObsException e) {
            fail("Scenario 7 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket, testDir);
            obsClient.dropFolder(dropFolderRequest);
        } catch (Exception e) {
            System.out.println("Cleanup failed: " + e.getMessage());
        }
    }

    @Test
    public void tc_alpha_java_sdk_snapshot_014() throws IOException {
        try {
            NewFolderRequest newFolderRequest2 = new NewFolderRequest(snapshotBucket, DIRECTORY_B);
            HeaderResponse createDirectoryResponse2 = obsClient.newFolder(newFolderRequest2);

            assertEquals("Scenario 2 should return 200", 200, createDirectoryResponse2.getStatusCode());

            assertNotNull("Request ID should not be null", createDirectoryResponse2.getRequestId());

        } catch (ObsException e) {
            fail("Scenario 2 failed: " + e.getResponseCode() + " - " + e.getMessage());
        }

        try {
            GetSnapshotListRequest listRequest = new GetSnapshotListRequest(snapshotBucket, DIRECTORY_B);
            GetSnapshotListResponse listResponse = obsClient.getSnapshotList(listRequest);

            assertEquals("Scenario 3 should fail with count of list 0, but got: ",0, listResponse.getSnapshotCount());

        } catch (ObsException e) {
            assertEquals("Scenario 3 should return 404", 404, e.getResponseCode());

            String errorMessage = e.getMessage().toLowerCase();
            assertTrue("Error message should indicate not found or not allowed",
                    errorMessage.contains("not found") ||
                            errorMessage.contains("no such") ||
                            errorMessage.contains("not allowed") ||
                            errorMessage.contains("permission") ||
                            errorMessage.contains("forbidden"));
        }

        DropFolderRequest dropFolderRequest = new DropFolderRequest(snapshotBucket,testDir);
        obsClient.dropFolder(dropFolderRequest);
    }
}
