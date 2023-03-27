package evaluator

import (
    "strconv"
    "gsubpy/ast"
    "gsubpy/object"
    "gsubpy/token"
)

// simplest way for now
var context object.Object

var (
    True = &object.BoolInst{Value: 1}
    False = &object.BoolInst{Value: 0}
)

var __builtins__ = map[string]object.Object{
    "object": object.PyObject,
    "True": True,
    "False": False,
    "None": object.Py_None,
    "print": object.Py_print,
    "super": &object.BuiltinClass{
        Name: "super",
        },
}

type Environment struct {
    store     map[string]object.Object
    parent    *Environment
}

func NewEnvironment() *Environment {
    builtinsEnv := &Environment{
        store: __builtins__,
        parent: nil,
    }

    return &Environment{
        store: map[string]object.Object{},
        parent: builtinsEnv,
    }
}

func (self *Environment) Set(key string, value object.Object) {
    self.store[key] = value
}

func (self *Environment) Get(key string) object.Object {
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

func (self *Environment) deriveEnv() *Environment {
    return &Environment{
        store: map[string]object.Object{},
        parent: self,
    }
}

func Exec(stmts []ast.Statement, env *Environment) *object.NoneInst {
    for _, stmt := range stmts {
        switch node := stmt.(type) {
        case *ast.AssignStatement:
            execAssignStatement(node, env)
        case *ast.IfStatement:
            execIfStatement(node, env)
        case *ast.WhileStatement:
            execWhileStatement(node, env)
        case *ast.DefStatement:
            execDefStatement(node, env)
        case *ast.ClassStatement:
            execClassStatement(node, env)
        case *ast.ExpressionStatement:
            Eval(node, env)
        }
    }

    return object.Py_None
}

func Eval(expression ast.Expression, env *Environment) object.Object {
    switch node := expression.(type) {
    case *ast.IdentifierExpression:
        return env.Get(node.Identifier.Literals)
    case *ast.PlusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if leftObj.Type() != rightObj.Type() {
            panic(&object.ExceptionInst{Msg: "TypeError: two different types"})
        }

        switch leftObj.(type) {
        case *object.IntegerInst:
            return &object.IntegerInst{
                Value: leftObj.(*object.IntegerInst).Value + rightObj.(*object.IntegerInst).Value,
                }
        case *object.PyStrInst:
            return &object.PyStrInst{
                Value: leftObj.(*object.PyStrInst).Value + rightObj.(*object.PyStrInst).Value,
                }
        }

    case *ast.MinusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return &object.IntegerInst{
            Value: leftObj.(*object.IntegerInst).Value - rightObj.(*object.IntegerInst).Value,
            }
    case *ast.MulExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return &object.IntegerInst{
            Value: leftObj.(*object.IntegerInst).Value * rightObj.(*object.IntegerInst).Value,
            }
    case *ast.DivideExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if rightObj.(*object.IntegerInst).Value == 0 {
            panic(&object.ExceptionInst{Msg: "ZeroDivisionError: division by zero"})
        }

        return &object.IntegerInst{
            Value: leftObj.(*object.IntegerInst).Value / rightObj.(*object.IntegerInst).Value,
            }
    case *ast.ComparisonExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        switch node.Operator.Type {
        case token.GT:
            if leftObj.(*object.IntegerInst).Value > rightObj.(*object.IntegerInst).Value {
                return True
            } else {
                return False
            }
        case token.LT:
            if leftObj.(*object.IntegerInst).Value < rightObj.(*object.IntegerInst).Value {
                return True
            } else {
                return False
            }
        }
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return &object.IntegerInst{Value: val}
    case *ast.StringExpression:
        return &object.PyStrInst{Value: node.Value.Literals}
    case *ast.ListExpression:
        listObj := &object.ListInst{}
        for _, item := range node.Items {
            listObj.Items = append(listObj.Items, Eval(item, env))
        }
        return listObj
    case *ast.DictExpression:
        dictObj := &object.DictInst{
            Map: map[object.Object]object.Object{},
        }
        for i := 0; i < len(node.Keys); i++ {
            k, v := node.Keys[i], node.Vals[i]
            dictObj.Map[Eval(k, env)] = Eval(v, env)
        }
        return dictObj
    case *ast.CallExpression:
        return evalCallExpression(node, env)
    case *ast.AttributeExpression:
        inst := Eval(node.Expr, env)
        return inst.Py__getattribute__(&object.PyStrInst{node.Attr.Literals})
    case *ast.ExpressionStatement:
        return Eval(node.Value, env)
    }
    return object.Py_None
}

func execAssignStatement(stmt *ast.AssignStatement, env *Environment) {
    switch attr := stmt.Target.(type) {
    case *ast.AttributeExpression:
        instObj := Eval(attr.Expr, env)
        valObj := Eval(stmt.Value, env)
        instObj.Py__setattr__(&object.PyStrInst{attr.Attr.Literals}, valObj)
    case *ast.IdentifierExpression:
        env.Set(attr.Identifier.Literals, Eval(stmt.Value, env))
    }
}

func execIfStatement(stmt ast.Statement, env *Environment) {
    if stmt != nil {
        ifstmt := stmt.(*ast.IfStatement)
        if ifstmt.Condition == nil || Eval(ifstmt.Condition, env) == True {
            Exec(ifstmt.Body, env)
        } else {
            execIfStatement(ifstmt.Else, env)
        }
    }
}

func execWhileStatement(stmt *ast.WhileStatement, env *Environment) {
    for Eval(stmt.Condition, env) == True {
        Exec(stmt.Body, env)
    }
}

func execDefStatement(stmt *ast.DefStatement, env *Environment) {
    funcObj := &object.FunctionInst{
        Name: stmt.Name.Literals,
        Body: stmt.Body,
    }

    var params []string
    for _, tok := range stmt.Params {
        params = append(params, tok.Literals)
    }
    funcObj.Params = params
    env.Set(funcObj.Name, funcObj)
}

func execClassStatement(node *ast.ClassStatement, env *Environment) {
    clsEnv := env.deriveEnv()
    Exec(node.Body, clsEnv)

    clsObj := &object.PyClass{
        Name: node.Name.Literals,
        Dict: clsEnv.store,
    }

    if env.Get(node.Parent.Literals) != nil {
        // FIXME: there would be issue if inherit object
        clsObj.Base = env.Get(node.Parent.Literals).(object.Class)
    } else {
        clsObj.Base = object.PyObject
    }

    env.Set(clsObj.Name, clsObj)
}

func evalCallExpression(callNode *ast.CallExpression, parentEnv *Environment) object.Object {
    callObj := Eval(callNode.Name, parentEnv)

    switch obj := callObj.(type) {
    case *object.Print:
        var paramObjs []object.Object
        for _, param := range callNode.Params {
            paramObjs = append(paramObjs, Eval(param, parentEnv))
        }
        obj.Py__call__(paramObjs)
        return object.Py_None
    case *object.BoundMethod:
        args := []object.Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        args = append([]object.Object{obj.Inst}, args...)

        return evalFuncCallExpr(obj.Func, args, parentEnv.deriveEnv())
    case object.Function:
        args := []object.Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        return evalFuncCallExpr(obj, args, parentEnv.deriveEnv())
    case *object.BuiltinClass:
        if obj.Name == "super" {
            return &object.SuperInst{Py__self__: context.(*object.PyInst)}
        }
    case object.Class:
        args := []object.Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        return evalClassCallExpr(obj, args, parentEnv.deriveEnv())
    }

    return object.Py_None
}

func evalSuperCallExpr() {
}

func evalFuncCallExpr(f object.Function, args []object.Object, env *Environment) object.Object {
    switch obj := f.(type) {
    case object.BuiltinFunction:
        switch obj := f.(type) {
        case *object.BuiltinObjectNew:
            return obj.Call(args[0].(object.Class))
        }
    case *object.FunctionInst:
        for i, _ := range obj.Params {
            env.Set(obj.Params[i], args[i])
        }

        for _, stmt := range obj.Body {
            switch node := stmt.(type) {
            case *ast.ReturnStatement:
                return Eval(node.Value, env)
            default:
                Exec([]ast.Statement{stmt}, env)
            }
        }
    }


    return object.Py_None
}

func evalClassCallExpr(cls object.Class, args []object.Object, env *Environment) object.Object {
    // TODO: by now, super() in __new__ is invalid

    __new__ := cls.Py__getattribute__(&object.PyStrInst{"__new__"})

    var instObj *object.PyInst
    if __new__ != nil {
        instObj = evalFuncCallExpr(__new__.(object.Function), []object.Object{cls}, env).(*object.PyInst)
    } else {
        instObj = cls.Py__new__(cls)
    }

    __init__ := cls.Py__getattribute__(&object.PyStrInst{"__init__"})
    if __init__ == nil {
        return instObj
    }

    args = append([]object.Object{instObj}, args...)

    context = instObj
    evalFuncCallExpr(__init__.(object.Function), args, env)
    context = nil

    return instObj
}

