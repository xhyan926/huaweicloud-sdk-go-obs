package com.obs.integrated_test.bucket_logging;

import com.obs.services.ObsClient;
import com.obs.services.model.AccessControlList;
import com.obs.services.model.BucketLoggingConfiguration;
import com.obs.services.model.GroupGrantee;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.Owner;
import com.obs.services.model.Permission;
import com.obs.services.model.TargetSortingTypeEnum;
import com.obs.test.TestTools;
import org.junit.Test;
import org.junit.rules.TestName;

import java.io.IOException;

import static org.junit.Assert.assertEquals;

public class BucketLoggingIT {
    protected static String bucketName = "test-bucket-logging-bucket";

    @org.junit.Rule
    public TestName testName = new TestName();

    @Test
    public void tc_alpha_java_bucket_logging_001() throws IOException, InterruptedException {
        ObsClient obsClient = TestTools.getPipelineEnvironment_OBS();
        //  1.设置桶log-delivery group WRITE and READ_ACP permissions
        String exampleOwnerId = "domainiddomainiddomainiddo000123";
        AccessControlList acl = new AccessControlList();
        Owner owner = new Owner();
        owner.setId(exampleOwnerId);
        acl.setOwner(owner);
        acl.grantPermission(GroupGrantee.LOG_DELIVERY, Permission.PERMISSION_WRITE);
        acl.grantPermission(GroupGrantee.LOG_DELIVERY, Permission.PERMISSION_READ_ACP);
        obsClient.setBucketAcl(bucketName, acl);
        //  2.指定日志转储归类方式为HOUR设置桶日志配置
        BucketLoggingConfiguration config = new BucketLoggingConfiguration(bucketName, "log_prefix/",
            TargetSortingTypeEnum.HOUR);
        config.setAgency("logtest");
        HeaderResponse result = obsClient.setBucketLogging(bucketName, config);
        assertEquals(200, result.getStatusCode());
        //  3.获取桶日志配置
        BucketLoggingConfiguration get_config = obsClient.getBucketLogging(bucketName);
        assertEquals(TargetSortingTypeEnum.HOUR, get_config.getTargetSorting());
        //  4.指定日志转储归类方式为DAY设置桶日志配置
        BucketLoggingConfiguration config2 = new BucketLoggingConfiguration();
        config2.setAgency("logtest");
        config2.setTargetBucketName(bucketName);
        config2.setLogfilePrefix("log_prefix/");
        config2.setTargetSorting(TargetSortingTypeEnum.DAY);
        HeaderResponse result2 = obsClient.setBucketLogging(bucketName, config2);
        assertEquals(200, result2.getStatusCode());
        //  5.获取桶日志配置
        BucketLoggingConfiguration get_config2 = obsClient.getBucketLogging(bucketName);
        assertEquals(TargetSortingTypeEnum.DAY, get_config2.getTargetSorting());
    }
}
