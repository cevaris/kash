package gache

import (
	"testing"
	"fmt"
)

func TestNewMapCache(t *testing.T) {
	testValue := 10
	mapCache := NewMapCache()
	mapCache.SetLoader(staticMapLoader(testValue))
	if actual, ok := mapCache.Get("test"); !ok || actual != testValue {
		t.Error("invalide cache data", actual, ok)
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

	if actual, ok := mapCache.Get("test"); ok || actual != nil {
		t.Error("not nil", actual, ok)
	}

	mapCache.sync("test")

	if actual, ok := mapCache.Get("test"); !ok || actual != testValueA {
		t.Error("invalide cache data", actual, ok)
	}

	mapCache.sync("test")

	if actual, ok := mapCache.Get("test"); !ok || actual != testValueB {
		t.Error("invalide cache data", actual)
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

func dumpMap(s map[interface{}]*Element) {
	builder := make([]string, len(s))

	var i = 0
	for k, v := range s {
		builder[i] = fmt.Sprintf("%+v:%+v, ", k, v)
		i++
	}
	println(fmt.Sprintf("{%v}", builder))
}
