package requests

import (
	"bytes"
	"crypto/tls"
	goJson "encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

var (
	argsPool = sync.Pool{
		New: func() interface{} {
			return new(args)
		},
	}

	httpClientPool = sync.Pool{
		New: func() interface{} {
			return new(http.Client)
		},
	}

	respPool = sync.Pool{
		New: func() interface{} {
			return new(Resp)
		},
	}
)

type args struct {
	allowRedirects bool
	verify         bool
	timeout        time.Duration
	params         map[string]string
	json           map[string]interface{}
	headers        map[string]string
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

func WithTimeout(timeout time.Duration) fakeArgs {
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
		err     error
		body    io.Reader
		jsonVal []byte
		req     *http.Request
		resp    *http.Response
	)

	ar := argsPool.Get().(*args)
	defer argsPool.Put(ar)
	ar.allowRedirects = true
	ar.verify = true
	ar.timeout = 0
	ar.params = nil
	ar.json = nil
	ar.headers = nil

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

	cli := httpClientPool.Get().(*http.Client)
	defer httpClientPool.Put(cli)
	cli.Timeout = ar.timeout

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

	resP := respPool.Get().(*Resp)
	defer respPool.Put(resP)
	resP.StatusCode = resp.StatusCode
	resP.header = resp.Header
	resP.Body, err = ioutil.ReadAll(resp.Body)

	return resP, err
}

func disableRedirect(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

type Resp struct {
	StatusCode int
	Body       []byte
	header     map[string][]string
}

func (cls *Resp) Header() map[string]string {
	newHeader := make(map[string]string)
	for k, v := range cls.header {
		newHeader[k] = strings.Join(v, ",")
	}

	return newHeader
}

func (cls *Resp) Json(arg ...interface{}) (val map[string]interface{}, err error) {
	if len(arg) == 0 {
		err = goJson.Unmarshal(cls.Body, &val)
	} else if len(arg) == 1 {
		err = goJson.Unmarshal(cls.Body, arg[0])
	} else {
		return nil, errors.New("too many arg, max len is 1")
	}

	return
}

func (cls *Resp) List(arg ...interface{}) (val []interface{}, err error) {
	if len(arg) == 0 {
		err = goJson.Unmarshal(cls.Body, &val)
	} else if len(arg) == 1 {
		err = goJson.Unmarshal(cls.Body, arg[0])
	} else {
		return nil, errors.New("too many arg, max len is 1")
	}

	return
}

func (cls *Resp) JsonAndValueIsString() (val map[string]string, err error) {
	err = goJson.Unmarshal(cls.Body, &val)
	return
}
