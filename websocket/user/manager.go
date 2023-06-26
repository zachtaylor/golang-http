package user

import (
	"taylz.io/http/websocket"
	"taylz.io/maps"
)

type Observer = maps.Observer[string, *T]

type ObserverFunc = maps.ObserverFunc[string, *T]

// Manager is a user manager
type Manager interface {
	HTTPReader
	HTTPWriter
	// Size returns the current len of the map
	Size() int
	// Get returns a user by name
	Get(string) *T
	// Must (re)authorizes a session for a websocket
	Must(*websocket.T, string) *T
	// GetWebsocket returns a user by websocket
	GetWebsocket(*websocket.T) *T
	// Observe adds a callback CacheObserver
	Observe(Observer)
}
