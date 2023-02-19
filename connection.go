package httplib

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type ConnectionOptions struct {
	HttpClient        *http.Client
	HttpClientTimeout time.Duration
	AuthMethod        AuthMethod
	BaseURL           string
	Verbose           bool
}

type Connection struct {
	httpClient *http.Client
	authMethod AuthMethod
	baseURL    string
	verbose    bool
}

func NewConnection(opts ConnectionOptions) *Connection {
	if opts.HttpClientTimeout <= 0 {
		opts.HttpClientTimeout = 50 * time.Second
	}

	if opts.HttpClient == nil {
		opts.HttpClient = &http.Client{
			Transport: defaultTransport(),
			Timeout:   50 * time.Second,
		}
	}

	return &Connection{
		httpClient: opts.HttpClient,
		authMethod: opts.AuthMethod,
		baseURL:    opts.BaseURL,
		verbose:    opts.Verbose,
	}
}

// Insecure controls whether a client verifies the server's
// certificate chain and host name.
func (c *Connection) Insecure(v bool) *Connection {
	t, ok := c.httpClient.Transport.(*http.Transport)
	if !ok {
		return c
	}

	if t.TLSClientConfig == nil {
		t.TLSClientConfig = &tls.Config{}
	}
	t.TLSClientConfig.InsecureSkipVerify = v

	return c
}

// Do calls the underlying http.Client and validates and handles any resulting response.
// The response body is closed after all validators and the handler run.
func (c *Connection) Do(req *http.Request, handleResponseFn HandleResponseFunc, validators ...HandleResponseFunc) (err error) {
	if c.authMethod != nil {
		c.authMethod.SetAuth(req)
	}

	if c.verbose {
		dumpRequest(req)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if c.verbose {
		dumpResponse(res, req.URL.Query().Get("watch") != "true")
	}

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
	if err = ChainHandlers(validators...)(res); err != nil {
		return err
	}

	handle := handleResponseFn
	if handle == nil {
		handle = consumeResponseBody
	}

	return handle(res)
}

// defaultTransport is the default implementation of Transport and is
// used by HTTPConnection when http.Client is <nil>.
// It uses HTTP proxies as directed by the $HTTP_PROXY and $NO_PROXY (or $http_proxy and
// $no_proxy) environment variables.
func defaultTransport() http.RoundTripper {
	return &http.Transport{
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
	}
}
