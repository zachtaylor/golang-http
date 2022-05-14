package user // import "taylz.io/http/user"

import (
	"sync"

	"taylz.io/http/session"
	"taylz.io/http/websocket"
	"taylz.io/yas"
)

// T is a user, bridges session and websocket
type T struct {
	session *session.T
	ws      yas.Set[*websocket.T]
	sync    sync.Mutex
	done    chan struct{}
	expired bool
}

// New creates a user
func New(session *session.T) (user *T) {
	return &T{
		session: session,
		ws:      yas.NewSet[*websocket.T](),
		done:    make(chan struct{}),
	}
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
	t.sync.Lock()
	sockets = t.ws.Slice()
	t.sync.Unlock()
	return
}

// AddSocket adds a socket id to the user
func (t *T) AddSocket(ws *websocket.T) {
	if t == nil || t.expired {
		return
	}
	t.sync.Lock()
	t.ws.Add(ws)
	t.sync.Unlock()
}

// RemoveSocket removes a socket id from the user
func (t *T) RemoveSocket(ws *websocket.T) {
	if t == nil || t.expired {
		return
	}
	t.sync.Lock()
	t.ws.Remove(ws)
	t.sync.Unlock()
}

func (t *T) WriteMessage(msg websocket.Message) error {
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
	for _, ws := range remove {
		t.ws.Remove(ws)
	}
	return
}
