package requests

import (
	"bytes"
	"crypto/tls"
	goJson "encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)


type args struct {
	allowRedirects bool
	verify bool
	timeout uint
	params map[string]string
	json map[string]interface{}
	headers map[string]string
}

type fakeArgs func(*args)


func Get(url string, arg ...fakeArgs) (*Resp, error) {
    return newRequest(url, "GET", arg...)
}

func Post(url string, arg ...fakeArgs) (*Resp, error) {
    return newRequest(url, "POST", arg...)
}

func Put(url string, arg ...fakeArgs) (*Resp, error) {
	return newRequest(url, "PUT", arg...)
}

func Delete(url string, arg ...fakeArgs) (*Resp, error) {
	return newRequest(url, "DELETE", arg...)
}

func Patch(url string, arg ...fakeArgs) (*Resp, error) {
	return newRequest(url, "PATCH", arg...)
}

func Options(url string, arg ...fakeArgs) (*Resp, error) {
	return newRequest(url, "OPTIONS", arg...)
}

func WithParams(params map[string]string) fakeArgs {
    return func(arg *args) {
    	arg.params = params
	}
}

func WithJson(json map[string]interface{}) fakeArgs {
	return func(arg *args) {
		arg.json = json
	}
}

func WithHeaders(headers map[string]string) fakeArgs {
	return func(arg *args) {
		arg.headers = headers
	}
}

func WithTimeout(timeout uint) fakeArgs {
	return func(arg *args) {
		arg.timeout = timeout
	}
}

func WithAllowRedirects(allowRedirects bool) fakeArgs {
	return func(arg *args) {
		arg.allowRedirects = allowRedirects
	}
}

func WithVerify(verify bool) fakeArgs {
	return func(arg *args) {
		arg.verify = verify
	}
}



func newRequest(url, method string, arg ...fakeArgs) (*Resp, error) {
    var (
    	err error
    	body io.Reader
    	jsonVal []byte
    	req *http.Request
    	resp *http.Response
	)

    ar := &args{
    	allowRedirects: true,
    	verify: true,
	}
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

    if ar.headers != nil {
    	for k, v := range ar.headers {
    		req.Header.Set(k, v)
		}
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
    	cli.CheckRedirect = disableRedirect
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
    	Header: getHeader(resp.Header),
	}

    resP.Body, err = ioutil.ReadAll(resp.Body)
    return resP, err
}

func getHeader(header map[string][]string) map[string]string {
	newHeader := make(map[string]string)
	for k, v := range header {
		newHeader[k] = strings.Join(v, ",")
	}

	return newHeader
}

func disableRedirect(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

type Resp struct {
	StatusCode int
	Body []byte
	Header map[string]string
}

func (cls *Resp) Json() (val map[string]interface{}, err error) {
	err = goJson.Unmarshal(cls.Body, &val)
	return
}
