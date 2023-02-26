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
            `val = 1`,
            map[string]*object.NumberObject{
                "val": &object.NumberObject{Value: 1},
            },
        },
        {
            `val = 1 + 1`,
            map[string]*object.NumberObject{
                "val": &object.NumberObject{Value: 2},
            },
        },
        {
            `val = 10 + 20 * 2`,
            map[string]*object.NumberObject{
                "val": &object.NumberObject{Value: 50},
            },
        },
    }

    for _, testCase := range testCases {
        l := lexer.New(testCase.input)
        p := parser.New(l)
        stmts := p.Parsing()
        run(stmts)
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
        run(stmts)
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

