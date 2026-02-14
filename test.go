package main

import (
	"github.com/FastSchnell/Teheran/requests"
	"fmt"
)

func main() {
	resp, err := requests.Get("https://httpbin.org/get",
		requests.WithParams(map[string]string{
			"key": "value",
		}),
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
	fmt.Println(json["url"])
}

