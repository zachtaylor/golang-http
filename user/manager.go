package user

import (
	"net/http"

	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

// Manager is a user manager
type Manager struct {
	settings Settings
	cache    *Cache
}

// NewManager creates a session server
func NewManager(settings Settings) (manager *Manager) {
	manager = &Manager{settings: settings, cache: NewCache()}
	settings.Sessions.Observe(manager.onSession)
	settings.Sockets.Observe(manager.onWebsocket)
	return
}

// onSession is in the session cache hot path
func (m *Manager) onSession(id string, oldSession, newSession *session.T) {
	if newSession == nil && oldSession != nil {
		if name := oldSession.Name(); len(name) < 1 {
		} else if user := m.cache.Get(name); user != nil {
			go m.onSessionRemoveUser(user)
		}
	} else if newSession != nil && oldSession == nil {
		if name := newSession.Name(); len(name) > 0 {
			go m.onSessionAddUser(name)
		}
	}
}
func (m *Manager) onSessionAddUser(name string) { m.cache.Set(name, New(name)) }
func (m *Manager) onSessionRemoveUser(user *T)  { m.cache.Set(user.close(m.settings.Sockets)) }

// onWebsocket is in the websocket cache hot path
func (m *Manager) onWebsocket(id string, oldSocket, newSocket *websocket.T) {
	if newSocket == nil && oldSocket != nil {
		if name := oldSocket.Name(); len(name) < 1 {
		} else if user := m.cache.Get(name); user != nil {
			go user.RemoveSocket(id)
		}
	} else if newSocket != nil && oldSocket == nil {
		if name := newSocket.Name(); len(name) < 1 {
		} else if user := m.cache.Get(name); user != nil {
			go user.AddSocket(newSocket)
		}
	}
}

// Observe adds a callback CacheObserver
func (m *Manager) Observe(f CacheObserver) { m.cache.Observe(f) }

// Must returns the Session and User from session.Manager.Must
func (m *Manager) Must(name string) (*session.T, *T) {
	return m.settings.Sessions.Must(name), m.cache.Get(name)
}

// GetName returns the Session and User for the given name
func (m *Manager) GetName(name string) (session *session.T, user *T, err error) {
	if session = m.settings.Sessions.GetName(name); session == nil {
		err = ErrNotFound
	} else if user = m.cache.Get(session.Name()); user == nil {
		session, err = nil, ErrSessionSync
	}
	return
}

// GetRequestCookie returns the Session and User
func (m *Manager) GetRequestCookie(r *http.Request) (session *session.T, user *T, err error) {
	if session, err = m.settings.Sessions.GetRequestCookie(r); err != nil {
	} else if user = m.cache.Get(session.Name()); user == nil {
		session, err = nil, ErrSessionSync
	}
	return
}

// Authorize links an unnamed websocket to a user using m.Must
func (m *Manager) Authorize(name string, ws *websocket.T) (session *session.T, user *T, err error) {
	if session, user = m.Must(name); user == nil {
		err = ErrSessionSync
	} else if m.settings.Sockets.Rename(ws, name) {
		user.AddSocket(ws)
	}
	return
}
