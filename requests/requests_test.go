package requests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestGet(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"message":"ok"}`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	val, err := resp.Json()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if val["message"] != "ok" {
		t.Errorf("expected message=ok, got %v", val["message"])
	}
}

func TestPost(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content-type, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(body, &data)

		w.WriteHeader(201)
		w.Write(body)
	})
	defer ts.Close()

	resp, err := Post(ts.URL, WithJson(map[string]interface{}{
		"name": "teheran",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	val, err := resp.Json()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if val["name"] != "teheran" {
		t.Errorf("expected name=teheran, got %v", val["name"])
	}
}

func TestPut(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"updated":true}`))
	})
	defer ts.Close()

	resp, err := Put(ts.URL, WithJson(map[string]interface{}{"key": "val"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(204)
	})
	defer ts.Close()

	resp, err := Delete(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestPatch(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"patched":true}`))
	})
	defer ts.Close()

	resp, err := Patch(ts.URL, WithJson(map[string]interface{}{"field": "new"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestOptions(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			t.Errorf("expected OPTIONS, got %s", r.Method)
		}
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(200)
	})
	defer ts.Close()

	resp, err := Options(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	h := resp.Header()
	if h["Allow"] != "GET, POST, OPTIONS" {
		t.Errorf("expected Allow header, got %v", h["Allow"])
	}
}

func TestWithParams(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("foo") != "bar" {
			t.Errorf("expected query param foo=bar, got %s", r.URL.Query().Get("foo"))
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL, WithParams(map[string]string{"foo": "bar"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestWithHeaders(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom") != "test-value" {
			t.Errorf("expected X-Custom header, got %s", r.Header.Get("X-Custom"))
		}
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL, WithHeaders(map[string]string{"X-Custom": "test-value"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestWithTimeout(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
	})
	defer ts.Close()

	_, err := Get(ts.URL, WithTimeout(50*time.Millisecond))
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestWithAllowRedirectsDisabled(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/redirected", http.StatusFound)
			return
		}
		w.WriteHeader(200)
	})
	defer ts.Close()

	resp, err := Get(ts.URL, WithAllowRedirects(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("expected status 302, got %d", resp.StatusCode)
	}
}

func TestWithAllowRedirectsEnabled(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/redirected", http.StatusFound)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"redirected":true}`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL, WithAllowRedirects(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200 after redirect, got %d", resp.StatusCode)
	}
}

func TestRespJsonWithStruct(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"name":"teheran","version":"1.0"}`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	type Info struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	var info Info
	_, err = resp.Json(&info)
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if info.Name != "teheran" || info.Version != "1.0" {
		t.Errorf("unexpected struct values: %+v", info)
	}
}

func TestRespJsonTooManyArgs(t *testing.T) {
	r := &Resp{Body: []byte(`{}`)}
	_, err := r.Json("a", "b")
	if err == nil {
		t.Error("expected error for too many args")
	}
}

func TestRespList(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`[1,2,3]`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, err := resp.List()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
}

func TestRespJsonAndValueIsString(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"key":"value","foo":"bar"}`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, err := resp.JsonAndValueIsString()
	if err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if val["key"] != "value" || val["foo"] != "bar" {
		t.Errorf("unexpected values: %v", val)
	}
}

func TestRespHeader(t *testing.T) {
	ts := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "hello")
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	})
	defer ts.Close()

	resp, err := Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h := resp.Header()
	if h["X-Test"] != "hello" {
		t.Errorf("expected X-Test=hello, got %s", h["X-Test"])
	}
}

func TestInvalidURL(t *testing.T) {
	_, err := Get("http://[::1]:0/invalid")
	if err == nil {
		t.Error("expected error for invalid url")
	}
}
