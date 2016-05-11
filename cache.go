package kash

import (
	"time"
	"math"
	"fmt"
)

// Element example; https://golang.org/src/container/list/list_test.go

const MaxDuration = time.Nanosecond * math.MaxInt64

type Cache interface{}

type element struct {
	AccessedAt time.Time
	WriteAt    time.Time
	Value      interface{}
}

func newElement(value interface{}) *element {
	return &element{
		AccessedAt: time.Now().UTC(),
		WriteAt: time.Now().UTC(),
		Value: value,
	}
}

func (e *element) AccessStale(now time.Time, ttl time.Duration) bool {
	return e.AccessedAt.Before(now.Add(-1 * ttl))
}

func (e *element) WriteStale(now time.Time, ttl time.Duration) bool {
	return e.WriteAt.Before(now.Add(-1 * ttl))
}

func (e *element) String() string {
	return fmt.Sprintf("%+v %+v", e.Value, e.WriteAt)
}

func (e1 *element) Compare(e2 *element) bool {
	return e1.Value == e2.Value
}