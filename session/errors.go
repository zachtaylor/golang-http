package session

import "errors"

var (
	// ErrNoID indicates that no identifying information was provided
	ErrNoID = errors.New("session id missing")
	// ErrExpired indicates any referenced session has since expired
	ErrExpired = errors.New("session expired")
)
