package file

import (
	"time"
	"os"
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Meta struct {
	Data map[string]time.Time `json:"list"`
	sync.RWMutex
}

func (m *Meta) List() map[string]time.Time {
	m.Lock()
	defer m.Unlock()

	return m.Data
}

func (m *Meta) Set(name string, createdAt time.Time) {
	m.Lock()
	defer m.Unlock()

	m.Data[name] = createdAt
}

func (m *Meta) Get(name string) (time.Time, bool) {
	m.Lock()
	defer m.Unlock()

	createdAt, ok := m.Data[name]

	return createdAt, ok
}

func (m *Meta) Delete(name string) {
	m.Lock()
	defer m.Unlock()

	delete(m.Data, name)
}

type FileMeta struct {
	Location string
	FileMode os.FileMode
	Meta *Meta
}

func (f *FileMeta) Boot() error {
	meta := &Meta{Data: make(map[string]time.Time) }

	file, err := os.Open(f.Location + ".meta.json")
	if os.IsNotExist(err) {
		if err := ioutil.WriteFile(f.Location + ".meta.json", []byte("{ \"list\": {} }"), 0777); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if err == nil {
		if err := json.NewDecoder(file).Decode(&meta); err != nil {
			return err
		}
	}
	f.Meta = meta

	return nil
}

func (f *FileMeta) store() error {
	file, err := os.Create(f.Location + ".meta.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(f.Meta)
}

func (f *FileMeta) List() map[string]time.Time {
	return f.Meta.List()
}

func (f *FileMeta) Put(name string) error {
	f.Meta.Set(name, time.Now())

	return f.store()
}

func (f *FileMeta) Remove(name string) error {
	f.Meta.Delete(name)

	return f.store()
}

func (f *FileMeta) IsExpired(name string, duration time.Duration) bool {
	if meta, ok := f.Meta.Get(name); ok {
		return time.Now().After(meta.Add(duration))
	}
	return false
}