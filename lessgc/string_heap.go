package lessgc

import (
	"reflect"
	"unsafe"
)

type String struct {
	reflect.StringHeader
}

func (s String) Get() string {
	if s.Data == 0 {
		return ""
	}
	//goland:noinspection ALL
	return unsafe.String((*byte)(unsafe.Pointer(s.Data)), s.Len)
}

type StringHeap struct {
	b SliceHeap[byte]
}

func NewStringHeap() *StringHeap {
	return &StringHeap{
		b: SliceHeap[byte]{
			b: Heap[[]byte]{
				nextInsertIndex: segmentCapacity,
			},
		},
	}
}

func (h *StringHeap) Store(v string) String {
	if len(v) == 0 {
		return String{}
	}
	s := h.b.Store([]byte(v))

	return String{
		StringHeader: reflect.StringHeader{
			Data: s.Data,
			Len:  s.Len,
		},
	}
}
