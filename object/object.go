/*
In inheritance:
class Base:
    def __str__(self):
        return "<%s object at %s>" % (type(self), id(self))

class SubClass(Base): pass

                            inherit
                    type <------------ metaclass
                    | ^                   |
                    | |                   |
        instantiate | |inherit            v
                    | |                  ...
                    v |     inherit
                   object <----------------------------------- class
                      ^ ^                                        |
                       \ \    inherit                            | inst 
                        \ ---------------------------- int       v
                      .  \ inherit                      |     instance
                      .   --------------- function      |
                      .                      |          | inst
                      .                 inst |          v
                                             |        number instance
                                             v
                                         func instance

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
    Py__getattribute__(string) Object
    Py__setattr__(string, Object)
    Py__init__(Object)
}

// type PyObject struct {
// }
// func (po *PyObject) __str__() {
// }

type PyObject struct {
}

func (o *PyObject) GetObjType() ObjType {return NONE}
func (o *PyObject) Py__getattribute__(attr string) Object {return o}
func (o *PyObject) Py__setattr__(attr string, valObj Object) {}
func (o *PyObject) Py__init__(Object) {}

type NoneObject struct {
    PyObject
    Value   int
}
func (no *NoneObject) GetObjType() ObjType {return NONE}
func (no NoneObject) String() string {
    return "None"
}

type BoolObject struct {
    PyObject
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
    PyObject
    Value   int
}

func (no *NumberObject) GetObjType() ObjType {return NUMBER}
func (no NumberObject) String() string {return fmt.Sprint(no.Value)}

type StringObject struct {
    PyObject
    Value   string
}

func (self *StringObject) GetObjType() ObjType {return STRING}
func (self StringObject) String() string {return "'" + self.Value + "'"}

type ListObject struct {
    PyObject
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

    PyObject
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
    PyObject
    Name    string
    Params  []string
    Body    []ast.Statement
}

func (fo *FunctionObject) GetObjType() ObjType {return FUNCTION}
func (fo *FunctionObject) String() string {
    return fmt.Sprintf("<function %s at %p>", fo.Name, fo)
}

type ClassObject struct {
    PyObject
    Name          string
    Py__base__    Object
    Py__dict__    map[string]Object
}

func (co *ClassObject) GetObjType() ObjType {return CLASS}
func (co *ClassObject) Py__getattribute__(attr string) Object {
    // FIXME: defualt to object
    val := co.Py__dict__[attr]
    if val != nil {
        return val
    }

    if co.Py__base__ == nil {
        return nil
    }

    return co.Py__base__.Py__getattribute__(attr)
}
func (co *ClassObject) Py__setattr__(attr string, val Object) {
    co.Py__dict__[attr] = val
}
func (co *ClassObject) Py__new__() Object {
    return &InstanceObject{
        Py__class__: co,
        Py__dict__: map[string]Object{},
    }
}
// func (co *ClassObject) Py__call__() Object {
//     instObj = co.Py__new__()
// }
// func (co *ClassObject) Call() Object {
//     return co.Py__call__()
// }

func (co *ClassObject) String() string {
    return fmt.Sprintf("<class %s at %p>", co.Name, co)
}

type BoundMethod struct {
    PyObject
    Func        *FunctionObject
    Inst        *InstanceObject
} 
func (bm *BoundMethod) String() string {
    return fmt.Sprintf("<bound method %s.%s object at %p>",
        bm.Inst.Py__class__.Name, bm.Func.Name, bm)
}

type InstanceObject struct {
    PyObject
    Py__class__ *ClassObject
    Py__dict__ map[string]Object
}

func (io *InstanceObject) Py__getattribute__(attr string) Object {
    targetObj, ok := io.Py__dict__[attr]
    if ok {
        return targetObj 
    }

    switch targetObj := io.Py__class__.Py__getattribute__(attr).(type) {
    case *FunctionObject:
        // FIXME: here supposed to be return the identical method everytime
        return &BoundMethod{
            Func: targetObj,
            Inst: io,
            }
    default:
        // not supposed to be nil
        return targetObj
    }
}
func (io *InstanceObject) Py__setattr__(attr string, val Object) {
    io.Py__dict__[attr] = val
}

func (io *InstanceObject) String() string {
    return fmt.Sprintf("<%v objects at %p>", io.Py__class__.Name, io)
}

// temporary
type Print struct {
    PyObject
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
    PyObject
    Msg string
}

func (eo *ExceptionObject) GetObjType() ObjType {return EXCEPTION}
func (eo *ExceptionObject) ErrorMsg() string {return eo.Msg}
func (eo *ExceptionObject) String() string {return "Exception"}

type BuiltinClass struct {
    PyObject
    Name string
}

func (bc *BuiltinClass) String() string {
    return bc.Name
}

type SuperInstance struct {
    PyObject
    Py__self__ *InstanceObject
}

func (si *SuperInstance) Py__getattribute__(attr string) Object {
    switch targetObj := si.Py__self__.Py__class__.Py__base__.Py__getattribute__(attr).(type) {
    case *FunctionObject:
        return &BoundMethod{
            Func: targetObj,
            Inst: si.Py__self__,
            }
    }
    return nil
}

func (si *SuperInstance) String() string {
    return "<unknown now>"
}

