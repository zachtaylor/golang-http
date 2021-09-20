package user

import (
	"sync"

	"taylz.io/http/websocket"
)

// T is a user, bridges session and websocket
type T struct {
	name    string
	sockets *Sockets
	buff    chan []byte
	once    sync.Once
	done    chan bool
}

// New creates a user
func New(name string) (user *T) {
	user = &T{
		name:    name,
		done:    make(chan bool),
		buff:    make(chan []byte),
		sockets: NewSockets(),
	}
	go user.watchBuff()
	return
}

// Name returns the name given during creation
func (t *T) Name() string { return t.name }

// Done returns the done channel for user
func (t *T) Done() <-chan bool { return t.done }

// Sockets returns the socket ids linked with the user
func (t *T) Sockets() []string { return t.sockets.Keys() }

// AddSocket adds a socket id to the user
func (t *T) AddSocket(ws *websocket.T) { t.sockets.Set(ws.ID(), ws) }

// RemoveSocket removes a socket id from the user
func (t *T) RemoveSocket(id string) { t.sockets.Remove(id) }

// Write starts a go routine to call WriteSync
func (t *T) Write(bytes []byte) { go t.WriteSync(bytes) }

// WriteSync waits to put a buffer into send queue
func (t *T) WriteSync(bytes []byte) { t.buff <- bytes }

// writeSync writes the buffer to all sockets
func (t *T) writeSync(bytes []byte) {
	t.sockets.Each(func(id string, ws *websocket.T) {
		ws.Write(bytes)
	})
}

// close is a special hook for the manager
func (t *T) close(man *websocket.Manager) (string, *T) {
	t.once.Do(func() {
		man.Unname(t.sockets.Keys())
		close(t.done)
		close(t.buff)
	})
	return t.name, nil
}

// watchBuff loops over buff until done
func (t *T) watchBuff() {
	for {
		select {
		case <-t.done:
			websocket.DrainChanBytes(t.buff)
			return
		case buff, ok := <-t.buff:
			if !ok {
				return
			}
			t.writeSync(buff)
		}
	}
}
