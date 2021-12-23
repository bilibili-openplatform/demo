package com.example.demo.util;

import org.apache.http.HttpEntity;
import org.apache.http.HttpStatus;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.ByteArrayEntity;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.StringEntity;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;

import java.io.File;
import java.net.URI;
import java.nio.charset.StandardCharsets;

public class HttpUtil {

    public static String doPost(URI uri) throws Exception {
        CloseableHttpClient httpClient = HttpClients.createDefault();
        CloseableHttpResponse response = null;
        try {
            HttpPost httpPost = new HttpPost(uri);
            response = httpClient.execute(httpPost);
            if (response.getStatusLine().getStatusCode() != HttpStatus.SC_OK) {
                throw new Exception("请求失败...");
            }
            return EntityUtils.toString(response.getEntity(), "UTF-8");
        } finally {
            if (response != null) {
                response.close();
            }
            httpClient.close();
        }
    }

    public static String doPostJson(URI uri, String json) throws Exception {
        CloseableHttpClient httpClient = HttpClients.createDefault();
        CloseableHttpResponse response = null;
        try {
            HttpPost httpPost = new HttpPost(uri);
            httpPost.setHeader("Content-Type", "application/json;charset=utf-8");
            StringEntity entity = new StringEntity(json, StandardCharsets.UTF_8);
            entity.setContentType("application/json;charset=utf-8");
            httpPost.setEntity(entity);

            response = httpClient.execute(httpPost);
            if (response.getStatusLine().getStatusCode() != HttpStatus.SC_OK) {
                throw new Exception("请求失败...");
            }
            return EntityUtils.toString(response.getEntity(), "UTF-8");
        } finally {
            if (response != null) {
                response.close();
            }
            httpClient.close();
        }
    }

    public static String doPostFile(URI uri, File file) throws Exception {
        CloseableHttpClient httpClient = HttpClients.createDefault();
        CloseableHttpResponse response = null;
        try {
            HttpPost httpPost = new HttpPost(uri);
            FileBody bin = new FileBody(file);
            HttpEntity reqEntity = MultipartEntityBuilder.create().addPart("file", bin).build();
            httpPost.setEntity(reqEntity);

            response = httpClient.execute(httpPost);
            if (response.getStatusLine().getStatusCode() != HttpStatus.SC_OK) {
                throw new Exception("请求失败...");
            }
            return EntityUtils.toString(response.getEntity(), "UTF-8");
        } finally {
            if (response != null) {
                response.close();
            }
            httpClient.close();
        }
    }

    public static String doPostStream(URI uri, byte[] bytes) throws Exception {
        CloseableHttpClient httpClient = HttpClients.createDefault();
        CloseableHttpResponse response = null;
        try {
            HttpPost httpPost = new HttpPost(uri);
            httpPost.setHeader("Content-Type", ContentType.APPLICATION_OCTET_STREAM.getMimeType());
            httpPost.setEntity(new ByteArrayEntity(bytes));
            response = httpClient.execute(httpPost);
            if (response.getStatusLine().getStatusCode() != HttpStatus.SC_OK) {
                throw new Exception("请求失败...");
            }
            return EntityUtils.toString(response.getEntity(), "UTF-8");
        } finally {
            if (response != null) {
                response.close();
            }
            httpClient.close();
        }
    }
}
