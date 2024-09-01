package safemap

import (
	"sync"
)

type SafeMap struct {
	mu sync.RWMutex
	m  map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[interface{}]interface{}),
	}
}

func (sm *SafeMap) Get(key interface{}) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	val, ok := sm.m[key]
	return val, ok
}

func (sm *SafeMap) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.m)
}

func (sm *SafeMap) Set(key interface{}, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.m[key] = value
}

func (sm *SafeMap) Delete(key interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.m, key)
}

func (sm *SafeMap) Range(f func(key, value interface{}) bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for k, v := range sm.m {
		if !f(k, v) {
			break
		}
	}
}
