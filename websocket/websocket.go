package websocket

import "golang.org/x/net/websocket"

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
	conn *Conn
	id   string
	name string
	in   <-chan *Message
	out  chan []byte
	done chan bool
}

// New creates a websocket wrapper T
func New(conn *Conn, id string, name string) *T {
	return &T{
		conn: conn,
		id:   id,
		name: name,
		in:   newChanMessage(conn),
		out:  make(chan []byte),
		done: make(chan bool),
	}
}

// ID returns the websocket ID
func (ws *T) ID() string { return ws.id }

// Name returns the name
func (ws *T) Name() string { return ws.name }

// SetName changes the name
func (ws *T) SetName(name string) { ws.name = name }

// Message calls Write using Message.JSON data format
func (ws *T) Message(uri string, data MsgData) {
	ws.Write(Message{URI: uri, Data: data}.EncodeToJSON())
}

// Write starts a goroutine to call WriteSync
func (ws *T) Write(bytes []byte) { go ws.WriteSync(bytes) }

// WriteSync waits to put a buffer into send queue
func (ws *T) WriteSync(bytes []byte) { ws.out <- bytes }

// Close closes the observable channel
func (ws *T) Close() {
	if ws.done != nil {
		close(ws.done)
		ws.done = nil
	}
}

func (ws *T) watchout() {
	for {
		select {
		case <-ws.done:
			close(ws.out)
			return
		case buff, ok := <-ws.out: // write to client
			if !ok {
				ws.Close()
				return
			}
			if err := Send(ws.conn, buff); err != nil {
				ws.Close()
				return
			}
		}
	}
}

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

// Writer is an API for *websocket.T
type Writer interface {
	Message(uri string, data MsgData)
	Write(bytes []byte)
}
