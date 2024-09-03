package chartree

var mapping = []uint8{
	// NUL 	SOH 	STX 	ETX 	EOT 	ENQ 	ACK 	BEL 	 BS  	 HT  	 LF  	 VT  	 FF  	 CR  	 SO  	 SI
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// DLE 	DC1 	DC2 	DC3 	DC4 	NAK 	SYN 	ETB 	CAN 	 EM  	SUB 	ESC 	 FS  	 GS  	 RS  	 US
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// SP  	! 	" 	# 	$ 	% 	& 	' 	( 	) 	* 	+ 	, 	- 	. 	/
	0, 1, 0, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 5, 6, 0,
	// 0 	1 	2 	3 	4 	5 	6 	7 	8 	9 	: 	; 	< 	= 	> 	?
	7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 0, 0, 0, 0,
	// @ 	A 	B 	C 	D 	E 	F 	G 	H 	I 	J 	K 	L 	M 	N 	O
	0, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33,
	// P 	Q 	R 	S 	T 	U 	V 	W 	X 	Y 	Z 	[ 	\ 	] 	^ 	_
	34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 0, 0, 0, 0, 45,
	// ` 	a 	b 	c 	d 	e 	f 	g 	h 	i 	j 	k 	l 	m 	n 	o
	0, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
	// p 	q 	r 	s 	t 	u 	v 	w 	x 	y 	z 	{ 	| 	} 	~ 	DEL
	61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 0, 0, 0, 0, 0,
}

type Data struct {
	s     string
	order uint32
}

type CharRoot struct {
	*CharTree
	reversed bool
}

type CharTree struct {
	m    [MaxLength]*CharTree
	data Data
}

const MaxLength = 71

// const MaxLength = 96
// const AsciiStartIndex = 32

type Match uint8

const (
	MatchPrefix = Match(iota)
	MatchSuffix
	MatchContain
)

func NewManager() *Manager {
	return &Manager{
		matches: CharRoot{
			CharTree: &CharTree{m: [MaxLength]*CharTree{}, data: Data{}},
			reversed: false,
		},
		prefixMatches: CharRoot{
			CharTree: &CharTree{m: [MaxLength]*CharTree{}, data: Data{}},
			reversed: false,
		},
		reversedSuffixMatches: CharRoot{
			CharTree: &CharTree{m: [MaxLength]*CharTree{}, data: Data{}},
			reversed: true,
		},
	}
}

type Manager struct {
	matches               CharRoot
	prefixMatches         CharRoot
	reversedSuffixMatches CharRoot
	count                 int
}

func (m *Manager) Add(match Match, path string, s string, order int) {
	t := []byte(path)
	switch match {
	case MatchPrefix:
		m.add(&m.prefixMatches, t, s, order)
	case MatchSuffix:
		m.add(&m.reversedSuffixMatches, t, s, order)
	case MatchContain:
		m.add(&m.matches, t, s, order)
	}
}

func (m *Manager) match(node *CharTree, parsed []byte) Data {
	n := node.m[parsed[0]]
	if n == nil {
		return Data{}
	}
	if len(parsed) == 1 {
		return n.data
	} else {
		return m.match(n, parsed[1:])
	}
}

func (m *Manager) Find(s string) Data {
	var data Data
	ts := []byte(s)
	parsed := make([]uint8, len(ts))
	// nextMatches1 := []*CharTree{m.matches.CharTree, m.prefixMatches.CharTree}
	// var nextMatches2 []*CharTree
	i := 0
	for ; i < len(parsed); i += 4 {
		parsed[i] = mapping[ts[i]-1]
		parsed[i+1] = mapping[ts[i+1]-1]
		parsed[i+2] = mapping[ts[i+2]-1]
		parsed[i+3] = mapping[ts[i+3]-1]
	}
	i -= 4
	for ; i < len(parsed); i++ {
		parsed[i] = mapping[ts[i]-1]
	}

	ndata := m.match(m.prefixMatches.CharTree, parsed)
	if ndata.s != "" && (data.s == "" || ndata.order < data.order) {
		data = ndata
	}

	nextMatch := m.reversedSuffixMatches.CharTree
	for i := len(parsed) - 1; i >= 0; i-- {
		c := parsed[i]
		nextMatch = nextMatch.m[c]
		if nextMatch == nil {
			break
		}
		if data.s == "" || nextMatch.data.order < data.order {
			data = nextMatch.data
		}
	}
	i = 0
	for ; i < len(parsed); i += 4 {
		ndata1 := m.match(m.matches.CharTree, parsed[i:])
		ndata2 := m.match(m.matches.CharTree, parsed[i+1:])
		ndata3 := m.match(m.matches.CharTree, parsed[i+2:])
		ndata4 := m.match(m.matches.CharTree, parsed[i+3:])
		if ndata1.s != "" && (data.s == "" || ndata1.order < data.order) {
			data = ndata1
		}
		if ndata2.s != "" && (data.s == "" || ndata2.order < data.order) {
			data = ndata2
		}
		if ndata3.s != "" && (data.s == "" || ndata3.order < data.order) {
			data = ndata3
		}
		if ndata4.s != "" && (data.s == "" || ndata4.order < data.order) {
			data = ndata4
		}
	}
	i -= 4
	for ; i < len(parsed); i++ {
		ndata1 := m.match(m.matches.CharTree, parsed[i:])
		if ndata1.s != "" && (data.s == "" || ndata1.order < data.order) {
			data = ndata1
		}
	}

	return data
}

func (m *Manager) add(t *CharRoot, path []uint8, s string, order int) {
	node := t.CharTree
	if t.reversed {
		for i := len(path) - 1; i >= 0; i-- {
			c := mapping[path[i]-1]
			if node.m[c] == nil {
				node.m[c] = &CharTree{
					m:    [MaxLength]*CharTree{},
					data: Data{},
				}
				m.count++
			}
			node = node.m[c]
		}
	} else {
		for _, c := range path {
			c = mapping[c-1]
			if node.m[c] == nil {
				node.m[c] = &CharTree{
					m:    [MaxLength]*CharTree{},
					data: Data{},
				}
				m.count++
			}
			node = node.m[c]
		}
	}

	if node.data.order == 0 || node.data.order > uint32(order) {
		node.data.s = s
		node.data.order = uint32(order)
	}
}
