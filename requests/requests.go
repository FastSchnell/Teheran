package requests

import (
	"bytes"
	"crypto/tls"
	goJson "encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	netUrl "net/url"
	"strings"
	"sync"
	"time"
)

const defaultTimeout = 30 * time.Second

var (
	argsPool = sync.Pool{
		New: func() interface{} {
			return new(args)
		},
	}

	insecureTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
)

type args struct {
	allowRedirects bool
	verify         bool
	timeout        time.Duration
	params         map[string]string
	json           map[string]interface{}
	data           map[string]string
	headers        map[string]string
}

type Option func(*args)

func Get(url string, arg ...Option) (*Resp, error) {
	return newRequest(url, "GET", arg...)
}

func Post(url string, arg ...Option) (*Resp, error) {
	return newRequest(url, "POST", arg...)
}

func Put(url string, arg ...Option) (*Resp, error) {
	return newRequest(url, "PUT", arg...)
}

func Delete(url string, arg ...Option) (*Resp, error) {
	return newRequest(url, "DELETE", arg...)
}

func Patch(url string, arg ...Option) (*Resp, error) {
	return newRequest(url, "PATCH", arg...)
}

func Options(url string, arg ...Option) (*Resp, error) {
	return newRequest(url, "OPTIONS", arg...)
}

func WithParams(params map[string]string) Option {
	return func(arg *args) {
		arg.params = params
	}
}

func WithJson(json map[string]interface{}) Option {
	return func(arg *args) {
		arg.json = json
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(arg *args) {
		arg.headers = headers
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(arg *args) {
		arg.timeout = timeout
	}
}

func WithAllowRedirects(allowRedirects bool) Option {
	return func(arg *args) {
		arg.allowRedirects = allowRedirects
	}
}

func WithVerify(verify bool) Option {
	return func(arg *args) {
		arg.verify = verify
	}
}

func newRequest(url, method string, arg ...Option) (*Resp, error) {
	var (
		err     error
		body    io.Reader
		jsonVal []byte
		req  *http.Request
		resp *http.Response
	)

	ar := argsPool.Get().(*args)
	defer argsPool.Put(ar)
	ar.allowRedirects = true
	ar.verify = true
	ar.timeout = defaultTimeout
	ar.params = nil
	ar.json = nil
	ar.data = nil
	ar.headers = nil

	for _, a := range arg {
		a(ar)
	}

	if method == "GET" {

	} else if ar.json != nil {
		jsonVal, err = goJson.Marshal(ar.json)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonVal)
	} else if ar.data != nil {
		vals := make(netUrl.Values)
		for k, v := range ar.data {
			vals.Set(k, v)
		}
		body = strings.NewReader(vals.Encode())
	}

	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Close = true

	if method == "GET" {

	} else if ar.json != nil {
		req.Header.Set("Content-Type", "application/json")
	} else if ar.data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if ar.headers != nil {
		for k, v := range ar.headers {
			req.Header.Add(k, v)
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
		Timeout: ar.timeout,
	}

	if !ar.allowRedirects {
		cli.CheckRedirect = disableRedirect
	}

	if !ar.verify {
		cookieJar, _ := cookiejar.New(nil)
		cli.Jar = cookieJar
		cli.Transport = insecureTransport
	}

	resp, err = cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resP := new(Resp)
	resP.StatusCode = resp.StatusCode
	resP.header = resp.Header
	resP.Body, err = io.ReadAll(resp.Body)

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

func (r *Resp) Header() map[string]string {
	newHeader := make(map[string]string)
	for k, v := range r.header {
		newHeader[k] = strings.Join(v, ",")
	}

	return newHeader
}

func (r *Resp) Json(arg ...interface{}) (val map[string]interface{}, err error) {
	if len(arg) == 0 {
		err = goJson.Unmarshal(r.Body, &val)
	} else if len(arg) == 1 {
		err = goJson.Unmarshal(r.Body, arg[0])
	} else {
		return nil, errors.New("too many arg, max len is 1")
	}

	return
}

func (r *Resp) List(arg ...interface{}) (val []interface{}, err error) {
	if len(arg) == 0 {
		err = goJson.Unmarshal(r.Body, &val)
	} else if len(arg) == 1 {
		err = goJson.Unmarshal(r.Body, arg[0])
	} else {
		return nil, errors.New("too many arg, max len is 1")
	}

	return
}

func (r *Resp) JsonAndValueIsString() (val map[string]string, err error) {
	err = goJson.Unmarshal(r.Body, &val)
	return
}
