package kash

import (
	"testing"
	"fmt"
)

func TestNewMapCache(t *testing.T) {
	testValue := 10
	mapCache := NewMapCache()
	mapCache.SetLoader(staticMapLoader(testValue))
	if actual, exists := mapCache.Get("test"); !exists || actual != testValue {
		t.Error("invalid cache data", actual, exists)
	}
}

func TestMapCacheKeyLoader(t *testing.T) {
	testValueA := 10
	testValueB := 20
	mapCache := NewMapCache()
	mapCache.SetLoader(mapKeyLoaderIter(
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
	mapCache.SetLoader(staticMapLoader(testValue))

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
		builder[i] = fmt.Sprintf("%+v:%+v, ", k, v)
		i++
	}
	println(fmt.Sprintf("{%v}", builder))
}
