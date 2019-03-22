Teheran
=======

Teheran是一个简单的HTTP库，能让请求rest api更加简洁。body入参格式只支持json，其他场景还是建议使用[Go](http://golang.org/)标准库http。

![image](teheran.jpg)


Installation
------------

Install Teheran using the "go get" command:

    go get github.com/FastSchnell/Teheran/requests
    
    
Usage
-----
```go
import "Teheran/requests"

func testReq() {
    body := map[string]interface{}{"11": "22"}
    resp, err := requests.Post("http://112.74.200.115", requests.WithJson(body))
    if err != nil {
        fmt.Println(err.Error())
    }

    if resp.StatusCode == 404 {
        json, err := resp.Json()
        if err != nil {
            fmt.Println("json err", err.Error())
        }

        for k, v := range json {
            fmt.println(k, v)
        }
    }
}

```


License
-------

Teheran is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
    


项目故事
-------
为了纪念最痛苦的一次离职经历，2019.01.31-2019.02.28。个人称之为"逃离德黑兰"计划。