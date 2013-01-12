package testflight

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type Requester struct {
	Server *httptest.Server
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
	return "http://" + requester.Server.Listener.Addr().String() + route
}
