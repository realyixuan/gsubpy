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
}


func brutalforce(input string) map[string]int {
    /*
    Here I also wanna feel the difference in difficulty between 
    doing it all at once and splitting them into more pieces of smaller function

    then, implement multi-char variable and multi-digit number

    */

    env := map[string]int{}
    var curvar string
    var curval string
    var location string
    var exprval int
    var operator string = "+"
    var tcurvar string

    for _, c := range input {
        // if we use this approach of if-condition, then our way of thought is 
        // splitting variable and the others.
        // but, it's not that normal, maybe we split it by left and right 
        // is more suitable for us using if-condition.
        // later, I will try it.

        // And this patching-way can just resolve simple structure, 
        // can't go far, and I also wanna see.

        if 'a' <= c && c <= 'z' {
            if location == "left" {
                curvar += string(c)
            } else if location == "right" {
                tcurvar += string(c)
            }
        } else {
            if location == "right" {
                if tcurvar != "" {
                    if operator == "+" {
                        exprval += env[tcurvar]
                    } else if operator == "-" {
                        exprval -= env[tcurvar]
                    }
                    tcurvar = ""
                }
            }
            if c == '=' {
                location = "right"
                continue
            }

            if '0' <= c && c <= '9' {
                curval += string(c)
            }

            if c == ' ' {
                continue
            }

            if c == '\n' {
                if curvar != "" {
                    if curval != "" {
                        v, err := strconv.Atoi(curval)
                        if err == nil {
                            env[curvar] = v
                        } else {
                            fmt.Println("there is a error:", err)
                        }
                    } else {
                        env[curvar] = exprval
                    }
                }

                location = "left"

                curvar = ""
                tcurvar = ""
                curval = ""
                operator = "+"
                exprval = 0
            }

            if c == '+' || c == '-' {
                operator = string(c)
            }
        }

    }

    return env
}

