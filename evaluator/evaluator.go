package evaluator

import (
    "strconv"
    "gsubpy/ast"
    "gsubpy/object"
    "gsubpy/token"
)

var (
    True = &object.BoolObject{Value: 1}
    False = &object.BoolObject{Value: 0}
    None = &object.NoneObject{Value: 0}
)

var __builtins__ = map[string]object.Object{
    "True": True,
    "False": False,
    "None": None,
    "print": &object.Print{},
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
            panic(&object.ExceptionObject{"TypeError: two different types"})
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
            panic(&object.ExceptionObject{"ZeroDivisionError: division by zero"})
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
    case *ast.FunctionCallExpression:
        return evalFunctionCallExpression(node, env)
    case *ast.ExpressionStatement:
        return Eval(node.Value, env)
    }
    return None
}

func execAssignStatement(stmt *ast.AssignStatement, env *Environment) {
    env.Set(stmt.Identifier.Literals, Eval(stmt.Value, env))
}

func execIfStatement(stmt *ast.IfStatement, env *Environment) {
    if stmt != nil {
        if stmt.Condition == nil || Eval(stmt.Condition, env) == True {
            Exec(stmt.Body, env)
        } else {
            execIfStatement(stmt.Else, env)
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

func evalFunctionCallExpression(funcNode *ast.FunctionCallExpression, parentEnv *Environment) object.Object {
    funcObj := Eval(funcNode.Name, parentEnv)

    switch obj := funcObj.(type) {
    case *object.Print:
        var paramObjs []object.Object
        for _, param := range funcNode.Params {
            paramObjs = append(paramObjs, Eval(param, parentEnv))
        }
        obj.Call(paramObjs)
        return None
    default:
        env := parentEnv.deriveEnv()

        for i, expr := range funcNode.Params {
            env.Set(obj.(*object.FunctionObject).Params[i], Eval(expr, parentEnv))
        }

        for _, stmt := range obj.(*object.FunctionObject).Body {
            switch node := stmt.(type) {
            case *ast.ReturnStatement:
                return Eval(node.Value, env)
            default:
                Exec([]ast.Statement{stmt}, env)
            }
        }
    }

    return None
}

