package cache

// Driver contains all the methods of managing the cache
type Driver interface {
	Meta() Meta

	Put(name string, contents interface{}) error
	Get(name string) (interface{}, error)
	Exists(name string) bool
	Remove(name string) error
}
