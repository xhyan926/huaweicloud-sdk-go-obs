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
import com.obs.services.model.BucketCors;
import com.obs.services.model.BucketCorsRule;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.OptionsInfoRequest;
import com.obs.services.model.OptionsInfoResult;
import com.obs.services.model.SetBucketCorsRequest;

/**
 * This sample demonstrates how to send a bucket preflight request (OPTIONS)
 * for CORS using the OBS SDK for Java.
 */
public class OptionsBucketSample {
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

            setBucketCors();

            optionsBucket();

        } catch (ObsException e) {
            System.out.println("Response Code: " + e.getResponseCode());
            System.out.println("Error Message: " + e.getErrorMessage());
            System.out.println("Request ID: " + e.getErrorRequestId());
        } finally {
            if (obsClient != null) {
                try {
                    obsClient.close();
                } catch (IOException e) {
                    // ignore
                }
            }
        }
    }

    private static void setBucketCors() throws ObsException {
        System.out.println("Setting bucket CORS\n");

        BucketCorsRule rule = new BucketCorsRule();
        rule.setId("cors-rule-001");
        rule.setAllowedOrigin(Arrays.asList("http://www.example.com", "http://www.test.com"));
        rule.setAllowedMethod(Arrays.asList("GET", "PUT", "POST"));
        rule.setAllowedHeader(Collections.singletonList("*"));
        rule.setExposeHeader(Arrays.asList("x-obs-request-id", "x-obs-id-2"));
        rule.setMaxAgeSecond(100);
        BucketCors bucketCors = new BucketCors(Collections.singletonList(rule));

        HeaderResponse response = obsClient.setBucketCors(new SetBucketCorsRequest(BUCKET_NAME, bucketCors));
        System.out.println("Set bucket CORS response status: " + response.getStatusCode() + "\n");
    }

    private static void optionsBucket() throws ObsException {
        System.out.println("Options bucket\n");

        OptionsInfoRequest request = new OptionsInfoRequest(BUCKET_NAME);
        request.setOrigin("http://www.example.com");
        request.setRequestMethod(Arrays.asList("GET", "PUT"));
        request.setRequestHeaders(Collections.singletonList("Authorization"));

        OptionsInfoResult result = obsClient.optionsBucket(request);

        System.out.println("\tAllowOrigin: " + result.getAllowOrigin());
        System.out.println("\tAllowMethods: " + result.getAllowMethods());
        System.out.println("\tAllowHeaders: " + result.getAllowHeaders());
        System.out.println("\tExposeHeaders: " + result.getExposeHeaders());
        System.out.println("\tMaxAge: " + result.getMaxAge());
    }
}
