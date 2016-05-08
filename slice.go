package gache

import "time"

type SliceCache struct {
	data     []*Element
	loader   func() []*Element
	duration time.Duration
}

func NewSliceCache(loader func() []*Element, duration time.Duration) *SliceCache {
	c := &SliceCache{duration: duration}
	c.SetLoader(loader)
	return c
}

func (c *SliceCache) SetLoader(loader func() []*Element) {
	if c.loader == nil {
		c.loader = loader
		// Launch background loader
		go func() {
			for range time.Tick(c.duration) {
				c.sync()
			}
		}()
	}
}

func (c *SliceCache) SetDuration(duration time.Duration) {
	c.duration = duration
}

func (c *SliceCache) sync() {
	c.data = c.loader()
}