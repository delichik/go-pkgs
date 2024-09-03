package chartree

import (
	"strings"
	"testing"

	"github.com/delichik/go-pkgs/utils"
)

func TestCharTree(t *testing.T) {

	m := NewManager()

	m.Add(MatchContain, "a", "a", 12121)
	for i := range 100 {
		r := utils.RandomStringN(12)
		m.Add(MatchContain, r, r, i+1)
	}
	m.Add(MatchContain, "ai", "ai", 102)
	m.Add(MatchContain, "main", "main", 103)
	m.Add(MatchPrefix, "mai", "mai", 1000)
	for i := range 100 {
		r := utils.RandomStringN(12)
		m.Add(MatchPrefix, r, r, i+1001)
	}
	m.Add(MatchPrefix, "ma", "ma", 1102)
	m.Add(MatchPrefix, "main", "main", 1103)
	m.Add(MatchSuffix, "in", "in", 2000)
	for i := range 100 {
		r := utils.RandomStringN(12)
		m.Add(MatchSuffix, r, r, i+2001)
	}
	m.Add(MatchSuffix, "ain", "ain", 2102)
	m.Add(MatchSuffix, "main", "main", 0)

	r := m.Find("main")
	if r.s != "main" || r.order != 0 {
		t.Error(r)
	}
}

func BenchmarkCharTree_100(b *testing.B) {
	_charTree(b, 100)
}

func BenchmarkArray_100(b *testing.B) {
	_array(b, 100)
}

func BenchmarkCharTree_1000(b *testing.B) {
	_charTree(b, 1000)
}

func BenchmarkArray_1000(b *testing.B) {
	_array(b, 1000)
}

func BenchmarkCharTree_10000(b *testing.B) {
	_charTree(b, 10000)
}

func BenchmarkArray_10000(b *testing.B) {
	_array(b, 10000)
}

func _charTree(b *testing.B, n int) {
	m := NewManager()

	for i := range n {
		t := utils.RandomStringN(12)
		m.Add(MatchContain, t, t, i+1)
	}
	for i := range n {
		t := utils.RandomStringN(12)
		m.Add(MatchPrefix, t, t, i+1+n*10)
	}
	for i := range n {
		t := utils.RandomStringN(12)
		m.Add(MatchSuffix, t, t, i+1+n*20)
	}
	m.Add(MatchSuffix, "main", "ain", 0)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Find("mainmainmainmain")
	}
}

type _T struct {
	match Match
	path  string
}

func _array(b *testing.B, n int) {
	var l []_T

	for range n {
		t := utils.RandomStringN(12)
		l = append(l, _T{MatchContain, t})
	}
	for range n {
		t := utils.RandomStringN(12)
		l = append(l, _T{MatchPrefix, t})
	}
	for range n {
		t := utils.RandomStringN(12)
		l = append(l, _T{MatchSuffix, t})
	}
	l = append(l, _T{MatchSuffix, "ain"})

	s := "mainmainmainmain"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, t := range l {
			switch t.match {
			case MatchContain:
				if strings.Contains(s, t.path) {
					break
				}
			case MatchPrefix:
				if strings.HasPrefix(s, t.path) {
					break
				}
			case MatchSuffix:
				if strings.HasSuffix(s, t.path) {
					break
				}
			}
		}
	}
}
