package gache

import (
	"time"
	"fmt"
)

type SliceCache struct {
	data     []*Element
	loader   func() []*Element
	duration time.Duration
}

func NewSliceCache() *SliceCache {
	var nilSliceLoader = func() []*Element {
		return nil
	}

	c := &SliceCache{
		loader: nilSliceLoader,
		duration: MaxDuration,
	}
	c.launchLoader()
	return c
}

func (c *SliceCache) SetLoader(loader func() []*Element) {
	c.loader = loader
}

func (c *SliceCache) RefreshAfterWrite(duration time.Duration) {
	c.duration = duration
}

func (c *SliceCache) Get() []*Element {
	if c.data == nil {
		fmt.Println("initializing cache")
		c.sync()
	}
	return c.data
}

func (c *SliceCache) String() string {
	builder := make([]string, len(c.data))

	for i, v := range c.data {
		builder[i] = fmt.Sprintf("%+v, ", v)
	}

	return fmt.Sprintf(
		"SliceCache(%+v,%+v,%+v)",
		builder,
		c.loader,
		c.duration,
	)
}

func (c *SliceCache) sync() {
	c.data = c.loader()
}

func (c *SliceCache) launchLoader() {
	// should only launch once
	go func() {
		for range time.Tick(c.duration) {
			c.sync()
		}
	}()
}

func dumpSlice(s []*Element) {
	builder := make([]string, len(s))

	for i, v := range s {
		builder[i] = fmt.Sprintf("%+v, ", v)
	}
	println(fmt.Sprintf("%v", builder))
}