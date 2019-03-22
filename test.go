package main

import (
	"Teheran/requests"
	"fmt"
)

func main() {
	resp, err := requests.Get("http://112.74.200.115/ip",
		requests.WithJson(map[string]interface{}{
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

