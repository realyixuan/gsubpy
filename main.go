package main

import (
    "os"
    "fmt"

    "gsubpy/repl"
    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/evaluator"
    "gsubpy/object"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            // get interface out from interface ?
            switch exception := r.(type) {
            case object.Exception:
                fmt.Println(exception.ErrorMsg())
            default:
                fmt.Println("interpreter error:", exception)
            }
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

