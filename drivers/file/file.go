package file

import (
	"os"
	"strings"
	"path/filepath"
	"io"
	"errors"
	"github.com/joostlawerman/cache"
)

type Driver struct {
	Location string
	meta     *Meta
}

func NewDriver(location string) (*Driver, error) {
	location = strings.TrimSuffix(location, "/")

	_, err := os.Stat(location)
	if os.IsNotExist(err) {
		os.MkdirAll(location, 0777)
	}

	meta, err := NewMeta(location + "/")
	if err != nil {
		return nil, err
	}
	return &Driver{
		Location: location + "/",
		meta: meta,
	}, nil
}

func (d *Driver) Meta() cache.Meta {
	return d.meta
}

func (d *Driver) Put(key string, contents interface{}) error {
	_, err := os.Stat(filepath.Dir(key))
	if os.IsNotExist(err) {
		os.Mkdir(d.Location + filepath.Dir(key), 0666)
	}

	reader, ok := contents.(io.Reader)
	if !ok {
		return errors.New("Interface should be a reader")
	}

	file, err := os.Create(d.Location + key)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (d *Driver) Get(key string) (interface{}, error) {
	return os.Open(d.Location + key)
}

func (d *Driver) Exists(key string) bool {
	_, err := os.Stat(d.Location + key)

	return !os.IsNotExist(err)
}

func (d *Driver) Remove(key string) error {
	return os.Remove(d.Location + key)
}
