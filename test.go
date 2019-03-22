package main

import (
	"Teheran/request"
	"fmt"
)

func main() {
	var req request.Request
	req.AllowRedirects = false
	req.Verify = false
	resp, err := req.Post("https://cfg.aiclk.com/hdjump?iclicashid=7414647", nil, nil, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp.StatusCode)
	json, err := resp.Json()
	if err != nil {
		fmt.Println("json err", err.Error())
	}
	fmt.Println(json["errmsg"])
}

