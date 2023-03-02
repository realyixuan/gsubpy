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
    return p.parsing(NO_INDENTS)
}

func (p *Parser)parsing(indents string) []ast.Statement {
    // return statements
    var stmts []ast.Statement

    for p.l.CurToken.TokenType != token.EOF {
        if p.isWhiteLine() {
            p.l.ReadNextToken()
            continue
        }
        if !isEQIndents(indents, p.l.Indents) {
            // precisely speaking, whether p.l.Indents == indents or p.l.Indents < indents
            // but there now isn't error handling
            // so omit the error
            break
        }
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
        } else {
            return p.parsingExpressionStatement()
        }
    case token.NUMBER:
        return p.parsingExpressionStatement()
    case token.IF:
        return p.parsingIfStatement()
    case token.WHILE:
        return p.parsingWhileStatement()
    case token.DEF:
        return p.parsingDefStatement()
    case token.RETURN:
        return p.parsingReturnStatement()
    }
    return nil
}

func (p *Parser)parsingExpressionStatement() *ast.ExpressionStatement {
    val := p.parsingExpression(0)
    p.l.ReadNextToken()
    return &ast.ExpressionStatement{
        Value: val,
    }
}

func (p *Parser)parsingIfStatement() *ast.IfStatement {
    curIndents := p.l.Indents

    ifStatement := &ast.IfStatement{}
    p.l.ReadNextToken()
    ifStatement.Condition = p.parsingExpression(0)

    if p.l.CurToken.TokenType == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken()
    }

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic("wrong indents")
    }
    
    ifStatement.Body = p.parsing(p.l.Indents)

    return ifStatement
}

func (p *Parser)parsingWhileStatement() *ast.WhileStatement {
    curIndents := p.l.Indents

    stmt := &ast.WhileStatement{}
    p.l.ReadNextToken()
    stmt.Condition = p.parsingExpression(0)

    if p.l.CurToken.TokenType == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken()
    }

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic("wrong indents")
    }
    
    stmt.Body = p.parsing(p.l.Indents)

    return stmt
}

func (p *Parser)parsingDefStatement() *ast.DefStatement {
    curIndents := p.l.Indents

    p.l.ReadNextToken()
    stmt := &ast.DefStatement{
        Name: p.l.CurToken,
    }

    p.l.ReadNextToken()
    if p.l.CurToken.TokenType != token.LPAREN {
        panic("wrong syntax")
    }

    p.l.ReadNextToken()
    stmt.Params = p.parsingDefParams()

    if p.l.CurToken.TokenType != token.RPAREN {
        panic("wrong syntax")
    }

    p.l.ReadNextToken()
    if p.l.CurToken.TokenType == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken() // skip over '\n'
    }

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic("wrong indents")
    }
    
    stmt.Body = p.parsing(p.l.Indents)

    return stmt
}

func (p *Parser)parsingReturnStatement() *ast.ReturnStatement {
    p.l.ReadNextToken()
    stmt := &ast.ReturnStatement{
        Value: p.parsingExpression(LOWEST),
    }
    return stmt
}

func (p *Parser)parsingDefParams() []token.Token {
    var params []token.Token

    for p.l.CurToken.TokenType != token.RPAREN {
        params = append(params, p.l.CurToken)
        p.l.ReadNextToken()
        if p.l.CurToken.TokenType == token.COMMA {
            p.l.ReadNextToken()
        }
    }

    return params
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

    for p.l.CurToken.TokenType != token.LINEFEED && p.l.CurToken.TokenType != token.COLON && p.l.CurToken.TokenType != token.EOF && getPrecedence(p.l.CurToken.Literals) > precedence {
        left = p.infixFn(left)
    }

    return left
}

func (p *Parser)parsingCallParams(precedence int) []ast.Expression {
    var params []ast.Expression

    for p.l.CurToken.TokenType != token.RPAREN && p.l.CurToken.TokenType != token.EOF{
        param := p.parsingExpression(LOWEST)
        params = append(params, param)
        if p.l.CurToken.TokenType == token.COMMA {
            p.l.ReadNextToken()
        }
    }

    p.l.ReadNextToken()

    return params
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
    case token.LPAREN:
        return &ast.FunctionCallExpression{
            Name: expression,
            Params: p.parsingCallParams(getPrecedence("(")),
        }
    }

    return nil
}

func (p *Parser)isWhiteLine() bool {
    if p.l.CurToken.TokenType != token.LINEFEED {
        return false
    }
    return true
}

func getPrecedence(literals string) int {
    switch literals {
    case "(":
        return CALL
    case "<":
        return COMPARISON
    case ">":
        return COMPARISON
    case "+":
        return SUM
    case "-":
        return SUM
    case "*":
        return PRODUCT
    case "/":
        return PRODUCT
    default:
        return LOWEST
    }

}

const (
    LOWEST int = iota
    COMPARISON
    SUM
    PRODUCT
    CALL
)

const (
    NO_INDENTS = ""
)

func isEQIndents(indents1, indents2 string) bool {
    return indents1 == indents2
}

func isGTIndents(indents1, indents2 string) bool {
    return len(indents1) > len(indents2) && indents1[:len(indents2)] == indents2
}

func isLTIndents(indents1, indents2 string) bool {
    return len(indents1) < len(indents2) && indents2[:len(indents1)] == indents1
}

