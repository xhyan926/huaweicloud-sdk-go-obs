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

import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.objectlock.ObjectRetention;
import com.obs.services.model.objectlock.SetObjectRetentionRequest;

/**
 * This sample demonstrates how to set object-level WORM protection policy
 * (object retention) using the OBS SDK for Java.
 */
public class ObjectRetentionSample {
    private static final String endPoint = "https://your-endpoint";

    private static final String ak = "*** Provide your Access Key ***";

    private static final String sk = "*** Provide your Secret Key ***";

    private static ObsClient obsClient;

    private static String bucketName = "my-obs-bucket-demo";

    private static String objectKey = "my-object-key";

    public static void main(String[] args) {
        ObsConfiguration config = new ObsConfiguration();
        config.setSocketTimeout(30000);
        config.setConnectionTimeout(10000);
        config.setEndPoint(endPoint);
        try {
            obsClient = new ObsClient(ak, sk, config);

            setObjectRetention();

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

    private static void setObjectRetention() throws ObsException {
        System.out.println("Setting object retention\n");

        // Set retention mode to COMPLIANCE with a future timestamp (milliseconds)
        long retainUntilDate = System.currentTimeMillis() + 30L * 24 * 60 * 60 * 1000;
        ObjectRetention retention = new ObjectRetention("COMPLIANCE", retainUntilDate);
        SetObjectRetentionRequest request =
            new SetObjectRetentionRequest(bucketName, objectKey, retention);

        HeaderResponse response = obsClient.setObjectRetention(request);
        System.out.println("Set object retention, status code: " + response.getStatusCode() + "\n");
    }
}
