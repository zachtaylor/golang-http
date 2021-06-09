package user

import "taylz.io/http/websocket"

// T is a user, bridges session and websocket
type T struct {
	man     *websocket.Manager
	name    string
	sockets *Sockets
}

// New creates a user
func New(man *websocket.Manager, name string) *T {
	return &T{
		man:     man,
		name:    name,
		sockets: NewSockets(),
	}
}

// Name returns the name given during creation
func (t *T) Name() string { return t.name }

// Sockets returns the socket ids linked with the user
func (t *T) Sockets() []string { return t.sockets.Keys() }

// AddSocket adds a socket id to the user
func (t *T) AddSocket(id string) { t.sockets.Set(id, true) }

// RemoveSocket removes a socket id from the user
func (t *T) RemoveSocket(id string) { t.sockets.Remove(id) }

// Message calls Write using websocket.Transport data format
func (t *T) Message(uri string, data map[string]interface{}) {
	t.Write(websocket.Message{URI: uri, Data: data}.EncodeToJSON())
}

// Write writes the buffer to all sockets
func (t *T) Write(bytes []byte) {
	for _, k := range t.Sockets() {
		if ws := t.man.Get(k); ws != nil {
			ws.Write(bytes)
		}
	}
}
