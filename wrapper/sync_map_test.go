package wrapper

import (
	"sync"
	"testing"

	"github.com/delichik/go-pkgs/utils"
)

func BenchmarkStoreOriginal(b *testing.B) {
	m := sync.Map{}
	for i := 0; i < b.N; i++ {
		m.Store(utils.RandomStringN(10), "value")
	}
}

func BenchmarkStoreWrapper(b *testing.B) {
	m := SyncMap[string, string]{}
	for i := 0; i < b.N; i++ {
		m.Store(utils.RandomStringN(10), "value")
	}
}

func BenchmarkLoadOriginal(b *testing.B) {
	m := sync.Map{}
	list := make([]string, 1024)
	for i := range list {
		list[i] = utils.RandomStringN(10)
		m.Store(list[i], "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(list[i%1024])
	}
}

func BenchmarkLoadWrapper(b *testing.B) {
	m := SyncMap[string, string]{}
	list := make([]string, 1024)
	for i := range list {
		list[i] = utils.RandomStringN(10)
		m.Store(list[i], "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(list[i%1024])
	}
}
