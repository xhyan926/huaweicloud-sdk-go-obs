/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.test.tools;

import com.obs.services.internal.Constants;

import okhttp3.Authenticator;
import okhttp3.Call;
import okhttp3.Credentials;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.InetSocketAddress;
import java.net.Proxy;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.cert.X509Certificate;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

public class TokenGetter
{
    private static String endPoint = "";

    /*
    IAM用户名
     */
    private static String userName = "";

    /*
    IAM用户的登录密码
     */
    private static String passWord = "";

    /*
    IAM用户所属帐号名称
     */
    private static String domainName = "";

    private static long durationSeconds = 3600L;

    static String token;

    private static final boolean proxyIsable = false;

    private static final String proxyHostAddress = "*** Provide your proxy host address ***";

    private static final int proxyPort = 8080;

    private static final String proxyUser = "*** Provide your proxy user ***";

    private static final String proxyPassword = "*** Provide your proxy password ***";


    //domainName is huawei account name
    public static void initCredentials(String IamEndPoint, String IamUserName, String IamUserPassWord, String HuaweiAccountName, long TokenDurationSeconds){
        endPoint = IamEndPoint;
        userName = IamUserName;
        passWord = IamUserPassWord;
        domainName = HuaweiAccountName;
        durationSeconds = TokenDurationSeconds;
    }
    public static void getToken() {
        Request.Builder builder = new Request.Builder();

        builder.addHeader("Content-Type", "application/json;charset=utf8");
        builder.url(endPoint + "/v3/auth/tokens");
        String mimeType = "application/json";
		/* request body sample
		   {
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"name": "***userName***",
							"password": "***passWord***",
							"domain": {
								"name": "***domainName***"
							}
						}
					}
				},
				"scope": {
					"domain": {
						"name": "***domainName***"
					}
				}
			}
		  }
		 */
        String content = "{\r\n" +
            "			\"auth\": {\r\n" +
            "				\"identity\": {\r\n" +
            "					\"methods\": [\"password\"],\r\n" +
            "					\"password\": {\r\n" +
            "						\"user\": {\r\n" +
            "							\"name\": \"" + userName + "\",\r\n" +
            "							\"password\": \"" + passWord + "\",\r\n" +
            "							\"domain\": {\r\n" +
            "								\"name\": \"" + domainName + "\"\r\n" +
            "							}\r\n" +
            "						}\r\n" +
            "					}\r\n" +
            "				},\r\n" +
            "				\"scope\": {\r\n" +
            "					\"domain\": {\r\n" +
            "						\"name\": \"" + domainName + "\"\r\n" +
            "					}\r\n" +
            "				}\r\n" +
            "			}\r\n" +
            "		  }";
        try {
            builder.post(createRequestBody(mimeType, content));
        } catch (UnsupportedEncodingException e1) {
            e1.printStackTrace();
        }
        try {
            getTokeResponse(builder.build());
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    private static void getTokeResponse(Request request) throws IOException {
        Call c = GetHttpClient().newCall(request);
        Response res = c.execute();
        String header = res.headers().toString();
        if (header.trim().equals("")) {
            // System.out.println("\n");
        } else {
            // System.out.println("getTokeResponse headers:"+header);
            // System.out.println("getTokeResponse body:"+res.body().string());
            String subjectToken = res.header("X-Subject-Token");
            // System.out.println("the Token :" + subjectToken);
            token = subjectToken;
        }
        res.close();
    }

    public static String getSecurityToken() {
        String token = TokenGetter.token;
        Request.Builder builder = new Request.Builder();
        builder.addHeader("Content-Type", "application/json;charset=utf8");
        builder.url(endPoint + "/v3.0/OS-CREDENTIAL/securitytokens");
        String mimeType = "application/json";

		/* request body sample
			 {
			    "auth": {
			        "identity": {
			            "methods": [
			                "token"
			            ],
			            "token": {
			                "id": "***yourToken***",
			                "duration_seconds": "***your-duration-seconds***"
			            }
			        }
			    }
			}
		 */
        String content = "{\r\n" +
            "    \"auth\": {\r\n" +
            "        \"identity\": {\r\n" +
            "            \"methods\": [\r\n" +
            "                \"token\"\r\n" +
            "            ],\r\n" +
            "            \"token\": {\r\n" +
            "                \"id\": \""+ token +"\",\r\n" +
            "                \"duration_seconds\": \""+ durationSeconds +"\"\r\n" +
            "\r\n" +
            "            }\r\n" +
            "        }\r\n" +
            "    }\r\n" +
            "}";

        try {
            builder.post(createRequestBody(mimeType, content));
        } catch (UnsupportedEncodingException e1) {
            e1.printStackTrace();
        }
        try {
            return getSecurityTokenResponse(builder.build());
        } catch (IOException e) {
            e.printStackTrace();
        }
        return "no SecurityToken";
    }

    public static String getAK(String tokenContent){
        return getTokenVal(tokenContent,"access");
    }
    public static String getSK(String tokenContent){
        return getTokenVal(tokenContent,"secret");
    }
    public static String getStsToken(String tokenContent){
        return getTokenVal(tokenContent,"securitytoken");
    }

    static String getTokenVal(String tokenContent, String variableName){
        String key = variableName+"\":\"";
        int startIndex = tokenContent.indexOf(key) + key.length();
        int endIndex = tokenContent.indexOf("\"", startIndex);
        String variable = tokenContent.substring(startIndex, endIndex);
        return variable;

    }

    private static String getSecurityTokenResponse(Request request) throws IOException {
        Call c = GetHttpClient().newCall(request);
        Response res = c.execute();
        String content = "";
        if (res.body() != null) {
            content = res.body().string();
            // System.out.println("getSecurityTokenResponse headers:\n" + res.headers() + "\n\n");
            if (content.trim().equals("")) {
                // System.out.println("\n");
            } else {
                // System.out.println("getSecurityTokenResponse body:\n" + content + "\n\n");
            }
        } else {
            // System.out.println("\n");
        }
        res.close();
        return content;
    }

    private static OkHttpClient GetHttpClient() {
        X509TrustManager xtm = new X509TrustManager() {
            @Override
            public void checkClientTrusted(X509Certificate[] chain, String authType) {
            }

            @Override
            public void checkServerTrusted(X509Certificate[] chain, String authType) {
            }

            @Override
            public X509Certificate[] getAcceptedIssuers() {
                X509Certificate[] x509Certificates = new X509Certificate[0];
                return x509Certificates;
            }
        };

        SSLContext sslContext = null;
        try {
            sslContext = SSLContext.getInstance("SSL");
            sslContext.init(null, new TrustManager[] { xtm }, new SecureRandom());
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        } catch (KeyManagementException e) {
            e.printStackTrace();
        }

        HostnameVerifier DO_NOT_VERIFY = (arg0, arg1) -> true;

        OkHttpClient.Builder builder = new OkHttpClient.Builder().followRedirects(false).retryOnConnectionFailure(false)
            .sslSocketFactory(sslContext.getSocketFactory(), xtm).hostnameVerifier(DO_NOT_VERIFY).cache(null);

        if(proxyIsable) {
            builder.proxy(new java.net.Proxy(Proxy.Type.HTTP, new InetSocketAddress(proxyHostAddress, proxyPort)));

            Authenticator proxyAuthenticator = (route, response) -> {
                String credential = Credentials.basic(proxyUser, proxyPassword);
                return response.request().newBuilder()
                    .header(Constants.CommonHeaders.PROXY_AUTHORIZATION, credential)
                    .build();
            };
            builder.proxyAuthenticator(proxyAuthenticator);
        }

        return builder.build();
    }

    private static RequestBody createRequestBody(String mimeType, String content) throws UnsupportedEncodingException {
        return RequestBody.create(MediaType.parse(mimeType), content.getBytes("UTF-8"));
    }
}
