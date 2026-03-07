// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of the
// License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

// examples/post_upload_sample.go
//
// POST 上传示例
//
// POST 上传是一种允许直接从浏览器向 OBS 上传文件的方法，文件不经过后端服务器。
// 这种方法有以下优势：
//   1. 降低服务器负载和网络带宽
//   2. 支持大文件上传
//   3. 上传速度更快
//
// 使用方法：
//   1. 配置环境变量（见下方）
//   2. 运行示例：go run examples/post_upload_sample.go
//   3. 使用生成的 HTML 表单上传文件
//
// 环境变量：
//   OBS_AK: 你的 Access Key ID
//   OBS_SK: 你的 Secret Access Key
//   OBS_ENDPOINT: OBS 端点，例如 https://obs.cn-north-1.myhuaweicloud.com
//
// 注意事项：
//   - Policy 和签名由后端生成，前端只使用生成的凭证
//   - 可以设置文件大小限制
//   - 可以设置文件类型限制
//   - 策略有过期时间，过期后需要重新生成

package main

import (
	"fmt"
	"log"
	"os"

	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func main() {
	// 1. 从环境变量获取 OBS 凭证
	ak := os.Getenv("OBS_AK")
	sk := os.Getenv("OBS_SK")
	endpoint := os.Getenv("OBS_ENDPOINT")

	if ak == "" || sk == "" || endpoint == "" {
		log.Fatal("请设置环境变量 OBS_AK、OBS_SK 和 OBS_ENDPOINT")
	}

	// 2. 创建 OBS 客户端
	obsClient, err := obs.New(ak, sk, endpoint)
	if err != nil {
		log.Fatalf("创建 OBS 客户端失败: %v", err)
	}

	// 3. 配置 POST 上传策略参数
	bucketName := "your-bucket-name"      // 替换为你的桶名
	objectKey := "uploads/test.jpg"        // 上传后的对象名
	expiresIn := int64(3600)            // 1 小时后过期

	// 4. 创建 POST 上传策略
	input := &obs.CreatePostPolicyInput{
		Bucket:    bucketName,
		Key:       objectKey,
		ExpiresIn: expiresIn,
		Acl:       "public-read", // 可选：设置对象 ACL
		Conditions: []obs.PostPolicyCondition{
			// 添加 content-length-range 条件，限制文件大小为 0-10MB
			obs.CreatePostPolicyCondition(
				"content-length-range",
				"$content-length",
				[]interface{}{0, 10 * 1024 * 1024},
			),
			// 可以添加更多条件，例如：
			// - content-type 限制
			// - 元数据条件
			// - 特定文件名前缀
		},
	}

	// 5. 生成策略和签名
	output, err := obsClient.CreatePostPolicy(input)
	if err != nil {
		log.Fatalf("创建 POST 策略失败: %v", err)
	}

	// 6. 输出策略信息
	fmt.Println("=== POST 上传策略已生成 ===")
	fmt.Printf("Access Key ID: %s\n", output.AccessKeyId)
	fmt.Printf("Policy (Base64): %s\n", output.Policy)
	fmt.Printf("Signature: %s\n", output.Signature)
	fmt.Printf("完整 Token: %s\n\n", output.Token)

	// 7. 生成前端 HTML 表单示例
	fmt.Println("=== 前端 HTML 表单示例 ===")
	fmt.Println(generateHTMLForm(endpoint, bucketName, objectKey, output.AccessKeyId, output.Signature, output.Policy))

	fmt.Println("\n=== 使用说明 ===")
	fmt.Println("1. 将上面的 HTML 表单保存为文件（例如 upload.html）")
	fmt.Println("2. 在浏览器中打开 upload.html")
	fmt.Println("3. 选择文件并点击上传按钮")
	fmt.Println("4. 文件将直接上传到 OBS")
	fmt.Println("\n注意事项：")
	fmt.Println("- 策略在 1 小时后过期，过期后需要重新生成")
	fmt.Println("- 文件大小限制为 10MB")
	fmt.Println("- 文件上传后将公开可读（public-read）")
}

// generateHTMLForm 生成用于 POST 上传的 HTML 表单
func generateHTMLForm(endpoint, bucket, key, ak, signature, policy string) string {
	// 构建完整的 POST URL
	postURL := fmt.Sprintf("%s/%s/", endpoint, bucket)

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OBS POST 上传示例</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 40px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            border-bottom: 3px solid #007bff;
            padding-bottom: 10px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: bold;
            color: #555;
        }
        input[type="text"], input[type="file"] {
            width: 100%%;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
            font-size: 14px;
        }
        input[readonly] {
            background-color: #f8f9fa;
            color: #6c757d;
        }
        button {
            background-color: #007bff;
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            font-weight: bold;
            width: 100%%;
        }
        button:hover {
            background-color: #0056b3;
        }
        h2 {
            color: #666;
            margin-top: 30px;
            font-size: 18px;
        }
        ul {
            background-color: #f8f9fa;
            padding: 20px;
            border-radius: 4px;
            line-height: 1.6;
        }
        li {
            margin-bottom: 8px;
        }
        .info-box {
            background-color: #d1ecf1;
            border: 1px solid #bee5eb;
            color: #0c5460;
            padding: 15px;
            border-radius: 4px;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>OBS POST 上传示例</h1>

        <div class="info-box">
            <strong>说明：</strong>文件将直接从浏览器上传到 OBS，不经过后端服务器。
        </div>

        <form action="%s" method="post" enctype="multipart/form-data">
            <div class="form-group">
                <label for="ak">AK (Access Key ID):</label>
                <input type="text" id="ak" name="AWSAccessKeyId" value="%s" readonly>
            </div>

            <div class="form-group">
                <label for="policy">Policy:</label>
                <input type="text" id="policy" name="policy" value="%s" readonly>
            </div>

            <div class="form-group">
                <label for="signature">Signature:</label>
                <input type="text" id="signature" name="signature" value="%s" readonly>
            </div>

            <div class="form-group">
                <label for="key">Key (对象名):</label>
                <input type="text" id="key" name="key" value="%s" readonly>
            </div>

            <div class="form-group">
                <label for="file">选择文件:</label>
                <input type="file" id="file" name="file" required>
            </div>

            <button type="submit">上传到 OBS</button>
        </form>

        <h2>注意事项：</h2>
        <ul>
            <li>策略和签名由后端服务器生成</li>
            <li>文件直接从浏览器上传到 OBS</li>
            <li>不经过后端服务器，降低服务器负载</li>
            <li>文件大小限制：10MB（由策略条件定义）</li>
            <li>策略过期时间：1 小时</li>
            <li>上传后的文件将公开可读（public-read）</li>
        </ul>

        <h2>技术细节：</h2>
        <ul>
            <li>POST URL: %s</li>
            <li>桶名称: %s</li>
            <li>对象名: %s</li>
        </ul>
    </div>
</body>
</html>`, postURL, ak, policy, signature, key, postURL, bucket, key)
}
