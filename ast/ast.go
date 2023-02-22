package ast

import (
    "gsubpy/token"
)

type AssignStatement struct {
    Identifier  token.Token
    Value       token.Token
}

