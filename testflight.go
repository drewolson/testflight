package testflight

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Response struct {
	RawResponse *http.Response
	Body        string
	StatusCode  int
}

func newResponse(response *http.Response) *Response {
	body, _ := ioutil.ReadAll(response.Body)
	return &Response{
		RawResponse: response,
		Body:        string(body),
		StatusCode:  response.StatusCode,
	}
}

type Requester struct {
	server *httptest.Server
}

func (requester *Requester) Get(route string) *Response {
	response, _ := http.Get(requester.url(route))
	return newResponse(response)
}

func (requester *Requester) Post(route string, contentType string, postBody string) *Response {
	response, _ := http.Post(requester.url(route), contentType, strings.NewReader(postBody))
	return newResponse(response)
}

func (requester *Requester) url(route string) string {
	return "http://" + requester.server.Listener.Addr().String() + route
}

func WithServer(handler http.Handler, context func(*Requester)) {
	server := httptest.NewServer(handler)
	defer server.Close()

	requester := &Requester{server: server}
	context(requester)
}
