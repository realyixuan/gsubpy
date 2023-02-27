package evaluator

import (
    "strconv"
    "gsubpy/ast"
    "gsubpy/object"
    "gsubpy/token"
)

var env = map[string]object.Object{
    "False": &object.BoolObject{Value: 0},
    "True": &object.BoolObject{Value: 1},
}

func Exec(stmts []ast.Statement) {
    for _, stmt := range stmts {
        switch node := stmt.(type) {
        case *ast.AssignStatement:
            execAssignStatement(node)
        case *ast.IfStatement:
            execIfStatement(node)
        case *ast.WhileStatement:
            execWhileStatement(node)
        case *ast.ExpressionStatement:
            Eval(node)
        }
    }
}

func Eval(expression ast.Expression) object.Object {
    switch node := expression.(type) {
    case *ast.IdentifierExpression:
        return env[node.Identifier.Literals]
    case *ast.PlusExpression:
        leftObj := Eval(node.Left)
        rightObj := Eval(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value + rightObj.(*object.NumberObject).Value,
            }
    case *ast.MinusExpression:
        leftObj := Eval(node.Left)
        rightObj := Eval(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value - rightObj.(*object.NumberObject).Value,
            }
    case *ast.MulExpression:
        leftObj := Eval(node.Left)
        rightObj := Eval(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value * rightObj.(*object.NumberObject).Value,
            }
    case *ast.DivideExpression:
        leftObj := Eval(node.Left)
        rightObj := Eval(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value / rightObj.(*object.NumberObject).Value,
            }
    case *ast.ComparisonExpression:
        leftObj := Eval(node.Left)
        rightObj := Eval(node.Right)
        switch node.Operator.TokenType {
        case token.GT:
            if leftObj.(*object.NumberObject).Value > rightObj.(*object.NumberObject).Value {
                return env["True"]
            } else {
                return env["False"]
            }
        case token.LT:
            if leftObj.(*object.NumberObject).Value < rightObj.(*object.NumberObject).Value {
                return env["True"]
            } else {
                return env["False"]
            }
        }
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return &object.NumberObject{Value: val}
    case *ast.ExpressionStatement:
        return Eval(node.Value)
    }
    return nil    // XXX: temporary solution
}

func execAssignStatement(stmt *ast.AssignStatement) {
    env[stmt.Identifier.Literals] = Eval(stmt.Value)
}

func execIfStatement(stmt *ast.IfStatement) {
    if Eval(stmt.Condition) == env["True"] {
        Exec(stmt.Body)
    }
}

func execWhileStatement(stmt *ast.WhileStatement) {
    for Eval(stmt.Condition) == env["True"] {
        Exec(stmt.Body)
    }
}

