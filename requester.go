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

func (requester *Requester) Get(route string) *Response {
	return requester.performRequest("GET", route, "", "")
}

func (requester *Requester) Post(route, contentType, body string) *Response {
	return requester.performRequest("POST", route, contentType, body)
}

func (requester *Requester) Put(route, contentType, body string) *Response {
	return requester.performRequest("PUT", route, contentType, body)
}

func (requester *Requester) Patch(route, contentType, body string) *Response {
	return requester.performRequest("PATCH", route, contentType, body)
}

func (requester *Requester) Delete(route, contentType, body string) *Response {
	return requester.performRequest("DELETE", route, contentType, body)
}

func (requester *Requester) DoWithClient(request *http.Request, client http.Client) *Response {
	fullUrl, err := url.Parse(requester.httpUrl(request.URL.String()))
	if err != nil {
		panic(err)
	}

	request.URL = fullUrl
	return requester.sendRequestWithClient(request, client)
}

func (requester *Requester) Do(request *http.Request) *Response {
	client := http.Client{}
	return requester.DoWithClient(request, client)
}

func (requester *Requester) Url(route string) string {
	return requester.server.Listener.Addr().String() + route
}

func (requester *Requester) performRequest(httpAction, route, contentType, body string) *Response {
	request, err := http.NewRequest(httpAction, requester.httpUrl(route), strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", contentType)
	return requester.sendRequest(request)
}

func (requester *Requester) sendRequest(request *http.Request) *Response {
	client := http.Client{}
	return requester.sendRequestWithClient(request, client)
}

func (requester *Requester) sendRequestWithClient(request *http.Request, client http.Client) *Response {
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	return newResponse(response)
}

func (requester *Requester) httpUrl(route string) string {
	return "http://" + requester.Url(route)
}
