package session

import "errors"

var (
	// ErrNoID indicates that no identifying information was provided
	ErrNoID = errors.New("session: no id")
	// ErrExpired indicates any referenced session has since expired
	ErrExpired = errors.New("session: expired")
)
