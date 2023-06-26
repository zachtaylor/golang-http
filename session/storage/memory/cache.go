package session_memory

import (
	"taylz.io/http/session"
	"taylz.io/maps"
)

type Cache = maps.Observable[string, *session.T]

func NewCache() *Cache { return maps.NewObservable[string, *session.T]() }
