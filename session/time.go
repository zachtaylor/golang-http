package session

import "time"

func now() time.Time { return time.Now() }

func GetID(expiry time.Time, cache *Cache, id string) (session *T) {
	if session = cache.dat[id]; session != nil && session.expired(expiry) {
		session = nil
	}
	return
}

func GetName(expiry time.Time, cache *Cache, name string) (session *T) {
	cache.mu.Lock()
	session = getName(expiry, cache, name)
	cache.mu.Unlock()
	return
}

func getName(expiry time.Time, cache *Cache, name string) (session *T) {
	for _, t := range cache.dat {
		if t.name == name {
			if !t.expired(expiry) {
				session = t
			}
			break
		}
	}
	return
}

func afterFunc(d time.Duration, f func()) { time.AfterFunc(d, f) }
