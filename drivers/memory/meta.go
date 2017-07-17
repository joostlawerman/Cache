package memory

import (
	"time"
	"sync"
)

type Meta struct {
	Data map[string]time.Time
	sync.RWMutex
}

func NewMeta()*Meta {
	return &Meta{ Data: make(map[string]time.Time) }
}

func (m *Meta) List() map[string]time.Time {
	m.Lock()
	defer m.Unlock()

	return m.Data
}

func (m *Meta) Set(key string, createdAt time.Time) {
	m.Lock()
	defer m.Unlock()

	m.Data[key] = createdAt
}

func (m *Meta) Get(key string) (time.Time, bool) {
	m.Lock()
	defer m.Unlock()

	createdAt, ok := m.Data[key]

	return createdAt, ok
}

func (m *Meta) Delete(key string) {
	m.Lock()
	defer m.Unlock()

	delete(m.Data, key)
}

func (m *Meta) IsExpired(key string, duration time.Duration) bool {
	if meta, ok := m.Get(key); ok {
		return time.Now().After(meta.Add(duration))
	}
	return false
}
