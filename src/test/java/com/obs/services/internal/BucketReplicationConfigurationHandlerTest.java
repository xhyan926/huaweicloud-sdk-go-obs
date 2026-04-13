/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.internal;

import com.obs.services.internal.handler.XmlResponsesSaxParser.BucketReplicationConfigurationHandler;
import com.obs.services.internal.utils.ServiceUtils;
import com.obs.services.model.ReplicationConfiguration;
import org.junit.Test;
import org.xml.sax.InputSource;
import org.xml.sax.SAXException;
import org.xml.sax.XMLReader;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;

import static org.junit.Assert.assertEquals;

public class BucketReplicationConfigurationHandlerTest {
    @Test
    public void should_parse_bucket_replication_configuration_with_deleteData_enabled_correctly() throws IOException, SAXException {

        final String testRuleId = "Rule-1";
        final String testRuleStatus = "Enabled";
        final String testPrefix = "testPrefix-1";
        final String testBucket = "test-bucket-name-1";
        final String testStorageClass = "STANDARD";
        final String testDeleteDataStatus = "Enabled";
        final String testHistoricalObjectReplicationStatus = "Enabled";
        final String testAgency = "testAgency-1";

        String xml = new StringBuilder("<ReplicationConfiguration xmlns=\"http://obs.myhuaweicloud.com/doc/2006-03-01/\">")
            .append("<Rule>")
            .append("<ID>" + testRuleId + "</ID>")
            .append("<Status>" + testRuleStatus + "</Status>")
            .append("<Prefix>" + testPrefix + "</Prefix>")
            .append("<Destination>")
            .append("<Bucket>" + testBucket + "</Bucket>")
            .append("<StorageClass>" + testStorageClass + "</StorageClass>")
            .append("<DeleteData>" + testDeleteDataStatus + "</DeleteData>")
            .append("</Destination>")
            .append("<HistoricalObjectReplication>" + testHistoricalObjectReplicationStatus
                + "</HistoricalObjectReplication>")
            .append("</Rule>")
            .append("<Agency>" + testAgency + "</Agency>")
            .append("</ReplicationConfiguration>")
            .toString();

        InputStream inputStream = new ByteArrayInputStream(xml.getBytes());

        BucketReplicationConfigurationHandler handler = new BucketReplicationConfigurationHandler();
        XMLReader xmlReader = ServiceUtils.loadXMLReader();
        xmlReader.setErrorHandler(handler);
        xmlReader.setContentHandler(handler);
        xmlReader.parse(new InputSource(inputStream));

        ReplicationConfiguration replicationConfiguration = handler.getReplicationConfiguration();

        assertEquals(testAgency, replicationConfiguration.getAgency());
        assertEquals(testRuleId, replicationConfiguration.getRules().get(0).getId());
        assertEquals(testRuleStatus, replicationConfiguration.getRules().get(0).getStatus().getCode());
        assertEquals(testPrefix, replicationConfiguration.getRules().get(0).getPrefix());
        assertEquals(testBucket, replicationConfiguration.getRules().get(0).getDestination().getBucket());
        assertEquals(testStorageClass,
            replicationConfiguration.getRules().get(0).getDestination().getObjectStorageClass().getCode());
        assertEquals(testDeleteDataStatus,
            replicationConfiguration.getRules().get(0).getDestination().getDeleteData().getCode());
    }
    @Test
    public void should_parse_bucket_replication_configuration_with_deleteData_disabled_correctly() throws IOException, SAXException {

        final String testRuleId = "Rule-1";
        final String testRuleStatus = "Enabled";
        final String testPrefix = "testPrefix-1";
        final String testBucket = "test-bucket-name-1";
        final String testStorageClass = "STANDARD";
        final String testDeleteDataStatus = "Disabled";
        final String testHistoricalObjectReplicationStatus = "Enabled";
        final String testAgency = "testAgency-1";

        String xml = new StringBuilder("<ReplicationConfiguration xmlns=\"http://obs.myhuaweicloud.com/doc/2006-03-01/\">")
            .append("<Rule>")
            .append("<ID>" + testRuleId + "</ID>")
            .append("<Status>" + testRuleStatus + "</Status>")
            .append("<Prefix>" + testPrefix + "</Prefix>")
            .append("<Destination>")
            .append("<Bucket>" + testBucket + "</Bucket>")
            .append("<StorageClass>" + testStorageClass + "</StorageClass>")
            .append("<DeleteData>" + testDeleteDataStatus + "</DeleteData>")
            .append("</Destination>")
            .append("<HistoricalObjectReplication>" + testHistoricalObjectReplicationStatus
                + "</HistoricalObjectReplication>")
            .append("</Rule>")
            .append("<Agency>" + testAgency + "</Agency>")
            .append("</ReplicationConfiguration>")
            .toString();

        InputStream inputStream = new ByteArrayInputStream(xml.getBytes());

        BucketReplicationConfigurationHandler handler = new BucketReplicationConfigurationHandler();
        XMLReader xmlReader = ServiceUtils.loadXMLReader();
        xmlReader.setErrorHandler(handler);
        xmlReader.setContentHandler(handler);
        xmlReader.parse(new InputSource(inputStream));

        ReplicationConfiguration replicationConfiguration = handler.getReplicationConfiguration();

        assertEquals(testAgency, replicationConfiguration.getAgency());
        assertEquals(testRuleId, replicationConfiguration.getRules().get(0).getId());
        assertEquals(testRuleStatus, replicationConfiguration.getRules().get(0).getStatus().getCode());
        assertEquals(testPrefix, replicationConfiguration.getRules().get(0).getPrefix());
        assertEquals(testBucket, replicationConfiguration.getRules().get(0).getDestination().getBucket());
        assertEquals(testStorageClass,
            replicationConfiguration.getRules().get(0).getDestination().getObjectStorageClass().getCode());
        assertEquals(testDeleteDataStatus,
            replicationConfiguration.getRules().get(0).getDestination().getDeleteData().getCode());
    }
}
