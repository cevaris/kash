package kash

import (
	"time"
)

/*
 Cache Notes
 - http://stackoverflow.com/questions/10153724/google-guavas-cacheloader-loadall-vs-reload-semantics
  */

type MapCache struct {
	data        map[interface{}]*element
	loadFunc    func(key interface{}) interface{}
	reloadFunc  func(key interface{}, prev interface{}) interface{}
	ttlAccessed time.Duration
	ttlRefresh  time.Duration
	ttlWrite    time.Duration
}

func NewMapCache() *MapCache {
	nilLoadFunc := func(key interface{}) interface{} {
		return nil
	}
	nilReloadFunc := func(key interface{}, prev interface{}) interface{} {
		return nil
	}
	c := &MapCache{
		data: make(map[interface{}]*element),
		loadFunc: nilLoadFunc,
		ttlWrite: MaxDuration,
		ttlAccessed: MaxDuration,
		ttlRefresh: MaxDuration,
		reloadFunc: nilReloadFunc,
	}
	return c
}

func (c *MapCache) SetLoader(keyLoader func(key interface{}) interface{}) {
	c.loadFunc = keyLoader
	c.reloadFunc = func(key interface{}, prev interface{}) interface{} {
		return keyLoader(key)
	}
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
	c.ttlRefresh = ttl
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
	if value, exists := c.data[key]; exists {

		if value.WriteStale(now, c.ttlRefresh) {
			c.Refresh(key)
		}

		defer c.updateAccessTime(key, now)
		return value.Value, true

	} else if value := c.loadFunc(key); value != nil {
		c.data[key] = newElement(value)
		return c.data[key].Value, true
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

func (c *MapCache) Refresh(key interface{}) {
	if value, exists := c.data[key]; exists {
		c.data[key] = newElement(c.reloadFunc(key, value))
	} else if value := c.loadFunc(key); value != nil {
		c.data[key] = newElement(value)
	}
}

func (c *MapCache) updateAccessTime(key interface{}, now time.Time) {
	if value, exists := c.data[key]; exists {
		value.AccessedAt = now
		c.data[key] = value
	}
}
