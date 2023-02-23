package evaluator

import (
    "strconv"
    "gsubpy/ast"
)

var E = map[string]int{}

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

func evalExpression(expression ast.Expression) int{
    switch node := expression.(type) {
    case *ast.PlusExpression:
        left := evalExpression(node.Left)
        right := evalExpression(node.Right)
        return left + right
    case *ast.NumberExpression:
        val, _ := strconv.Atoi(node.Value.Literals)
        return val
    }
    return 0    // XXX: temporary solution
}

