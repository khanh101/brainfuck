package util

import (
	"sync"
)

type MutexSyncMap struct {
	mu sync.Mutex
	m  map[string]interface{}
}

func NewMutexSyncMap() *MutexSyncMap {
	return &MutexSyncMap{
		m: make(map[string]interface{}),
	}
}

func (sm *MutexSyncMap) Load(key string) (value interface{}, ok bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, ok = sm.m[key]
	return
}

func (sm *MutexSyncMap) Store(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

func (sm *MutexSyncMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
}

func (sm *MutexSyncMap) Range(f func(key interface{}, value interface{}) bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for k, v := range sm.m {
		if !f(k, v) {
			break
		}
	}
}
