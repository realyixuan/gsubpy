package evaluator

import (
    "strconv"
    "gsubpy/ast"
    "gsubpy/token"
)

// simplest way for now
var context Object

func Exec(stmts []ast.Statement, env *Environment) *NoneInst {
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

    return Py_None
}

func Eval(expression ast.Expression, env *Environment) Object {
    switch node := expression.(type) {
    case *ast.IdentifierExpression:
        return env.Get(node.Identifier.Literals)
    case *ast.PlusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if leftObj.Type() != rightObj.Type() {
            panic(&ExceptionInst{Msg: "TypeError: two different types"})
        }

        switch leftObj.(type) {
        case *IntegerInst:
            return &IntegerInst{
                Value: leftObj.(*IntegerInst).Value + rightObj.(*IntegerInst).Value,
                }
        case *PyStrInst:
            return &PyStrInst{
                Value: leftObj.(*PyStrInst).Value + rightObj.(*PyStrInst).Value,
                }
        }

    case *ast.MinusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return &IntegerInst{
            Value: leftObj.(*IntegerInst).Value - rightObj.(*IntegerInst).Value,
            }
    case *ast.MulExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return &IntegerInst{
            Value: leftObj.(*IntegerInst).Value * rightObj.(*IntegerInst).Value,
            }
    case *ast.DivideExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if rightObj.(*IntegerInst).Value == 0 {
            panic(&ExceptionInst{Msg: "ZeroDivisionError: division by zero"})
        }

        return &IntegerInst{
            Value: leftObj.(*IntegerInst).Value / rightObj.(*IntegerInst).Value,
            }
    case *ast.ComparisonExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        switch node.Operator.Type {
        case token.GT:
            if leftObj.(*IntegerInst).Value > rightObj.(*IntegerInst).Value {
                return Py_True
            } else {
                return Py_False
            }
        case token.LT:
            if leftObj.(*IntegerInst).Value < rightObj.(*IntegerInst).Value {
                return Py_True
            } else {
                return Py_False
            }
        }
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return &IntegerInst{Value: val}
    case *ast.StringExpression:
        return &PyStrInst{Value: node.Value.Literals}
    case *ast.ListExpression:
        listObj := &ListInst{}
        for _, item := range node.Items {
            listObj.Items = append(listObj.Items, Eval(item, env))
        }
        return listObj
    case *ast.DictExpression:
        dictObj := &DictInst{
            Map: map[PyStrInst]Object{},
        }
        for i := 0; i < len(node.Keys); i++ {
            k, v := node.Keys[i], node.Vals[i]
            dictObj.Py__setitem__(Eval(k, env), Eval(v, env))
        }
        return dictObj
    case *ast.CallExpression:
        return evalCallExpression(node, env)
    case *ast.AttributeExpression:
        inst := Eval(node.Expr, env)
        return inst.Py__getattribute__(&PyStrInst{node.Attr.Literals})
    case *ast.ExpressionStatement:
        return Eval(node.Value, env)
    }
    return Py_None
}

func execAssignStatement(stmt *ast.AssignStatement, env *Environment) {
    switch attr := stmt.Target.(type) {
    case *ast.AttributeExpression:
        instObj := Eval(attr.Expr, env)
        valObj := Eval(stmt.Value, env)
        instObj.Py__setattr__(&PyStrInst{attr.Attr.Literals}, valObj)
    case *ast.IdentifierExpression:
        env.Set(attr.Identifier.Literals, Eval(stmt.Value, env))
    }
}

func execIfStatement(stmt ast.Statement, env *Environment) {
    if stmt != nil {
        ifstmt := stmt.(*ast.IfStatement)
        if ifstmt.Condition == nil || Eval(ifstmt.Condition, env) == Py_True {
            Exec(ifstmt.Body, env)
        } else {
            execIfStatement(ifstmt.Else, env)
        }
    }
}

func execWhileStatement(stmt *ast.WhileStatement, env *Environment) {
    for Eval(stmt.Condition, env) == Py_True {
        Exec(stmt.Body, env)
    }
}

func execDefStatement(stmt *ast.DefStatement, env *Environment) {
    funcObj := &FunctionInst{
        Name: stmt.Name.Literals,
        Body: stmt.Body,
        env: env,
    }

    var params []string
    for _, tok := range stmt.Params {
        params = append(params, tok.Literals)
    }
    funcObj.Params = params
    env.Set(funcObj.Name, funcObj)
}

func execClassStatement(node *ast.ClassStatement, env *Environment) {
    clsEnv := env.DeriveEnv()
    Exec(node.Body, clsEnv)

    clsObj := &PyClass{
        Name: node.Name.Literals,
        Dict: clsEnv.Store(),
    }

    if env.Get(node.Parent.Literals) != nil {
        // FIXME: there would be issue if inherit object
        clsObj.Base = env.Get(node.Parent.Literals).(Class)
    } else {
        clsObj.Base = Py_object
    }

    env.Set(clsObj.Name, clsObj)
}

func evalCallExpression(callNode *ast.CallExpression, parentEnv *Environment) Object {
    callObj := Eval(callNode.Name, parentEnv)

    switch obj := callObj.(type) {
    case BuiltinFunction:
        var paramObjs []Object
        for _, param := range callNode.Params {
            paramObjs = append(paramObjs, Eval(param, parentEnv))
        }

        switch o := obj.(type) {
        case *Print:
            o.Py__call__(paramObjs...)
            return Py_None
        case *Len:
            return o.Py__call__(paramObjs...)
        case *PyNew:
            return o.Py__call__(paramObjs...)
        case *Super:
            return &SuperInst{Py__self__: context.(*PyInst)}
        }
    case Class:
        args := []Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        return evalClassCallExpr(obj, args, parentEnv.DeriveEnv())
    case *BoundMethod:
        args := []Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        args = append([]Object{obj.Inst}, args...)

        return obj.Func.Py__call__(args...)
    case Function:
        args := []Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }

        return obj.Py__call__(args...)
    }

    return Py_None
}

func evalClassCallExpr(cls Class, args []Object, env *Environment) Object {
    // A special branch. remove it when __call__ added into class
    switch o := cls.(type) {
    case *Pytype:
        return o.Py__call__(args...)
    }

    // TODO: by now, super() in __new__ is invalid

    __new__ := cls.Py__getattribute__(&PyStrInst{"__new__"})

    var instObj Object
    if __new__ != nil {
        instObj = __new__.(Function).Py__call__(cls)
    } else {
        instObj = cls.Py__new__(cls)
    }

    __init__ := cls.Py__getattribute__(&PyStrInst{"__init__"})
    if __init__ == nil {
        return instObj
    }

    args = append([]Object{instObj}, args...)

    context = instObj
    __init__.(Function).Py__call__(args...)
    context = nil

    return instObj
}

