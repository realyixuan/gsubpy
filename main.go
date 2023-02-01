/*
Here I wonder why we should write in some order of lexer, parser, and ast ...
I have been wanting to write code in my own thinking way for a long time.
ALSO, this can be as a try to conquer hard unknown problem.
So I still want to try it in a irregular way.
Then, I will write my own first interpreter as if I never meet the right solution before
and totally have no idea how to do.
*/

/*

Actually, I totally have no idea what to do to get this done.
But I think, in this situation, trying to conclude a generic way is almost impossible.
Conversely, brutal force, which now could be the best way, can make me,
in the processing, figure clearly out what the problems truly are.
Yeah, I don't know the viable solution, nor the problems. 

My first thought is to write the rough code which can process simple python code.

1. Maybe I can first write the code processing simplest code like this:
    a = 1
    b = 1
    c = a + b

2. Then the code processing a little bit harder code like this:
    a = 1
    b = 1
    def foo(a, b):
        return a + b
    c = foo(a, b)

After that, I try to and may be possibly come up with a relatively generic way.

eventually, I should be able to interpreter the code: 

    def return_even_or_zero(n):
        if i % 2 == 0:
            return i
        return 0
    total = 0
    i = 0
    while i < 10:
        val = return_even_or_zero(i)
        total = total + val
        i = i + 1
    print(total)

for brutalforce, the order of implementation is: simplest > multi-char-variable > multi-digit-number > minus-operator > single-condition-if > while loop > expr including number > function > nested indent > if and while



*/

package main

import (
    "fmt"
    "strings"
    "strconv"
)

func main() {
    test_normalapproach()
}

func test_brutalforce() {
    input1 := `
a = 1
b = 2
c = a + b
`
    if env := brutalforce(input1); env["c"] != 3 {
        fmt.Println("simplest impl failed:", env["c"], env)
    }

    input2 := `
one = 1
two = 2
sum = one + two
`
    if env := brutalforce(input2); env["sum"] != 3 {
        fmt.Println("multi-char variable failed:", env["sum"], env)
    }

    input3 := `
twelve = 12
twentyone = 21
sum = twelve + twentyone
`
    if env := brutalforce(input3); env["sum"] != 33 {
        fmt.Println("multi-digit number failed:", env["sum"], env)
    }

    input4 := `
twelve = 12
twentyone = 21
sum = twentyone - twelve
`
    if env := brutalforce(input4); env["sum"] != 9 {
        fmt.Println("minus-operator failed:", env["sum"], env)
    }

    input5 := `
twelve = 12
twentyone = 21
sum = -twentyone - twelve
`
    if env := brutalforce(input5); env["sum"] != -33 {
        fmt.Println("pre-minus-operator failed:", env["sum"], env)
    }


    input6 := `
twelve = 12
twentyone = 21
sum = -twentyone - twelve
`
    if env := brutalforce(input6); env["sum"] != -33 {
        fmt.Println("pre-minus-operator failed:", env["sum"], env)
    }

    input7 := `
twelve = 12
twentyone = 21
difference = twentyone - twelve
if difference < 0:
    difference = 0
`
    if env := brutalforce(input7); env["difference"] != 9 {
        fmt.Println("single-if-condition-1 failed:", env["difference"], env)
    }

    input8 := `
twelve = 12
twentyone = 21
difference = twelve - twentyone
if difference < 0:
    difference = 0
`
    if env := brutalforce(input8); env["difference"] != 0 {
        fmt.Println("single-if-condition-2 failed:", env["difference"], env)
    }

    input9 := `
total = 0
i = 0
k = 1
while i < 20: 
    total = total + i
    i = i + k
`
    // instead of previous implementations, here I can't iterate
    // each char once, I need to return to origin of while-block, 
    // and still collecting while computing, or collecting and computing
    // apparently, latter is easier, but former is what I wanna try once,
    // and figure out why it is so. 

    // the obvious way is: record the origin the while-block, of course extra variable
    // is necessary, then set up the loop.

    if env := brutalforce(input9); env["total"] != 190 {
        fmt.Println("single-condition-while-loop failed:", env["total"], env)
    }

    input10 := `
sum = 12 + 21
`
    if env := brutalforce(input10); env["sum"] != 33 {
        fmt.Println("number-in-expression-1 failed:", env["sum"], env)
    }

    input11 := `
twelve = 12
sum = twelve + 21
`
    if env := brutalforce(input11); env["sum"] != 33 {
        fmt.Println("number-in-expression-2 failed:", env["sum"], env)
    }

    input12 := `
def add(a, b):
    return a + b

ten = 10
nine = 9

`
    // just initialized function
    brutalforce(input12)

}


func brutalforce(input string) map[string]int {
    /*
    Here I also wanna feel the difference in difficulty between 
    doing it all at once and splitting them into more pieces of smaller function

    */

    is_alphabet := func(c byte) bool { return 'a' <= c && c <= 'z' }
    is_digit := func(c byte) bool { return '0' <= c && c <= '9' }

    env := map[string]int{}
    var word string
    var number string
    var tovar string
    var curval string
    var location string
    var exprval int
    var operator string = ""
    var fromvar string
    var keyword string
    var lcond int
    var rcond int
    var ifcondresult bool
    var curindent int
    var columncount int
    var whileorigin int = -1
    var whilecondresult bool
    var curstate string

    type function struct {
        name string
        args []string
        begin int
        end int
        env map[string]int
    }
    // var funcs [string]function
    var curfunc function


    for i := 0; i < len(input); i++ {
        c := input[i]
        columncount += 1

        if is_alphabet(c) {
            word += string(c)
        } else if is_digit(c) {
            number += string(c)
        } else {
            if keyword != "" {
                if keyword == "if" {
                    // TODO: if set condition by operator? try it later
                    if word != "" {
                        if operator == "" {
                            v, err := env[word]
                            if err == true {
                                lcond = v
                            } else {
                                fmt.Println("if condition err:", err, env, word)
                            }
                        } else {
                            v, err := env[word]
                            if err == true {
                                rcond = v
                            } else {
                                fmt.Println("if condition err:", err, env, word)
                            }
                        }
                    }

                    if number != "" {
                        if operator == "" {
                            v, err := strconv.Atoi(number)
                            if err == nil {
                                lcond = v
                            } else {
                                fmt.Println("while condition err: ", err, env, word)
                            }
                        } else {
                            v, err := strconv.Atoi(number)
                            if err == nil {
                                rcond = v
                            } else {
                                fmt.Println("while condition err: ", word, env, word)
                            }
                        }
                    }

                    if c == ':' {
                        if operator == ">" {
                            if lcond > rcond {
                                ifcondresult = true
                            } else {
                                ifcondresult = false
                            }
                        } else if operator == "<" {
                            if lcond < rcond {
                                ifcondresult = true
                            } else {
                                ifcondresult = false
                            }
                        }
                    }

                } else if keyword == "while" {
                    if word != "" {
                        if operator == "" {
                            v, err := env[word]
                            if err == true {
                                lcond = v
                            } else {
                                fmt.Println("err:", lcond)
                            }
                        } else {
                            v, err := env[word]
                            if err == true {
                                rcond = v
                            } else {
                                fmt.Println("err:", lcond)
                            }
                        }
                    }

                    if number != "" {
                        if operator == "" {
                            v, err := strconv.Atoi(number)
                            if err == nil {
                                lcond = v
                            } else {
                                fmt.Println("err: ", err)
                            }
                        } else {
                            v, err := strconv.Atoi(number)
                            if err == nil {
                                rcond = v
                            } else {
                                fmt.Println("err: ", err)
                            }
                        }
                    }

                    if c == ':' {
                        if operator == ">" {
                            if lcond > rcond {
                                whilecondresult = true
                            } else {
                                whilecondresult = false
                            }
                        } else if operator == "<" {
                            if lcond < rcond {
                                whilecondresult = true
                            } else {
                                whilecondresult = false
                            }
                        }
                    }
                } else if keyword == "def" {
                    if curstate == "funcname" {
                        curfunc.name = word
                    } else if curstate == "funcargs" {
                        if word != "" {
                            curfunc.args = append(curfunc.args, word)
                        }
                    }

                    if c == '(' {
                        curstate = "funcargs"
                    } else if c == ')' {
                        curstate = "funcblock"
                    } else if c == ',' {
                    } else if c == ':' {
                    }
                }
            }


            {
                if word == "if" {
                    // Here would add more variable to record current situation
                    // Also, here should a variable to record the indents level
                    keyword = "if"
                } else if word == "while" {
                    keyword = "while"
                    whileorigin = i - len("while")
                } else if word == "def" {
                    keyword = "def"
                    curstate = "funcname"
                }
            }


            {
                if c == ' ' {
                    if columncount - curindent == 1 {
                        curindent += 1
                    }
                }

                if c == '=' {
                    location = "right"
                }

                if word != "" {
                    if location == "left" {
                        tovar = word
                    } else if location == "right" {
                        fromvar = word
                    }
                }

                if location == "right" {
                    if fromvar != "" {
                        if operator == "+" || operator == "" {
                            exprval += env[fromvar]
                        } else if operator == "-" {
                            exprval -= env[fromvar]
                        }
                        fromvar = ""
                    } else if number != "" {
                        v, err := strconv.Atoi(number)
                        if err != nil {
                            fmt.Println("strconv.Atoi(number) err", v, err)
                        } else {
                            if operator == "+" || operator == "" {
                                exprval += v
                            } else if operator == "-" {
                                exprval -= v
                            }
                        }
                    }
                }

                if c == '\n' {
                    if keyword == "ifblock" && ifcondresult == true ||
                        keyword == "" ||
                        keyword == "whileblock" && whilecondresult == true {
                        if tovar != "" {
                            if curval != "" {
                                v, err := strconv.Atoi(curval)
                                if err == nil {
                                    env[tovar] = v
                                } else {
                                    fmt.Println("there is a error:", err)
                                }
                            } else {
                                env[tovar] = exprval
                            }
                        }
                    }

                    if keyword == "if" {
                        keyword = "ifblock"
                    } else if keyword == "ifblock" {
                        if i+1+curindent >= len(input) || i >= len(input) {
                            keyword = ""
                        } else if strings.Repeat(" ", curindent) != input[i+1:i+1+curindent] {
                            keyword = ""
                        }
                    } else if keyword == "while" {
                        keyword = "whileblock"
                    } else if keyword == "whileblock" {
                        if whilecondresult == true {
                            if i+1+curindent >= len(input) || i >= len(input) {
                                i = whileorigin - 1
                                keyword = ""
                            } else if strings.Repeat(" ", curindent) != input[i+1:i+1+curindent] {
                                i = whileorigin - 1
                                keyword = ""
                            }
                        } else {
                            if i+1+curindent >= len(input) || i >= len(input) {
                                keyword = ""
                                whileorigin = -1
                            } else if strings.Repeat(" ", curindent) != input[i+1:i+1+curindent] {
                                keyword = ""
                                whileorigin = -1
                            }
                        }
                    } else if keyword == "def" {
                        keyword = "defblock"
                        curfunc.begin = i + 1
                    } else if keyword == "defblock" {
                        if i+1+curindent >= len(input) || i >= len(input) {
                            keyword = ""
                            curfunc.end = i
                        } else if strings.Repeat(" ", curindent) != input[i+1:i+1+curindent] {
                            keyword = ""
                            curfunc.end = i
                        }
                    }

                    location = "left"

                    tovar = ""
                    fromvar = ""
                    curval = ""
                    operator = ""
                    exprval = 0
                    columncount = 0

                    curindent = 0
                }

                if c == '+' || c == '-' || c == '<' || c == '>' {
                    operator = string(c)
                }

            }
            word = ""
            number = ""
        }

    }

    return env
}

func test_normalapproach() {
    input := `
abc = 10+ 1111
`
    normalapproach(input)

//     input1 := `
// a = 1
// b = 2
// `
//     if env := normalapproach(input1); env["c"] != 3 {
//         fmt.Println("simplest impl failed:", env["c"], env)
//     }
}

func normalapproach(input string) map[string]int {
    /*
    In this approach, separate the process into difference parts.
    some difference parts from above brutalforce:
        - define sematics, statement and expression
        - separate lexer and eval
    */

    tokens := lexer(input)

    ret := parser(tokens)
    // fmt.Println(reflect.TypeOf(ret[0]))

    fmt.Println(ret)

    return map[string]int{"Hello": 1}
}

type Statement struct {
    tokens []string
}

type Assignment struct {
    variable string
    expression Statement
}

func parser(tokens []string) []interface{} {
    var stmt []interface{}
    statements := get_statements(tokens)
    for _, statement := range statements {
        if idx := contain(statement.tokens, "="); idx != -1 {
            var assignment = Assignment{
                variable: statement.tokens[idx-1],
                expression: Statement{statement.tokens[idx+1:]},
            }
            stmt = append(stmt, assignment)
        }
    }

    return stmt
}

func get_statements(tokens []string) []Statement {
    var statements []Statement

    var idx int = -1
    for i := 0; i < len(tokens); i++ {
        if tokens[i] == "END" {
            if i-idx > 1 {
                statements = append(statements, Statement{tokens[idx+1:i]})
            }
            idx = i
        }
    }

    return statements
}

func lexer(s string) []string {
    /*
    converting code text into tokens, and
        - linefeed becoming string `END`
    */
    is_digit := func(c rune) bool {return '0' <= c && c <= '9'}
    is_alphabet := func(c rune) bool {return 'a' <= c && c <= 'z'}

    var tokens []string
    var word string
    for _, c := range s {
        if is_digit(c) {
            word += string(c)
        } else if is_alphabet(c) {
            word += string(c)
        } else {
            if word != "" {
                tokens = append(tokens, word)
            }

            if c == ' ' {
                // pass
            } else if c == '\n'{
                tokens = append(tokens, "END")
            } else {
                tokens = append(tokens, string(c))
            }

            word = ""
        }
    }

    if word != "" {
        tokens = append(tokens, word)
    }

    return tokens
}

func contain(arr []string, target string) int {
    for i, s := range arr {
        if s == target {
            return i
        }
    }

    return -1
}

