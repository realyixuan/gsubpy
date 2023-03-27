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

type Type int

const (
    LIST Type = iota
    DICT
    STRING
    INTEGER
    BOOL
    FUNCTION
    NONE
    EXCEPTION
    CLASS
    INSTANCE
    SUPER
    METHOD
    TYPE
)

type Object interface {
    Type() Type
    Py__getattribute__(*PyStrInst) Object
    Py__setattr__(*PyStrInst, Object)
    Py__repr__() *PyStrInst
}

type Class interface {
    Object
    Py__init__(*PyInst)
    Py__new__(Class) *PyInst
    Py__name__() *PyStrInst
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
    Func    func(Class) *PyInst
}
func (b *BuiltinObjectNew) Type() Type {return METHOD}
func (b *BuiltinObjectNew) Py__call__() Object {return Py_None}
func (b *BuiltinObjectNew) isBuiltinFunction() {}   // TODO: builtin type interface rename
func (b *BuiltinObjectNew) Call(cls Class) *PyInst {
    return b.Func(cls)
}
func (b *BuiltinObjectNew) Py__getattribute__(*PyStrInst) Object {return nil}
func (b *BuiltinObjectNew) Py__setattr__(*PyStrInst, Object) {}
func (b *BuiltinObjectNew) Py__repr__() *PyStrInst {
    return &PyStrInst{"<builtin __new__>"}
}

type ObjectClass struct {
}

func (o *ObjectClass) Py__repr__() *PyStrInst {
    return &PyStrInst{"object"}
}
func (o *ObjectClass) Type() Type {return NONE}
func (o *ObjectClass) Py__getattribute__(attr *PyStrInst) Object {
    if attr.Value == "__new__" {
        return &BuiltinObjectNew{
            Func: o.Py__new__,
            }
    }
    return nil
}
func (o *ObjectClass) Py__setattr__(attr *PyStrInst, valObj Object) {}
func (o *ObjectClass) Py__new__(cls Class) *PyInst {
    return &PyInst{
        Py__class__: cls,
        Py__dict__: map[string]Object{},
        }
}
func (o *ObjectClass) Py__init__(*PyInst) {}
func (o *ObjectClass) Py__name__() *PyStrInst {return &PyStrInst{"object"}}
func (o *ObjectClass) Py__base__() Class {return nil}

var PyObject = &ObjectClass{}

type Pytype struct {
}
func (t *Pytype) Type() Type {return TYPE}
func (t *Pytype) Py__init__(*PyInst) {}
func (t *Pytype) Py__name__() *PyStrInst {return &PyStrInst{"type"}}
func (t *Pytype) Py__base__() Class {return PyObject}
func (t *Pytype) Py__repr__() *PyStrInst {
    return &PyStrInst{"type"}
}
func (t *Pytype) Py__new__(Class) *PyInst {return nil}
func (t *Pytype) Py__pnew__(mcs *Pytype, name string, base *PyClass, attrs map[string]Object) *PyClass {
    return &PyClass{
        Name: name,
        Base: base,
        Dict: attrs,
        }
}
func (t *Pytype) Py__getattribute__(*PyStrInst) Object {return nil}
func (t *Pytype) Py__setattr__(*PyStrInst, Object) {}

var Py_type = &Pytype{}

type NoneInst struct {
    Value   int
}
func (ni *NoneInst) Py__repr__() *PyStrInst {
    return &PyStrInst{"None"}
}
func (ni *NoneInst) Type() Type {return NONE}
func (ni *NoneInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (ni *NoneInst) Py__setattr__(*PyStrInst, Object) {}

var Py_None = &NoneInst{Value: 0}

type BoolInst struct {
    Value   int
}
func (bi *BoolInst) Type() Type {return BOOL}
func (bi *BoolInst) Py__repr__() *PyStrInst {
    if bi.Value == 1 {
        return &PyStrInst{"True"}
    } else {
        return &PyStrInst{"False"}
    }
}
func (bi *BoolInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (bi *BoolInst) Py__setattr__(*PyStrInst, Object) {}

type IntegerInst struct {
    Value   int
}

func (ni *IntegerInst) Type() Type {return INTEGER}
func (ni *IntegerInst) Py__repr__() *PyStrInst {
    return &PyStrInst{fmt.Sprint(ni.Value)}
}
func (ni *IntegerInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (ni *IntegerInst) Py__setattr__(*PyStrInst, Object) {}

type PyStr struct {
    Name          string
    Base          Class
}
func (ps *PyStr) Py__init__(*PyInst) {}
func (ps *PyStr) Py__name__() *PyStrInst {return &PyStrInst{"str"}}
func (ps *PyStr) Type() Type {return TYPE}
func (ps *PyStr) Py__getattribute__(attr *PyStrInst) Object {return nil}
func (ps *PyStr) Py__setattr__(*PyStrInst, Object) {}
func (ps *PyStr) Py__new__(cls Class) *PyStrInst {return &PyStrInst{""}}
func (pc *PyStr) Py__base__() Class {return pc.Base}
func (pc *PyStr) Py__repr__() *PyStrInst {
    return &PyStrInst{fmt.Sprintf("<class '%s'>", "str")}
}

var Py_str = &PyStr{
    Name: "str",
    Base: PyObject,
}

type PyStrInst struct {
    Value   string
}

func (si *PyStrInst) Type() Type {return STRING}
func (si *PyStrInst) Py__repr__() *PyStrInst {
    return &PyStrInst{"'" + si.Value + "'"}
}
func (si *PyStrInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (si *PyStrInst) Py__setattr__(*PyStrInst, Object) {}

type ListInst struct {
    Items   []Object
}

func (li *ListInst) Type() Type {return LIST}
func (li *ListInst) Py__repr__() *PyStrInst {
    var s string
    s += "["
    if len(li.Items) > 0 {
        s += fmt.Sprintf("%v", li.Items[0])
    }
    for _, item := range li.Items[1:] {
        s += ", "
        s += fmt.Sprintf("%v", item)
    }
    s += "]"
    return &PyStrInst{s}
}

func (li *ListInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (li *ListInst) Py__setattr__(*PyStrInst, Object) {}

type DictInst struct {
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

func (di *DictInst) Type() Type {return DICT}
func (di *DictInst) Py__repr__() *PyStrInst {
    var s string
    s += "{"
    var i = 0
    for k, v := range di.Map {
        s += fmt.Sprintf("%v:%v", k, v)
        if i == len(di.Map) - 1 {
            break
        }
        s += ", "
        i++
    }
    s += "}"
    return &PyStrInst{s}
}

func (di *DictInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (di *DictInst) Py__setattr__(*PyStrInst, Object) {}

// TODO: rename instance
type FunctionInst struct {
    Name    string
    Params  []string
    Body    []ast.Statement
}

func (fi *FunctionInst) Type() Type {return FUNCTION}

func (fi *FunctionInst) Py__call__() Object {return Py_None}

func (fi *FunctionInst) Py__repr__() *PyStrInst {
    return &PyStrInst{fmt.Sprintf("<function %s at %p>", fi.Name, fi)}
}
func (fi *FunctionInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (fi *FunctionInst) Py__setattr__(*PyStrInst, Object) {}

type PyClass struct {
    Name          string
    Base          Class
    Dict    map[string]Object
}

func (pc *PyClass) Py__init__(*PyInst) {}
func (pc *PyClass) Py__name__() *PyStrInst {return &PyStrInst{pc.Name}}
func (pc *PyClass) Type() Type {return CLASS}
func (pc *PyClass) Py__getattribute__(attr *PyStrInst) Object {
    // FIXME: defualt to object
    val := pc.Dict[attr.Value]
    if val != nil {
        return val
    }

    if pc.Base == nil {
        return nil
    }

    return pc.Base.Py__getattribute__(attr)
}
func (pc *PyClass) Py__setattr__(attr *PyStrInst, val Object) {
    pc.Dict[attr.Value] = val
}
func (pc *PyClass) Py__new__(cls Class) *PyInst {
    // FIXME: haven't execute customized __new__
    return pc.Base.Py__new__(cls)
}
func (pc *PyClass) Py__base__() Class {
    // FIXME: haven't execute customized __new__
    return pc.Base
}

func (pc *PyClass) Py__repr__() *PyStrInst {
    return &PyStrInst{fmt.Sprintf("<class %s at %p>", pc.Name, pc)}
}

type BoundMethod struct {
    Func        *FunctionInst
    Inst        *PyInst
} 
func (bm *BoundMethod) Py__repr__() *PyStrInst {
    s := fmt.Sprintf("<bound method %s.%s object at %p>",
            bm.Inst.Py__class__.Py__name__().Value, bm.Func.Name, bm)

    return &PyStrInst{s}
}
func (bm *BoundMethod) Py__call__() Object {return Py_None}
func (bm *BoundMethod) Type() Type {return METHOD}
func (bm *BoundMethod) Py__getattribute__(*PyStrInst) Object {return nil}
func (bm *BoundMethod) Py__setattr__(*PyStrInst, Object) {}

type PyInst struct {
    Py__class__ Class
    Py__dict__ map[string]Object
}

func (pi *PyInst) Type() Type {return INSTANCE}
func (pi *PyInst) Py__getattribute__(attr *PyStrInst) Object {
    targetObj, ok := pi.Py__dict__[attr.Value]
    if ok {
        return targetObj 
    }

    switch targetObj := pi.Py__class__.Py__getattribute__(attr).(type) {
    case *FunctionInst:
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
func (pi *PyInst) Py__setattr__(attr *PyStrInst, val Object) {
    pi.Py__dict__[attr.Value] = val
}

func (pi *PyInst) Py__repr__() *PyStrInst {
    s := fmt.Sprintf("<%v objects at %p>", pi.Py__class__.Py__name__().Value, pi)
    return &PyStrInst{s}
}

// TODO: temporary
type Print struct {
}

func (p *Print) Type() Type {return FUNCTION}
func (p *Print) Py__repr__() *PyStrInst {return &PyStrInst{"print"}}
func (p *Print) Py__call__(objs []Object) {
    for _, obj := range objs {
        fmt.Print(obj.Py__repr__().Value)
        fmt.Print(" ")
    }
    fmt.Println()
}
func (p *Print) Py__getattribute__(*PyStrInst) Object {return nil}
func (p *Print) Py__setattr__(*PyStrInst, Object) {}

var Py_print = &Print{}

type Exception interface {
    Object
    ErrorMsg() string
}

type ExceptionInst struct {
    Msg string
}

func (ei *ExceptionInst) Type() Type {return EXCEPTION}
func (ei *ExceptionInst) ErrorMsg() string {return ei.Msg}
func (ei *ExceptionInst) String() string {return "Exception"}
func (ei *ExceptionInst) Py__repr__() *PyStrInst {return &PyStrInst{"Exception"}}
func (ei *ExceptionInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (ei *ExceptionInst) Py__setattr__(*PyStrInst, Object) {}

type BuiltinClass struct {
    ObjectClass
    Name string
}

func (bc *BuiltinClass) String() string {
    return bc.Name
}

type SuperInst struct {
    Py__self__ *PyInst
}

func (si *SuperInst) Py__getattribute__(attr *PyStrInst) Object {
    switch targetObj := si.Py__self__.Py__class__.Py__base__().Py__getattribute__(attr).(type) {
    case *FunctionInst:
        return &BoundMethod{
            Func: targetObj,
            Inst: si.Py__self__,
            }
    }
    return nil
}

func (si *SuperInst) Py__repr__() *PyStrInst {
    return &PyStrInst{"<super object>"}
}
func (si *SuperInst) Type() Type {return SUPER}
func (si *SuperInst) Py__setattr__(*PyStrInst, Object) {}

