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
    Literals    string
}

type NumberExpression struct {
    Value   token.Token
}

func (ne *NumberExpression) getExpression() {}

type PlusExpression struct {
    Left    Expression
    Right   Expression
}

func (pe *PlusExpression) getExpression() {}

type MulExpression struct {
    Left    Expression
    Right   Expression
}

func (me *MulExpression) getExpression() {}

