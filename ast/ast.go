package ast

import (
    // "bytes"

    "gsubpy/token"
)

type Node interface {
    TokenLiteral() string
    // String() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type ExpressionStatement struct {
    Token       token.Token
    Expression  Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {return es.Token.Literal}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

// func (p *Program) String() string {
//     var out bytes.Buffer
//
//     for _, s := range p.Statements {
//         out.WriteString(s.String())
//     }
// }

type AssignmentStatement struct {
    Token token.Token
    Name *Identifier
    Value Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) TokenLiteral() string {return as.Token.Literal}
// func (as *AssignmentStatement) String() string {
//     var out bytes.Buffer
// 
//     out.WriteString(as.TokenLiteral() + " ")
//     out.WriteString(as.Name.String())
//     out.WriteString(" = ")
// 
//     if as.Value != nil {
//         out.WriteString(as.Value.String())
//     }
// }

type Identifier struct {
    Token token.Token
    Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {return i.Token.Literal}

type InfixExpression struct {
    Token       token.Token
    Left        Expression
    Operator    string
    Right       Expression
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string {return oe.Token.Literal}

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) expressionNode()  {}
func (il *IntegerLiteral) TokenLiteral() string  {return il.Token.Literal}

