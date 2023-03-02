package evaluator

import (
    "testing"
    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/object"
)

func TestOneLineAssignStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.NumberObject
    }{
        {
            `val = 10 + 20 * 10 / 2 - 50`,
            map[string]*object.NumberObject{
                "val": &object.NumberObject{Value: 60},
            },
        },
    }

    for _, testCase := range testCases {
        l := lexer.New(testCase.input)
        p := parser.New(l)
        stmts := p.Parsing()
        Exec(stmts)
        for target, expectedObj := range testCase.expected {
            if resultedObj := env[target].(*object.NumberObject); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resultedObj)
            }
        }
    }
}

func TestMultiLineAssignStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.NumberObject
    }{
        {
            "vala = 1\n" +
            "valb = 2\n",
            map[string]*object.NumberObject{
                "vala": &object.NumberObject{Value: 1},
                "valb": &object.NumberObject{Value: 2},
            },
        },
        {
            "a = 1\n" +
            "b = 2\n" +
            "c = a + b\n",
            map[string]*object.NumberObject{
                "a": &object.NumberObject{Value: 1},
                "b": &object.NumberObject{Value: 2},
                "c": &object.NumberObject{Value: 3},
            },
        },
    }

    for _, testCase := range testCases {
        l := lexer.New(testCase.input)
        p := parser.New(l)
        stmts := p.Parsing()
        Exec(stmts)
        for varname, expectedObj := range testCase.expected {
            res, ok := env[varname]
            if !ok {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*object.NumberObject); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    varname, *expectedObj, varname, *resultedObj)
            }
        }
    }
}

func TestIfStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.NumberObject
    }{
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a > b:\n" +
            "    a = a * 10\n" +
            "    b = b * 10\n" +
            "a = a + 10\n",
            map[string]*object.NumberObject{
                "a": &object.NumberObject{Value: 20},
                "b": &object.NumberObject{Value: 20},
            },
        },
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a < b:\n" +
            "    a = a * 10\n" +
            "    b = b * 10\n" +
            "a = a + 10\n",
            map[string]*object.NumberObject{
                "a": &object.NumberObject{Value: 110},
                "b": &object.NumberObject{Value: 200},
            },
        },
    }

    for _, testCase := range testCases {
        l := lexer.New(testCase.input)
        p := parser.New(l)
        stmts := p.Parsing()
        Exec(stmts)
        for varname, expectedObj := range testCase.expected {
            res, ok := env[varname]
            if !ok {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*object.NumberObject); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    varname, *expectedObj, varname, *resultedObj)
            }
        }
    }
}

func TestWhileStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.NumberObject
    }{
        {
            "i = 0\n" +
            "total = 0\n" +
            "while i < 10:\n" +
            "    total = total + i\n" +
            "    i = i + 1\n",
            map[string]*object.NumberObject{
                "total": &object.NumberObject{Value: 45},
            },
        },
        {
            "i = 10\n" +
            "total = 0\n" +
            "while i > 10:\n" +
            "    total = total + i\n" +
            "    i = i + 1\n",
            map[string]*object.NumberObject{
                "total": &object.NumberObject{Value: 0},
            },
        },
    }

    for _, testCase := range testCases {
        l := lexer.New(testCase.input)
        p := parser.New(l)
        stmts := p.Parsing()
        Exec(stmts)
        for varname, expectedObj := range testCase.expected {
            res, ok := env[varname]
            if !ok {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*object.NumberObject); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    varname, *expectedObj, varname, *resultedObj)
            }
        }
    }
}

func TestExpressionStatement(t *testing.T) {
    // should have no error
    input := ""+
    "a = 1 + 1\n" + 
    "1 + 1\n" +
    "a + 1\n"
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    Exec(stmts)
}

func TestBlankLineStatement(t *testing.T) {
    // should have no error
    input := ""+
    "a = 1 + 1\n" +
    "     \n" +
    "1 + 1\n" +
    "\n"      +
    "     \n" +
    "a + 1\n" +
    "b = a + 1\n"
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    Exec(stmts)
    if obj, _ := env["b"]; obj.(*object.NumberObject).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*object.NumberObject).Value)
    }
}

func TestEOFLineStatement(t *testing.T) {
    // should have no error
    input := ""+
    "a = 1 + 1\n" +
    "a + 1\n" +
    "b = a + 1"
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    Exec(stmts)
    if obj, _ := env["b"]; obj.(*object.NumberObject).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*object.NumberObject).Value)
    }
}

func TestFunctionDefStatement(t *testing.T) {
    // should have no error
    input := ""+
    "def foo(a, b):\n" +
    "    c = a + b\n"
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    Exec(stmts)
    if obj, _ := env["foo"]; obj == nil {
        t.Errorf("func 'foo' does not exists")
    }
}

func TestReturnStatement(t *testing.T) {
    // should have no error
    input := ""+
    "def foo(a, b):\n" +
    "    return a + b\n"
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    Exec(stmts)
    if obj, _ := env["foo"]; obj == nil {
        t.Errorf("func 'foo' does not exists")
    }
}

func TestFunctionCallStatement(t *testing.T) {
    // should have no error
    input := ""+
    "def foo(a, b):\n" +
    "    return a + b\n" +
    "res = foo(1, 1)\n"
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    Exec(stmts)
    if obj, _ := env["res"]; obj.(*object.NumberObject).Value != 2 {
        t.Errorf("expected %v, got %v", 2, obj.(*object.NumberObject).Value)
    }
}

