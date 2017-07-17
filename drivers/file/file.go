package file

import (
	"io"
	"os"
	"strings"
	"github.com/joostlawerman/cache"
	"errors"
	"path/filepath"
)

func Driver(location string, fileMode os.FileMode) cache.Driver {
	location = strings.TrimSuffix(location, "/")

	_, err := os.Stat(location)
	if os.IsNotExist(err) {
		os.Mkdir(location, fileMode)
	}
	return &File{
		Location: location + "/",
		FileMode: fileMode,
		FileMeta: &FileMeta{
			Location: location + "/",
			FileMode: fileMode,
		},
	}
}

type File struct {
	Location string
	FileMode os.FileMode
	FileMeta *FileMeta
}

func (m *File) Meta() cache.Meta {
	return m.FileMeta
}

func (m *File) Put(name string, contents interface{}) error {
	_, err := os.Stat(filepath.Dir(name))
	if os.IsNotExist(err) {
		os.Mkdir(m.Location + filepath.Dir(name), m.FileMode)
	}

	reader, ok := contents.(io.Reader)
	if !ok {
		return errors.New("Interface should be a reader")
	}

	file, err := os.Create(m.Location + name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (m *File) Get(name string) (interface{}, error) {
	return os.Open(m.Location + name)
}

func (m *File) Exists(name string) bool {
	_, err := os.Stat(m.Location + name)

	return !os.IsNotExist(err)
}

func (m *File) Remove(name string) error {
	return os.Remove(m.Location + name)
}
