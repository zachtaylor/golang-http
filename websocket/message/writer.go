package message

import (
	"context"
)

// Writer is an interface for writing data to websockets
type Writer interface {
	// ID returns the websocket ID
	ID() string

	// Context returns the websocket connection context
	Context() context.Context

	// Subprotocol returns the negotiated connection subprotocol
	Subprotocol() string

	// WriteMessage encodes Message
	WriteMessage(*T) error

	// WriteMessageBytes writes an encoded Message
	WriteMessageBytes([]byte) error
}

// WriterFunc emulates Writer
func WriterFunc(subprocotol string, f func([]byte) error) Writer {
	return writerFunc{
		ctx:  context.Background(),
		subp: subprocotol,
		f:    f,
	}
}

type writerFunc struct {
	ctx  context.Context
	subp string
	f    func([]byte) error
}

func (w writerFunc) ID() string { return "" }

func (w writerFunc) Context() context.Context { return w.ctx }

func (w writerFunc) Subprotocol() string { return w.subp }

// WriteMessage calls WriteMessageBytes
func (w writerFunc) WriteMessage(msg *T) error { return w.WriteMessageBytes(msg.Marshal()) }

// WriteMessageBytes calls w.f
func (w writerFunc) WriteMessageBytes(bytes []byte) error { return w.f(bytes) }
