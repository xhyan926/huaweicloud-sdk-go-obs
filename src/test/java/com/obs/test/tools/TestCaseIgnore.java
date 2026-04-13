/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.test.tools;

public class TestCaseIgnore {
    public static boolean needIgnore() {
        StackTraceElement[] stackTrace = Thread.currentThread().getStackTrace();
        StackTraceElement element = stackTrace[2]; // 获取第三个堆栈元素（当前测试用例名）
        String testCaseName = element.getClassName() + "_" + element.getMethodName();
        testCaseName = testCaseName.replace('.', '_');
        return  "ignore".equals(System.getenv(testCaseName));
    }
}
