package websocket

import "net/http"

type Manager struct {
	Settings
	*Cache
}

// NewManager creates a session server
func NewManager(settings Settings, cache *Cache) *Manager {
	return &Manager{Settings: settings, Cache: cache}
}

// Upgrader returns a new http.Handler that adds upgrades request to add a Websocket to this cache
func (m *Manager) Upgrader() http.Handler {
	return Upgrader(m.connect)
}

func (m *Manager) connect(conn *Conn) {
	var name string
	if session := m.Sessions.RequestSessionCookie(conn.Request()); session != nil {
		name = session.Name()
	}
	m.mu.Lock()
	var id string
	for ok := true; ok; ok = m.Get(id) != nil {
		id = m.Keygen()
	}
	ws := New(conn, id, name)
	m.set(id, ws)
	m.mu.Unlock()
	m.Watch(ws)
	m.Remove(ws.id)
}

func (m *Manager) Watch(ws *T) {
	go ws.watchin(m.Handler)
	go ws.watchout()
	<-ws.done
}
