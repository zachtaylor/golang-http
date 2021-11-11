package user

import (
	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

// Settings is for Manager
type Settings struct {
	Sessions *session.Manager
	Sockets  *websocket.Manager
}

// NewSettings creates Settings
func NewSettings(sessions *session.Manager, sockets *websocket.Manager) Settings {
	return Settings{
		Sessions: sessions,
		Sockets:  sockets,
	}
}
