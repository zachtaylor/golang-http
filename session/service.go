package session

import "taylz.io/http"

// Service implements Manager with a *Cache
type Service struct {
	Settings
	cache *Cache
}

// NewService creates a Service
func NewService(settings Settings) *Service {
	return &Service{
		Settings: settings,
		cache:    NewCache(),
	}
}

// Must refreshes and returns the Session with the given username if one exists, and creates one if necessary
func (s *Service) Must(name string) (session *T) {
	now := now()
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	if session = getName(now, s.cache, name); session != nil {
		session.time = now.Add(s.Lifetime)
		return
	}
	var id string
	for ok := true; ok; ok = GetID(now, s.cache, id) != nil {
		id = s.Keygen()
	}
	session = New(id, name, now.Add(s.Lifetime))
	s.cache.set(id, session)
	return
}

// Get returns a Session by ID
func (s *Service) Get(id string) *T { return GetID(now(), s.cache, id) }

// Count returns the current len of the map
func (s *Service) Count() int { return len(s.cache.dat) }

// Update changes the internal expiry time of a Session
func (s *Service) Update(sessionID string) error {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	session := s.cache.dat[sessionID]
	if session == nil {
		return ErrExpired
	}
	session.time = now().Add(s.Lifetime)
	return nil
}

// Remove removes a Session
func (s *Service) Remove(id string) { s.cache.Set(id, nil) }

// Observe adds a CacheObserver
func (s *Service) Observe(f CacheObserver) { s.cache.Observe(f) }

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
	expiry, list := now(), make([]string, 0)
	s.cache.mu.Lock()
	for k, t := range s.cache.dat {
		if t.expired(expiry) {
			list = append(list, k)
		}
	}
	for _, key := range list {
		s.cache.set(key, nil)
	}
	s.cache.mu.Unlock()
	afterFunc(s.GC, s.collectgarbage)
}
