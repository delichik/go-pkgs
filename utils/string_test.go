package utils

import (
	"fmt"
	"testing"
)

func TestRandomStringN(t *testing.T) {
	fmt.Println(RandomStringN(10))
}

func BenchmarkRandomStringN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandomStringN(10)
	}
}
