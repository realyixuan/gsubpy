package evaluator

var __builtins__ = map[string] Object{
    "object": Py_object,
    "True": Py_True,
    "False": Py_False,
    "None": Py_None,
    "print": Py_print,
    "type": Py_type,
    "int": Py_int,
    "str": Py_str,
    "len": Py_len,
    "bool": Py_bool,
    "hash": Py_hash,
    "Exception": Py_Exception,

    "iter": Py_iter,
    "next": Py_next,
    "range": Py_range,
}

type Environment struct {
    store     *DictInst
    parent    *Environment
}

func NewEnvironment() *Environment {
    builtinsEnv := &Environment{
        store: newDictInst(),
        parent: nil,
    }

    for k, v := range __builtins__ {
        builtinsEnv.Set(newStringInst(k), v)
    }

    return &Environment{
        store: newDictInst(),
        parent: builtinsEnv,
    }
}

func (e *Environment) SetFromString(key string, value Object) {
    e.Set(newStringInst(key), value)
}

func (self *Environment) Set(key Object, value Object) {
    self.store.Set(key, value)
}

func (e *Environment) GetFromString(key string) Object {
    return e.Get(newStringInst(key))
}

func (self *Environment) Get(key Object) Object {
    if self.parent == nil {
        val := self.store.Get(key)
        return val
    }

    if obj := self.store.Get(key); obj != nil {
        return obj
    } else {
        return self.parent.Get(key)
    }
}

func (self *Environment) DeriveEnv() *Environment {
    return &Environment{
        store: newDictInst(),
        parent: self,
    }
}

func (self *Environment) Store() *DictInst {
    return self.store
}

