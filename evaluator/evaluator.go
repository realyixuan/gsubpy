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

func exec(stmts []ast.Statement) {
    for _, stmt := range stmts {
        switch node := stmt.(type) {
        case *ast.AssignStatement:
            execAssignStatement(node)
        case *ast.IfStatement:
            execIfStatement(node)
        }
    }
}

func evalExpression(expression ast.Expression) object.Object {
    switch node := expression.(type) {
    case *ast.IdentifierExpression:
        return env[node.Identifier.Literals]
    case *ast.PlusExpression:
        leftObj := evalExpression(node.Left)
        rightObj := evalExpression(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value + rightObj.(*object.NumberObject).Value,
            }
    case *ast.MinusExpression:
        leftObj := evalExpression(node.Left)
        rightObj := evalExpression(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value - rightObj.(*object.NumberObject).Value,
            }
    case *ast.MulExpression:
        leftObj := evalExpression(node.Left)
        rightObj := evalExpression(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value * rightObj.(*object.NumberObject).Value,
            }
    case *ast.DivideExpression:
        leftObj := evalExpression(node.Left)
        rightObj := evalExpression(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value / rightObj.(*object.NumberObject).Value,
            }
    case *ast.ComparisonExpression:
        leftObj := evalExpression(node.Left)
        rightObj := evalExpression(node.Right)
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
    }
    return nil    // XXX: temporary solution
}

func execAssignStatement(stmt *ast.AssignStatement) {
    env[stmt.Identifier.Literals] = evalExpression(stmt.Value)
}

func execIfStatement(stmt *ast.IfStatement) {
    if evalExpression(stmt.Condition) == env["True"] {
        exec(stmt.Body)
    }
}
