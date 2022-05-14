package user

import (
	"taylz.io/http"
	"taylz.io/http/session"
	"taylz.io/http/websocket"
)

// Manager is a user manager
type Manager interface {
	// Count returns the current len of the map
	Count() int
	// Get returns a user by name
	Get(string) *T
	// Must (re)authorizes a session for a websocket
	Must(*websocket.T, string) *T
	// GetWebsocket returns a user by websocket
	GetWebsocket(*websocket.T) *T
	// Observe adds a callback CacheObserver
	Observe(Observer)
	// ReadHTTP returns the User and Session
	ReadHTTP(*http.Request) (*T, *session.T, error)
	// WriteHTTP writes the Set-Cookie header using the session.Manager
	WriteHTTP(http.ResponseWriter, *T) error
}
