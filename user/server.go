package user

import (
	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

// Server is another dependency type for Manager
type Server interface {
	// GetSessionManager returns the *session.Manager instance
	GetSessionManager() *session.Manager
	// GetWebsocketManager returns the *websocket.Manager instance
	GetWebsocketManager() *websocket.Manager
}
