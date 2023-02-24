package httplib

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"
)

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:   true,
			MaxConnsPerHost:     100,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,

			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 50 * time.Second,
	}
}

// InsecureSkipVerify controls whether a http.Client verifies
// the server's certificate chain and host name.
func InsecureSkipVerify(c *http.Client, v bool) {
	t, ok := c.Transport.(*http.Transport)
	if !ok {
		return
	}

	if t.TLSClientConfig == nil {
		t.TLSClientConfig = &tls.Config{}
	}
	t.TLSClientConfig.InsecureSkipVerify = v
}

// HandleResponse composes the *http.Response validators with the handler.
func HandleResponse(res *http.Response, handler HandleResponseFunc, validators ...HandleResponseFunc) (err error) {
	if len(validators) == 0 {
		validators = []HandleResponseFunc{
			CheckStatus(
				http.StatusOK,
				http.StatusCreated,
				http.StatusAccepted,
				http.StatusNonAuthoritativeInfo,
				http.StatusNoContent,
			),
		}
	}
	err = ChainHandlers(validators...)(res)
	if err != nil {
		return err
	}

	if handler == nil {
		handler = consumeResponseBody
	}

	return handler(res)
}

type FireOptions struct {
	AuthMethod      AuthMethod
	ResponseHandler HandleResponseFunc
	Validators      []HandleResponseFunc
	Verbose         bool
}

// Fire calls the http.Client.Do() and validates and handles any resulting response.
// The response body is closed after all validators and the handler run.
func Fire(c *http.Client, req *http.Request, opts FireOptions) (err error) {
	if opts.AuthMethod != nil {
		opts.AuthMethod.SetAuth(req)
	}

	if opts.Verbose {
		DumpRequest(req, os.Stderr)
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if opts.Verbose {
		DumpResponse(res, os.Stderr, req.URL.Query().Get("watch") != "true")
	}

	return HandleResponse(res, opts.ResponseHandler, opts.Validators...)
}
