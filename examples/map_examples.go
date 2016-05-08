package main

import (
	cache "github.com/cevaris/gache"
	"time"
	"fmt"
)

func mapCacheStatic() {
	fmt.Println("mapCacheStatic")
	mapCache := cache.NewMapCache()
	mapCache.SetLoader(func(key interface{}) interface{} {
		someExternalSource := map[interface{}]interface{}{
			"keyA": 1,
			"keyB": 2,
		}
		return someExternalSource[key]
	})

	if actual, exists := mapCache.Get("keyA"); exists {
		fmt.Println("Found value", actual)
	}

	if actual, exists := mapCache.Get("keyB"); exists {
		fmt.Println("Found value", actual)
	}

	if _, exists := mapCache.Get("doesNotExist"); !exists {
		fmt.Println("No value found for key 'doesNotExist'")
	}

	fmt.Println("\n")
}

func mapCacheWithRefresh() {
	fmt.Println("mapCacheWithRefresh")
	mapCache := cache.NewMapCache()
	mapCache.RefreshAfterWrite(500 * time.Millisecond)
	mapCache.SetLoader(func(key interface{}) interface{} {
		someExternalSource := map[interface{}]interface{}{
			"keyA": time.Now().UTC(),
			"keyB": time.Now().UTC().Second(),
		}
		return someExternalSource[key]
	})

	time.Sleep(time.Second)

	if actual, exists := mapCache.Get("keyA"); exists {
		fmt.Println("Found value", actual)
	}

	if actual, exists := mapCache.Get("keyB"); exists {
		fmt.Println("Found value", actual)
	}

	time.Sleep(2 * time.Second)

	if actual, exists := mapCache.Get("keyA"); exists {
		fmt.Println("Found value", actual)
	}

	if actual, exists := mapCache.Get("keyB"); exists {
		fmt.Println("Found value", actual)
	}

	fmt.Println("\n")
}

func main() {
	mapCacheStatic()
	mapCacheWithRefresh()
}