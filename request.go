package httplib

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

func NewGetRequest(ub URLBuilder) (req *http.Request, err error) {
	return newRequest(http.MethodGet, ub, nil)
}

func NewPostRequest(ub URLBuilder, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	return newRequest(http.MethodPost, ub, getBodyFn)
}

func NewPutRequest(ub URLBuilder, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	return newRequest(http.MethodPut, ub, getBodyFn)
}

func NewDeleteRequest(ub URLBuilder) (req *http.Request, err error) {
	return newRequest(http.MethodDelete, ub, nil)
}

// newRequest creates a new http.Request with specified method, uri and request body.
func newRequest(method string, ub URLBuilder, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	url, err := ub.Build()
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if getBodyFn != nil {
		if body, err = getBodyFn(); err != nil {
			return nil, err
		}
		if nopper, ok := body.(nopCloser); ok {
			body = nopper.Reader
		}
	}

	req, err = http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.GetBody = getBodyFn

	return req, nil
}

// dumpRequest dumps http.Request to os.Stderr.
func dumpRequest(req *http.Request) {
	buf, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		os.Stderr.WriteString("error dumping *http.Request")
		return
	}
	os.Stderr.Write(buf)
	os.Stderr.Write([]byte{'\n'})
}
