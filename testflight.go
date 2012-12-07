package testflight

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

const (
	JSON         = "application/json"
	FORM_ENCODED = "application/x-www-form-urlencoded"
)

type Response struct {
	Body        string
	RawResponse *http.Response
	StatusCode  int
}

func newResponse(response *http.Response) *Response {
	body, _ := ioutil.ReadAll(response.Body)
	return &Response{
		Body:        string(body),
		RawResponse: response,
		StatusCode:  response.StatusCode,
	}
}

type Requester struct {
	server *httptest.Server
}

func (requester *Requester) Get(route string) *Response {
	return requester.performRequest("GET", route, "", "")
}

func (requester *Requester) Post(route, contentType, body string) *Response {
	return requester.performRequest("POST", route, contentType, body)
}

func (requester *Requester) Put(route, contentType, body string) *Response {
	return requester.performRequest("PUT", route, contentType, body)
}

func (requester *Requester) Delete(route, contentType, body string) *Response {
	return requester.performRequest("DELETE", route, contentType, body)
}

func (requester *Requester) Do(request *http.Request) *Response {
	fullUrl, _ := url.Parse(requester.url(request.URL.String()))
	request.URL = fullUrl
	return requester.sendRequest(request)
}

func (requester *Requester) performRequest(httpAction, route, contentType, body string) *Response {
	request, _ := http.NewRequest(httpAction, requester.url(route), strings.NewReader(body))
	request.Header.Add("Content-Type", contentType)
	return requester.sendRequest(request)
}

func (requester *Requester) sendRequest(request *http.Request) *Response {
	client := http.Client{}
	response, _ := client.Do(request)
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
