package httplib

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

// HandleResponseFunc is used to validate or handle the response to a request.
type HandleResponseFunc func(*http.Response) error

// ChainHandlers allows for the composing of validators or response handlers.
func ChainHandlers(handlers ...HandleResponseFunc) HandleResponseFunc {
	return func(r *http.Response) error {
		for _, h := range handlers {
			if h == nil {
				continue
			}
			if err := h(r); err != nil {
				return err
			}
		}
		return nil
	}
}

// FromJSON decodes a response as a JSON object.
func FromJSON(v any) HandleResponseFunc {
	return func(res *http.Response) error {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, v); err != nil {
			return err
		}
		return nil
	}
}

func consumeResponseBody(res *http.Response) (err error) {
	const maxDiscardSize = 640 * 1 << 10
	if _, err = io.CopyN(io.Discard, res.Body, maxDiscardSize); err == io.EOF {
		err = nil
	}
	return err
}

// dumpResponse dumps http.Response to os.Stderr.
func dumpResponse(res *http.Response, body bool) {
	buf, err := httputil.DumpResponse(res, body)
	if err != nil {
		os.Stderr.WriteString("error dumping *http.Response")
		return
	}
	os.Stderr.Write(buf)
	os.Stderr.Write([]byte{'\n'})
}
