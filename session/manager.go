package session

import (
	"net/http"
	"time"
)

// Manager is a session manager
type Manager struct {
	Settings
	*Cache
}

// NewManager creates a session server
func NewManager(settings Settings) (manager *Manager) {
	manager = &Manager{Settings: settings, Cache: NewCache()}
	time.AfterFunc(manager.Lifetime, func() { manager.collectgarbage() })
	return
}

// Grant returns a new Session granted to the username
func (m *Manager) Grant(name string) (session *T) {
	m.mu.Lock()
	var id string
	for ok := true; ok; ok = m.Get(id) != nil {
		id = m.Keygen()
	}
	session = New(time.Now(), id, name)
	m.set(id, session)
	m.mu.Unlock()
	return
}

// GetName returns the Session with the username
func (m *Manager) GetName(name string) (session *T) {
	m.mu.Lock()
	for _, t := range m.dat {
		if t.name == name {
			session = t
			break
		}
	}
	m.mu.Unlock()
	return
}

// RequestSessionCookie returns Session associated to the Request via Session cookie
func (m *Manager) RequestSessionCookie(r *http.Request) *T {
	cookie, err := r.Cookie(m.CookieID)
	if err != nil {
		return nil
	}
	return m.Get(cookie.Value)
}

// WriteSessionCookie writes the session cookie to the ResponseWriter
func (m *Manager) WriteSessionCookie(w http.ResponseWriter, session *T) {
	header := m.CookieID + "=" + session.id + "; Path=/; "
	if m.Secure {
		header += "Secure; "
	}
	if m.Strict {
		header += "SameSite=Strict;"
	} else {
		header += "SameSite=Lax;"
	}
	w.Header().Set("Set-Cookie", header)
}

// WriteSessionCookieExpired writes an expired session cookie to the ResponseWriter
func (s *Manager) WriteSessionCookieExpired(w http.ResponseWriter) {
	w.Header().Set("Set-Cookie", s.CookieID+"=; Path=/; Expires==Thu, 01 Jan 1970 00:00:00 GMT;")
}

func (m *Manager) collectgarbage() {
	list := make([]string, 0)
	expirey := time.Now().Add(-m.Lifetime)
	m.mu.Lock()
	for k, t := range m.dat {
		if t.time.Before(expirey) {
			list = append(list, k)
		}
	}
	for _, key := range list {
		m.set(key, nil)
	}
	m.mu.Unlock()
	time.AfterFunc(m.GC, func() { m.collectgarbage() })
}
