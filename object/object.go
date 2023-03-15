/*
In inheritance:
class Base:
    def __str__(self):
        return "<%s object at %s>" % (type(self), id(self))

class SubClass: pass

with composition not inheritance, there are some issues for
reuse behaviours

*/



package object

import (
    "fmt"
    
    "gsubpy/ast"
)

type ObjType int

const (
    LIST ObjType = iota
    DICT
    STRING
    NUMBER
    BOOL
    FUNCTION
    NONE
    EXCEPTION
    CLASS
)

type Object interface {
    GetObjType() ObjType
}

// type PyObject struct {
// }
// func (po *PyObject) __str__() {
// }

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

type DictObject struct {
    /*
        FIX: if key is the pointer of Objects, 
        then there must be some issues, such as 
        that a string even can't match the key of
        another equivalent string Object
        like:
        if d['a'] = 1
        then call it again, d['a'] will raise key-not-exist error

    */

    Map   map[Object]Object 
                            
}

func (do *DictObject) GetObjType() ObjType {return DICT}
func (do *DictObject) String() string {
    var s string
    s += "{"
    var i = 0
    for k, v := range do.Map {
        s += fmt.Sprintf("%v:%v", k, v)
        if i == len(do.Map) - 1 {
            break
        }
        s += ", "
        i++
    }
    s += "}"
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

type ClassObject struct {
    Name    string
    Dict    map[string]Object
}

func (co *ClassObject) GetObjType() ObjType {return CLASS}
func (co *ClassObject) String() string {
    return fmt.Sprintf("<class %s at %p>", co.Name, co)
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

type Exception interface {
    Object
    ErrorMsg() string
}

type ExceptionObject struct {
    Msg string
}

func (eo *ExceptionObject) GetObjType() ObjType {return EXCEPTION}
func (eo *ExceptionObject) ErrorMsg() string {return eo.Msg}
func (eo *ExceptionObject) String() string {return "Exception"}

