package main

import (
	cache "github.com/cevaris/gache"
	"time"
	"fmt"
)

func sliceCacheStatic() {
	fmt.Println("exampleStaticCache")
	sliceCache := cache.NewSliceCache()
	sliceCache.SetLoader(func() []*cache.Element {
		return []*cache.Element{
			cache.NewElement(time.Now().UTC()),
			cache.NewElement(time.Now().UTC().Add(-5 * time.Second)),
		}
	})

	var slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.Value.(time.Time))
	}

	time.Sleep(1 * time.Second)

	slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.Value.(time.Time))
	}
	fmt.Println("\n")
}

func sliceCacheWithRefresh() {
	fmt.Println("exampleLoaderWithRefresh")
	sliceCache := cache.NewSliceCache()
	sliceCache.RefreshAfterWrite(100 * time.Millisecond)
	sliceCache.SetLoader(func() []*cache.Element {
		return []*cache.Element{
			cache.NewElement(time.Now().UTC()),
			cache.NewElement(time.Now().UTC().Add(-5 * time.Second)),
		}
	})

	var slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.Value.(time.Time))
	}

	time.Sleep(1 * time.Second)

	slice = sliceCache.Get()
	for i, e := range slice {
		fmt.Println(i, e.Value.(time.Time))
	}
	fmt.Println("\n")
}

func main() {
	sliceCacheStatic()
	sliceCacheWithRefresh()
}