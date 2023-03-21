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
    True = &object.BoolObject{Value: 1}
    False = &object.BoolObject{Value: 0}
    None = &object.NoneObject{Value: 0}
)

var __builtins__ = map[string]object.Object{
    "object": object.Pyobject,
    "True": True,
    "False": False,
    "None": None,
    "print": &object.Print{},
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

func Exec(stmts []ast.Statement, env *Environment) *object.NoneObject {
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

    return None
}

func Eval(expression ast.Expression, env *Environment) object.Object {
    switch node := expression.(type) {
    case *ast.IdentifierExpression:
        return env.Get(node.Identifier.Literals)
    case *ast.PlusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if leftObj.GetObjType() != rightObj.GetObjType() {
            panic(&object.ExceptionObject{Msg: "TypeError: two different types"})
        }

        switch leftObj.(type) {
        case *object.NumberObject:
            return &object.NumberObject{
                Value: leftObj.(*object.NumberObject).Value + rightObj.(*object.NumberObject).Value,
                }
        case *object.StringObject:
            return &object.StringObject{
                Value: leftObj.(*object.StringObject).Value + rightObj.(*object.StringObject).Value,
                }
        }

    case *ast.MinusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value - rightObj.(*object.NumberObject).Value,
            }
    case *ast.MulExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value * rightObj.(*object.NumberObject).Value,
            }
    case *ast.DivideExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if rightObj.(*object.NumberObject).Value == 0 {
            panic(&object.ExceptionObject{Msg: "ZeroDivisionError: division by zero"})
        }

        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value / rightObj.(*object.NumberObject).Value,
            }
    case *ast.ComparisonExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        switch node.Operator.Type {
        case token.GT:
            if leftObj.(*object.NumberObject).Value > rightObj.(*object.NumberObject).Value {
                return True
            } else {
                return False
            }
        case token.LT:
            if leftObj.(*object.NumberObject).Value < rightObj.(*object.NumberObject).Value {
                return True
            } else {
                return False
            }
        }
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return &object.NumberObject{Value: val}
    case *ast.StringExpression:
        return &object.StringObject{Value: node.Value.Literals}
    case *ast.ListExpression:
        listObj := &object.ListObject{}
        for _, item := range node.Items {
            listObj.Items = append(listObj.Items, Eval(item, env))
        }
        return listObj
    case *ast.DictExpression:
        dictObj := &object.DictObject{
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
        return inst.Py__getattribute__(node.Attr.Literals)
    case *ast.ExpressionStatement:
        return Eval(node.Value, env)
    }
    return None
}

func execAssignStatement(stmt *ast.AssignStatement, env *Environment) {
    switch attr := stmt.Target.(type) {
    case *ast.AttributeExpression:
        instObj := Eval(attr.Expr, env)
        valObj := Eval(stmt.Value, env)
        instObj.Py__setattr__(attr.Attr.Literals, valObj)
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
    funcObj := &object.FunctionObject{
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
        Py__dict__: clsEnv.store,
    }

    if env.Get(node.Parent.Literals) != nil {
        // FIXME: there would be issue if inherit object
        clsObj.Py__base__ = env.Get(node.Parent.Literals).(*object.PyClass)
        
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
        obj.Call(paramObjs)
        return None
    case *object.PyClass:
        args := []object.Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        return evalClassCallExpr(obj, args, parentEnv.deriveEnv())
    case *object.BoundMethod:
        args := []object.Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        args = append([]object.Object{obj.Inst}, args...)

        return evalFuncCallExpr(obj.Func, args, parentEnv.deriveEnv())
    case *object.FunctionObject:
        args := []object.Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        return evalFuncCallExpr(obj, args, parentEnv.deriveEnv())
    case *object.BuiltinClass:
        if obj.Name == "super" {
            return &object.SuperInstance{Py__self__: context.(*object.PyInstance)}
        }
    }

    return None
}

func evalSuperCallExpr() {
}

func evalFuncCallExpr(funcObj *object.FunctionObject, args []object.Object, env *Environment) object.Object {
    for i, _ := range funcObj.Params {
        env.Set(funcObj.Params[i], args[i])
    }

    for _, stmt := range funcObj.Body {
        switch node := stmt.(type) {
        case *ast.ReturnStatement:
            return Eval(node.Value, env)
        default:
            Exec([]ast.Statement{stmt}, env)
        }
    }

    return None
}

func evalClassCallExpr(clsObj *object.PyClass, args []object.Object, env *Environment) object.Object {
    instObj := clsObj.Py__new__(clsObj)

    __init__ := clsObj.Py__getattribute__("__init__")
    if __init__ == nil {
        return instObj
    }

    args = append([]object.Object{instObj}, args...)

    context = instObj
    evalFuncCallExpr(__init__.(*object.FunctionObject), args, env)
    context = nil

    return instObj
}

