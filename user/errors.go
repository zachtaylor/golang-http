package user

import (
	"errors"

	"taylz.io/http/session"
)

var (
	// ErrNoID(=session.ErrNoID) is returned by the Manager when user information is missing
	ErrNoID = session.ErrNoID
	// ErrExpired(=session.ErrExpired) is returned by the Manager when the session is expired
	ErrExpired = session.ErrExpired
	// ErrSessionSync is returned by the Manager when a caching issue arises with session.Manager
	ErrSessionSync = errors.New("session out of sync")
)
