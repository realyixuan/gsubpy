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

type IdentifierExpression struct {
    Identifier  token.Token
}

func (ie *IdentifierExpression) getExpression() {}

type NumberExpression struct {
    Value   token.Token
}

func (ne *NumberExpression) getExpression() {}

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

