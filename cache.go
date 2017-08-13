package cache

import (
	"errors"
	"time"
)

const ErrorItemIsExpired = "Item is expired"

// When the duration is zero the items in the cache will never expire.
// If the interval is zero the cache will not automatically remove
// expired items from the cache.
type Cache struct {
	duration time.Duration
	interval time.Duration
	close    chan struct{}
	driver   Driver
}

// New creates a new *Cache
func New(driver Driver, duration time.Duration) *Cache {
	return &Cache{
		duration: duration,
		driver:   driver,
	}
}

// Open creates a new *Cache and triggers the schedule method
func Open(driver Driver, duration, interval time.Duration) *Cache {
	cache := New(driver, duration)

	if interval > 0 {
		cache.Schedule(interval)
	}

	return cache
}

// Schedule schedules the cache clean method for with the interval
func (c *Cache) Schedule(interval time.Duration) {
	c.close = schedule(interval, func() {
		c.Clean()
	})
}

// Clean removes all expired items from the cache
func (c *Cache) Clean() {
	for key := range c.driver.Meta().List() {
		if c.driver.Meta().IsExpired(key, c.duration) {
			c.Remove(key)
		}
	}
}

// Put puts a item in the cache
func (c *Cache) Put(key string, contents interface{}) error {
	// Put into meta
	if err := c.driver.Meta().Put(key); err != nil {
		return err
	}
	return c.driver.Put(key, contents)
}

// Get retrieves a item from the cache it will return a "ErrorItemIsExpired"
// when the item is expired
func (c *Cache) Get(key string) (interface{}, error) {
	// Check if item is expired
	if c.IsExpired(key) {
		return nil, errors.New(ErrorItemIsExpired)
	}
	return c.driver.Get(key)
}

// IsExpired checks if the item is expired
func (c *Cache) IsExpired(key string) bool {
	// Do not expire when duration is 0
	if c.duration == 0 {
		return false
	}
	return c.driver.Meta().IsExpired(key, c.duration)
}

// Exists checks if an item exist within the cache
func (c *Cache) Exists(key string) bool {
	return c.driver.Exists(key)
}

// Remove removes an item from the cache
func (c *Cache) Remove(key string) error {
	c.driver.Meta().Remove(key)

	return c.driver.Remove(key)
}

// Remember checks if a item exists or is expires if it does not exist
// it will use the callback to set the value
func (c *Cache) Remember(key string, callback func() (interface{}, error)) (interface{}, error) {
	if !c.Exists(key) || c.IsExpired(key) {
		// If it does not exist retrieve value
		value, err := callback()
		if err != nil {
			return nil, err
		}
		if err := c.Put(key, value); err != nil {
			return nil, err
		}
	}
	return c.Get(key)
}

// Close stop the scheduled clean method
func (c *Cache) Close() {
	close(c.close)
}
