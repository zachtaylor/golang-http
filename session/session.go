package session // import "taylz.io/http/session"

import (
	"errors"
	"time"

	"taylz.io/maps"
)

var (
	// ErrNoID indicates that no identifying information was provided
	ErrNoID = errors.New("session: id missing")
	// ErrExpired indicates any referenced session has since expired
	ErrExpired = errors.New("session: expired")
)

// Manager is an interface for dealing with sessions
type Manager interface {
	HTTPReader
	HTTPWriter
	// Size returns the current len of the map
	Size() int
	// Get returns a Session by ID
	Get(id string) *T
	// Must refreshes and returns the Session with the given username if one exists, and creates one if necessary
	Must(name string) *T
	// Observe adds an Observer
	Observe(Observer)
	// Update resets the internal expiry time of a Session
	Update(id string) error
	// Delete removes a Session
	Delete(id string)
}

type Observer = maps.Observer[string, *T]

type ObserverFunc = maps.ObserverFunc[string, *T]

// T is a session
type T struct {
	time time.Time
	id   string
	name string
}

// New creates a session
func New(id, name string, time time.Time) *T {
	return &T{time: time, id: id, name: name}
}

// ID returns the SessionID
func (t *T) ID() string { return t.id }

// Name returns the name given during creation
func (t *T) Name() string { return t.name }

// Expires returns the time this session expires
func (t *T) Expires() time.Time { return t.time }

// SetExpires updates the time this session expires
func (t *T) SetExpires(time time.Time) { t.time = time }

// Expired returns t.Expires().Before(time.Now())
func (t *T) Expired() bool { return t.time.Before(time.Now()) }
