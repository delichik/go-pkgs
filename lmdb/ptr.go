package lmdb

import (
	"unsafe"
)

// Ptr 使用 Ptr 替换普通的 Go 指针，可以防止 GC 扫描
// 但是需要小心的是这个 Ptr 底层的内容可能
type Ptr[T any] struct {
	p unsafe.Pointer
}

func PtrFrom[T any](v *T) Ptr[T] {
	return Ptr[T]{
		p: unsafe.Pointer(v),
	}
}

func (p *Ptr[T]) V() *T {
	return (*T)(p.p)
}
