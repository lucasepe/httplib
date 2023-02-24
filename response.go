package httplib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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

// DumpResponse dumps http.Response to the specified io.Writer.
func DumpResponse(res *http.Response, wri io.Writer, body bool) {
	buf, err := httputil.DumpResponse(res, body)
	if err != nil {
		fmt.Fprintf(wri, "error dumping *http.Response: %s\n", err.Error())
		return
	}
	wri.Write(buf)
	wri.Write([]byte{'\n'})
}

func consumeResponseBody(res *http.Response) (err error) {
	const maxDiscardSize = 640 * 1 << 10
	if _, err = io.CopyN(io.Discard, res.Body, maxDiscardSize); err == io.EOF {
		err = nil
	}
	return err
}
