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
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resultedObj := env.Get(target).(*object.NumberObject); *resultedObj != *expectedObj {
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
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
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
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*object.NumberObject); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    varname, *expectedObj, varname, *resultedObj)
            }
        }
    }
}

func TestIfElifElseStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.NumberObject
    }{
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a > b:\n" +
            "    res = 1\n" +
            "elif a < b:\n" +
            "    res = 2\n" +
            "else:\n" +
            "    res = 3\n",
            map[string]*object.NumberObject{
                "res": &object.NumberObject{Value: 2},
            },
        },
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a > b:\n" +
            "    res = 1\n" +
            "else:\n" +
            "    res = 2\n",
            map[string]*object.NumberObject{
                "res": &object.NumberObject{Value: 2},
            },
        },
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a < b:\n" +
            "    res = 1\n" +
            "else:\n" +
            "    res = 2\n",
            map[string]*object.NumberObject{
                "res": &object.NumberObject{Value: 1},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
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
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
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
    testRunProgram(input)
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
    env := testRunProgram(input)
    if obj := env.Get("b"); obj.(*object.NumberObject).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*object.NumberObject).Value)
    }
}

func TestEOFLineStatement(t *testing.T) {
    // should have no error
    input := ""+
    "a = 1 + 1\n" +
    "a + 1\n" +
    "b = a + 1"
    env := testRunProgram(input)
    if obj := env.Get("b"); obj.(*object.NumberObject).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*object.NumberObject).Value)
    }
}

func TestFunctionDefStatement(t *testing.T) {
    // should have no error
    input := ""+
    "def foo(a, b):\n" +
    "    c = a + b\n"
    env := testRunProgram(input)
    if obj := env.Get("foo"); obj == nil {
        t.Errorf("func 'foo' does not exists")
    }
}

func TestClassStatement(t *testing.T) {
    // should have no error
    input := ""+
    "class Foo:\n" +
    "   a = 1\n"
    env := testRunProgram(input)
    if obj := env.Get("Foo"); obj == nil {
        t.Errorf("class 'Foo' does not exists")
    }
}

func TestInstanceMethod(t *testing.T) {
    // should have no error
    input := `
class Foo:
    c = 3
    def __init__(self):
        self.a = 1
        self.b = 2
    def sum(self):
        return self.a + self.b + self.c

foo = Foo()
res = foo.sum()
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*object.NumberObject).Value != 6 {
        t.Errorf("instance method wrong: expected 6, got %v", obj.(*object.NumberObject).Value)
    }
}

func TestDotGetExpression(t *testing.T) {
    input := ""+
    "class Foo:\n" +
    "   a = 1\n" +
    "res = Foo.a"
    env := testRunProgram(input)

    if obj := env.Get("res"); obj.(*object.NumberObject).Value != 1 {
        t.Errorf("res should be %d, got %v", 1, obj.(*object.NumberObject).Value)
    }
}

func TestDotSetExpression(t *testing.T) {
    input := ""+
    "class Foo:\n" +
    "   a = 1\n" +
    "Foo.a = 2\n" +
    "res = Foo.a"
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*object.NumberObject).Value != 2 {
        t.Errorf("res should be %d, got %v", 2, obj.(*object.NumberObject).Value)
    }
}

func TestReturnStatement(t *testing.T) {
    // should have no error
    input := ""+
    "def foo(a, b):\n" +
    "    return a + b\n"
    env := testRunProgram(input)
    if obj := env.Get("foo"); obj == nil {
        t.Errorf("func 'foo' does not exists")
    }
}

func TestFunctionCallStatement(t *testing.T) {
    // should have no error
    input := ""+
    "def foo(a, b):\n" +
    "    return a + b\n" +
    "res = foo(1, 1)\n"
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*object.NumberObject).Value != 2 {
        t.Errorf("expected %v, got %v", 2, obj.(*object.NumberObject).Value)
    }
}

func TestStringAssignStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.StringObject
    }{
        {
            `val = "abc"`,
            map[string]*object.StringObject{
                "val": &object.StringObject{Value: "abc"},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resultedObj := env.Get(target).(*object.StringObject); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resultedObj)
            }
        }
    }
}

func TestStringPlusStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.StringObject
    }{
        {
            "a = 'abc'\n" +
            "b = 'def'\n" +
            "c = a + b\n",
            map[string]*object.StringObject{
                "c": &object.StringObject{Value: "abcdef"},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resultedObj := env.Get(target).(*object.StringObject); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resultedObj)
            }
        }
    }
}

func TestListStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.ListObject
    }{
        {
            "a = 'abc'\n" +
            "c = [1, a, 'd']\n",
            map[string]*object.ListObject{
                "c": &object.ListObject{
                        Items: []object.Object{
                                &object.NumberObject{Value: 1},
                                &object.StringObject{Value: "abc"},
                                &object.StringObject{Value: "d"},
                            },
                    },
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resObj := env.Get(target).(*object.ListObject); len(resObj.Items) != len(expectedObj.Items) {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resObj)
            }
        }
    }
}

func TestDictStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*object.DictObject
    }{
        {
            "a = 'abc'\n" +
            "d = {a: 'abc'}\n",
            map[string]*object.DictObject{
                "d": &object.DictObject{
                        Map: map[object.Object]object.Object{
                                &object.StringObject{Value: "abc"}: &object.StringObject{Value: "d"},
                            },
                    },
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resObj := env.Get(target).(*object.DictObject); len(resObj.Map) != len(expectedObj.Map) {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resObj)
            }
        }
    }
}

func TestZeroDivisionError(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            expr := r.(*object.ExceptionObject)
            if expr.Msg != "ZeroDivisionError: division by zero" {
                t.Errorf("expected 'ZeroDivisionError' got %v", expr.Msg)
            }
        }
    } ()

    testRunProgram("1 / 0")
    
}

func TestTypeError(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            expr := r.(*object.ExceptionObject)
            if expr.Msg != "TypeError: two different types" {
                t.Errorf("expected 'TypeError' got %v", expr.Msg)
            }
        }
    } ()

    testRunProgram("'a' + 1")
    
}

func testRunProgram(input string) *Environment{
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    env := NewEnvironment()
    Exec(stmts, env)
    return env
}

