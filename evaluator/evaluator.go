package evaluator

import (
    "strconv"
    "gsubpy/ast"
)

var E = map[string]int{}

func eval(stmts []ast.AssignStatement) {
    for _, stmt := range stmts {
        evalAssignStatement(stmt)
    }
}

func evalAssignStatement(stmt ast.AssignStatement) {
    val, _ := strconv.Atoi(stmt.Value.Literals)
    E[stmt.Identifier.Literals] = val
}

