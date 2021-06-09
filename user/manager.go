package user

import (
	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

// Manager is a user manager
type Manager struct {
	Settings
	*Cache
}

// NewManager creates a session server
func NewManager(settings Settings, cache *Cache) (manager *Manager) {
	manager = &Manager{Settings: settings, Cache: cache}
	settings.Sessions.Observe(manager.onSession)
	settings.Sockets.Observe(manager.onWebsocket)
	return
}

func (m *Manager) onSession(id string, oldSession, newSession *session.T) {
	if newSession == nil {
		m.Set(oldSession.Name(), nil)
	} else {
		m.Set(newSession.Name(), New(m.Sockets))
	}
}

func (m *Manager) onWebsocket(id string, oldSocket, newSocket *websocket.T) {
	if newSocket == nil {
		if name := newSocket.Name(); len(name) < 1 {
		} else if user := m.Get(name); user != nil {
			user.AddSocket(id)
		}
	} else if oldSocket != nil {
		if name := oldSocket.Name(); len(name) < 1 {
		} else if user := m.Get(name); user != nil {
			user.RemoveSocket(id)
		}
	}
}

// Authorize links a websocket to a session
func (m *Manager) Authorize(session *session.T, ws *websocket.T) {
	if user := m.Get(session.Name()); user != nil {
		user.AddSocket(ws.ID())
	}
}
