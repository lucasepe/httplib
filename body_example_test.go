package httplib_test

import (
	"fmt"
	"io"

	"github.com/lucasepe/httplib"
)

func ExampleToJSON() {
	type LoginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	login := LoginData{
		Username: "pinco.pallo@gmail.com",
		Password: "abbracadabbra",
	}

	inp, err := httplib.ToJSON(&login)()
	if err != nil {
		panic(err)
	}

	dat, err := io.ReadAll(inp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", dat)

	// Output:
	// {"username":"pinco.pallo@gmail.com","password":"abbracadabbra"}
}
