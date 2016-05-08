package kash

import (
	"time"
)

type MapCache struct {
	data                 map[interface{}]*element
	intervalCacheRefresh time.Duration
	loadFunc             func(interface{}) interface{}
	refreshInterval      time.Duration
	refreshFunc          func(key interface{}, prev interface{}) interface{}
	ttlAccessed          time.Duration
	ttlWrite             time.Duration
}

func NewMapCache() *MapCache {
	nilSliceLoader := func(interface{}) interface{} {
		return nil
	}

	c := &MapCache{
		data: make(map[interface{}]*element),
		intervalCacheRefresh: MaxDuration,
		ttlWrite: MaxDuration,
		loadFunc: nilSliceLoader,
		ttlAccessed: MaxDuration,
	}
	return c
}

func (c *MapCache) SetLoader(keyLoader func(interface{}) interface{}) {
	c.loadFunc = keyLoader
}

func (c *MapCache) Get(key interface{}) (interface{}, bool) {
	return c.get(key, time.Now().UTC())
}

func (c *MapCache) GetIfPresent(key interface{}) (interface{}, bool) {
	return c.getIfPresent(key, time.Now().UTC())
}

func (c *MapCache) Invalidate(key interface{}) {
	delete(c.data, key)
}

func (c *MapCache) Put(key interface{}, value interface{}) {
	c.data[key] = newElement(value)
}

func (c *MapCache) PutAll(values map[interface{}]interface{}) {
	for k, v := range values {
		c.data[k] = newElement(v)
	}
}

func (c *MapCache) RefreshAfterWrite(ttl time.Duration) {
	c.intervalCacheRefresh = ttl
}

func (c *MapCache) ExpireAfterAccess(ttl time.Duration) {
	c.ttlAccessed = ttl
}

func (c *MapCache) ExpireAfterWrite(ttl time.Duration) {
	c.ttlWrite = ttl
}

func (c *MapCache) Len() int {
	return len(c.data)
}

func (c *MapCache) get(key interface{}, now time.Time) (interface{}, bool) {
	if value, exists := c.data[key]; exists && !value.WriteStale(now, c.ttlWrite) {
		defer c.updateAccessTime(key, now)
		return value.Value, true
	}

	// Cache miss or data older than ttl
	c.sync(key)

	if value, exists := c.data[key]; exists {
		defer c.updateAccessTime(key, now)
		return value.Value, true
	} else {
		return nil, false
	}
}

func (c *MapCache) getIfPresent(key interface{}, now time.Time) (interface{}, bool) {
	if value, exists := c.data[key]; exists && !value.WriteStale(now, c.ttlWrite) {
		return value.Value, true
	} else {
		return nil, false
	}
}

//func (c *MapCache) launchCacheRefresher() {
//	go func() {
//		for range time.Tick(c.refreshInterval) {
//			now := time.Now().UTC()
//			c.refreshKeys(now)
//		}
//	}()
//}

func (c *MapCache) sync(key interface{}) {
	if value := c.loadFunc(key); value != nil {
		c.data[key] = newElement(value)
	}
}

func (c *MapCache) updateAccessTime(key interface{}, now time.Time) {
	if value, exists := c.data[key]; exists {
		value.AccessedAt = now
		c.data[key] = value
	}
}
