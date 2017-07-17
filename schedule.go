package cache

import "time"

// schedule schedules a closure for a certain interval
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