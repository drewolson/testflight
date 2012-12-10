# testflight

## Installation

```bash
go get github.com/drewolson/testflight
```

```go
import "github.com/drewolson/testflight"
```

## Usage

testflight makes it simple to test your http servers in Go. Suppose you're using [pat](https://github.com/bmizerany/pat) to create a simple http handler, like so:

```go
func Handler() http.Handler {
	m := pat.New()

	m.Get("/hello/:name", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hello, "+req.URL.Query().Get(":name"))
	}))

	m.Post("/post/form", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		name := req.Form.Get("name")
		w.WriteHeader(201)
		io.WriteString(w, name+" created")
	}))

	return m
}
```

Let's use testflight to test our handler. Keep in mind that testflight is doing full-stack http tests. We're also using assert for test assertions.

```go
func TestGet(t *testing.T) {
	testflight.WithServer(Handler(), func(r *testflight.Requester) {
		response := r.Get("/hello/drew")

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "hello, drew", response.Body)
	})
}

func TestPostWithForm(t *testing.T) {
	testflight.WithServer(handler(), func(r *testflight.Requester) {
		response := r.Post("/post/form", testflight.FORM_ENCODED, "name=Drew")

		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "Drew created", response.Body)
	})
}
```

The testflight.Requester class has the following methods: Get, Post, Put, Delete and Do. Do accepts an *http.Request for times when you need more explicit control of our request. See testflight_test.go for more usage information.

## Contributing

First, run the tests.

```bash
mkdir testflight
cd testflight

export GOPATH=`pwd`

go get github.com/drewolson/testflight
go get github.com/kr/pretty
go get github.com/bmizerany/assert
go get github.com/bmizerany/pat

go test github.com/drewolson/testflight
```

Now write new tests, fix them and send me a pull request!
