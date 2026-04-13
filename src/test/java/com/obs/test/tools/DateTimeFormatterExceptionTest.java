package com.obs.test.tools;

import com.obs.services.ObsClient;
import com.obs.services.internal.utils.ServiceUtils;
import com.obs.services.model.PutObjectRequest;
import com.obs.test.TestTools;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.util.Date;
import java.util.Locale;

public class DateTimeFormatterExceptionTest
{
    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Test
    public void test_date_with_wrong_DayOfWeek(){
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "_exampleObjectKey";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        PutObjectRequest putObjectRequest = new PutObjectRequest();
        putObjectRequest.setBucketName(bucketName);
        putObjectRequest.setObjectKey(objectKey);
        putObjectRequest.addUserHeaders("x-obs-date",getDateStringWithWrongDayOfWeek(ServiceUtils.formatRfc822Date(new Date())));
        obsClient.putObject(putObjectRequest);
    }
    @Test
    public void test_date_with_right_DayOfWeek(){
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        String objectKey = bucketName + "_exampleObjectKey";
        ObsClient obsClient = TestTools.getPipelineEnvironment();
        assert obsClient != null;
        PutObjectRequest putObjectRequest = new PutObjectRequest();
        putObjectRequest.setBucketName(bucketName);
        putObjectRequest.setObjectKey(objectKey);
        putObjectRequest.addUserHeaders("x-obs-date", ServiceUtils.formatRfc822Date(new Date()));
        obsClient.putObject(putObjectRequest);
    }

    protected String getDateStringWithWrongDayOfWeek(String dateString){
        if (dateString.startsWith("Mon")) {
            return "Tue" + dateString.substring(3);
        } else{
            return "Mon" + dateString.substring(3);
        }
    }
}
