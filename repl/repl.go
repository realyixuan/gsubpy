package repl

import (
    "os"
    "bufio"
    "fmt"

    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/evaluator"
    "gsubpy/object"
    "gsubpy/ast"
    "gsubpy/token"
)

func REPLRunning() {
    for {
        fmt.Print(">>> ")
        reader := bufio.NewReader(os.Stdin)

        line, err := reader.ReadString('\n')

        if err != nil {
            // if EOF
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
            //
        } else if len(stmts) > 1{
            panic("invalid syntax")
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
