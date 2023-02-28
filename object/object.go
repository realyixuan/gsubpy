package object

import (
    "gsubpy/ast"
)

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

type FunctionObject struct {
    Name    string
    Params  []string
    Body    []ast.Statement
}

func (no *FunctionObject) isObject() {}


