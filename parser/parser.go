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
    case token.IF:
        return p.parsingIfStatement()
    }
    return nil
}

func (p *Parser)parsingIfStatement() *ast.IfStatement {
    ifStatement := &ast.IfStatement{}
    p.l.ReadNextToken()
    ifStatement.Condition = p.parsingExpression(0)

    if p.l.CurToken.TokenType == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken()
    }

    ifStatement.Body = p.Parsing()

    return ifStatement
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

    infixPrecedence := getPrecedence(p.l.CurToken.Literals)

    for p.l.CurToken.TokenType != token.LINEFEED && p.l.CurToken.TokenType != token.COLON && p.l.CurToken.TokenType != token.EOF && infixPrecedence > precedence {
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
            Right: p.parsingExpression(getPrecedence("+")),
        }
    case token.MINUS:
        return &ast.MinusExpression{
            Left: expression,
            Right: p.parsingExpression(getPrecedence("-")),
        }
    case token.MUL:
        return &ast.MulExpression{
            Left: expression,
            Right: p.parsingExpression(getPrecedence("*")),
        }
    case token.DIVIDE:
        return &ast.DivideExpression{
            Left: expression,
            Right: p.parsingExpression(getPrecedence("/")),
        }
    case token.GT:
        return &ast.ComparisonExpression{
            Operator: token.Token{token.GT, ">"},
            Left: expression,
            Right: p.parsingExpression(getPrecedence(">")),
        }
    case token.LT:
        return &ast.ComparisonExpression{
            Operator: token.Token{token.LT, "<"},
            Left: expression,
            Right: p.parsingExpression(getPrecedence("<")),
        }
    }

    return nil
}

func getPrecedence(literals string) int {
    switch literals {
    case "<":
        return 1
    case ">":
        return 1
    case "+":
        return 2
    case "-":
        return 2
    case "*":
        return 3
    case "/":
        return 3
    default:
        return 0
    }

}

