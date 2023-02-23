package parser

import (
    "gsubpy/ast"
    "gsubpy/lexer"
    "gsubpy/token"
)

type Parser struct {
    l       *lexer.Lexer
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}
    return p
}

func (p *Parser)Parsing() []ast.Statement {
    // return statements
    var stmts []ast.Statement

    stmt := p.parsingStatement()

    stmts = append(stmts, stmt)

    return stmts

}

func (p *Parser)parsingStatement() ast.Statement {
    switch p.l.CurToken.TokenType {
    case token.IDENTIFIER:
        if p.l.PeekNextToken().TokenType == token.ASSIGN {
            return p.parsingAssignStatement()
        }
    }
    return nil
}

func (p *Parser)parsingAssignStatement() *ast.AssignStatement {
    assignment := ast.AssignStatement{Identifier: p.l.CurToken}
    p.l.ReadNextToken()
    p.l.ReadNextToken()
    assignment.Value = p.parsingExpression()
    p.l.ReadNextToken()
    return &assignment
}

func (p *Parser)parsingExpression() ast.Expression {
    if p.l.PeekNextToken().TokenType == token.PLUS {
        left := p.l.CurToken
        p.l.ReadNextToken()
        p.l.ReadNextToken()
        right := p.l.CurToken
        rtExpression := &ast.PlusExpression{
            Left: &ast.NumberExpression{left},
            Right: &ast.NumberExpression{right},
        }
        p.l.ReadNextToken()
        return rtExpression
    } else if p.l.PeekNextToken().TokenType == token.EOF {
        rtExpression := &ast.NumberExpression{p.l.CurToken}
        p.l.ReadNextToken()
        return rtExpression
    }
    return nil
}

