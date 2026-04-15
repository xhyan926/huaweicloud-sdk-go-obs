/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026-2026. All rights reserved.
 */

package com.obs.integrated_test.objectlock;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;

import com.obs.services.ObsClient;
import com.obs.services.model.AuthTypeEnum;
import com.obs.services.model.BucketVersioningConfiguration;
import com.obs.services.model.HeaderResponse;
import com.obs.services.model.VersioningStatusEnum;
import com.obs.services.model.objectlock.DefaultRetention;
import com.obs.services.model.objectlock.GetObjectLockConfigurationRequest;
import com.obs.services.model.objectlock.GetObjectLockConfigurationResult;
import com.obs.services.model.objectlock.ObjectLockConfiguration;
import com.obs.services.model.objectlock.ObjectLockRule;
import com.obs.services.model.objectlock.SetObjectLockConfigurationRequest;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;

import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

import java.util.Arrays;
import java.util.Collection;
import java.util.Locale;

@RunWith(Parameterized.class)
public class ObjectLockConfigurationIT {
    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();

    @Parameterized.Parameter()
    public String authTypeName;

    @Parameterized.Parameter(1)
    public AuthTypeEnum authType;

    @Parameterized.Parameters(name = "{0}")
    public static Collection<Object[]> authTypeData() {
        return Arrays.asList(new Object[][] {
            {"OBS", AuthTypeEnum.OBS},
            {"V2", AuthTypeEnum.V2}
        });
    }

    private ObsClient obsClient;
    private String bucketName;

    @Before
    public void setUp() {
        if (authType == AuthTypeEnum.OBS) {
            obsClient = TestTools.getPipelineEnvironment();
        } else {
            obsClient = TestTools.getPipelineEnvironment_V2();
        }
        Assert.assertNotNull("ObsClient should not be null", obsClient);
        bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT)
            .replace("[", "").replace("]", "");
    }

    @After
    public void tearDown() {
        // 桶的创建和删除由 PrepareTestBucket @Rule 统一管理
    }

    // SET + GET + 边界值测试
    @Test
    public void test_SDK_objectlock_001() {
        // 开启多版本
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // SET: 配置 Days 模式
        DefaultRetention retention = new DefaultRetention("COMPLIANCE", 30, null);
        ObjectLockRule rule = new ObjectLockRule(retention);
        ObjectLockConfiguration config = new ObjectLockConfiguration("Enabled", rule);
        SetObjectLockConfigurationRequest setRequest =
            new SetObjectLockConfigurationRequest(bucketName, config);
        HeaderResponse headerResponse = obsClient.setObjectLockConfiguration(setRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());

        // GET: 验证配置
        GetObjectLockConfigurationRequest getRequest =
            new GetObjectLockConfigurationRequest(bucketName);
        GetObjectLockConfigurationResult result = obsClient.getObjectLockConfiguration(getRequest);
        assertNotNull(result.getObjectLockConfiguration());
        assertEquals("Enabled", result.getObjectLockConfiguration().getObjectLockEnabled());
        assertNotNull(result.getObjectLockConfiguration().getRule());
        assertEquals("COMPLIANCE",
            result.getObjectLockConfiguration().getRule().getDefaultRetention().getMode());
        assertEquals(Integer.valueOf(30),
            result.getObjectLockConfiguration().getRule().getDefaultRetention().getDays());

        // SET: 配置 Years 模式
        retention = new DefaultRetention("COMPLIANCE", null, 1);
        rule = new ObjectLockRule(retention);
        config = new ObjectLockConfiguration("Enabled", rule);
        setRequest = new SetObjectLockConfigurationRequest(bucketName, config);
        headerResponse = obsClient.setObjectLockConfiguration(setRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());

        // GET: 验证 Years 配置
        result = obsClient.getObjectLockConfiguration(getRequest);
        assertNotNull(result.getObjectLockConfiguration());
        assertEquals(Integer.valueOf(1),
            result.getObjectLockConfiguration().getRule().getDefaultRetention().getYears());
        assertEquals(Integer.valueOf(0),
            result.getObjectLockConfiguration().getRule().getDefaultRetention().getDays());
    }

    // SET + GET + 清空配置 + 再 GET
    @Test
    public void test_SDK_objectlock_002() {
        // 开启多版本
        obsClient.setBucketVersioning(bucketName,
            new BucketVersioningConfiguration(VersioningStatusEnum.ENABLED));

        // SET: 配置 Days
        DefaultRetention retention = new DefaultRetention("COMPLIANCE", 10, null);
        ObjectLockRule rule = new ObjectLockRule(retention);
        ObjectLockConfiguration config = new ObjectLockConfiguration("Enabled", rule);
        SetObjectLockConfigurationRequest setRequest =
            new SetObjectLockConfigurationRequest(bucketName, config);
        HeaderResponse headerResponse = obsClient.setObjectLockConfiguration(setRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());

        // GET: 验证配置
        GetObjectLockConfigurationRequest getRequest =
            new GetObjectLockConfigurationRequest(bucketName);
        GetObjectLockConfigurationResult result = obsClient.getObjectLockConfiguration(getRequest);
        assertNotNull(result.getObjectLockConfiguration().getRule());

        // SET: 清空配置（不携带 Rule）
        config = new ObjectLockConfiguration("Enabled", null);
        setRequest = new SetObjectLockConfigurationRequest(bucketName, config);
        headerResponse = obsClient.setObjectLockConfiguration(setRequest);
        Assert.assertEquals(200, headerResponse.getStatusCode());

        // GET: 验证配置已清空
        result = obsClient.getObjectLockConfiguration(getRequest);
        assertNotNull(result.getObjectLockConfiguration());
        assertEquals("Enabled", result.getObjectLockConfiguration().getObjectLockEnabled());
        assertNull(result.getObjectLockConfiguration().getRule());
    }
}
