package session

import "time"

// T is a session
type T struct {
	time time.Time
	id   string
	name string
}

// New creates a session
func New(t time.Time, id, name string) *T {
	return &T{time: t, id: id, name: name}
}

// Time returns the last updated time
func (t *T) Time() time.Time { return t.time }

// Update refreshes the session
func (t *T) Update() { t.time = time.Now() }

// ID returns the SessionID
func (t *T) ID() string { return t.id }

// Name returns the name given during creation
func (t *T) Name() string { return t.name }
