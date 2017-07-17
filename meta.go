package cache

import "time"

type Meta interface {
	Boot() (error)

	List() map[string]time.Time
	Put(key string) error
	IsExpired(key string, duration time.Duration) bool
	Remove(key string) error
}