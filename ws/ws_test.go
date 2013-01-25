package ws

import (
	"code.google.com/p/go.net/websocket"
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"net/http"
	"testing"
	"time"
)

func multiresponsewebsocketHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/websocket", websocket.Handler(func(ws *websocket.Conn) {
		for i := 0; i < 2; i++ {
			var name string
			websocket.Message.Receive(ws, &name)
			websocket.Message.Send(ws, "Hello, "+name)
		}
	}))

	return mux
}

func websocketHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/websocket", websocket.Handler(func(ws *websocket.Conn) {
		var name string
		websocket.Message.Receive(ws, &name)
		websocket.Message.Send(ws, "Hello, "+name)
	}))

	return mux
}


func donothingwebsocketHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/websocket", websocket.Handler(func(ws *websocket.Conn) {
		var name string
		websocket.Message.Receive(ws, &name)
	}))

	return mux
}

func TestWebSocket(t *testing.T) {
	testflight.WithServer(websocketHandler(), func(r *testflight.Requester) {
		connection := Connect(r, "/websocket")

		connection.WriteMessage("Drew")
		message, _ := connection.ReceiveMessage()
		assert.Equal(t, "Hello, Drew", message)
	})
}

func TestWebSocketReceiveMessageTimesOut(t *testing.T) {

	testflight.WithServer(donothingwebsocketHandler(), func(r *testflight.Requester) {
		connection := Connect(r, "/websocket")

		connection.WriteMessage("Drew")
		_, err := connection.ReceiveMessage()
		assert.Equal(t, TimeoutError{}, *err)
	})
}

func TestWebSocketTimeoutIsConfigurable(t *testing.T) {
	testflight.WithServer(websocketHandler(), func(r *testflight.Requester) {
		connection := Connect(r, "/websocket")
		connection.Timeout = 2 * time.Second

		go func () {
			time.Sleep(1 * time.Second)
			connection.WriteMessage("Drew")
		}()

		message, _ := connection.ReceiveMessage()
		assert.Equal(t, "Hello, Drew", message)
	})
}

func TestWebSocketRecordsReceivedMessages(t *testing.T) {
	testflight.WithServer(websocketHandler(), func(r *testflight.Requester) {
		connection := Connect(r, "/websocket")

		connection.WriteMessage("Drew")
		connection.ReceiveMessage()
		assert.Equal(t, "Hello, Drew", connection.ReceivedMessages[0])
	})
}

func TestWebSocketFlushesMessages(t *testing.T){
	testflight.WithServer(multiresponsewebsocketHandler(), func(r *testflight.Requester) {
		connection := Connect(r, "/websocket")

		connection.WriteMessage("Drew")
		connection.WriteMessage("Bob")
		connection.FlushMessages(2)
		assert.Equal(t, 2, len(connection.ReceivedMessages))
	})
}

func TestWebSocketTimesOutWhileFlushingMessages(t *testing.T){
	testflight.WithServer(donothingwebsocketHandler(), func(r *testflight.Requester) {
		connection := Connect(r, "/websocket")
		connection.WriteMessage("Drew")

		err := connection.FlushMessages(2)
		assert.Equal(t, TimeoutError{}, *err)
	})
}
