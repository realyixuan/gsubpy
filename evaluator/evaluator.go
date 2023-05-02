package evaluator

import (
    "fmt"
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
        case *ast.ForStatement:
            rv, isReturn := execForStatement(node, env)
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
        case *ast.AssertStatement:
            execAssertStatement(node, env)
        }
    }
    return Py_None, false
}

func Eval(expression ast.Expression, env *Environment) Object {
    switch node := expression.(type) {
    case *ast.IdentifierExpression:
        return env.Get(newStringInst(node.Identifier.Literals))
    case *ast.PlusExpression:
        left := Eval(node.Left, env)
        right := Eval(node.Right, env)
        return op_ADD(left, right)
    case *ast.MinusExpression:
        left := Eval(node.Left, env)
        right := Eval(node.Right, env)
        return op_SUB(left, right)
    case *ast.MulExpression:
        left := Eval(node.Left, env)
        right := Eval(node.Right, env)
        return op_MUL(left, right)
    case *ast.DivideExpression:
        left := Eval(node.Left, env)
        right := Eval(node.Right, env)
        return op_DIV(left, right)
    case *ast.ComparisonExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)
        switch node.Operator.Type {
        case token.GT:
            return op_GT(leftObj, rightObj)
        case token.LT:
            return op_LT(leftObj, rightObj)
        case token.EQ:
            return op_EQ(leftObj, rightObj)
        case token.IN:
            return op_IN(leftObj, rightObj)
        case token.NIN:
            return op_NIN(leftObj, rightObj)
        case token.IS:
            return op_IS(leftObj, rightObj)
        case token.ISN:
            return op_ISN(leftObj, rightObj)
        }
        return Py_True
    case *ast.NotExpression:
        obj := Eval(node.Expr, env)
        if op_CALL(Py_bool, obj) == Py_True {
            return Py_False
        } else {
            return Py_True
        }
    case *ast.AndExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if op_CALL(Py_bool, leftObj) == Py_False {
            return leftObj
        }

        return rightObj
    case *ast.OrExpression:
        leftObj := Eval(node.Left, env)
        rightObj := Eval(node.Right, env)

        if op_CALL(Py_bool, leftObj) == Py_True {
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
        dict := newDictInst()
        for i := 0; i < len(node.Keys); i++ {
            k, v := Eval(node.Keys[i], env), Eval(node.Vals[i], env)
            op_SUBSCR_SET(dict, k, v)
        }
        return dict
    case *ast.SubscriptExpression:
        return op_SUBSCR_GET(Eval(node.Target, env), Eval(node.Val, env))
    case *ast.CallExpression:
        return evalCallExpression(node, env)
    case *ast.AttributeExpression:
        inst := Eval(node.Expr, env)
        return op_GETATTR(inst, newStringInst(node.Attr.Literals))
    case *ast.ExpressionStatement:
        return Eval(node.Value, env)
    }
    return Py_None
}

func execAssignStatement(stmt *ast.AssignStatement, env *Environment) {
    switch attr := stmt.Target.(type) {
    case *ast.AttributeExpression:
        inst := Eval(attr.Expr, env)
        val := Eval(stmt.Value, env)

        op_SETATTR(inst, newStringInst(attr.Attr.Literals), val)
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

func execForStatement(stmt *ast.ForStatement, env *Environment) (Object, bool) {
    target := Eval(stmt.Target, env)
    iterator := op_CALL(Py_iter, target)

    for val := op_CALL(Py_next, iterator); val != nil; {
        // not considering multi-values currently,
        env.SetFromString(stmt.Identifiers[0].Literals, val)
        rv, isReturn := Exec(stmt.Body, env)
        if isReturn {
            return rv, isReturn
        }
        val = op_CALL(Py_next, iterator)
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

    clsObj := op_CALL(
        Py_type,
        newStringInst(node.Name.Literals),
        base,
        clsEnv.Store(),
    )

    env.SetFromString(node.Name.Literals, clsObj)
}

func execAssertStatement(stmt *ast.AssertStatement, env *Environment) {
    for op_CALL(Py_bool, Eval(stmt.Condition, env)) == Py_False {
        panic(Error(fmt.Sprintf("assert error:  %v", StringOf(Eval(stmt.Msg, env)))))
    }
}

func evalCallExpression(callNode *ast.CallExpression, parentEnv *Environment) Object {
    callObj := Eval(callNode.Name, parentEnv)

    var args []Object
    for _, param := range callNode.Params {
        args = append(args, Eval(param, parentEnv))
    }
    
    return op_CALL(callObj, args...)
}

func op_ADD(left Object, right Object) Object {
    return typeCall(__add__, left, right)
}

func op_SUB(left Object, right Object) Object {
    return typeCall(__sub__, left, right)
}

func op_MUL(left Object, right Object) Object {
    return typeCall(__mul__, left, right)
}

func op_DIV(left Object, right Object) Object {
    return typeCall(__floordiv__, left, right)
}

func op_IN(left Object, right Object) Object {
    return typeCall(__contains__, right, left)
}

func op_NIN(left Object, right Object) Object {
    if op_IN(left, right) == Py_True {
        return Py_False
    } else {
        return Py_True
    }
}

func op_IS(left Object, right Object) Object {
    if left.id() == right.id() {
        return Py_True
    } else {
        return Py_False
    }
}

func op_ISN(left Object, right Object) Object {
    if op_IS(left, right) == Py_True {
        return Py_False
    } else {
        return Py_True
    }
}

func op_EQ(left Object, right Object) Object {
    return typeCall(__eq__, left, right)
}

func op_GT(left Object, right Object) Object {
    return typeCall(__gt__, left, right)
}

func op_LT(left Object, right Object) Object {
    return typeCall(__lt__, left, right)
}

func op_SETATTR(inst Object, attr *StringInst, value Object) {
    typeCall(__setattr__, inst, attr, value)
}

func op_GETATTR(inst Object, attr *StringInst) Object {
    return typeCall(__getattribute__, inst, attr)
}

func op_SUBSCR_GET(inst Object, item Object) Object {
    return typeCall(__getitem__, inst, item)
}

func op_SUBSCR_SET(inst Object, key Object, item Object) Object {
    return typeCall(__setitem__, inst, key, item)
}

func op_CALL(obj Object, args ...Object) Object {
    // eliminate obj.otype() here ?
    __call__Fn := attrItself(obj.otype(), __call__)
    args = append([]Object{obj}, args...)

    if __call__Fn != PyBuiltinFunction__call__ {
        return op_CALL(__call__Fn, args...)
    } else {
        return __call__Fn.(Function).call(args...)
    }
}

func typeCall(attrName *StringInst, obj Object, args ...Object) Object {
    attr := attrItself(obj.otype(), attrName)
    if attr == nil {
        panic(Error(fmt.Sprintf("%v object is not callable", StringOf(obj.otype()))))
    }

    fn, ok := attr.(Function) 
    if !ok {
        panic(Error(fmt.Sprintf("%v object is not callable", StringOf(obj.otype()))))
    }
    
    args = append([]Object{obj}, args...)
    return op_CALL(fn, args...)
}

func StringOf(obj Object) Object {
    __str__Fn := attrItself(obj.otype(), __str__)
    return op_CALL(__str__Fn, obj)
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
