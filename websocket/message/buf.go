package message

import (
	"context"

	"golang.org/x/time/rate"
	"taylz.io/http/websocket"
)

// buf is buffer value type Writer for *websocket.T
type buf struct {
	ws          *websocket.T
	readLimiter *rate.Limiter
}

func newBuf(ws *websocket.T, limit rate.Limit) buf {
	return buf{ws: ws, readLimiter: rate.NewLimiter(limit, 3)}
}

func (b buf) ID() string { return b.ws.ID() }

func (b buf) Context() context.Context { return b.ws.Context() }

func (b buf) Subprotocol() string { return b.ws.Subprotocol() }

func (b buf) WriteMessage(msg *T) error { return writeMessage(b.ws, msg) }

func (b buf) WriteMessageBytes(buf []byte) error { return writeMessageBytes(b.ws, buf) }
