package evaluator

var (
    Py_True = &BoolInst{Value: 1}
    Py_False = &BoolInst{Value: 0}
)

var __builtins__ = map[string] Object{
    "object": Py_object,
    "True": Py_True,
    "False": Py_False,
    "None": Py_None,
    "print": Py_print,
    "len": Py_len,
    "super": Py_super,
    "type": Py_type,
}

type Environment struct {
    store     map[string]Object
    parent    *Environment
}

func NewEnvironment() *Environment {
    builtinsEnv := &Environment{
        store: __builtins__,
        parent: nil,
    }

    return &Environment{
        store: map[string]Object{},
        parent: builtinsEnv,
    }
}

func (self *Environment) Set(key string, value Object) {
    self.store[key] = value
}

func (self *Environment) Get(key string) Object {
    // omit the condition of key not being existing
    if self.parent == nil {
        return self.store[key]
    }

    if obj, ok := self.store[key]; ok {
        return obj
    } else {
        return self.parent.Get(key)
    }
}

func (self *Environment) DeriveEnv() *Environment {
    return &Environment{
        store: map[string] Object{},
        parent: self,
    }
}

func (self *Environment) Store() map[string]Object {
    return self.store
}

