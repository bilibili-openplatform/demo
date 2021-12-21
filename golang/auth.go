package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AuthLoginResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccessToken string  `json:"access_token"`
		RefreshToken string  `json:"refresh_token"`
		ExpiresIn int64  `json:"expires_in"`
	}
}

func authLogin(code string) (res *AuthLoginResp, err error) {
	params := url.Values{}
	params.Set("client_id", _client_id)
	params.Set("client_secret", _client_secret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	uri := _openapi_auth + "?" + params.Encode()
	fmt.Printf("authLogin uri(%v)\n", uri)
	req, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		fmt.Printf("authLogin http.NewRequest error(%v)\n", err)
		return
	}
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("authLogin c.Do error(%v), uri(%s)\n", err, uri)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("authLogin read resp err(%v)\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("authLogin uri(%v) statecode(%v), failed\n", uri, resp.StatusCode)
		return
	}
	fmt.Printf("authLogin uri(%v) respbody(%v) success\n", uri, string(bs))
	json.Unmarshal(bs, &res)
	return
}
