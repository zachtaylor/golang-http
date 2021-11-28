package websocket

import (
	"net/http"

	"nhooyr.io/websocket"
)

// Manager is a websocket manager
type Manager struct {
	settings Settings
	cache    *Cache
	accept   *websocket.AcceptOptions
}

// NewManager creates a websocket manager
func NewManager(settings Settings) *Manager {
	return &Manager{
		settings: settings,
		cache:    NewCache(),
		accept: &websocket.AcceptOptions{
			Subprotocols:         settings.Protocol.GetSubprotocols(),
			InsecureSkipVerify:   settings.InsecureSkipVerify,
			OriginPatterns:       settings.OriginPatterns,
			CompressionMode:      settings.CompressionMode,
			CompressionThreshold: settings.CompressionThreshold,
		},
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

// ServeHTTP implements http.Handler by attempting websocket upgrade, handled by known subprotocol
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if conn, f := m.Accept(w, r); conn == nil {
		// skip
	} else if err := f(m.new(r, conn)); err == nil {
		conn.Close(StatusNormalClosure, "done")
	} else if err == ErrTooFast {
		conn.Close(StatusPolicyViolation, err.Error())
	} else if err == ErrUnsupportedDataType {
		conn.Close(StatusInvalidFramePayloadData, err.Error())
	} else {
		conn.Close(StatusProtocolError, err.Error())
	}
}

// Accept creates a connection
func (m *Manager) Accept(w http.ResponseWriter, r *http.Request) (*Conn, Subprotocol) {
	if conn, err := websocket.Accept(w, r, m.accept); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else if f := m.settings.Protocol.GetSubprotocol(conn.Subprotocol()); f == nil {
		conn.Close(StatusNormalClosure, `unknown subprotocol: "`+conn.Subprotocol()+`"`)
	} else {
		return conn, f
	}
	return nil, nil
}

// New saves a websocket
func (m *Manager) New(conn *Conn, sessionID string) (ws *T) {
	var id string
	// var  sessionID string
	// if s, _ := session.FromContext(r.Context()); s != nil {
	// 	sessionID = s.ID()
	// }
	m.cache.mu.Lock()
	for {
		id = m.settings.Keygen()
		if m.Get(id) == nil {
			break
		}
	}
	ws = New(r.Context(), conn, id, sessionID)
	m.cache.set(id, ws)
	m.cache.mu.Unlock()
	return
}
