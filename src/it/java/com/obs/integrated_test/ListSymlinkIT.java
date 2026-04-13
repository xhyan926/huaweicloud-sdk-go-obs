/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.integrated_test;

import com.obs.services.internal.handler.XmlResponsesSaxParser;
import com.obs.services.model.ObjectTypeEnum;

import org.junit.Assert;
import org.junit.Test;
import org.junit.rules.ExpectedException;

import java.io.ByteArrayInputStream;

public class ListSymlinkIT {
    @org.junit.Rule
    public ExpectedException expectedException = ExpectedException.none();
    public static String listObjectsForTest = "<ListBucketResult xmlns=\"http://obs.cn-north-4.myhuaweicloud.com/doc/2015-06-30/\">\n"
        + "<Name>examplebucket</Name>\n" + "<Prefix>obj</Prefix>\n" + "<Marker>obj001</Marker>\n"
        + "<MaxKeys>1000</MaxKeys>\n" + "<IsTruncated>false</IsTruncated>\n"
        // object 1 SYMLINK
        + "  <Contents>\n" + "    <Key>obj001</Key>\n" + "    <LastModified>2015-07-01T02:11:19.775Z</LastModified>\n"
        + "    <ETag>\"a72e382246ac83e86bd203389849e71d\"</ETag>\n" + "    <Size>9</Size>\n" + "    <Owner>\n"
        + "      <ID>b4bf1b36d9ca43d984fbcb9491b6fce9</ID>\n" + "      <DisplayName>ObjectOwnerName</DisplayName>\n"
        + "    </Owner>\n" + "    <StorageClass>STANDARD</StorageClass>\n" + "<Type>SYMLINK</Type>\n"+"  </Contents>\n"
        // object 2 APPENDABLE
        + "  <Contents>\n" + "    <Key>obj002</Key>\n" + "    <LastModified>2015-07-01T02:11:19.775Z</LastModified>\n"
        + "    <ETag>\"a72e382246ac83e86bd203389849e71d\"</ETag>\n" + "    <Size>9</Size>\n" + "    <Owner>\n"
        + "      <ID>b4bf1b36d9ca43d984fbcb9491b6fce9</ID>\n" + "      <DisplayName>ObjectOwnerName</DisplayName>\n"
        + "    </Owner>\n" + "    <StorageClass>STANDARD</StorageClass>\n" + "<Type>APPENDABLE</Type>\n"+"  </Contents>\n"
        // object 3 NORMAL（no <Type> in xml）
        + "  <Contents>\n" + "    <Key>obj003</Key>\n" + "    <LastModified>2015-07-01T02:11:19.775Z</LastModified>\n"
        + "    <ETag>\"a72e382246ac83e86bd203389849e71d\"</ETag>\n" + "    <Size>9</Size>\n" + "    <Owner>\n"
        + "      <ID>b4bf1b36d9ca43d984fbcb9491b6fce9</ID>\n" + "      <DisplayName>ObjectOwnerName</DisplayName>\n"
        + "    </Owner>\n" + "    <StorageClass>STANDARD</StorageClass>\n"+"  </Contents>\n"
        + "</ListBucketResult>";

    public static String listVersionsForTest = "<ListVersionsResult xmlns=\"http://obs.cn-north-4.myhuaweicloud.com/doc/2015-06-30/\">\n"
        + " <Name>bucket02</Name>\n" + "  <Prefix/>\n" + "  <KeyMarker/>\n" + "  <VersionIdMarker/>\n"
        + "  <MaxKeys>1000</MaxKeys>\n" + "  <IsTruncated>false</IsTruncated>\n"
        // object 1 SYMLINK
        + "  <Version>\n"+ "    <Key>object001</Key>\n"
        + "    <VersionId>00011000000000013F16000001643A22E476FFFF9046024ECA3655445346485a</VersionId>\n"
        + "    <IsLatest>true</IsLatest>\n" + "    <LastModified>2015-07-01T00:32:16.482Z</LastModified>\n"
        + "    <ETag>\"2fa3bcaaec668adc5da177e67a122d7c\"</ETag>\n" + "    <Size>12041</Size>\n" + "    <Owner>\n"
        + "      <ID>b4bf1b36d9ca43d984fbcb9491b6fce9</ID>\n" + "      <DisplayName>ObjectOwnerName</DisplayName>\n"
        + "    </Owner>\n" + "    <StorageClass>STANDARD</StorageClass>\n"  + "<Type>SYMLINK</Type>\n" + "  </Version>\n"
        // object 2 APPENDABLE
        + "  <Version>\n"+ "    <Key>object002</Key>\n"
        + "    <VersionId>00011000000000013F16000001643A22E476FFFF9046024ECA3655445346485a</VersionId>\n"
        + "    <IsLatest>true</IsLatest>\n" + "    <LastModified>2015-07-01T00:32:16.482Z</LastModified>\n"
        + "    <ETag>\"2fa3bcaaec668adc5da177e67a122d7c\"</ETag>\n" + "    <Size>12041</Size>\n" + "    <Owner>\n"
        + "      <ID>b4bf1b36d9ca43d984fbcb9491b6fce9</ID>\n" + "      <DisplayName>ObjectOwnerName</DisplayName>\n"
        + "    </Owner>\n" + "    <StorageClass>STANDARD</StorageClass>\n"  + "<Type>APPENDABLE</Type>\n"+ "  </Version>\n"
        // object 3 NORMAL（no <Type> in xml）
        + "  <Version>\n"+ "    <Key>object003</Key>\n"
        + "    <VersionId>00011000000000013F16000001643A22E476FFFF9046024ECA3655445346485a</VersionId>\n"
        + "    <IsLatest>true</IsLatest>\n" + "    <LastModified>2015-07-01T00:32:16.482Z</LastModified>\n"
        + "    <ETag>\"2fa3bcaaec668adc5da177e67a122d7c\"</ETag>\n" + "    <Size>12041</Size>\n" + "    <Owner>\n"
        + "      <ID>b4bf1b36d9ca43d984fbcb9491b6fce9</ID>\n" + "      <DisplayName>ObjectOwnerName</DisplayName>\n"
        + "    </Owner>\n" + "    <StorageClass>STANDARD</StorageClass>\n" + "  </Version>\n"
        + "</ListVersionsResult>";

    @Test
    public void should_parse_objectType_successfully_when_list_object() {
        XmlResponsesSaxParser.ListObjectsHandler listObjectsHandler = (new XmlResponsesSaxParser()).parse(
            new ByteArrayInputStream(listObjectsForTest.getBytes()), XmlResponsesSaxParser.ListObjectsHandler.class, true);
        Assert.assertEquals(ObjectTypeEnum.SYMLINK, listObjectsHandler.getObjects().get(0).getObjectType());
        Assert.assertEquals(ObjectTypeEnum.APPENDABLE, listObjectsHandler.getObjects().get(1).getObjectType());
        Assert.assertEquals(ObjectTypeEnum.NORMAL, listObjectsHandler.getObjects().get(2).getObjectType());
    }

    @Test
    public void should_parse_objectType_successfully_when_list_version() {
        XmlResponsesSaxParser.ListVersionsHandler listVersionsHandler = (new XmlResponsesSaxParser()).parse(
            new ByteArrayInputStream(listVersionsForTest.getBytes()), XmlResponsesSaxParser.ListVersionsHandler.class, true);
        Assert.assertEquals(ObjectTypeEnum.SYMLINK, listVersionsHandler.getItems().get(0).getObjectType());
        Assert.assertEquals(ObjectTypeEnum.APPENDABLE, listVersionsHandler.getItems().get(1).getObjectType());
        Assert.assertEquals(ObjectTypeEnum.NORMAL, listVersionsHandler.getItems().get(2).getObjectType());
    }
}
