# Teheran

A lightweight, high-performance HTTP client library for Go. Teheran provides a clean, functional-options API for making REST API calls with minimal boilerplate. It is optimized for high-throughput scenarios using object pooling to reduce GC pressure.

![image](teheran.jpg)

## Features

- Simple, chainable API with functional options pattern
- Support for all standard HTTP methods (GET, POST, PUT, DELETE, PATCH, OPTIONS)
- JSON request body serialization
- Query parameter and custom header support
- Configurable request timeout
- Redirect control
- TLS verification toggle
- Object pooling via `sync.Pool` for high-concurrency performance

## Installation

```bash
go get github.com/FastSchnell/Teheran/requests
```

Requires Go 1.21 or later.

## Usage

### Basic GET Request

```go
package main

import (
    "fmt"
    "github.com/FastSchnell/Teheran/requests"
)

func main() {
    resp, err := requests.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(resp.StatusCode)
    json, err := resp.Json()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(json["url"])
}
```

### POST with JSON Body

```go
body := map[string]interface{}{
    "username": "admin",
    "email":    "admin@example.com",
}
resp, err := requests.Post("https://httpbin.org/post", requests.WithJson(body))
```

### Query Parameters

```go
resp, err := requests.Get("https://httpbin.org/get",
    requests.WithParams(map[string]string{
        "page":  "1",
        "limit": "20",
    }),
)
```

### Custom Headers

```go
resp, err := requests.Get("https://api.example.com/data",
    requests.WithHeaders(map[string]string{
        "Authorization": "Bearer <token>",
        "Accept":        "application/json",
    }),
)
```

### Request Timeout

```go
resp, err := requests.Get("https://api.example.com/data",
    requests.WithTimeout(5 * time.Second),
)
```

### Disable Redirects

```go
resp, err := requests.Get("https://example.com/redirect",
    requests.WithAllowRedirects(false),
)
// resp.StatusCode will be 301/302 instead of following the redirect
```

### Skip TLS Verification

```go
resp, err := requests.Get("https://self-signed.example.com",
    requests.WithVerify(false),
)
```

### Response Handling

```go
// Parse response as map[string]interface{}
json, err := resp.Json()

// Parse response into a struct
var user User
_, err := resp.Json(&user)

// Parse response as a JSON array
list, err := resp.List()

// Parse response as map[string]string
strMap, err := resp.JsonAndValueIsString()

// Access response headers
headers := resp.Header()

// Access raw response body
raw := resp.Body
```

## License

Teheran is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

## Story

This project was created to commemorate one of the most difficult resignation experiences in my career, which took place from January 31, 2019 to February 28, 2019. I personally refer to it as the "Escape from Tehran" plan.
