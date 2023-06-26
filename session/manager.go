package session

import "taylz.io/maps"

// Manager is an interface for dealing with sessions
type Manager interface {
	HTTPReader
	HTTPWriter
	// Size returns the current len of the map
	Size() int
	// Get returns a Session by ID
	Get(id string) *T
	// Must refreshes and returns the Session with the given username if one exists, and creates one if necessary
	Must(name string) *T
	// Observe adds an Observer
	Observe(Observer)
	// Update resets the internal expiry time of a Session
	Update(id string) error
	// Delete removes a Session
	Delete(id string)
}

type Observer = maps.Observer[string, *T]

type ObserverFunc = maps.ObserverFunc[string, *T]
