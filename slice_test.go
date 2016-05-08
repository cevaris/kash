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

	actualLen := len(sliceCache.Get())
	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", actualLen)
	}
}

func TestSliceCacheString(t *testing.T) {
	testData := []*Element{NewElement("alpha"), NewElement("bravo")}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	actualLen := len(sliceCache.Get())
	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", actualLen)
	}
}

func TestSliceCacheGenericInterface(t *testing.T) {
	testData := []*Element{NewElement(&TestObject{}), NewElement(&TestObject{})}
	sliceCache := NewSliceCache()
	sliceCache.SetLoader(buildSliceLoader(testData))

	actualLen := len(sliceCache.Get())
	if ! reflect.DeepEqual(sliceCache.Get(), testData) {
		t.Error("invalid cache data", actualLen)
	}
}

func buildSliceLoader(values []*Element) func() []*Element {
	return func() []*Element {
		return values
	}
}