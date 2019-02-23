package main

import (
	"Teheran/request"
	"fmt"
)

func main() {
	var req request.Request
	resp, err := req.Post("http://112.74.200.115", nil, nil, 0)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(resp.StatusCode)
	json, err := resp.Json()
	if err != nil {
		fmt.Println("json err", err.Error())
	}
	fmt.Println(json["errmsg"])
}

