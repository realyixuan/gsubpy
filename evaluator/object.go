/*
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

*/

package evaluator

import (
    "fmt"
    "strconv"
    "unsafe"
    "crypto/sha1"
    "encoding/binary"
    
    "github.com/realyixuan/gsubpy/ast"
)

var __name__ = newStringInst("__name__")
var __new__ = newStringInst("__new__")
var __init__ = newStringInst("__init__")
var __repr__ = newStringInst("__repr__")
var __str__ = newStringInst("__str__")
var __eq__ = newStringInst("__eq__")
var __lt__ = newStringInst("__lt__")
var __gt__ = newStringInst("__gt__")
var __hash__ = newStringInst("__hash__")
var __getattribute__ = newStringInst("__getattribute__")
var __setattr__ = newStringInst("__setattr__")
var __call__ = newStringInst("__call__")
var __len__ = newStringInst("__len__")
var __bool__ = newStringInst("__bool__")

var __getitem__ = newStringInst("__getitem__")

var __iter__ = newStringInst("__iter__")
var __next__ = newStringInst("__next__")

type Object interface {
    Type()          Class
    Id()            int64
    Attr(*StringInst)    Object
    attrs()         *DictInst
}

type Class interface {
    Object
    Base() Class
}

type Function interface {
    Object
    Call(...Object) Object
}

type objectData struct {
    d       *DictInst
}
func (o *objectData) attrs() *DictInst { return o.d }

var Pyobject__new__ = newBuiltinFunc(
    __new__,
    func(objs ...Object) Object {
        cls := objs[0].(Class)
        return newPyInst(cls)
    },
)

var Pyobject__init__ = newBuiltinFunc(
    __init__,
    func(objs ...Object) Object {return nil},
)

var Pyobject__repr__ = newBuiltinFunc(
    __repr__,
    func(objs ...Object) Object {
        self := objs[0]
        s := fmt.Sprintf("<%v object at 0x%x>", self.Type().Attr(__name__).(*StringInst).Value, self.Id())
        return newStringInst(s)
    },
)

var Pyobject__str__ = newBuiltinFunc(
    __str__,
    func(objs ...Object) Object {
        self := objs[0]

        var s string
        switch self.(type) {
        case Class:
            s = fmt.Sprintf("<class '%v'>", self.Attr(__name__).(*StringInst).Value)
        default:
            s = Pyobject__repr__.Call(self).(*StringInst).Value
        }

        return newStringInst(s)
    },
)

var Pyobject__eq__ = newBuiltinFunc(__eq__,
    func(objs ...Object) Object {
        self := objs[0]
        other := objs[1]
        if self.Id() == other.Id() {
            return Py_True
        } else {
            return Py_False
        }
    },
)

var Pyobject__lt__ = newBuiltinFunc(__lt__,
    func(objs ...Object) Object {return nil},
)


var Pyobject__gt__ = newBuiltinFunc(__gt__,
    func(objs ...Object) Object {return nil},
)

var Pyobject__hash__ = newBuiltinFunc(__hash__,
    func(objs ...Object) Object {
        self := objs[0]
        id := self.Id()

        buf := make([]byte, 8)
        binary.LittleEndian.PutUint64(buf, uint64(id))

        return newIntegerInst(hash(buf))
    },
)

var Pyobject__getattribute__ = newBuiltinFunc(__getattribute__,
    func(objs ...Object) Object {
        self := objs[0]
        name := objs[1]
        attr := attrFromAll(self, name.(*StringInst))
        return attr
    },
)

var Pyobject__setattr__ = newBuiltinFunc(__setattr__,
    func(objs ...Object) Object {
        self := objs[0]
        name := objs[1]
        val := objs[2]
        self.attrs().Set(name, val)
        return Py_None
    },
)

type Pyobject struct {
    *objectData
}

func newPyobject() *Pyobject {
    o := &Pyobject{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (po *Pyobject) Type() Class { return Py_type }
func (po *Pyobject) Base() Class { return nil }
func (po *Pyobject) Id() int64 { return int64(uintptr(unsafe.Pointer(po))) }
func (po *Pyobject) Attr(name *StringInst) Object { return po.d.Get(name) }

var Py_object = newPyobject()
func init() {
    Py_object.attrs().Set(__name__, newStringInst("object"))
    Py_object.attrs().Set(__hash__, Pyobject__hash__)
    Py_object.attrs().Set(__new__, Pyobject__new__)
    Py_object.attrs().Set(__init__, Pyobject__init__)
    Py_object.attrs().Set(__repr__, Pyobject__repr__)
    Py_object.attrs().Set(__str__, Pyobject__str__)
    Py_object.attrs().Set(__eq__, Pyobject__eq__)
    Py_object.attrs().Set(__lt__, Pyobject__lt__)
    Py_object.attrs().Set(__gt__, Pyobject__gt__)
    Py_object.attrs().Set(__getattribute__, Pyobject__getattribute__)
    Py_object.attrs().Set(__setattr__, Pyobject__setattr__)
}

type Pytype struct {
    *objectData
}

func newPytype() *Pytype {
    o := &Pytype{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (pt *Pytype) Type() Class { return Py_type }
func (pt *Pytype) Base() Class { return Py_object }
func (pt *Pytype) Id() int64 { return int64(uintptr(unsafe.Pointer(pt))) }
func (pt *Pytype) Attr(name *StringInst) Object { return Getattr(pt, name) }

var Py_type = newPytype()
func init()  {
    Py_type.attrs().Set(__name__, newStringInst("type"))
    Py_type.attrs().Set(__new__, newBuiltinFunc(__new__,
            func(objs ...Object) Object {
                var name *StringInst = objs[0].(*StringInst)
                var base *Pyclass = objs[1].(*Pyclass)
                var attrs *DictInst = objs[2].(*DictInst)
                return newPyclass(name, base, attrs)
            },
        ),
    )

    Py_type.attrs().Set(__init__, newBuiltinFunc(__init__,
            func(objs ...Object) Object {
                return Py_None
            },
        ),
    )

    Py_type.attrs().Set(__call__, newBuiltinFunc(__call__,
            func(objs ...Object) Object {
                cls := objs[0]
                switch cls {
                case Py_type:
                    return objs[1].Type()
                default:
                    self := Call(attrItself(cls, __new__), objs...)

                    args := append([]Object{self}, objs[1:]...)
                    Call(attrItself(cls, __init__), args...)
                    return self
                }
            },
        ),
    )

    Py_type.attrs().Set(__getattribute__, newBuiltinFunc(__getattribute__,
            func(objs ...Object) Object {
                cls := objs[0]
                name := objs[1]
                attr := attrFromAll(cls, name.(*StringInst))
                return attr
            },
        ),
    )
}

type Pyclass struct {
    *objectData
    base    Class
    name    *StringInst
}

func newPyclass(
    name    *StringInst,
    base    Class,
    dict    *DictInst,
) *Pyclass {
    o := &Pyclass{
        objectData: &objectData{
            d: dict,
        },
        name: name,
        base: base,
    }
    o.init()
    return o
}

func (pc *Pyclass) init() {
    pc.attrs().Set(__name__, pc.name)
}

func (pc *Pyclass) Type() Class { return Py_type }
func (pc *Pyclass) Base() Class { return pc.base }
func (pc *Pyclass) Id() int64 { return int64(uintptr(unsafe.Pointer(pc))) }
func (pc *Pyclass) Attr(name *StringInst) Object { return Getattr(pc, name) }

type PyInst struct {
    *objectData
    class   Class
}

func newPyInst(cls Class) *PyInst {
    return &PyInst{
        objectData: &objectData{
            d: newDictInst(),
        },
        class: cls,
    }
}

func (i *PyInst) Type() Class { return i.class }
func (i *PyInst) Id() int64 { return int64(uintptr(unsafe.Pointer(i))) }
func (i *PyInst) Attr(name *StringInst) Object { return Getattr(i, name) }

type PyFunction struct {
    *objectData
}

func newPyFunction() *PyFunction {
    o := &PyFunction{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (pf *PyFunction) Type() Class { return Py_type }
func (pf *PyFunction) Base() Class { return Py_object }
func (pf *PyFunction) Id() int64 { return int64(uintptr(unsafe.Pointer(pf))) }
func (pf *PyFunction) Attr(name *StringInst) Object { return Getattr(pf, name) }

var Py_function = newPyFunction()
func init() {
    Py_function.attrs().Set(__name__, newStringInst("function"))
    Py_function.attrs().Set(__call__, newBuiltinFunc(__call__,
            func(objs ...Object) Object {
                self := objs[0]
                return self.(Function).Call(objs[1:]...)
            },
        ),
    )
}

var PyBuiltinFunction__call__ = newBuiltinFunc(__call__,
    func(objs ...Object) Object {
        self := objs[0]
        return self.(Function).Call(objs[1:]...)
    },
)

type PyBuiltinFunction struct {
    *objectData
}

func newPyBuiltinFunction() *PyBuiltinFunction {
    o := &PyBuiltinFunction{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (bf *PyBuiltinFunction) Type() Class { return Py_type }
func (bf *PyBuiltinFunction) Base() Class { return Py_object }
func (bf *PyBuiltinFunction) Id() int64 { return int64(uintptr(unsafe.Pointer(bf))) }
func (bf *PyBuiltinFunction) Attr(name *StringInst) Object { return Getattr(bf, name) }

var Py_builtin_function = newPyBuiltinFunction()
func init() {
    Py_builtin_function.attrs().Set(__name__, newStringInst("builtin_function"))
    Py_builtin_function.attrs().Set(__call__, PyBuiltinFunction__call__)
}

type FunctionInst struct {
    *objectData
    Name    *StringInst
    Params  []*StringInst
    Body    []ast.Statement
    env     *Environment
}

func newFunctionInst(
    name    *StringInst,
    params  []*StringInst,
    body    []ast.Statement,
    env     *Environment,
) *FunctionInst {
    o := &FunctionInst{
        objectData: &objectData{
            d: newDictInst(),
            },
        Name: name,
        Params: params,
        Body: body,
        env: env,
    }
    o.init()
    return o
}

func (f *FunctionInst) init() {
    f.attrs().Set(__name__, f.Name)
}

func (f *FunctionInst) Call(objs ...Object) Object {
    env := f.env.DeriveEnv()
    for i := 0; i < len(f.Params); i++ {
        env.Set(f.Params[i], objs[i])
    }
    rv, _ := Exec(f.Body, env)
    return rv
}

func (f *FunctionInst) Type() Class { return Py_function }
func (f *FunctionInst) Id() int64 { return int64(uintptr(unsafe.Pointer(f))) }
func (f *FunctionInst) Attr(name *StringInst) Object { return Getattr(f, name) }

type builtinFn func(...Object) Object

type BuiltinFunctionInst struct {
    *objectData
    name    *StringInst
    gfunc   builtinFn
}

func newBuiltinFunc(name *StringInst, f builtinFn) *BuiltinFunctionInst {
    return &BuiltinFunctionInst{
        objectData: &objectData{
            d: newDictInst(),
        },
        name: name,
        gfunc: f,
    }
}

func (f *BuiltinFunctionInst) Call(objs ...Object) Object { return f.gfunc(objs...) }
func (f *BuiltinFunctionInst) Type() Class { return Py_builtin_function }
func (f *BuiltinFunctionInst) Id() int64 { return int64(uintptr(unsafe.Pointer(f))) }
func (f *BuiltinFunctionInst) Attr(name *StringInst) Object { return Getattr(f, name) }

var Py_print = newBuiltinFunc(
    newStringInst("print"),
    func(objs ...Object) Object {
        for _, obj := range objs {
            fmt.Print(StringOf(obj))
            fmt.Print(" ")
        }
        fmt.Println()
        return Py_None
    },
)

var Py_len = newBuiltinFunc(
    newStringInst("len"),
    func(objs ...Object) Object {
        lenFn := attrItself(objs[0].Type(), __len__)
        return Call(lenFn, objs[0])
    },
)

var Py_hash = newBuiltinFunc(
    newStringInst("hash"),
    func(objs ...Object) Object {
        hashFn := attrItself(objs[0].Type(), __hash__).(Function)
        return Call(hashFn, objs[0])
    },
)

var Py_iter = newBuiltinFunc(
    newStringInst("iter"),
    func(objs ...Object) Object {
        iterFn := attrItself(objs[0].Type(), __iter__)
        return Call(iterFn, objs[0])
    },
)

var Py_next = newBuiltinFunc(
    newStringInst("next"),
    func(objs ...Object) Object {
        nextFn := attrItself(objs[0].Type(), __next__)
        return Call(nextFn, objs[0])
    },
)

type MethodInst struct {
    *objectData
    inst       Object
    f          Function
} 

func newMethod(inst Object, f Function) *MethodInst {
    return &MethodInst{
        objectData: &objectData{
            d: newDictInst(),
            },
        f: f,
        inst: inst,
    }
}

func (m *MethodInst) Type() Class { return Py_function }
func (m *MethodInst) Id() int64 { return int64(uintptr(unsafe.Pointer(m))) }
func (m *MethodInst) Attr(name *StringInst) Object { return Getattr(m, name) }
func (m *MethodInst) Call(objs ...Object) Object {
    objs = append([]Object{m.inst}, objs...)
    return Call(m.f, objs...)
}

type PyNoneType struct {
    *objectData
}

func newPyNoneType() *PyNoneType {
    o := &PyNoneType{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (pn *PyNoneType) Type() Class { return Py_type }
func (pn *PyNoneType) Base() Class { return Py_object }
func (pn *PyNoneType) Id() int64 { return int64(uintptr(unsafe.Pointer(pn))) }
func (pn *PyNoneType) Attr(name *StringInst) Object { return Getattr(pn, name) }

var Py_NoneType = newPyNoneType()
func init() {
    Py_NoneType.attrs().Set(__name__, newStringInst("NoneType"))
    Py_NoneType.attrs().Set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                return newStringInst("None")
            },
        ),
    )
}

type PyNone struct {
    *objectData
}

func newPyNone() *PyNone {
    return &PyNone{
        objectData: &objectData{},
    }
}

func (n *PyNone) Type() Class { return Py_NoneType }
func (n *PyNone) Id() int64 { return int64(uintptr(unsafe.Pointer(n))) }
func (n *PyNone) Attr(name *StringInst) Object { return Getattr(n, name) }

var Py_None = newPyNone()

type Pyint struct {
    *objectData
}

func newPyint() *Pyint {
    o := &Pyint{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (pi *Pyint) Type() Class { return Py_type }
func (pi *Pyint) Base() Class { return Py_object }
func (pi *Pyint) Id() int64 { return int64(uintptr(unsafe.Pointer(pi))) }
func (pi *Pyint) Attr(name *StringInst) Object { return Getattr(pi, name) }

var Py_int = newPyint()
func init() {
    Py_int.attrs().Set(__name__, newStringInst("int"))

    Py_int.attrs().Set(__new__, newBuiltinFunc(__new__,
            func(objs ...Object) Object {
                if len(objs[1:]) == 0 {
                    return newIntegerInst(int64(0))
                }

                switch o := objs[1].(type) {
                case *StringInst:
                    v, _ := strconv.Atoi(o.Value)
                    return newIntegerInst(int64(v))
                case *IntegerInst:
                    return o
                }

                return nil
            },
        ),
    )

    Py_int.attrs().Set(__hash__, newBuiltinFunc(__hash__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    Py_int.attrs().Set(__repr__, newBuiltinFunc(__repr__,
            func(objs ...Object) Object {
                v := strconv.FormatInt(objs[0].(*IntegerInst).Value, 10)
                return newStringInst(v)
            },
        ),
    )

    Py_int.attrs().Set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                v := strconv.FormatInt(objs[0].(*IntegerInst).Value, 10)
                return newStringInst(v)
            },
        ),
    )

    Py_int.attrs().Set(__eq__, newBuiltinFunc(__eq__,
            func(objs ...Object) Object {
                v1 := objs[0].(*IntegerInst).Value
                v2 := objs[1].(*IntegerInst).Value
                
                if v1 == v2 {
                    return Py_True
                } else {
                    return Py_False
                }
            },
        ),
    )

    Py_int.attrs().Set(__gt__, newBuiltinFunc(__gt__,
            func(objs ...Object) Object {
                v1 := objs[0].(*IntegerInst).Value
                v2 := objs[1].(*IntegerInst).Value
                
                if v1 > v2 {
                    return Py_True
                } else {
                    return Py_False
                }

            },
        ),
    )

    Py_int.attrs().Set(__lt__, newBuiltinFunc(__lt__,
            func(objs ...Object) Object {
                v1 := objs[0].(*IntegerInst).Value
                v2 := objs[1].(*IntegerInst).Value
                
                if v1 < v2 {
                    return Py_True
                } else {
                    return Py_False
                }
            },
        ),
    )

    Py_int.attrs().Set(__bool__, newBuiltinFunc(__bool__,
            func(objs ...Object) Object {
                self := objs[0].(*IntegerInst)
                if self.Value == 0 {
                    return Py_False
                } else {
                    return Py_True
                }
            },
        ),
    )
}


type IntegerInst struct {
    *objectData
    class   Class
    Value   int64
}

func newIntegerInst(v int64) *IntegerInst {
    return &IntegerInst{
        objectData: &objectData{
            d: newDictInst(),
        },
        class: Py_int,
        Value: v,
    }
}

func (i *IntegerInst) Type() Class { return i.class }
func (i *IntegerInst) Id() int64 { return int64(uintptr(unsafe.Pointer(i))) }
func (i *IntegerInst) Attr(name *StringInst) Object { return Getattr(i, name) }

type Pystr_iterator struct {
    *objectData
}

func newPystr_iterator() *Pystr_iterator {
    o := &Pystr_iterator{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    o.init()
    return o
}

func (psi *Pystr_iterator) init() {
    psi.attrs().Set(__name__, newStringInst("str_iterator"))

    psi.attrs().Set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    psi.attrs().Set(__next__, newBuiltinFunc(__next__,
            func(objs ...Object) Object {
                self := objs[0].(*StringIteratorInst)
                if self.idx >= Call(Py_len, self.stringInst).(*IntegerInst).Value {
                    // TODO: replace it with StopIteration Exception
                    return nil
                }
                self.idx += 1
                return typeCall(__getitem__, self.stringInst, newIntegerInst(self.idx-1))
            },
        ),
    )

}

func (psi *Pystr_iterator) Type() Class { return Py_type }
func (psi *Pystr_iterator) Base() Class { return Py_object }
func (psi *Pystr_iterator) Id() int64 { return int64(uintptr(unsafe.Pointer(psi))) }
func (psi *Pystr_iterator) Attr(name *StringInst) Object { return Getattr(psi, name) }

var Py_str_iterator = newPystr_iterator()

type StringIteratorInst struct {
    *objectData
    idx     int64
    stringInst    *StringInst
}

func newStringIteratorInst(t *StringInst) *StringIteratorInst {
    return &StringIteratorInst{
        objectData: &objectData{d: newDictInst()},
        idx: 0,
        stringInst: t,
    }
}

func (lsi *StringIteratorInst) Type() Class { return Py_str_iterator }
func (lsi *StringIteratorInst) Id() int64 { return int64(uintptr(unsafe.Pointer(lsi))) }
func (lsi *StringIteratorInst) Attr(name *StringInst) Object { return Getattr(lsi, name) }


var Pystr__hash__ = newBuiltinFunc(__hash__,
    func(objs ...Object) Object {
        strInst := objs[0].(*StringInst)
        return newIntegerInst(hash([]byte(strInst.Value)))
    },
)

var Pystr__eq__ = newBuiltinFunc(__eq__,
    func(objs ...Object) Object {
        if s2, ok := objs[1].(*StringInst); ok {
            s1 := objs[0].(*StringInst)
            if s1.Value == s2.Value {
                return Py_True
            }
        }

        return Py_False
    },
)

type Pystr struct {
    *objectData
}

func newPystr() *Pystr {
    o := &Pystr{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (ps *Pystr) Type() Class { return Py_type }
func (ps *Pystr) Base() Class { return Py_object }
func (ps *Pystr) Id() int64 { return int64(uintptr(unsafe.Pointer(ps))) }
func (ps *Pystr) Attr(name *StringInst) Object { return Getattr(ps, name) }

var Py_str = newPystr()
func init() {
    Py_str.attrs().Set(__name__, newStringInst("str"))

    Py_str.attrs().Set(__new__, newBuiltinFunc(__new__,
            func(objs ...Object) Object {
                if len(objs[1:]) == 0 {
                    return newStringInst("")
                }

                switch o := objs[1].(type) {
                case *IntegerInst:
                    v := strconv.FormatInt(o.Value, 10)
                    return newStringInst(v)
                case *StringInst:
                    return o
                }

                return nil
            },
        ),
    )

    Py_str.attrs().Set(__hash__, Pystr__hash__)
    Py_str.attrs().Set(__eq__, Pystr__eq__)

    Py_str.attrs().Set(__repr__, newBuiltinFunc(__repr__,
            func(objs ...Object) Object {
                return newStringInst("'" + objs[0].(*StringInst).Value + "'")
            },
        ),
    )

    Py_str.attrs().Set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    Py_str.attrs().Set(__gt__, newBuiltinFunc(__gt__,
            func(objs ...Object) Object {
                s1 := objs[0].(*StringInst)
                s2 := objs[1].(*StringInst)

                if s1.Value > s2.Value {
                    return Py_True
                } else {
                    return Py_False
                }
            },
        ),
    )

    Py_str.attrs().Set(__lt__, newBuiltinFunc(__lt__,
            func(objs ...Object) Object {
                s1 := objs[0].(*StringInst)
                s2 := objs[1].(*StringInst)

                if s1.Value < s2.Value {
                    return Py_True
                } else {
                    return Py_False
                }
            },
        ),
    )

    Py_str.attrs().Set(__len__, newBuiltinFunc(__len__,
            func(objs ...Object) Object {
                return newIntegerInst(int64(len(objs[0].(*StringInst).Value)))
            },
        ),
    )

    Py_str.attrs().Set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                self := objs[0].(*StringInst)
                return newStringIteratorInst(self)
            },
        ),
    )

    Py_str.attrs().Set(__getitem__, newBuiltinFunc(__getitem__,
            func(objs ...Object) Object {
                self := objs[0].(*StringInst)
                idx := objs[1].(*IntegerInst)
                return newStringInst(string(self.Value[idx.Value]))
            },
        ),
    )

}

type StringInst struct {
    *objectData
    Value   string
}

func newStringInst(s string) *StringInst {
    return &StringInst{
        objectData: &objectData{
            d: newDictInst(),
        },
        Value: s,
    }
}

func (s *StringInst) Type() Class { return Py_str }
func (s *StringInst) Id() int64 { return int64(uintptr(unsafe.Pointer(s))) }
func (s *StringInst) Attr(name *StringInst) Object { return Getattr(s, name) }
func (s *StringInst) String() string { return s.Value }

type Pybool struct {
    *objectData
}

func newPybool() *Pybool {
    o := &Pybool{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (pb *Pybool) Type() Class { return Py_type }
func (pb *Pybool) Base() Class { return Py_int }
func (pb *Pybool) Id() int64 { return int64(uintptr(unsafe.Pointer(pb))) }
func (pb *Pybool) Attr(name *StringInst) Object { return Getattr(pb, name) }

var Py_bool = newPybool()
func init() {
    Py_bool.attrs().Set(__new__, newBuiltinFunc(__new__,
            func(objs ...Object) Object {
                if boolFn := attrItself(objs[1].Type(), __bool__); boolFn != nil {
                    return Call(boolFn, objs[1])
                } else if lenFn := attrItself(objs[1].Type(), __len__); lenFn != nil {
                    l := Call(lenFn, objs[1]).(*IntegerInst)
                    if l.Value != 0 {
                        return Py_True
                    } else {
                        return Py_False
                    }
                }
                return Py_True
            },
        ),
    )

    Py_bool.attrs().Set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                o := objs[0].(*IntegerInst)
                if o.Value == 1 {
                    return newStringInst("True")
                } else {
                    return newStringInst("False")
                }
            },
        ),
    )
}

var Py_True = &IntegerInst{
    objectData: &objectData{
        d: newDictInst(),
    },
    class: Py_bool,
    Value: 1,
}

var Py_False = &IntegerInst{
    objectData: &objectData{
        d: newDictInst(),
    },
    class: Py_bool,
    Value: 0,
}


type Pylist_iterator struct {
    *objectData
}

func newPylist_iterator() *Pylist_iterator {
    o := &Pylist_iterator{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    o.init()
    return o
}

func (pli *Pylist_iterator) init() {
    pli.attrs().Set(__name__, newStringInst("list_iterator"))

    pli.attrs().Set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    pli.attrs().Set(__next__, newBuiltinFunc(__next__,
            func(objs ...Object) Object {
                self := objs[0].(*ListIteratorInst)
                if self.idx >= Call(Py_len, self.listInst).(*IntegerInst).Value {
                    // TODO: replace it with StopIteration Exception
                    return nil
                }
                self.idx += 1
                return typeCall(__getitem__, self.listInst, newIntegerInst(self.idx-1))
            },
        ),
    )

}

func (pli *Pylist_iterator) Type() Class { return Py_type }
func (pli *Pylist_iterator) Base() Class { return Py_object }
func (pli *Pylist_iterator) Id() int64 { return int64(uintptr(unsafe.Pointer(pli))) }
func (pli *Pylist_iterator) Attr(name *StringInst) Object { return Getattr(pli, name) }

var Py_list_iterator = newPylist_iterator()

type ListIteratorInst struct {
    *objectData
    idx     int64
    listInst    *ListInst
}

func newListIteratorInst(t *ListInst) *ListIteratorInst {
    return &ListIteratorInst{
        objectData: &objectData{d: newDictInst()},
        idx: 0,
        listInst: t,
    }
}

func (lii *ListIteratorInst) Type() Class { return Py_list_iterator }
func (lii *ListIteratorInst) Id() int64 { return int64(uintptr(unsafe.Pointer(lii))) }
func (lii *ListIteratorInst) Attr(name *StringInst) Object { return Getattr(lii, name) }

type Pylist struct {
    *objectData
}

func newPylist() *Pylist {
    return &Pylist{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
}

func (pl *Pylist) Type() Class { return Py_type }
func (pl *Pylist) Base() Class { return Py_object }
func (pl *Pylist) Id() int64 { return int64(uintptr(unsafe.Pointer(pl))) }
func (pl *Pylist) Attr(name *StringInst) Object { return Getattr(pl, name) }

var Py_list = newPylist()
func init() {
    Py_list.attrs().Set(__name__, newStringInst("list"))

    Py_list.attrs().Set(__len__, newBuiltinFunc(__len__,
            func(objs ...Object) Object {
                return newIntegerInst(int64(len(objs[0].(*ListInst).items)))
            },
        ),
    )

    Py_list.attrs().Set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                li := objs[0].(*ListInst)

                var s string
                s += "["
                if len(li.items) > 0 {
                    strFn := attrItself(li.items[0].Type(), __str__).(Function)
                    s += fmt.Sprintf("%v", Call(strFn, li.items[0]))
                }
                for _, item := range li.items[1:] {
                    s += ", "
                    strFn := attrItself(item.Type(), __str__).(Function)
                    s += fmt.Sprintf("%v", Call(strFn, item))
                }
                s += "]"
                return newStringInst(s)
            },
        ),
    )

    Py_list.attrs().Set(__getitem__, newBuiltinFunc(__getitem__,
            func(objs ...Object) Object {
                self := objs[0].(*ListInst)
                idx := objs[1].(*IntegerInst)
                return self.items[idx.Value]
            },
        ),
    )

    Py_list.attrs().Set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                self := objs[0].(*ListInst)
                return newListIteratorInst(self)
            },
        ),
    )

}

type ListInst struct {
    *objectData
    items []Object
}

func newListInst() *ListInst {
    return &ListInst{
        objectData: &objectData{
            d: newDictInst(),
        },
        items: []Object{},
    }
}

func (l *ListInst) Type() Class { return Py_list }
func (l *ListInst) Id() int64 { return int64(uintptr(unsafe.Pointer(l))) }
func (l *ListInst) Attr(name *StringInst) Object { return Getattr(l, name) }

func hash(bv []byte) int64 {
    // sha1 just fine
    bs := sha1.Sum(bv)
    var v int64
    for i := 0; i < 8; i++ {
        v += int64(bs[i]) << (8*i)
    }
    return v
}

type Pydict struct {
    *objectData
}

func newPydict() *Pydict {
    return &Pydict{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
}

func (pd *Pydict) Type() Class { return Py_type }
func (pd *Pydict) Base() Class { return Py_object }
func (pd *Pydict) Id() int64 { return int64(uintptr(unsafe.Pointer(pd))) }
func (pd *Pydict) Attr(name *StringInst) Object { return Getattr(pd, name) }

var Py_dict = newPydict()

type pair struct {
    Key     Object
    Value   Object
}

type DictInst struct {
    d *DictInst
    store map[int64][]*pair
}

func newDictInst() *DictInst {
    return &DictInst{
        d: &DictInst{
            store: map[int64][]*pair{},
            },
        store: map[int64][]*pair{},
    }

}

func (d *DictInst) attrs() *DictInst { return d.d }
func (d *DictInst) Type() Class { return Py_dict }
func (d *DictInst) Id() int64 { return int64(uintptr(unsafe.Pointer(d))) }
func (d *DictInst) Attr(name *StringInst) Object { return Getattr(d, name) }
func (d *DictInst) Get(key Object) Object {
    var hashVal Object
    switch key.(type) {
    case *StringInst:
        hashVal = Pystr__hash__.Call(key)
    default:
        hashVal = typeCall(__hash__, key)
    }
    if pairs, ok := d.store[hashVal.(*IntegerInst).Value]; ok {
        for _, pair := range pairs {
            var isMatch Object
            switch pair.Key.(type) {
            case *StringInst:
                isMatch = Pystr__eq__.Call(pair.Key, key)
            default:
                isMatch = typeCall(__eq__, pair.Key, key)
            }

            if isMatch == Py_True {
                return pair.Value
            }

        }
    }

    return nil
}

func (d *DictInst) Set(key Object, val Object) {
    var hashVal Object
    switch key.(type) {
    case *StringInst:
        hashVal = Pystr__hash__.Call(key)
    default:
        hashVal = typeCall(__hash__, key)
    }

    var flag bool = false

    for _, pair := range d.store[hashVal.(*IntegerInst).Value] {
        var isMatch Object
        switch pair.Key.(type) {
        case *StringInst:
            isMatch = Pystr__eq__.Call(pair.Key, key)
        default:
            isMatch = typeCall(__eq__, pair.Key, key)
        }

        if isMatch == Py_True {
            pair.Value = val
            flag = true
        }
    }

    if flag == false {
        d.store[hashVal.(*IntegerInst).Value] = append(d.store[hashVal.(*IntegerInst).Value], &pair{Key: key, Value: val})
    }
}

type Pyrange_iterator struct {
    *objectData
}

func newPyrange_iterator() *Pyrange_iterator {
    o := &Pyrange_iterator{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    o.init()
    return o
}

func (pri *Pyrange_iterator) init() {
    pri.attrs().Set(__name__, newStringInst("range_iterator"))

    pri.attrs().Set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    pri.attrs().Set(__next__, newBuiltinFunc(__next__,
            func(objs ...Object) Object {
                self := objs[0].(*RangeIteratorInst)
                if self.rangeInst.step > 0 {
                    if self.curV >= self.rangeInst.end {
                        return nil
                    }
                } else if self.rangeInst.step < 0 {
                    if self.curV <= self.rangeInst.end {
                        return nil
                    }
                }
                res := self.curV
                self.curV += self.rangeInst.step
                return newIntegerInst(res)
            },
        ),
    )
}

func (pri *Pyrange_iterator) Type() Class { return Py_type }
func (pri *Pyrange_iterator) Base() Class { return Py_object }
func (pri *Pyrange_iterator) Id() int64 { return int64(uintptr(unsafe.Pointer(pri))) }
func (pri *Pyrange_iterator) Attr(name *StringInst) Object { return Getattr(pri, name) }

var Py_range_iterator = newPyrange_iterator()

type RangeIteratorInst struct {
    *objectData
    curV        int64
    rangeInst    *RangeInst
}

func newRangeIteratorInst(t *RangeInst) *RangeIteratorInst {
    return &RangeIteratorInst{
        objectData: &objectData{d: newDictInst()},
        curV: t.start,
        rangeInst: t,
    }
}

func (rii *RangeIteratorInst) Type() Class { return Py_range_iterator }
func (rii *RangeIteratorInst) Id() int64 { return int64(uintptr(unsafe.Pointer(rii))) }
func (rii *RangeIteratorInst) Attr(name *StringInst) Object { return Getattr(rii, name) }

type Pyrange struct {
    *objectData
}

func newPyrange() *Pyrange {
    o := &Pyrange{objectData: &objectData{d: newDictInst()}}
    o.init()
    return o
}

func (pr *Pyrange) init() {
    pr.attrs().Set(__name__, newStringInst("range"))

    pr.attrs().Set(__new__, newBuiltinFunc(__new__,
            func (objs ...Object) Object {
                if len(objs) == 2 {
                    return newRangeInst(
                        0,
                        objs[1].(*IntegerInst).Value,
                        1,
                        )
                } else if len(objs) == 3 {
                    return newRangeInst(
                        objs[1].(*IntegerInst).Value,
                        objs[2].(*IntegerInst).Value,
                        1,
                        )
                } else if len(objs) == 4 {
                    return newRangeInst(
                        objs[1].(*IntegerInst).Value,
                        objs[2].(*IntegerInst).Value,
                        objs[3].(*IntegerInst).Value,
                        )
                }

                return nil
            },
        ),
    )

    pr.attrs().Set(__iter__, newBuiltinFunc(__iter__,
            func (objs ...Object) Object {
                return newRangeIteratorInst(objs[0].(*RangeInst))
            },
        ),
    )

}

func (pr *Pyrange) Type() Class { return Py_type }
func (pr *Pyrange) Base() Class { return Py_object }
func (pr *Pyrange) Id() int64 { return int64(uintptr(unsafe.Pointer(pr))) }
func (pr *Pyrange) Attr(name *StringInst) Object { return Getattr(pr, name) }

var Py_range = newPyrange()

type RangeInst struct {
    *objectData
    start   int64
    end     int64
    step    int64
}

func newRangeInst(start, end, step int64) *RangeInst {
    return &RangeInst{
        objectData: &objectData{d: newDictInst()},
        start: start,
        end: end,
        step: step,
    }
}

func (ri *RangeInst) Type() Class { return Py_range }
func (ri *RangeInst) Id() int64 { return int64(uintptr(unsafe.Pointer(ri))) }
func (ri *RangeInst) Attr(name *StringInst) Object { return Getattr(ri, name) }


type PyException struct {
    *objectData
}

func (pe *PyException) Type() Class { return Py_type }
func (pe *PyException) Id() int64 { return int64(uintptr(unsafe.Pointer(pe))) }
func (pe *PyException) Base() Class { return Py_object }
func (pe *PyException) Attr(name *StringInst) Object { return Getattr(pe, name) }

var Py_Exception = &PyException{
    &objectData{
        d: newDictInst(),
    },
}
func init() {
    Py_Exception.attrs().Set(__name__, newStringInst("Exception"))
}

type ExceptionInst struct {
    *objectData
    Payload     Object
}

func newExceptionInst(obj Object) *ExceptionInst {
    return &ExceptionInst{
        objectData: &objectData{
            d: newDictInst(),
        },
        Payload: obj,
    }
}

func Error(s string) *ExceptionInst {
    return newExceptionInst(newStringInst(s))
}

func (e *ExceptionInst) Type() Class { return Py_Exception }
func (e *ExceptionInst) Id() int64 { return int64(uintptr(unsafe.Pointer(e))) }
func (e *ExceptionInst) Attr(name *StringInst) Object { return Getattr(e, name) }

func Getattr(obj Object, name *StringInst) Object {
    __getattribute__ := attrFromAll(obj.Type(), __getattribute__).(Function)
    return Call(__getattribute__, obj, name)
}

func attrItself(obj Object, name *StringInst) Object {
    switch obj.(type) {
    case Class:
        cls := obj.(Class)
        for c := cls; c != nil; c = c.Base() {
            if rv := c.attrs().Get(name); rv != nil {
                return rv
            }
        }
    default:
        if rv := obj.attrs().Get(name); rv != nil {
            return rv
        }
    }

    return nil
}

func attrFromAll(obj Object, name *StringInst) Object {
    if rv := attrItself(obj, name); rv != nil {
        return rv
    }

    if rv := attrItself(obj.Type(), name); rv != nil {
        switch v := rv.(type) {
        case Function:
            return newMethod(obj, v)
        default:
            return rv
        }
    }

    return nil
}

func typeCall(attrName *StringInst, obj Object, args ...Object) Object {
    attr := attrItself(obj.Type(), attrName)
    if attr == nil {
        panic(Error(fmt.Sprintf("%v object is not callable", StringOf(obj.Type()))))
    }

    fn, ok := attr.(Function) 
    if ok == false {
        panic(Error(fmt.Sprintf("%v object is not callable", StringOf(attr.Type()))))
    }
    
    args = append([]Object{obj}, args...)
    return Call(fn, args...)
}

func Call(obj Object, args ...Object) Object {
    __call__Fn := attrItself(obj.Type(), __call__)
    args = append([]Object{obj}, args...)

    if __call__Fn != PyBuiltinFunction__call__ {
        return Call(__call__Fn, args...)
    } else {
        return __call__Fn.(Function).Call(args...)
    }
}

