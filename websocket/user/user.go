package user // import "taylz.io/http/user"

import (
	"errors"
	"sync"

	"taylz.io/http"
	"taylz.io/http/session"
	"taylz.io/http/websocket"
	"taylz.io/maps"
)

var (
	// ErrSessionSync indicates a caching issue with session.Manager
	ErrSessionSync = errors.New("user: session out of sync")

	// ErrMissingConn indicates a write has no available websockets
	ErrMissingConn = errors.New("user: missing conn instance")
)

// Manager is a user manager
type Manager interface {
	Reader
	Writer
	// Size returns the current len of the map
	Size() int
	// Get returns a user by name
	Get(string) *T
	// Must (re)authorizes a session for a websocket
	Must(*websocket.T, string) *T
	// GetWebsocket returns a user by websocket
	GetWebsocket(*websocket.T) *T
	// Observe adds a callback CacheObserver
	Observe(Observer)
}

type Observer = maps.Observer[string, *T]

type ObserverFunc = maps.ObserverFunc[string, *T]

// Reader is an interface for recognizing Users in http.Request
type Reader interface {
	// ReadHTTP returns the User and Session
	ReadHTTP(*http.Request) (*T, error)
}

// Writer is an interface for writing Users to http.Writer
type Writer interface {
	// WriteHTTP writes the Set-Cookie header using the session.Manager
	WriteHTTP(http.Writer, *T)
}

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
