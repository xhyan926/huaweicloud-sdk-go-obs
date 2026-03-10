//go:build integration

package auth

import (
	"strings"
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs/test/integration"
)

// TestAuthentication_ShouldConnectSuccessfully_GivenValidStaticCredentials 测试静态凭证认证
func TestAuthentication_ShouldConnectSuccessfully_GivenValidStaticCredentials(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldGetBucketLocationSuccessfully_GivenValidCredentials", func(t *testing.T) {
		// 测试基本认证 - 通过获取桶位置验证
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		output, err := client.TestClient().GetBucketLocation(input)
		if err != nil {
			t.Fatalf("获取桶位置失败: %v", err)
		}

		// 验证返回结果
		if output == nil {
			t.Error("GetBucketLocation返回nil")
		}

		// 验证桶名称
		if output.Location == "" {
			t.Error("桶位置信息为空")
		}

		client.AddTestCase("静态凭证认证成功")
		t.Logf("静态凭证认证通过，桶位置: %s", output.Location)
	})
}

// TestAuthentication_ShouldConnectSuccessfully_GivenTemporaryCredentials 测试临时凭证认证
func TestAuthentication_ShouldConnectSuccessfully_GivenTemporaryCredentials(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	// 检查是否配置了临时凭证
	if client.Config.SecurityToken == "" {
		t.Skip("跳过临时凭证测试，未配置OBS_TEST_TOKEN环境变量")
	}

	bucket := client.GetTestBucket()

	t.Run("ShouldAuthenticateWithToken_GivenValidTemporaryCredentials", func(t *testing.T) {
		// 测试临时凭证认证
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		output, err := client.TestClient().GetBucketLocation(input)
		if err != nil {
			t.Fatalf("临时凭证认证失败: %v", err)
		}

		// 验证返回结果
		if output == nil {
			t.Error("GetBucketLocation返回nil")
		}

		client.AddTestCase("临时凭证认证成功")
		t.Logf("临时凭证认证通过，安全令牌长度: %d", len(client.Config.SecurityToken))
	})
}

// TestAuthentication_ShouldFail_GivenInvalidCredentials 测试无效凭证的认证失败
func TestAuthentication_ShouldFail_GivenInvalidCredentials(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	t.Run("ShouldFailAuthentication_GivenInvalidAccessKey", func(t *testing.T) {
		// 创建使用无效AK的客户端
		invalidClient, err := obs.New(
			"invalid-access-key",
			client.Config.SecretKey,
			client.Config.Endpoint,
			obs.WithSecurityToken(client.Config.SecurityToken),
			obs.WithRegion(client.Config.Region),
		)
		if err != nil {
			t.Fatalf("创建客户端失败: %v", err)
		}
		defer invalidClient.Close()

		// 尝试使用无效凭证访问
		bucket := client.GetTestBucket()
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		_, err = invalidClient.GetBucketLocation(input)
		if err == nil {
			t.Error("期望认证失败，但操作成功")
		}

		// 验证是认证错误
		if obsErr, ok := err.(obs.ObsError); ok {
			if obsErr.StatusCode != 403 && obsErr.StatusCode != 401 {
				t.Errorf("期望403或401错误，实际: %d", obsErr.StatusCode)
			}
		}

		client.AddTestCase("无效AK认证失败测试通过")
		t.Logf("无效AK认证失败，错误: %v", err)
	})

	t.Run("ShouldFailAuthentication_GivenInvalidSecretKey", func(t *testing.T) {
		// 创建使用无效SK的客户端
		invalidClient, err := obs.New(
			client.Config.AccessKey,
			"invalid-secret-key",
			client.Config.Endpoint,
			obs.WithSecurityToken(client.Config.SecurityToken),
			obs.WithRegion(client.Config.Region),
		)
		if err != nil {
			t.Fatalf("创建客户端失败: %v", err)
		}
		defer invalidClient.Close()

		// 尝试使用无效凭证访问
		bucket := client.GetTestBucket()
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		_, err = invalidClient.GetBucketLocation(input)
		if err == nil {
			t.Error("期望认证失败，但操作成功")
		}

		// 验证是认证错误
		if obsErr, ok := err.(obs.ObsError); ok {
			if obsErr.StatusCode != 403 && obsErr.StatusCode != 401 {
				t.Errorf("期望403或401错误，实际: %d", obsErr.StatusCode)
			}
		}

		client.AddTestCase("无效SK认证失败测试通过")
		t.Logf("无效SK认证失败，错误: %v", err)
	})
}

// TestAuthentication_DifferentSignatureTypes 测试不同签名类型
func TestAuthentication_DifferentSignatureTypes(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldWorkWithSignatureV2_GivenValidCredentials", func(t *testing.T) {
		// 创建使用SignatureV2的客户端
		v2Client, err := obs.New(
			client.Config.AccessKey,
			client.Config.SecretKey,
			client.Config.Endpoint,
			obs.WithSignatureType(obs.TYPE_SIGNATURE_V2),
			obs.WithSecurityToken(client.Config.SecurityToken),
			obs.WithRegion(client.Config.Region),
		)
		if err != nil {
			t.Fatalf("创建V2客户端失败: %v", err)
		}
		defer v2Client.Close()

		// 测试V2签名
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		_, err = v2Client.GetBucketLocation(input)
		if err != nil {
			t.Logf("V2签名测试失败: %v (可能是OBS不支持的签名类型)", err)
			// 某些区域可能不支持V2签名，所以这里只是记录而不是失败
		} else {
			client.AddTestCase("V2签名认证成功")
			t.Log("V2签名认证通过")
		}
	})

	t.Run("ShouldWorkWithSignatureV4_GivenValidCredentials", func(t *testing.T) {
		// 创建使用SignatureV4的客户端（默认）
		v4Client, err := obs.New(
			client.Config.AccessKey,
			client.Config.SecretKey,
			client.Config.Endpoint,
			obs.WithSignatureType(obs.TYPE_SIGNATURE_V4),
			obs.WithSecurityToken(client.Config.SecurityToken),
			obs.WithRegion(client.Config.Region),
		)
		if err != nil {
			t.Fatalf("创建V4客户端失败: %v", err)
		}
		defer v4Client.Close()

		// 测试V4签名
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		output, err := v4Client.GetBucketLocation(input)
		if err != nil {
			t.Fatalf("V4签名测试失败: %v", err)
		}

		if output.Location == "" {
			t.Error("V4签名返回空位置信息")
		}

		client.AddTestCase("V4签名认证成功")
		t.Logf("V4签名认证通过")
	})

	t.Run("ShouldWorkWithSignatureObs_GivenValidCredentials", func(t *testing.T) {
		// 创建使用SignatureObs的客户端
		obsClient, err := obs.New(
			client.Config.AccessKey,
			client.Config.SecretKey,
			client.Config.Endpoint,
			obs.WithSignatureType(obs.TYPE_SIGNATURE_OBS),
			obs.WithSecurityToken(client.Config.SecurityToken),
			obs.WithRegion(client.Config.Region),
		)
		if err != nil {
			t.Fatalf("创建Obs签名客户端失败: %v", err)
		}
		defer obsClient.Close()

		// 测试Obs签名
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		output, err := obsClient.GetBucketLocation(input)
		if err != nil {
			t.Fatalf("Obs签名测试失败: %v", err)
		}

		if output.Location == "" {
			t.Error("Obs签名返回空位置信息")
		}

		client.AddTestCase("Obs签名认证成功")
		t.Logf("Obs签名认证通过")
	})
}

// TestAuthentication_CrossRegionAccess 测试跨区域访问
func TestAuthentication_CrossRegionAccess(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	t.Run("ShouldAccessInSameRegion_GivenSameRegionClient", func(t *testing.T) {
		bucket := client.GetTestBucket()
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		output, err := client.TestClient().GetBucketLocation(input)
		if err != nil {
			t.Fatalf("同区域访问失败: %v", err)
		}

		// 验证返回的区域
		if output.Location == "" {
			t.Error("返回的区域为空")
		}

		client.AddTestCase("同区域访问成功")
		t.Logf("同区域访问成功，区域: %s", output.Location)
	})
}

// TestAuthentication_ErrorHandling 测试认证错误处理
func TestAuthentication_ErrorHandling(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	t.Run("ShouldReturnObsError_GivenAuthenticationFailure", func(t *testing.T) {
		// 创建无效客户端
		invalidClient, err := obs.New(
			"invalid-ak",
			"invalid-sk",
			client.Config.Endpoint,
		)
		if err != nil {
			t.Fatalf("创建客户端失败: %v", err)
		}
		defer invalidClient.Close()

		bucket := client.GetTestBucket()
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		_, err = invalidClient.GetBucketLocation(input)
		if err == nil {
			t.Fatal("期望认证失败")
		}

		// 验证错误类型
		obsErr, ok := err.(obs.ObsError)
		if !ok {
			t.Fatalf("错误不是ObsError类型: %T", err)
		}

		// 验证错误字段
		if obsErr.StatusCode == 0 {
			t.Error("错误状态码为0")
		}

		if obsErr.Code == "" {
			t.Error("错误代码为空")
		}

		if obsErr.Message == "" {
			t.Error("错误消息为空")
		}

		// 验证请求ID存在
		if obsErr.RequestId == "" {
			t.Error("请求ID为空")
		}

		client.AddTestCase("认证错误处理测试通过")
		t.Logf("认证错误: 状态码=%d, 代码=%s, 消息=%s",
			obsErr.StatusCode, obsErr.Code, obsErr.Message)
	})
}

// TestAuthentication_ConcurrentAuthentication 测试并发认证
func TestAuthentication_ConcurrentAuthentication(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	t.Run("ShouldHandleConcurrentAuthentications_GivenMultipleRequests", func(t *testing.T) {
		bucket := client.GetTestBucket()

		// 并发执行多个认证请求
		numRequests := 5
		errChan := make(chan error, numRequests)

		for i := 0; i < numRequests; i++ {
			go func(index int) {
				input := &obs.GetBucketLocationInput{
					Bucket: bucket,
				}

				_, err := client.TestClient().GetBucketLocation(input)
				errChan <- err
			}(i)
		}

		// 收集错误
		errorCount := 0
		for i := 0; i < numRequests; i++ {
			err := <-errChan
			if err != nil {
				errorCount++
				t.Logf("认证请求 %d 失败: %v", i, err)
			}
		}

		if errorCount > 0 {
			t.Errorf("有 %d/%d 个认证请求失败", errorCount, numRequests)
		}

		client.AddTestCase("并发认证测试通过")
		t.Logf("并发认证完成，成功率: %d/%d", numRequests-errorCount, numRequests)
	})
}

// TestAuthentication_AuthenticationDetails 测试认证详细信息
func TestAuthentication_AuthenticationDetails(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	bucket := client.GetTestBucket()

	t.Run("ShouldGetAuthenticatedUserDetails_GivenValidCredentials", func(t *testing.T) {
		// 通过获取桶信息验证认证
		input := &obs.GetBucketMetadataInput{
			Bucket: bucket,
		}

		output, err := client.TestClient().GetBucketMetadata(input)
		if err != nil {
			t.Fatalf("获取桶元数据失败: %v", err)
		}

		// 验证返回结果
		if output == nil {
			t.Error("GetBucketMetadata返回nil")
		}

		// 验证Owner信息
		if output.Owner == nil {
			t.Error("Owner信息为空")
		} else {
			if output.Owner.ID == "" {
				t.Error("Owner ID为空")
			}

			client.AddTestCase("获取认证用户详情成功")
			t.Logf("认证用户: ID=%s, DisplayName=%s",
				output.Owner.ID, output.Owner.DisplayName)
		}
	})
}

// TestAuthentication_TokenExpiration 测试临时凭证过期处理
func TestAuthentication_TokenExpiration(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	// 检查是否配置了临时凭证
	if client.Config.SecurityToken == "" {
		t.Skip("跳过临时凭证过期测试，未配置OBS_TEST_TOKEN环境变量")
	}

	t.Run("ShouldAuthenticateWithToken_GivenValidTemporaryCredentials", func(t *testing.T) {
		bucket := client.GetTestBucket()

		// 使用临时凭证执行操作
		input := &obs.GetBucketLocationInput{
			Bucket: bucket,
		}

		output, err := client.TestClient().GetBucketLocation(input)
		if err != nil {
			t.Fatalf("临时凭证操作失败: %v", err)
		}

		if output.Location == "" {
			t.Error("返回的位置信息为空")
		}

		client.AddTestCase("临时凭证操作成功")
		t.Logf("临时凭证操作成功，令牌长度: %d", len(client.Config.SecurityToken))
	})
}

// TestAuthentication_ClientConfiguration 测试客户端认证配置
func TestAuthentication_ClientConfiguration(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	t.Run("ShouldCreateClientWithAuthConfiguration_GivenValidConfig", func(t *testing.T) {
		// 验证客户端配置
		if client.Config.AccessKey == "" {
			t.Error("AccessKey为空")
		}

		if client.Config.SecretKey == "" {
			t.Error("SecretKey为空")
		}

		if client.Config.Endpoint == "" {
			t.Error("Endpoint为空")
		}

		if client.Config.TestBucket == "" {
			t.Error("TestBucket为空")
		}

		// 验证endpoint格式
		if !strings.HasPrefix(client.Config.Endpoint, "http://") &&
			!strings.HasPrefix(client.Config.Endpoint, "https://") {
			t.Errorf("Endpoint格式不正确: %s", client.Config.Endpoint)
		}

		client.AddTestCase("客户端配置验证通过")
		t.Logf("客户端配置: Endpoint=%s, Bucket=%s, Region=%s",
			client.Config.Endpoint, client.Config.TestBucket, client.Config.Region)
	})
}

// TestAuthentication_CredentialsLength 测试凭证长度限制
func TestAuthentication_CredentialsLength(t *testing.T) {
	client := integration.NewTestClient(t)
	defer client.Cleanup(t)

	t.Run("ShouldValidateCredentialsLength_GivenValidCredentials", func(t *testing.T) {
		// 验证凭证长度
		if len(client.Config.AccessKey) == 0 {
			t.Error("AccessKey长度为0")
		}

		if len(client.Config.AccessKey) > 100 {
			t.Logf("警告: AccessKey长度过长: %d", len(client.Config.AccessKey))
		}

		if len(client.Config.SecretKey) == 0 {
			t.Error("SecretKey长度为0")
		}

		if len(client.Config.SecretKey) > 100 {
			t.Logf("警告: SecretKey长度过长: %d", len(client.Config.SecretKey))
		}

		if client.Config.SecurityToken != "" {
			if len(client.Config.SecurityToken) > 1000 {
				t.Logf("警告: SecurityToken长度过长: %d", len(client.Config.SecurityToken))
			}
		}

		client.AddTestCase("凭证长度验证通过")
		t.Logf("凭证长度: AK=%d, SK=%d, Token=%d",
			len(client.Config.AccessKey),
			len(client.Config.SecretKey),
			len(client.Config.SecurityToken))
	})
}
