package testflight

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type person struct {
	Name string `json:"name"`
}

func handler() http.Handler {
	m := pat.New()

	m.Get("/hello/:name", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hello, "+req.URL.Query().Get(":name"))
	}))

	m.Post("/post/json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		person := &person{}
		body, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(body, person)
		w.WriteHeader(201)
		io.WriteString(w, person.Name+" created")
	}))

	return m
}

func TestGet(t *testing.T) {
	WithServer(handler(), func(r *Requester) {
		response := r.Get("/hello/drew")

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "hello, drew", response.Body)
	})
}

func TestPostWithJson(t *testing.T) {
	WithServer(handler(), func(r *Requester) {
		response := r.Post("/post/json", "application/json", `{"name": "Drew"}`)

		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "Drew created", response.Body)
	})
}
