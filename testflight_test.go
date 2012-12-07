package testflight

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/bmizerany/pat"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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

	m.Post("/post/form", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		name := req.Form.Get("name")
		w.WriteHeader(201)
		io.WriteString(w, name+" created")
	}))

	m.Put("/put/json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		person := &person{}
		body, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(body, person)
		w.WriteHeader(200)
		io.WriteString(w, person.Name+" updated")
	}))

	m.Del("/delete/json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		person := &person{}
		body, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(body, person)
		w.WriteHeader(200)
		io.WriteString(w, person.Name+" deleted")
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
		response := r.Post("/post/json", JSON, `{"name": "Drew"}`)

		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "Drew created", response.Body)
	})
}

func TestPostWithForm(t *testing.T) {
	WithServer(handler(), func(r *Requester) {
		response := r.Post("/post/form", FORM_ENCODED, "name=Drew")

		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "Drew created", response.Body)
	})
}

func TestPut(t *testing.T) {
	WithServer(handler(), func(r *Requester) {
		response := r.Put("/put/json", JSON, `{"name": "Drew"}`)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Drew updated", response.Body)
	})
}

func TestDelete(t *testing.T) {
	WithServer(handler(), func(r *Requester) {
		response := r.Delete("/delete/json", JSON, `{"name": "Drew"}`)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Drew deleted", response.Body)
	})
}

func TestDo(t *testing.T) {
	WithServer(handler(), func(r *Requester) {
		request, _ := http.NewRequest("DELETE", "/delete/json", strings.NewReader(`{"name": "Drew"}`))
		request.Header.Add("Content-Type", JSON)

		response := r.Do(request)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Drew deleted", response.Body)
	})
}
