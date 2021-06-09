package websocket

import (
	"taylz.io/http/keygen"
	"taylz.io/http/session"
)

type Settings struct {
	Sessions *session.Manager
	Keygen   keygen.Func
	Handler  Handler
}

// DefaultSettings uses keygen.Default
func DefaultSettings(sessions *session.Manager, handler Handler) Settings {
	return Settings{
		Sessions: sessions,
		Keygen:   keygen.Default,
		Handler:  handler,
	}
}
