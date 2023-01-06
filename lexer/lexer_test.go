package lexer

import (
    "testing"
)

func TestLexer(t *testing.T) {
    code := `
a = 1
b = 1

def foo(a, b):
    return a + b

c = foo(a, b)
`
}

