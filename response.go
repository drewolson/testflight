package testflight

import (
	"io/ioutil"
	"net/http"
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
