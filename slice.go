package gache

import (
	"time"
)

type SliceCache struct {
	data     []*Element
	loader   func() []*Element
	duration time.Duration
}

func NewSliceCache() *SliceCache {
	var nilSliceLoader = func() []*Element {
		return nil
	}

	c := &SliceCache{
		loader: nilSliceLoader,
		duration: MaxDuration,
	}
	c.launchLoader()
	return c
}

func (c *SliceCache) SetLoader(loader func() []*Element) {
	c.loader = loader
}

func (c *SliceCache) RefreshAfterWrite(duration time.Duration) {
	c.duration = duration
}

func (c *SliceCache) Get() []*Element {
	if c.data == nil {
		c.sync()
	}
	return c.data
}

func (c *SliceCache) sync() {
	c.data = c.loader()
}

func (c *SliceCache) launchLoader() {
	// should only launch once
	go func() {
		for range time.Tick(c.duration) {
			c.sync()
		}
	}()
}

