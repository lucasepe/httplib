package httplib_test

import (
	"fmt"

	"github.com/lucasepe/httplib"
)

func ExampleNewURLBuilder() {
	ub := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: "https://dev.my.site",
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
	// https://dev.my.site/my-great-org/projects?api_version=7.0
}
