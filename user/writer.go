package user

import "taylz.io/http/websocket"

// Writer is an interface for writing data to users
type Writer interface {
	// Name returns the unique non-empty username of the client
	Name() string
	websocket.Writer
}

// WriterFunc emulates Writer, but returning any error closes the Writer
func WriterFunc(name string, f func([]byte) error) Writer {
	return writerFunc{
		name: name,
		ws:   websocket.WriterFunc(f),
	}
}

type writerFunc struct {
	name string
	ws   websocket.Writer
}

// Name returns the given name
func (w writerFunc) Name() string { return w.name }

// Done returns the done channel
func (w writerFunc) Done() <-chan bool { return w.ws.Done() }

// Write calls go WriteSync
func (w writerFunc) Write(bytes []byte) { w.ws.Write(bytes) }

// WriteSync calls the func
func (w writerFunc) WriteSync(bytes []byte) { w.ws.WriteSync(bytes) }
