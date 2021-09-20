package websocket

import "net/http"

// Manager is a websocket manager
type Manager struct {
	settings Settings
	cache    *Cache
}

// NewManager creates a websocket manager
func NewManager(settings Settings, cache *Cache) *Manager {
	return &Manager{settings: settings, cache: NewCache()}
}

// Upgrader returns a new http.Handler that upgrades requests to add Websockets to Manager.Cache
func (m *Manager) Upgrader() http.Handler {
	return Upgrader(m.connect)
}

func (m *Manager) connect(conn *Conn) {
	var name string
	if session, _ := m.settings.Sessions.GetRequestCookie(conn.Request()); session != nil {
		name = session.Name()
	}
	m.cache.mu.Lock()
	var id string
	for ok := true; ok; ok = m.Get(id) != nil {
		id = m.settings.Keygen()
	}
	ws := New(conn, id, name)
	m.cache.set(id, ws)
	m.cache.mu.Unlock()
	ws.watch(m.settings.Handler)
	m.cache.Remove(ws.id)
}

// Rename changes the internal name of a managed websocket
func (m *Manager) Rename(ws *T, name string) (ok bool) {
	m.cache.mu.Lock()
	if (len(name) < 1) != (len(ws.name) < 1) && ws != m.cache.dat[ws.id] {
		ok, ws.name = true, name
	}
	m.cache.mu.Unlock()
	return
}

// Unname wipes the internal name of managed websockets
func (m *Manager) Unname(ids []string) {
	m.cache.mu.Lock()
	for _, id := range ids {
		if ws := m.Get(id); ws != nil {
			ws.name = ""
		}
	}
	m.cache.mu.Unlock()
}

// Get returns the websocket by id
func (m *Manager) Get(id string) *T { return m.cache.dat[id] }

// Observe adds a CacheObserver
func (m *Manager) Observe(f CacheObserver) { m.cache.Observe(f) }
