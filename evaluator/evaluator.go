package evaluator

import (
    "strconv"
    "gsubpy/ast"
    "gsubpy/token"
)

// simplest way for now
var context Object

func Exec(stmts []ast.Statement, env *Environment) (Object, bool) {
    var s ast.Statement
    defer func() {
        if r := recover(); r != nil {
            // TODO: add method of getting literals
            var f Frame
            switch o := s.(type) {
            case *ast.AssignStatement:
                f = Frame{Literals: o.Literals}
            case *ast.ExpressionStatement:
                f = Frame{Literals: o.Literals}
            }
            Py_traceback.append(f)
            panic(r)
        }
    }()

    for _, stmt := range stmts {
        s = stmt
        switch node := stmt.(type) {
        case *ast.AssignStatement:
            execAssignStatement(node, env)
        case *ast.IfStatement:
            rv, isReturn := execIfStatement(node, env)
            if isReturn {
                return rv, isReturn
            }
        case *ast.WhileStatement:
            rv, isReturn := execWhileStatement(node, env)
            if isReturn {
                return rv, isReturn
            }
        case *ast.DefStatement:
            execDefStatement(node, env)
        case *ast.ClassStatement:
            execClassStatement(node, env)
        case *ast.ExpressionStatement:
            Eval(node, env)
        case *ast.ReturnStatement:
            return Eval(node.Value, env), true
        }
    }
    return Py_None, false
}

func Eval(expression ast.Expression, env *Environment) Object {
    switch node := expression.(type) {
    case *ast.IdentifierExpression:
        return env.Get(node.Identifier.Literals)
    case *ast.PlusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if leftObj.Py__class__() != rightObj.Py__class__() {
            panic(NewException("TypeError: two different types"))
        }

        switch leftObj.(type) {
        case *IntegerInst:
            return NewInteger(leftObj.(*IntegerInst).Value + rightObj.(*IntegerInst).Value)
        case *PyStrInst:
            return NewStrInst(leftObj.(*PyStrInst).Value + rightObj.(*PyStrInst).Value)
        }

    case *ast.MinusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return NewInteger(leftObj.(*IntegerInst).Value - rightObj.(*IntegerInst).Value)
    case *ast.MulExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return NewInteger(leftObj.(*IntegerInst).Value * rightObj.(*IntegerInst).Value)
    case *ast.DivideExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if rightObj.(*IntegerInst).Value == 0 {
            panic(NewException("ZeroDivisionError: division by zero"))
        }

        return NewInteger(leftObj.(*IntegerInst).Value / rightObj.(*IntegerInst).Value)
    case *ast.ComparisonExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        switch node.Operator.Type {
        case token.GT:
            return leftObj.Py__gt__(rightObj)
        case token.LT:
            return leftObj.Py__lt__(rightObj)
        case token.EQ:
            return leftObj.Py__eq__(rightObj)
        }
    case *ast.NotExpression:
        obj := Eval(node.Expr, env)
        // TODO: need to add __bool__ for every type
        // now, temporarily apply this to comparison expression
        if obj == Py_True {
            return Py_False
        } else {
            return Py_True
        }
    case *ast.AndExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        if leftObj == Py_True && rightObj == Py_True {
            return Py_True
        } else {
            return Py_False
        }
    case *ast.OrExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        if leftObj == Py_True || rightObj == Py_True {
            return Py_True
        } else {
            return Py_False
        }
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return NewInteger(val)
    case *ast.StringExpression:
        return NewStrInst(node.Value.Literals)
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
        return inst.Py__getattribute__(NewStrInst(node.Attr.Literals))
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
        instObj.Py__setattr__(NewStrInst(attr.Attr.Literals), valObj)
    case *ast.IdentifierExpression:
        env.Set(attr.Identifier.Literals, Eval(stmt.Value, env))
    }
}

func execIfStatement(stmt ast.Statement, env *Environment) (Object, bool) {
    if stmt != nil {
        ifstmt := stmt.(*ast.IfStatement)
        if ifstmt.Condition == nil || Eval(ifstmt.Condition, env) == Py_True {
            rv, isReturn := Exec(ifstmt.Body, env)
            if isReturn {
                return rv, true
            }
        } else {
            rv, isReturn := execIfStatement(ifstmt.Else, env)
            if isReturn {
                return rv, true
            }
        }
    }

    return nil, false
}

func execWhileStatement(stmt *ast.WhileStatement, env *Environment) (Object, bool) {
    for Eval(stmt.Condition, env) == Py_True {
        rv, isReturn := Exec(stmt.Body, env)
        if isReturn {
            return rv, isReturn
        }
    }
    return nil, false
}

func execDefStatement(stmt *ast.DefStatement, env *Environment) {
    funcObj := &FunctionInst{
        Name: NewStrInst(stmt.Name.Literals),
        Body: stmt.Body,
        env: env,
    }

    var params []string
    for _, tok := range stmt.Params {
        params = append(params, tok.Literals)
    }
    funcObj.Params = params
    env.Set(funcObj.Name.Value, funcObj)
}

func execClassStatement(node *ast.ClassStatement, env *Environment) {
    clsEnv := env.DeriveEnv()
    Exec(node.Body, clsEnv)

    clsObj := &PyClass{
        Name: NewStrInst(node.Name.Literals),
        Dict: &DictInst{Map: map[PyStrInst]Object{}},
    }

    for k, v := range clsEnv.Store() {
        clsObj.Dict.Py__setitem__(&k, v)
    }

    if env.Get(node.Parent.Literals) != nil {
        // FIXME: there would be issue if inherit object
        clsObj.Base = env.Get(node.Parent.Literals).(Class)
    } else {
        clsObj.Base = Py_object
    }

    env.Set(clsObj.Name.Value, clsObj)
}

func evalCallExpression(callNode *ast.CallExpression, parentEnv *Environment) Object {
    callObj := Eval(callNode.Name, parentEnv)

    var paramObjs []Object
    for _, param := range callNode.Params {
        paramObjs = append(paramObjs, Eval(param, parentEnv))
    }

    switch obj := callObj.(type) {
    // refactor it after add builtin-type interface
    case *Pytype:
        return obj.Call(paramObjs...)
    case BuiltinFunction:
        switch o := obj.(type) {
        case *Print:
            o.Py__call__(paramObjs...)
            return Py_None
        case *Len:
            return o.Py__call__(paramObjs...)
        case *PyNew:
            return o.Py__call__(paramObjs...)
        }
    case Class:
        args := []Object{}
        for _, param := range callNode.Params {
            args = append(args, Eval(param, parentEnv))
        }
        inst := Py_type.Py__call__(obj, args...)
        return inst
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

type Frame struct {
    ast.Literals
    context     string
}

type Traceback struct {
    Frames      []Frame
}

func (tb *Traceback) append(f Frame) {
    tb.Frames = append(tb.Frames, f)
}

var Py_traceback = &Traceback{}

