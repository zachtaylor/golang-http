package user

import (
	"errors"

	"taylz.io/http/session"
)

var (
	// ErrNotFound is returned by Manager.Get when a named user is missing in cache
	ErrNotFound = errors.New("user not found")
	// ErrExpired(=session.ErrExpired) is returned by the Manager when the session is expired
	ErrExpired = session.ErrExpired
	// ErrSessionSync is returned by Manager.GetRequestCookie when a caching issue arises with session.Manager
	ErrSessionSync = errors.New("session out of sync")
)
