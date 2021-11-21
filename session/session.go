package session // import "taylz.io/http/session"

import "time"

// T is a session
type T struct {
	time time.Time
	id   string
	name string
}

// New creates a session
func New(id, name string, time time.Time) *T {
	return &T{id: id, name: name, time: time}
}

// ID returns the SessionID
func (t *T) ID() string { return t.id }

// Name returns the name given during creation
func (t *T) Name() string { return t.name }

// Expires returns the time this session expires
func (t *T) Expires() time.Time { return t.time }

// Expired returns t.Expires().Before(time.Now())
func (t *T) Expired() bool              { return t.expired(time.Now()) }
func (t *T) expired(now time.Time) bool { return t.time.Before(now) }
