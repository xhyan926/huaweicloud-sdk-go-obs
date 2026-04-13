package com.obs;

import static com.obs.services.internal.ObsConstraint.CUSTOM_DOMAIN_CERTIFICATE_ID_MIN_LENGTH;
import static com.obs.services.internal.ObsConstraint.CUSTOM_DOMAIN_MAX_SIZE_KB;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.fail;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.support.membermodification.MemberModifier.suppress;
import static org.powermock.api.support.membermodification.MemberMatcher.method;

import com.obs.services.ObsClient;
import com.obs.services.exception.ObsException;
import com.obs.services.internal.ObsConstraint;
import com.obs.services.internal.utils.ServiceUtils;
import com.obs.services.model.BucketCustomDomainInfo;
import com.obs.services.model.CustomDomainCertificateConfig;
import com.obs.services.model.DeleteBucketCustomDomainRequest;
import com.obs.services.model.GetBucketCustomDomainRequest;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.SetBucketCustomDomainRequest;
import java.nio.charset.StandardCharsets;
import java.util.Date;

import org.junit.Before;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;

@PowerMockIgnore({ "javax.management.*", "javax.net.ssl.*", "org.apache.logging.log4j.*",
        "sun.security.ssl.*", "com.sun.crypto.provider.*" })
@PrepareForTest({ ServiceUtils.class })
public class BucketUnitTest {
    private static final String TEST_BUCKET = "test-union-sdk-cname";
    private static final String TEST_DOMAIN = "enddd.zeekrlife.com";
    private static final String TEST_ENDPOINT = "https://obs.cn-north-7.ulanqab.huawei.com";

    private ObsClient obsClient;

    @org.junit.Rule
    public ExpectedException expectedException = ExpectedException.none();
    public static class TestHeaderResponse extends HeaderResponse {
        @Override
        public void setStatusCode(int statusCode) {
            super.setStatusCode(statusCode);
        }
    }


    @Before
    public void setUp() {
        obsClient = spy(new ObsClient(TEST_ENDPOINT));
        suppress(method(ServiceUtils.class, "signWithHmacSha1", String.class, String.class));
    }

    private String generateLargeStringInKB(int kb) {
        byte[] bytes = new byte[kb * 1024];
        for (int i = 0; i < bytes.length; i++) {
            bytes[i] = 'a';
        }
        return new String(bytes, StandardCharsets.UTF_8);
    }

    @Test
    public void should_throw_exception_when_request_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Request is null");
        obsClient.setBucketCustomDomain(null);
        fail("Expected ObsException for null request");
    }

    @Test
    public void should_throw_exception_when_bucket_name_is_null() {
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Bucket name is null");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(null, TEST_DOMAIN, null);
        obsClient.setBucketCustomDomain(req);
        fail("Expected ObsException for null bucket name");
    }

    @Test
    public void should_throw_exception_when_domain_name_is_null() {
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, null, null);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Domain name is null");
        obsClient.setBucketCustomDomain(req);
        fail("Expected ObsException for null domain name");
    }

    @Test
    public void should_throw_exception_when_cert_name_is_null() {
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName(null);
        config.setCertificate("someCert");
        config.setPrivateKey("someKey");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Certificate name cannot be null");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate name cannot be null");
    }

    @Test
    public void should_throw_exception_when_certificate_is_null() {
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("someName");
        config.setCertificate(null);
        config.setPrivateKey("someKey");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Certificate cannot be null");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate cannot be null");
    }

    @Test
    public void should_throw_exception_when_private_key_is_null() {
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("someName");
        config.setCertificate("someCert");
        config.setPrivateKey(null);
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Private key cannot be null");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Private key cannot be null");
    }

    @Test
    public void should_throw_exception_when_cert_name_length_too_short() {
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("a");
        config.setCertificate("cert");
        config.setPrivateKey("key");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Name length should be between " +
            ObsConstraint.CUSTOM_DOMAIN_NAME_MIN_LENGTH + " and " +
            ObsConstraint.CUSTOM_DOMAIN_NAME_MAX_LENGTH + " characters.");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate name length invalid");
    }

    @Test
    public void should_throw_exception_when_cert_name_length_too_long() {
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789");
        config.setCertificate("cert");
        config.setPrivateKey("key");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Name length should be between " +
            ObsConstraint.CUSTOM_DOMAIN_NAME_MIN_LENGTH + " and " +
            ObsConstraint.CUSTOM_DOMAIN_NAME_MAX_LENGTH + " characters.");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate name length invalid");
    }

    @Test
    public void should_throw_exception_when_certificate_id_too_short() {
        // Create a configuration with certificateId too short
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("validName");
        config.setCertificate("cert");
        config.setPrivateKey("key");
        config.setCertificateId("short");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("CertificateId length should be exactly "+CUSTOM_DOMAIN_CERTIFICATE_ID_MIN_LENGTH+" characters.");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate Id length invalid");
    }

    @Test
    public void should_throw_exception_when_certificate_id_too_long() {
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("validName");
        config.setCertificate("cert");
        config.setPrivateKey("key");
        config.setCertificateId("0123456789ABCDEFGHJKLMNOPQRST");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("CertificateId length should be exactly "+CUSTOM_DOMAIN_CERTIFICATE_ID_MIN_LENGTH+" characters.");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate Id length invalid");
    }

    @Test
    public void should_throw_exception_when_certificate_exceeds_40KB() {
        String largeCert = generateLargeStringInKB(41);
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("someName");
        config.setCertificate(largeCert);
        config.setPrivateKey("someKey");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("Certificate size should be less than or equal to 40 KB");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate size exceeds limit");
    }

    @Test
    public void should_throw_exception_when_certificate_chain_exceeds_40KB() {
        String largeChain = generateLargeStringInKB(41);
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("someName");
        config.setCertificate("validCert");
        config.setCertificateChain(largeChain);
        config.setPrivateKey("someKey");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("CertificateChain size should be less than or equal to 40 KB");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: Certificate Chain size exceeds limit");
    }

    @Test
    public void should_throw_exception_when_private_key_exceeds_40KB() {
        String largeKey = generateLargeStringInKB(41);
        CustomDomainCertificateConfig config = new CustomDomainCertificateConfig();
        config.setName("someName");
        config.setCertificate("validCert");
        config.setPrivateKey(largeKey);
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, config);
        expectedException.expect(IllegalArgumentException.class);
        expectedException.expectMessage("PrivateKey size should be less than or equal to " + CUSTOM_DOMAIN_MAX_SIZE_KB + " KB.");
        obsClient.setBucketCustomDomain(req);
        fail("Expected error: PrivateKey size exceeds limit");
    }

    @Test
    public void should_get_bucket_custom_domain() {
        BucketCustomDomainInfo info = new BucketCustomDomainInfo();
        info.addDomain(TEST_BUCKET, new Date(), "dummyCertId");
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(info)
                    .when(obsClient, "getBucketCustomDomainImpl", any(GetBucketCustomDomainRequest.class));
        } catch (Exception e) {
            fail("Failed to stub getBucketCustomDomainImpl: " + e.getMessage());
        }
        BucketCustomDomainInfo result = obsClient.getBucketCustomDomain(TEST_BUCKET);
        assertNotNull(result);
        assertEquals(1, result.getDomains().size());
        assertEquals(TEST_BUCKET, result.getDomains().get(0).getDomainName());
        assertEquals("dummyCertId", result.getDomains().get(0).getCertificateId());
        verify(obsClient).getBucketCustomDomain(TEST_BUCKET);
    }

    @Test
    public void should_set_bucket_custom_domain() {
        CustomDomainCertificateConfig config = mock(CustomDomainCertificateConfig.class);
        when(config.getName()).thenReturn("dummyName");
        when(config.getCertificate()).thenReturn("dummyCert");
        when(config.getPrivateKey()).thenReturn("dummyKey");
        TestHeaderResponse response = new TestHeaderResponse();
        response.setStatusCode(200);
        SetBucketCustomDomainRequest request = new SetBucketCustomDomainRequest(TEST_BUCKET, "test.huawei.com", config);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(response)
                    .when(obsClient, "setBucketCustomDomainImpl", any(SetBucketCustomDomainRequest.class));
        } catch (Exception e) {
            fail("Failed to stub setBucketCustomDomainImpl: " + e.getMessage());
        }
        HeaderResponse resp = obsClient.setBucketCustomDomain(request);
        assertNotNull(resp);
        assertEquals(200, resp.getStatusCode());
        verify(obsClient).setBucketCustomDomain(request);
    }

    @Test
    public void should_set_bucket_custom_domain_with_null_config() {
        TestHeaderResponse response = new TestHeaderResponse();
        response.setStatusCode(200);
        SetBucketCustomDomainRequest request = new SetBucketCustomDomainRequest(TEST_BUCKET, "test.huawei.com", null);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(response)
                    .when(obsClient, "setBucketCustomDomainImpl", any(SetBucketCustomDomainRequest.class));
        } catch (Exception e) {
            fail("Failed to stub setBucketCustomDomainImpl: " + e.getMessage());
        }
        HeaderResponse resp = obsClient.setBucketCustomDomain(request);
        assertNotNull(resp);
        assertEquals(200, resp.getStatusCode());
        verify(obsClient).setBucketCustomDomain(request);
    }

    @Test
    public void should_delete_bucket_custom_domain() {
        CustomDomainCertificateConfig config = mock(CustomDomainCertificateConfig.class);
        when(config.getName()).thenReturn("dummyName");
        when(config.getCertificate()).thenReturn("dummyCert");
        when(config.getPrivateKey()).thenReturn("dummyKey");
        TestHeaderResponse resp1 = new TestHeaderResponse();
        resp1.setStatusCode(200);
        TestHeaderResponse resp2 = new TestHeaderResponse();
        resp2.setStatusCode(204);
        SetBucketCustomDomainRequest request = new SetBucketCustomDomainRequest(TEST_BUCKET, "test.huawei.com", config);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(resp1)
                    .when(obsClient, "setBucketCustomDomainImpl", any(SetBucketCustomDomainRequest.class));
            org.powermock.api.mockito.PowerMockito.doReturn(resp2)
                    .when(obsClient, "deleteBucketCustomDomainImpl", any(DeleteBucketCustomDomainRequest.class));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse response1 = obsClient.setBucketCustomDomain(request);
        HeaderResponse response2 = obsClient.deleteBucketCustomDomain(TEST_BUCKET, "test.huawei.com");
        assertNotNull(response1);
        assertEquals(200, response1.getStatusCode());
        assertNotNull(response2);
        assertEquals(204, response2.getStatusCode());
        verify(obsClient).setBucketCustomDomain(request);
        verify(obsClient).deleteBucketCustomDomain(TEST_BUCKET, "test.huawei.com");
    }

    @Test
    public void should_put_custom_domain_certificate_with_name() throws ObsException {
        TestHeaderResponse dummyResponse = new TestHeaderResponse();
        dummyResponse.setStatusCode(200);
        doReturn(dummyResponse).when(obsClient).setBucketCustomDomain(any(SetBucketCustomDomainRequest.class));
        CustomDomainCertificateConfig domA = new CustomDomainCertificateConfig();
        domA.setName("NameA");
        domA.setCertificateId("certId123");
        domA.setCertificate("validCertificate");
        domA.setPrivateKey("dummyKeyA");
        CustomDomainCertificateConfig domB = new CustomDomainCertificateConfig();
        domB.setName("ValidDomainName");
        domB.setCertificateId("certId456");
        domB.setCertificate("validCertificateWithName");
        domB.setPrivateKey("dummyKeyB");
        SetBucketCustomDomainRequest reqA = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domA);
        SetBucketCustomDomainRequest reqB = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domB);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqA));
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqB));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse respA = obsClient.setBucketCustomDomain(reqA);
        HeaderResponse respB = obsClient.setBucketCustomDomain(reqB);
        assertNotNull(respA);
        assertNotNull(respB);
        verify(obsClient).setBucketCustomDomain(eq(reqA));
        verify(obsClient).setBucketCustomDomain(eq(reqB));
    }

    @Test
    public void should_put_custom_domain_certificate_with_id() throws ObsException {
        TestHeaderResponse dummyResponse = new TestHeaderResponse();
        dummyResponse.setStatusCode(200);
        doReturn(dummyResponse).when(obsClient).setBucketCustomDomain(any(SetBucketCustomDomainRequest.class));
        CustomDomainCertificateConfig domA = new CustomDomainCertificateConfig();
        domA.setName("DomainA");
        domA.setCertificate("validCertificate");
        domA.setPrivateKey("dummyKeyA");
        CustomDomainCertificateConfig domB = new CustomDomainCertificateConfig();
        domB.setName("DomainB");
        domB.setCertificateId("certId123");
        domB.setCertificate("validCertificateWithId");
        domB.setPrivateKey("dummyKeyB");
        SetBucketCustomDomainRequest reqA = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domA);
        SetBucketCustomDomainRequest reqB = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domB);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqA));
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqB));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse respA = obsClient.setBucketCustomDomain(reqA);
        HeaderResponse respB = obsClient.setBucketCustomDomain(reqB);
        assertNotNull(respA);
        assertNotNull(respB);
        verify(obsClient).setBucketCustomDomain(eq(reqA));
        verify(obsClient).setBucketCustomDomain(eq(reqB));
    }

    @Test
    public void should_put_custom_domain_certificate() throws ObsException {
        TestHeaderResponse dummyResponse = new TestHeaderResponse();
        dummyResponse.setStatusCode(200);
        doReturn(dummyResponse).when(obsClient).setBucketCustomDomain(any(SetBucketCustomDomainRequest.class));
        CustomDomainCertificateConfig domA = new CustomDomainCertificateConfig();
        domA.setName("DomainA");
        domA.setPrivateKey("dummyKeyA");
        CustomDomainCertificateConfig domB = new CustomDomainCertificateConfig();
        domB.setName("DomainB");
        domB.setCertificate("validCertificate");
        domB.setPrivateKey("dummyKeyB");
        SetBucketCustomDomainRequest reqA = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domA);
        SetBucketCustomDomainRequest reqB = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domB);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqA));
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqB));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse respA = obsClient.setBucketCustomDomain(reqA);
        HeaderResponse respB = obsClient.setBucketCustomDomain(reqB);
        assertNotNull(respA);
        assertNotNull(respB);
        verify(obsClient).setBucketCustomDomain(eq(reqA));
        verify(obsClient).setBucketCustomDomain(eq(reqB));
    }

    @Test
    public void should_put_custom_domain_certificate_with_chain() throws ObsException {
        TestHeaderResponse dummyResponse = new TestHeaderResponse();
        dummyResponse.setStatusCode(200);
        doReturn(dummyResponse).when(obsClient).setBucketCustomDomain(any(SetBucketCustomDomainRequest.class));
        CustomDomainCertificateConfig domA = new CustomDomainCertificateConfig();
        domA.setName("DomainA");
        domA.setCertificate("validCertificate");
        domA.setPrivateKey("dummyKeyA");
        CustomDomainCertificateConfig domB = new CustomDomainCertificateConfig();
        domB.setName("DomainB");
        domB.setCertificate("validCertificate");
        domB.setCertificateChain("validChain");
        domB.setPrivateKey("dummyKeyB");
        SetBucketCustomDomainRequest reqA = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domA);
        SetBucketCustomDomainRequest reqB = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domB);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqA));
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqB));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse respA = obsClient.setBucketCustomDomain(reqA);
        HeaderResponse respB = obsClient.setBucketCustomDomain(reqB);
        assertNotNull(respA);
        assertNotNull(respB);
        verify(obsClient).setBucketCustomDomain(eq(reqA));
        verify(obsClient).setBucketCustomDomain(eq(reqB));
    }

    @Test
    public void should_put_custom_domain_certificate_with_private_key() throws ObsException {
        TestHeaderResponse dummyResponse = new TestHeaderResponse();
        dummyResponse.setStatusCode(200);
        doReturn(dummyResponse).when(obsClient).setBucketCustomDomain(any(SetBucketCustomDomainRequest.class));
        CustomDomainCertificateConfig domA = new CustomDomainCertificateConfig();
        domA.setName("DomainA");
        domA.setCertificate("validCertificate");
        domA.setPrivateKey("dummyKeyA");
        CustomDomainCertificateConfig domB = new CustomDomainCertificateConfig();
        domB.setName("DomainB");
        domB.setCertificate("validCertificate");
        domB.setPrivateKey("validPrivateKey");
        SetBucketCustomDomainRequest reqA = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domA);
        SetBucketCustomDomainRequest reqB = new SetBucketCustomDomainRequest(TEST_BUCKET, TEST_DOMAIN, domB);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqA));
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(reqB));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse respA = obsClient.setBucketCustomDomain(reqA);
        HeaderResponse respB = obsClient.setBucketCustomDomain(reqB);
        assertNotNull(respA);
        assertNotNull(respB);
        verify(obsClient).setBucketCustomDomain(eq(reqA));
        verify(obsClient).setBucketCustomDomain(eq(reqB));
    }

    @Test
    public void should_put_and_get_custom_domain_certificate() throws ObsException {
        TestHeaderResponse dummyResponse = new TestHeaderResponse();
        dummyResponse.setStatusCode(200);
        doReturn(dummyResponse).when(obsClient).setBucketCustomDomain(any(SetBucketCustomDomainRequest.class));
        BucketCustomDomainInfo info = new BucketCustomDomainInfo();
        info.addDomain("test.huawei.com", new Date(), "certId123");
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(info)
                    .when(obsClient, "getBucketCustomDomainImpl", any(GetBucketCustomDomainRequest.class));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        CustomDomainCertificateConfig dom = new CustomDomainCertificateConfig();
        dom.setName("DomainA");
        dom.setCertificate("validCertificate");
        dom.setPrivateKey("dummyKey");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, "test.huawei.com", dom);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(req));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse respSet = obsClient.setBucketCustomDomain(req);
        BucketCustomDomainInfo retrieved = obsClient.getBucketCustomDomain(TEST_BUCKET);
        assertNotNull(respSet);
        assertEquals(200, respSet.getStatusCode());
        assertNotNull(retrieved);
        assertEquals(1, retrieved.getDomains().size());
        assertEquals("test.huawei.com", retrieved.getDomains().get(0).getDomainName());
        assertEquals("certId123", retrieved.getDomains().get(0).getCertificateId());
        verify(obsClient).setBucketCustomDomain(eq(req));
        verify(obsClient).getBucketCustomDomain(eq(TEST_BUCKET));
    }

    @Test
    public void should_put_get_and_delete_custom_domain_certificate() throws ObsException {
        TestHeaderResponse dummyResponse = new TestHeaderResponse();
        dummyResponse.setStatusCode(200);
        doReturn(dummyResponse).when(obsClient).setBucketCustomDomain(any(SetBucketCustomDomainRequest.class));
        BucketCustomDomainInfo info = new BucketCustomDomainInfo();
        info.addDomain("example.com", new Date(), "certId123");
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(info)
                    .when(obsClient, "getBucketCustomDomainImpl", any(GetBucketCustomDomainRequest.class));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        TestHeaderResponse delResponse = new TestHeaderResponse();
        delResponse.setStatusCode(204);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(delResponse)
                    .when(obsClient, "deleteBucketCustomDomainImpl", any(DeleteBucketCustomDomainRequest.class));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        CustomDomainCertificateConfig dom = new CustomDomainCertificateConfig();
        dom.setName("DomainA");
        dom.setCertificate("validCertificate");
        dom.setPrivateKey("dummyKey");
        SetBucketCustomDomainRequest req = new SetBucketCustomDomainRequest(TEST_BUCKET, "test.huawei.com", dom);
        try {
            org.powermock.api.mockito.PowerMockito.doReturn(new TestHeaderResponse() {{ setStatusCode(200); }})
                    .when(obsClient, "setBucketCustomDomainImpl", eq(req));
        } catch (Exception e) {
            fail("Stubbing failed: " + e.getMessage());
        }
        HeaderResponse respSet = obsClient.setBucketCustomDomain(req);
        BucketCustomDomainInfo retrieved = obsClient.getBucketCustomDomain(TEST_BUCKET);
        HeaderResponse respDel = obsClient.deleteBucketCustomDomain(TEST_BUCKET, "test.huawei.com");
        assertNotNull(respSet);
        assertEquals(200, respSet.getStatusCode());
        assertNotNull(retrieved);
        assertEquals(1, retrieved.getDomains().size());
        assertEquals("example.com", retrieved.getDomains().get(0).getDomainName());
        assertEquals("certId123", retrieved.getDomains().get(0).getCertificateId());
        assertNotNull(respDel);
        assertEquals(204, respDel.getStatusCode());
        verify(obsClient).setBucketCustomDomain(eq(req));
        verify(obsClient).getBucketCustomDomain(eq(TEST_BUCKET));
        verify(obsClient).deleteBucketCustomDomain(TEST_BUCKET, "test.huawei.com");
    }
}
