package lmdb

const defaultFragmentCount = 1024
const fragmentSize = 128

type Index struct {
	Fragment uint
	Index    uint
}

type Item[T any] struct {
	realItem T
	deleted  bool
}

type Items[T any] struct {
	items        []Item[T]
	deletedCount uint
}

type LMDB[K comparable, V any] struct {
	length  uint
	m       map[K]Index
	storage []Items[V]
}

func NewLMDB[K comparable, V any]() *LMDB[K, V] {
	return &LMDB[K, V]{
		length:  0,
		m:       make(map[K]Index),
		storage: make([]Items[V], 0, defaultFragmentCount),
	}
}

func (m *LMDB[K, V]) Get(key K) (*V, bool) {
	index, ok := m.m[key]
	if !ok {
		return nil, false
	}

	t := m.storage[index.Fragment].items[index.Index]
	return &t.realItem, true
}

func (m *LMDB[K, V]) getNextIndex() Index {
	for fi, fragment := range m.storage {
		if fragment.deletedCount == 0 {
			continue
		}

		for ii := range fragmentSize {
			if ii >= len(fragment.items) {
				fragment.items = append(fragment.items, Item[V]{deleted: true})
				m.storage[fi] = fragment
				return Index{
					Fragment: uint(fi),
					Index:    uint(ii - 1),
				}
			}
			item := fragment.items[ii]
			if !item.deleted {
				continue
			}
			return Index{
				Fragment: uint(fi),
				Index:    uint(ii - 1),
			}
		}
	}

	m.storage = append(m.storage, Items[V]{
		items:        make([]Item[V], 1, fragmentSize),
		deletedCount: fragmentSize,
	})

	return Index{
		Fragment: uint(len(m.storage) - 1),
		Index:    0,
	}
}

func (m *LMDB[K, V]) Set(key K, value V) {
	index, ok := m.m[key]
	if ok {
		m.storage[index.Fragment].items[index.Index] = Item[V]{
			realItem: value,
			deleted:  false,
		}
	} else {
		index = m.getNextIndex()
		m.m[key] = index
		fragment := m.storage[index.Fragment]
		fragment.items[index.Index] = Item[V]{
			realItem: value,
			deleted:  false,
		}
		m.storage[index.Fragment] = fragment
	}
	m.length++
}

func (m *LMDB[K, V]) Delete(key K) {
	index, ok := m.m[key]
	if !ok {
		return
	}
	fragment := m.storage[index.Fragment]
	fragment.items[index.Index] = Item[V]{
		deleted: true,
	}
	fragment.deletedCount++
	m.storage[index.Fragment] = fragment
	m.length--
	delete(m.m, key)
}

func (m *LMDB[K, V]) Len() uint {
	return m.length
}
