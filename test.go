package main

import (
	"Teheran/request"
	"fmt"
)

func main() {
	var req request.Request
	resp, err := req.Get("http://112.74.200.115/ip",
		req.WithAllowRedirects(true),
		req.WithTimeout(1),
		req.WithJson(map[string]interface{}{
			"11": "22",}),
		req.WithParams(map[string]string{
			"11": "22",}),
		)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp.StatusCode)
	json, err := resp.Json()
	if err != nil {
		fmt.Println("json err", err.Error())
	}
	fmt.Println(json["ip"])
}

