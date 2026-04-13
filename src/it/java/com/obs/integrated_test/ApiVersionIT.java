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

package com.obs.integrated_test;

import com.obs.services.ObsClient;
import com.obs.services.internal.service.AbstractRequestConvertor;
import com.obs.services.model.AuthTypeEnum;
import com.obs.test.TestTools;
import com.obs.test.tools.PrepareTestBucket;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TestName;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.Locale;

import static org.junit.Assert.assertEquals;

public class ApiVersionIT {


    @Rule
    public TestName testName = new TestName();

    @Rule
    public PrepareTestBucket prepareTestBucket = new PrepareTestBucket();
    @Test
    public void test_while_bucketname_is_null() throws NoSuchMethodException, SecurityException, IllegalAccessException, IllegalArgumentException, InvocationTargetException {
        String bucketName = "";
        
        ObsClient obsClient = TestTools.getRequestPaymentEnvironment_User1();

        //        Class<?> clazz_2 = clazz_1.getSuperclass();
        Class<?> clazz_2 = AbstractRequestConvertor.class;
        
        Method method = clazz_2.getDeclaredMethod("getApiVersion", new Class[]{String.class});
        method.setAccessible(true);
        Object result = method.invoke(obsClient, new Object[]{bucketName});
        
        System.out.println(result);
        
        assertEquals(result, AuthTypeEnum.OBS);
    }
    
    @Test
    public void test_while_bucketname_not_null() throws NoSuchMethodException, SecurityException, IllegalAccessException, IllegalArgumentException, InvocationTargetException {
        String bucketName = testName.getMethodName().replace("_", "-").toLowerCase(Locale.ROOT);
        
        ObsClient obsClient = TestTools.getRequestPaymentEnvironment_User1();

        Class<?> clazz_2 = AbstractRequestConvertor.class;
        
        Method method = clazz_2.getDeclaredMethod("getApiVersion", new Class[]{String.class});
        method.setAccessible(true);
        Object result = method.invoke(obsClient, new Object[]{bucketName});
        
        System.out.println(result);
        
        assertEquals(result, AuthTypeEnum.OBS);
    }
}
