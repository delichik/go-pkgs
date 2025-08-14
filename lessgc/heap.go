package lessgc

import "unsafe"

type Pointer[T any] uintptr

func (p Pointer[T]) Get() *T {
	if p == 0 {
		return nil
	}
	//goland:noinspection ALL
	return (*T)(unsafe.Pointer(p))
}

type Heap[T any] struct {
	noCopy          noCopy
	segments        []segment[T]
	nextInsertIndex int
}

func NewHeap[T any]() *Heap[T] {
	return &Heap[T]{
		nextInsertIndex: segmentCapacity,
	}
}

func (h *Heap[T]) getEmpty() (segment[T], int) {
	if h.nextInsertIndex == segmentCapacity {
		h.segments = append(h.segments, make(segment[T], segmentCapacity))
		h.nextInsertIndex = 0
	}
	sg := h.segments[len(h.segments)-1]
	i := h.nextInsertIndex
	h.nextInsertIndex++
	return sg, i
}

func (h *Heap[T]) Store(v T) Pointer[T] {
	sg, i := h.getEmpty()
	sg[i] = v
	return Pointer[T](unsafe.Pointer(&(sg[i])))
}
