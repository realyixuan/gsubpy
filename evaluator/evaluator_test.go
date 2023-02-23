package evaluator

import (
    "testing"
    "gsubpy/lexer"
    "gsubpy/parser"
)

func TestEvaluator(t *testing.T) {
    input := `val = 1`

    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    eval(stmts)
    if E["val"] != 1 {
        t.Errorf("expected %d, got %d", 1, E["val"])
    }
}

func TestEvaluator2(t *testing.T) {
    input := `val = 1 + 1`

    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    eval(stmts)
    if E["val"] != 2 {
        t.Errorf("expected %d, got %d", 2, E["val"])
    }
}

