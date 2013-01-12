package ws

import (
	"code.google.com/p/go.net/websocket"
)

type Connection struct {
	RawConn *websocket.Conn
}

func newConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		RawConn: conn,
	}
}

func (connection *Connection) ReceiveMessage() string {
	var message string
	websocket.Message.Receive(connection.RawConn, &message)
	return message
}

func (connection *Connection) WriteMessage(message string) {
	websocket.Message.Send(connection.RawConn, message)
}
