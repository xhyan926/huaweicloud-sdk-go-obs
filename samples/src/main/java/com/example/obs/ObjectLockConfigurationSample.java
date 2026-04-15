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
import com.obs.services.model.objectlock.DefaultRetention;
import com.obs.services.model.objectlock.GetObjectLockConfigurationRequest;
import com.obs.services.model.objectlock.GetObjectLockConfigurationResult;
import com.obs.services.model.objectlock.ObjectLockConfiguration;
import com.obs.services.model.objectlock.ObjectLockRule;
import com.obs.services.model.objectlock.SetObjectLockConfigurationRequest;

/**
 * This sample demonstrates how to set and get the default WORM (ObjectLock)
 * policy of a bucket using the OBS SDK for Java.
 */
public class ObjectLockConfigurationSample {
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

            setObjectLockConfiguration();

            getObjectLockConfiguration();

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

    private static void setObjectLockConfiguration() throws ObsException {
        System.out.println("Setting bucket object lock configuration\n");

        DefaultRetention retention = new DefaultRetention("COMPLIANCE", 30, null);
        ObjectLockRule rule = new ObjectLockRule(retention);
        ObjectLockConfiguration config = new ObjectLockConfiguration("Enabled", rule);
        SetObjectLockConfigurationRequest request =
            new SetObjectLockConfigurationRequest(bucketName, config);

        HeaderResponse response = obsClient.setObjectLockConfiguration(request);
        System.out.println("Set object lock configuration, status code: " + response.getStatusCode() + "\n");
    }

    private static void getObjectLockConfiguration() throws ObsException {
        System.out.println("Getting bucket object lock configuration\n");

        GetObjectLockConfigurationRequest request =
            new GetObjectLockConfigurationRequest(bucketName);
        GetObjectLockConfigurationResult result = obsClient.getObjectLockConfiguration(request);

        ObjectLockConfiguration config = result.getObjectLockConfiguration();
        System.out.println("ObjectLockEnabled: " + config.getObjectLockEnabled());
        if (config.getRule() != null) {
            DefaultRetention retention = config.getRule().getDefaultRetention();
            System.out.println("Mode: " + retention.getMode());
            System.out.println("Days: " + retention.getDays());
            System.out.println("Years: " + retention.getYears());
        }
        System.out.println();
    }
}
