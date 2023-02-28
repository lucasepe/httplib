package httplib

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func Get(url string) (req *http.Request, err error) {
	return newRequest(http.MethodGet, url, nil)
}

func Post(url string, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	return newRequest(http.MethodPost, url, getBodyFn)
}

func Put(url string, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	return newRequest(http.MethodPut, url, getBodyFn)
}

func Patch(url string, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	return newRequest(http.MethodPatch, url, getBodyFn)
}

func Delete(url string) (req *http.Request, err error) {
	return newRequest(http.MethodDelete, url, nil)
}

func AddHeader(req *http.Request, vv ...string) {
	if len(vv)%2 != 0 {
		return
	}
	for i := 0; i < len(vv); i += 2 {
		req.Header.Add(vv[i], vv[i+1])
	}
}

func SetHeader(req *http.Request, vv ...string) {
	if len(vv)%2 != 0 {
		return
	}
	for i := 0; i < len(vv); i += 2 {
		req.Header.Set(vv[i], vv[i+1])
	}
}

// DumpRequest dumps http.Request to the specified io.Writer.
func DumpRequest(req *http.Request, wri io.Writer) {
	buf, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Fprintf(wri, "error dumping *http.Request: %s", err.Error())
		return
	}
	wri.Write(buf)
	wri.Write([]byte{'\n'})
}

// newRequest creates a new http.Request with specified method, uri and request body.
func newRequest(method string, url string, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	var body io.Reader
	if getBodyFn != nil {
		if body, err = getBodyFn(); err != nil {
			return nil, err
		}
		if nopper, ok := body.(nopCloser); ok {
			body = nopper.Reader
		}
	}

	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.GetBody = getBodyFn

	return req, nil
}
