package gache

type Cache interface{}

type Element struct {
	Value interface{}
}

func NewElement(value interface{}) *Element {
	return &Element{Value: value}
}
