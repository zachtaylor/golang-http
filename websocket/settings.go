package websocket

import "taylz.io/http/session"

// Settings is a websocket manager Settings
type Settings struct {
	Sessions *session.Manager
	Keygen   func() string
	Handler  Handler
}

// NewSettings returns a new Settings
func NewSettings(sessions *session.Manager, keygen func() string, handler Handler) Settings {
	return Settings{
		Sessions: sessions,
		Keygen:   keygen,
		Handler:  handler,
	}
}
