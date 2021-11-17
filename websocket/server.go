package websocket

import "taylz.io/http/session"

// Server is another dependency type for Manager
type Server interface {
	// GetSessionManager returns the *session.Manager instance
	GetSessionManager() *session.Manager
	// GetSessionManager returns the websocket.Handler
	GetWebsocketHandler() Handler
}
