package gache

import (
	"testing"
	"time"
)

const TestCacheSize int = 10

func TestNewSliceCache(t *testing.T) {
	sliceCache := NewSliceCache(
		func() []interface{} {
			return make([]interface{}, TestCacheSize)
		},
		time.Second,
	)

	actualLen := len(sliceCache.Data)
	if len(sliceCache.Data) != TestCacheSize {
		t.Error("invalid default cache size", actualLen)
	}
}

func TestNoUpCacheLoader(t *testing.T) {
	noopLoader := func() []interface{} {
		return nil
	}

	sliceCache := NewSliceCache(noopLoader, time.Second)
	if sliceCache.Data != nil {
		t.Error("cache not nil")
	}
}

func TestSliceCacheInt(t *testing.T) {
	sliceCache := NewSliceCache(
		func() []interface{} {
			newData := []string{"test", "bravo"}
			return newData
		},
		time.Second,
	)

	actualLen := len(sliceCache.Data)
	if len(sliceCache.Data) != TestCacheSize {
		t.Error("invalid default cache size", actualLen)
	}
}