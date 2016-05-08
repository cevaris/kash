package gache

type Cache interface{}

type Element struct{}
type Integer struct {
	Element
	Value int
}
type String struct {
	Element
	Value int
}
