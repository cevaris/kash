package gache

//import "time"
//
//type MapCache struct {
//	data     map[interface{}]interface{}
//	loader   func(interface{}) interface{}
//	duration time.Duration
//}
//
//func NewMapCache(loader func() []*Element, duration time.Duration) *MapCache {
//	c := &MapCache{duration: duration}
//	c.SetLoader(loader)
//	return c
//}
//
//func (c *MapCache) SetLoader(loader func(interface{}) interface{}) {
//	if c.loader == nil {
//		c.loader = loader
//		// Launch background loader
//		go func() {
//			for range time.Tick(c.duration) {
//				c.sync()
//			}
//		}()
//	}
//}
//
//func (c *MapCache) SetDuration(duration time.Duration) {
//	c.duration = duration
//}
//
//func (c *MapCache) sync() {
//	c.data = c.loader()
//}