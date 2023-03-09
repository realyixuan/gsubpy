package object

import (
    "fmt"
    
    "gsubpy/ast"
)

type ObjType int

const (
    LIST ObjType = iota
    STRING
    NUMBER
    BOOL
    FUNCTION
    NONE
)

type Object interface {
    GetObjType() ObjType
}

type NoneObject struct {
    Value   int
}
func (no *NoneObject) GetObjType() ObjType {return NONE}
func (no NoneObject) String() string {
    return "None"
}

type BoolObject struct {
    Value   int
}
func (bo *BoolObject) GetObjType() ObjType {return BOOL}
func (bo BoolObject) String() string {
    if bo.Value == 1 {
        return "True"
    } else {
        return "False"
    }
}

type NumberObject struct {
    Value   int
}

func (no *NumberObject) GetObjType() ObjType {return NUMBER}
func (no NumberObject) String() string {return fmt.Sprint(no.Value)}

type StringObject struct {
    Value   string
}

func (self *StringObject) GetObjType() ObjType {return STRING}
func (self StringObject) String() string {return "'" + self.Value + "'"}

type ListObject struct {
    Items   []Object
}

func (self *ListObject) GetObjType() ObjType {return LIST}
func (self ListObject) String() string {
    var s string
    s += "["
    if len(self.Items) > 0 {
        s += fmt.Sprintf("%v", self.Items[0])
    }
    for _, item := range self.Items[1:] {
        s += ", "
        s += fmt.Sprintf("%v", item)
    }
    s += "]"
    return s
}

type FunctionObject struct {
    Name    string
    Params  []string
    Body    []ast.Statement
}

func (fo *FunctionObject) GetObjType() ObjType {return FUNCTION}
func (fo FunctionObject) String() string {
    return fmt.Sprintf("<function %s at %p>", fo.Name, &fo)
}

// temporary
type Print struct {
}

func (p *Print) GetObjType() ObjType {return FUNCTION}
func (p Print) String() string {return "print"}
func (p *Print) Call(objs []Object) {
    for _, obj := range objs {
        fmt.Print(obj)
        fmt.Print(" ")
    }
    fmt.Println()
}

