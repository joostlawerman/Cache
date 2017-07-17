package memory

import (
	"errors"
	"time"
	"github.com/joostlawerman/cache"
)

func Driver() cache.Driver {
	return &Memory{ Cache: make(map[string]interface{}), Meta: make(map[string]time.Time) }
}

type Memory struct {
	Meta map[string]time.Time
	Cache map[string]interface{}
}

func (m *Memory) GetMeta() (map[string]time.Time, error) {
	return m.Meta, nil
}

func (m *Memory) Put(name string, contents interface{}) error {
	m.Meta[name] = time.Now()

	m.Cache[name] = contents

	return nil
}

func (m *Memory) Get(name string) (interface{}, error) {
	if item, ok := m.Cache[name]; ok {
		return item, nil
	}
	return nil, errors.New("Item does not exist")
}

func (m *Memory) Exists(name string) bool {
	if _, ok := m.Cache[name]; ok {
		return true
	}
	return false
}

func (m *Memory) Remove(name string) error {
	delete(m.Meta, name)

	delete(m.Cache, name)

	return nil
}
