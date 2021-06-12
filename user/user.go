package user

import "taylz.io/http/websocket"

// T is a user, bridges session and websocket
type T struct {
	name    string
	sockets *Sockets
}

// New creates a user
func New(name string) *T {
	return &T{
		name:    name,
		sockets: NewSockets(),
	}
}

// Name returns the name given during creation
func (t *T) Name() string { return t.name }

// Sockets returns the socket ids linked with the user
func (t *T) Sockets() []string { return t.sockets.Keys() }

// AddSocket adds a socket id to the user
func (t *T) AddSocket(ws *websocket.T) { t.sockets.Set(ws.ID(), ws) }

// RemoveSocket removes a socket id from the user
func (t *T) RemoveSocket(id string) { t.sockets.Remove(id) }

// Message starts a goroutine to call WriteSync with websocket.Mesage.EncodeToJson
func (t *T) Message(uri string, data map[string]interface{}) {
	go t.WriteSync(websocket.Message{URI: uri, Data: data}.EncodeToJSON())
}

// Write starts a go routine to call WriteSync
func (t *T) Write(bytes []byte) { go t.WriteSync(bytes) }

// WriteSync writes the buffer to all sockets
func (t *T) WriteSync(bytes []byte) {
	t.sockets.Each(func(id string, ws *websocket.T) {
		ws.Write(bytes)
	})
}
