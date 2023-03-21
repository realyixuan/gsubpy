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
    Py__new__(*PyClass) *PyInstance
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
func (o *PyObject) Py__new__(cls *PyClass) *PyInstance {
    return &PyInstance{
        Py__class__: cls,
        Py__dict__: map[string]Object{},
        }
}
func (o *PyObject) Py__init__(Object) {}

var Pyobject = &PyObject{}

type Pytype struct {
    PyObject
}
func (t *Pytype) Py__new__(mcs *Pytype, name string, base *PyClass, attrs map[string]Object) *PyClass {
    return &PyClass{
        Name: name,
        Py__base__: base,
        Py__dict__: attrs,
        }
}

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

type PyClass struct {
    PyObject
    Name          string
    Py__base__    *PyClass
    Py__dict__    map[string]Object
}

func (pc *PyClass) GetObjType() ObjType {return CLASS}
func (pc *PyClass) Py__getattribute__(attr string) Object {
    // FIXME: defualt to object
    fmt.Println(attr, pc.Py__dict__)
    val := pc.Py__dict__[attr]
    if val != nil {
        return val
    }

    if pc.Py__base__ == nil {
        return nil
    }

    return pc.Py__base__.Py__getattribute__(attr)
}
func (pc *PyClass) Py__setattr__(attr string, val Object) {
    pc.Py__dict__[attr] = val
}
func (pc *PyClass) Py__new__(cls *PyClass) *PyInstance {
    // FIXME: unify PyObject and Pyclass
    // Maybe add a new interface: class-interface
    
    if pc.Py__base__ == nil {
        return Pyobject.Py__new__(cls)
    } else {
        return pc.Py__base__.Py__new__(cls)
    }
}
// func (co *ClassObject) Py__call__() Object {
//     instObj = co.Py__new__()
// }
// func (co *ClassObject) Call() Object {
//     return co.Py__call__()
// }

func (pc *PyClass) String() string {
    return fmt.Sprintf("<class %s at %p>", pc.Name, pc)
}

type BoundMethod struct {
    PyObject
    Func        *FunctionObject
    Inst        *PyInstance
} 
func (bm *BoundMethod) String() string {
    return fmt.Sprintf("<bound method %s.%s object at %p>",
        bm.Inst.Py__class__.Name, bm.Func.Name, bm)
}

type PyInstance struct {
    PyObject
    Py__class__ *PyClass
    Py__dict__ map[string]Object
}

func (pi *PyInstance) Py__getattribute__(attr string) Object {
    targetObj, ok := pi.Py__dict__[attr]
    if ok {
        return targetObj 
    }

    switch targetObj := pi.Py__class__.Py__getattribute__(attr).(type) {
    case *FunctionObject:
        // FIXME: here supposed to be return the identical method everytime
        return &BoundMethod{
            Func: targetObj,
            Inst: pi,
            }
    default:
        // not supposed to be nil
        return targetObj
    }
}
func (pi *PyInstance) Py__setattr__(attr string, val Object) {
    pi.Py__dict__[attr] = val
}

func (pi *PyInstance) String() string {
    return fmt.Sprintf("<%v objects at %p>", pi.Py__class__.Name, pi)
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
    Py__self__ *PyInstance
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

