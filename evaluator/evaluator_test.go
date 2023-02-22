package evaluator

import (
    "testing"
    "gsubpy/lexer"
    "gsubpy/parser"
)

func TestEvaluator(t *testing.T) {
    minput := `val = 1`

    l := &lexer.Lexer{Input: minput}
    stmts := parser.Parsing(l)
    eval(stmts)
    if E["val"] != 1 {
        t.Errorf("expected %d, got %d", 1, E["val"])
    }
}

