package gache

import "time"

const DefaultSliceCacheSize int = 1000

type SliceCache struct {
	Cache
	Data     []interface{}
	loader   func() []interface{}
	duration time.Duration
}

func NewSliceCache(loader func() []interface{}, duration time.Duration) *SliceCache {
	c := &SliceCache{
		loader: loader,
		duration: duration,
	}

	// Launch background loader
	go func() {
		for range time.Tick(c.duration) {
			c.Sync()
		}
	}()

	//Force sync now
	c.Sync()

	return c
}

func (c *SliceCache) Sync() {
	c.Data = c.loader()
}