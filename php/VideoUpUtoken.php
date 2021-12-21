<?php
/**
 * PHP Version 7.3.24
 */
const CLIENT_ID = "xxx"; // 应用client_id，请设置成真实信息
const CLIENT_SECRET = "xxx"; // 应用secret，请设置成真实信息
const URL_AUTH = "https://api.bilibili.com/x/account-oauth2/v1/token";
const URL_VIDEO_INIT = "http://member.bilibili.com/arcopen/fn/archive/video/init";
const URL_VIDEO_UP = "http://openupos.bilivideo.com/video/v2/part/upload";
const URL_VIDEO_COMPLETE = "http://member.bilibili.com/arcopen/fn/archive/video/complete";
const URL_COVER_UP = "http://member.bilibili.com/arcopen/fn/archive/cover/upload";
const URL_ARCHIVE_ADD = "http://member.bilibili.com/arcopen/fn/archive/add-by-utoken";
const URL_TYPE_LIST = "http://member.bilibili.com/arcopen/fn/archive/type/list";


// 代码示例
{
    /* 授权登录，入参code通过授权SDK获取
     * 获取到的access_token等结果建议保存起来，以免每次交互都重新请求
     * 根据expires_in识别access_token的过期时间
     * 可通过refresh_token来续期access_token
     */
    $code = "xxx"; // 授权code，请通过授权SDK获取
    $authData = authLogin($code);
    if (!$authData) {
        return;
    }
    $aToken = $authData["access_token"];

    // 文件基本信息
    $filePath = "xxx"; // 视频文件路径，请设置成真实信息
    $name = basename($filePath);

    // 视频上传预处理
    $initData = videoUpInit($aToken, $name);
    if (!$initData) {
        return;
    }
    $uToken = $initData["upload_token"];

    // 文件切片
    try {
        $partNum = 0;
        $file_pointer = fopen($filePath, 'rb');
        while (!feof($file_pointer)) {
            $partNum++; // 分片号从1开始
            $tmp = fread($file_pointer, 8364032); // 每次8M
            echo "文件分片".$partNum.", 大小".strlen($tmp)."\n";
            // 分片上传
            $upResp = videoUpPart($uToken, $partNum, $tmp);
            if (!$upResp) {
                return;
            }
        }
        fclose($file_pointer);
    } catch (Exception $e) {
        echo "文件切片失败: code=".$e->getCode().", message=".$e->getMessage()."\n";
    }

    // 合片
    $completeResp = videoUpComplete($uToken);
    if (!$completeResp) {
        return;
    }

    // 图片上传
    $coverPath = "xxx"; // 封面图路径，请设置成真实信息
    $coverResp = coverUp($aToken, $coverPath);
    if (!$coverResp) {
        return;
    }
    $coverUrl = $coverResp["url"];

    // 稿件提交
    $title = "xxx"; // 稿件标题，请设置成真实信息
    $tid = 21; // 稿件分区，请使用分区查询接口找到真实的分区信息，可参考typeList()方法
    $tag = "xxx"; // 稿件标签，请设置成真实信息
    $addResp = archiveAdd($aToken, $title, $coverUrl, $tid, $tag, $uToken);
    if (!$addResp) {
        return;
    }
    echo "稿件提交成功: resource_id=".$addResp["resource_id"];
}


/**
 * 授权登录，获取access_token等信息
 * @param $code
 * @return array|bool
 */
function authLogin($code)
{
    // 获取access_token
    $url = URL_AUTH."?client_id=".CLIENT_ID."&client_secret=".CLIENT_SECRET. "&grant_type=authorization_code&code=".$code;

    $resp =  curl_post($url);
    if(!isset($resp['code']) || $resp['code'] != 0) {
        echo "授权登陆失败:".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
        return false;
    }
    echo "authLogin: ".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
    return $resp['data'];
}


/**
 * 视频上传预处理
 * @param $aToken
 * @param $name
 * @return array|bool
 */
function videoUpInit($aToken, $name)
{
    $url = URL_VIDEO_INIT."?client_id=".CLIENT_ID."&access_token=".$aToken;
    $body = [
        'name'=>$name
    ];
    $req_headers    = array();
    $req_headers[]  = 'Content-Type: application/json;';
    $resp = curl_post($url, json_encode($body,JSON_UNESCAPED_UNICODE), $req_headers);
    if(!isset($resp['code']) || $resp['code'] != 0) {
        echo "文件预处理失败:".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
        return false;
    }
    echo "videoUpInit: ".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
    return $resp['data'];
}


/**
 * 视频文件分片
 * @param $uToken
 * @param $partNum
 * @param $partData
 * @return bool
 */
function videoUpPart($uToken, $partNum, $partData): bool
{
    $url = URL_VIDEO_UP."?upload_token=".$uToken."&part_number=".$partNum;
    $resp = curl_post($url, $partData);
    if(!isset($resp['code']) || $resp['code'] != 0) {
        echo "上传分片失败".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
        return false;
    }
    echo "videoUpPart: ".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
    return true;
}


/**
 * 视频分片合片
 * @param $uToken
 * @return bool
 */
function videoUpComplete($uToken): bool
{
    $url = URL_VIDEO_COMPLETE."?upload_token=".$uToken;
    $resp = curl_post($url);
    if(!isset($resp['code']) || $resp['code'] != 0) {
        echo "合并分片失败".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
        return false;
    }
    echo "videoUpComplete: ".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
    return true;
}


/**
 * 封面上传
 * @param $aToken
 * @param $coverFile
 * @return array|bool
 */
function coverUp($aToken, $coverFile)
{
    $url = URL_COVER_UP."?client_id=".CLIENT_ID."&access_token=".$aToken;
    $data = [
        'file'=>new \CURLFile($coverFile)
    ];
    $req_headers    = array();
    $req_headers[]  = 'Content-Type:multipart/form-data';
    $resp = curl_post($url, $data, $req_headers);
    if(!isset($resp['code']) || $resp['code'] != 0) {
        echo "上传封面失败".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
        return false;
    }
    echo "coverUp: ".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
    return $resp['data'];
}


/**
 * 稿件提交
 * @param $aToken
 * @param $title
 * @param $coverUrl
 * @param $tid
 * @param $tag
 * @param $uToken
 * @return array|bool
 */
function archiveAdd($aToken, $title, $coverUrl, $tid, $tag, $uToken)
{
    $url = URL_ARCHIVE_ADD."?client_id=".CLIENT_ID."&access_token=".$aToken."&upload_token=".$uToken;
    $body = [
        'title'=>$title,
        'cover'=>$coverUrl,
        'tid'=>$tid,
        'tag'=>$tag,
        'copyright'=>1
    ];
    $req_headers    = array();
    $req_headers[]  = 'Content-Type: application/json;';
    $resp = curl_post($url, json_encode($body,JSON_UNESCAPED_UNICODE), $req_headers);
    if(!isset($resp['code']) || $resp['code'] != 0) {
        echo "稿件提交失败".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
        return false;
    }
    echo "archiveAdd: ".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
    return $resp['data'];
}


/**
 * 分区查询
 * @param $aToken
 * @return array|bool
 */
function typeList($aToken)
{
    $url = URL_TYPE_LIST."?client_id=".CLIENT_ID."&access_token=".$aToken;
    $resp = curl_get($url);
    if(!isset($resp['code']) || $resp['code'] != 0) {
        echo "分区查询失败".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
        return false;
    }
    echo "typeList: ".json_encode($resp,JSON_UNESCAPED_UNICODE)."\n";
    return $resp['data'];
}


// POST请求
function curl_post($url, $data = array(), $header = array())
{
    // 初始化
    $curl = curl_init();
    if (!empty($header)) {
        curl_setopt($curl, CURLOPT_HTTPHEADER, $header);
    }
    // 设置抓取的url
    curl_setopt($curl, CURLOPT_URL, $url);
    // 设置头文件的信息作为数据流输出
    curl_setopt($curl, CURLOPT_HEADER, 0);
    // 设置获取的信息以文件流的形式返回，而不是直接输出。
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
    // 设置post方式提交
    curl_setopt($curl, CURLOPT_POST, 1);
    // 设置post数据
    curl_setopt($curl, CURLOPT_POSTFIELDS, $data);

    curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, false); // 不验证证书下同
    curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, false);
    // 执行命令
    $json = curl_exec($curl);

    // 关闭URL请求
    curl_close($curl);

    $result = json_decode($json, true);

    return $result;
}


// GET请求
function curl_get($url, $data = array())
{
    // 初始化
    $ch = curl_init();
    // 设置选项，包括URL
    if(!empty($data)){
        $query = http_build_query($data);
        $url = $url . '?' . $query;
    }
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($ch, CURLOPT_HEADER, 0);
    // 执行并获取HTML文档内容
    $output = curl_exec($ch);
    // 释放curl句柄
    curl_close($ch);

    $result = json_decode($output, true);

    return $result;
}
