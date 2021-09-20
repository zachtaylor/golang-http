package session

import "errors"

var (
	// ErrNoCookie is returned by Manager.GetRequestCookie when no cookie is found
	ErrNoCookie = errors.New("cookie not found")
	// ErrExpired is returned by Manager.GetRequestCookie when cookie is expired
	ErrExpired = errors.New("session expired")
)
