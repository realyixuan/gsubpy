package evaluator

import (
    "strconv"
    "gsubpy/ast"
    "gsubpy/object"
)

var E = map[string]object.Object{}

func eval(stmts []ast.Statement) {
    for _, stmt := range stmts {
        switch node := stmt.(type) {
        case *ast.AssignStatement:
            evalAssignStatement(node)
        }
    }
}

func evalAssignStatement(stmt *ast.AssignStatement) {
    E[stmt.Identifier.Literals] = evalExpression(stmt.Value)
}

func evalExpression(expression ast.Expression) object.Object {
    switch node := expression.(type) {
    case *ast.PlusExpression:
        leftObj := evalExpression(node.Left)
        rightObj := evalExpression(node.Right)
        return &object.NumberObject{
            Value: leftObj.(*object.NumberObject).Value + rightObj.(*object.NumberObject).Value,
            }
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return &object.NumberObject{Value: val}
    }
    return nil    // XXX: temporary solution
}

