package com.obs.integrated_test;

import com.obs.services.internal.ServiceException;
import org.junit.Before;
import org.junit.Test;
import java.util.HashMap;
import java.util.Map;

import static org.junit.Assert.*;

public class ServiceExceptionIT {

    private ServiceException exception;
    private static final String XML_WITH_AUTH_MESSAGE = "<Error><Code>403</Code><Message>Forbidden</Message><EncodedAuthorizationMessage>testEncodedMessage</EncodedAuthorizationMessage></Error>";
    private static final String EXPECTED_AUTH_MESSAGE = "testEncodedMessage";

    @Before
    public void setUp() {
        exception = new ServiceException("Test message", XML_WITH_AUTH_MESSAGE);
    }

    @Test
    public void testEncodedAuthorizationMessageParsing() {
        assertEquals(EXPECTED_AUTH_MESSAGE, exception.getEncodedAuthorizationMessage());
    }

    @Test
    public void testSetAndGetEncodedAuthorizationMessage() {
        ServiceException e = new ServiceException("Test");
        e.setEncodedAuthorizationMessage("testMessage");
        assertEquals("testMessage", e.getEncodedAuthorizationMessage());
    }

    @Test
    public void testToStringIncludesEncodedAuthorizationMessage() {
        ServiceException e = new ServiceException("Test", XML_WITH_AUTH_MESSAGE);
        String toStringResult = e.toString();
        assertTrue(toStringResult.contains("<EncodedAuthorizationMessage>testEncodedMessage</EncodedAuthorizationMessage>"));
    }

    @Test
    public void testExceptionWithEncodedAuthorizationMessage() {
        try {
            throw new ServiceException("Test", XML_WITH_AUTH_MESSAGE);
        } catch (ServiceException e) {
            assertEquals(EXPECTED_AUTH_MESSAGE, e.getEncodedAuthorizationMessage());
        }
    }
}