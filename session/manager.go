package session

import (
	"net/http"
	"time"
)

// Manager is a session manager
type Manager struct {
	settings Settings
	cache    *Cache
}

// NewManager creates a session server
func NewManager(settings Settings) (manager *Manager) {
	manager = &Manager{settings: settings, cache: NewCache()}
	time.AfterFunc(manager.settings.Lifetime, func() { manager.collectgarbage() })
	return
}

// ageLimit returns the age limit relative to time.Now()
func (m *Manager) ageLimit() time.Time { return time.Now().Add(-m.settings.Lifetime) }

// Must refreshes and returns the Session with the given username if one exists, and creates one if necessary
func (m *Manager) Must(name string) (session *T) {
	expiry := m.ageLimit()
	m.cache.mu.Lock()
	defer m.cache.mu.Unlock()
	if session = m.getName(name, expiry); session != nil {
		session.Update()
		return
	}
	var id string
	for ok := true; ok; ok = m.get(id, expiry) != nil {
		id = m.settings.Keygen()
	}
	session = New(id, name)
	m.cache.set(id, session)
	return
}

// Get returns a Session by ID
func (m *Manager) Get(id string) *T { return m.get(id, m.ageLimit()) }

// get checks expiry
func (m *Manager) get(id string, expiry time.Time) (session *T) {
	if session = m.cache.Get(id); session.time.Before(expiry) {
		session = nil
	}
	return
}

// Count returns the current len of the map
func (m *Manager) Count() int { return len(m.cache.dat) }

// Remove removes a Session
func (m *Manager) Remove(id string) { m.cache.Set(id, nil) }

// Observe adds a CacheObserver
func (m *Manager) Observe(f CacheObserver) { m.cache.Observe(f) }

// GetName returns Session by username
func (m *Manager) GetName(name string) (session *T) {
	if len(name) > 1 {
		expiry := m.ageLimit()
		m.cache.mu.Lock()
		session = m.getName(name, expiry)
		m.cache.mu.Unlock()
	}
	return
}

// getName iterates m.cache.dat without locking m.cache.mu and check expiry
func (m *Manager) getName(name string, expiry time.Time) (session *T) {
	for _, t := range m.cache.dat {
		if t.name == name && t.time.After(expiry) {
			session = t
			break
		}
	}
	return
}

// GetRequestCookie returns Session by Request.Cookie
func (m *Manager) GetRequestCookie(r *http.Request) (session *T, err error) {
	if cookie, _err := r.Cookie(m.settings.CookieID); _err != nil {
		err = ErrNoCookie
	} else if session = m.Get(cookie.Value); session == nil {
		err = ErrExpired
	}
	return
}

// WriteSetCookie writes the Set-Cookie header
func (m *Manager) WriteSetCookie(w http.ResponseWriter, session *T) {
	if session == nil {
		w.Header().Set("Set-Cookie", m.settings.CookieID+"=; Path=/; Expires==Thu, 01 Jan 1970 00:00:00 GMT;")
		return
	}
	header := m.settings.CookieID + "=" + session.id + "; Path=/; "
	if m.settings.Secure {
		header += "Secure; "
	}
	if m.settings.Strict {
		header += "SameSite=Strict;"
	} else {
		header += "SameSite=Lax;"
	}
	w.Header().Set("Set-Cookie", header)
}

func (m *Manager) collectgarbage() {
	expiry := m.ageLimit()
	list := make([]string, 0)
	m.cache.mu.Lock()
	for k, t := range m.cache.dat {
		if t.time.Before(expiry) {
			list = append(list, k)
		}
	}
	for _, key := range list {
		m.cache.set(key, nil)
	}
	m.cache.mu.Unlock()
	time.AfterFunc(m.settings.GC, func() { m.collectgarbage() })
}
