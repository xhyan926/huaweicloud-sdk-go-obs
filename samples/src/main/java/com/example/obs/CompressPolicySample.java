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
import com.obs.services.model.compress.CompressPolicyConfiguration;
import com.obs.services.model.compress.CompressPolicyRule;
import com.obs.services.model.compress.DeleteBucketCompressPolicyRequest;
import com.obs.services.model.compress.GetBucketCompressPolicyRequest;
import com.obs.services.model.compress.GetBucketCompressPolicyResult;
import com.obs.services.model.compress.SetBucketCompressPolicyRequest;

/**
 * This sample demonstrates how to set, get, and delete the online decompression
 * (compress) policy of a bucket using the OBS SDK for Java.
 */
public class CompressPolicySample {
    private static final String endPoint = "https://your-endpoint";

    private static final String ak = "*** Provide your Access Key ***";

    private static final String sk = "*** Provide your Secret Key ***";

    private static ObsClient obsClient;

    private static String bucketName = "my-obs-bucket-demo";

    public static void main(String[] args) {
        ObsConfiguration config = new ObsConfiguration();
        config.setSocketTimeout(30000);
        config.setConnectionTimeout(10000);
        config.setEndPoint(endPoint);
        try {
            obsClient = new ObsClient(ak, sk, config);

            setBucketCompressPolicy();

            getBucketCompressPolicy();

            deleteBucketCompressPolicy();

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
                } catch (IOException e) {
                }
            }
        }
    }

    private static void setBucketCompressPolicy() throws ObsException {
        System.out.println("Setting bucket compress policy\n");

        // Step 1: Create a decompression rule
        CompressPolicyRule rule = new CompressPolicyRule();
        rule.setId("rule-001");
        rule.setProject("your-project-id");
        rule.setAgency("your-agency");
        rule.setEvents(Arrays.asList("ObjectCreated:*"));
        rule.setPrefix("input/");
        rule.setSuffix(".zip");
        rule.setOverwrite(0);
        rule.setDecompresspath("output/");
        rule.setPolicytype("decompress");

        // Step 2: Create configuration with rules
        CompressPolicyConfiguration policyConfig = new CompressPolicyConfiguration();
        policyConfig.setRules(Arrays.asList(rule));

        // Step 3: Send request
        SetBucketCompressPolicyRequest request =
                new SetBucketCompressPolicyRequest(bucketName, policyConfig);
        HeaderResponse response = obsClient.setBucketCompressPolicy(request);
        System.out.println("Set compress policy, status code: " + response.getStatusCode() + "\n");
    }

    private static void getBucketCompressPolicy() throws ObsException {
        System.out.println("Getting bucket compress policy\n");

        GetBucketCompressPolicyRequest request = new GetBucketCompressPolicyRequest(bucketName);
        GetBucketCompressPolicyResult result = obsClient.getBucketCompressPolicy(request);

        CompressPolicyConfiguration config = result.getCompressPolicyConfiguration();
        if (config != null && config.getRules() != null) {
            for (CompressPolicyRule rule : config.getRules()) {
                System.out.println("Rule ID: " + rule.getId());
                System.out.println("Project: " + rule.getProject());
                System.out.println("Agency: " + rule.getAgency());
                System.out.println("Events: " + rule.getEvents());
                System.out.println("Prefix: " + rule.getPrefix());
                System.out.println("Suffix: " + rule.getSuffix());
                System.out.println("Overwrite: " + rule.getOverwrite());
                System.out.println("Decompress Path: " + rule.getDecompresspath());
                System.out.println();
            }
        }
    }

    private static void deleteBucketCompressPolicy() throws ObsException {
        System.out.println("Deleting bucket compress policy\n");

        DeleteBucketCompressPolicyRequest request = new DeleteBucketCompressPolicyRequest(bucketName);
        HeaderResponse response = obsClient.deleteBucketCompressPolicy(request);
        System.out.println("Delete compress policy, status code: " + response.getStatusCode() + "\n");
    }
}
