package ast

import (
    "github.com/realyixuan/gsubpy/token"
)

type Statement interface {
    getStatement()
    GetLiterals() Literals
}

type Expression interface {
    getExpression()
}

type Program struct {
    stmts []Statement
}

type Literals struct {
    LineNum     int
    Line        string
}

type AssignStatement struct {
    Target      Expression
    Value       Expression
    Literals
}

func (as *AssignStatement) getStatement() {}
func (as *AssignStatement) GetLiterals() Literals {return as.Literals}

type ExpressionStatement struct {
    Value       Expression
    Literals
}

func (es *ExpressionStatement) getStatement() {}
func (es *ExpressionStatement) getExpression() {}
func (es *ExpressionStatement) GetLiterals() Literals {return es.Literals}

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

type DictExpression struct {
    Keys   []Expression
    Vals   []Expression
}

func (de *DictExpression) getExpression() {}

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

type AndExpression struct {
    Left    Expression
    Right   Expression
}
func (ae *AndExpression) getExpression() {}

type OrExpression struct {
    Left    Expression
    Right   Expression
}
func (oe *OrExpression) getExpression() {}

type NotExpression struct {
    Expr   Expression
}
func (ne *NotExpression) getExpression() {}

type IfStatement struct {
    Condition   Expression
    Body        []Statement
    Else        Statement
    Literals
}

func (ie *IfStatement) getStatement() {}
func (ie *IfStatement) GetLiterals() Literals {return ie.Literals}

type WhileStatement struct {
    Condition   Expression
    Body        []Statement
    Literals
}

func (ws *WhileStatement) getStatement() {}
func (ws *WhileStatement) GetLiterals() Literals {return ws.Literals}

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
    Literals
}

func (fs *DefStatement) getStatement() {}
func (fs *DefStatement) GetLiterals() Literals {return fs.Literals}

type ReturnStatement struct {
    Value   Expression
    Literals
}

func (rs *ReturnStatement) getStatement() {}
func (rs *ReturnStatement) GetLiterals() Literals {return rs.Literals}

type CallExpression struct {
    Name        Expression
    Params      []Expression
}

func (ce *CallExpression) getExpression() {}

type ClassStatement struct {
    Name    token.Token
    Body    []Statement
    Parent  token.Token
    Literals
}

func (cs *ClassStatement) getStatement() {}
func (cs *ClassStatement) GetLiterals() Literals {return cs.Literals}

type AttributeExpression struct {
    Expr    Expression
    Attr    token.Token
}

func (de *AttributeExpression) getExpression() {}

