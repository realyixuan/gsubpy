package evaluator

import (
    "testing"
    "gsubpy/lexer"
    "gsubpy/parser"
)

func TestBuiltinLen(t *testing.T) {
    // should have no error
    input := `
lt = [1, 2, 3]
res = len(lt)
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*IntegerInst).Value)
    }
}

func TestOneLineAssignStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*IntegerInst
    }{
        {
            `val = 10 + 20 * 10 / 2 - 50`,
            map[string]*IntegerInst{
                "val": &IntegerInst{Value: 60},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resultedObj := env.Get(target).(*IntegerInst); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resultedObj)
            }
        }
    }
}

func TestMultiLineAssignStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*IntegerInst
    }{
        {
            "vala = 1\n" +
            "valb = 2\n",
            map[string]*IntegerInst{
                "vala": &IntegerInst{Value: 1},
                "valb": &IntegerInst{Value: 2},
            },
        },
        {
            "a = 1\n" +
            "b = 2\n" +
            "c = a + b\n",
            map[string]*IntegerInst{
                "a": &IntegerInst{Value: 1},
                "b": &IntegerInst{Value: 2},
                "c": &IntegerInst{Value: 3},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*IntegerInst); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    varname, *expectedObj, varname, *resultedObj)
            }
        }
    }
}

func TestIfStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*IntegerInst
    }{
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a > b:\n" +
            "    a = a * 10\n" +
            "    b = b * 10\n" +
            "a = a + 10\n",
            map[string]*IntegerInst{
                "a": &IntegerInst{Value: 20},
                "b": &IntegerInst{Value: 20},
            },
        },
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a < b:\n" +
            "    a = a * 10\n" +
            "    b = b * 10\n" +
            "a = a + 10\n",
            map[string]*IntegerInst{
                "a": &IntegerInst{Value: 110},
                "b": &IntegerInst{Value: 200},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*IntegerInst); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    varname, *expectedObj, varname, *resultedObj)
            }
        }
    }
}

func TestIfElifElseStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*IntegerInst
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
            map[string]*IntegerInst{
                "res": &IntegerInst{Value: 2},
            },
        },
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a > b:\n" +
            "    res = 1\n" +
            "else:\n" +
            "    res = 2\n",
            map[string]*IntegerInst{
                "res": &IntegerInst{Value: 2},
            },
        },
        {
            "a = 10\n" +
            "b = 20\n" +
            "if a < b:\n" +
            "    res = 1\n" +
            "else:\n" +
            "    res = 2\n",
            map[string]*IntegerInst{
                "res": &IntegerInst{Value: 1},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*IntegerInst); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    varname, *expectedObj, varname, *resultedObj)
            }
        }
    }
}

func TestWhileStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*IntegerInst
    }{
        {
            "i = 0\n" +
            "total = 0\n" +
            "while i < 10:\n" +
            "    total = total + i\n" +
            "    i = i + 1\n",
            map[string]*IntegerInst{
                "total": &IntegerInst{Value: 45},
            },
        },
        {
            "i = 10\n" +
            "total = 0\n" +
            "while i > 10:\n" +
            "    total = total + i\n" +
            "    i = i + 1\n",
            map[string]*IntegerInst{
                "total": &IntegerInst{Value: 0},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for varname, expectedObj := range testCase.expected {
            res := env.Get(varname)
            if res == nil {
                t.Errorf("no variable %v", varname)
            } else if resultedObj := res.(*IntegerInst); *resultedObj != *expectedObj {
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
    if obj := env.Get("b"); obj.(*IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*IntegerInst).Value)
    }
}

func TestEOFLineStatement(t *testing.T) {
    // should have no error
    input := ""+
    "a = 1 + 1\n" +
    "a + 1\n" +
    "b = a + 1"
    env := testRunProgram(input)
    if obj := env.Get("b"); obj.(*IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*IntegerInst).Value)
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
    if obj := env.Get("res"); obj.(*IntegerInst).Value != 6 {
        t.Errorf("instance method wrong: expected 6, got %v", obj.(*IntegerInst).Value)
    }
}

func TestClassInheritanceForClass(t *testing.T) {
    // should have no error
    input := `
class Base:
    x = 10
class Foo(Base):
    factor = 30

foo = Foo()
res1 = Foo.x 
res2 = foo.x
    `
    env := testRunProgram(input)
    if obj := env.Get("res1"); obj.(*IntegerInst).Value != 10 {
        t.Errorf("instance method wrong: expected 10, got %v", obj.(*IntegerInst).Value)
    }
    if obj := env.Get("res2"); obj.(*IntegerInst).Value != 10 {
        t.Errorf("instance method wrong: expected 10, got %v", obj.(*IntegerInst).Value)
    }
}

func TestClassInheritanceForInstance(t *testing.T) {
    input := `
class Base:
    def __init__(self, a):
        self.a = a

class Foo(Base):
    def __init__(self, a, b):
        self.b = b
        super().__init__(a)

foo = Foo(1, 2)
res1 = foo.a
res2 = foo.b
    `
    env := testRunProgram(input)
    if obj := env.Get("res1"); obj.(*IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*IntegerInst).Value)
    }
    if obj := env.Get("res2"); obj.(*IntegerInst).Value != 2 {
        t.Errorf("instance attr wrong: expected 2, got %v", obj.(*IntegerInst).Value)
    }
}

func TestClass__new__(t *testing.T) {
    input := `
class Foo:
    def __new__(cls):
        return object.__new__(cls)
        
    def __init__(self, a):
        self.a = a

res = Foo(1).a
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*IntegerInst).Value)
    }
}

func TestInheritanceWithoutObject(t *testing.T) {
    input := `
class Foo:
    def __new__(cls):
        return object.__new__(cls)
        
    def __init__(self, a):
        self.a = a

res = Foo(1).a
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*IntegerInst).Value)
    }
}

func TestInheritanceWithObject(t *testing.T) {
    input := `
class Foo(object):
    def __init__(self, a):
        self.a = a

res = Foo(1).a
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*IntegerInst).Value)
    }
}

func TestTypeReturnType(t *testing.T) {
    input := `
s = "hello world"
res = type(s)
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.Type() != TYPE {
        t.Errorf("type() wrong: expected %v, got %v", TYPE, obj.Type())
    }
}

func TestDotGetExpression(t *testing.T) {
    input := ""+
    "class Foo:\n" +
    "   a = 1\n" +
    "res = Foo.a"
    env := testRunProgram(input)

    if obj := env.Get("res"); obj.(*IntegerInst).Value != 1 {
        t.Errorf("res should be %d, got %v", 1, obj.(*IntegerInst).Value)
    }
}

func TestDotSetExpression(t *testing.T) {
    input := ""+
    "class Foo:\n" +
    "   a = 1\n" +
    "Foo.a = 2\n" +
    "res = Foo.a"
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*IntegerInst).Value != 2 {
        t.Errorf("res should be %d, got %v", 2, obj.(*IntegerInst).Value)
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
    if obj := env.Get("res"); obj.(*IntegerInst).Value != 2 {
        t.Errorf("expected %v, got %v", 2, obj.(*IntegerInst).Value)
    }
}

func TestStringAssignStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*PyStrInst
    }{
        {
            `val = "abc"`,
            map[string]*PyStrInst{
                "val": &PyStrInst{Value: "abc"},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resultedObj := env.Get(target).(*PyStrInst); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resultedObj)
            }
        }
    }
}

func TestStringPlusStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*PyStrInst
    }{
        {
            "a = 'abc'\n" +
            "b = 'def'\n" +
            "c = a + b\n",
            map[string]*PyStrInst{
                "c": &PyStrInst{Value: "abcdef"},
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resultedObj := env.Get(target).(*PyStrInst); *resultedObj != *expectedObj {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resultedObj)
            }
        }
    }
}

func TestListStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*ListInst
    }{
        {
            "a = 'abc'\n" +
            "c = [1, a, 'd']\n",
            map[string]*ListInst{
                "c": &ListInst{
                        Items: []Object{
                                &IntegerInst{Value: 1},
                                &PyStrInst{Value: "abc"},
                                &PyStrInst{Value: "d"},
                            },
                    },
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resObj := env.Get(target).(*ListInst); len(resObj.Items) != len(expectedObj.Items) {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resObj)
            }
        }
    }
}

func TestDictStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*DictInst
    }{
        {
            "a = 'abc'\n" +
            "d = {a: 'abc'}\n",
            map[string]*DictInst{
                "d": &DictInst{
                        Map: map[Object]Object{
                                &PyStrInst{Value: "abc"}: &PyStrInst{Value: "d"},
                            },
                    },
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resObj := env.Get(target).(*DictInst); len(resObj.Map) != len(expectedObj.Map) {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resObj)
            }
        }
    }
}

func TestZeroDivisionError(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            expr := r.(*ExceptionInst)
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
            expr := r.(*ExceptionInst)
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

