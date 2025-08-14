package lessgc

import (
	"reflect"
	"unsafe"
)

type Slice[T any] struct {
	reflect.SliceHeader
}

func (s Slice[T]) Get() []T {
	if s.Data == 0 {
		return nil
	}
	//goland:noinspection ALL
	return unsafe.Slice((*T)(unsafe.Pointer(s.Data)), s.Len)
}

type SliceHeap[T any] struct {
	b Heap[[]T]
}

func NewSliceHeap[T any]() *SliceHeap[T] {
	return &SliceHeap[T]{
		b: Heap[[]T]{
			nextInsertIndex: segmentCapacity,
		},
	}
}

func (h *SliceHeap[T]) Store(v []T) Slice[T] {
	if len(v) == 0 {
		return Slice[T]{}
	}
	sg, i := h.b.getEmpty()
	sg[i] = v

	return Slice[T]{
		SliceHeader: reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(&v[0])),
			Len:  len(v),
			Cap:  cap(v),
		},
	}
}
