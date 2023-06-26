package user // import "taylz.io/http/user"

import (
	"sync"

	"taylz.io/http/session"
	"taylz.io/http/websocket"
	"taylz.io/maps"
)

// T is a user, bridges session and websocket
type T struct {
	session *session.T
	ws      maps.Set[*websocket.T]
	sync    sync.Mutex
	done    chan struct{}
	expired bool
}

// New creates a user
func New(session *session.T) (user *T) {
	return &T{
		session: session,
		ws:      maps.NewSet[*websocket.T](),
		done:    make(chan struct{}),
	}
}

// Session returns the session id
func (t *T) Session() *session.T { return t.session }

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

func (t *T) WriteText(data []byte) error { return t.Write(websocket.MessageText, data) }

func (t *T) WriteBinary(data []byte) error { return t.Write(websocket.MessageBinary, data) }

func (t *T) Write(typ websocket.MessageType, data []byte) (err error) {
	if t == nil || t.expired {
		return session.ErrExpired
	}
	for ws := range t.ws {
		if ws.Write(typ, data) != nil {
			t.ws.Remove(ws)
		} else {
			err = nil
		}
	}
	return
}
