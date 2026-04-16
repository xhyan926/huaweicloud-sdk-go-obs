/**
 * Copyright 2026 Huawei Technologies Co.,Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
package com.example.obs;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collections;

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.mirrorback.DeleteBucketMirrorBackToSourceRequest;
import com.obs.services.model.mirrorback.GetBucketMirrorBackToSourceRequest;
import com.obs.services.model.mirrorback.GetBucketMirrorBackToSourceResult;
import com.obs.services.model.mirrorback.MirrorBackCondition;
import com.obs.services.model.mirrorback.MirrorBackHttpHeader;
import com.obs.services.model.mirrorback.MirrorBackHttpHeaderSet;
import com.obs.services.model.mirrorback.MirrorBackPublicSource;
import com.obs.services.model.mirrorback.MirrorBackRedirect;
import com.obs.services.model.mirrorback.MirrorBackSourceEndpoint;
import com.obs.services.model.mirrorback.MirrorBackToSourceConfiguration;
import com.obs.services.model.mirrorback.MirrorBackToSourceRule;
import com.obs.services.model.mirrorback.SetBucketMirrorBackToSourceRequest;

/**
 * This sample demonstrates how to set, get, and delete the mirror back to source
 * policy of a bucket using the OBS SDK for Java.
 */
public class MirrorBackToSourceSample {
    private static final String END_POINT = "https://your-endpoint";

    private static final String AK = "*** Provide your Access Key ***";

    private static final String SK = "*** Provide your Secret Key ***";

    private static final String BUCKET_NAME = "my-obs-bucket-demo";

    private static ObsClient obsClient;

    public static void main(String[] args) {
        ObsConfiguration config = new ObsConfiguration();
        config.setSocketTimeout(30000);
        config.setConnectionTimeout(10000);
        config.setEndPoint(END_POINT);
        try {
            obsClient = new ObsClient(AK, SK, config);

            setBucketMirrorBackToSource();

            getBucketMirrorBackToSource();

            deleteBucketMirrorBackToSource();

        } catch (ObsException e) {
            System.out.println("Response Code: " + e.getResponseCode());
            System.out.println("Error Message: " + e.getErrorMessage());
            System.out.println("Error Code:       " + e.getErrorCode());
            System.out.println("Request ID:      " + e.getErrorRequestId());
            System.out.println("Host ID:           " + e.getErrorHostId());
        } finally {
            if (obsClient != null) {
                try {
                    obsClient.close();
                } catch (IOException ignored) {
                }
            }
        }
    }

    private static void setBucketMirrorBackToSource() throws ObsException {
        System.out.println("Setting bucket mirror back to source policy\n");

        // Step 1: Create condition
        MirrorBackCondition condition = new MirrorBackCondition();
        condition.setHttpErrorCodeReturnedEquals("404");
        condition.setObjectKeyPrefixEquals("images/");

        // Step 2: Create source endpoint
        MirrorBackSourceEndpoint endpoint = new MirrorBackSourceEndpoint();
        endpoint.setMaster(Collections.singletonList("https://source-bucket.obs.example.com"));
        endpoint.setSlave(Collections.singletonList("https://backup-bucket.obs.example.com"));

        // Step 3: Create public source
        MirrorBackPublicSource publicSource = new MirrorBackPublicSource();
        publicSource.setSourceEndpoint(endpoint);

        // Step 4: Create HTTP header rules
        MirrorBackHttpHeader httpHeader = new MirrorBackHttpHeader();
        httpHeader.setPassAll(false);
        httpHeader.setPass(Arrays.asList("content-type", "cache-control"));
        httpHeader.setRemove(Collections.singletonList("authorization"));
        httpHeader.setSet(Collections.singletonList(
            new MirrorBackHttpHeaderSet("x-custom-header", "custom-value")));

        // Step 5: Create redirect configuration
        MirrorBackRedirect redirect = new MirrorBackRedirect();
        redirect.setAgency("your-agency");
        redirect.setPublicSource(publicSource);
        redirect.setPassQueryString(true);
        redirect.setMirrorFollowRedirect(false);
        redirect.setMirrorHttpHeader(httpHeader);
        redirect.setReplaceKeyPrefixWith("backup/");
        redirect.setMirrorAllowHttpMethod(Arrays.asList("GET", "HEAD"));

        // Step 6: Create rule
        MirrorBackToSourceRule rule = new MirrorBackToSourceRule();
        rule.setId("mirror-rule-001");
        rule.setCondition(condition);
        rule.setRedirect(redirect);

        // Step 7: Create configuration with rules
        MirrorBackToSourceConfiguration policyConfig = new MirrorBackToSourceConfiguration();
        policyConfig.setRules(Collections.singletonList(rule));

        // Step 8: Send request
        SetBucketMirrorBackToSourceRequest request =
                new SetBucketMirrorBackToSourceRequest(BUCKET_NAME, policyConfig);
        HeaderResponse response = obsClient.setBucketMirrorBackToSource(request);
        System.out.println("Set mirror back to source policy, status code: "
                + response.getStatusCode() + "\n");
    }

    private static void getBucketMirrorBackToSource() throws ObsException {
        System.out.println("Getting bucket mirror back to source policy\n");

        GetBucketMirrorBackToSourceRequest request =
                new GetBucketMirrorBackToSourceRequest(BUCKET_NAME);
        GetBucketMirrorBackToSourceResult result = obsClient.getBucketMirrorBackToSource(request);

        MirrorBackToSourceConfiguration config = result.getMirrorBackToSourceConfiguration();
        if (config != null && config.getRules() != null) {
            for (MirrorBackToSourceRule rule : config.getRules()) {
                System.out.println("Rule ID: " + rule.getId());
                if (rule.getCondition() != null) {
                    System.out.println("Condition HTTP Error Code: "
                            + rule.getCondition().getHttpErrorCodeReturnedEquals());
                    System.out.println("Condition Object Key Prefix: "
                            + rule.getCondition().getObjectKeyPrefixEquals());
                }
                if (rule.getRedirect() != null) {
                    MirrorBackRedirect redirect = rule.getRedirect();
                    System.out.println("Agency: " + redirect.getAgency());
                    System.out.println("Pass Query String: " + redirect.getPassQueryString());
                    System.out.println("Mirror Follow Redirect: " + redirect.getMirrorFollowRedirect());
                    if (redirect.getPublicSource() != null
                            && redirect.getPublicSource().getSourceEndpoint() != null) {
                        MirrorBackSourceEndpoint ep =
                                redirect.getPublicSource().getSourceEndpoint();
                        System.out.println("Master Endpoints: " + ep.getMaster());
                        System.out.println("Slave Endpoints: " + ep.getSlave());
                    }
                }
                System.out.println();
            }
        }
    }

    private static void deleteBucketMirrorBackToSource() throws ObsException {
        System.out.println("Deleting bucket mirror back to source policy\n");

        DeleteBucketMirrorBackToSourceRequest request =
                new DeleteBucketMirrorBackToSourceRequest(BUCKET_NAME);
        HeaderResponse response = obsClient.deleteBucketMirrorBackToSource(request);
        System.out.println("Delete mirror back to source policy, status code: "
                + response.getStatusCode() + "\n");
    }
}
