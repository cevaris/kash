package gache

import (
	"testing"
	"time"
	"reflect"
)

const TestCacheSize int = 10

type TestObject struct{}

func TestNewSliceCache(t *testing.T) {
	sliceCache := NewSliceCache(
		buildSliceLoader(make([]*Element, TestCacheSize)),
		time.Second,
	)
	sliceCache.sync()

	actualLen := len(sliceCache.data)
	if len(sliceCache.data) != TestCacheSize {
		t.Error("invalid default cache size", actualLen)
	}
}

func TestNilSliceCacheLoader(t *testing.T) {
	sliceCache := NewSliceCache(buildSliceLoader(nil), time.Second)
	sliceCache.sync()

	if sliceCache.data != nil {
		t.Error("cache not nil")
	}
}

func TestSliceCacheInt(t *testing.T) {
	testData := []*Element{NewElement(10), NewElement(100)}
	sliceCache := NewSliceCache(
		buildSliceLoader(testData),
		time.Second,
	)
	sliceCache.sync()

	actualLen := len(sliceCache.data)
	if ! reflect.DeepEqual(sliceCache.data, testData) {
		t.Error("invalid cache data", actualLen)
	}
}

func TestSliceCacheString(t *testing.T) {
	testData := []*Element{NewElement("alpha"), NewElement("bravo")}
	sliceCache := NewSliceCache(
		buildSliceLoader(testData),
		time.Second,
	)
	sliceCache.sync()

	actualLen := len(sliceCache.data)
	if ! reflect.DeepEqual(sliceCache.data, testData) {
		t.Error("invalid cache data", actualLen)
	}
}

func TestSliceCacheGenericInterface(t *testing.T) {
	testData := []*Element{NewElement(&TestObject{}), NewElement(&TestObject{})}
	sliceCache := NewSliceCache(
		buildSliceLoader(testData),
		time.Second,
	)
	sliceCache.sync()

	actualLen := len(sliceCache.data)
	if ! reflect.DeepEqual(sliceCache.data, testData) {
		t.Error("invalid cache data", actualLen)
	}
}

func buildSliceLoader(values []*Element) func() []*Element {
	return func() []*Element {
		return values
	}
}