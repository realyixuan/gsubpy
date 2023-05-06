package parser

import (
    "fmt"

    "github.com/realyixuan/gsubpy/ast"
    "github.com/realyixuan/gsubpy/lexer"
    "github.com/realyixuan/gsubpy/token"
    "github.com/realyixuan/gsubpy/evaluator"
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
    p.registerStatementParsingFn(token.BREAK, p.parsingBreakStatement)
    p.registerStatementParsingFn(token.CONTINUE, p.parsingContinueStatement)
    p.registerStatementParsingFn(token.FOR, p.parsingForStatement)
    p.registerStatementParsingFn(token.DEF, p.parsingDefStatement)
    p.registerStatementParsingFn(token.CLASS, p.parsingClassStatement)
    p.registerStatementParsingFn(token.RETURN, p.parsingReturnStatement)
    p.registerStatementParsingFn(token.RAISE, p.parsingRaiseStatement)
    p.registerStatementParsingFn(token.ASSERT, p.parsingAssertStatement)

    // a trick, if the a statement doesn't belong to any one above, then
    // it default to the expression-statement, using token IDENTIFIER to 
    // denote it.
    p.registerStatementParsingFn(token.IDENTIFIER, p.parsingExpressionStatement)

    // register expression-parsing function
    p.registerPrefixFn(token.IDENTIFIER, p.getIDENTIFIERPrefix)
    p.registerPrefixFn(token.INTEGER, p.getINTEGERPrefix)
    p.registerPrefixFn(token.STRING, p.getSTRINGPrefix)
    p.registerPrefixFn(token.LBRACKET, p.getLBRACKETPrefix)
    p.registerPrefixFn(token.LBRACE, p.getLBRACEPrefix)
    p.registerPrefixFn(token.LPAREN, p.getLPARENPrefix)
    p.registerPrefixFn(token.NOT, p.getNOTPrefix)

    p.registerInfixFn(token.DOT, p.getDOTInfix)
    p.registerInfixFn(token.PLUS, p.getPLUSInfix)
    p.registerInfixFn(token.MINUS, p.getMINUSInfix)
    p.registerInfixFn(token.MUL, p.getMULInfix)
    p.registerInfixFn(token.DIVIDE, p.getDIVIDEInfix)
    p.registerInfixFn(token.GT, p.getGTInfix)
    p.registerInfixFn(token.LT, p.getLTInfix)
    p.registerInfixFn(token.EQ, p.getEQInfix)
    p.registerInfixFn(token.NEQ, p.getNEQInfix)
    p.registerInfixFn(token.IN, p.getINInfix)
    p.registerInfixFn(token.NIN, p.getNINInfix)
    p.registerInfixFn(token.IS, p.getISInfix)
    p.registerInfixFn(token.ISN, p.getISNInfix)
    p.registerInfixFn(token.LPAREN, p.getLPARENInfix)
    p.registerInfixFn(token.LBRACKET, p.getLBRACKETInfix)
    p.registerInfixFn(token.AND, p.getANDInfix)
    p.registerInfixFn(token.OR, p.getORInfix)

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
        // don't have to care here, in this level.
        if isGTIndents(indents, p.l.Indents) {
            break
        } else if isLTIndents(indents, p.l.Indents) {
            panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nIndentError: wrong indents", p.l.LineNum, p.l.Line)))
        }

        stmt := p.parsingStatement()
        stmts = append(stmts, stmt)
    }

    return stmts

}

func (p *Parser)parsingStatement() ast.Statement {
    stmtParsingFn := p.getStmtParsingFn()
    if stmtParsingFn == nil {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: ...", p.l.LineNum, p.l.Line)))
    }
    return stmtParsingFn()
}

func (p *Parser) registerStatementParsingFn(tokenType token.TokenType, fn statementParsingFn) {
    p.statementParsingFns[tokenType] = fn
}

func (p *Parser) getStmtParsingFn() statementParsingFn {
    // Because there is no keyword to identify the assignment statement
    // so have to make a judgement for it
    if _, ok := token.Keywords[p.l.CurToken.Literals]; ok {
        return p.statementParsingFns[p.l.CurToken.Type]
    }

    if p.isAssignStatement() {
        return p.statementParsingFns[token.ASSIGN]
    } else {
        return p.statementParsingFns[token.IDENTIFIER]
    }
}


func (p *Parser)parsingExpressionStatement() ast.Statement {
    expr := &ast.ExpressionStatement{
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }
    expr.Value = p.parsingExpression(0)
    p.l.ReadNextToken()
    p.l.ReadNextToken()
    return expr
}

func (p *Parser)parsingIfStatement() ast.Statement {
    curIndents := p.l.Indents

    ifStatement := &ast.IfStatement{
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }
    if p.l.CurToken.Type == token.ELSE {
        ifStatement.Condition = nil
        p.l.ReadNextToken()
    } else {
        p.l.ReadNextToken()
        ifStatement.Condition = p.parsingExpression(0)
        p.l.ReadNextToken()
    }

    if p.l.CurToken.Type == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken()
    }

    if !isGTIndents(p.l.Indents, curIndents) {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nIndentError: wrong Indents", p.l.LineNum, p.l.Line)))
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

    stmt := &ast.WhileStatement{
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }
    p.l.ReadNextToken()
    stmt.Condition = p.parsingExpression(0)
    p.l.ReadNextToken()

    if p.l.CurToken.Type == token.COLON {
        p.l.ReadNextToken()
        p.l.ReadNextToken()
    }

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nIndentError: wrong Indents", p.l.LineNum, p.l.Line)))
    }
    
    stmt.Body = p.parsing(p.l.Indents)

    return stmt
}

func (p *Parser)parsingBreakStatement() ast.Statement {
    stmt := &ast.BreakStatement{
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }
    p.l.ReadNextToken()
    p.readNotLineFeedToken()
    return stmt
}

func (p *Parser)parsingContinueStatement() ast.Statement {
    stmt := &ast.ContinueStatement{
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }
    p.l.ReadNextToken()
    p.readNotLineFeedToken()
    return stmt
}

func (p *Parser)parsingForStatement() ast.Statement {
    curIndents := p.l.Indents

    stmt := &ast.ForStatement{
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }

    p.l.ReadNextToken()
    var idents []token.Token
    for ; p.l.CurToken.Type != token.IN; p.l.ReadNextToken() {
        if p.l.CurToken.Type != token.IDENTIFIER {
            panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: ", p.l.LineNum, p.l.Line)))
        }

        idents = append(idents, p.l.CurToken)

        if p.l.PeekNextToken().Type == token.COMMA {
            p.l.ReadNextToken()
        }
    }

    if p.l.CurToken.Type != token.IN {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: expect in", p.l.LineNum, p.l.Line)))
    }

    stmt.Identifiers = idents
    p.l.ReadNextToken()

    stmt.Target = p.parsingExpression(0)

    p.l.ReadNextToken()
    if p.l.CurToken.Type != token.COLON {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: expect :", p.l.LineNum, p.l.Line)))
    }
    p.l.ReadNextToken()
    p.l.ReadNextToken()

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nIndentError: wrong Indents", p.l.LineNum, p.l.Line)))
    }
    
    stmt.Body = p.parsing(p.l.Indents)

    return stmt
}

func (p *Parser)parsingDefStatement() ast.Statement {
    curIndents := p.l.Indents

    p.l.ReadNextToken()
    stmt := &ast.DefStatement{
        Name: p.l.CurToken,
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }

    p.l.ReadNextToken()
    if p.l.CurToken.Type != token.LPAREN {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: wrong syntax", p.l.LineNum, p.l.Line)))
    }

    p.l.ReadNextToken()
    stmt.Params = p.parsingDefParams()

    if p.l.CurToken.Type != token.RPAREN {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: wrong syntax", p.l.LineNum, p.l.Line)))
    }

    p.l.ReadNextToken()
    if p.l.CurToken.Type == token.COLON {
        p.l.ReadNextToken()
        p.readNotLineFeedToken()
    }

    if isLTIndents(p.l.Indents, curIndents) && isEQIndents(p.l.Indents, curIndents) {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nIndentError: wrong Indents", p.l.LineNum, p.l.Line)))
    }
    
    stmt.Body = p.parsing(p.l.Indents)

    return stmt
}

func (p *Parser)parsingClassStatement() ast.Statement {
    classIndents := p.l.Indents

    p.l.ReadNextToken()
    stmt := &ast.ClassStatement{
        Name: p.l.CurToken,
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }

    // inheritance
    if p.l.ReadNextToken(); p.l.CurToken.Type == token.LPAREN {
        p.l.ReadNextToken()

        // supposed to be identifier token
        stmt.Parent = p.l.CurToken

        if p.l.ReadNextToken(); p.l.CurToken.Type != token.RPAREN {
            panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: class define wrong syntax", p.l.LineNum, p.l.Line)))
        }
        p.l.ReadNextToken()
    }

    if p.l.CurToken.Type != token.COLON {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: class define wrong syntax", p.l.LineNum, p.l.Line)))
    }

    p.readNotLineFeedToken()

    internalIndents := p.l.Indents
    
    if !isGTIndents(internalIndents, classIndents) {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nIndentError: in class wrong Indents", p.l.LineNum, p.l.Line)))
    }

    for isEQIndents(internalIndents, p.l.Indents) {
        for _, st := range p.parsing(internalIndents) {
            stmt.Body = append(stmt.Body, st)
        }
    }

    return stmt
}

func (p *Parser)parsingReturnStatement() ast.Statement {
    p.l.ReadNextToken()
    stmt := &ast.ReturnStatement{
        Value: p.parsingExpression(LOWEST),
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }
    p.l.ReadNextToken()
    return stmt
}

func (p *Parser)parsingRaiseStatement() ast.Statement {
    p.l.ReadNextToken()
    stmt := &ast.RaiseStatement{
        Value: p.parsingExpression(LOWEST),
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }
    p.l.ReadNextToken()
    p.l.ReadNextToken()
    return stmt
}

func (p *Parser)parsingAssertStatement() ast.Statement {
    p.l.ReadNextToken()
    stmt := &ast.AssertStatement{
        Condition: p.parsingExpression(LOWEST),
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }

    p.l.ReadNextToken()
    if p.l.CurToken.Type == token.COMMA {
        p.l.ReadNextToken()
        stmt.Msg = p.parsingExpression(LOWEST)
        p.l.ReadNextToken()
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
    expr := p.parsingExpression(LOWEST)
    switch expr.(type) {
    case *ast.IdentifierExpression:
    case *ast.AttributeExpression:
    case *ast.SubscriptExpression:
    default:
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: invalid syntax", p.l.LineNum, p.l.Line)))
    }

    assignment := ast.AssignStatement{
        Target: expr,
        Literals: ast.Literals{LineNum: p.l.LineNum, Line: p.l.Line},
    }

    p.l.ReadNextToken()

    symbol := p.l.CurToken.Type

    p.l.ReadNextToken()

    var val ast.Expression
    rExpr := p.parsingExpression(0)
    if symbol == token.PLUSASSIGN {
        val = &ast.PlusExpression{
            Left: assignment.Target,
            Right: rExpr,
            }
    } else if symbol == token.MINUSASSIGN {
        val = &ast.MinusExpression{
            Left: assignment.Target,
            Right: rExpr,
            }
    } else if symbol == token.MULASSIGN {
        val = &ast.MulExpression{
            Left: assignment.Target,
            Right: rExpr,
            }
    } else if symbol == token.DIVIDEASSIGN {
        val = &ast.DivideExpression{
            Left: assignment.Target,
            Right: rExpr,
            }
    } else {
        val = rExpr
    }

    assignment.Value = val
    p.l.ReadNextToken()
    p.l.ReadNextToken()
    return &assignment
}

func (p *Parser)parsingExpression(precedence int) ast.Expression {
    prefixFn := p.getPrefixFn()
    left := prefixFn()

    // get the corresponding precedence of Token, and 
    // if it doesn't be in definition, like EOF, COLON, and others
    // that means the ending of the expression, it
    // will return LOWEST precedence
    for getPrecedence(p.l.PeekNextToken().Type) > precedence {
        p.l.ReadNextToken()

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
        p.readNotLineFeedToken()
        params = append(params, param)
        if p.l.CurToken.Type == token.COMMA {
            p.readNotLineFeedToken()
        }
    }

    return params
}

func (p *Parser)parsingSubscript(precedence int) ast.Expression {
    var val ast.Expression
    for p.l.CurToken.Type != token.RBRACKET && p.l.CurToken.Type != token.EOF {
        val = p.parsingExpression(LOWEST)
        p.readNotLineFeedToken()
    }

    return val
}

func (p *Parser)isWhiteLine() bool {
    if p.l.CurToken.Type != token.LINEFEED {
        return false
    }
    return true
}

func (p *Parser) getPrefixFn() prefPrefixFn {
    prefixFn, ok := p.prefixFns[p.l.CurToken.Type] 
    if !ok {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: invalid syntax", p.l.LineNum, p.l.Line)))
    }
    return prefixFn
}

func (p *Parser) getInfixFn() prefInfixFn {
    infixFn := p.infixFns[p.l.CurToken.Type] 
    return infixFn
}

func getPrecedence(tok token.TokenType) int {
    switch tok {
    case token.DOT:
        return ATTR
    case token.LPAREN:
        return CALL
    case token.LBRACKET:
        return CALL
    case token.LT:
        return COMPARISON
    case token.GT:
        return COMPARISON
    case token.EQ:
        return COMPARISON
    case token.NEQ:
        return COMPARISON
    case token.IN:
        return COMPARISON
    case token.NIN:
        return COMPARISON
    case token.IS:
        return COMPARISON
    case token.ISN:
        return COMPARISON
    case token.PLUS:
        return SUM
    case token.MINUS:
        return SUM
    case token.MUL:
        return PRODUCT
    case token.DIVIDE:
        return PRODUCT
    case token.AND:
        return AND
    case token.OR:
        return OR
    case token.NOT:
        return NOT
    default:
        return LOWEST
    }

}

func (p *Parser) getIDENTIFIERPrefix() ast.Expression {
    return &ast.IdentifierExpression{p.l.CurToken}
}

func (p *Parser) isAssignStatement() bool {
    cl := *p.l
    defer func() {
        p.l = &cl
    }()

    p.parsingExpression(LOWEST)

    p.l.ReadNextToken()

    if p.isAssignType(p.l.CurToken.Type) {
        return true
    } else {
        return false
    }
}

func (p *Parser) isAssignType(tokType token.TokenType) bool {
    if tokType == token.ASSIGN ||
       tokType == token.PLUSASSIGN ||
       tokType == token.MINUSASSIGN ||
       tokType == token.MULASSIGN ||
       tokType == token.DIVIDEASSIGN {
        return true
    } else {
        return false
    }
}

func (p *Parser) getINTEGERPrefix() ast.Expression {
    return &ast.NumberExpression{p.l.CurToken}
}

func (p *Parser) getSTRINGPrefix() ast.Expression {
    return &ast.StringExpression{p.l.CurToken}
}

func (p *Parser) getLBRACKETPrefix() ast.Expression {
    expr := &ast.ListExpression{}

    p.readNotLineFeedToken()
    for p.l.CurToken.Type != token.EOF && p.l.CurToken.Type != token.RBRACKET {
        expr.Items = append(expr.Items, p.parsingExpression(LOWEST))
        p.readNotLineFeedToken()

        if p.l.CurToken.Type == token.COMMA {
            p.readNotLineFeedToken()
        }
    }

    return expr
}

func (p *Parser) getLBRACEPrefix() ast.Expression {
    expr := &ast.DictExpression{}

    p.readNotLineFeedToken()
    for p.l.CurToken.Type != token.EOF && p.l.CurToken.Type != token.RBRACE {
        expr.Keys = append(expr.Keys, p.parsingExpression(LOWEST))
        p.readNotLineFeedToken()

        if p.l.CurToken.Type != token.COLON {
            panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: there is a syntax error in dict", p.l.LineNum, p.l.Line)))
        }
        p.readNotLineFeedToken()

        expr.Vals = append(expr.Vals, p.parsingExpression(LOWEST))
        p.readNotLineFeedToken()

        if p.l.CurToken.Type == token.COMMA {
            p.readNotLineFeedToken()
        }
    }

    if p.l.CurToken.Type != token.RBRACE {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: there is a syntax error in dict, expect '}'", p.l.LineNum, p.l.Line)))
    }

    return expr
}

func (p *Parser) getLPARENPrefix() ast.Expression {
    p.readNotLineFeedToken()

    expr := p.parsingExpression(LOWEST)

    if p.l.PeekNextToken().Type == token.RPAREN {
        p.readNotLineFeedToken()
    } else {
        panic(evaluator.Error(fmt.Sprintf("line %v\n\t%s\nSyntaxError: expect ')'", p.l.LineNum, p.l.Line)))
    }
    return expr
}

func (p *Parser) getNOTPrefix() ast.Expression {
    expr := &ast.NotExpression{}

    p.l.ReadNextToken()
    expr.Expr = p.parsingExpression(getPrecedence(token.NOT))

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

func (p *Parser) getDOTInfix(left ast.Expression) ast.Expression {
    expr := &ast.AttributeExpression{Expr: left}

    expr.Attr = p.l.CurToken

    return expr
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

func (p *Parser) getEQInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.EQ, "=="},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.EQ)),
    }
}

func (p *Parser) getNEQInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.NEQ, "!="},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.NEQ)),
    }
}

func (p *Parser) getINInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.IN, "in"},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.IN)),
    }
}

func (p *Parser) getNINInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.NIN, "not in"},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.NIN)),
    }
}

func (p *Parser) getISInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.IS, "is"},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.IS)),
    }
}

func (p *Parser) getISNInfix(left ast.Expression) ast.Expression {
    return &ast.ComparisonExpression{
        Operator: token.Token{token.ISN, "is not"},
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.ISN)),
    }
}

func (p *Parser) getLPARENInfix(left ast.Expression) ast.Expression {
    return &ast.CallExpression{
        Name: left,
        Params: p.parsingCallParams(getPrecedence(token.LPAREN)),
    }
}

func (p *Parser) getLBRACKETInfix(left ast.Expression) ast.Expression {
    return &ast.SubscriptExpression{
        Target: left,
        Val: p.parsingSubscript(getPrecedence(token.LBRACKET)),
    }
}

func (p *Parser) getANDInfix(left ast.Expression) ast.Expression {
    return &ast.AndExpression{
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.AND)),
    }
}

func (p *Parser) getORInfix(left ast.Expression) ast.Expression {
    return &ast.OrExpression{
        Left: left,
        Right: p.parsingExpression(getPrecedence(token.OR)),
    }
}

func (p *Parser) registerInfixFn(tok token.TokenType, fn prefInfixFn) {
    p.infixFns[tok] = fn
}

func (p *Parser) readNotLineFeedToken() {
    for p.l.ReadNextToken(); p.l.CurToken.Type == token.LINEFEED; p.l.ReadNextToken() {}
}

const (
    LOWEST int = iota
    OR
    AND
    NOT
    COMPARISON
    SUM
    PRODUCT
    CALL
    ATTR
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

