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
	Value     interface{}
	CreatedAt time.Time
}

func newElement(value interface{}) *element {
	return &element{
		Value: value,
		CreatedAt: time.Now().UTC(),
	}
}

func (e *element) Stale(now time.Time, ttl time.Duration) bool {
	return e.CreatedAt.Before(now.Add(-1 * ttl))
}

func (e *element) String() string {
	return fmt.Sprintf("%+v %+v", e.Value, e.CreatedAt)
}

func (e1 *element) Compare(e2 *element) bool {
	return e1.Value == e2.Value
}