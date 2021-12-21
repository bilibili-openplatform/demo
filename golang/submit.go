package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ArcAddReq struct {
	Mid       int64             `json:"mid" validate:"required"`
	Title     string            `json:"title" validate:"required"`
	Cover     string            `json:"cover" validate:"required"`
	TypeID    int16             `json:"tid" validate:"required"`
	Tag       string            `json:"tag" validate:"required"`
	Desc      string            `json:"desc"`
	Copyright int8              `json:"copyright" validate:"required"`
}

type ArcAddResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ResourceId string `json:"resource_id"`
	} `json:"data"`
}

func arcAdd(accessToken string, uploadToken string, p *ArcAddReq) (res *ArcAddResp, err error) {
	params := url.Values{}
	params.Set("client_id", _client_id)
	params.Set("access_token", accessToken)
	params.Set("upload_token", uploadToken)
	uri := _openapi_archive_add + "?" + params.Encode()
	arcAddReqStr, _ := json.Marshal(p)
	fmt.Printf("arcAdd uri(%v) reqbody(%v) ready\n", uri, string(arcAddReqStr))
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(arcAddReqStr))
	if err != nil {
		fmt.Printf("arcAdd http.NewRequest error(%v)\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("arcAdd c.Do error(%v), uri(%s)\n", err, uri)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("arcAdd read resp err(%v)\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("arcAdd uri(%v) statecode(%v), failed\n", uri, resp.StatusCode)
		return
	}
	fmt.Printf("arcAdd uri(%v) respbody(%v) success\n", uri, string(bs))
	res = &ArcAddResp{}
	if err = json.Unmarshal(bs, res); err != nil {
		fmt.Printf("arcAdd json.Unmarshal error(%v)", err)
		return
	}
	return
}
