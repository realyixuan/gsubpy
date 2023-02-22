package parser

import (
    "gsubpy/ast"
    "gsubpy/lexer"
)

func Parsing(l *lexer.Lexer) []ast.AssignStatement {
    // return statements
    var stmts []ast.AssignStatement

    stmt := parsingStatement(l)

    stmts = append(stmts, stmt)

    return stmts

}

func parsingStatement(l *lexer.Lexer) ast.AssignStatement {
    assignment := ast.AssignStatement{Identifier: l.NextToken()}

    l.NextToken()

    assignment.Value = l.NextToken()

    return assignment
}

