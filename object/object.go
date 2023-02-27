package object

type Object interface {
    isObject()
}

type BoolObject struct {
    Value   int
}
func (bo *BoolObject) isObject() {}

type NumberObject struct {
    Value   int
}

func (no *NumberObject) isObject() {}

