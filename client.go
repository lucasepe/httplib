package httplib

import (
	"crypto/tls"
	"net"
	"net/http"
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

type FireOptions struct {
	Request         *http.Request
	AuthMethod      AuthMethod
	ResponseHandler HandleResponseFunc
	Validators      []HandleResponseFunc
	Verbose         bool
}

// Fire calls the http.Client.Do() and validates and handles any resulting response.
// The response body is closed after all validators and the handler run.
func Fire(c *http.Client, opts FireOptions) (err error) {
	if opts.AuthMethod != nil {
		opts.AuthMethod.SetAuth(opts.Request)
	}

	if opts.Verbose {
		dumpRequest(opts.Request)
	}

	res, err := c.Do(opts.Request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if opts.Verbose {
		dumpResponse(res, opts.Request.URL.Query().Get("watch") != "true")
	}

	if len(opts.Validators) == 0 {
		opts.Validators = []HandleResponseFunc{
			CheckStatus(
				http.StatusOK,
				http.StatusCreated,
				http.StatusAccepted,
				http.StatusNonAuthoritativeInfo,
				http.StatusNoContent,
			),
		}
	}
	err = ChainHandlers(opts.Validators...)(res)
	if err != nil {
		return err
	}

	handle := opts.ResponseHandler
	if handle == nil {
		handle = consumeResponseBody
	}

	return handle(res)
}
