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

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.dis.DeleteBucketDisPolicyRequest;
import com.obs.services.model.dis.DisPolicyConfiguration;
import com.obs.services.model.dis.DisPolicyRule;
import com.obs.services.model.dis.GetBucketDisPolicyRequest;
import com.obs.services.model.dis.GetBucketDisPolicyResult;
import com.obs.services.model.dis.SetBucketDisPolicyRequest;

/**
 * This sample demonstrates how to set, get, and delete the DIS notification
 * policy of a bucket using the OBS SDK for Java.
 */
public class DisPolicySample {
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

            setBucketDisPolicy();

            getBucketDisPolicy();

            deleteBucketDisPolicy();

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

    private static void setBucketDisPolicy() throws ObsException {
        System.out.println("Setting bucket DIS notification policy\n");

        // Step 1: Create a DIS notification rule
        DisPolicyRule rule = new DisPolicyRule();
        rule.setId("rule-001");
        rule.setStream("your-dis-stream-name");
        rule.setProject("your-project-id");
        rule.setAgency("your-agency");
        rule.setEvents(Arrays.asList("ObjectCreated:*", "ObjectRemoved:*"));
        rule.setPrefix("input/");
        rule.setSuffix(".txt");

        // Step 2: Create configuration with rules
        DisPolicyConfiguration policyConfig = new DisPolicyConfiguration();
        policyConfig.setRules(Arrays.asList(rule));

        // Step 3: Send request
        SetBucketDisPolicyRequest request =
                new SetBucketDisPolicyRequest(BUCKET_NAME, policyConfig);
        HeaderResponse response = obsClient.setBucketDisPolicy(request);
        System.out.println("Set DIS notification policy, status code: " + response.getStatusCode() + "\n");
    }

    private static void getBucketDisPolicy() throws ObsException {
        System.out.println("Getting bucket DIS notification policy\n");

        GetBucketDisPolicyRequest request = new GetBucketDisPolicyRequest(BUCKET_NAME);
        GetBucketDisPolicyResult result = obsClient.getBucketDisPolicy(request);

        DisPolicyConfiguration config = result.getDisPolicyConfiguration();
        if (config != null && config.getRules() != null) {
            for (DisPolicyRule rule : config.getRules()) {
                System.out.println("Rule ID: " + rule.getId());
                System.out.println("Stream: " + rule.getStream());
                System.out.println("Project: " + rule.getProject());
                System.out.println("Agency: " + rule.getAgency());
                System.out.println("Events: " + rule.getEvents());
                System.out.println("Prefix: " + rule.getPrefix());
                System.out.println("Suffix: " + rule.getSuffix());
                System.out.println();
            }
        }
    }

    private static void deleteBucketDisPolicy() throws ObsException {
        System.out.println("Deleting bucket DIS notification policy\n");

        DeleteBucketDisPolicyRequest request = new DeleteBucketDisPolicyRequest(BUCKET_NAME);
        HeaderResponse response = obsClient.deleteBucketDisPolicy(request);
        System.out.println("Delete DIS notification policy, status code: " + response.getStatusCode() + "\n");
    }
}
