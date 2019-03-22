package request

import (
	"bytes"
	"crypto/tls"
	goJson "encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)


type args struct {
	allowRedirects bool
	verify bool
	timeout uint
	params map[string]string
	json map[string]interface{}
}

type fakeArgs func(*args)

type Request struct {}

func (cls *Request) Get(url string, arg ...fakeArgs) (*Resp, error) {
    return cls.newRequest(url, "GET", arg...)
}

func (cls *Request) Post(url string, arg ...fakeArgs) (*Resp, error) {
    return cls.newRequest(url, "POST", arg...)
}

func (cls *Request) Put(url string, arg ...fakeArgs) (*Resp, error) {
	return cls.newRequest(url, "PUT", arg...)
}

func (cls *Request) Delete(url string, arg ...fakeArgs) (*Resp, error) {
	return cls.newRequest(url, "DELETE", arg...)
}

func (cls *Request) Patch(url string, arg ...fakeArgs) (*Resp, error) {
	return cls.newRequest(url, "PATCH", arg...)
}

func (cls *Request) Options(url string, arg ...fakeArgs) (*Resp, error) {
	return cls.newRequest(url, "OPTIONS", arg...)
}

func (cls *Request) WithParams(params map[string]string) fakeArgs {
    return func(arg *args) {
    	arg.params = params
	}
}

func (cls *Request) WithJson(json map[string]interface{}) fakeArgs {
	return func(arg *args) {
		arg.json = json
	}
}

func (cls *Request) WithTimeout(timeout uint) fakeArgs {
	return func(arg *args) {
		arg.timeout = timeout
	}
}

func (cls *Request) WithAllowRedirects(allowRedirects bool) fakeArgs {
	return func(arg *args) {
		arg.allowRedirects = allowRedirects
	}
}

func (cls *Request) WithVerify(verify bool) fakeArgs {
	return func(arg *args) {
		arg.verify = verify
	}
}



func (cls *Request) newRequest(url, method string, arg ...fakeArgs) (*Resp, error) {
    var (
    	err error
    	body io.Reader
    	jsonVal []byte
    	req *http.Request
    	resp *http.Response
	)

    ar := new(args)
    for _, a := range arg {
    	a(ar)
	}

    if ar.json != nil && method != "GET" {
    	jsonVal, err = goJson.Marshal(ar.json)
    	if err != nil {
    		return nil, err
		}

    	body = bytes.NewBuffer(jsonVal)
	}

    req, err = http.NewRequest(method, url, body)
    if err != nil {
    	return nil, err
	}

    if ar.json != nil && method != "GET" {
    	req.Header.Set("Content-Type", "application/json")
	}

    if ar.params != nil {
    	q := req.URL.Query()
    	for k, v := range ar.params {
            q.Add(k, v)
		}
    	req.URL.RawQuery = q.Encode()
	}

    cli := &http.Client{
    	Timeout: time.Second * time.Duration(ar.timeout),
	}

    if !ar.allowRedirects {
    	cli.CheckRedirect = cls.disableRedirect
	}

    if !ar.verify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		cookieJar, _ := cookiejar.New(nil)

		cli.Jar = cookieJar
		cli.Transport = tr
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

func (cls *Request) disableRedirect(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
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
