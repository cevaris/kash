package kash

import (
	"testing"
	"reflect"
	"fmt"
)

const TestCacheSize int = 10

type TestObject struct{}

func TestNewSliceCache(t *testing.T) {
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(make([]interface{}, TestCacheSize)))

	actualLen := len(sliceCache.Get())
	if len(sliceCache.Get()) != TestCacheSize {
		t.Error("invalid default cache size", actualLen)
	}
}

func TestNilSliceCacheLoader(t *testing.T) {
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(nil))

	if sliceCache.Get() != nil {
		t.Error("cache not nil")
	}
}

func TestSliceCacheInt(t *testing.T) {
	testData := []interface{}{newElement(10), newElement(100)}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", sliceCache.Get())
	}
}

func TestSliceCacheString(t *testing.T) {
	testData := []interface{}{newElement("alpha"), newElement("bravo")}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", sliceCache.Get())
	}
}

func TestSliceCacheGenericInterface(t *testing.T) {
	testData := []interface{}{newElement(&TestObject{}), newElement(&TestObject{})}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", sliceCache.Get())
	}
}

func TestSliceCacheLoader(t *testing.T) {
	testData := [][]interface{}{
		[]interface{}{newElement("alpha"), newElement("bravo")},
		[]interface{}{newElement("charile"), newElement("delta")},
		[]interface{}{newElement("echo"), newElement("foxtrot")},
	}

	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoaderIter(testData))

	if ! reflect.DeepEqual(sliceCache.Get(), testData[0]) {
		t.Error("invalid cache data", sliceCache)
	}

	sliceCache.sync()

	if ! reflect.DeepEqual(sliceCache.Get(), testData[1]) {
		t.Error("invalid cache data", sliceCache)
	}

	sliceCache.sync()

	if ! reflect.DeepEqual(sliceCache.Get(), testData[2]) {
		t.Error("invalid cache data", sliceCache)
	}
}

func buildSliceLoader(values []interface{}) func() []interface{} {
	return func() []interface{} {
		return values
	}
}

func buildSliceLoaderIter(values [][]interface{}) func() []interface{} {
	var i = 0
	return func() []interface{} {
		for i < len(values) {
			p := i
			i = i + 1
			return values[p]
		}
		return nil
	}
}

func dumpSlice(s []interface{}) {
	builder := make([]string, len(s))

	for i, v := range s {
		builder[i] = fmt.Sprintf("%+v, ", v)
	}
	println(fmt.Sprintf("%v", builder))
}