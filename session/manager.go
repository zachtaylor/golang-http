package session

import "taylz.io/http"

// Manager is an interface for dealing with sessions
type Manager interface {
	Reader
	Writer
	// Count returns the current len of the map
	Count() int
	// Get returns a Session by ID
	Get(id string) *T
	// Must refreshes and returns the Session with the given username if one exists, and creates one if necessary
	Must(name string) *T
	// Observe adds a CacheObserver
	Observe(CacheObserver)
	// Update changes the internal expiry time of a Session
	Update(id string) error
	// Remove removes a Session
	Remove(id string)
}

// type Reader is an interface for recognizing Sessions in http.Request
type Reader interface {
	// ReadHTTP returns Session by *http.Request
	ReadHTTP(*http.Request) (*T, error)
}

type ReaderFunc func(*http.Request) (*T, error)

func (f ReaderFunc) ReadHTTP(r *http.Request) (*T, error) { return f(r) }

func ContextReader() Reader {
	return ReaderFunc(func(r *http.Request) (t *T, err error) {
		if session, ok := FromContext(r.Context()); ok {
			if session == nil {
				err = ErrExpired
			} else {
				t = session
			}
		} else {
			err = ErrNoID
		}
		return
	})
}

// type Writer is an interface for writing Sessions to http.ResponseWriter
type Writer interface {
	// WriteCookie writes the Set-Cookie header in http.ResponseWriter
	WriterHTTP(http.ResponseWriter, *T) error
}
