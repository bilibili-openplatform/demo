package main

import (
	"bytes"
	"fmt"
	"os"
)

const (
	_client_id    = "xxx" // 应用client_id，请设置成真实信息
	_client_secret = "xxx" // 应用secret，请设置成真实信息
	_openapi_auth   = "https://api.bilibili.com/x/account-oauth2/v1/token"
	_openapi_videoup_init   = "http://member.bilibili.com/arcopen/fn/archive/video/init"
	_openapi_videoup_upload = "http://openupos.bilivideo.com/video/v2/part/upload"
	_openapi_videoup_merge  = "http://member.bilibili.com/arcopen/fn/archive/video/complete"
	_openapi_archive_upcover = "http://member.bilibili.com/arcopen/fn/archive/cover/upload"
	_openapi_archive_add     = "http://member.bilibili.com/arcopen/fn/archive/add-by-utoken"
	_openapi_type_list = "http://member.bilibili.com/arcopen/fn/archive/type/list"
)

func main() {
	/* step 0: 授权登录，入参code通过授权SDK获取
	 * 获取到的access_token等结果建议保存起来，以免每次交互都重新请求
	 * 根据expires_in识别access_token的过期时间
	 * 可通过refresh_token来续期access_token
	 */
	code := "xxx" // 授权code，请通过授权SDK获取
	authLoginResp, err := authLogin(code)
	if err != nil || authLoginResp.Code != 0 {
		fmt.Printf("登录失败, code=%s\n", code)
		return
	}
	fmt.Printf("登录成功, resp=(%+v)\n", authLoginResp.Data)
	aToken := authLoginResp.Data.AccessToken

	// 文件信息获取
	filePath := "xxx" // 文件路径，请设置成真实信息
	fs, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("读取文件信息失败, path=%s\n", filePath)
		return
	}

	//step 1.1: 文件上传初始化
	initReq := &VideoupInitReq{
		Name: fs.Name(),
	}
	videoupInitResp := videoupInit(aToken, initReq)
	uToken := videoupInitResp.Data.UploadToken
	if videoupInitResp.Code != 0 {
		fmt.Printf("文件预处理失败, resp=%+v\n", videoupInitResp)
		return
	}

	//step 1.2: 文件分片信息计算，8M一个分片
	chunkInfos := preChunk(fs, 8*1024*1024)
	fmt.Printf("分片信息 %+v\n", chunkInfos)

	//step 1.3: 读取本地文件并分片上传
	f, _ := os.Open(filePath)
	defer f.Close()
	for _, v := range chunkInfos {
		bs := make([]byte, v.Size)
		f.ReadAt(bs, v.Start)
		videoUploadResp, err := videoUpload(uToken, v.PartNum, v, bytes.NewReader(bs))
		if err != nil || videoUploadResp.Code != 0 {
			fmt.Printf("文件分片上传失败, part_num=%d, resp=%+v, err=%+v\n", v.PartNum, videoupInitResp, err)
			return
		}
	}

	//step 1.4 分片合并
	videoupMerge(uToken)

	//step 2.1 稿件封面上传
	coverPath := "xxx" // 封面路径，请设置成真实信息
	coverUploadResp, err := uploadCover(aToken, coverPath)
	if err != nil || coverUploadResp.Code != 0 {
		fmt.Printf("封面上传失败, resp=%+v, err=%+v\n", coverUploadResp, err)
		return
	}
	coverUrl := coverUploadResp.Data.Url

	//// 稿件分区列表
	//typeListResp, err := typeList(aToken)
	//if err != nil || typeListResp.Code != 0 {
	//	fmt.Printf("获取分区列表失败, resp=%+v, err=%+v\n", typeListResp, err)
	//	return
	//}
	//fmt.Printf("获取分区列表成功, resp=(%+v)\n", authLoginResp.Data)

	//step 3.1 稿件提交
	arcAddReq := &ArcAddReq{
		Title:     "xxx", // 稿件标题，请设置成真实信息
		Cover:     coverUrl,
		TypeID:    21, // 稿件分区，请使用分区查询接口找到真实的分区信息，可参考typeList()方法
		Tag:       "xxx", // 稿件标签，请设置成真实信息
		Desc:      "xxx", // 稿件描述，请设置成真实信息
		Copyright: 1,
	}
	arcAddResp, err := arcAdd(aToken, uToken, arcAddReq)
	if err != nil || arcAddResp.Code != 0 {
		fmt.Printf("稿件提交失败, resp=%+v, err=%+v\n", arcAddResp, err)
		return
	}
	resourceId := arcAddResp.Data.ResourceId
	fmt.Printf("稿件提交成功 resource_id %s\n", resourceId)
}
