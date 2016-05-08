package gache
//
//import (
//	"testing"
//	"time"
//	"reflect"
//)
//
//const TestCacheSize int = 10
//
//type TestObject struct{}
//
//func TestNewMapCache(t *testing.T) {
//	sliceCache := NewMapCache(
//		buildMapLoader(10),
//	)
//	sliceCache.sync()
//
//	actualLen := len(sliceCache.data)
//	if len(sliceCache.data) != TestCacheSize {
//		t.Error("invalid default cache size", actualLen)
//	}
//}
//
//func TestNilSliceCacheLoader(t *testing.T) {
//	sliceCache := NewSliceCache(buildMapLoader(nil), time.Second)
//	sliceCache.sync()
//
//	if sliceCache.data != nil {
//		t.Error("cache not nil")
//	}
//}
//
//func TestSliceCacheInt(t *testing.T) {
//	testData := []*Element{NewElement(10), NewElement(100)}
//	sliceCache := NewSliceCache(
//		buildMapLoader(testData),
//	)
//	sliceCache.sync()
//
//	actualLen := len(sliceCache.data)
//	if ! reflect.DeepEqual(sliceCache.data, testData) {
//		t.Error("invalid cache data", actualLen)
//	}
//}
//
//func TestSliceCacheString(t *testing.T) {
//	testData := []*Element{NewElement("alpha"), NewElement("bravo")}
//	sliceCache := NewSliceCache(
//		buildMapLoader(testData),
//	)
//	sliceCache.sync()
//
//	actualLen := len(sliceCache.data)
//	if ! reflect.DeepEqual(sliceCache.data, testData) {
//		t.Error("invalid cache data", actualLen)
//	}
//}
//
//func TestSliceCacheGenericInterface(t *testing.T) {
//	testData := []*Element{NewElement(&TestObject{}), NewElement(&TestObject{})}
//	sliceCache := NewSliceCache(
//		buildMapLoader(testData),
//	)
//	sliceCache.sync()
//
//	actualLen := len(sliceCache.data)
//	if ! reflect.DeepEqual(sliceCache.data, testData) {
//		t.Error("invalid cache data", actualLen)
//	}
//}
//
//func buildMapLoader(value interface{}) func() func(interface{}) interface{} {
//	return func() func(interface{}) interface{} {
//		return value
//	}
//}