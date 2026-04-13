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

package com.obs.services;


import com.obs.test.objects.BaseObjectTest;
import org.junit.Assert;
import org.junit.Test;

import java.io.IOException;

public class DefaultCredentialsProviderChainTest extends BaseObjectTest {

    @Test
    public void should_add_all_credentialsProviders_in_this_credentialsProviders() throws IOException {
        IObsCredentialsProvider provider1 = new BasicObsCredentialsProvider("ak1", "sk1");
        IObsCredentialsProvider provider2 = new BasicObsCredentialsProvider("ak2", "sk2");
        DefaultCredentialsProviderChain chain = new DefaultCredentialsProviderChain(provider1, provider2);
        Assert.assertEquals("provider1 error", provider1.getSecurityKey().getAccessKey(),
            chain.getSecurityKey().getAccessKey());
    }



}
