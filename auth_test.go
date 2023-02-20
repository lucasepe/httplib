package httplib_test

import (
	"fmt"
	"net/http"

	"github.com/lucasepe/httplib"
)

func ExampleTokenAuth() {
	auth := httplib.TokenAuth{
		Token: "H-E-L-L-O",
	}

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		panic(err)
	}
	auth.SetAuth(req)

	fmt.Print(req.Header)

	// Output:
	// map[Authorization:[Bearer H-E-L-L-O]]
}

func ExampleBasicAuth() {
	auth := httplib.BasicAuth{
		Username: "gopher",
		Password: "abbracadabbra!",
	}

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		panic(err)
	}
	auth.SetAuth(req)

	fmt.Print(req.Header)

	// Output:
	// map[Authorization:[Basic Z29waGVyOmFiYnJhY2FkYWJicmEh]]
}
