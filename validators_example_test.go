package httplib_test

import (
	"fmt"

	"github.com/lucasepe/httplib"
)

func ExampleIsNotFoundError() {
	ub := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: "http://example.com",
		Path:    "404",
	})

	req, err := httplib.NewGetRequest(ub)
	if err != nil {
		panic(err)
	}

	err = httplib.Fire(httplib.NewClient(), req, httplib.FireOptions{
		Validators: []httplib.HandleResponseFunc{
			httplib.CheckStatus(200),
		},
	})
	if err != nil {
		if httplib.IsNotFoundError(err) {
			fmt.Println("got a 404")
		}
	}

	// Output:
	// got a 404
}
