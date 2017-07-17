package file_test

import (
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
	"github.com/joostlawerman/cache/drivers/file"
	"github.com/joostlawerman/cache"
)

var staticStash *cache.Cache

func init() {
	staticStash, _ = cache.Open(file.Driver(".static-cache", 0777), time.Second * 1, 0)
}

func TestFileMeta_IsExpired(t *testing.T) {
	for name, contents := range testPutCases {
		assert.NoError(t, staticStash.Put(name, contents))

		assert.False(t, staticStash.IsExpired(name))
	}

	time.Sleep(2 * time.Second)

	for name := range testPutCases {
		assert.True(t, staticStash.IsExpired(name))
	}
}