package session

import (
	"time"

	"taylz.io/http"
)

// Service implements Manager with a *Cache
type Service struct {
	Settings
	cache  *Cache
	keygen func() string
	expiry TestExpired
}

// NewService creates a Service
func NewService(settings Settings, keygen func() string) *Service {
	service := &Service{
		Settings: settings,
		cache:    NewCache(),
		keygen:   keygen,
	}
	time.AfterFunc(service.Lifetime, service.collectgarbage)
	return service
}

// Must refreshes and returns the Session with the given username if one exists, and creates one if necessary
func (s *Service) Must(name string) (session *T) {
	if _, session = s.cache.TestFirst(TestName(name)); session != nil {
		session.time = time.Now().Add(s.Lifetime)
	} else {
		var id string
		for ok := true; ok; ok = (s.Get(id) != nil) {
			id = s.keygen()
		}
		session = New(id, name, time.Now().Add(s.Lifetime))
		s.cache.Set(id, session)
	}
	return
}

// Get returns a Session by ID
func (s *Service) Get(id string) *T { return s.cache.Get(id) }

// Count returns the current len of the map
func (s *Service) Count() int { return s.cache.Count() }

// Update changes the internal expiry time of a Session
func (s *Service) Update(sessionID string) (err error) {
	s.cache.Sync(func() {
		if session := s.Get(sessionID); session == nil {
			err = ErrExpired
		} else {
			session.time = time.Now().Add(s.Lifetime)
		}
	})
	return
}

// Remove removes a Session
func (s *Service) Remove(id string) { s.cache.Remove(id) }

// Observe adds an Observer
func (s *Service) Observe(f Observer) { s.cache.Observe(f) }

// ReadHTTP implements Reader
func (s *Service) ReadHTTP(r *http.Request) (session *T, err error) {
	if cookie, _err := r.Cookie(s.CookieID); _err != nil {
		err = ErrNoID
	} else if session = s.Get(cookie.Value); session == nil {
		err = ErrExpired
	}
	return
}

// WriterHTTP writes the Set-Cookie header
func (s *Service) WriterHTTP(w http.ResponseWriter, session *T) {
	if session == nil {
		w.Header().Set("Set-Cookie", s.CookieID+"=; Path=/; Expires==Thu, 01 Jan 1970 00:00:00 GMT;")
		return
	}
	header := s.CookieID + "=" + session.id + "; Path=/; "
	if s.Secure {
		header += "Secure; "
	}
	if s.Strict {
		header += "SameSite=Strict;"
	} else {
		header += "SameSite=Lax;"
	}
	w.Header().Set("Set-Cookie", header)
}

func (s *Service) collectgarbage() {
	s.expiry = TestExpired(time.Now())
	s.cache.RemoveTest(s.expiry)
	time.AfterFunc(s.GC, s.collectgarbage)
}
