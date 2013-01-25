package ws

import (
	"code.google.com/p/go.net/websocket"
	"time"
)

type Connection struct {
	RawConn *websocket.Conn
	ReceivedMessages []string
	unreadMessages []string
}

type TooShortError struct {

}

type TimeoutError struct {

}

func (e TimeoutError) Error() string {
	return "No Message Received in 30 seconds"
}

func (e TooShortError) Error() string {
	return "Unread Messages too Short"
}

func newConnection(conn *websocket.Conn) *Connection {
	connection := &Connection{
		RawConn: conn,
	}
	go connection.pollForMessages()
	return connection
}

func (connection *Connection) pollForMessages() {
	for {
		message := connection.receiveMessage()
		if message != "" {
			connection.ReceivedMessages = append(connection.ReceivedMessages, message)
			connection.unreadMessages = append(connection.unreadMessages, message)
		}
	}
}

func (connection *Connection) ReceiveMessage() (string, *TimeoutError) {
	for i := 0; i <= 1; i++ {
		message, poperr := connection.popUnreadMessage()
		if poperr == nil {
			return message, nil
		}
		time.Sleep(1 * time.Second)
	}
	return "", &TimeoutError{}
}

func(connection *Connection) popUnreadMessage() (string, *TooShortError) {
	if len(connection.unreadMessages) < 1 {
		return "", &TooShortError{}
	} else {
		message := connection.unreadMessages[0]
		connection.unreadMessages = connection.unreadMessages[1:]
		return message, nil
	}
	return "", nil
}

func (connection *Connection) receiveMessage() (message string) {
	websocket.Message.Receive(connection.RawConn, &message)
	return message
}

func (connection *Connection) WriteMessage(message string) {
	websocket.Message.Send(connection.RawConn, message)
}
