package com.obs.integrated_test;

import com.obs.services.exception.ObsException;
import org.junit.Before;
import org.junit.Test;
import java.util.HashMap;
import java.util.Map;

import static org.junit.Assert.*;

public class ObsExceptionIT {

    private ObsException exception;
    private static final String XML_WITH_AUTH_MESSAGE = "<Error><Code>403</Code><Message>Forbidden</Message><EncodedAuthorizationMessage>testEncodedMessage</EncodedAuthorizationMessage></Error>";
    private static final String EXPECTED_AUTH_MESSAGE = "testEncodedMessage";

    @Before
    public void setUp() {
        exception = new ObsException("Test message", XML_WITH_AUTH_MESSAGE);
    }

    @Test
    public void testEncodedAuthorizationMessageParsing() {
        assertEquals(EXPECTED_AUTH_MESSAGE, exception.getEncodedAuthorizationMessage());
    }

    @Test
    public void testSetAndGetEncodedAuthorizationMessage() {
        ObsException e = new ObsException("Test");
        e.setEncodedAuthorizationMessage("testMessage");
        assertEquals("testMessage", e.getEncodedAuthorizationMessage());
    }

    @Test
    public void testToStringIncludesEncodedAuthorizationMessage() {
        ObsException e = new ObsException("Test", "<Error><EncodedAuthorizationMessage>testMessage</EncodedAuthorizationMessage></Error>");
        String toStringResult = e.toString();
        assertTrue(toStringResult.contains("<EncodedAuthorizationMessage>testMessage</EncodedAuthorizationMessage>"));
    }

    @Test
    public void testExceptionWithEncodedAuthorizationMessage() {
        try {
            throw new ObsException("Test", XML_WITH_AUTH_MESSAGE);
        } catch (ObsException e) {
            assertEquals(EXPECTED_AUTH_MESSAGE, e.getEncodedAuthorizationMessage());
        }
    }
}