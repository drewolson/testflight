package testflight

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type Requester struct {
	server *httptest.Server
}

func (requester *Requester) Get(route string, mods ...func(*http.Request) error) *Response {
	return requester.performRequest("GET", route, "", "", mods...)
}

func (requester *Requester) Post(route, contentType, body string, mods ...func(*http.Request) error) *Response {
	return requester.performRequest("POST", route, contentType, body, mods...)
}

func (requester *Requester) Put(route, contentType, body string, mods ...func(*http.Request) error) *Response {
	return requester.performRequest("PUT", route, contentType, body, mods...)
}

func (requester *Requester) Patch(route, contentType, body string, mods ...func(*http.Request) error) *Response {
	return requester.performRequest("PATCH", route, contentType, body, mods...)
}

func (requester *Requester) Delete(route, contentType, body string, mods ...func(*http.Request) error) *Response {
	return requester.performRequest("DELETE", route, contentType, body, mods...)
}

func (requester *Requester) Do(request *http.Request) *Response {
	fullUrl, err := url.Parse(requester.httpUrl(request.URL.String()))
	if err != nil {
		panic(err)
	}

	request.URL = fullUrl
	return requester.sendRequest(request)
}

func (requester *Requester) Url(route string) string {
	return requester.server.Listener.Addr().String() + route
}

func (requester *Requester) performRequest(httpAction, route, contentType, body string, mods ...func(*http.Request) error) *Response {
	request, err := http.NewRequest(httpAction, requester.httpUrl(route), strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	// call each modifier so they can manipulate our request before sending
	// if a modifier fails we panic on its error code
	for _, op := range mods {
		err := op(request)
		if err != nil {
			panic(err)
		}
	}

	request.Header.Add("Content-Type", contentType)
	return requester.sendRequest(request)
}

func (requester *Requester) sendRequest(request *http.Request) *Response {
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	return newResponse(response)
}

func (requester *Requester) httpUrl(route string) string {
	return "http://" + requester.Url(route)
}
