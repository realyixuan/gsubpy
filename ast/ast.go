package ast

import (
    "gsubpy/token"
)

type Statement interface {
    getStatement()
}

type Expression interface {
    getExpression()
}

type Program struct {
    stmts []Statement
}

type AssignStatement struct {
    Identifier  token.Token
    Value       Expression
}

func (as *AssignStatement) getStatement() {}

type ExpressionStatement struct {
    Value       Expression
}

func (es *ExpressionStatement) getStatement() {}
func (es *ExpressionStatement) getExpression() {}

type IdentifierExpression struct {
    Identifier  token.Token
}

func (ie *IdentifierExpression) getExpression() {}

type NumberExpression struct {
    Value   token.Token
}

func (ne *NumberExpression) getExpression() {}

type StringExpression struct {
    Value   token.Token
}

func (se *StringExpression) getExpression() {}

type ListExpression struct {
    Items   []Expression
}

func (le *ListExpression) getExpression() {}

type PlusExpression struct {
    Left    Expression
    Right   Expression
}

func (pe *PlusExpression) getExpression() {}

type MinusExpression struct {
    Left    Expression
    Right   Expression
}

func (me *MinusExpression) getExpression() {}

type MulExpression struct {
    Left    Expression
    Right   Expression
}

func (me *MulExpression) getExpression() {}

type DivideExpression struct {
    Left    Expression
    Right   Expression
}

func (de *DivideExpression) getExpression() {}

type IfStatement struct {
    Condition   Expression
    Body        []Statement
    Else        Statement
}

func (ie *IfStatement) getStatement() {}

type WhileStatement struct {
    Condition   Expression
    Body        []Statement
}

func (ws *WhileStatement) getStatement() {}

type ComparisonExpression struct {
    Operator    token.Token
    Left        Expression
    Right       Expression
}

func (ce *ComparisonExpression) getExpression() {}

type DefStatement struct {
    Name    token.Token
    Params  []token.Token
    Body    []Statement
}

func (fs *DefStatement) getStatement() {}

type ReturnStatement struct {
    Value   Expression
}

func (rs *ReturnStatement) getStatement() {}

type FunctionCallExpression struct {
    Name        Expression
    Params      []Expression
}

func (ce *FunctionCallExpression) getExpression() {}

