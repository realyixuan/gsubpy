package evaluator

import (
    "testing"
    // "reflect"

    "gsubpy/lexer"
    "gsubpy/object"
    "gsubpy/parser"
)


func TestSimpleAssignment1(t *testing.T) {
    input := `sum = 10`
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    env := NewEnvironment()
    Eval(program, env)

    if res := env.store["sum"]; res.(*object.Integer).Value != 10 {
        t.Errorf("assignment statement is wrong. got='%d', expected='%d'",
            res.(*object.Integer).Value, 10)
    }
}

func TestSimpleAssignment2(t *testing.T) {
    input := `sum = 1 + 10`
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    env := NewEnvironment()
    Eval(program, env)

    t.Log(env.store["sum"].(*object.Integer).Value)
    if res := env.store["sum"]; res.(*object.Integer).Value != 11 {
        t.Errorf("assignment statement is wrong. got='%d', expected='%d'",
            res.(*object.Integer).Value, 10)
    }
}

func TestSimpleAssignment3(t *testing.T) {
    input := `sum = 10 - 1`
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    env := NewEnvironment()
    Eval(program, env)

    t.Log(env.store["sum"].(*object.Integer).Value)
    if res := env.store["sum"]; res.(*object.Integer).Value != 9 {
        t.Errorf("assignment statement is wrong. got='%d', expected='%d'",
            res.(*object.Integer).Value, 10)
    }
}

// func TestEvalIntegerExpression(t *testing.T) {
//     tests := []struct {
//         input       string
//         expected    int64
//     } {
//         {"5", 5},
//         {"10", 10},
//     }
//     for _, tt := range tests {
//         evaluated := testEval(tt.input)
//         testIntegerObject(t, evaluated, tt.expected)
//     }
// }

// func testEval(input string) object.Object {
//     l := lexer.New(input)
//     p := parser.New(l)
//     program := p.ParseProgram()
// 
//     return Eval(program)
// }

// func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
//     result, ok := obj.(*object.Integer)
//     if !ok {
//         t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
//         return false
//     }
//     if result.Value != expected {
//         t.Errorf("object has wrong value. got=%d, want=%d",
//             result.Value, expected)
//         return false
//     }
//     return true
// }
// 
// func testAssignmentStatements(t *testing.T) {
//     tests := []struct {
//         input string
//         expected int64
//     }{
//         {"a = 10", 10},
//     }
// 
//     for _, tt := range tests {
//         testIntegerObject(t, testEval(tt.input), tt.expected)
//     }
// }

