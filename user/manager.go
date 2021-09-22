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
			go m.onSessionAddUser(name, id)
		}
	}
}
func (m *Manager) onSessionAddUser(name, session string) {
	m.cache.Set(name, New(name, session))
}
func (m *Manager) onSessionRemoveUser(user *T) {
	m.settings.Sockets.RemoveSessionWebsockets(user.Sockets())
	m.cache.Set(user.name, nil)
	user.close()
}

// onWebsocket is in the websocket cache hot path
func (m *Manager) onWebsocket(id string, oldSocket, newSocket *websocket.T) {
	if newSocket == nil && oldSocket != nil {
		if sessionID := oldSocket.SessionID(); len(sessionID) < 1 {
		} else if session := m.settings.Sessions.Get(sessionID); session == nil {
		} else if user := m.cache.Get(session.Name()); user != nil {
			go user.RemoveSocket(id)
		}
	} else if newSocket != nil && oldSocket == nil {
		if sessionID := newSocket.SessionID(); len(sessionID) < 1 {
		} else if session := m.settings.Sessions.Get(sessionID); session == nil {
		} else if user := m.cache.Get(session.Name()); user != nil {
			go user.AddSocket(newSocket)
		}
	}
}

// Must returns the User and Session from session.Manager.Must
func (m *Manager) Must(name string) (user *T, session *session.T) {
	session = m.settings.Sessions.Must(name)
	user = m.cache.Get(name)
	return
}

// Count returns the current len of the map
func (m *Manager) Count() int { return len(m.cache.dat) }

// Observe adds a callback CacheObserver
func (m *Manager) Observe(f CacheObserver) { m.cache.Observe(f) }

// Get returns the User and Session for the given username
func (m *Manager) Get(name string) (user *T, session *session.T, err error) {
	if user = m.cache.Get(name); user == nil {
		err = ErrNotFound
	} else if session = m.settings.Sessions.Get(user.session); session == nil {
		user, err = nil, ErrExpired
	}
	return
}

// GetSession returns the User and Session for the given sessionID
func (m *Manager) GetSession(id string) (user *T, session *session.T, err error) {
	if session = m.settings.Sessions.Get(id); session == nil {
		err = ErrNotFound
	} else if user = m.cache.Get(session.Name()); user == nil {
		session, err = nil, ErrSessionSync
	}
	return
}

// GetRequestCookie returns the User and Session
func (m *Manager) GetRequestCookie(r *http.Request) (user *T, session *session.T, err error) {
	if session, err = m.settings.Sessions.GetRequestCookie(r); err != nil {
	} else if user = m.cache.Get(session.Name()); user == nil {
		session, err = nil, ErrSessionSync
	}
	return
}

// Authorize links a websocket to a [new] user using m.Must
func (m *Manager) Authorize(name string, ws *websocket.T) (user *T, session *session.T, err error) {
	if oldSessionID := ws.SessionID(); len(oldSessionID) < 1 {
	} else if oldSession := m.settings.Sessions.Get(oldSessionID); oldSession == nil {
	} else if oldUser := m.cache.Get(oldSession.Name()); oldUser != nil {
		m.settings.Sockets.SetSessionID(ws, "")
		oldUser.RemoveSocket(ws.ID())
	}
	if user, session = m.Must(name); m.settings.Sockets.SetSessionID(ws, session.ID()) {
		user.AddSocket(ws)
	}
	return
}
