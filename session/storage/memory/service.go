package session_memory

import (
	"strings"
	"time"

	"taylz.io/http"
	"taylz.io/http/security"
	"taylz.io/http/session"
)

// Service implements Manager with a *Cache
type Service struct {
	Settings
	cache  *Cache
	keygen func() string
	reader session.HTTPReader
}

// NewService creates a Service
func NewService(settings Settings, keygen func() string, readerBuilder func(*Service) session.HTTPReader) *Service {
	service := &Service{
		Settings: settings,
		cache:    NewCache(),
		keygen:   keygen,
	}
	service.reader = readerBuilder(service)
	time.AfterFunc(service.Lifetime, service.gc)
	return service
}

// Must refreshes and returns the Session with the given username if one exists, and creates one if necessary
func (s *Service) Must(name string) (t *session.T) {
	find := func(k string, v *session.T) bool { return v.Name() == name }
	if _, t = s.cache.Find(find); t != nil {
		t.SetExpires(time.Now().Add(s.Lifetime))
	} else {
		var id string
		for ok := true; ok; ok = (s.Get(id) != nil) {
			id = s.keygen()
		}
		t = session.New(id, name, time.Now().Add(s.Lifetime))
		s.cache.Set(id, t)
	}
	return
}

// Get returns a Session by ID
func (s *Service) Get(id string) *session.T { return s.cache.Get(id) }

// Size returns the current len of the map
func (s *Service) Size() int { return s.cache.Size() }

// Update changes the internal expiry time of a Session
func (s *Service) Update(id string) (err error) {
	s.cache.WithLock(func() {
		if t := s.Get(id); t == nil {
			err = session.ErrExpired
		} else {
			t.SetExpires(time.Now().Add(s.Lifetime))
		}
	})
	return
}

// Delete removes a Session
func (s *Service) Delete(id string) { s.cache.Delete(id) }

// Observe adds an Observer
func (s *Service) Observe(f session.Observer) { s.cache.Observe(f) }

func CookieReader(s *Service) session.HTTPReader {
	return session.HTTPReaderFunc(func(r *http.Request) (t *session.T, err error) {
		if cookie, _err := r.Cookie(s.CookieID); _err != nil {
			err = session.ErrNoID
		} else if t = s.Get(cookie.Value); t == nil {
			err = session.ErrExpired
		}
		return
	})
}

var errBearerTokenFormat = http.StatusError(http.StatusBadRequest, "Authorization format bad")

func BearerReader(s *Service) session.HTTPReader {
	return session.HTTPReaderFunc(func(r *http.Request) (t *session.T, err error) {
		if bearer := r.Header.Get("Authorization"); bearer == "" {
			err = session.ErrNoID
		} else if len(bearer) != 43 {
			err = errBearerTokenFormat
		} else if bearer[:6] != "Bearer" {
			err = errBearerTokenFormat
		} else if token := bearer[7:]; strings.Trim(token, security.CHARS_UUID) != "" {
			err = errBearerTokenFormat
		} else if t = s.Get(token); t == nil {
			err = session.ErrExpired
		}
		return
	})
}

func (s *Service) ReadHTTP(r *http.Request) (t *session.T, err error) {
	return s.reader.ReadHTTP(r)
}

// WriteHTTP writes the Set-Cookie header
func (s *Service) WriteHTTP(w http.ResponseWriter, t *session.T) {
	if t == nil {
		w.Header().Set("Set-Cookie", s.CookieID+"=; Path=/; Expires==Thu, 01 Jan 1970 00:00:00 GMT;")
		return
	}
	header := s.CookieID + "=" + t.ID()
	if s.Domain != "" {
		header += "; Domain=" + s.Domain
	}
	if s.Path != "" {
		header += "; Path=" + s.Path
	}
	if s.Secure {
		header += "; Secure"
	}
	if s.SameSite != "" {
		header += "; SameSite=" + s.SameSite
	}
	if s.HttpOnly {
		header += "; HttpOnly"
	}
	w.Header().Set("Set-Cookie", header)
}

func (s *Service) gc() {
	for now := range time.Tick(s.GC) {
		s.cache.DeleteFunc(func(k string, v *session.T) bool { return v.Expires().Before(now) })
	}
}
