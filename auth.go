package httplib

import (
	"fmt"
	"net/http"
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
