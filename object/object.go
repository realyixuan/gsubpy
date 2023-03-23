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
    Py__repr__() string
}

type Class interface {
    Object
    Py__init__(Object)
    Py__new__(Class) *PyInstance
    Py__name__() string
    Py__base__() Class
}

type Function interface {
    Object
    Py__call__() Object
}

type BuiltinFunction interface {
    Function
    isBuiltinFunction()
}

type BuiltinObjectNew struct {
    ObjectClass
    Func    func(Class) *PyInstance
}
func (b *BuiltinObjectNew) Py__call__() Object {return nil}
func (b *BuiltinObjectNew) isBuiltinFunction() {}
func (b *BuiltinObjectNew) Call(cls Class) *PyInstance {
    return b.Func(cls)
}

type ObjectClass struct {
}

func (o *ObjectClass) Py__repr__() string {return "object"}
func (o *ObjectClass) GetObjType() ObjType {return NONE}
func (o *ObjectClass) Py__getattribute__(attr string) Object {
    if attr == "__new__" {
        return &BuiltinObjectNew{
            Func: o.Py__new__,
            }
    }
    return nil
}
func (o *ObjectClass) Py__setattr__(attr string, valObj Object) {}
func (o *ObjectClass) Py__new__(cls Class) *PyInstance {
    return &PyInstance{
        Py__class__: cls,
        Py__dict__: map[string]Object{},
        }
}
func (o *ObjectClass) Py__init__(Object) {}
func (o *ObjectClass) Py__name__() string {return "object"}
func (o *ObjectClass) Py__base__() Class {return nil}

var PyObject = &ObjectClass{}

type Pytype struct {
    ObjectClass
}
func (t *Pytype) Py__repr__() string {return "type"}
func (t *Pytype) Py__new__(mcs *Pytype, name string, base *PyClass, attrs map[string]Object) *PyClass {
    return &PyClass{
        Name: name,
        Base: base,
        Dict: attrs,
        }
}

type NoneObject struct {
    ObjectClass
    Value   int
}
func (no *NoneObject) Py__repr__() string {return "None"}
func (no *NoneObject) GetObjType() ObjType {return NONE}
func (no NoneObject) String() string {
    return "None"
}

type BoolObject struct {
    ObjectClass
    Value   int
}
func (bo *BoolObject) GetObjType() ObjType {return BOOL}
func (bo *BoolObject) Py__repr__() string {
    if bo.Value == 1 {
        return "True"
    } else {
        return "False"
    }
}

type NumberObject struct {
    ObjectClass
    Value   int
}

func (no *NumberObject) GetObjType() ObjType {return NUMBER}
func (no *NumberObject) Py__repr__() string {return fmt.Sprint(no.Value)}

type StringObject struct {
    ObjectClass
    Value   string
}

func (self *StringObject) GetObjType() ObjType {return STRING}
func (self *StringObject) Py__repr__() string {return "'" + self.Value + "'"}

type ListObject struct {
    ObjectClass
    Items   []Object
}

func (self *ListObject) GetObjType() ObjType {return LIST}
func (self *ListObject) Py__repr__() string {
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

    ObjectClass
    Map   map[Object]Object 
                            
}

func (do *DictObject) GetObjType() ObjType {return DICT}
func (do *DictObject) Py__repr__() string {
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
    ObjectClass
    Name    string
    Params  []string
    Body    []ast.Statement
}

func (fo *FunctionObject) GetObjType() ObjType {return FUNCTION}

func (fo *FunctionObject) Py__call__() Object {return nil}

func (fo *FunctionObject) Py__repr__() string {
    return fmt.Sprintf("<function %s at %p>", fo.Name, fo)
}

type PyClass struct {
    ObjectClass
    Name          string
    Base          Class
    Dict    map[string]Object
}

func (pc *PyClass) Py__name__() string {return pc.Name}
func (pc *PyClass) GetObjType() ObjType {return CLASS}
func (pc *PyClass) Py__getattribute__(attr string) Object {
    // FIXME: defualt to object
    val := pc.Dict[attr]
    if val != nil {
        return val
    }

    if pc.Base == nil {
        return nil
    }

    return pc.Base.Py__getattribute__(attr)
}
func (pc *PyClass) Py__setattr__(attr string, val Object) {
    pc.Dict[attr] = val
}
func (pc *PyClass) Py__new__(cls Class) *PyInstance {
    // FIXME: haven't execute customized __new__
    return pc.Base.Py__new__(cls)
}
func (pc *PyClass) Py__base__() Class {
    // FIXME: haven't execute customized __new__
    return pc.Base
}

func (pc *PyClass) Py__repr__() string {
    return fmt.Sprintf("<class %s at %p>", pc.Name, pc)
}

type BoundMethod struct {
    ObjectClass
    Func        *FunctionObject
    Inst        *PyInstance
} 
func (bm *BoundMethod) Py__repr__() string {
    return fmt.Sprintf("<bound method %s.%s object at %p>",
        bm.Inst.Py__class__.Py__name__(), bm.Func.Name, bm)
}
func (bm *BoundMethod) Py__call__() Object {return nil}

type PyInstance struct {
    ObjectClass
    Py__class__ Class
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

func (pi *PyInstance) Py__repr__() string {
    return fmt.Sprintf("<%v objects at %p>", pi.Py__class__.Py__name__(), pi)
}

// temporary
type Print struct {
    ObjectClass
}

func (p *Print) GetObjType() ObjType {return FUNCTION}
func (p *Print) Py__repr__() string {return "print"}
func (p *Print) Call(objs []Object) {
    for _, obj := range objs {
        fmt.Print(obj.Py__repr__())
        fmt.Print(" ")
    }
    fmt.Println()
}

type Exception interface {
    Object
    ErrorMsg() string
}

type ExceptionObject struct {
    ObjectClass
    Msg string
}

func (eo *ExceptionObject) GetObjType() ObjType {return EXCEPTION}
func (eo *ExceptionObject) ErrorMsg() string {return eo.Msg}
func (eo *ExceptionObject) String() string {return "Exception"}
func (eo *ExceptionObject) Py__repr__() string {return "Exception"}

type BuiltinClass struct {
    ObjectClass
    Name string
}

func (bc *BuiltinClass) String() string {
    return bc.Name
}

type SuperInstance struct {
    ObjectClass
    Py__self__ *PyInstance
}

func (si *SuperInstance) Py__getattribute__(attr string) Object {
    switch targetObj := si.Py__self__.Py__class__.Py__base__().Py__getattribute__(attr).(type) {
    case *FunctionObject:
        return &BoundMethod{
            Func: targetObj,
            Inst: si.Py__self__,
            }
    }
    return nil
}

func (si *SuperInstance) Py__repr__() string {
    return "<unknown now>"
}

