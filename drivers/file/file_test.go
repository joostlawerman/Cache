package file_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"strings"
	"io"
	"github.com/joostlawerman/cache/drivers/file"
	"github.com/joostlawerman/cache"
)

var expiringStash *cache.Cache

func init() {
	expiringStash, _ = cache.Open(file.Driver(".expiring-cache", 0777), time.Second * 2, time.Second * 2)
}


var testPutCases = map[string]io.Reader{
	"lorum": strings.NewReader("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. Nam eget dui."),
}

func TestFile_Put(t *testing.T) {
	for name, contents := range testPutCases {
		assert.NoError(t, expiringStash.Put(name, contents))
	}

	time.Sleep(5 * time.Second)

	for name := range testPutCases {
		assert.False(t, expiringStash.Exists(name))
	}
}

//func TestFile_Remove(t *testing.T) {
//	for name := range testPutCases {
//		assert.NoError(t, stash.Remove(name))
//	}
//}
