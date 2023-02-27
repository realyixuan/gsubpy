package main

import (
    "os"
    "bufio"
    "fmt"

    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/evaluator"
    "gsubpy/object"
    "gsubpy/ast"
)

func main() {
    if len(os.Args) == 1 {
        for {
            reader := bufio.NewReader(os.Stdin)
            fmt.Print(">>> ")
            line , _ := reader.ReadString('\n')
            l := lexer.New(line)
            p := parser.New(l)
            stmts := p.Parsing()

            if len(stmts) == 0 {
                //
            } else {
                stmt := stmts[0]
                // ? How to compare interfaces
                switch node := stmt.(type) {
                case *ast.ExpressionStatement:
                    obj := evaluator.Eval(node)
                    fmt.Println(obj.(*object.NumberObject).Value)
                case ast.Statement:
                    evaluator.Exec(stmts)
                }
            }
        }
    }
}

