package websocket

import "net/http"

// Manager is a websocket manager
type Manager struct {
	server   Server
	settings Settings
	cache    *Cache
}

// NewManager creates a websocket manager
func NewManager(server Server, settings Settings) *Manager {
	return &Manager{
		server:   server,
		settings: settings,
		cache:    NewCache(),
	}
}

// Get returns the websocket by id
func (m *Manager) Get(id string) *T { return m.cache.dat[id] }

// Count returns the current len of the map
func (m *Manager) Count() int { return len(m.cache.dat) }

// Observe adds a CacheObserver
func (m *Manager) Observe(f CacheObserver) { m.cache.Observe(f) }

// Each calls the func for each websocket (under lock)
func (m *Manager) Each(f func(string, *T)) { m.cache.Each(f) }

// NewUpgrader returns a new http.Handler that upgrades requests to add Websockets to Manager.Cache
func (m *Manager) NewUpgrader() http.Handler { return Upgrader(m.connect) }

// connect is called by the websocket api connection upgrader
func (m *Manager) connect(conn *Conn) {
	var sessionID string
	if session, _ := m.server.GetSessionManager().GetRequestCookie(conn.Request()); session != nil {
		sessionID = session.ID()
	}
	m.cache.mu.Lock()
	var id string
	for ok := true; ok; ok = m.Get(id) != nil {
		id = m.settings.Keygen()
	}
	ws := New(conn, id, sessionID)
	m.cache.set(id, ws)
	m.cache.mu.Unlock()
	ws.watch(m.server.GetWebsocketHandler())
	m.cache.Remove(ws.id)
}

// SetSessionID changes the internal SessionID of a managed websocket
func (m *Manager) SetSessionID(ws *T, sessionID string) (ok bool) {
	m.cache.mu.Lock()
	if ws == m.cache.dat[ws.id] {
		ok, ws.session = true, sessionID
	}
	m.cache.mu.Unlock()
	return
}

// RemoveSessionWebsockets wipes the internal SessionID of managed websockets by socketID
func (m *Manager) RemoveSessionWebsockets(ids []string) {
	m.cache.mu.Lock()
	for _, id := range ids {
		if ws := m.Get(id); ws != nil {
			ws.session = ""
		}
	}
	m.cache.mu.Unlock()
}
