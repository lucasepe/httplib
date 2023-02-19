package httplib

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

// AuthMethod is concrete implementation of common.AuthMethod for HTTP services
type AuthMethod interface {
	SetAuth(r *http.Request)
}

// BasicAuth represent a HTTP basic auth
type BasicAuth struct {
	Username, Password string
}

func (a *BasicAuth) SetAuth(r *http.Request) {
	if a == nil {
		return
	}

	r.SetBasicAuth(a.Username, a.Password)
}

// TokenAuth implements an http.AuthMethod that can be used with http transport
// to authenticate with HTTP token authentication (also known as bearer
// authentication).
type TokenAuth struct {
	Token string
}

func (a *TokenAuth) SetAuth(r *http.Request) {
	if a == nil {
		return
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

// GetBodyFunc provides a Builder with a source for a request body.
type GetBodyFunc func() (io.ReadCloser, error)

// NewRequest creates a new http.Request with specified method, uri and request body.
func NewRequest(method, uri string, getBodyFn GetBodyFunc) (req *http.Request, err error) {
	var body io.Reader
	if getBodyFn != nil {
		if body, err = getBodyFn(); err != nil {
			return nil, err
		}
		if nopper, ok := body.(nopCloser); ok {
			body = nopper.Reader
		}
	}

	req, err = http.NewRequest(method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.GetBody = getBodyFn

	return req, nil
}

// nopCloser is like io.NopCloser(),
// but it is a concrete type so we can strip it out
// before setting a body on a request.
// See https://github.com/carlmjohnson/requests/discussions/49
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

var _ io.ReadCloser = nopCloser{}

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
