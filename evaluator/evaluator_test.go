package evaluator

import (
    "testing"
    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/object"
)

func TestEvaluator(t *testing.T) {
    input := `val = 1`

    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    eval(stmts)
    if node := E["val"].(*object.NumberObject); node.Value != 1 {
        t.Errorf("expected %d, got %d", 1, node.Value)
    }
}

func TestEvaluator2(t *testing.T) {
    input := `val = 1 + 1`

    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    eval(stmts)
    if node := E["val"].(*object.NumberObject); node.Value != 2 {
        t.Errorf("expected %d, got %d", 2, node.Value)
    }
}

func TestEvaluator3(t *testing.T) {
    input := `val = 10 + 20 * 2`

    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    eval(stmts)
    if node := E["val"].(*object.NumberObject); node.Value != 50 {
        t.Errorf("expected %d, got %d", 50, node.Value)
    }
}

