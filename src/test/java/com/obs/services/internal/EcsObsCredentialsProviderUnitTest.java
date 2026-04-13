/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
 */

package com.obs.services.internal;

import com.obs.services.EcsObsCredentialsProvider;
import org.junit.Assert;
import org.junit.Test;

public class EcsObsCredentialsProviderUnitTest {

    EcsObsCredentialsProvider ecsObsCredentialsProvider = new EcsObsCredentialsProvider();
    @Test
    public void tc_should_EcsObsCredentialsProvider_set_get_MetadataTokenTTLSeconds_successfully() {
        int testMetadataTokenTTLSeconds = 60;
        ecsObsCredentialsProvider.setMetadataTokenTTLSeconds(testMetadataTokenTTLSeconds);
        Assert.assertEquals(testMetadataTokenTTLSeconds, ecsObsCredentialsProvider.getMetadataTokenTTLSeconds());
    }
}
