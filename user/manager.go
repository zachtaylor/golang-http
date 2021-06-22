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
func NewManager(settings Settings) (manager *Manager) {
	manager = &Manager{Settings: settings, Cache: NewCache()}
	settings.Sessions.Observe(manager.onSession)
	settings.Sockets.Observe(manager.onWebsocket)
	return
}

func (m *Manager) onSession(id string, oldSession, newSession *session.T) {
	if newSession == nil && oldSession != nil {
		if name := oldSession.Name(); len(name) > 0 {
			go m.Set(name, nil)
		}
	} else if newSession != nil && oldSession == nil {
		if name := oldSession.Name(); len(name) > 0 {
			go m.Set(name, New(name))
		}
	}
}

func (m *Manager) onWebsocket(id string, oldSocket, newSocket *websocket.T) {
	if newSocket == nil && oldSocket != nil {
		if name := oldSocket.Name(); len(name) < 1 {
		} else if user := m.Get(name); user != nil {
			go user.RemoveSocket(id)
		}
	} else if newSocket != nil && oldSocket == nil {
		if name := newSocket.Name(); len(name) < 1 {
		} else if user := m.Get(name); user != nil {
			go user.AddSocket(newSocket)
		}
	}
}

// Authorize links a websocket to a session
func (m *Manager) Authorize(session *session.T, ws *websocket.T) {
	if user := m.Get(session.Name()); user != nil {
		user.AddSocket(ws)
	}
}
