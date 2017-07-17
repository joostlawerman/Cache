package cache

import "time"

// Meta contains the time a item has been created
type Meta interface {
	List() map[string]time.Time
	Put(key string) error
	IsExpired(key string, duration time.Duration) bool
	Remove(key string) error
}