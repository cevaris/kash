package gache

import (
	"testing"
	"reflect"
)

const TestCacheSize int = 10

type TestObject struct{}

func TestNewSliceCache(t *testing.T) {
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(make([]*Element, TestCacheSize)))

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
	testData := []*Element{NewElement(10), NewElement(100)}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", sliceCache.Get())
	}
}

func TestSliceCacheString(t *testing.T) {
	testData := []*Element{NewElement("alpha"), NewElement("bravo")}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", sliceCache.Get())
	}
}

func TestSliceCacheGenericInterface(t *testing.T) {
	testData := []*Element{NewElement(&TestObject{}), NewElement(&TestObject{})}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", sliceCache.Get())
	}
}

func TestSliceCacheLoader(t *testing.T) {
	testData := [][]*Element{
		[]*Element{NewElement("alpha"), NewElement("bravo")},
		[]*Element{NewElement("charile"), NewElement("delta")},
		[]*Element{NewElement("echo"), NewElement("foxtrot")},
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

func buildSliceLoader(values []*Element) func() []*Element {
	return func() []*Element {
		return values
	}
}

func buildSliceLoaderIter(values [][]*Element) func() []*Element {
	var i = 0
	return func() []*Element {
		for i < len(values) {
			p := i
			i = i + 1
			return values[p]
		}
		return nil
	}
}