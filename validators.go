package httplib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// CheckStatus validates the response has an acceptable status code.
func CheckStatus(acceptStatuses ...int) HandleResponseFunc {
	return func(res *http.Response) error {
		for _, code := range acceptStatuses {
			if res.StatusCode == code {
				return nil
			}
		}

		return fmt.Errorf("%w: unexpected status: %d",
			&StatusError{StatusCode: res.StatusCode}, res.StatusCode)
	}
}

// ErrorJSON validates the response has an acceptable status
// code and if it's bad, attempts to marshal the JSON
// into the error object provided.
func ErrorJSON(v error, acceptStatuses ...int) HandleResponseFunc {
	return func(res *http.Response) error {
		for _, code := range acceptStatuses {
			if res.StatusCode == code {
				return nil
			}
		}

		if res.Body == nil {
			return &StatusError{StatusCode: res.StatusCode}
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return &StatusError{StatusCode: res.StatusCode, Inner: err}
		}

		if err = json.Unmarshal(data, &v); err != nil {
			return &StatusError{StatusCode: res.StatusCode, Inner: err}
		}

		return &StatusError{StatusCode: res.StatusCode, Inner: v}
	}
}

type StatusError struct {
	StatusCode int
	Inner      error
}

func (e *StatusError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("unexpected status: %d: %v", e.StatusCode, e.Inner)
	}
	return fmt.Sprintf("unexpected status: %d:", e.StatusCode)
}

func (e *StatusError) Unwrap() error {
	return e.Inner
}

// HasStatusErr returns true if err is a ResponseError caused by any of the codes given.
func HasStatusErr(err error, codes ...int) bool {
	if err == nil {
		return false
	}
	if se := new(StatusError); errors.As(err, &se) {
		for _, code := range codes {
			if se.StatusCode == code {
				return true
			}
		}
	}
	return false
}

func IsNotFoundError(err error) bool {
	return HasStatusErr(err, http.StatusNotFound)
}
