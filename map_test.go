package kash

import (
	"testing"
	"fmt"
	"time"
)

func TestNewMapCache(t *testing.T) {
	testValue := 10
	mapCache := NewMapCache()
	mapCache.SetLoad(staticMapLoader(testValue))
	if actual, exists := mapCache.Get("test"); !exists || actual != testValue {
		t.Error("invalid cache data", actual, exists)
	}
}

func TestMapCacheKeyLoader(t *testing.T) {
	testValueA := 10
	testValueB := 20
	mapCache := NewMapCache()
	mapCache.SetLoad(mapKeyLoaderIter(
		[]map[interface{}]interface{}{
			map[interface{}]interface{}{},
			map[interface{}]interface{}{"test": testValueA},
			map[interface{}]interface{}{"test": testValueB},
		},
	))

	if actual, exists := mapCache.Get("test"); exists || actual != nil {
		t.Error("not nil", actual, exists)
	}

	mapCache.Refresh("test")

	if actual, exists := mapCache.Get("test"); !exists || actual != testValueA {
		t.Error("invalid cache data", actual, exists)
	}

	mapCache.Refresh("test")

	if actual, exists := mapCache.Get("test"); !exists || actual != testValueB {
		t.Error("invalid cache data", actual)
	}
}

func TestNewMapCacheGetIfPresent(t *testing.T) {
	testValue := 10
	testKey := "test"
	mapCache := NewMapCache()
	mapCache.SetLoad(staticMapLoader(testValue))

	if actual, exists := mapCache.Get(testKey); !exists || actual != testValue {
		t.Error("invalid cache data", actual, exists)
	}

	mapCache.Invalidate(testKey)

	if actual, exists := mapCache.GetIfPresent(testKey); exists || actual != nil {
		t.Error("not nil", actual, exists)
	}

	if actual, exists := mapCache.Get(testKey); !exists || actual != testValue {
		t.Error("invalid cache data", actual, exists)
	}

}

func TestMapCacheInvalidate(t *testing.T) {
	testMap := map[interface{}]interface{}{"a": 1, "b": 2}
	mapCache := NewMapCache()
	mapCache.PutAll(testMap)

	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}

	mapCache.Invalidate("a")

	if actual, exists := mapCache.Get("a"); exists || actual != nil {
		t.Error("not nil", actual, exists)
	}
}

func TestMapCacheLen(t *testing.T) {
	testMap := map[interface{}]interface{}{"a": 1, "b": 2}
	mapCache := NewMapCache()

	if actual := mapCache.Len(); actual != 0 {
		t.Error("invalid cache len", actual)
	}

	mapCache.PutAll(testMap)

	if actual := mapCache.Len(); actual != len(testMap) {
		t.Error("invalid cache len", actual)
	}
}

func TestMapCachePut(t *testing.T) {
	testKey := "test"
	testValue := 10
	mapCache := NewMapCache()

	if actual, exists := mapCache.Get(testKey); exists || actual != nil {
		t.Error("not nil", actual, exists)
	}

	mapCache.Put(testKey, testValue)

	if actual, exists := mapCache.Get("test"); !exists || actual != testValue {
		t.Error("invalid cache data", actual, exists)
	}
}

func TestMapCachePutNil(t *testing.T) {
	testKey := "test"
	mapCache := NewMapCache()

	if actual, exists := mapCache.Get(testKey); exists || actual != nil {
		t.Error("not nil", actual, exists)
	}

	mapCache.Put(testKey, nil)

	if actual, exists := mapCache.Get(testKey); exists || actual != nil {
		t.Error("not nil", actual, exists)
	}
}

func TestMapCachePutAll(t *testing.T) {
	testMap := map[interface{}]interface{}{"a": 1, "b": 2}
	mapCache := NewMapCache()

	if actual, exists := mapCache.Get("a"); exists || actual != nil {
		t.Error("not nil", actual, exists)
	}

	mapCache.PutAll(testMap)

	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
}

func TestMapCacheRefreshAfterWriteAsync(t *testing.T) {
	testMap := map[interface{}]interface{}{"a": 1}
	mapCache := NewMapCache()
	mapCache.RefreshAfterWrite(5 * time.Millisecond)
	mapCache.SetLoad(func(key interface{}) interface{} {
		someExternalSource := map[interface{}]interface{}{
			"a": 1,
		}
		return someExternalSource[key]
	})
	mapCache.SetAsyncReload(func(key interface{}, prev interface{}) {
		someExternalSource := map[interface{}]interface{}{
			"a": time.Now().Unix(),
		}
		go func() {
			time.Sleep(5 * time.Millisecond)
			mapCache.ReloadChan <- KeyValue{Key: key, Value: someExternalSource[key]}
		}()
	})


	// Cache.Load
	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
	time.Sleep(10 * time.Millisecond)

	// Cache.Reload (kicks of async job, returns old values)
	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
	time.Sleep(10 * time.Millisecond)

	// Cache.Reload (returns new async loaded returned values)
	if actual, exists := mapCache.Get("a"); !exists || actual == testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
}

func TestMapCacheRefreshAfterWriteBlocking(t *testing.T) {
	testMap := map[interface{}]interface{}{"a": 1}
	mapCache := NewMapCache()
	mapCache.RefreshAfterWrite(1 * time.Millisecond)
	mapCache.SetLoad(func(key interface{}) interface{} {
		someExternalSource := map[interface{}]interface{}{
			"a": 1,
		}
		return someExternalSource[key]
	})
	mapCache.SetReload(func(key interface{}, prev interface{}) interface{} {
		someExternalSource := map[interface{}]interface{}{
			"a": time.Now().Unix(),
		}
		time.Sleep(5 * time.Millisecond)
		return someExternalSource[key]
	})


	// Cache.Load
	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
	time.Sleep(10 * time.Millisecond)

	// Cache.Reload (blocks and returns new value)
	if actual, exists := mapCache.Get("a"); !exists || actual == testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
}

func TestMapCacheExpireAfterAccess(t *testing.T) {
	testMap := map[interface{}]interface{}{"a": 1}
	mapCache := NewMapCache()
	mapCache.ExpireAfterAccess(5 * time.Millisecond)
	mapCache.PutAll(testMap)

	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
	time.Sleep(10 * time.Millisecond)

	if actual, exists := mapCache.Get("a"); exists {
		t.Error("invalid cache data", actual, exists)
	}

	mapCache.Put("a", testMap["a"])

	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
	time.Sleep(10 * time.Millisecond)

	if actual, exists := mapCache.Get("a"); exists {
		t.Error("invalid cache data", actual, exists)
	}

	mapCache.Put("a", testMap["a"])

	// Keep cache and keeping "a" alive
	for _, _ = range [10]struct{}{} {
		mapCache.Get("a")
		time.Sleep(2 * time.Millisecond)
	}

	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}

}

func TestMapCacheExpireAfterWrite(t *testing.T) {
	testMap := map[interface{}]interface{}{"a": 1}
	mapCache := NewMapCache()
	mapCache.ExpireAfterWrite(1 * time.Millisecond)

	mapCache.PutAll(testMap)
	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
	time.Sleep(5 * time.Millisecond)
	if actual, exists := mapCache.Get("a"); exists {
		t.Error("invalid cache data", actual, exists)
	}

	mapCache.Put("a", testMap["a"])
	if actual, exists := mapCache.Get("a"); !exists || actual != testMap["a"] {
		t.Error("invalid cache data", actual, exists)
	}
	time.Sleep(5 * time.Millisecond)
	if actual, exists := mapCache.Get("a"); exists {
		t.Error("invalid cache data", actual, exists)
	}
}

func staticMapLoader(value interface{}) func(interface{}) interface{} {
	return func(interface{}) interface{} {
		return value
	}
}

func mapKeyLoaderIter(values []map[interface{}]interface{}) func(interface{}) interface{} {
	var i = 0
	return func(key interface{}) interface{} {
		for i < len(values) {
			curr := i
			i++
			currData := values[curr]
			return currData[key]
		}
		return nil
	}
}

func dumpMap(s map[interface{}]*element) {
	builder := make([]string, len(s))

	var i = 0
	for k, v := range s {
		builder[i] = fmt.Sprintf("%+v:%+v, ", k, v.Value)
		i++
	}
	println(fmt.Sprintf("{%v}", builder))
}
