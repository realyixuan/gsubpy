package evaluator

import (
    "strconv"

    "github.com/realyixuan/gsubpy/ast"
    "github.com/realyixuan/gsubpy/token"
)

func Exec(stmts []ast.Statement, env *Environment) (Object, bool) {
    var s ast.Statement
    defer func() {
        if r := recover(); r != nil {
            f := Frame{Literals: s.GetLiterals()}
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
        return env.Get(newStringInst(node.Identifier.Literals))
    case *ast.PlusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if leftObj.Type() != rightObj.Type() {
            panic(Error("TypeError: two different types"))
        }

        switch leftObj.(type) {
        case *IntegerInst:
            return newIntegerInst(leftObj.(*IntegerInst).Value + rightObj.(*IntegerInst).Value)
        case *StringInst:
            return newStringInst(leftObj.(*StringInst).Value + rightObj.(*StringInst).Value)
        }

    case *ast.MinusExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return newIntegerInst(leftObj.(*IntegerInst).Value - rightObj.(*IntegerInst).Value)
    case *ast.MulExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        return newIntegerInst(leftObj.(*IntegerInst).Value * rightObj.(*IntegerInst).Value)
    case *ast.DivideExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if rightObj.(*IntegerInst).Value == 0 {
            panic(Error("ZeroDivisionError: division by zero"))
        }

        return newIntegerInst(leftObj.(*IntegerInst).Value / rightObj.(*IntegerInst).Value)
    case *ast.ComparisonExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        switch node.Operator.Type {
        case token.GT:
            return typeCall(__gt__, leftObj, rightObj)
        case token.LT:
            return typeCall(__lt__, leftObj, rightObj)
        case token.EQ:
            return typeCall(__eq__, leftObj, rightObj)
        }
        return Py_True
    case *ast.NotExpression:
        obj := Eval(node.Expr, env)
        if toPy_True(obj) {
            return Py_False
        } else {
            return Py_True
        }
    case *ast.AndExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if !toPy_True(leftObj) {
            return leftObj
        }

        return rightObj
    case *ast.OrExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if toPy_True(leftObj) {
            return leftObj
        }

        return rightObj
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return newIntegerInst(int64(val))
    case *ast.StringExpression:
        return newStringInst(node.Value.Literals)
    case *ast.ListExpression:
        listObj := newListInst()
        for _, item := range node.Items {
            listObj.items = append(listObj.items, Eval(item, env))
        }
        return listObj
    case *ast.DictExpression:
        dictObj := newDictInst()
        for i := 0; i < len(node.Keys); i++ {
            k, v := node.Keys[i], node.Vals[i]
            dictObj.Set(Eval(k, env), Eval(v, env))
        }
        return dictObj
    case *ast.CallExpression:
        return evalCallExpression(node, env)
    case *ast.AttributeExpression:
        inst := Eval(node.Expr, env)
        return inst.Attr(newStringInst(node.Attr.Literals))
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

        typeCall(__setattr__, instObj, newStringInst(attr.Attr.Literals), valObj)
    case *ast.IdentifierExpression:
        env.SetFromString(attr.Identifier.Literals, Eval(stmt.Value, env))
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

    var params []*StringInst
    for _, tok := range stmt.Params {
        params = append(params, newStringInst(tok.Literals))
    }

    funcObj := newFunctionInst(
        newStringInst(stmt.Name.Literals),
        params,
        stmt.Body,
        env,
    )

    env.Set(funcObj.Name, funcObj)
}

func execClassStatement(node *ast.ClassStatement, env *Environment) {
    clsEnv := env.DeriveEnv()
    Exec(node.Body, clsEnv)

    var base Class
    if env.GetFromString(node.Parent.Literals) != nil {
        base = env.GetFromString(node.Parent.Literals).(Class)
    } else {
        base = Py_object
    }

    clsObj := newPyclass(
        newStringInst(node.Name.Literals),
        base,
        clsEnv.Store(),
    )

    env.Set(clsObj.name, clsObj)
}

func evalCallExpression(callNode *ast.CallExpression, parentEnv *Environment) Object {
    callObj := Eval(callNode.Name, parentEnv)

    var args []Object
    for _, param := range callNode.Params {
        args = append(args, Eval(param, parentEnv))
    }
    
    return Call(callObj, args...)
}

func toPy_True(obj Object) bool {
    if rv := Call(Py_bool, obj); rv == Py_True {
        return true
    } else {
        return false
    }
}

func StringOf(obj Object) Object {
    __str__Fn := attrItself(obj.Type(), __str__)
    return Call(__str__Fn, obj)
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

