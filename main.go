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
    // "fmt"
)

func main() {
    input := `
a = 1
b = 2
c = a + b
`
    brutalforce(input)

}


func brutalforce(input string) {
    /*
    Here I also wanna feel the difference in difficulty between 
    doing it all at once and splitting them into more pieces of smaller function

    */

    env := map[string]int{}
    var curvar string
    var curval int
    var location string
    var exprval int
    var operator string = "+"

    for _, c := range input {
        // There are still a lot of mess 
        // So I decide to keep direction right, and simplify problems
        // continuely simplify the rules:
        //      single-char variable and single-digit number of value

        if c == '=' {
            location = "right"
            continue
        }

        if '0' <= c && c <= '9' {
            curval = int(c - '0')
        }

        if c == ' ' {
            continue
        }

        if c == '\n' {
            if exprval > 0 {
                env[curvar] = exprval
            } else {
                env[curvar] = curval
            }

            location = "left"
        }

        if 'a' <= c && c <= 'z' {
            if location == "left" {
                curvar = string(c)
            } else if location == "right" {
                if operator == "+" {
                    exprval += env[string(c)]
                }
            }
        }

        if c == '+' {
            operator = "+"
        }

    }

    println(env["c"])

}

