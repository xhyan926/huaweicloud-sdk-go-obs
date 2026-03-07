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

// main/post_upload_sample.go
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
//   2. 运行示例：go run main/post_upload_sample.go
//   3. 使用生成的 HTML 表单上传文件
//
// 环境变量：
//   OBS_AK: 你的 Access Key ID
//   OBS_SK: 你的 Secret Access Key
//   OBS_ENDPOINT: OBS 端点，例如 https://obs.cn-north-1.myhuaweicloud.com
//
// 注意事项：
//   - Policy 和签名由后端生成，前端只使用生成的凭证
//   - 策略有过期时间，过期后需要重新生成
//   - 对于高级用法（如自定义条件、文件大小限制），请使用 CreateBrowserBasedSignature 接口
//

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
	bucketName := "your-bucket-name" // 替换为你的桶名
	objectKey := "uploads/test.jpg"  // 上传后的对象名
	expiresIn := int64(3600)      // 1 小时后过期

	// 4. 创建 POST 上传策略（简化版本）
	input := &obs.CreatePostPolicyInput{
		Bucket: bucketName,
		Key:    objectKey,
		Expires: expiresIn,
		Acl:     "public-read", // 可选：设置对象 ACL
	}

	// 5. 生成策略和签名
	output, err := obsClient.CreatePostPolicy(input)
	if err != nil {
		log.Fatalf("创建 POST 策略失败: %v", err)
	}

	// 6. 输出策略信息
	fmt.Println("=== POST 上传策略已生成 ===")
	fmt.Printf("Policy (Base64): %s\n", output.Policy)
	fmt.Printf("Signature: %s\n", output.Signature)
	fmt.Println()

	// 7. 生成 HTML 表单
	htmlForm := generateHTMLForm(bucketName, objectKey, output.Policy, output.Signature, endpoint)

	// 8. 保存 HTML 表单到文件
	htmlFile := "post_upload_form.html"
	err = os.WriteFile(htmlFile, []byte(htmlForm), 0644)
	if err != nil {
		log.Fatalf("保存 HTML 文件失败: %v", err)
	}

	fmt.Printf("HTML 表单已保存到: %s\n", htmlFile)
	fmt.Println("请在浏览器中打开该文件并上传测试文件")
	fmt.Println()
	fmt.Println("=== 高级用法 ===")
	fmt.Println("如需自定义条件（如文件大小限制、内容类型限制等），请使用 CreateBrowserBasedSignature 接口：")
	fmt.Println("  input := &obs.CreateBrowserBasedSignatureInput{")
	fmt.Println("    Bucket: bucketName,")
	fmt.Println("    Key: objectKey,")
	fmt.Println("    FormParams: map[string]string{\"content-type\": \"image/jpeg\"},")
	fmt.Println("    RangeParams: []obs.RangeParams{")
	fmt.Println("      {RangeName: \"content-length-range\", Lower: 1, Upper: 10 * 1024 * 1024},")
	fmt.Println("    },")
	fmt.Println("    Expires: 3600,")
	fmt.Println("  }")
	fmt.Println("  output, _ := obsClient.CreateBrowserBasedSignature(input)")
}

// generateHTMLForm 生成 HTML 上传表单
func generateHTMLForm(bucket, objectKey, policy, signature, endpoint string) string {
	// 从 endpoint 中提取主机名
	host := ""
	if len(endpoint) > 8 {
		host = endpoint[8:] // 去掉 "https://"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>OBS POST 上传示例</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            border-bottom: 2px solid #4CAF50;
            padding-bottom: 10px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        input[type="file"] {
            width: 100%%;
            padding: 10px;
            border: 2px solid #ddd;
            border-radius: 4px;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover {
            background-color: #45a049;
        }
        .info {
            background-color: #e7f3ff;
            border-left: 4px solid #2196F3;
            padding: 15px;
            margin-bottom: 20px;
        }
        code {
            background-color: #f4f4f4;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            word-break: break-all;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>OBS POST 上传示例</h1>

        <div class="info">
            <p><strong>上传配置：</strong></p>
            <p>桶名称: <code>%s</code></p>
            <p>对象键: <code>%s</code></p>
            <p>Endpoint: <code>%s</code></p>
        </div>

        <form action="https://%s/" method="post" enctype="multipart/form-data">
            <div class="form-group">
                <label for="key">对象键（文件名）:</label>
                <input type="text" name="key" value="%s" readonly />
            </div>

            <div class="form-group">
                <label for="acl">ACL（可选）:</label>
                <select name="acl">
                    <option value="private">private</option>
                    <option value="public-read" selected>public-read</option>
                    <option value="public-read-write">public-read-write</option>
                </select>
            </div>

            <div class="form-group" style="display: none;">
                <label>Policy (Base64):</label>
                <input type="text" name="policy" value="%s" readonly />
            </div>

            <div class="form-group" style="display: none;">
                <label>Signature:</label>
                <input type="text" name="signature" value="%s" readonly />
            </div>

            <div class="form-group">
                <label for="file">选择文件:</label>
                <input type="file" name="file" required />
            </div>

            <div class="form-group">
                <button type="submit">上传文件</button>
            </div>
        </form>

        <div class="info">
            <h3>说明：</h3>
            <ul>
                <li>Policy 和 Signature 已自动填充到表单中</li>
                <li>可以直接点击"上传文件"按钮上传任意文件</li>
                <li>上传成功后，文件将保存在指定桶的指定路径下</li>
                <li>策略有过期时间（1小时），过期后需要重新生成</li>
            </ul>
        </div>

        <div class="info">
            <h3>高级用法：</h3>
            <p>如需自定义上传条件（如文件大小限制、内容类型限制），请使用 <code>CreateBrowserBasedSignature</code> 接口。</p>
        </div>
    </div>
</body>
</html>`,
		bucket, objectKey, host, objectKey, policy, signature)
}
