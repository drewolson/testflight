package ws

import (
	"code.google.com/p/go.net/websocket"
	"time"
)

type Connection struct {
	RawConn *websocket.Conn
	ReceivedMessages []string
	Timeout time.Duration
	unreadMessages []string
}

func newConnection(conn *websocket.Conn) *Connection {
	connection := &Connection{
		RawConn: conn,
		Timeout: 1 * time.Second,
	}
	return connection
}

func (connection *Connection) ReceiveMessage() (string, *TimeoutError) {
	messageChan := make(chan string)

	go connection.receiveMessage(messageChan)

	select {
		case  <-time.After(connection.Timeout):
			return "", &TimeoutError{}
		case message := <-messageChan:
			connection.ReceivedMessages = append(connection.ReceivedMessages, message)
			return message, nil
	}

	return "", nil
}

func (connection *Connection) FlushMessages(number int) *TimeoutError {
	for i := 0; i < number; i++ {
		_, err := connection.ReceiveMessage()
		if err != nil {
			return err
		}
	}
	return nil
}

func (connection *Connection) receiveMessage(messageChan chan string) {
	var message string
	websocket.Message.Receive(connection.RawConn, &message)
	messageChan <- message
}

func (connection *Connection) WriteMessage(message string) {
	websocket.Message.Send(connection.RawConn, message)
}
