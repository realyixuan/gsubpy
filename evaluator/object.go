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

// TODO: class should also have __call__
// TODO: distinguish class and function although both have __call__

package evaluator

import (
    "fmt"
    "strconv"
    
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
    Py__str__() *PyStrInst
    Py__class__() Class
}

type Class interface {
    Object
    Py__init__(Object, ...Object)
    Py__new__(Class, ...Object) Object
    Py__name__() *PyStrInst
    Py__base__() Class
}

type Function interface {
    Object
    Py__name__() *PyStrInst
    Py__call__(...Object) Object
}

type BuiltinFunction interface {
    Function
    Call()
}

// TODO: should be class
type PyNew struct {
    Func    func(Class, ...Object) Object
}
func (n *PyNew) Py__class__() Class {return Py_Function}
func (n *PyNew) Type() Type {return METHOD}
func (n *PyNew) Call() {}
func (n *PyNew) Py__call__(objs ...Object) Object {
    return n.Func(objs[0].(Class))
}
func (n *PyNew) Py__name__() *PyStrInst {return NewStrInst("<builtin __new__>")}
func (n *PyNew) Py__getattribute__(*PyStrInst) Object {return nil}
func (n *PyNew) Py__setattr__(*PyStrInst, Object) {}
func (n *PyNew) Py__repr__() *PyStrInst {
    return NewStrInst("<builtin __new__>")
}
func (n *PyNew) Py__str__() *PyStrInst {
    return n.Py__repr__()
}

type ObjectClass struct {
    Base Class
}

func (o *ObjectClass) Py__class__() Class {return Py_type}
func (o *ObjectClass) Py__repr__() *PyStrInst {
    return NewStrInst("object")
}
func (o *ObjectClass) Py__str__() *PyStrInst {
    return o.Py__repr__()
}
func (o *ObjectClass) Type() Type {return NONE}
func (o *ObjectClass) Py__getattribute__(attr *PyStrInst) Object {
    if attr.Value == "__new__" {
        return &PyNew{
            Func: o.Py__new__,
            }
    }
    return nil
}
func (o *ObjectClass) Py__setattr__(attr *PyStrInst, valObj Object) {}
func (o *ObjectClass) Py__new__(cls Class, objs ...Object) Object {
    return &PyInst{
        Py_class: cls,
        Py__dict__: map[string]Object{},
        }
}
func (o *ObjectClass) Py__init__(Object, ...Object) {}
func (o *ObjectClass) Py__name__() *PyStrInst {return NewStrInst("object")}
func (o *ObjectClass) Py__base__() Class {return nil}

var Py_object = &ObjectClass{}

type Pytype struct {}
func (t *Pytype) Py__class__() Class {return Py_type}
func (t *Pytype) Type() Type {return TYPE}
func (t *Pytype) Call(objs ...Object) Object {
    if len(objs) == 1 {
        return objs[0].Py__class__()
    }

    // then, o is supposed to be the name of class or cls
    return t.Py__new__(t, objs...)
}
func (t *Pytype) Py__call__(o Object, objs ...Object) Object {
    switch cls := o.(type) {
    case *Super:
        inst := cls.Py__new__(cls)
        return inst
    case Class:
        // temporary solution: context should be a stack
        context = cls
        inst := cls.Py__new__(cls, objs...)
        context = inst
        cls.Py__init__(inst, objs...)
        context = nil
        return inst
    }
    return nil
}
func (t *Pytype) Py__init__(Object, ...Object) {}
func (t *Pytype) Py__name__() *PyStrInst {return t.Py__repr__()}
func (t *Pytype) Py__base__() Class {return Py_object}
func (t *Pytype) Py__repr__() *PyStrInst {
    return NewStrInst("<class 'type'>")
}
func (t *Pytype) Py__str__() *PyStrInst {
    return t.Py__repr__()
}
func (t *Pytype) Py__new__(cls Class, objs ...Object) Object {
    var name *PyStrInst = objs[0].(*PyStrInst)
    var base *PyClass = objs[1].(*PyClass)
    var attrs *DictInst = objs[2].(*DictInst)

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
func (ni *NoneInst) Py__class__() Class {return nil}
func (ni *NoneInst) Py__repr__() *PyStrInst {
    return NewStrInst("None")
}
func (ni *NoneInst) Py__str__() *PyStrInst {
    return ni.Py__repr__()
}
func (ni *NoneInst) Type() Type {return NONE}
func (ni *NoneInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (ni *NoneInst) Py__setattr__(*PyStrInst, Object) {}

var Py_None = &NoneInst{
    Value: 0,
}

type BoolInst struct {
    Value   int
}
func (bi *BoolInst) Py__class__() Class {return nil}
func (bi *BoolInst) Type() Type {return BOOL}
func (bi *BoolInst) Py__repr__() *PyStrInst {
    if bi.Value == 1 {
        return NewStrInst("True")
    } else {
        return NewStrInst("False")
    }
}
func (bi *BoolInst) Py__str__() *PyStrInst {
    return bi.Py__repr__()
}
func (bi *BoolInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (bi *BoolInst) Py__setattr__(*PyStrInst, Object) {}

type Pyint struct {}
func (pi *Pyint) Py__class__() Class {return Py_object}
func (pi *Pyint) Py__init__(inst Object, objs ...Object) {}
func (pi *Pyint) Py__name__() *PyStrInst {return NewStrInst("int")}
func (pi *Pyint) Type() Type {return CLASS}
func (pi *Pyint) Py__getattribute__(attr *PyStrInst) Object {return nil}
func (pi *Pyint) Py__setattr__(attr *PyStrInst, val Object) {}
func (pi *Pyint) Py__new__(cls Class, objs ...Object) Object {
    inst := &IntegerInst{
        Py_class: cls,
    }

    if len(objs) == 1 {
        switch obj := objs[0].(type) {
        case *PyStrInst:
            v, _ := strconv.Atoi(obj.Value)
            inst.Value = v
        case *IntegerInst:
            inst.Value = obj.Value
        }
    } else if len(objs) == 0 {
        inst.Value = 0
    }

    return inst
}
func (pi *Pyint) Py__base__() Class {return Py_object}
func (pi *Pyint) Py__repr__() *PyStrInst {
    return NewStrInst(fmt.Sprint("<class 'int'>"))
}
func (pi *Pyint) Py__str__() *PyStrInst {
    return pi.Py__repr__()
}

var Py_int = &Pyint{}

type IntegerInst struct {
    Value int
    Py_class Class
}

func (i *IntegerInst) Py__class__() Class {return i.Py_class}
func (i *IntegerInst) Type() Type {return INSTANCE}
func (i *IntegerInst) Py__getattribute__(attr *PyStrInst) Object {return nil}
func (i *IntegerInst) Py__setattr__(attr *PyStrInst, val Object) {}
func (i *IntegerInst) Py__repr__() *PyStrInst {
    return NewStrInst(fmt.Sprint(i.Value))
}
func (i *IntegerInst) Py__str__() *PyStrInst {return i.Py__repr__()}

func NewInteger(n int) *IntegerInst {
    return &IntegerInst{
        Value: n,
        Py_class: Py_int,
    }
}

type PyStr struct {}
func (ps *PyStr) Py__class__() Class {return Py_type}
func (ps *PyStr) Py__init__(Object, ...Object) {}
func (ps *PyStr) Py__name__() *PyStrInst {return NewStrInst("str")}
func (ps *PyStr) Type() Type {return TYPE}
func (ps *PyStr) Py__getattribute__(attr *PyStrInst) Object {return nil}
func (ps *PyStr) Py__setattr__(*PyStrInst, Object) {}
func (ps *PyStr) Py__new__(cls Class, objs ...Object) Object {
    inst := &PyStrInst{
        Py_class: cls,
    }

    if len(objs) == 1 {
        switch obj := objs[0].(type) {
        case *PyStrInst:
            inst.Value = obj.Value
        case *IntegerInst:
            inst.Value = strconv.Itoa(obj.Value)
        }
    } else if len(objs) == 0 {
        inst.Value = ""
    }

    return inst

}
func (pc *PyStr) Py__base__() Class {return Py_object}
func (pc *PyStr) Py__repr__() *PyStrInst {
    return NewStrInst(fmt.Sprintf("<class '%s'>", "str"))
}
func (pc *PyStr) Py__str__() *PyStrInst {
    return pc.Py__repr__()
}

var Py_str = &PyStr{}

type PyStrInst struct {
    Value   string
    Py_class Class
}

func (si *PyStrInst) Py__class__() Class {return si.Py_class}
func (si *PyStrInst) Type() Type {return STRING}
func (si *PyStrInst) Py__repr__() *PyStrInst {
    return NewStrInst("'" + si.Value + "'")
}
func (si *PyStrInst) Py__str__() *PyStrInst {
    return NewStrInst(fmt.Sprint(si.Value))
}
func (si *PyStrInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (si *PyStrInst) Py__setattr__(*PyStrInst, Object) {}
func (si *PyStrInst) String() string {
    return fmt.Sprint(si.Value)
}

func NewStrInst(s string) *PyStrInst {
    return &PyStrInst{
        Value: s,
        Py_class: Py_str,
    }
}

type ListInst struct {
    Items   []Object
}

func (li *ListInst) Py__class__() Class {return nil}
func (li *ListInst) Type() Type {return LIST}
func (li *ListInst) Py__repr__() *PyStrInst {
    var s string
    s += "["
    if len(li.Items) > 0 {
        s += fmt.Sprintf("%v", li.Items[0].Py__repr__())
    }
    for _, item := range li.Items[1:] {
        s += ", "
        s += fmt.Sprintf("%v", item.Py__repr__())
    }
    s += "]"
    return NewStrInst(s)
}
func (li *ListInst) Py__str__() *PyStrInst {
    return li.Py__repr__()
}

func (li *ListInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (li *ListInst) Py__setattr__(*PyStrInst, Object) {}
func (li *ListInst) Py__len__() *IntegerInst {
    return &IntegerInst{
        Value: len(li.Items),
        Py_class: Py_int,
        }
}

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

    Map   map[PyStrInst]Object 
                            
}

func (di *DictInst) Py__class__() Class {return nil}
func (di *DictInst) Type() Type {return DICT}
func (di *DictInst) Py__repr__() *PyStrInst {
    var s string
    s += "{"
    var i = 0
    for k, v := range di.Map {
        s += fmt.Sprintf("%v:%v", k.Py__str__(), v.Py__str__())
        if i == len(di.Map) - 1 {
            break
        }
        s += ", "
        i++
    }
    s += "}"
    return NewStrInst(s)
}
func (di *DictInst) Py__str__() *PyStrInst {
    return di.Py__repr__()
}

func (di *DictInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (di *DictInst) Py__setattr__(*PyStrInst, Object) {}
func (di *DictInst) Py__len__() *IntegerInst {
    return &IntegerInst{
        Value: len(di.Map),
        Py_class: Py_int,
        }
}
func (di *DictInst) Py__getitem__(k Object) Object {
    switch obj := k.(type) {
    case *PyStrInst:
        return di.Map[*obj]
    }
    return nil
}
func (di *DictInst) Py__setitem__(k Object, v Object) {
    switch obj := k.(type) {
    case *PyStrInst:
        di.Map[*obj] = v
    default:
        panic("Do not yet support keys other than str")
    }
}

// TODO: rename instance
type FunctionInst struct {
    Name    *PyStrInst
    Params  []string
    Body    []ast.Statement
    env     *Environment
}

func (fi *FunctionInst) Py__class__() Class {return Py_Function}
func (fi *FunctionInst) Type() Type {return FUNCTION}
func (fi *FunctionInst) Py__call__(objs ...Object) Object {
    env := fi.env.DeriveEnv()
    for i := 0; i < len(fi.Params); i++ {
        env.Set(fi.Params[i], objs[i])
    }

    for _, stmt := range fi.Body {
        switch node := stmt.(type) {
        case *ast.ReturnStatement:
            return Eval(node.Value, env)
        default:
            Exec([]ast.Statement{stmt}, env)
        }
    }
    return Py_None
}
func (fi *FunctionInst) Py__name__() *PyStrInst {return fi.Name}
func (fi *FunctionInst) Py__repr__() *PyStrInst {
    return NewStrInst(fmt.Sprintf("<function %v at %p>", fi.Name.Py__str__(), fi))
}
func (fi *FunctionInst) Py__str__() *PyStrInst {
    return fi.Py__repr__()
}
func (fi *FunctionInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (fi *FunctionInst) Py__setattr__(*PyStrInst, Object) {}

type PyFunction struct {}
func (F *PyFunction) Py__class__() Class {return Py_type}
func (F *PyFunction) Py__init__(Object, ...Object) {}
func (F *PyFunction) Py__name__() *PyStrInst {return NewStrInst("function")}
func (F *PyFunction) Type() Type {return CLASS}
func (F *PyFunction) Py__getattribute__(attr *PyStrInst) Object {return nil}
func (F *PyFunction) Py__setattr__(attr *PyStrInst, val Object) {}
func (F *PyFunction) Py__new__(cls Class, objs ...Object) Object {return nil}
func (F *PyFunction) Py__base__() Class {return Py_object}
func (F *PyFunction) Py__repr__() *PyStrInst {return NewStrInst("<class 'function'>")}
func (F *PyFunction) Py__str__() *PyStrInst {return F.Py__repr__()}

var Py_Function = &PyFunction{}

type PyClass struct {
    Name        *PyStrInst
    Base        Class
    Dict        *DictInst    
}

func (pc *PyClass) Py__class__() Class {return Py_type}
func (pc *PyClass) Py__init__(o Object, objs ...Object) {
    if __init__ := pc.Dict.Py__getitem__(NewStrInst("__init__")); __init__ != nil {
        objs = append([]Object{o}, objs...)
        __init__.(Function).Py__call__(objs...)
    }
}
func (pc *PyClass) Py__name__() *PyStrInst {return pc.Name}
func (pc *PyClass) Type() Type {return CLASS}
func (pc *PyClass) Py__getattribute__(attr *PyStrInst) Object {
    // FIXME: defualt to object
    val := pc.Dict.Py__getitem__(attr)
    if val != nil {
        return val
    }

    if pc.Base == nil {
        return nil
    }

    return pc.Base.Py__getattribute__(attr)
}
func (pc *PyClass) Py__setattr__(attr *PyStrInst, val Object) {
    pc.Dict.Py__setitem__(attr, val)
}
func (pc *PyClass) Py__new__(cls Class, objs ...Object) Object {
    objs = append([]Object{cls}, objs...)
    if __new__ := pc.Dict.Py__getitem__(NewStrInst("__new__")); __new__ != nil {
        return __new__.(Function).Py__call__(objs...)
    }
    return pc.Base.Py__new__(cls, objs...)
}
func (pc *PyClass) Py__base__() Class {
    // FIXME: haven't execute customized __new__
    return pc.Base
}

func (pc *PyClass) Py__repr__() *PyStrInst {
    return NewStrInst(fmt.Sprintf("<class '%v'>", pc.Name.Py__str__()))
}
func (pc *PyClass) Py__str__() *PyStrInst {
    return pc.Py__repr__()
}

type BoundMethod struct {
    Func        *FunctionInst
    Inst        *PyInst
} 
func (bm *BoundMethod) Py__class__() Class {return Py_Function}
func (bm *BoundMethod) Py__repr__() *PyStrInst {
    s := fmt.Sprintf("<bound method %s.%s object at %p>",
            bm.Inst.Py__class__().Py__name__().Value, bm.Func.Name.Py__str__(), bm)

    return NewStrInst(s)
}
func (bm *BoundMethod) Py__str__() *PyStrInst {return bm.Py__repr__()}
func (bm *BoundMethod) Py__call__(objs ...Object) Object {
    objs = append([]Object{bm}, objs...)
    return bm.Func.Py__call__(objs...)
}
func (bm *BoundMethod) Py__name__() *PyStrInst {return bm.Func.Py__name__()}
func (bm *BoundMethod) Type() Type {return METHOD}
func (bm *BoundMethod) Py__getattribute__(*PyStrInst) Object {return nil}
func (bm *BoundMethod) Py__setattr__(*PyStrInst, Object) {}

type PyInst struct {
    Py_class Class
    Py__dict__ map[string]Object
}

func (pi *PyInst) Py__class__() Class {return pi.Py_class}
func (pi *PyInst) Type() Type {return INSTANCE}
func (pi *PyInst) Py__getattribute__(attr *PyStrInst) Object {
    targetObj, ok := pi.Py__dict__[attr.Value]
    if ok {
        return targetObj 
    }

    switch targetObj := pi.Py__class__().Py__getattribute__(attr).(type) {
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
    s := fmt.Sprintf("<%v objects at %p>", pi.Py__class__().Py__name__().Value, pi)
    return NewStrInst(s)
}
func (pi *PyInst) Py__str__() *PyStrInst {return pi.Py__repr__()}

// TODO: temporary
type Print struct {}
func (p *Print) Py__class__() Class {return Py_Function}
func (p *Print) Type() Type {return FUNCTION}
func (p *Print) Call() {}
func (p *Print) Py__repr__() *PyStrInst {return NewStrInst("print")}
func (p *Print) Py__str__() *PyStrInst {return p.Py__repr__()}
func (p *Print) Py__name__() *PyStrInst {return p.Py__repr__()}
func (p *Print) Py__call__(objs ...Object) Object {
    for _, obj := range objs {
        fmt.Print(obj.Py__str__().Value)
        fmt.Print(" ")
    }
    fmt.Println()
    return Py_None
}
func (p *Print) Py__getattribute__(*PyStrInst) Object {return nil}
func (p *Print) Py__setattr__(*PyStrInst, Object) {}

var Py_print = &Print{}

type Len struct {}
func (l *Len) Py__class__() Class {return Py_Function}
func (l *Len) Type() Type {return FUNCTION}
func (l *Len) Call() {}
func (l *Len) Py__repr__() *PyStrInst {return NewStrInst("<function 'len'>")}
func (l *Len) Py__str__() *PyStrInst {return l.Py__repr__()}
func (l *Len) Py__name__() *PyStrInst {return l.Py__repr__()}
func (l *Len) Py__call__(objs ...Object) Object {
    // supposed to have only one arguments
    obj := objs[0]

    switch o := obj.(type) {
    case *ListInst:
        return o.Py__len__()
    case *DictInst:
        return o.Py__len__()
    case *PyInst:
        // check customized __len__ of instance
        return nil
    default:
        // not supposed to run into here
        return nil
    }
}
func (l *Len) Py__getattribute__(*PyStrInst) Object {return nil}
func (l *Len) Py__setattr__(*PyStrInst, Object) {}

var Py_len = &Len{}

type Exception interface {
    Object
    ErrorMsg() string
}

type ExceptionInst struct {
    Msg string
}

func (ei *ExceptionInst) Py__class__() Class {return nil}
func (ei *ExceptionInst) Type() Type {return EXCEPTION}
func (ei *ExceptionInst) ErrorMsg() string {return ei.Msg}
func (ei *ExceptionInst) String() string {return "Exception"}
func (ei *ExceptionInst) Py__repr__() *PyStrInst {return NewStrInst("Exception")}
func (ei *ExceptionInst) Py__str__() *PyStrInst {return ei.Py__repr__()}
func (ei *ExceptionInst) Py__getattribute__(*PyStrInst) Object {return nil}
func (ei *ExceptionInst) Py__setattr__(*PyStrInst, Object) {}

type Super struct {}
func (s *Super) Py__class__() Class {return Py_type}
func (s *Super) Type() Type {return CLASS}
func (s *Super) Call() {}
func (s *Super) Py__getattribute__(*PyStrInst) Object {return nil}
func (s *Super) Py__setattr__(*PyStrInst, Object) {}
func (s *Super) Py__repr__() *PyStrInst {return NewStrInst("<class 'super'>")}
func (s *Super) Py__str__() *PyStrInst {return s.Py__repr__()}
func (s *Super) Py__init__(Object, ...Object) {}
func (s *Super) Py__new__(cls Class, objs ...Object) Object {
    switch o := context.(type) {
    case *PyInst:
        return &SuperInst{
            class: o.Py__class__(),
            inst: o,
        }
    case Class:
        return &SuperInst{
            class: o,
        }
    }
    return nil
}
func (s *Super) Py__name__() *PyStrInst {return NewStrInst("super")}
func (s *Super) Py__base__() Class {return Py_object}

var Py_super = &Super{}

type SuperInst struct {
    class Class
    inst *PyInst
}

func (si *SuperInst) Py__class__() Class {return Py_super}
func (si *SuperInst) Py__getattribute__(attr *PyStrInst) Object {
    if si.inst != nil {
        switch targetObj := si.class.Py__base__().Py__getattribute__(attr).(type) {
        case *FunctionInst:
            return &BoundMethod{
                Func: targetObj,
                Inst: si.inst,
                }
        }
    } else {
        return si.class.Py__base__().Py__getattribute__(attr)
    } 
    return nil
}

func (si *SuperInst) Py__repr__() *PyStrInst {
    return NewStrInst("<super object>")
}
func (si *SuperInst) Py__str__() *PyStrInst {
    return si.Py__repr__()
}
func (si *SuperInst) Type() Type {return SUPER}
func (si *SuperInst) Py__setattr__(*PyStrInst, Object) {}
