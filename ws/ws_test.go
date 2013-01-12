package ws

import (
	"code.google.com/p/go.net/websocket"
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"net/http"
	"testing"
)

func websocketHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/websocket", websocket.Handler(func(ws *websocket.Conn) {
		var name string
		websocket.Message.Receive(ws, &name)
		websocket.Message.Send(ws, "Hello, "+name)
	}))

	return mux
}

func TestWebSocket(t *testing.T) {
	testflight.WithServer(websocketHandler(), func(r *testflight.Requester) {
		connection := Connect(r, "/websocket")

		connection.WriteMessage("Drew")
		assert.Equal(t, "Hello, Drew", connection.ReceiveMessage())
	})
}
