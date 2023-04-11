package main

import (
    "os"
    "fmt"

    "gsubpy/repl"
    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/evaluator"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            for _, f := range evaluator.Py_traceback.Frames {
                fmt.Println("line", f.LineNum)
                fmt.Println("\t", f.Line)
            }
            fmt.Println(r)
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

