package gache

import (
	"time"
)

type MapCache struct {
	data            map[interface{}]*Element
	ttl             time.Duration
	loader          func(interface{}) interface{}
	refreshInterval time.Duration
}

func NewMapCache() *MapCache {
	nilSliceLoader := func(interface{}) interface{} {
		return nil
	}

	c := &MapCache{
		data: make(map[interface{}]*Element),
		ttl: MaxDuration,
		loader: nilSliceLoader,
		refreshInterval: time.Second,
	}
	c.launchLoader()
	return c
}

func (c *MapCache) SetLoader(keyLoader func(interface{}) interface{}) {
	c.loader = keyLoader
}

func (c *MapCache) Get(key interface{}) (interface{}, bool) {
	if value, ok := c.data[key]; ok {
		return value.Value, true
	}

	// Cache miss, fetch new value for key
	c.sync(key)

	if value, ok := c.data[key]; ok {
		return value.Value, true
	}

	return nil, false
}

func (c *MapCache) RefreshAfterWrite(ttl time.Duration) {
	c.ttl = ttl
}

func (c *MapCache) refreshKeys(now time.Time) {
	for k, v := range c.data {
		if v.CreatedAt.Before(now.Add(c.ttl)) {
			c.sync(k)
		}
	}
}

func (c *MapCache) sync(key interface{}) {
	if value := c.loader(key); value != nil {
		c.data[key] = NewElement(value)
	}
}

func (c *MapCache) launchLoader() {
	// should only launch once
	go func() {
		for range time.Tick(c.refreshInterval) {
			now := time.Now().UTC()
			c.refreshKeys(now)
		}
	}()
}