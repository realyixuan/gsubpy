package object

import (
    "gsubpy/ast"
)

type ObjType string

const (
    STRING      = "string"
    NUMBER      = "number"
    BOOL        = "bool"
    FUNCTION    = "function"
)

type Object interface {
    GetObjType() ObjType
}

type BoolObject struct {
    Value   int
}
func (bo *BoolObject) GetObjType() ObjType {return BOOL}

type NumberObject struct {
    Value   int
}

func (no *NumberObject) GetObjType() ObjType {return NUMBER}

type StringObject struct {
    Value   string
}

func (self *StringObject) GetObjType() ObjType {return STRING}

type FunctionObject struct {
    Name    string
    Params  []string
    Body    []ast.Statement
}

func (no *FunctionObject) GetObjType() ObjType {return FUNCTION}

