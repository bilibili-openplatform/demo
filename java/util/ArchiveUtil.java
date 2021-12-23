package com.example.demo.util;

import com.example.demo.entity.ArcAddParam;
import com.example.demo.entity.ArcBizException;
import com.example.demo.entity.CommonReply;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.http.client.utils.URIBuilder;

import java.io.File;
import java.net.URI;
import java.util.HashMap;
import java.util.Map;

public class ArchiveUtil {

    private String clientId;
    private String secret;

    public ArchiveUtil(String clientId, String secret) {
        this.clientId = clientId;
        this.secret = secret;
    }

    // apis
    private final static String url_token = "https://api.bilibili.com/x/account-oauth2/v1/token";
    private final static String url_arc_init = "https://member.bilibili.com/arcopen/fn/archive/video/init";
    private final static String url_arc_up = "https://openupos.bilivideo.com/video/v2/part/upload";
    private final static String url_arc_complete = "https://member.bilibili.com/arcopen/fn/archive/video/complete";
    private final static String url_arc_cover_up = "https://member.bilibili.com/arcopen/fn/archive/cover/upload";
    private final static String url_arc_submit = "https://member.bilibili.com/arcopen/fn/archive/add-by-utoken";

    /**
     * 获取 accessToken
     *
     * @param code 授权拿到的临时票据 code
     * @return upload_token
     */
    public String getAccessToken(String code) throws Exception {
        URI uri = new URIBuilder(url_token)
                .setParameter("client_id", clientId)
                .setParameter("client_secret", secret)
                .setParameter("grant_type", "authorization_code")
                .setParameter("code", code)
                .build();
        ObjectMapper mapper = new ObjectMapper();
        String res = HttpUtil.doPost(uri);
        CommonReply reply = mapper.readValue(res, CommonReply.class);
        if (reply.getCode() != 0) {
            throw new ArcBizException(reply);
        }
        return (String) reply.getData().get("access_token");
    }

    /**
     * 文件上传预处理
     *
     * @param accessToken access_token
     * @param name        文件名字，需携带正确的扩展名，例如test.mp4
     * @return upload_token
     */
    public String arcInit(String accessToken, String name) throws Exception {
        URI uri = new URIBuilder(url_arc_init)
                .setParameter("client_id", clientId)
                .setParameter("access_token", accessToken)
                .build();
        ObjectMapper mapper = new ObjectMapper();
        Map<String, String> param = new HashMap<>();
        param.put("name", name);
        String res = HttpUtil.doPostJson(uri, mapper.writeValueAsString(param));
        CommonReply reply = mapper.readValue(res, CommonReply.class);
        if (reply.getCode() != 0) {
            throw new ArcBizException(reply);
        }
        return (String) reply.getData().get("upload_token");
    }

    /**
     * 文件分片上传
     *
     * @param uploadToken upload_token
     * @param partNum     分片编号
     * @param bytes       字节数组
     */
    public void arcUp(String uploadToken, int partNum, byte[] bytes) throws Exception {
        URI uri = new URIBuilder(url_arc_up)
                .setParameter("upload_token", uploadToken)
                .setParameter("part_number", partNum + "")
                .build();
        ObjectMapper mapper = new ObjectMapper();
        String res = HttpUtil.doPostStream(uri, bytes);
        CommonReply reply = mapper.readValue(res, CommonReply.class);
        if (reply.getCode() != 0) {
            throw new ArcBizException(reply);
        }
    }

    /**
     * 文件分片合片
     *
     * @param uploadToken upload_token
     */
    public void arcComplete(String uploadToken) throws Exception {
        URI uri = new URIBuilder(url_arc_complete)
                .setParameter("upload_token", uploadToken)
                .build();
        ObjectMapper mapper = new ObjectMapper();
        String res = HttpUtil.doPostJson(uri, "");
        CommonReply reply = mapper.readValue(res, CommonReply.class);
        if (reply.getCode() != 0) {
            throw new ArcBizException(reply);
        }
    }


    /**
     * 稿件封面上传
     *
     * @param accessToken access_token
     * @param file        封面图片文件
     * @return 封面图片地址
     */
    public String uploadCover(String accessToken, File file) throws Exception {
        URI uri = new URIBuilder(url_arc_cover_up)
                .setParameter("client_id", clientId)
                .setParameter("access_token", accessToken)
                .build();
        String res = HttpUtil.doPostFile(uri, file);
        ObjectMapper mapper = new ObjectMapper();
        CommonReply reply = mapper.readValue(res, CommonReply.class);
        if (reply.getCode() != 0) {
            throw new ArcBizException(reply);
        }
        return (String) reply.getData().get("url");
    }

    /**
     * 视频稿件提交
     *
     * @param accessToken access_token
     * @param arcAdd      稿件提交参数
     * @return 稿件ID
     */
    public String arcSubmit(String accessToken, String uploadToken, ArcAddParam arcAdd) throws Exception {
        URI uri = new URIBuilder(url_arc_submit)
                .setParameter("client_id", clientId)
                .setParameter("access_token", accessToken)
                .setParameter("upload_token", uploadToken)
                .build();
        ObjectMapper mapper = new ObjectMapper();
        String res = HttpUtil.doPostJson(uri, mapper.writeValueAsString(arcAdd));
        CommonReply reply = mapper.readValue(res, CommonReply.class);
        if (reply.getCode() != 0) {
            throw new ArcBizException(reply);
        }
        return (String) reply.getData().get("resource_id");
    }
}
