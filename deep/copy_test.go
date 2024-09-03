package deep

import (
	"encoding/json"
	"testing"

	"github.com/bytedance/sonic"
)

type _T1 struct {
	T1V   int            `json:"t_1_v"`
	T1P2  *_T2           `json:"t_1_p_2"`
	T1P3  any            `json:"t_1_p_3"`
	T1P4  [3]int         `json:"t_1_p_4"`
	T1P5  [3]int         `json:"t_1_p_5"`
	T1P6  *_T1           `json:"t_1_p_6"`
	T1P7  []int          `json:"t_1_p_7"`
	T1P8  []int          `json:"t_1_p_8"`
	T1P9  []*_T2         `json:"t_1_p_9"`
	T1P10 map[string]int `json:"t_1_p_10"`
	T1P11 map[string]int `json:"t_1_p_11"`
	T1P12 []int          `json:"t_1_p_12"`
}

type _T2 struct {
	T2P1 *_T1 `json:"t_2_p_1"`
	T2V  int  `json:"t_2_v"`
}

func _assert(t *testing.T, a bool) {
	if !a {
		t.Fail()
		panic("deepcopy failed")
	}
}

func newSrc() *_T1 {
	src := &_T1{
		T1V:  1,
		T1P2: nil,
		T1P3: nil,
	}
	src.T1P2 = &_T2{
		T2P1: src,
		T2V:  2,
	}
	src.T1P3 = src.T1P2
	src.T1P4 = [3]int{1, 2, 3}
	src.T1P5 = src.T1P4
	src.T1P6 = src
	src.T1P7 = []int{1, 2, 3}
	src.T1P8 = src.T1P7
	src.T1P9 = []*_T2{src.T1P2}
	src.T1P10 = map[string]int{"a": 1, "b": 2}
	src.T1P11 = src.T1P10
	src.T1P12 = src.T1P4[:]
	return src
}

func Test_deep_copy(t *testing.T) {
	src := newSrc()
	dst := &_T1{}
	Copy(src, dst)
	_assert(t, dst.T1V == src.T1V)
	_assert(t, dst.T1P2.T2V == src.T1P2.T2V)
	_assert(t, dst.T1P2.T2P1.T1V == src.T1P2.T2P1.T1V)
	_assert(t, dst.T1P3.(*_T2).T2V == src.T1P3.(*_T2).T2V)
	_assert(t, dst.T1P3.(*_T2).T2P1.T1V == src.T1P3.(*_T2).T2P1.T1V)
	_assert(t, dst.T1P4[0] == src.T1P4[0])
	_assert(t, dst.T1P4[1] == src.T1P4[1])
	_assert(t, dst.T1P4[2] == src.T1P4[2])
	_assert(t, dst.T1P5[0] == src.T1P5[0])
	_assert(t, dst.T1P5[1] == src.T1P5[1])
	_assert(t, dst.T1P5[2] == src.T1P5[2])
	_assert(t, dst.T1P7[0] == src.T1P7[0])
	_assert(t, dst.T1P7[1] == src.T1P7[1])
	_assert(t, dst.T1P7[2] == src.T1P7[2])
	_assert(t, dst.T1P8[0] == src.T1P8[0])
	_assert(t, dst.T1P8[1] == src.T1P8[1])
	_assert(t, dst.T1P8[2] == src.T1P8[2])
	_assert(t, dst.T1P10["a"] == src.T1P10["a"])
	_assert(t, dst.T1P10["b"] == src.T1P10["b"])
	_assert(t, dst.T1P11["a"] == src.T1P11["a"])
	_assert(t, dst.T1P11["b"] == src.T1P11["b"])

	dst.T1V = 3
	_assert(t, dst.T1V == dst.T1P2.T2P1.T1V)
	_assert(t, dst.T1V == dst.T1P3.(*_T2).T2P1.T1V)
	dst.T1P4[1] = 6
	_assert(t, dst.T1P5[0] == dst.T1P4[0])
	_assert(t, dst.T1P5[1] != dst.T1P4[1])
	_assert(t, dst.T1P5[2] == dst.T1P4[2])
	_assert(t, dst.T1P12[1] != dst.T1P4[1])
	dst.T1P7[1] = 6
	_assert(t, dst.T1P7[0] == dst.T1P8[0])
	_assert(t, dst.T1P7[1] == dst.T1P8[1])
	_assert(t, dst.T1P7[2] == dst.T1P8[2])

	dst.T1P10["b"] = 3
	_assert(t, dst.T1P11["a"] == 1)
	_assert(t, dst.T1P11["b"] == 3)
}

func Benchmark_copy_json(b *testing.B) {
	src := newSrc()
	for i := 0; i < b.N; i++ {
		dst := &_T1{}
		mid, _ := json.Marshal(src)
		_ = json.Unmarshal(mid, dst)
	}
}

func Benchmark_copy_sonic(b *testing.B) {
	src := newSrc()
	for i := 0; i < b.N; i++ {
		dst := &_T1{}
		mid, _ := sonic.Marshal(src)
		_ = sonic.Unmarshal(mid, dst)
	}
}

func Benchmark_copy_deep(b *testing.B) {
	src := newSrc()
	for i := 0; i < b.N; i++ {
		dst := &_T1{}
		Copy(src, dst)
	}
}
