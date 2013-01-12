package testflight

import (
	"code.google.com/p/go.net/websocket"
)

type WebsocketResponse struct {
	RawConn *websocket.Conn
}

func newWebsocketResponse(conn *websocket.Conn) *WebsocketResponse {
	return &WebsocketResponse{
		RawConn: conn,
	}
}

func (response *WebsocketResponse) ReceiveMessage() string {
	var message string
	websocket.Message.Receive(response.RawConn, &message)
	return message
}

func (response *WebsocketResponse) WriteMessage(message string) {
	websocket.Message.Send(response.RawConn, message)
}
