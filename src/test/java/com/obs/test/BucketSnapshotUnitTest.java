package com.obs.test;
import com.obs.services.ObsClient;
import com.obs.services.internal.utils.ServiceUtils;
import com.obs.services.model.CreateSnapshotRequest;
import com.obs.services.model.CreateSnapshotResponse;
import com.obs.services.model.DeleteSnapshotRequest;
import com.obs.services.model.GetSnapshotListRequest;
import com.obs.services.model.GetSnapshotListResponse;
import com.obs.services.model.GetSnapshottableDirListRequest;
import com.obs.services.model.GetSnapshottableDirListResult;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.Owner;
import com.obs.services.model.RenameSnapshotRequest;
import com.obs.services.model.RenameSnapshotResponse;
import com.obs.services.model.SetDisallowSnapshotRequest;
import com.obs.services.model.SetSnapshotAllowRequest;
import com.obs.services.model.Snapshot;
import com.obs.services.model.SnapshottableDir;
import org.junit.AfterClass;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.mockserver.integration.ClientAndServer;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;
import static org.mockserver.model.HttpRequest.request;
import static org.mockserver.model.HttpResponse.response;
import static org.junit.Assert.fail;

@PowerMockIgnore({ "javax.management.*", "javax.net.ssl.*", "org.apache.logging.log4j.*",
        "sun.security.ssl.*", "com.sun.crypto.provider.*" })
@PrepareForTest({ ServiceUtils.class })
public class BucketSnapshotUnitTest {
    private static final String TEST_BUCKET = "test-bucket";

    private static final String OBJECT_KEY = "objectKey";
    private static final String SNAPSHOT_NAME = "snapshotName";
    private static final String OLD_SNAPSHOT_NAME = "old-snapshot-001";
    private static final String NEW_SNAPSHOT_NAME = "new-snapshot-001";

    public static final String PROXY_HOST_PROPERTY_NAME = "http.proxyHost";
    public static final String PROXY_PORT_PROPERTY_NAME = "http.proxyPort";
    public static final String PROXY_HOST_S_PROPERTY_NAME = "https.proxyHost";
    public static final String PROXY_PORT_S_PROPERTY_NAME = "https.proxyPort";

    public static final Integer responseCodeForTest = 200;
    private static ClientAndServer mockServer;

    private ObsClient obsClient;

    @org.junit.Rule
    public ExpectedException expectedException = ExpectedException.none();

    @BeforeClass
    public static void setMockServer(){
        // 启动 MockServer
        mockServer = ClientAndServer.startClientAndServer();
        System.setProperty(PROXY_HOST_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_PROPERTY_NAME, "" + mockServer.getLocalPort());
        System.setProperty(PROXY_HOST_S_PROPERTY_NAME, "localhost");
        System.setProperty(PROXY_PORT_S_PROPERTY_NAME, "" + mockServer.getLocalPort());
    }

    @Before
    public void setUpClient(){
        obsClient = TestTools.getPipelineEnvironment();
    }

    @AfterClass
    public static void clearEnv() {
        // 关闭 MockServer
        mockServer.close();
        System.clearProperty(PROXY_HOST_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_PROPERTY_NAME);
        System.clearProperty(PROXY_HOST_S_PROPERTY_NAME);
        System.clearProperty(PROXY_PORT_S_PROPERTY_NAME);
    }

    private String generateLargeString(int length) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < length; i++) {
            sb.append('a');
        }
        return sb.toString();
    }

    @Test
    public void snapshot_name_too_long_fail(){
        expectedException.expect(IllegalArgumentException.class);
        CreateSnapshotRequest createSnapshotRequest = new CreateSnapshotRequest(TEST_BUCKET, OBJECT_KEY, generateLargeString(256));
        obsClient.createSnapshot(createSnapshotRequest);
        fail("Expected ObsException for create snapshot request");
    }

    @Test
    public void object_key_too_long_fail(){
        expectedException.expect(IllegalArgumentException.class);
        CreateSnapshotRequest createSnapshotRequest = new CreateSnapshotRequest(TEST_BUCKET, generateLargeString(1025), SNAPSHOT_NAME);
        obsClient.createSnapshot(createSnapshotRequest);
        fail("Expected ObsException for create snapshot request");
    }

    @Test
    public void max_keys_negative(){
        expectedException.expect(IllegalArgumentException.class);
        GetSnapshottableDirListRequest getSnapshottableDirListRequest = new GetSnapshottableDirListRequest(TEST_BUCKET, -1);
        obsClient.getSnapshottableDirList(getSnapshottableDirListRequest);
        fail("Expected ObsException for create snapshot request");
    }

    @Test
    public void set_snapshot_allow_should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("request is null");
        obsClient.setSnapshotAllow(null);
        fail("Expected ObsException for null request");
    }

    @Test
    public void set_snapshot_allow_should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        SetSnapshotAllowRequest req = new SetSnapshotAllowRequest(null, OBJECT_KEY);
        obsClient.setSnapshotAllow(req);
        fail("Expected ObsException for null bucket name");
    }

    @Test
    public void set_snapshot_allow_should_throw_exception_when_object_key_is_null() {
        SetSnapshotAllowRequest req = new SetSnapshotAllowRequest(TEST_BUCKET, null);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("objectKey is null");
        obsClient.setSnapshotAllow(req);
        fail("Expected ObsException for null domain name");
    }

    @Test
    public void set_snapshot_allow_success() {
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("POST").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));
        SetSnapshotAllowRequest req = new SetSnapshotAllowRequest(TEST_BUCKET, OBJECT_KEY);
        HeaderResponse result = obsClient.setSnapshotAllow(req);
        assertNotNull(result);
        assertEquals(200, result.getStatusCode());
    }

    @Test
    public void delete_snapshot_should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("request is null");
        obsClient.deleteSnapshot(null);
        fail("Expected ObsException for null request");
    }

    @Test
    public void delete_snapshot_should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        DeleteSnapshotRequest req = new DeleteSnapshotRequest(null, OBJECT_KEY, SNAPSHOT_NAME);
        obsClient.deleteSnapshot(req);
        fail("Expected ObsException for null bucket name");
    }

    @Test
    public void delete_snapshot_should_throw_exception_when_object_key_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("objectKey is null");
        DeleteSnapshotRequest req = new DeleteSnapshotRequest(TEST_BUCKET, null, SNAPSHOT_NAME);
        obsClient.deleteSnapshot(req);
        fail("Expected ObsException for null domain name");
    }

    @Test
    public void delete_snapshot_should_throw_exception_when_snapshot_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Snapshot name cannot be null or empty");
        DeleteSnapshotRequest req = new DeleteSnapshotRequest(TEST_BUCKET, OBJECT_KEY, null);
        obsClient.deleteSnapshot(req);
        fail("Expected ObsException for null domain name");
    }

    @Test
    public void delete_snapshot_success() {
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("DELETE").withPath(""))
                .respond(response().withStatusCode(200));
        DeleteSnapshotRequest req = new DeleteSnapshotRequest(TEST_BUCKET, OBJECT_KEY, SNAPSHOT_NAME);
        HeaderResponse result = obsClient.deleteSnapshot(req);
        assertNotNull(result);
        assertEquals(200, result.getStatusCode());
    }

    @Test
    public void set_snapshot_disallow_should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("request is null");
        obsClient.setDisallowSnapshot(null);
        fail("Expected ObsException for null request");
    }

    @Test
    public void set_snapshot_disallow_should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        SetDisallowSnapshotRequest req = new SetDisallowSnapshotRequest(null, OBJECT_KEY);
        obsClient.setDisallowSnapshot(req);
        fail("Expected ObsException for null bucket name");
    }

    @Test
    public void set_snapshot_disallow_should_throw_exception_when_object_key_is_null() {
        SetDisallowSnapshotRequest req = new SetDisallowSnapshotRequest(TEST_BUCKET, null);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("objectKey is null");
        obsClient.setDisallowSnapshot(req);
        fail("Expected ObsException for null domain name");
    }

    @Test
    public void set_snapshot_disallow_success() {
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("DELETE").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));
        SetDisallowSnapshotRequest req = new SetDisallowSnapshotRequest(TEST_BUCKET, OBJECT_KEY);
        HeaderResponse result = obsClient.setDisallowSnapshot(req);
        assertNotNull(result);
        assertEquals(200, result.getStatusCode());
    }

    @Test
    public void get_snapshottable_dir_should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("request is null");
        obsClient.getSnapshottableDirList((GetSnapshottableDirListRequest) null);
        fail("Expected ObsException for null request");
    }

    @Test
    public void get_snapshottable_dir_should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        GetSnapshottableDirListRequest req = new GetSnapshottableDirListRequest(null);
        obsClient.getSnapshottableDirList(req);
        fail("Expected ObsException for null bucket name");
    }

    @Test
    public void get_snapshottable_dir_success_empty_list() {
        String xml = "<SnapshottableDirListBody xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">\n" +
                "  <IsTruncated>false</IsTruncated>\n" +
                "  <MaxKeys>1000</MaxKeys>\n" +
                "  <SnapshottableDirCount>0</SnapshottableDirCount>\n" +
                "</SnapshottableDirListBody>";
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest).withHeader("Content-Type", "application/xml").withBody(xml));
        GetSnapshottableDirListRequest req = new GetSnapshottableDirListRequest(TEST_BUCKET);
        GetSnapshottableDirListResult result = obsClient.getSnapshottableDirList(req);
        assertNotNull(result);
        assertEquals(0, result.getSnapshottableDirCount());
        assertEquals(0, result.getSnapshottableDir().size());
        assertNull(result.getMarker());
        assertNull(result.getNextMarker());
    }

    @Test
    public void get_snapshottable_dir_success() {
        String xml = "<SnapshottableDirListBody xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">\n" +
                "  <IsTruncated>false</IsTruncated>\n" +
                "  <MaxKeys>1000</MaxKeys>\n" +
                "  <SnapshottableDirCount>2</SnapshottableDirCount>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "</SnapshottableDirListBody>";
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest).withHeader("Content-Type", "application/xml").withBody(xml));
        GetSnapshottableDirListRequest req = new GetSnapshottableDirListRequest(TEST_BUCKET, "marker", 1000);
        GetSnapshottableDirListResult result = obsClient.getSnapshottableDirList(req);
        List<SnapshottableDir> snapshottableDirs = result.getSnapshottableDir();
        assertNotNull(result);
        assertEquals(2, result.getSnapshottableDirCount());
        assertEquals(2, result.getSnapshottableDir().size());
        assertEquals(17549011444542L, snapshottableDirs.get(0).getModificationTime().getTime());
        assertEquals(new Owner().getDisplayName(), snapshottableDirs.get(0).getOwner().getDisplayName());
        assertEquals(new Owner().getId(), snapshottableDirs.get(0).getOwner().getId());
        assertEquals("0", snapshottableDirs.get(0).getGroup());
        assertEquals("46111925724881354752", snapshottableDirs.get(0).getFileId());
        assertEquals("16832", snapshottableDirs.get(0).getPermission());
        assertEquals(new Integer(50), snapshottableDirs.get(0).getSnapshotQuota());
        assertEquals("dir11/dir22/dir33/", snapshottableDirs.get(0).getParentFullPath());

        snapshottableDirs.get(1).setPermission("0");
        assertEquals("0", snapshottableDirs.get(1).getPermission());
    }

    @Test
    public void get_snapshottable_dir_no_max_keys_success() {
        String xml = "<SnapshottableDirListBody xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">\n" +
                "  <IsTruncated>false</IsTruncated>\n" +
                "  <MaxKeys>1000</MaxKeys>\n" +
                "  <SnapshottableDirCount>2</SnapshottableDirCount>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "</SnapshottableDirListBody>";
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest).withHeader("Content-Type", "application/xml").withBody(xml));
        GetSnapshottableDirListRequest req = new GetSnapshottableDirListRequest(TEST_BUCKET, "marker");
        GetSnapshottableDirListResult result = obsClient.getSnapshottableDirList(req);
        assertNotNull(result);
        assertEquals(2, result.getSnapshottableDirCount());
        assertEquals(2, result.getSnapshottableDir().size());
    }

    @Test
    public void get_snapshottable_dir_no_marker_success() {
        String xml = "<SnapshottableDirListBody xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">\n" +
                "  <IsTruncated>false</IsTruncated>\n" +
                "  <MaxKeys>1000</MaxKeys>\n" +
                "  <SnapshottableDirCount>2</SnapshottableDirCount>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "</SnapshottableDirListBody>";
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest).withHeader("Content-Type", "application/xml").withBody(xml));
        GetSnapshottableDirListRequest req = new GetSnapshottableDirListRequest(TEST_BUCKET, 1000);
        GetSnapshottableDirListResult result = obsClient.getSnapshottableDirList(req);
        assertNotNull(result);
        assertEquals(2, result.getSnapshottableDirCount());
        assertEquals(2, result.getSnapshottableDir().size());
    }

    @Test
    public void get_snapshottable_dir_no_marker_or_max_keys_success() {
        String xml = "<SnapshottableDirListBody xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">\n" +
                "  <IsTruncated>false</IsTruncated>\n" +
                "  <MaxKeys>1000</MaxKeys>\n" +
                "  <SnapshottableDirCount>2</SnapshottableDirCount>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "  <SnapshottableDir>\n" +
                "    <ModificationTime>17549011444542</ModificationTime>\n" +
                "    <Owner>0</Owner>\n" +
                "    <Group>0</Group>\n" +
                "    <Permission>16832</Permission>\n" +
                "    <SnapshotQuota>50</SnapshotQuota>\n" +
                "    <ParentFullPath>dir11/dir22/dir33/</ParentFullPath>\n" +
                "    <FileId>46111925724881354752</FileId>\n" +
                "  </SnapshottableDir>\n" +
                "</SnapshottableDirListBody>";
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest).withHeader("Content-Type", "application/xml").withBody(xml));
        GetSnapshottableDirListRequest req = new GetSnapshottableDirListRequest(TEST_BUCKET);
        GetSnapshottableDirListResult result = obsClient.getSnapshottableDirList(req);
        assertNotNull(result);
        assertEquals(2, result.getSnapshottableDirCount());
        assertEquals(2, result.getSnapshottableDir().size());
        assertFalse(result.isTruncated());
        assertEquals(1000, result.getMaxKeys());
    }

    @Test
    public void get_snapshottable_dir_list_result_max_keys_zero_when_not_specified() {
        GetSnapshottableDirListResult getSnapshottableDirListResult = new GetSnapshottableDirListResult("marker", "nextMarker", false, 50);
        assertEquals(0, getSnapshottableDirListResult.getSnapshottableDirCount());
    }

    @Test
    public void get_snapshot_list_should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("request is null");
        obsClient.getSnapshotList(null);
        fail("Expected IllegalArgumentException for null request");
    }

    @Test
    public void get_snapshot_list_should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        GetSnapshotListRequest req = new GetSnapshotListRequest(null, null, null);
        req.setMaxKeys(3);
        req.setObjectKey("objKey");
        obsClient.getSnapshotList(req);
        fail("Expected IllegalArgumentException for null bucket name");
    }

    @Test
    public void get_snapshot_list_success() {
        String xml = "<SnapshotListBody xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">\n" +
                "  <IsTruncated>false</IsTruncated>\n" +
                "  <MaxKeys>50</MaxKeys>\n" +
                "  <SnapshotCount>3</SnapshotCount>\n" +
                "</SnapshotListBody>";
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest).withHeader("Content-Type", "application/xml").withBody(xml));
        GetSnapshotListRequest req = new GetSnapshotListRequest(TEST_BUCKET, null, 50);
        req.setMarker("marker");
        GetSnapshotListResponse result = obsClient.getSnapshotList(req);
        assertNotNull(result);
        assertEquals(3, result.getSnapshotCount());
        assertEquals(50, result.getMaxKeys());
        assertNull(result.getMarker());
        assertNull(result.getNextMarker());
        assertFalse(result.isTruncated());
    }

    @Test
    public void get_snapshot_list_empty_result() {
        String xml = "<SnapshotListBody xmlns=\"http://obs.myhwclouds.com/doc/2015-06-30/\">\n" +
                "  <IsTruncated>false</IsTruncated>\n" +
                "  <MaxKeys>50</MaxKeys>\n" +
                "  <SnapshotCount>0</SnapshotCount>\n" +
                "</SnapshotListBody>";
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("GET").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest).withHeader("Content-Type", "application/xml").withBody(xml));

        GetSnapshotListRequest req = new GetSnapshotListRequest(TEST_BUCKET, null, null, 50);
        GetSnapshotListResponse result = obsClient.getSnapshotList(req);
        assertNotNull(result);
        assertEquals(0, result.getSnapshotCount());
        assertTrue(result.getSnapshotList().isEmpty());
    }

    @Test
    public void get_snapshot_list_set(){
        GetSnapshotListResponse getSnapshotListResponse = new GetSnapshotListResponse(null, null, false, 0, 0, null);
        getSnapshotListResponse.setMarker("marker");
        getSnapshotListResponse.setNextMarker("nextMarker");
        getSnapshotListResponse.setTruncated(true);
        getSnapshotListResponse.setMaxKeys(50);
        getSnapshotListResponse.setSnapshotCount(30);
        getSnapshotListResponse.setSnapshotList(new ArrayList<>());

        String expectedString = "GetSnapshotListResponse [marker=" + "marker" + ", nextMarker=" + "nextMarker"
                + ", truncated=" + "true" + ", maxKeys=" + "50" + ", snapshotCount=" + "30"
                + ", snapshotList=" + "[]" + "]";
        assertEquals(expectedString, getSnapshotListResponse.toString());
    }

    // CreateSnapshot Tests
    @Test
    public void create_snapshot_should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("request is null");
        obsClient.createSnapshot(null);
        fail("Expected IllegalArgumentException for null request");
    }

    @Test
    public void create_snapshot_should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        CreateSnapshotRequest req = new CreateSnapshotRequest(null, OBJECT_KEY, SNAPSHOT_NAME);
        obsClient.createSnapshot(req);
        fail("Expected IllegalArgumentException for null bucket name");
    }

    @Test
    public void create_snapshot_should_throw_exception_when_object_key_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("objectKey is null");
        CreateSnapshotRequest req = new CreateSnapshotRequest(TEST_BUCKET, null, SNAPSHOT_NAME);
        obsClient.createSnapshot(req);
        fail("Expected IllegalArgumentException for null object key");
    }

    @Test
    public void create_snapshot_should_throw_exception_when_snapshot_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Snapshot name cannot be null or empty");
        CreateSnapshotRequest req = new CreateSnapshotRequest(TEST_BUCKET, OBJECT_KEY, null);
        obsClient.createSnapshot(req);
        fail("Expected IllegalArgumentException for null snapshot name");
    }

    @Test
    public void create_snapshot_success() {
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("POST").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));
        CreateSnapshotRequest req = new CreateSnapshotRequest(TEST_BUCKET, OBJECT_KEY, SNAPSHOT_NAME);
        CreateSnapshotResponse result = obsClient.createSnapshot(req);
        assertNotNull(result);
    }

    @Test
    public void create_snapshot_success_escape_characters_xml() {
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("POST").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));
        char tab = '\t';    // tab
        char privateUse = '\uE000';
        String stringWithEscapeChars = "<>&\"\\" + tab + privateUse;
        CreateSnapshotRequest req = new CreateSnapshotRequest(TEST_BUCKET, OBJECT_KEY, stringWithEscapeChars);
        CreateSnapshotResponse result = obsClient.createSnapshot(req);
        assertNotNull(result);
    }

    // RenameSnapshot Tests
    @Test
    public void rename_snapshot_should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("request is null");
        obsClient.renameSnapshot(null);
        fail("Expected IllegalArgumentException for null request");
    }

    @Test
    public void rename_snapshot_should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("bucketName is null");
        RenameSnapshotRequest req = new RenameSnapshotRequest(null, OBJECT_KEY, OLD_SNAPSHOT_NAME, NEW_SNAPSHOT_NAME);
        obsClient.renameSnapshot(req);
        fail("Expected IllegalArgumentException for null bucket name");
    }

    @Test
    public void rename_snapshot_should_throw_exception_when_object_key_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("objectKey is null");
        RenameSnapshotRequest req = new RenameSnapshotRequest(TEST_BUCKET, null, OLD_SNAPSHOT_NAME, NEW_SNAPSHOT_NAME);
        obsClient.renameSnapshot(req);
        fail("Expected IllegalArgumentException for null object key");
    }

    @Test
    public void rename_snapshot_should_throw_exception_when_old_snapshot_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("oldSnapshotName is null");
        RenameSnapshotRequest req = new RenameSnapshotRequest(TEST_BUCKET, OBJECT_KEY, null, NEW_SNAPSHOT_NAME);
        obsClient.renameSnapshot(req);
        fail("Expected IllegalArgumentException for null old snapshot name");
    }

    @Test
    public void rename_snapshot_should_throw_exception_when_new_snapshot_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("newSnapshotName is null");
        RenameSnapshotRequest req = new RenameSnapshotRequest(TEST_BUCKET, OBJECT_KEY, OLD_SNAPSHOT_NAME, null);
        obsClient.renameSnapshot(req);
        fail("Expected IllegalArgumentException for null new snapshot name");
    }

    @Test
    public void rename_snapshot_success() {
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        RenameSnapshotRequest req = new RenameSnapshotRequest(TEST_BUCKET, OBJECT_KEY, OLD_SNAPSHOT_NAME, NEW_SNAPSHOT_NAME);
        RenameSnapshotResponse result = obsClient.renameSnapshot(req);
        assertNotNull(result);
    }

    @Test
    public void rename_snapshot_conflict_error() {
        obsClient = TestTools.getPipelineEnvironment();
        mockServer.reset();
        mockServer
                .when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));

        RenameSnapshotRequest req = new RenameSnapshotRequest(TEST_BUCKET, OBJECT_KEY, OLD_SNAPSHOT_NAME, NEW_SNAPSHOT_NAME);
        RenameSnapshotResponse result = obsClient.renameSnapshot(req);
        assertNotNull(result);
    }

    @Test
    public void rename_snapshot_not_found_error() {
        mockServer.reset();
        mockServer
                .when(request().withMethod("PUT").withPath(""))
                .respond(response().withStatusCode(responseCodeForTest));
        RenameSnapshotRequest req = new RenameSnapshotRequest(TEST_BUCKET, OBJECT_KEY, OLD_SNAPSHOT_NAME, NEW_SNAPSHOT_NAME);
        RenameSnapshotResponse result = obsClient.renameSnapshot(req);
        assertNotNull(result);
    }

    @Test
    public void rename_snapshot_to_string() {
        RenameSnapshotRequest renameSnapshotRequest = new RenameSnapshotRequest(TEST_BUCKET, null, null, null);
        renameSnapshotRequest.setObjectKey(OBJECT_KEY);
        renameSnapshotRequest.setOldSnapshotName(OLD_SNAPSHOT_NAME);
        renameSnapshotRequest.setNewSnapshotName(NEW_SNAPSHOT_NAME);

        String expectedString = "RenameSnapshotRequest [bucketName=" + TEST_BUCKET + ", objectKey=" + OBJECT_KEY
                + ", oldSnapshotName=" + OLD_SNAPSHOT_NAME + ", newSnapshotName=" + NEW_SNAPSHOT_NAME + "]";

        assertEquals(expectedString, renameSnapshotRequest.toString());
    }

    @Test
    public void snapshot_initialization(){
        long timestamp1 = 1735084800000L;
        long timestamp2 = 1735084800005L;
        Date date1 = new Date(timestamp1);
        Date date2 = new Date(timestamp2);
        Snapshot snapshot = new Snapshot("snapshot", date1, "snapshotId");
        snapshot.setSnapshotName("snapshotName2");
        snapshot.setModifyTime(date2);
        snapshot.setSnapshotID("snapshotId2");

        assertEquals("snapshotName2", snapshot.getSnapshotName());
        assertEquals(date2, snapshot.getModifyTime());
        assertEquals("snapshotId2", snapshot.getSnapshotID());

        Snapshot snapshot2 = new Snapshot("snapshotName2", date2, "snapshotId2");
        assertEquals(snapshot2, snapshot);
        assertEquals(snapshot.hashCode(), snapshot2.hashCode());
    }

}