package session

import "time"

// T is a session
type T struct {
	time time.Time
	id   string
	name string
}

func New(t time.Time, id, name string) *T {
	return &T{time: t, id: id, name: name}
}

// Update refreshes the session
func (t *T) Update() { t.time = time.Now() }

// Name returns the name given during creation
func (t *T) Name() string { return t.name }
