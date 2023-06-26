package user

import "taylz.io/http"

// type HTTPReader is an interface for recognizing Users in http.Request
type HTTPReader interface {
	// ReadHTTP returns the User and Session
	ReadHTTP(*http.Request) (*T, error)
}

// type HTTPWriter is an interface for writing Users to http.ResponseWriter
type HTTPWriter interface {
	// WriteHTTP writes the Set-Cookie header using the session.Manager
	WriteHTTP(http.ResponseWriter, *T)
}

func (s *Service) ReadHTTP(r *http.Request) (*T, error) {
	if session, err := s.sessions.ReadHTTP(r); session == nil {
		return nil, err
	} else if user := s.Get(session.Name()); user != nil {
		return user, nil
	}
	return nil, ErrSessionSync
}

func (s *Service) WriteHTTP(w http.ResponseWriter, user *T) {
	s.sessions.WriteHTTP(w, user.session)
}
