package websocket // import "taylz.io/http/websocket"

import (
	"sync"

	"golang.org/x/net/websocket"
)

// Conn = websocket.Conn
type Conn = websocket.Conn

// Upgrader = websocket.Handler
type Upgrader = websocket.Handler

// Send calls the websocket send API
func Send(conn *Conn, bytes []byte) error { return websocket.Message.Send(conn, bytes) }

// Receive calls the websocket receive API
func Receive(conn *Conn) (buf string, err error) {
	err = websocket.Message.Receive(conn, &buf)
	return
}

// T is a Websocket
type T struct {
	conn    *Conn
	id      string
	session string
	in      <-chan *Message
	out     chan []byte
	once    sync.Once
	done    chan bool
}

// New creates a websocket wrapper T
func New(conn *Conn, id string, sessionID string) *T {
	return &T{
		conn:    conn,
		id:      id,
		session: sessionID,
		in:      newChanMessage(conn),
		out:     make(chan []byte),
		done:    make(chan bool),
	}
}

// ID returns the websocket ID
func (ws *T) ID() string { return ws.id }

// SessionID returns the associated SessionID, if available
func (ws *T) SessionID() string { return ws.session }

// Done returns the done channel
func (ws *T) Done() <-chan bool { return ws.done }

// Write starts a goroutine to call WriteSync
func (ws *T) Write(bytes []byte) { go ws.WriteSync(bytes) }

// WriteSync waits to put a buffer into send queue
func (ws *T) WriteSync(bytes []byte) { ws.out <- bytes }

// Close closes the done channel, returns success
func (ws *T) Close() (ok bool) {
	ws.once.Do(func() {
		ok = true
		close(ws.done)
		close(ws.out)
	})
	return
}

// watch starts goroutines for in and out, awaits <-done
func (ws *T) watch(handler Handler) {
	go ws.watchin(handler)
	go ws.watchout()
	<-ws.done
}

// watchout loops over out until done
func (ws *T) watchout() {
	for {
		select {
		case <-ws.done:
			go DrainChanBytes(ws.out)
			return
		case buff, ok := <-ws.out:
			if !ok {
				ws.Close()
				return
			}
			if err := Send(ws.conn, buff); err != nil {
				go DrainChanBytes(ws.out)
				ws.Close()
				return
			}
		}
	}
}

// watchin loops over in until done
func (ws *T) watchin(handler Handler) {
	for {
		select {
		case <-ws.done:
			go drainChanMessage(ws.in)
			return
		case msg, ok := <-ws.in:
			if !ok {
				ws.Close()
				return
			}

			go handler.ServeWS(ws, msg)
		}
	}
}
