package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cirello.io/bloomfilterd/internal/filter"
)

type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("errReader errored")
}

type errWriter struct {
	headers      http.Header
	responseCode int
}

func (e *errWriter) Header() http.Header {
	return e.headers
}
func (e *errWriter) WriteHeader(i int) {
	e.responseCode = i
}
func (errWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("errWriter errored")
}

func TestHttp(t *testing.T) {
	d := &daemon{
		filters: make(map[string]*filter.Bloomfilter),
	}

	r, _ := http.NewRequest("POST", "/add", strings.NewReader(`{"name":"default", "size":4096, "hashcount":8}`))
	rec := httptest.NewRecorder()
	d.ServeHTTP(rec, r)
	if _, ok := d.filters["default"]; !ok {
		t.Error("default filter not created")
	}

	rErr, _ := http.NewRequest("POST", "/add", strings.NewReader(`{"name":"default", "size":4096, "hashcount":8`))
	recErr := httptest.NewRecorder()
	d.ServeHTTP(recErr, rErr)
	if recErr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("invalid JSON not dealt with properly")
	}

	rOverwrite, _ := http.NewRequest("POST", "/add", strings.NewReader(`{"name":"default", "size":4096, "hashcount":8}`))
	recOverwrite := httptest.NewRecorder()
	d.ServeHTTP(recOverwrite, rOverwrite)
	if recOverwrite.Result().StatusCode != http.StatusBadRequest {
		t.Error("Overwrites must always return StatusBadRequest")
	}

	rEncErr, _ := http.NewRequest("POST", "/add", strings.NewReader(`{"name":"default-2", "size":4096, "hashcount":8}`))
	errW := &errWriter{headers: make(map[string][]string)}
	d.ServeHTTP(errW, rEncErr)
	if errW.responseCode != http.StatusInternalServerError {
		t.Error("JSON encoding errors must always return StatusInternalServerError", errW.responseCode)
	}

	rList, _ := http.NewRequest("GET", "/list", nil)
	recList := httptest.NewRecorder()
	d.ServeHTTP(recList, rList)
	b, _ := ioutil.ReadAll(recList.Result().Body)
	if fmt.Sprintf("%s", b) != `{"default":0,"default-2":0}`+"\n" {
		t.Errorf("did not get valid filter list: %s %v", b, len(b))
	}

	rListErr, _ := http.NewRequest("GET", "/list", nil)
	errListW := &errWriter{headers: make(map[string][]string)}
	d.ServeHTTP(errListW, rListErr)
	if errListW.responseCode != http.StatusInternalServerError {
		t.Error("JSON encoding errors must always return StatusInternalServerError", errW.responseCode)
	}
}

func TestHttpFilter(t *testing.T) {
	d := &daemon{
		filters: map[string]*filter.Bloomfilter{
			"default": filter.New(4096, 8),
		},
	}

	calls := []struct {
		url        string
		method     string
		body       io.Reader
		statusCode int
		result     *string
	}{
		{"/filter/invalid", "POST", strings.NewReader("test"), http.StatusNotFound, nil},
		{"/filter/default", "POST", strings.NewReader("test"), http.StatusNoContent, nil},
		{"/filter/default", "POST", new(errReader), http.StatusInternalServerError, nil},
		{"/filter/default?body=test", "GET", nil, http.StatusOK, pstr(`{"Name":"default","OK":true}`)},
		{"/filter/default", "DELETE", strings.NewReader("test"), http.StatusNoContent, nil},
		{"/filter/default", "DELETE", new(errReader), http.StatusInternalServerError, nil},
		{"/filter/default?body=test", "GET", nil, http.StatusOK, pstr(`{"Name":"default","OK":false}`)},
		{"/filter/default", "PATCH", nil, http.StatusMethodNotAllowed, nil},
	}

	for _, call := range calls {
		r, _ := http.NewRequest(call.method, call.url, call.body)
		rec := httptest.NewRecorder()
		d.ServeHTTP(rec, r)
		if sc := rec.Result().StatusCode; sc != call.statusCode {
			t.Logf("call %#v - status code - got: %v expected: %v", call, sc, call.statusCode)
		}
		if call.result != nil {
			b, _ := ioutil.ReadAll(rec.Result().Body)
			if s := strings.TrimSpace(string(b)); s != *call.result {
				t.Logf("call %#v - body - got: %v expected: %v", call, s, *call.result)
			}
		}
	}

	r, _ := http.NewRequest("GET", "/filter/default?body=test", nil)
	errW := &errWriter{headers: make(map[string][]string)}
	d.ServeHTTP(errW, r)
	if errW.responseCode != http.StatusInternalServerError {
		t.Error("JSON encoding errors must always return StatusInternalServerError", errW.responseCode)
	}
}

func pstr(str string) *string {
	return &str
}
