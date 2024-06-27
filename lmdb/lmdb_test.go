package lmdb

import (
	"runtime"
	"strconv"
	"testing"
)

type TestStruct struct {
	Str string
	Num int
}

func TestLMDB_Get(t *testing.T) {
	lmdb := NewLMDB[string, TestStruct]()
	_, ok := lmdb.Get("test")
	if ok {
		t.Errorf("should not be ok")
		t.FailNow()
	}
}

func TestLMDB_Set(t *testing.T) {
	lmdb := NewLMDB[string, TestStruct]()
	lmdb.Set("test", TestStruct{
		Str: "123",
		Num: 123,
	})
	_, ok := lmdb.Get("test")
	if !ok {
		t.Errorf("should be ok")
		t.FailNow()
	}
	if lmdb.Len() != 1 {
		t.Errorf("should be 1")
		t.FailNow()
	}
}

func TestLMDB_Delete(t *testing.T) {
	lmdb := NewLMDB[string, TestStruct]()
	lmdb.Set("test", TestStruct{
		Str: "123",
		Num: 123,
	})
	_, ok := lmdb.Get("test")
	if !ok {
		t.FailNow()
	}
	lmdb.Delete("test")
	_, ok = lmdb.Get("test")
	if ok {
		t.Errorf("delete should not be ok")
		t.FailNow()
	}
	if lmdb.Len() != 0 {
		t.Errorf("length should be 0")
		t.FailNow()
	}
}

func TestLMDB_Fragment(t *testing.T) {
	lmdb := NewLMDB[string, TestStruct]()
	for i := range fragmentSize + 1 {
		lmdb.Set("test"+strconv.Itoa(i), TestStruct{
			Str: "123",
			Num: 123,
		})
	}
	if lmdb.Len() != fragmentSize+1 {
		t.Errorf("length should be %d", fragmentSize+1)
		t.FailNow()
	}
}

func BenchmarkMap_Get(b *testing.B) {
	db := map[string]*TestStruct{}
	for i := range fragmentSize + 1 {
		db["test"+strconv.Itoa(i)] = &TestStruct{
			Str: "123",
			Num: 123,
		}
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, ok := db["test"+strconv.Itoa(i%fragmentSize+1)]
		if !ok {
			b.FailNow()
		}
	}
}

func BenchmarkLMDB_Get(b *testing.B) {
	lmdb := NewLMDB[string, TestStruct]()
	for i := range fragmentSize + 1 {
		lmdb.Set("test"+strconv.Itoa(i), TestStruct{
			Str: "123",
			Num: 123,
		})
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, ok := lmdb.Get("test" + strconv.Itoa(i%fragmentSize+1))
		if !ok {
			b.FailNow()
		}
	}
}

func BenchmarkMap_Set(b *testing.B) {
	db := map[string]*TestStruct{}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		db["test"+strconv.Itoa(i)] = &TestStruct{
			Str: "123",
			Num: 123,
		}
	}
}

func BenchmarkLMDB_Set(b *testing.B) {
	lmdb := NewLMDB[string, TestStruct]()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lmdb.Set("test"+strconv.Itoa(i), TestStruct{
			Str: "123",
			Num: 123,
		})
	}
}

func BenchmarkMap_GC(b *testing.B) {
	db := map[string]*TestStruct{}
	for i := range 100000 {
		db["test"+strconv.Itoa(i)] = &TestStruct{
			Str: "123",
			Num: 123,
		}
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		runtime.GC()
	}
	runtime.KeepAlive(db)
}

func BenchmarkLMDB_GC(b *testing.B) {
	lmdb := NewLMDB[string, TestStruct]()
	for i := range 100000 {
		lmdb.Set("test"+strconv.Itoa(i), TestStruct{
			Str: "123",
			Num: 123,
		})
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		runtime.GC()
	}
	runtime.KeepAlive(lmdb)
}

func BenchmarkMap_Empty_GC(b *testing.B) {
	db := map[string]*TestStruct{}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		runtime.GC()
	}
	runtime.KeepAlive(db)
}

func BenchmarkLMDB_Empty_GC(b *testing.B) {
	lmdb := NewLMDB[string, TestStruct]()
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		runtime.GC()
	}
	runtime.KeepAlive(lmdb)
}
