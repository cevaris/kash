package main

import (
	cache "github.com/cevaris/kash"
	"time"
	"fmt"
)

func mapCacheStatic() {
	fmt.Println("mapCacheStatic")
	mapCache := cache.NewMapCache()
	mapCache.SetLoad(func(key interface{}) interface{} {
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
	mapCache.RefreshAfterWrite(5 * time.Millisecond)
	mapCache.SetLoad(func(key interface{}) interface{} {
		someExternalSource := map[interface{}]interface{}{
			"keyA": 1,
			"keyB": 2,
		}
		return someExternalSource[key]
	})
	mapCache.SetReload(func(key interface{}, prev interface{}) {
		someExternalSource := map[interface{}]interface{}{
			"keyA": time.Now().Unix(),
		}
		go func() {
			time.Sleep(10 * time.Millisecond)
			mapCache.ReloadChan <- cache.KeyValue{Key: key, Value: someExternalSource[key]}
		}()
	})


	// Cache.Load
	if actual, exists := mapCache.Get("keyA"); exists {
		fmt.Println(fmt.Sprintf("Get keyA -> %+v", actual))
	}
	if actual, exists := mapCache.Get("keyB"); exists {
		fmt.Println(fmt.Sprintf("Get keyB -> %+v", actual))
	}
	time.Sleep(10 * time.Millisecond)

	// Cache.Reload (kicks of async job, returns old values)
	if actual, exists := mapCache.Get("keyA"); exists {
		fmt.Println(fmt.Sprintf("Get keyA -> %+v", actual))
	}
	if actual, exists := mapCache.Get("keyB"); exists {
		fmt.Println(fmt.Sprintf("Get keyB -> %+v", actual))
	}
	time.Sleep(20 * time.Millisecond)

	// Cache.Reload (returns new async returned values)
	if actual, exists := mapCache.Get("keyA"); exists {
		fmt.Println(fmt.Sprintf("Get keyA -> %+v", actual))
	}

	fmt.Println("\n")
}

func main() {
	mapCacheStatic()
	mapCacheWithRefresh()
}