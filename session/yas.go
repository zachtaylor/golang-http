package session

import (
	"time"

	"taylz.io/yas"
)

type Cache = yas.Observatory[*T]

func NewCache() *Cache { return yas.NewObservatory[*T]() }

type Observer = yas.Observer[*T]

type ObserverFunc = yas.ObserverFunc[*T]

// TestExpired wraps a time.Time to a yas.Tester[*T]
type TestExpired time.Time

// Test checks expired with now=this
func (t TestExpired) Test(s *T) bool { return s.expired(time.Time(t)) }

// TestName wraps a string to a yas.Tester[*T]
type TestName string

// Test checks T.Name == this
func (t TestName) Test(s *T) bool { return s != nil && s.name == string(t) }
