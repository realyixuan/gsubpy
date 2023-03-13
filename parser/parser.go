package parser

import (
    "gsubpy/ast"
    "gsubpy/lexer"
    "gsubpy/token"
    "gsubpy/object"
)

type (
    prefPrefixFn func() ast.Expression
    prefInfixFn func(ast.Expression) ast.Expression
    statementParsingFn func() ast.Statement
)

type Parser struct {
    l                       *lexer.Lexer
    prefixFns               map[token.TokenType]prefPrefixFn
    infixFns                map[token.TokenType]prefInfixFn
    statementParsingFns     map[token.TokenType]statementParsingFn
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l:                      l,
        prefixFns:              make(map[token.TokenType]prefPrefixFn),
        infixFns:               make(map[token.TokenType]prefInfixFn),
        statementParsingFns:    make(map[token.TokenType]statementParsingFn),
    }

    // register statement-parsing function
    p.registerStatementParsingFn(token.ASSIGN, p.parsingAssignStatement)
    p.registerStatementParsingFn(token.IF, p.parsingIfStatement)
    p.registerStatementParsingFn(token.WHILE, p.parsingWhileStatement)
    p.registerStatementParsingFn(token.DEF, p.parsingDefStatement)
    p.registerStatementParsingFn(token.RETURN, p.parsingReturnStatement)

    // a trick, if the a statement doesn't belong to any one above, then
    // it default to the expression-statement, using token IDENTIFIER to 
    // denote it.
    p.registerStatementParsingFn(token.IDENTIFIER, p.parsingExpressionStatement)

    // register expression-parsing function
    p.registerPrefixFn(token.IDENTIFIER, p.getIDENTIFIERPrefix)
    p.registerPrefixFn(token.NUMBER, p.getNUMBERPrefix)
    p.registerPrefixFn(token.STRING, p.getSTRINGPrefix)
    p.registerPrefixFn(token.LBRACKET, p.getLBRACKETPrefix)

    p.registerInfixFn(token.PLUS, p.getPLUSInfix)
    p.registerInfixFn(token.MINUS, p.getMINUSInfix)
    p.registerInfixFn(token.MUL, p.getMULInfix)
    p.registerInfixFn(token.DIVIDE, p.getDIVIDEInfix)
    p.registerInfixFn(token.GT, p.getGTInfix)
    p.registerInfixFn(token.LT, p.getLTInfix)
    p.registerInfixFn(token.LPAREN, p.getLPARENInfix)

    return p
}

func (p *Parser)Parsing() []ast.Statement {
    return p.parsing(NO_INDENTS)
}

func (p *Parser)parsing(indents string) []ast.Statement {
    // return statements
    var stmts []ast.Statement

    for p.l.CurToken.Type != token.EOF {
        if p.isWhiteLine() {
            p.l.ReadNextToken()
            continue
        }

        // In same block code, all statements should have same
        // indents, and having shorter indents MAY BE also legitimate,
        // perhaps it belongs to upper block, or not. In either case,
        // don't have to care here, in the level.
        if isGTIndents(indents, p.l.Indents) {
            break
        } else if isLTIndents(indents, p.l.Indents) {
            panic(&object.ExceptionObject{"IndentError: wrong indents"})
        }

        stmt := p.parsingStatement()
        stmts = append(stmts, stmt)
    }

    return stmts

}

func (p *Parser)parsingStatement() ast.Statement {
    // get corresponding statement parser, otherwise return nil
    stmtParsingFn := p.getStmtParsingFn()
    if stmtParsingFn == nil {
        panic(&object.ExceptionObject{"SyntaxError: ..."})
    }
    // then call it to return statement
    return stmtParsingFn()


    // switch p.l.CurToken.Type {
    // case token.IDENTIFIER:
    //     if p.l.PeekNextToken().Type == token.ASSIGN {
    //         return p.parsingAssignStatement()
    //     } else {
    //         return p.parsingExpressionStatement()
    //     }
    // case token.IF:
    //     return p.parsingIfStatement()
    // case token.WHILE:
    //     return p.parsingWhileStatement()
    // case token.DEF:
    //     return p.parsingDefStatement()
    // case token.RETURN:
    //     return p.parsingReturnStatement()
    // default:
    //     return p.parsingExpressionStatement()
    // }
    // return nil
}

func (p *Parser) registerStatementParsingFn(tokenType token.TokenType, fn statementParsingFn) {
    p.statementParsingFns[tokenType] = fn
}

func (p *Parser) getStmtParsingFn() statementParsingFn {
    // Because there is no keyword to identify the assignment statement
    // so have to make a judgement for it
    if p.l.CurToken.Type == token.IDENTIFIER && p.l.PeekNextToken().Type == token.ASSIGN {
        return p.statementParsingFns[token.ASSIGN]
    } else if _, ok := p.statementParsingFns[p.l.CurToken.Type]; ok {
        return p.statementParsingFns[p.l.CurToken.Type]
    } else {
        return p.statementParsingFns[token.IDENTIFIER]
    }
}


func (p *Parser)parsingExpressionStatement() ast.Statement {
    val := p.parsingExpression(0)
    p.l.ReadNextToken()
    return &ast.ExpressionStatement{
        Value: val,
    }
}

func (p *Parser)parsingIfStatement() ast.Statement {
    curIndents := p.l.Indents

    ifStatement := &ast.IfStatement{}
    if p.l.CurToken.Type == token.ELSE {
        ifStatement.Condition = nil
        p.l.ReadNextToken()
    } else {
        p.l.ReadNextToken()
        ifStatement.Condition = p.parsingExpression(0)
    }

    if p.l.CurToken.Type == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken()
    }

    if !isGTIndents(p.l.Indents, curIndents) {
        panic(&object.ExceptionObject{"IndentError: wrong Indents"})
    }
    
    ifStatement.Body = p.parsing(p.l.Indents)

    if isEQIndents(p.l.Indents, curIndents) {
        ifStatement.Else = p.parsingElifOrElseStatement()
    }

    return ifStatement
}

func (p *Parser)parsingElifOrElseStatement() ast.Statement {
    if p.l.CurToken.Type == token.ELIF || p.l.CurToken.Type == token.ELSE{
        return p.parsingIfStatement()
    }

    return nil
}

func (p *Parser)parsingWhileStatement() ast.Statement {
    curIndents := p.l.Indents

    stmt := &ast.WhileStatement{}
    p.l.ReadNextToken()
    stmt.Condition = p.parsingExpression(0)

    if p.l.CurToken.Type == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken()
    }

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic(&object.ExceptionObject{"IndentError: wrong Indents"})
    }
    
    stmt.Body = p.parsing(p.l.Indents)

    return stmt
}

func (p *Parser)parsingDefStatement() ast.Statement {
    curIndents := p.l.Indents

    p.l.ReadNextToken()
    stmt := &ast.DefStatement{
        Name: p.l.CurToken,
    }

    p.l.ReadNextToken()
    if p.l.CurToken.Type != token.LPAREN {
        panic(&object.ExceptionObject{"SyntaxError: wrong syntax"})
    }

    p.l.ReadNextToken()
    stmt.Params = p.parsingDefParams()

    if p.l.CurToken.Type != token.RPAREN {
        panic(&object.ExceptionObject{"SyntaxError: wrong syntax"})
    }

    p.l.ReadNextToken()
    if p.l.CurToken.Type == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken() // skip over '\n'
    }

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic(&object.ExceptionObject{"IndentError: wrong Indents"})
    }
    
    stmt.Body = p.parsing(p.l.Indents)

    return stmt
}

func (p *Parser)parsingReturnStatement() ast.Statement {
    p.l.ReadNextToken()
    stmt := &ast.ReturnStatement{
        Value: p.parsingExpression(LOWEST),
    }
    return stmt
}

func (p *Parser)parsingDefParams() []token.Token {
    var params []token.Token

    for p.l.CurToken.Type != token.RPAREN {
        params = append(params, p.l.CurToken)
        p.l.ReadNextToken()
        if p.l.CurToken.Type == token.COMMA {
            p.l.ReadNextToken()
        }
    }

    return params
}

func (p *Parser)parsingAssignStatement() ast.Statement {
    assignment := ast.AssignStatement{Identifier: p.l.CurToken}
    p.l.ReadNextToken()
    p.l.ReadNextToken()
    assignment.Value = p.parsingExpression(0)
    p.l.ReadNextToken()
    return &assignment
}

func (p *Parser)parsingExpression(precedence int) ast.Expression {
    prefixFn := p.getPrefixFn()
    left := prefixFn()

    p.l.ReadNextToken()

    // get the corresponding precedence of Token, and 
    // if it doesn't be in definition, like EOF, COLON, and others
    // that means the ending of the expression, it
    // will return LOWEST precedence
    for getPrecedence(p.l.CurToken.Type) > precedence {
        infixFn := p.getInfixFn()

        p.l.ReadNextToken()

        left = infixFn(left)
    }

    return left
}

func (p *Parser)parsingCallParams(precedence int) []ast.Expression {
    var params []ast.Expression

    for p.l.CurToken.Type != token.RPAREN && p.l.CurToken.Type != token.EOF {
        param := p.parsingExpression(LOWEST)
        params = append(params, param)
        if p.l.CurToken.Type == token.COMMA {
            p.l.ReadNextToken()
        }
    }

    p.l.ReadNextToken()

    return params
}

func (p *Parser)isWhiteLine() bool {
    if p.l.CurToken.Type != token.LINEFEED {
        return false
    }
    return true
}

func (p *Parser) getPrefixFn() prefPrefixFn {
    prefixFn := p.prefixFns[p.l.CurToken.Type] 
    return prefixFn
}

func (p *Parser) getInfixFn() prefInfixFn {
    infixFn := p.infixFns[p.l.CurToken.Type] 
    return infixFn
}

func getPrecedence(tok token.TokenType) int {
    switch tok {
    case token.LPAREN:
        return CALL
    case token.LT:
        return COMPARISON
    case token.GT:
        return COMPARISON
    case token.PLUS:
        return SUM
    case token.MINUS:
        return SUM
    case token.MUL:
        return PRODUCT
    case token.DIVIDE:
        return PRODUCT
    default:
        return LOWEST
    }

}

func (p *Parser) getIDENTIFIERPrefix() ast.Expression {
    return &ast.IdentifierExpression{p.l.CurToken}
    
}

func (p *Parser) getNUMBERPrefix() ast.Expression {
    return &ast.NumberExpression{p.l.CurToken}
}

func (p *Parser) getSTRINGPrefix() ast.Expression {
    return &ast.StringExpression{p.l.CurToken}
}

func (p *Parser) getLBRACKETPrefix() ast.Expression {
    expr := &ast.ListExpression{}

    p.l.ReadNextToken()
    for p.l.CurToken.Type != token.EOF && p.l.CurToken.Type != token.RBRACKET {
        expr.Items = append(expr.Items, p.parsingExpression(LOWEST))

        if p.l.CurToken.Type == token.COMMA {
            p.l.ReadNextToken()
        }
    }

    return expr
}

func (p *Parser) registerPrefixFn(tok token.TokenType, fn prefPrefixFn) {
    p.prefixFns[tok] = fn
}

func (p *Parser) getPLUSInfix(left ast.Expression) ast.Expression {
    return &ast.PlusExpression{
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.PLUS)),
    }
}

func (p *Parser) getMINUSInfix(left ast.Expression) ast.Expression {
    return &ast.MinusExpression{
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.MINUS)),
    }
}

func (p *Parser) getMULInfix(left ast.Expression) ast.Expression {
    return &ast.MulExpression{
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.MUL)),
    }
}

func (p *Parser) getDIVIDEInfix(left ast.Expression) ast.Expression {
    return &ast.DivideExpression{
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.DIVIDE)),
    }
}

func (p *Parser) getGTInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.GT, ">"},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.GT)),
    }
}

func (p *Parser) getLTInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.LT, "<"},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.LT)),
    }
}

func (p *Parser) getLPARENInfix(left ast.Expression) ast.Expression {
    return &ast.FunctionCallExpression{
        Name: left,
        Params: p.parsingCallParams(getPrecedence(token.LPAREN)),
    }
}

func (p *Parser) registerInfixFn(tok token.TokenType, fn prefInfixFn) {
    p.infixFns[tok] = fn
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

