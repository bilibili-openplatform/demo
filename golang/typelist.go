package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TypeListResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccessToken string  `json:"access_token"`
		RefreshToken string  `json:"refresh_token"`
		ExpiresIn int64  `json:"expires_in"`
	}
}

func typeList(aToken string) (res *AuthLoginResp, err error) {
	params := url.Values{}
	params.Set("client_id", _client_id)
	params.Set("access_token", aToken)
	uri := _openapi_type_list + "?" + params.Encode()
	fmt.Printf("typeList uri(%v)\n", uri)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Printf("typeList http.NewRequest error(%v)\n", err)
		return
	}
	c := http.DefaultClient
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("typeList c.Do error(%v), uri(%s)\n", err, uri)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("typeList read resp err(%v)\n", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("typeList uri(%v) statecode(%v), failed\n", uri, resp.StatusCode)
		return
	}
	fmt.Printf("typeList uri(%v) respbody(%v) success\n", uri, string(bs))
	json.Unmarshal(bs, &res)
	return
}
