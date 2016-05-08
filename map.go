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
	return c
}

func (c *MapCache) SetLoader(keyLoader func(interface{}) interface{}) {
	c.loader = keyLoader
}

func (c *MapCache) Get(key interface{}) (interface{}, bool) {
	return c.get(key, time.Now().UTC())
}

func (c *MapCache) get(key interface{}, now time.Time) (interface{}, bool) {
	if value, exists := c.data[key]; exists && !value.Stale(now, c.ttl) {
		return value.Value, true
	}

	// Cache miss or data older than ttl
	c.sync(key)

	if value, exists := c.data[key]; exists {
		return value.Value, true
	} else {
		return nil, false
	}
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