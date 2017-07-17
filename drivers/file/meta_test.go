package file_test

import (
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
	"github.com/joostlawerman/cache/drivers/file"
	"github.com/joostlawerman/cache"
	"strings"
)

var staticStash *cache.Cache

func init() {
	driver, err := file.NewDriver("../../.test/.static-cache")
	if err != nil {
		panic(err)
	}

	staticStash = cache.New(driver, time.Second * 1)
}

func TestFileMeta_IsExpired(t *testing.T) {
	for name, contents := range testPutCases {
		assert.NoError(t, staticStash.Put(name, strings.NewReader(contents)))

		assert.False(t, staticStash.IsExpired(name))
	}

	time.Sleep(2 * time.Second)

	for name := range testPutCases {
		assert.True(t, staticStash.IsExpired(name))
	}
}