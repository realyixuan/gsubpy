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

for brutalforce, the order of implementation is: simplest > multi-char-variable > multi-digit-number > minus-operator > single-condition-if > while loop



*/

package main

import (
    "fmt"
    "strconv"
)

func main() {
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
        fmt.Println("single-if-condition failed:", env["difference"], env)
    }

    input8 := `
twelve = 12
twentyone = 21
difference = twelve - twentyone
if difference < 0:
    different = 0
`
    if env := brutalforce(input8); env["res"] != 0 {
        fmt.Println("single-if-condition failed:", env["difference"], env)
    }
}


func brutalforce(input string) map[string]int {
    /*
    Here I also wanna feel the difference in difficulty between 
    doing it all at once and splitting them into more pieces of smaller function

    */

    is_alphabet := func(c rune) bool { return 'a' <= c && c <= 'z' }
    is_digit := func(c rune) bool { return '0' <= c && c <= '9' }

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


    for _, c := range input {
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
                        keyword = "ifblock"
                    }

                }
            }

            if word == "if" {
                // Here would add more variable to record current situation
                // Also, here should a variable to record the indents level
                keyword = "if"
            }

            {
                if c == ' ' {
                    if columncount - curindent == 1 {
                        curindent += 1
                    }
                }

                if columncount == 1 && curindent == 0 {
                    keyword = ""
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
                        if operator == "+" || operator == ""{
                            exprval += env[fromvar]
                        } else if operator == "-" {
                            exprval -= env[fromvar]
                        }
                        fromvar = ""
                    }
                }
                if number != "" {
                    curval = number
                }

                if c == '\n' {
                    if keyword == "ifblock" && ifcondresult == false {
                    } else {
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

                    location = "left"

                    tovar = ""
                    fromvar = ""
                    curval = ""
                    operator = ""
                    exprval = 0
                    curindent = 0
                    columncount = 0
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

