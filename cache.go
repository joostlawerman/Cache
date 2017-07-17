package cache

import (
	"time"
	"errors"
)

const ErrorItemIsExpired = "Item is expired"

type Cache struct {
	duration time.Duration
	interval time.Duration
	close chan struct{}
	driver Driver
}

func (c *Cache) boot() error {
	if err := c.driver.Meta().Boot(); err != nil {
		return err
	}

	if c.interval > 0 {
		c.close = schedule(c.interval, func() {
			meta := c.driver.Meta().List()

			for name, createdAt := range meta {
				if time.Now().After(createdAt.Add(c.duration)) {
					c.Remove(name)
				}
			}
		})
	}
	return nil
}

func (c *Cache) Put(key string, contents interface{}) error {
	if err := c.driver.Meta().Put(key); err != nil {
		return err
	}
	return c.driver.Put(key, contents)
}

func (c *Cache) Get(key string) (interface{}, error) {
	if c.IsExpired(key) {
		return nil, errors.New(ErrorItemIsExpired)
	}
	return c.driver.Get(key)
}

func (c *Cache) IsExpired(key string) bool {
	if c.duration == 0 {
		return false
	}
	return c.driver.Meta().IsExpired(key, c.duration)
}

func (c *Cache) Exists(key string) bool {
	return c.driver.Exists(key)
}

func (c *Cache) Remove(key string) error {
	c.driver.Meta().Remove(key)

	return c.driver.Remove(key)
}

func (c *Cache) Remember(key string, callback func() (interface{}, error)) (interface{}, error) {
	if !c.Exists(key) || c.IsExpired(key) {
		contents, err := callback()
		if err != nil {
			return nil, err
		}
		if err := c.Put(key, contents); err != nil {
			return nil, err
		}
	}
	return c.Get(key)
}

func (c *Cache) Close() {
	close(c.close)
}

func New(driver Driver, duration, interval time.Duration) *Cache {
	return &Cache{
		duration: duration,
		interval: interval,
		driver: driver,
	}
}

func Open(driver Driver, duration, interval time.Duration) (*Cache, error) {
	cache := New(driver, duration, interval)

	if err := cache.boot(); err != nil {
		return nil, err
	}
	return cache, nil
}

func schedule(interval time.Duration, closure func()) chan struct{} {
	ticker := time.NewTicker(interval)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <- ticker.C:
				closure()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}


// File
//
//func init() {
//	if ok, err := exists(CacheLoc); !ok && err == nil {
//		os.Mkdir(CacheLoc, 0777)
//	} else {
//		exitErr(err)
//	}
//
//	history = make(map[string]time.Time)
//}
//
//func historySet(key string) {
//	history[key] = time.Now().Add(5 * time.Minute)
//}
//
//func historyCheck(key string) bool {
//	if expire, ok := history[key]; ok {
//		return expire.After(time.Now())
//	}
//	return false
//}
//
//func historyUpdate() {
//	for key, expire := range history {
//		if expire.Before(time.Now()) {
//			os.Remove(cacheLocation(key))
//		}
//	}
//}
//
//func cacheOr(key string, callback func() []byte) []byte {
//	if data, err := cacheGet(key); err == nil {
//		return data
//	}
//
//	data := callback()
//
//	go func() {
//		err := cacheSet(key, data)
//
//		exitErr(err)
//	}()
//
//	return data
//}
//
//func cacheGet(key string) ([]byte, error) {
//	if historyCheck(key) {
//		data, err := ioutil.ReadFile(cacheLocation(key))
//		if os.IsNotExist(err) {
//			return data, errors.New("File does not exists")
//		}
//		if err != nil {
//			return data, err
//		}
//
//		if err != nil {
//			return data, err
//		}
//		return data, nil
//	}
//
//	go historyUpdate()
//
//	return nil, errors.New("Cache file expired")
//}
//
//func cacheSet(key string, data []byte) error {
//	historySet(key)
//
//	return cacheWrite(cacheLocation(key), data)
//}
//
//func cacheWrite(location string, data []byte) error {
//	if ok, err := exists(location); !ok && err == nil {
//		os.Remove(location)
//	} else {
//		exitErr(err)
//	}
//
//	return ioutil.WriteFile(location, data, 0777)
//}
//
//func exists(path string) (bool, error) {
//	_, err := os.Stat(path)
//
//	if err == nil {
//		return true, nil
//	}
//
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//
//	return true, err
//}
//
//func cacheLocation(key string) string {
//	return "./.cache/" + key + ".json"
//}
//
