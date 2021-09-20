package user

import "errors"

var (
	// ErrNotFound is returned by the Manager.GetName when a named user is missing in cache
	ErrNotFound = errors.New("user not found")
	// ErrSessionSync is returned by the Manager when a caching issue arises with session.Manager
	ErrSessionSync = errors.New("session out of sync")
)
