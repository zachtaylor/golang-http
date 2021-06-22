package user

import (
	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

type Settings struct {
	Sessions *session.Manager
	Sockets  *websocket.Manager
}

func NewSettings(sessions *session.Manager, sockets *websocket.Manager) Settings {
	return Settings{
		Sessions: sessions,
		Sockets:  sockets,
	}
}
