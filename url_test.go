package httplib_test

import (
	"fmt"

	"github.com/lucasepe/httplib"
)

func ExampleUrlBuilder() {
	ub := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: "https://dev.azure.com",
		Path:    "my-great-org/projects",
		Params: []string{
			"api_version", "7.0",
		},
	})
	url, err := ub.Build()
	if err != nil {
		panic(err)
	}

	fmt.Print(url.String())

	// Output:
	// https://dev.azure.com/my-great-org/projects?api_version=7.0
}
