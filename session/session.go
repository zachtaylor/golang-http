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
