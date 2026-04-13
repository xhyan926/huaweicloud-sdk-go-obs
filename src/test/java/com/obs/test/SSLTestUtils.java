/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.test;

import javax.net.ssl.X509TrustManager;
import java.security.cert.X509Certificate;

/**
 * SSL测试工具类
 */
public class SSLTestUtils {

    /**
     * 信任所有证书的TrustManager，用于测试环境
     */
    public static X509TrustManager trustAllManager = new X509TrustManager() {
        @Override
        public void checkClientTrusted(X509Certificate[] chain, String authType) {
            // 客户端证书验证
        }

        @Override
        public void checkServerTrusted(X509Certificate[] chain, String authType) {
            // 服务端证书验证
        }

        @Override
        public X509Certificate[] getAcceptedIssuers() {
            return new X509Certificate[0];
        }
    };
}