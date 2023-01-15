package store

import (
	"golang.org/x/exp/constraints"
)

type Store[K constraints.Ordered, V any] interface {
	Get(K) (V, bool)
	Put(K, V)
	GetAll() interface{}
}
