package user

import (
	"net/http"

	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

// Manager is a user manager
type Manager struct {
	sessions *session.Manager
	sockets  *websocket.Manager
	cache    *Cache
}

// NewManager creates a user manager
func NewManager(sessions *session.Manager, sockets *websocket.Manager) (manager *Manager) {
	manager = &Manager{
		sessions: sessions,
		sockets:  sockets,
		cache:    NewCache(),
	}
	sessions.Observe(manager.onSession)
	sockets.Observe(manager.onWebsocket)
	return
}

// onSession is in the session cache hot path
func (m *Manager) onSession(id string, oldSession, newSession *session.T) {
	if newSession == nil && oldSession != nil {
		if user := m.cache.Get(oldSession.Name()); user != nil {
			go m.onSessionRemoveUser(user)
		}
	} else if newSession != nil && oldSession == nil {
		go m.onSessionAddUser(newSession.Name(), id)
	}
}
func (m *Manager) onSessionAddUser(name, session string) {
	m.cache.Set(name, New(name, session))
}
func (m *Manager) onSessionRemoveUser(user *T) {
	m.sockets.RemoveSessionWebsockets(user.Sockets())
	m.cache.Set(user.name, nil)
	// user.close()
}

// onWebsocket is in the websocket cache hot path
func (m *Manager) onWebsocket(id string, oldSocket, newSocket *websocket.T) {
	if newSocket == nil && oldSocket != nil {
		if user, _, _ := m.GetSocket(oldSocket); user != nil {
			go user.RemoveSocket(id)
		}
	} else if newSocket != nil && oldSocket == nil {
		if user, _, _ := m.GetSocket(newSocket); user != nil {
			go user.AddSocket(newSocket)
		}
	}
}

// Must returns the User and Session from session.Manager.Must
func (m *Manager) Must(name string) (user *T, session *session.T) {
	session = m.sessions.Must(name)
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
		err = ErrExpired
	} else if session = m.sessions.Get(user.session); session == nil {
		user, err = nil, ErrSessionSync
	}
	return
}

// GetSession returns the User and Session for the given sessionID
func (m *Manager) GetSession(id string) (user *T, session *session.T, err error) {
	if session = m.sessions.Get(id); session == nil {
		err = ErrExpired
	} else if user = m.cache.Get(session.Name()); user == nil {
		session, err = nil, ErrSessionSync
	}
	return
}

// GetSocket returns the User and Session for the given websocket
func (m *Manager) GetSocket(ws *websocket.T) (*T, *session.T, error) {
	if sessionID := ws.SessionID(); len(sessionID) < 1 {
		return nil, nil, ErrNoID
	} else {
		return m.GetSession(sessionID)
	}
}

// GetRequestCookie returns the User and Session
func (m *Manager) GetRequestCookie(r *http.Request) (user *T, session *session.T, err error) {
	if session, err = m.sessions.GetRequestCookie(r); err != nil {
	} else if user = m.cache.Get(session.Name()); user == nil {
		session, err = nil, ErrSessionSync
	}
	return
}

// Authorize links a websocket to a [new] user using m.Must
func (m *Manager) Authorize(name string, ws *websocket.T) (user *T, session *session.T, err error) {
	if oldUser, _, _ := m.GetSocket(ws); oldUser != nil {
		m.sockets.SetSessionID(ws, "")
		oldUser.RemoveSocket(ws.ID())
	}
	if user, session = m.Must(name); m.sockets.SetSessionID(ws, session.ID()) {
		user.AddSocket(ws)
	}
	return
}
