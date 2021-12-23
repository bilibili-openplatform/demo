package com.example.demo;

import com.example.demo.entity.ArcAddParam;
import com.example.demo.util.ArchiveUtil;

import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.util.Arrays;

public class Demo {

    private final static String clientId = "这里填入您的 client_id";
    private final static String secret = "这里填入对应的 client_secret";

    // 视频文件名
    private final static String fileName = "video_test.MP4";
    // 视频文件路径
    private final static String filePath = "/Users/bilibili/files/video_test.MP4";
    // 稿件封面图片路径
    private final static String coverPath = "/Users/bilibili/files/cover_test.png";

    public static void main(String[] args) throws Exception {
        ArchiveUtil arcUtil = new ArchiveUtil(clientId, secret);

        // 通过临时 code 获取 access_token，详见账号授权文档：
        // https://openhome.bilibili.com/doc/4/eaf0e2b5-bde9-b9a0-9be1-019bb455701c
        String accessToken = arcUtil.getAccessToken("这里填入您获取的code");
        System.out.println("access_token: " + accessToken);

        // 文件上传预处理
        // https://openhome.bilibili.com/doc/4/0c532c6a-e6fb-0aff-8021-905ae2409095
        String uploadToken = arcUtil.arcInit(accessToken, fileName);
        System.out.println("arcInit success: " + uploadToken);

        // 文件分片上传
        // https://openhome.bilibili.com/doc/4/733a520a-c50f-7bb4-17cb-35338ba20500
        File file = new File(filePath);
        InputStream is = new FileInputStream(file);
        int partSize = 8 * 1024 * 1024;
        byte[] part = new byte[partSize];
        int len;
        int partNum = 1;
        // 可适当使用并发上传，线程数不宜过大
        while ((len = is.read(part)) != -1) {
            if (len < partSize) {
                part = Arrays.copyOf(part, len);
            }
            arcUtil.arcUp(uploadToken, partNum, part);
            System.out.println("arcUpload success: " + partNum);
            partNum++;
        }

        // 文件分片合片
        // https://openhome.bilibili.com/doc/4/0828e499-38d8-9e58-2a70-a7eaebf9dd64
        arcUtil.arcComplete(uploadToken);
        System.out.println("arcComplete success");

        // 上传封面
        // https://openhome.bilibili.com/doc/4/8243399e-50e3-4058-7f01-1ebe4c632cf8
        String coverUrl = arcUtil.uploadCover(accessToken, new File(coverPath));
        System.out.println("uploadCover success: " + coverUrl);

        // 构造稿件提交参数，提交稿件
        // https://openhome.bilibili.com/doc/4/f7fc57dd-55a1-5cb1-cba4-61fb2994bf0f
        ArcAddParam arcAddParam = new ArcAddParam();
        arcAddParam.setTitle("测试投稿-" + System.currentTimeMillis());
        arcAddParam.setCover(coverUrl);
        // 调用分区查询接口获取，选择合适的分区
        // https://openhome.bilibili.com/doc/4/4f13299b-5316-142f-df6a-87313eaf85a9
        arcAddParam.setTid(75);
        arcAddParam.setNoReprint(1);
        arcAddParam.setDesc("测试投稿-描述");
        arcAddParam.setTag("生活,搞笑,游戏");
        arcAddParam.setCopyright(1);
        arcAddParam.setSource("");

        String resourceId = arcUtil.arcSubmit(accessToken, uploadToken, arcAddParam);
        System.out.println("arcSubmit success: " + resourceId);
    }
}
