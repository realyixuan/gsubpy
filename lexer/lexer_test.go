package lexer

import (
    "gsubpy/token"
)

import (
    "testing"
)

func TestLexer(t *testing.T) {
    input := `val = 1`

    l := &Lexer{input: input}

    tk1 := l.nextToken()
    if tk1.TokenType != token.IDENTIFIER || tk1.Literals != "val" {
        t.Errorf("expected %s, got %s", "val", tk1.Literals)
    }

    tk2 := l.nextToken()
    if tk2.TokenType != token.ASSIGN || tk2.Literals != "=" {
        t.Errorf("expected %s, got %s", "=", tk2.Literals)
    }

    tk3 := l.nextToken()
    if tk3.TokenType != token.NUMBER || tk3.Literals != "1" {
        t.Errorf("expected %s, got %s", "1", tk3.Literals)
    }

    tk4 := l.nextToken()
    if tk4.TokenType != token.EOF || tk4.Literals != "" {
        t.Errorf("expected %s, got %s", "1", string(tk4.Literals))
    }
}

