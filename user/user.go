package user // import "taylz.io/http/user"

import (
	"sync"

	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

// T is a user, bridges session and websocket
type T struct {
	session *session.T
	wsocks  map[*websocket.T]struct{}
	wsmu    sync.Mutex
	done    chan struct{}
	expired bool
}

// New creates a user
func New(session *session.T) (user *T) {
	user = &T{
		session: session,
		wsocks:  make(map[*websocket.T]struct{}),
		done:    make(chan struct{}),
	}
	go user.watch()
	return
}

// Name returns the session name
func (t *T) Name() string { return t.session.Name() }

// Done returns the done channel for user
func (t *T) Done() <-chan struct{} {
	if t == nil || t.expired {
		return nil
	}
	return t.done
}

// Sockets returns the websockets linked with the user
func (t *T) Sockets() (sockets []*websocket.T) {
	if t == nil || t.expired {
		return
	}
	t.wsmu.Lock()
	sockets = t.asyncSockets()
	t.wsmu.Unlock()
	return
}

func (t *T) asyncSockets() []*websocket.T {
	i, sockets := 0, make([]*websocket.T, len(t.wsocks))
	for ws := range t.wsocks {
		sockets[i] = ws
		i++
	}
	return sockets
}

var wsFound = struct{}{}

// AddSocket adds a socket id to the user
func (t *T) AddSocket(ws *websocket.T) {
	if t == nil || t.expired {
		return
	}
	t.wsmu.Lock()
	t.wsocks[ws] = wsFound
	t.wsmu.Unlock()
}

// RemoveSocket removes a socket id from the user
func (t *T) RemoveSocket(ws *websocket.T) {
	if t == nil || t.expired {
		return
	}
	t.wsmu.Lock()
	delete(t.wsocks, ws)
	t.wsmu.Unlock()
}

func (t *T) WriteMessage(msg *websocket.Message) error {
	return t.Write(websocket.MessageText, msg.ShouldMarshal())
}

func (t *T) Write(typ websocket.MessageType, data []byte) (err error) {
	if t == nil || t.expired {
		return session.ErrExpired
	}
	remove := []*websocket.T{}
	err = ErrMissingConn
	for _, ws := range t.Sockets() {
		if ws.Write(typ, data) != nil { // erases err
			remove = append(remove, ws)
		} else {
			err = nil
		}
	}
	t.wsmu.Lock()
	for _, ws := range remove {
		delete(t.wsocks, ws)
	}
	t.wsmu.Unlock()
	return
}

func (t *T) watch() {
	<-t.done
	t.wsmu.Lock()
	t.expired = true
	t.wsocks = nil
	t.wsmu.Unlock()
}
