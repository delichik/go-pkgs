package lessgc

import (
	"runtime"
	"strings"
	"testing"
)

// 全局变量用于防止编译器优化
var (
	benchSinkInt int
)

// ---------- StringHeap 对比 ----------

func BenchmarkStringHeap_vs_Direct(b *testing.B) {
	const size = 1024
	sh := NewStringHeap()
	original := strings.Repeat("x", size)
	ref := sh.Store(original)

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s := ref.Get()
			benchSinkInt += int(s[len(s)-1])
		}
	})

	b.Run("Direct", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s := original
			benchSinkInt += int(s[len(s)-1])
		}
	})
	runtime.KeepAlive(sh)
}

// ---------- SliceHeap 对比 ----------

func BenchmarkSliceHeap_vs_Direct(b *testing.B) {
	const size = 1024
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = i
	}

	sh := NewSliceHeap[int]()
	ref := sh.Store(data)

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s := ref.Get()
			benchSinkInt += s[0] + s[len(s)-1]
		}
	})

	b.Run("Direct", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s := data
			benchSinkInt += s[0] + s[len(s)-1]
		}
	})
	runtime.KeepAlive(sh)
}

// ---------- MapHeap 对比 ----------

func BenchmarkMapHeap_vs_Direct(b *testing.B) {
	const size = 1024
	m := make(map[int]int, size)
	for i := 0; i < size; i++ {
		m[i] = i + 1
	}

	mh := NewMapHeap[int, int]()
	ref := mh.Store(m)

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mm := ref.Get()
			benchSinkInt += mm[i%size]
		}
	})

	b.Run("Direct", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mm := m
			benchSinkInt += mm[i%size]
		}
	})
	runtime.KeepAlive(mh)
}

// ---------- 通用 Heap[T] 对比 ----------

type benchStruct struct {
	A int
	B int
	C int
}

func BenchmarkHeapStruct_vs_Direct(b *testing.B) {
	v := benchStruct{A: 1, B: 2, C: 3}
	h := NewHeap[benchStruct]()
	p := h.Store(v)

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pv := p.Get()
			benchSinkInt += pv.A + pv.B + pv.C
		}
	})

	b.Run("Direct", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchSinkInt += v.A + v.B + v.C
		}
	})
	runtime.KeepAlive(h)
}
