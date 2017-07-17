package file

import (
	"time"
	"os"
	"encoding/json"
	"github.com/joostlawerman/cache/drivers/memory"
)

type Meta struct {
	Location string
	data     *memory.Meta
}

func NewMeta(location string) (*Meta, error) {
	data := memory.NewMeta()

	file, err := os.Open(location + ".meta.json")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		if err := json.NewDecoder(file).Decode(&data); err != nil {
			return nil, err
		}
	}
	return &Meta{Location: location, data: data}, nil
}

func (m *Meta) store() error {
	file, err := os.Create(m.Location + ".meta.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(m.data)
}

func (m *Meta) List() map[string]time.Time {
	return m.data.List()
}

func (m *Meta) Put(key string) error {
	m.data.Set(key, time.Now())

	return m.store()
}

func (m *Meta) Remove(key string) error {
	m.data.Delete(key)

	return m.store()
}

func (m *Meta) IsExpired(key string, duration time.Duration) bool {
	if meta, ok := m.data.Get(key); ok {
		return time.Now().After(meta.Add(duration))
	}
	return false
}
