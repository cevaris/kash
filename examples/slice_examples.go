package main

import (
	cache "github.com/cevaris/kash"
	"time"
	"fmt"
)

func sliceCacheStatic() {
	fmt.Println("sliceCacheStatic")
	sliceCache := cache.NewSliceCache()
	sliceCache.SetLoader(func() []interface{} {
		return []interface{}{
			time.Now().UTC(),
			time.Now().UTC().Add(-5 * time.Second),
		}
	})

	var slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.(time.Time))
	}

	time.Sleep(1 * time.Second)

	slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.(time.Time))
	}
	fmt.Println("\n")
}

func sliceCacheWithRefresh() {
	fmt.Println("sliceCacheWithRefresh")
	sliceCache := cache.NewSliceCache()
	sliceCache.RefreshAfterWrite(100 * time.Millisecond)
	sliceCache.SetLoader(func() []interface{} {
		return []interface{}{
			time.Now().UTC(),
			time.Now().UTC().Add(-5 * time.Second),
		}
	})

	var slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.(time.Time))
	}

	time.Sleep(1 * time.Second)

	slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.(time.Time))
	}
	fmt.Println("\n")
}

func main() {
	sliceCacheStatic()
	sliceCacheWithRefresh()
}