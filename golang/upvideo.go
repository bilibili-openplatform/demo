package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type VideoupInitReq struct {
	Name    string    `json:"name"`
}

type VideoInitResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UploadToken      string  `json:"upload_token"`
	}
}

type VideoUploadResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type VideoMergeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UploadToken      string  `json:"upload_token"`
	}
}

func videoupInit(accessToken string, p *VideoupInitReq) (res *VideoInitResp) {
	params := url.Values{}
	params.Set("client_id", _client_id)
	params.Set("access_token", accessToken)
	uri := _openapi_videoup_init + "?" + params.Encode()
	initReqStr, _ := json.Marshal(p)
	fmt.Printf("videoupInit uri(%v) reqbody(%v) ready\n", uri, string(initReqStr))
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(initReqStr))
	if err != nil {
		fmt.Printf("videoupInit http.NewRequest error(%v)\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("videoupInit c.Do error(%v), uri(%s)\n", err, uri)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("videoupInit read resp err(%v)\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("videoupInit uri(%v) statecode(%v), failed\n", uri, resp.StatusCode)
		return
	}
	fmt.Printf("videoupInit uri(%v) respbody(%v) success\n", uri, string(bs))
	_ = json.Unmarshal(bs, &res)
	return
}

func videoUpload(uploadToken string, partNum int64, chunkInfo ChunkInfo, body io.Reader) (res *VideoUploadResp, err error) {
	params := url.Values{}
	params.Set("upload_token", uploadToken)
	params.Set("part_number", strconv.FormatInt(partNum, 10))
	uri := _openapi_videoup_upload + "?" + params.Encode()
	fmt.Printf("videoupLoad(%v) put uri(%v) ready\n", chunkInfo, uri)
	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		fmt.Printf("videoupLoad http.NewRequest error(%v)\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("videoUpload c.Do error(%v), uri(%s)\n", err, uri)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("videoUpload read resp err(%v)\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("videoupLoad(%v) put uri(%v) statecode(%v), failed\n", chunkInfo, uri, resp.StatusCode)
		return
	}
	fmt.Printf("videoupLoad(%v) put uri(%v) respbody(%v) success\n", chunkInfo, uri, string(bs))
	_ = json.Unmarshal(bs, &res)
	return
}

func videoupMerge(uploadToken string) (res *VideoMergeResp) {
	params := url.Values{}
	params.Set("upload_token", uploadToken)
	uri := _openapi_videoup_merge + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		fmt.Printf("videoupMerge http.NewRequest error(%v)\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("videoupMerge c.Do error(%v), uri(%s)\n", err, uri)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("videoupMerge read resp err(%v)", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("videoupMerge uri(%v) statecode(%v), failed\n", uri, resp.StatusCode)
		return
	}
	fmt.Printf("videoupMerge uri(%v) respbody(%v) success\n", uri, string(bs))
	_ = json.Unmarshal(bs, &res)
	return
}

func preChunk(fs os.FileInfo, chunksize int64) (res []ChunkInfo) {
	n := fs.Size() / chunksize
	var i int64
	for i = 0; i < n; i++ {
		v := ChunkInfo{
			PartNum: i + 1,
			Chunk:   i,
			Chunks:  n + 1,
			Size:    chunksize,
			Start:   i * chunksize,
			End:     (i + 1) * chunksize,
			Total:   fs.Size(),
		}
		res = append(res, v)
	}
	lastsize := fs.Size() - n*chunksize
	if lastsize == 0 {
		return
	}
	v := ChunkInfo{
		PartNum: i + 1,
		Chunk:   i,
		Chunks:  n + 1,
		Size:    lastsize,
		Start:   i * chunksize,
		End:     i*chunksize + lastsize,
		Total:   fs.Size(),
	}
	res = append(res, v)
	return
}

type ChunkInfo struct {
	PartNum int64
	Chunk   int64
	Chunks  int64
	Size    int64
	Start   int64
	End     int64
	Total   int64
}
