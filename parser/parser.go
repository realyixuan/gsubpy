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

    for p.l.PeekNextToken().TokenType != token.EOF {
        stmt := p.parsingStatement()
        stmts = append(stmts, stmt)
    }
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
    assignment.Value = p.parsingExpression(0)
    p.l.ReadNextToken()
    return &assignment
}

func (p *Parser)parsingExpression(precedence int) ast.Expression {
    left := p.prefixFn()
    p.l.ReadNextToken()

    var infixPrecedence int
    if p.l.CurToken.Literals == "+" {
        infixPrecedence = 1
    } else {
        infixPrecedence = 2
    }

    for p.l.CurToken.TokenType != token.LINEFEED && p.l.CurToken.TokenType != token.EOF && infixPrecedence > precedence {
        left = p.infixFn(left)
    }

    return left
}

func (p *Parser) prefixFn() ast.Expression {
    if p.l.CurToken.TokenType == token.IDENTIFIER {
        return &ast.IdentifierExpression{p.l.CurToken}
    } else if p.l.CurToken.TokenType == token.NUMBER {
        return &ast.NumberExpression{p.l.CurToken}
    }
    return nil
}

func (p *Parser) infixFn(expression ast.Expression) ast.Expression {
    curTokenType := p.l.CurToken.TokenType
    p.l.ReadNextToken()
    switch curTokenType {
    case token.PLUS:
        return &ast.PlusExpression{
            Left: expression,
            Right: p.parsingExpression(1),
        }
    case token.MUL:
        return &ast.MulExpression{
            Left: expression,
            Right: p.parsingExpression(2),
        }
    }

    return nil
}

