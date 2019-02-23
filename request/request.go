package request

import (
	"bytes"
	goJson "encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {}

func (cls *Request) Get(url string, params map[string]string, timeout uint) (*Resp, error) {
    return cls.newRequest(url, "GET", params, nil, timeout)
}

func (cls *Request) Post(url string, params map[string]string, json map[string]interface{}, timeout uint) (*Resp, error) {
    return cls.newRequest(url, "POST", params, json, timeout)
}

func (cls *Request) Put(url string, params map[string]string, json map[string]interface{}, timeout uint) (*Resp, error) {
	return cls.newRequest(url, "PUT", params, json, timeout)
}

func (cls *Request) Delete(url string, params map[string]string, json map[string]interface{}, timeout uint) (*Resp, error) {
	return cls.newRequest(url, "DELETE", params, json, timeout)
}

func (cls *Request) Patch(url string, params map[string]string, json map[string]interface{}, timeout uint) (*Resp, error) {
	return cls.newRequest(url, "PATCH", params, json, timeout)
}

func (cls *Request) Options(url string, params map[string]string, json map[string]interface{}, timeout uint) (*Resp, error) {
	return cls.newRequest(url, "OPTIONS", params, json, timeout)
}

func (cls *Request) newRequest(url, method string, params map[string]string, json map[string]interface{}, timeout uint) (*Resp, error) {
    var (
    	err error
    	body io.Reader
    	jsonVal []byte
    	req *http.Request
    	resp *http.Response
	)

    if json != nil {
    	jsonVal, err = goJson.Marshal(json)
    	if err != nil {
    		return nil, err
		}

    	body = bytes.NewBuffer(jsonVal)
	}

    req, err = http.NewRequest(method, url, body)
    if err != nil {
    	return nil, err
	}

    if json != nil {
    	req.Header.Set("Content-Type", "application/json")
	}

    if params != nil {
    	q := req.URL.Query()
    	for k, v := range params {
            q.Add(k, v)
		}
    	req.URL.RawQuery = q.Encode()
	}

    cli := &http.Client{
    	Timeout: time.Second * time.Duration(timeout),
	}

    resp, err = cli.Do(req)
    if err != nil {
    	return nil, err
	}
    defer resp.Body.Close()

    resP := &Resp{
    	StatusCode: resp.StatusCode,
    	Header: resp.Header,
	}

    resP.Body, err = ioutil.ReadAll(resp.Body)
    return resP, err
}

type Resp struct {
	StatusCode int
	Body []byte
	Header map[string][]string
}

func (cls *Resp) Json() (val map[string]interface{}, err error) {
	err = goJson.Unmarshal(cls.Body, &val)
	return
}
