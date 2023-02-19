package httplib_test

import (
	"fmt"

	"github.com/lucasepe/httplib"
)

func ExampleNewURL() {
	url, err := httplib.NewURL(
		"https://dev.azure.com",
		"my-great-org/projects",
		"api_version", "7.0")
	if err != nil {
		panic(err)
	}

	fmt.Print(url.String())

	// Output:
	// https://dev.azure.com/my-great-org/projects?api_version=7.0
}
