package subject

import "sync"

type SubKey uint

type T[D any] struct {
	data    D
	version int
	Name    string
	subs    map[SubKey]func(D) error
	reclaim []SubKey
	sync    sync.RWMutex
}

func (t *T[D]) Current() (d D) {
	t.sync.RLock()
	d = t.data
	t.sync.RUnlock()
	return
}

func (t *T[D]) Version() int {
	return t.version
}

func (t *T[D]) Update(data D) {
	t.sync.RLock()
	t.version++
	t.data = data
	for k, sub := range t.subs {
		if sub(data) != nil {
			k := k
			go t.Unsubscribe(k)
		}
	}
	t.sync.RUnlock()
}

func (t *T[D]) Subscribe(callback func(D) error) (k SubKey) {
	t.sync.Lock()
	if len(t.reclaim) > 0 {
		k, t.reclaim = t.reclaim[0], t.reclaim[1:]
	} else {
		k = SubKey(len(t.subs))
	}
	t.subs[k] = callback
	t.sync.Unlock()
	return k
}

func (t *T[D]) Unsubscribe(k SubKey) {
	t.sync.Lock()
	if _, ok := t.subs[k]; ok {
		delete(t.subs, k)
		t.reclaim = append(t.reclaim, k)
	}
	t.sync.Unlock()
}
