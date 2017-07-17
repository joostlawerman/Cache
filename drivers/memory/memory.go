package memory

import (
	"sync"
)

type Driver struct {
	meta     *Meta
	data    map[string]interface{}
	sync.RWMutex
}

func NewDriver() *Driver {
	return &Driver{meta: NewMeta()}
}

func (d *Driver) Meta() *Meta {
	return d.meta
}

func (d *Driver) Put(key string, value interface{}) error {
	d.Lock()
	defer d.Unlock()

	d.data[key] = value
	return nil
}

func (d *Driver) Get(key string) (interface{}, error) {
	d.Lock()
	defer d.Unlock()

	return d.data[key], nil
}

func (d *Driver) Exists(key string) bool {
	d.Lock()
	defer d.Unlock()

	_, ok := d.data[key]
	return ok
}

func (d *Driver) Remove(key string) error {
	d.Lock()
	defer d.Unlock()

	delete(d.data, key)
	return nil
}
