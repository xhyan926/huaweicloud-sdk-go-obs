# 样例代码模板

本文档提供 OBS SDK 样例代码的编写约定和完整模板，基于 `BucketOperationsSample.java` 提取。

## 文件规范

| 规范项 | 要求 |
|--------|------|
| 文件位置 | `samples/src/main/java/com/example/obs/<Feature>Sample.java` |
| 类命名 | `<Feature>Sample.java`（如 `BucketTrashSample.java`） |
| 包名 | `com.example.obs` |
| 许可证头 | Apache License 2.0（与现有样例保持一致） |

---

## 完整样例模板

```java
/**
 * Copyright 2019 Huawei Technologies Co.,Ltd.
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
import com.obs.services.model.<feature>.<Feature>Configuration;
import com.obs.services.model.<feature>.Set<Feature>Request;
import com.obs.services.model.<feature>.Get<Feature>Request;
import com.obs.services.model.<feature>.Get<Feature>Result;
import com.obs.services.model.<feature>.Delete<Feature>Request;
import com.obs.services.model.HeaderResponse;

/**
 * This sample demonstrates how to do <feature> operations
 * (set/get/delete <feature> configuration) on OBS using the OBS SDK for Java.
 */
public class <Feature>Sample {
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
            // Constructs a obs client instance with your account for accessing OBS
            obsClient = new ObsClient(ak, sk, config);

            // Set <feature> configuration
            set<Feature>();

            // Get <feature> configuration
            get<Feature>();

            // Delete <feature> configuration
            delete<Feature>();

        } catch (ObsException e) {
            System.out.println("Response Code: " + e.getResponseCode());
            System.out.println("Error Message: " + e.getErrorMessage());
            System.out.println("Error Code:       " + e.getErrorCode());
            System.out.println("Request ID:      " + e.getErrorRequestId());
            System.out.println("Host ID:           " + e.getErrorHostId());
        } finally {
            if (obsClient != null) {
                try {
                    // Close obs client
                    obsClient.close();
                } catch (IOException e) {
                }
            }
        }
    }

    private static void set<Feature>() throws ObsException {
        System.out.println("Setting <feature> configuration\n");

        <Feature>Configuration config = new <Feature>Configuration(/* params */);
        Set<Feature>Request request = new Set<Feature>Request(bucketName, config);
        HeaderResponse response = obsClient.set<Feature>(request);
        System.out.println("Set <feature> response status: " + response.getStatusCode() + "\n");
    }

    private static void get<Feature>() throws ObsException {
        System.out.println("Getting <feature> configuration\n");

        Get<Feature>Request request = new Get<Feature>Request(bucketName);
        Get<Feature>Result result = obsClient.get<Feature>(request);
        <Feature>Configuration config = result.get<Feature>Configuration();

        // 打印配置字段
        System.out.println("Configuration: " + config + "\n");
    }

    private static void delete<Feature>() throws ObsException {
        System.out.println("Deleting <feature> configuration\n");

        Delete<Feature>Request request = new Delete<Feature>Request(bucketName);
        HeaderResponse response = obsClient.delete<Feature>(request);
        System.out.println("Delete <feature> response status: " + response.getStatusCode() + "\n");
    }
}
```

---

## 编写规则

1. **使用占位符**: endPoint/ak/sk 使用占位字符串，不包含真实凭证
2. **每个操作独立方法**: set/get/delete 分别为独立 `private static` 方法
3. **添加步骤注释**: 每个操作前添加注释说明操作目的
4. **打印关键结果**: 使用 `System.out.println()` 输出操作结果和状态码
5. **异常处理统一**: 在 `main()` 中统一捕获 `ObsException` 并打印错误信息
6. **资源清理**: `finally` 块中关闭 ObsClient
7. **ObsConfiguration 设置**: 设置 socketTimeout 和 connectionTimeout
8. **类 JavaDoc**: 使用英文描述样例功能
