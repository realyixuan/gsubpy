package repl

import (
    "os"
    "bufio"
    "fmt"

    "github.com/realyixuan/gsubpy/lexer"
    "github.com/realyixuan/gsubpy/parser"
    "github.com/realyixuan/gsubpy/evaluator"
    "github.com/realyixuan/gsubpy/ast"
    "github.com/realyixuan/gsubpy/token"
)

func REPLRunning() {
    env := evaluator.NewEnvironment()
    for {
        fmt.Print(">>> ")
        reader := bufio.NewReader(os.Stdin)

        line, err := reader.ReadString('\n')

        if err != nil {
            print("\n")
            break
        }

        tl := lexer.New(line)

        var input string
        if _, ok := token.Keywords[tl.CurToken.Literals]; !ok {
            input = line
        } else {
            input += line
            for {
                fmt.Print("... ")
                line, _ := reader.ReadString('\n')

                if line == "\n" {
                    break
                } 

                input += line
            }
        }

        l := lexer.New(input)
        p := parser.New(l)
        stmts := p.Parsing()

        if len(stmts) == 0 {
        } else if len(stmts) > 1{
            panic("invalid syntax")
        } else {
            stmt := stmts[0]
            switch node := stmt.(type) {
            case *ast.ExpressionStatement:
                obj := evaluator.Eval(node, env)
                if obj != evaluator.Py_None {
                    fmt.Println(evaluator.StringOf(obj))
                }
            case ast.Statement:
                evaluator.Exec(stmts, env)
            }
        }
    }
}

