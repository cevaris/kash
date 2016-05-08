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
	AccessedAt      time.Time
	CreatedAt       time.Time
	RefreshEligible bool
	Value           interface{}
}

func newElement(value interface{}) *element {
	return &element{
		AccessedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
		RefreshEligible: false,
		Value: value,
	}
}

func (e *element) WriteStale(now time.Time, ttl time.Duration) bool {
	return e.CreatedAt.Before(now.Add(-1 * ttl))
}

func (e *element) AccessStale(now time.Time, ttl time.Duration) bool {
	return e.AccessedAt.Before(now.Add(-1 * ttl))
}

func (e *element) FlagToRefresh() {
	e.RefreshEligible = true
}

func (e *element) ShouldRefresh() bool {
	return e.RefreshEligible
}

func (e *element) String() string {
	return fmt.Sprintf("%+v %+v", e.Value, e.CreatedAt)
}

func (e1 *element) Compare(e2 *element) bool {
	return e1.Value == e2.Value
}