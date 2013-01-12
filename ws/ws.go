package ws

import (
	"code.google.com/p/go.net/websocket"
	"github.com/drewolson/testflight"
)

func Connect(r *testflight.Requester, route string) *Connection {
	connection, _ := websocket.Dial(websocketRoute(r, route), "", "http://localhost/")
	return newConnection(connection)
}

func websocketRoute(r *testflight.Requester, route string) string {
	return "ws://" + r.Server.Listener.Addr().String() + route
}
