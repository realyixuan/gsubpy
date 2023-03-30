package main

import (
    "os"

    "gsubpy/repl"
    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/evaluator"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            panic(r)
        }
    }()

    if len(os.Args) == 1 {
        repl.REPLRunning()
    } else {
        data, _ := os.ReadFile(os.Args[1])
        l := lexer.New(string(data))
        p := parser.New(l)
        stmts := p.Parsing()
        evaluator.Exec(stmts, evaluator.NewEnvironment())
    }
}

