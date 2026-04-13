/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.internal.util;

import static com.obs.services.internal.utils.ServiceUtils.formatIso8601Date;

import com.obs.services.internal.handler.XmlResponsesSaxParser.BucketLifecycleConfigurationHandler;
import com.obs.services.model.StorageClassEnum;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;

import java.text.ParseException;
import java.util.Date;

public class BucketLifecycleConfigurationHandlerUnitTest {

    private final String[] invalidCharacters = {"\n\t ", "\n", "\t", " "};
    @Rule
    public ExpectedException expectedException = ExpectedException.none();
    @Test
    public void shouldNotThrowException_WhenEndDaysAfterInitiation() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "10";
        testHandler.startRule();
        testHandler.endDaysAfterInitiation(testContent);
        testHandler.endRule("");
        Integer parsedDaysAfterInitiation =
            testHandler.getLifecycleConfig().getRules().get(0).
                getAbortIncompleteMultipartUpload().getDaysAfterInitiation();
        // should equal
        Assert.assertEquals(testContent, "" + parsedDaysAfterInitiation);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endDaysAfterInitiation(testContent + invalidCharacter);
            parsedDaysAfterInitiation =
                testHandler.getLifecycleConfig().getRules().get(0).
                    getAbortIncompleteMultipartUpload().getDaysAfterInitiation();
            Assert.assertNotNull(parsedDaysAfterInitiation);
            Assert.assertEquals(testContent, "" + parsedDaysAfterInitiation);
        }
    }
    @Test
    public void shouldNotParseWrong_WhenEndStorageClass() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "STANDARD";
        testHandler.startRule();
        testHandler.startTransition();
        testHandler.endRule("");
        testHandler.endStorageClass(testContent);
        StorageClassEnum parsedStorageClassEnum =
            testHandler.getLifecycleConfig().getRules().get(0).getTransitions().get(0).getObjectStorageClass();
        // should equal
        Assert.assertNotNull(parsedStorageClassEnum);
        Assert.assertEquals(testContent, parsedStorageClassEnum.getCode());
        // should equal
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endStorageClass(testContent + invalidCharacter);
            parsedStorageClassEnum =
                testHandler.getLifecycleConfig().getRules().get(0).getTransitions().get(0).getObjectStorageClass();
            Assert.assertNotNull(parsedStorageClassEnum);
            Assert.assertEquals(testContent, parsedStorageClassEnum.getCode());
        }
    }
    @Test
    public void shouldNotParseWrong_WhenEndExpiredObjectDeleteMarker() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "true";
        testHandler.startRule();
        testHandler.startExpiration();
        testHandler.endRule("");
        testHandler.endExpiredObjectDeleteMarker(testContent);
        Boolean parsedExpiredObjectDeleteMarker =
            testHandler.getLifecycleConfig().getRules().get(0).getExpiration().getExpiredObjectDeleteMarker();
        // should equal
        Assert.assertNotNull(parsedExpiredObjectDeleteMarker);
        Assert.assertEquals(testContent, parsedExpiredObjectDeleteMarker.toString());

        // should equal
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endExpiredObjectDeleteMarker(testContent + invalidCharacter);
            parsedExpiredObjectDeleteMarker =
                testHandler.getLifecycleConfig().getRules().get(0).getExpiration().getExpiredObjectDeleteMarker();
            Assert.assertNotNull(parsedExpiredObjectDeleteMarker);
            Assert.assertEquals(testContent, parsedExpiredObjectDeleteMarker.toString());
        }
    }

    @Test
    public void shouldNotParseWrong_WhenEndDate() throws ParseException {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "2018-01-01T00:00:00.000Z";
        testHandler.startRule();
        testHandler.startExpiration();
        testHandler.endRule("");
        testHandler.endDate(testContent);
        Date parsedDate =
            testHandler.getLifecycleConfig().getRules().get(0).getExpiration().getDate();
        // should equal
        Assert.assertNotNull(parsedDate);
        Assert.assertEquals(testContent, formatIso8601Date(parsedDate));
        // should equal
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endDate(testContent + invalidCharacter);
            parsedDate =
                testHandler.getLifecycleConfig().getRules().get(0).getExpiration().getDate();
            Assert.assertNotNull(parsedDate);
            Assert.assertEquals(testContent, formatIso8601Date(parsedDate));
        }
    }

    @Test
    public void shouldNotThrowException_WhenEndNoncurrentDays() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "10";
        testHandler.startRule();
        testHandler.startTransition();
        testHandler.endRule("");
        testHandler.endNoncurrentDays(testContent);
        Integer parsedDaysAfterInitiation =
            testHandler.getLifecycleConfig().getRules().get(0).
                getTransitions().get(0).getDays();
        // should equal
        Assert.assertEquals(testContent, "" + parsedDaysAfterInitiation);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endNoncurrentDays(testContent + invalidCharacter);
            parsedDaysAfterInitiation =
                testHandler.getLifecycleConfig().getRules().get(0).getTransitions().get(0).getDays();
            Assert.assertNotNull(parsedDaysAfterInitiation);
            Assert.assertEquals(testContent, "" + parsedDaysAfterInitiation);
        }
    }
    @Test
    public void shouldNotThrowException_WhenEndDays() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "10";
        testHandler.startRule();
        testHandler.startTransition();
        testHandler.endRule("");
        testHandler.endDays(testContent);
        Integer parsedDaysAfterInitiation =
            testHandler.getLifecycleConfig().getRules().get(0).
                getTransitions().get(0).getDays();
        // should equal
        Assert.assertEquals(testContent, "" + parsedDaysAfterInitiation);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endDays(testContent + invalidCharacter);
            parsedDaysAfterInitiation =
                testHandler.getLifecycleConfig().getRules().get(0).getTransitions().get(0).getDays();
            Assert.assertNotNull(parsedDaysAfterInitiation);
            Assert.assertEquals(testContent, "" + parsedDaysAfterInitiation);
        }
    }
    @Test
    public void shouldNotThrowException_WhenEndID() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "test-id";
        testHandler.startRule();
        testHandler.endRule("");
        testHandler.endID(testContent);
        String parsedID =
            testHandler.getLifecycleConfig().getRules().get(0).getId();
        // should equal
        Assert.assertEquals(testContent, parsedID);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endID(testContent + invalidCharacter);
            parsedID = testHandler.getLifecycleConfig().getRules().get(0).getId();
            Assert.assertNotNull(parsedID);
            Assert.assertEquals(testContent, parsedID);
        }
    }
    @Test
    public void shouldNotThrowException_WhenEndPrefix() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "test-Prefix";
        testHandler.startRule();
        testHandler.endRule("");
        testHandler.endPrefix(testContent);
        String parsedPrefix =
            testHandler.getLifecycleConfig().getRules().get(0).getPrefix();
        // should equal
        Assert.assertEquals(testContent, parsedPrefix);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endPrefix(testContent + invalidCharacter);
            parsedPrefix = testHandler.getLifecycleConfig().getRules().get(0).getPrefix();
            Assert.assertNotNull(parsedPrefix);
            Assert.assertEquals(testContent, parsedPrefix);
        }
    }
    @Test
    public void shouldNotThrowException_WhenEndStatus() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "Enabled";
        Boolean testBoolean = true;
        testHandler.startRule();
        testHandler.endRule("");
        testHandler.endStatus(testContent);
        Boolean parsedStatus =
            testHandler.getLifecycleConfig().getRules().get(0).getEnabled();
        // should equal
        Assert.assertEquals(testBoolean, parsedStatus);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.endStatus(testContent + invalidCharacter);
            parsedStatus = testHandler.getLifecycleConfig().getRules().get(0).getEnabled();
            Assert.assertNotNull(parsedStatus);
            Assert.assertEquals(testBoolean, parsedStatus);
        }
    }
    @Test
    public void shouldNotThrowException_WhenEndKey() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "test-Key";
        testHandler.startRule();
        testHandler.endRule("");
        testHandler.endKey(testContent);
        String parsedKey =
            testHandler.getLifecycleConfig().getRules().get(0).getTagSet().getTags().get(0).getKey();
        // should equal
        Assert.assertEquals(testContent, parsedKey);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.getLifecycleConfig().getRules().get(0).getTagSet().getTags().clear();
            testHandler.endKey(testContent + invalidCharacter);
            parsedKey =
                testHandler.getLifecycleConfig().getRules().get(0).getTagSet().getTags().get(0).getKey();
            Assert.assertNotNull(parsedKey);
            Assert.assertEquals(testContent, parsedKey);
        }
    }
    @Test
    public void shouldNotThrowException_WhenEndValue() {
        BucketLifecycleConfigurationHandler testHandler = new BucketLifecycleConfigurationHandler(null);
        String testContent = "test-Value";
        testHandler.startRule();
        testHandler.endRule("");
        testHandler.endKey(testContent);
        testHandler.endValue(testContent);
        String parsedValue =
            testHandler.getLifecycleConfig().getRules().get(0).getTagSet().getTags().get(0).getValue();
        // should equal
        Assert.assertEquals(testContent, parsedValue);
        // shouldNotThrowException
        for (String invalidCharacter : invalidCharacters) {
            testHandler.getLifecycleConfig().getRules().get(0).getTagSet().getTags().clear();
            testHandler.endKey(testContent);
            testHandler.endValue(testContent + invalidCharacter);
            parsedValue =
                testHandler.getLifecycleConfig().getRules().get(0).getTagSet().getTags().get(0).getValue();
            Assert.assertNotNull(parsedValue);
            Assert.assertEquals(testContent, parsedValue);
        }
    }
}
