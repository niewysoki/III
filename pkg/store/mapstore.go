package store

import (
	"golang.org/x/exp/constraints"
	"sync"
)

type ThreadSafeMapStore[K constraints.Ordered, V any] struct {
	sync.RWMutex
	m map[K]V
}

func _[K constraints.Ordered, V any]() {
	var _ Store[K, V] = &ThreadSafeMapStore[K, V]{}
}

func NewThreadSafeMapStore[K constraints.Ordered, V any]() *ThreadSafeMapStore[K, V] {
	return &ThreadSafeMapStore[K, V]{
		m: make(map[K]V, 0),
	}
}

func (tsms *ThreadSafeMapStore[K, V]) Put(k K, v V) {
	tsms.Lock()
	tsms.m[k] = v
	tsms.Unlock()
}

func (tsms *ThreadSafeMapStore[K, V]) Get(k K) (V, bool) {
	tsms.RLock()
	v, ok := tsms.m[k]
	tsms.RUnlock()

	return v, ok
}

func (tsms *ThreadSafeMapStore[K, V]) GetAll() interface{} {
	m := make(map[K]V, 0)

	tsms.RLock()
	for k, v := range tsms.m {
		m[k] = v
	}
	tsms.RUnlock()

	return m
}
