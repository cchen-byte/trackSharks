package utils

import (
	"sync"
)


type MemorySet struct {
	set map[string]struct{}
	mx   sync.Mutex
}

func (m *MemorySet) Add(items ...string) (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	var count int
	l := len(m.set) // 原来的大小
	for _, item := range items {
		m.set[item] = struct{}{}
		nl := len(m.set) // 新的大小
		if nl != l {   // 如果是新加入, 那么新的大小一定和原来的大小不同
			count++
			l = nl
		}
	}
	return count, nil
}

func (m *MemorySet) HasItem(item string) bool {
	m.mx.Lock()
	defer m.mx.Unlock()

	_, ok := m.set[item]
	return ok
}

func (m *MemorySet) Remove(items ...string) (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	l := len(m.set) // 原来的大小
	if l == 0 {
		return 0, nil
	}

	var count int
	for _, item := range items {
		delete(m.set, item)
		nl := len(m.set) // 新的大小
		if nl == l {
			continue
		}

		// 如果真的删除了, 那么新的大小一定和原来的大小不同
		count++
		l = nl

		if l == 0 {
			break
		}
	}

	return count, nil
}

func (m *MemorySet) DeleteSet(value string) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	delete(m.set, value)
	return nil
}

func (m *MemorySet) GetSetSize() (int, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	return len(m.set), nil
}

func (m *MemorySet) Close() error { return nil }


func NewMemorySet() *MemorySet {
	return &MemorySet{
		set: make(map[string]struct{}),
	}
}

