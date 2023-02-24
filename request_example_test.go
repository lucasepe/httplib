package httplib_test

import (
	"fmt"

	"github.com/lucasepe/httplib"
)

func ExampleGet() {
	url, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: "http://httpbin.org",
		Path:    "user-agent",
	}).Build()
	if err != nil {
		panic(err)
	}
	req, err := httplib.Get(url.String())
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

func ExamplePost() {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	url, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: "http://httpbin.org",
		Path:    "post",
	}).Build()
	if err != nil {
		panic(err)
	}

	bodyFn := httplib.ToJSON(&Login{
		Username: "pinco.pallo@gmail.com",
		Password: "abbracadabbra",
	})

	req, err := httplib.Post(url.String(), bodyFn)
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
