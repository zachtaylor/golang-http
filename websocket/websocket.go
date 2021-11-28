package websocket // import "taylz.io/http/websocket"

import (
	"context"
	"io"
)

// T is a Websocket
type T struct {
	ctx     context.Context
	conn    *Conn
	id      string
	session string
}

// New creates a websocket wrapper T
func New(ctx context.Context, conn *Conn, id, session string) *T {
	return &T{
		ctx:     ctx,
		conn:    conn,
		id:      id,
		session: session,
	}
}

// ID returns the websocket ID
func (ws *T) ID() string { return ws.id }

// Context returns the original request Context
func (ws *T) Context() context.Context { return ws.ctx }

// Subprotocol returns the name of the negotiated subprotocol
func (ws *T) Subprotocol() string { return ws.conn.Subprotocol() }

// SessionID returns the associated SessionID, if available
func (ws *T) SessionID() string { return ws.session }

// CloseRead exposes the websocket Reader API and must only be called synchronously
//
// This is intended for use by a Protocol
func (ws *T) CloseRead() { ws.conn.CloseRead(ws.ctx) }

// Reader exposes the websocket Reader API and must only be called synchronously
//
// This is intended for use by a Protocol
func (ws *T) Reader() (MessageType, io.Reader, error) { return ws.conn.Reader(ws.ctx) }

// Writer exposes the websocket API and may be called asynchronously
//
// This is intended for use by a Protocol
func (ws *T) Writer(typ MessageType) (io.WriteCloser, error) { return ws.conn.Writer(ws.ctx, typ) }

// Close closes the connection, returns error
func (ws *T) Close(code StatusCode, reason string) error { return ws.conn.Close(code, reason) }
