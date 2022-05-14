package user

import "taylz.io/yas"

type Cache = yas.Observatory[*T]

func NewCache() *Cache { return yas.NewObservatory[*T]() }

type Observer = yas.Observer[*T]

type ObserverFunc = yas.ObserverFunc[*T]
