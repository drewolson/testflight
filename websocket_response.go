package testflight

import (
	"code.google.com/p/go.net/websocket"
	"io/ioutil"
)

type WebsocketResponse struct {
	RawConn *websocket.Conn
}

func newWebsocketResponse(conn *websocket.Conn) *WebsocketResponse {
	return &WebsocketResponse{
		RawConn: conn,
	}
}

func (response *WebsocketResponse) Read() string {
	message, _ := ioutil.ReadAll(response.RawConn)
	return string(message)
}
