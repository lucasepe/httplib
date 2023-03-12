package httplib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

// GetBodyFunc provides a Builder with a source for a request body.
type GetBodyFunc func() (io.ReadCloser, error)

// ToJSON is a GetBodyFunc that marshals a JSON object.
func ToJSON(v any) GetBodyFunc {
	return func() (io.ReadCloser, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return nopCloser{bytes.NewReader(b)}, nil
	}
}

// FormData is a GetBodyFunc that builds an encoded form body.
func FormData(data url.Values) GetBodyFunc {
	return func() (r io.ReadCloser, err error) {
		return nopCloser{strings.NewReader(data.Encode())}, nil
	}
}

// nopCloser is like io.NopCloser(),
// but it is a concrete type so we can strip it out
// before setting a body on a request.
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

var _ io.ReadCloser = nopCloser{}
