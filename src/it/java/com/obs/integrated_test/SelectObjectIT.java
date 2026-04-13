/**
 * Copyright 2019 Huawei Technologies Co.,Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package com.obs.integrated_test;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URISyntaxException;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.nio.file.Paths;
import java.util.Objects;

import com.obs.services.IObsClient;
import com.obs.services.ObsClient;
import com.obs.services.ObsConfiguration;
import com.obs.services.exception.ObsException;
import com.obs.services.model.UploadFileRequest;
import com.obs.services.model.select.CsvInput;
import com.obs.services.model.select.ExpressionType;
import com.obs.services.model.select.FileHeaderInfo;
import com.obs.services.model.select.InputSerialization;
import com.obs.services.model.select.SelectInputStream;
import com.obs.services.model.select.SelectObjectRequest;
import com.obs.services.model.select.SelectObjectResult;

/**
 * Test MessageDecoder behavior with SELECT requests under different settings.
 *
 * Example run:
 *
 * mvn exec:java -f pom-java.xml \ -Dexec.mainClass=com.obs.test.objects.SelectObjectSample \ -Dexec.classpathScope=test \ -Dexec.args='http://ip:port secretAK secretSK test-bucket narrow-1M.csv "select * from s3object" ./out.csv'
 */
public class SelectObjectIT {
    private static void printObsException(ObsException e) {
        System.out.println("Response Code:   " + e.getResponseCode());
        System.out.println("Response Status: " + e.getResponseStatus());
        System.out.println("Error Message:   " + e.getErrorMessage());
        System.out.println("Error Code:      " + e.getErrorCode());
        System.out.println("Request ID:      " + e.getErrorRequestId());
        System.out.println("Host ID:         " + e.getErrorHostId());
        System.out.println("Message:         " + e.getMessage());
    }

    public static void main(String[] args) throws IOException, URISyntaxException {
        if (args.length != 7) {
            System.err.print("Usage: <endpoint> <ak> <sk> <bucket> <object> <query> <out-file>");
            return;
        }

        final String endpoint = args[0];
        final String ak = args[1];
        final String sk = args[2];
        final String bucket = args[3];
        final String key = args[4];
        final String query = args[5];
        final String outPath = args[6];

        File outFile = new File(outPath);
        outFile.createNewFile();

        ObsConfiguration config = new ObsConfiguration();
        config.setSocketTimeout(30000);
        config.setConnectionTimeout(10000);
        config.setEndPoint(endpoint);

        System.out.println("-- Creating ObsClient instance ...");

        IObsClient obsClient = new ObsClient(ak, sk, config);

        try {
            System.out.printf("-- Checking if bucket '%s' exists ...\n", bucket);
            boolean b = obsClient.headBucket(bucket);
            System.out.printf("-- Bucket '%s' %s\n", bucket, b ? "exists" : "is missing");
            if (!b) {
                obsClient.createBucket(bucket);
            }
        } catch (ObsException e) {
            printObsException(e);

            System.out.printf("-- The bucket '%s' does not exist. Creating ...\n", bucket);
            obsClient.createBucket(bucket);
        }

        try {
            System.out.printf("-- Checking if object '%s' exists ...\n", key);

            if (!obsClient.doesObjectExist(bucket, key)) {
                System.out.printf("-- The object '%s' does not exist. Uploading ...", key);

                URL fileUrl = Objects.requireNonNull(SelectObjectIT.class.getClassLoader().getResource(key));

                obsClient.uploadFile(new UploadFileRequest(bucket, key, Paths.get(fileUrl.toURI()).toString()));
            }

            System.out.println("-- Running SELECT ...");

            SelectObjectRequest selectRequest = new SelectObjectRequest().withBucketName(bucket).withKey(key)
                    .withExpression(query).withExpressionType(ExpressionType.SQL)
                    .withInputSerialization(new InputSerialization().withCsv(new CsvInput().withFieldDelimiter(',')
                            .withRecordDelimiter('\n').withFileHeaderInfo(FileHeaderInfo.USE)));

            SelectObjectResult selectResult = obsClient.selectObjectContent(selectRequest);

            System.out.println("-- Processing the results ...");

            SelectInputStream recordsStream = selectResult.getInputStream();
            FileWriter fileWriter = new FileWriter(outFile);

            BufferedReader reader = new BufferedReader(new InputStreamReader(recordsStream, StandardCharsets.UTF_8));
            long start = System.nanoTime();
            while (true) {
                String line = reader.readLine();
                if (line == null) {
                    break;
                }
                fileWriter.write(line);
                fileWriter.write('\n');
                fileWriter.flush();
            }

            reader.close();
            fileWriter.close();

            System.out.printf("Took %f seconds...\n", (System.nanoTime() - start) / Math.pow(10, 9));
        } catch (ObsException e) {
            printObsException(e);
        }
    }
}
