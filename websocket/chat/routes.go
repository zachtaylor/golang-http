package chat

import (
	"time"

	"golang.org/x/time/rate"
	"taylz.io/http/websocket"
	"taylz.io/http/websocket/message"
)

// ChatProtocol returns a websocket.Protocol using message.NewSubprotocol
func ChatProtocol(server Server) websocket.Protocol {
	return websocket.ProtocolFunc(message.NewSubprotocol(rate.Every(time.Second), ChatHandler(server)))
}

func ChatHandler(server Server) message.Handler {
	return message.HandlerFunc(func(ws message.Writer, msg *message.T) {

	})
}
