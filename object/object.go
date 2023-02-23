package object

type Object interface {
    isObject()
}

type NumberObject struct {
    Value   int
}

func (no *NumberObject) isObject() {}

