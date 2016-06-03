package kash

import (
	"time"
	"sync"
)

/*
 Cache Notes
 - http://stackoverflow.com/questions/10153724/google-guavas-cacheloader-loadall-vs-reload-semantics
 - http://stackoverflow.com/questions/31983275/how-to-change-external-variables-value-inside-a-goroutine-closure
 - http://google.github.io/guava/releases/snapshot/api/docs/com/google/common/cache/CacheBuilder.html
 - http://google.github.io/guava/releases/snapshot/api/docs/com/google/common/cache/LoadingCache.html
 - https://github.com/google/guava/wiki/CachesExplained
  */

type KeyValue  struct {
	Key   interface{}
	Value interface{}
}

type MapCache struct {
	asyncReloadFunc func(key interface{}, prev interface{})
	data            map[interface{}]*element
	loadFunc        func(key interface{}) interface{}
	sizePruneLock   sync.Mutex
	reloadFunc      func(key interface{}, prev interface{}) interface{}
	ReloadChan      chan KeyValue
	ttlAccessed     time.Duration
	ttlRefresh      time.Duration
	ttlWrite        time.Duration
	weigherFunc     func(key interface{}, value interface{}) int64
	weightCurr      int64
	weightMax       int64
}

func NewMapCache() *MapCache {
	nilLoadFunc := func(key interface{}) interface{} {
		return nil
	}
	c := &MapCache{
		data: make(map[interface{}]*element),
		loadFunc: nilLoadFunc,
		sizePruneLock: sync.Mutex{},
		ttlWrite: MaxDuration,
		ttlAccessed: MaxDuration,
		ttlRefresh: MaxDuration,
		ReloadChan: make(chan KeyValue),
		weightCurr: 0,
	}

	c.launchReloadListener()

	c.SetReload(func(key interface{}, prev interface{}) interface{} {
		return c.loadFunc(key)
	})

	c.SetWeigher(func(key interface{}, prev interface{}) int64 {
		return 1
	})

	return c
}

func (c *MapCache) SetLoad(loadFunc func(key interface{}) interface{}) {
	c.loadFunc = loadFunc
}

func (c *MapCache) SetReload(reloadFunc func(key interface{}, prev interface{}) interface{}) {
	c.reloadFunc = reloadFunc
}

func (c *MapCache) SetAsyncReload(reloadAsyncFunc func(key interface{}, prev interface{})) {
	c.asyncReloadFunc = reloadAsyncFunc
}

func (c *MapCache) SetWeigher(weigherFunc func(key interface{}, value interface{}) int64) {
	c.weigherFunc = weigherFunc
}

func (c *MapCache) SetMaxWeight(weight int64) {
	c.weightMax = weight
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
	if value != nil {
		c.data[key] = newElementWithWeight(value, c.weigherFunc(key, value))
	}
}

func (c *MapCache) PutAll(values map[interface{}]interface{}) {
	for k, v := range values {
		c.Put(k, v)
	}
}

func (c *MapCache) Refresh(key interface{}) {
	if elem, exists := c.data[key]; exists {
		if c.asyncReloadFunc != nil {
			c.asyncReloadFunc(key, elem.Value)
		} else {
			c.Put(key, c.reloadFunc(key, elem.Value))
		}
	} else if value := c.loadFunc(key); value != nil {
		c.Put(key, value)
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
	if elem, exists := c.data[key]; exists {

		if elem.AccessStale(now, c.ttlAccessed) {
			c.Invalidate(key)
			return nil, false
		}

		if elem.WriteStale(now, c.ttlWrite) {
			c.Invalidate(key)
			return nil, false
		}

		if elem.WriteStale(now, c.ttlRefresh) {
			c.Refresh(key)

			defer c.updateAccessTime(key, now)
			return c.data[key].Value, true
		}

		defer c.updateAccessTime(key, now)
		return elem.Value, true
	} else if value := c.loadFunc(key); value != nil {
		c.Put(key, value)
		return c.data[key].Value, true
	} else {
		return nil, false
	}
}

func (c *MapCache) getIfPresent(key interface{}, now time.Time) (interface{}, bool) {
	if elem, exists := c.data[key]; exists {
		if !elem.WriteStale(now, c.ttlRefresh) {
			return elem.Value, true
		}
	}
	return nil, false
}

	func (c *MapCache) updateAccessTime(key interface{}, now time.Time) {
		if elem, exists := c.data[key]; exists {
			elem.AccessedAt = now
			c.data[key] = elem
		}
	}

func (c *MapCache) launchReloadListener() {
	go func() {
		for {
			select {
			case kv := <-c.ReloadChan:
				c.Put(kv.Key, kv.Value)
			}
		}
	}()
}
