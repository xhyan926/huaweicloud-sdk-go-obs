/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.internal;

import com.obs.services.internal.service.AbstractRequestConvertor;
import okhttp3.Response;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class AbstractRequestConvertorUnitTest {
    class TestClass extends AbstractRequestConvertor {
        TestClass(ObsProperties obsProperties) {
            this.obsProperties = obsProperties;
        }
        public void verifyResponseContentType(Response response) throws ServiceException {
            super.verifyResponseContentType(response);
        }
    }
    @Rule
    public ExpectedException expectedException = ExpectedException.none();
    @Test
    public void shouldNotThrowException_WhenContentTypeIsAllowed_InVerifyResponseContentType() {
        // Arrange
        Response mockResponse = mock(Response.class);
        ObsProperties obsProperties = mock(ObsProperties.class);
        when(obsProperties.getBoolProperty(ObsConstraint.VERIFY_RESPONSE_CONTENT_TYPE, true)).thenReturn(true);
        TestClass testClass = new TestClass(obsProperties);
        // shouldNotThrowException
        for (String allowedResponseHttpContentType : Constants.ALLOWED_RESPONSE_HTTP_CONTENT_TYPES_FOR_XML) {
            when(mockResponse.header(Constants.CommonHeaders.CONTENT_TYPE)).thenReturn(allowedResponseHttpContentType);
            testClass.verifyResponseContentType(mockResponse);
        }
    }

    @Test
    public void shouldThrowException_WhenContentTypeIsNotAllowed_InVerifyResponseContentType() {
        // Arrange
        String testWrongContentType = "text/plain";
        Response mockResponse = mock(Response.class);
        when(mockResponse.header(Constants.CommonHeaders.CONTENT_TYPE)).thenReturn(testWrongContentType);
        ObsProperties obsProperties = mock(ObsProperties.class);
        when(obsProperties.getBoolProperty(ObsConstraint.VERIFY_RESPONSE_CONTENT_TYPE, true)).thenReturn(true);
        TestClass testClass = new TestClass(obsProperties);

        // Act & Assert
        expectedException.expect(ServiceException.class);
        expectedException.expectMessage("Expected XML document response from OBS but received content type "
            + testWrongContentType);
        testClass.verifyResponseContentType(mockResponse);
    }

    @Test
    public void shouldNotThrowException_WhenVerifyResponseContentTypeIsFalse_InVerifyResponseContentType() {
        // Arrange
        Response mockResponse = mock(Response.class);
        when(mockResponse.header(Constants.CommonHeaders.CONTENT_TYPE)).thenReturn("text/plain");
        ObsProperties obsProperties = mock(ObsProperties.class);
        when(obsProperties.getBoolProperty(ObsConstraint.VERIFY_RESPONSE_CONTENT_TYPE, true)).thenReturn(false);
        TestClass testClass = new TestClass(obsProperties);
        testClass.verifyResponseContentType(mockResponse);

    }
}
