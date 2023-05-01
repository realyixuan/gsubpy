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
var __setitem__ = newStringInst("__setitem__")
var __contains__ = newStringInst("__contains__")

var __iter__ = newStringInst("__iter__")
var __next__ = newStringInst("__next__")

var __add__ = newStringInst("__add__")
var __sub__ = newStringInst("__sub__")
var __mul__ = newStringInst("__mul__")
var __floordiv__ = newStringInst("__floordiv__")

type Object interface {
    otype()          Class
    id()            int64
    attrs()         *DictInst
}

type Class interface {
    Object
    cbase() Class
}

type Function interface {
    Object
    call(...Object) Object
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
        s := fmt.Sprintf("<%v object at 0x%x>",
            op_GETATTR(self.otype(), __name__), self.id())
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
            s = fmt.Sprintf("<class '%v'>", op_GETATTR(self.otype(), __name__))
        default:
            s = Pyobject__repr__.call(self).(*StringInst).Value
        }

        return newStringInst(s)
    },
)

var Pyobject__eq__ = newBuiltinFunc(__eq__,
    func(objs ...Object) Object {
        self := objs[0]
        other := objs[1]
        if self.id() == other.id() {
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
        id := self.id()

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
        name := objs[1].(*StringInst)
        val := objs[2]
        self.attrs().set(name, val)
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

func (po *Pyobject) otype() Class { return Py_type }
func (po *Pyobject) cbase() Class { return nil }
func (po *Pyobject) id() int64 { return int64(uintptr(unsafe.Pointer(po))) }

var Py_object = newPyobject()
func init() {
    Py_object.attrs().set(__name__, newStringInst("object"))
    Py_object.attrs().set(__hash__, Pyobject__hash__)
    Py_object.attrs().set(__new__, Pyobject__new__)
    Py_object.attrs().set(__init__, Pyobject__init__)
    Py_object.attrs().set(__repr__, Pyobject__repr__)
    Py_object.attrs().set(__str__, Pyobject__str__)
    Py_object.attrs().set(__eq__, Pyobject__eq__)
    Py_object.attrs().set(__lt__, Pyobject__lt__)
    Py_object.attrs().set(__gt__, Pyobject__gt__)
    Py_object.attrs().set(__getattribute__, Pyobject__getattribute__)
    Py_object.attrs().set(__setattr__, Pyobject__setattr__)
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

func (pt *Pytype) otype() Class { return Py_type }
func (pt *Pytype) cbase() Class { return Py_object }
func (pt *Pytype) id() int64 { return int64(uintptr(unsafe.Pointer(pt))) }

var Py_type = newPytype()
func init()  {
    Py_type.attrs().set(__name__, newStringInst("type"))
    Py_type.attrs().set(__new__, newBuiltinFunc(__new__,
            func(objs ...Object) Object {
                if len(objs[1:]) == 1 {
                    return objs[1].otype()
                }
                var name *StringInst = objs[1].(*StringInst)
                var base = objs[2].(Class)
                var attrs = (objs[3]).(*DictInst)
                return newPyclass(name, base, attrs)
            },
        ),
    )

    Py_type.attrs().set(__init__, newBuiltinFunc(__init__,
            func(objs ...Object) Object {
                return Py_None
            },
        ),
    )

    Py_type.attrs().set(__call__, newBuiltinFunc(__call__,
            func(objs ...Object) Object {
                cls := objs[0]
                self := op_CALL(attrItself(cls, __new__), objs...)
                args := append([]Object{self}, objs[1:]...)
                op_CALL(attrItself(cls, __init__), args...)
                return self
            },
        ),
    )

    Py_type.attrs().set(__getattribute__, newBuiltinFunc(__getattribute__,
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
    pc.attrs().set(__name__, pc.name)
}

func (pc *Pyclass) otype() Class { return Py_type }
func (pc *Pyclass) cbase() Class { return pc.base }
func (pc *Pyclass) id() int64 { return int64(uintptr(unsafe.Pointer(pc))) }

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

func (i *PyInst) otype() Class { return i.class }
func (i *PyInst) id() int64 { return int64(uintptr(unsafe.Pointer(i))) }

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

func (pf *PyFunction) otype() Class { return Py_type }
func (pf *PyFunction) cbase() Class { return Py_object }
func (pf *PyFunction) id() int64 { return int64(uintptr(unsafe.Pointer(pf))) }

var Py_function = newPyFunction()
func init() {
    Py_function.attrs().set(__name__, newStringInst("function"))
    Py_function.attrs().set(__call__, newBuiltinFunc(__call__,
            func(objs ...Object) Object {
                self := objs[0]
                return self.(Function).call(objs[1:]...)
            },
        ),
    )
}

var PyBuiltinFunction__call__ = newBuiltinFunc(__call__,
    func(objs ...Object) Object {
        self := objs[0]
        return self.(Function).call(objs[1:]...)
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

func (bf *PyBuiltinFunction) otype() Class { return Py_type }
func (bf *PyBuiltinFunction) cbase() Class { return Py_object }
func (bf *PyBuiltinFunction) id() int64 { return int64(uintptr(unsafe.Pointer(bf))) }

var Py_builtin_function = newPyBuiltinFunction()
func init() {
    Py_builtin_function.attrs().set(__name__, newStringInst("builtin_function"))
    Py_builtin_function.attrs().set(__call__, PyBuiltinFunction__call__)
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
    f.attrs().set(__name__, f.Name)
}

func (f *FunctionInst) call(objs ...Object) Object {
    env := f.env.DeriveEnv()
    for i := 0; i < len(f.Params); i++ {
        env.Set(f.Params[i], objs[i])
    }
    rv, _ := Exec(f.Body, env)
    return rv
}

func (f *FunctionInst) otype() Class { return Py_function }
func (f *FunctionInst) id() int64 { return int64(uintptr(unsafe.Pointer(f))) }

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

func (f *BuiltinFunctionInst) call(objs ...Object) Object { return f.gfunc(objs...) }
func (f *BuiltinFunctionInst) otype() Class { return Py_builtin_function }
func (f *BuiltinFunctionInst) id() int64 { return int64(uintptr(unsafe.Pointer(f))) }

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
        lenFn := attrItself(objs[0].otype(), __len__)
        return op_CALL(lenFn, objs[0])
    },
)

var Py_hash = newBuiltinFunc(
    newStringInst("hash"),
    func(objs ...Object) Object {
        hashFn := attrItself(objs[0].otype(), __hash__).(Function)
        return op_CALL(hashFn, objs[0])
    },
)

var Py_id = newBuiltinFunc(
    newStringInst("id"),
    func(objs ...Object) Object {
        return newIntegerInst(objs[0].id())
    },
)

var Py_iter = newBuiltinFunc(
    newStringInst("iter"),
    func(objs ...Object) Object {
        iterFn := attrItself(objs[0].otype(), __iter__)
        return op_CALL(iterFn, objs[0])
    },
)

var Py_next = newBuiltinFunc(
    newStringInst("next"),
    func(objs ...Object) Object {
        nextFn := attrItself(objs[0].otype(), __next__)
        return op_CALL(nextFn, objs[0])
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

func (m *MethodInst) otype() Class { return Py_function }
func (m *MethodInst) id() int64 { return int64(uintptr(unsafe.Pointer(m))) }
func (m *MethodInst) call(objs ...Object) Object {
    objs = append([]Object{m.inst}, objs...)
    return op_CALL(m.f, objs...)
}

type PyNoneotype struct {
    *objectData
}

func newPyNoneotype() *PyNoneotype {
    o := &PyNoneotype{
        objectData: &objectData{
            d: newDictInst(),
        },
    }
    return o
}

func (pn *PyNoneotype) otype() Class { return Py_type }
func (pn *PyNoneotype) cbase() Class { return Py_object }
func (pn *PyNoneotype) id() int64 { return int64(uintptr(unsafe.Pointer(pn))) }

var Py_Noneotype = newPyNoneotype()
func init() {
    Py_Noneotype.attrs().set(__name__, newStringInst("Noneotype"))
    Py_Noneotype.attrs().set(__str__, newBuiltinFunc(__str__,
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

func (n *PyNone) otype() Class { return Py_Noneotype }
func (n *PyNone) id() int64 { return int64(uintptr(unsafe.Pointer(n))) }

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

func (pi *Pyint) otype() Class { return Py_type }
func (pi *Pyint) cbase() Class { return Py_object }
func (pi *Pyint) id() int64 { return int64(uintptr(unsafe.Pointer(pi))) }

var Py_int = newPyint()
func init() {
    Py_int.attrs().set(__name__, newStringInst("int"))

    Py_int.attrs().set(__new__, newBuiltinFunc(__new__,
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

    Py_int.attrs().set(__hash__, newBuiltinFunc(__hash__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    Py_int.attrs().set(__repr__, newBuiltinFunc(__repr__,
            func(objs ...Object) Object {
                v := strconv.FormatInt(objs[0].(*IntegerInst).Value, 10)
                return newStringInst(v)
            },
        ),
    )

    Py_int.attrs().set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                v := strconv.FormatInt(objs[0].(*IntegerInst).Value, 10)
                return newStringInst(v)
            },
        ),
    )

    Py_int.attrs().set(__eq__, newBuiltinFunc(__eq__,
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

    Py_int.attrs().set(__gt__, newBuiltinFunc(__gt__,
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

    Py_int.attrs().set(__lt__, newBuiltinFunc(__lt__,
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

    Py_int.attrs().set(__bool__, newBuiltinFunc(__bool__,
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

    Py_int.attrs().set(__add__, newBuiltinFunc(__add__,
            func(objs ...Object) Object {
                if objs[0].otype() != objs[1].otype() {
                    panic(Error("otypeError: two different types"))
                }

                self, other := objs[0].(*IntegerInst), objs[1].(*IntegerInst)
                return newIntegerInst(self.Value + other.Value)
            },
        ),
    )

    Py_int.attrs().set(__sub__, newBuiltinFunc(__sub__,
            func(objs ...Object) Object {
                self, other := objs[0].(*IntegerInst), objs[1].(*IntegerInst)
                return newIntegerInst(self.Value - other.Value)
            },
        ),
    )

    Py_int.attrs().set(__mul__, newBuiltinFunc(__mul__,
            func(objs ...Object) Object {
                self, other := objs[0].(*IntegerInst), objs[1].(*IntegerInst)
                return newIntegerInst(self.Value * other.Value)
            },
        ),
    )

    Py_int.attrs().set(__floordiv__, newBuiltinFunc(__floordiv__,
            func(objs ...Object) Object {
                self, other := objs[0].(*IntegerInst), objs[1].(*IntegerInst)

                if other.Value == 0 {
                    panic(Error("ZeroDivisionError: division by zero"))
                }

                return newIntegerInst(self.Value / other.Value)
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

func (i *IntegerInst) otype() Class { return i.class }
func (i *IntegerInst) id() int64 { return int64(uintptr(unsafe.Pointer(i))) }

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
    psi.attrs().set(__name__, newStringInst("str_iterator"))

    psi.attrs().set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    psi.attrs().set(__next__, newBuiltinFunc(__next__,
            func(objs ...Object) Object {
                self := objs[0].(*StringIteratorInst)
                if self.idx >= op_CALL(Py_len, self.stringInst).(*IntegerInst).Value {
                    // TODO: replace it with StopIteration Exception
                    return nil
                }
                self.idx += 1
                return op_SUBSCR_GET(self.stringInst, newIntegerInst(self.idx-1))
            },
        ),
    )

}

func (psi *Pystr_iterator) otype() Class { return Py_type }
func (psi *Pystr_iterator) cbase() Class { return Py_object }
func (psi *Pystr_iterator) id() int64 { return int64(uintptr(unsafe.Pointer(psi))) }

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

func (lsi *StringIteratorInst) otype() Class { return Py_str_iterator }
func (lsi *StringIteratorInst) id() int64 { return int64(uintptr(unsafe.Pointer(lsi))) }


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

func (ps *Pystr) otype() Class { return Py_type }
func (ps *Pystr) cbase() Class { return Py_object }
func (ps *Pystr) id() int64 { return int64(uintptr(unsafe.Pointer(ps))) }

var Py_str = newPystr()
func init() {
    Py_str.attrs().set(__name__, newStringInst("str"))

    Py_str.attrs().set(__new__, newBuiltinFunc(__new__,
            func(objs ...Object) Object {
                if len(objs[1:]) == 0 {
                    return newStringInst("")
                }

                o := objs[1]

                return typeCall(__str__, o)
            },
        ),
    )

    Py_str.attrs().set(__hash__, Pystr__hash__)
    Py_str.attrs().set(__eq__, Pystr__eq__)

    Py_str.attrs().set(__repr__, newBuiltinFunc(__repr__,
            func(objs ...Object) Object {
                return newStringInst("'" + objs[0].(*StringInst).Value + "'")
            },
        ),
    )

    Py_str.attrs().set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    Py_str.attrs().set(__gt__, newBuiltinFunc(__gt__,
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

    Py_str.attrs().set(__lt__, newBuiltinFunc(__lt__,
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

    Py_str.attrs().set(__len__, newBuiltinFunc(__len__,
            func(objs ...Object) Object {
                return newIntegerInst(int64(len(objs[0].(*StringInst).Value)))
            },
        ),
    )

    Py_str.attrs().set(__contains__, newBuiltinFunc(__contains__,
            func(objs ...Object) Object {
                self := objs[0].(*StringInst)
                item := objs[1]

                for _, ch := range self.Value {
                    chStr := newStringInst(string(ch))
                    if op_EQ(chStr, item) == Py_True {
                        return Py_True
                    }
                }

                return Py_False
            },
        ),
    )

    Py_str.attrs().set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                self := objs[0].(*StringInst)
                return newStringIteratorInst(self)
            },
        ),
    )

    Py_str.attrs().set(__getitem__, newBuiltinFunc(__getitem__,
            func(objs ...Object) Object {
                self := objs[0].(*StringInst)
                idx := objs[1].(*IntegerInst)
                return newStringInst(string(self.Value[idx.Value]))
            },
        ),
    )

    Py_str.attrs().set(__add__, newBuiltinFunc(__add__,
            func(objs ...Object) Object {
                if objs[0].otype() != objs[1].otype() {
                    panic(Error("otypeError: two different types"))
                }

                self, other := objs[0].(*StringInst), objs[1].(*StringInst)
                return newStringInst(self.Value + other.Value)
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

func (s *StringInst) otype() Class { return Py_str }
func (s *StringInst) id() int64 { return int64(uintptr(unsafe.Pointer(s))) }
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

func (pb *Pybool) otype() Class { return Py_type }
func (pb *Pybool) cbase() Class { return Py_int }
func (pb *Pybool) id() int64 { return int64(uintptr(unsafe.Pointer(pb))) }

var Py_bool = newPybool()
func init() {
    Py_bool.attrs().set(__new__, newBuiltinFunc(__new__,
            func(objs ...Object) Object {
                if boolFn := attrItself(objs[1].otype(), __bool__); boolFn != nil {
                    return op_CALL(boolFn, objs[1])
                } else if lenFn := attrItself(objs[1].otype(), __len__); lenFn != nil {
                    l := op_CALL(lenFn, objs[1]).(*IntegerInst)
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

    Py_bool.attrs().set(__str__, newBuiltinFunc(__str__,
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
    pli.attrs().set(__name__, newStringInst("list_iterator"))

    pli.attrs().set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    pli.attrs().set(__next__, newBuiltinFunc(__next__,
            func(objs ...Object) Object {
                self := objs[0].(*ListIteratorInst)
                if self.idx >= op_CALL(Py_len, self.listInst).(*IntegerInst).Value {
                    // TODO: replace it with StopIteration Exception
                    return nil
                }
                self.idx += 1
                return op_SUBSCR_GET(self.listInst, newIntegerInst(self.idx-1))
            },
        ),
    )

}

func (pli *Pylist_iterator) otype() Class { return Py_type }
func (pli *Pylist_iterator) cbase() Class { return Py_object }
func (pli *Pylist_iterator) id() int64 { return int64(uintptr(unsafe.Pointer(pli))) }

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

func (lii *ListIteratorInst) otype() Class { return Py_list_iterator }
func (lii *ListIteratorInst) id() int64 { return int64(uintptr(unsafe.Pointer(lii))) }

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

func (pl *Pylist) otype() Class { return Py_type }
func (pl *Pylist) cbase() Class { return Py_object }
func (pl *Pylist) id() int64 { return int64(uintptr(unsafe.Pointer(pl))) }

var Py_list = newPylist()
func init() {
    Py_list.attrs().set(__name__, newStringInst("list"))

    Py_list.attrs().set(__len__, newBuiltinFunc(__len__,
            func(objs ...Object) Object {
                return newIntegerInst(int64(len(objs[0].(*ListInst).items)))
            },
        ),
    )

    Py_list.attrs().set(__str__, newBuiltinFunc(__str__,
            func(objs ...Object) Object {
                li := objs[0].(*ListInst)

                var s string
                s += "["
                if len(li.items) > 0 {
                    s += fmt.Sprintf("%v", op_CALL(Py_str, li.items[0]))
                }
                for _, item := range li.items[1:] {
                    s += fmt.Sprintf(", %v", op_CALL(Py_str, item))
                }
                s += "]"

                return newStringInst(s)
            },
        ),
    )

    Py_list.attrs().set(__getitem__, newBuiltinFunc(__getitem__,
            func(objs ...Object) Object {
                self := objs[0].(*ListInst)
                idx := objs[1].(*IntegerInst)
                return self.items[idx.Value]
            },
        ),
    )

    Py_list.attrs().set(__contains__, newBuiltinFunc(__contains__,
            func(objs ...Object) Object {
                self := objs[0].(*ListInst)
                item := objs[1]
                for val := op_CALL(Py_next, self); val != nil; {
                    if op_EQ(val, item) == Py_True {
                        return Py_True
                    }
                    val = op_CALL(Py_next, self)
                }
                return Py_False
            },
        ),
    )

    Py_list.attrs().set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                self := objs[0].(*ListInst)
                return newListIteratorInst(self)
            },
        ),
    )

    Py_list.attrs().set(__add__, newBuiltinFunc(__add__,
            func(objs ...Object) Object {
                self, other := objs[0].(*ListInst), objs[1].(*ListInst)
                li := newListInst()
                li.items = append(self.items, other.items...)
                return li
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

func (l *ListInst) otype() Class { return Py_list }
func (l *ListInst) id() int64 { return int64(uintptr(unsafe.Pointer(l))) }

type Pydict_keyiterator struct {
    *objectData
}

func newPydict_keyiterator() *Pydict_keyiterator {
    o := &Pydict_keyiterator{
        objectData: &objectData{d: newDictInst()},
    }
    o.init()
    return o
}

func (dki *Pydict_keyiterator) init() {
    dki.attrs().set(__name__, newStringInst("dict_keyiterator"))

    dki.attrs().set(__iter__, newBuiltinFunc(__iter__, 
            func (objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    dki.attrs().set(__next__, newBuiltinFunc(__next__, 
            func (objs ...Object) Object {
                self := objs[0].(*DictKeyiteratorInst)

                if self.idx >= int64(len(self.keys)) {
                    return nil
                }

                res := self.keys[self.idx]
                self.idx += 1
                return res
            },
        ),
    )

}

func (dki *Pydict_keyiterator) otype() Class { return Py_type }
func (dki *Pydict_keyiterator) cbase() Class { return Py_object }
func (dki *Pydict_keyiterator) id() int64 { return int64(uintptr(unsafe.Pointer(dki))) }

var Py_dict_keyiterator = newPydict_keyiterator()

type DictKeyiteratorInst struct {
    *objectData
    idx     int64
    keys    []Object
}

func newDictKeyiteratorInst(t *DictInst) *DictKeyiteratorInst {
    // XXX: no way remembering location of keys,
    //  so have to put keys in list initially
    var dks []Object
    for _, v := range t.store {
        for _, pair := range v {
            dks = append(dks, pair.Key)
        }
    }
    return &DictKeyiteratorInst{
        objectData: &objectData{d: newDictInst()},
        keys: dks,
        idx: 0,
    }
}

func (ki *DictKeyiteratorInst) otype() Class { return Py_dict_keyiterator }
func (ki *DictKeyiteratorInst) id() int64 { return int64(uintptr(unsafe.Pointer(ki))) }

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

func (pd *Pydict) otype() Class { return Py_type }
func (pd *Pydict) cbase() Class { return Py_object }
func (pd *Pydict) id() int64 { return int64(uintptr(unsafe.Pointer(pd))) }

var Py_dict = newPydict()
func init() {
    Py_dict.attrs().set(__name__, newStringInst("dict"))

    Py_dict.attrs().set(__str__, newBuiltinFunc(__str__,
            func (objs ...Object) Object {
                d := objs[0].(*DictInst)

                var s string
                s += "{"
                for _, pairs := range d.store {
                    for _, pair := range pairs {
                        s += fmt.Sprintf("%v: %v, ",
                            op_CALL(Py_str, pair.Key),
                            op_CALL(Py_str, pair.Value),
                            )
                    }
                }
                s += "}"

                return newStringInst(s)
            },
        ),
    )

    Py_dict.attrs().set(__iter__, newBuiltinFunc(__iter__,
            func (objs ...Object) Object {
                self := objs[0].(*DictInst)
                return newDictKeyiteratorInst(self)
            },
        ),
    )

    Py_dict.attrs().set(__getitem__, newBuiltinFunc(__getitem__,
            func (objs ...Object) Object {
                self, key := objs[0].(*DictInst), objs[1]
                hashVal := op_CALL(Py_hash, key)
                if pairs, ok := self.store[hashVal.(*IntegerInst).Value]; ok {
                    for _, pair := range pairs {
                        if op_EQ(pair.Key, key) == Py_True {
                            return pair.Value
                        }
                    }
                }

                return nil
            },
        ),
    )

    Py_dict.attrs().set(__setitem__, newBuiltinFunc(__getitem__,
            func (objs ...Object) Object {
                self, key, val := objs[0].(*DictInst), objs[1], objs[2]
                hashVal := op_CALL(Py_hash, key)

                var flag bool = false
                for _, pair := range self.store[hashVal.(*IntegerInst).Value] {
                    if op_EQ(pair.Key, key) == Py_True {
                        pair.Value = val
                        flag = true
                    }
                }

                if !flag {
                    self.store[hashVal.(*IntegerInst).Value] = append(self.store[hashVal.(*IntegerInst).Value], &pair{Key: key, Value: val})
                }

                return Py_None
            },
        ),
    )

    Py_dict.attrs().set(__contains__, newBuiltinFunc(__contains__,
            func (objs ...Object) Object {
                self, item := objs[0].(*DictInst), objs[1]
                for val := op_CALL(Py_next, self); val != nil; {
                    if op_EQ(val, item) == Py_True {
                        return Py_True
                    }
                    val = op_CALL(Py_next, self)
                }
                return Py_False
            },
        ),
    )

}

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
func (d *DictInst) otype() Class { return Py_dict }
func (d *DictInst) id() int64 { return int64(uintptr(unsafe.Pointer(d))) }

func (d *DictInst) get(key *StringInst) Object {
    hashVal := Pystr__hash__.call(key)
    if pairs, ok := d.store[hashVal.(*IntegerInst).Value]; ok {
        for _, pair := range pairs {
            if Pystr__eq__.call(key, pair.Key) == Py_True {
                return pair.Value
            }
        }
    }

    return nil
}

func (d *DictInst) set(key *StringInst, val Object) {
    hashVal := Pystr__hash__.call(key)

    var flag bool = false
    for _, pair := range d.store[hashVal.(*IntegerInst).Value] {
        if Pystr__eq__.call(key, pair.Key) == Py_True {
            pair.Value = val
            flag = true
        }
    }

    if !flag {
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
    pri.attrs().set(__name__, newStringInst("range_iterator"))

    pri.attrs().set(__iter__, newBuiltinFunc(__iter__,
            func(objs ...Object) Object {
                return objs[0]
            },
        ),
    )

    pri.attrs().set(__next__, newBuiltinFunc(__next__,
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

func (pri *Pyrange_iterator) otype() Class { return Py_type }
func (pri *Pyrange_iterator) cbase() Class { return Py_object }
func (pri *Pyrange_iterator) id() int64 { return int64(uintptr(unsafe.Pointer(pri))) }

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

func (rii *RangeIteratorInst) otype() Class { return Py_range_iterator }
func (rii *RangeIteratorInst) id() int64 { return int64(uintptr(unsafe.Pointer(rii))) }

type Pyrange struct {
    *objectData
}

func newPyrange() *Pyrange {
    o := &Pyrange{objectData: &objectData{d: newDictInst()}}
    o.init()
    return o
}

func (pr *Pyrange) init() {
    pr.attrs().set(__name__, newStringInst("range"))

    pr.attrs().set(__new__, newBuiltinFunc(__new__,
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

    pr.attrs().set(__iter__, newBuiltinFunc(__iter__,
            func (objs ...Object) Object {
                return newRangeIteratorInst(objs[0].(*RangeInst))
            },
        ),
    )

}

func (pr *Pyrange) otype() Class { return Py_type }
func (pr *Pyrange) cbase() Class { return Py_object }
func (pr *Pyrange) id() int64 { return int64(uintptr(unsafe.Pointer(pr))) }

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

func (ri *RangeInst) otype() Class { return Py_range }
func (ri *RangeInst) id() int64 { return int64(uintptr(unsafe.Pointer(ri))) }


type PyException struct {
    *objectData
}

func (pe *PyException) otype() Class { return Py_type }
func (pe *PyException) id() int64 { return int64(uintptr(unsafe.Pointer(pe))) }
func (pe *PyException) cbase() Class { return Py_object }

var Py_Exception = &PyException{
    &objectData{
        d: newDictInst(),
    },
}
func init() {
    Py_Exception.attrs().set(__name__, newStringInst("Exception"))
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

func (e *ExceptionInst) otype() Class { return Py_Exception }
func (e *ExceptionInst) id() int64 { return int64(uintptr(unsafe.Pointer(e))) }

func hash(bv []byte) int64 {
    // sha1 just fine
    bs := sha1.Sum(bv)
    var v int64
    for i := 0; i < 8; i++ {
        v += int64(bs[i]) << (8*i)
    }
    return v
}

func Getattr(obj Object, name *StringInst) Object {
    __getattribute__ := attrFromAll(obj.otype(), __getattribute__).(Function)
    return op_CALL(__getattribute__, obj, name)
}

func attrItself(obj Object, name *StringInst) Object {
    switch obj.(type) {
    case Class:
        cls := obj.(Class)
        for c := cls; c != nil; c = c.cbase() {
            if rv := c.attrs().get(name); rv != nil {
                return rv
            }
        }
    default:
        if rv := obj.attrs().get(name); rv != nil {
            return rv
        }
    }

    return nil
}

func attrFromAll(obj Object, name *StringInst) Object {
    if rv := attrItself(obj, name); rv != nil {
        return rv
    }

    if rv := attrItself(obj.otype(), name); rv != nil {
        switch v := rv.(type) {
        case Function:
            return newMethod(obj, v)
        default:
            return rv
        }
    }

    return nil
}
