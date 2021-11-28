package user

// import (
// 	"context"

// 	"taylz.io/http/websocket/message"
// )

// // Writer is an interface for writing data to users
// type Writer interface {
// 	// Name returns the unique non-empty username of the client
// 	Name() string
// 	message.Writer
// }

// // WriterFunc emulates Writer, but returning any error closes the Writer
// func WriterFunc(name string, f func([]byte) error) Writer {
// 	return writerFunc{
// 		name: name,
// 		ws:   message.WriterFunc(f),
// 	}
// }

// type writerFunc struct {
// 	name string
// 	ws   func([]byte) error
// }

// // Name returns the given name
// func (w writerFunc) Name() string { return w.name }

// // Done returns the done channel
// func (w writerFunc) Context() context.Context { return w.ws.Context() }

// func (w writerFunc) Subprotocol() string { return w.ws.Subprotocol() }

// // WriteMessagez calls go WriteSync
// func (w writerFunc) WriteMessage(msg *message.T) { w.ws.WriteMessage(msg) }

// // WriteSync calls the func
// func (w writerFunc) WriteMessageBytes(bytes []byte) { w.ws.WriteMessageBytes(bytes) }
