package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type UPloadCoverResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Url string `json:"url"`
	} `json:"data"`
}

func uploadCover(accessToken string, cover string) (res *UPloadCoverResp, err error) {
	params := url.Values{}
	params.Set("client_id", _client_id)
	params.Set("access_token", accessToken)
	uri := _openapi_archive_upcover + "?" + params.Encode()
	fmt.Printf("uploadCover uri(%v) ready\n", uri)
	return generateFileAndUpload(cover, uri)
}

func generateFileAndUpload(filePath string, uri string) (res *UPloadCoverResp, err error) {
	res = &UPloadCoverResp{}
	fs, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("uploadCover stat file(%v) error(%v)\n", filePath, err)
		return
	}
	file1, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("uploadCover open file(%v) error(%v)\n", filePath, err)
		return
	}
	defer file1.Close()
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter1, err := bodyWriter.CreateFormFile("file", fs.Name())
	if err != nil {
		fmt.Printf("uploadCover bodyWriter.CreateFormFile error(%v)\n", err)
		return
	}
	_, err = io.Copy(fileWriter1, file1)
	if err != nil {
		fmt.Printf("uploadCover io.Copy error(%v)\n", err)
		return
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	req, err := http.NewRequest(http.MethodPost, uri, bodyBuffer)
	if err != nil {
		fmt.Printf("uploadCover http.NewRequest error(%v)\n", err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("uploadCover c.Do error(%v)\n", err)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("uploadCover ioutil.ReadAll(resp.Body) error(%v)\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("uploadCover uri(%v) statecode(%v), failed\n", uri, resp.StatusCode)
		return
	}
	fmt.Printf("uploadCover uri(%v) respbody(%v) success\n", uri, string(bs))
	if err = json.Unmarshal(bs, res); err != nil {
		fmt.Printf("uploadCover json.Unmarshal error(%v)", err)
		return
	}
	return
}
