package httplib_test

import (
	"fmt"

	"github.com/lucasepe/httplib"
)

func ExampleNewGetRequest() {
	ub := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: "http://httpbin.org",
		Path:    "user-agent",
	})
	req, err := httplib.NewGetRequest(ub)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "httplib.Client Example")

	var res map[string]string
	err = httplib.Fire(httplib.NewClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(&res),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res["user-agent"])
	// Output:
	// httplib.Client Example
}

func ExampleNewPostRequest() {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	ub := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: "http://httpbin.org",
		Path:    "post",
	})

	bodyFn := httplib.ToJSON(&Login{
		Username: "pinco.pallo@gmail.com",
		Password: "abbracadabbra",
	})

	req, err := httplib.NewPostRequest(ub, bodyFn)
	if err != nil {
		panic(err)
	}

	var res map[string]any
	err = httplib.Fire(httplib.NewClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(&res),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res["json"])
	// Output:
	// map[password:abbracadabbra username:pinco.pallo@gmail.com]
}
