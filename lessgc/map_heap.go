package lessgc

import (
	"cmp"
	"unsafe"
)

type Map[K cmp.Ordered, V any] uintptr

func (m Map[K, V]) Get() map[K]V {
	if m == 0 {
		return nil
	}
	//goland:noinspection ALL
	return *(*map[K]V)(unsafe.Pointer(m))
}

type MapHeap[K cmp.Ordered, V any] struct {
	b Heap[map[K]V]
}

func NewMapHeap[K cmp.Ordered, V any]() *MapHeap[K, V] {
	return &MapHeap[K, V]{
		b: Heap[map[K]V]{
			nextInsertIndex: segmentCapacity,
		},
	}
}

func (h *MapHeap[K, V]) Store(v map[K]V) Map[K, V] {
	if len(v) == 0 {
		return Map[K, V](0)
	}
	p := h.b.Store(v)
	return Map[K, V](p)
}
