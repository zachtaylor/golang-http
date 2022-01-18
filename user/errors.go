package user

import "errors"

var (
	// ErrSessionSync indicates a caching issue with session.Manager
	ErrSessionSync = errors.New("user: session out of sync")

	// ErrMissingConn indicates a write has no available websockets
	ErrMissingConn = errors.New("user: missing conn instance")
)
