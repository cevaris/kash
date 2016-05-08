package gache

import (
	"time"
	"math"
	"fmt"
)

// Element example; https://golang.org/src/container/list/list_test.go

const MaxDuration = time.Nanosecond * math.MaxInt64

type Cache interface{}

type Element struct {
	Value     interface{}
	CreatedAt time.Time
}

func NewElement(value interface{}) *Element {
	return &Element{
		Value: value,
		CreatedAt: time.Now().UTC(),
	}
}

func (e *Element) Stale(now time.Time, ttl time.Duration) bool {
	return e.CreatedAt.Before(now.Add(-1 * ttl))
}

func (e *Element) String() string {
	return fmt.Sprintf("%+v %+v", e.Value, e.CreatedAt)
}

func (e1 *Element) Compare(e2 *Element) bool {
	return e1.Value == e2.Value
}