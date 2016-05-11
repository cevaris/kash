package kash

import (
	"time"
	"fmt"
)

type SliceCache struct {
	data   []interface{}
	loader func() []interface{}
	ttl    time.Duration
}

func NewSliceCache() *SliceCache {
	nilSliceLoader := func() []interface{} {
		return nil
	}

	c := &SliceCache{
		loader: nilSliceLoader,
		ttl: MaxDuration,
	}
	c.launchLoader()
	return c
}

func (c *SliceCache) SetLoader(loader func() []interface{}) {
	c.loader = loader
}

func (c *SliceCache) SetCacheTtl(duration time.Duration) {
	c.ttl = duration
}

func (c *SliceCache) Get() []interface{} {
	if c.data == nil {
		c.sync()
	}
	return c.data
}

func (c *SliceCache) String() string {
	builder := make([]string, len(c.data))

	for i, v := range c.data {
		builder[i] = fmt.Sprintf("%+v, ", v)
	}

	return fmt.Sprintf(
		"SliceCache(%+v,%+v,%+v)",
		builder,
		c.loader,
		c.ttl,
	)
}

func (c *SliceCache) sync() {
	c.data = c.loader()
}

func (c *SliceCache) launchLoader() {
	// should only launch once
	go func() {
		for _ = range time.Tick(c.ttl) {
			c.sync()
		}
	}()
}
