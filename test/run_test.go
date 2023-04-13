package test

import (
    "testing"
    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/evaluator"
)

func TestBuiltinLen(t *testing.T) {
    // should have no error
    input := `
lt = [1, 2, 3]
res = len(lt)
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestOneLineAssignStatement(t *testing.T) {
    input := `val = 10 + 20 * 10 / 2 - 50`
    env := testRunProgram(input)

    if resultedObj := env.Get("val").(*evaluator.IntegerInst); resultedObj.Value != 60 {
        t.Errorf("expected %v, got %v", 60, resultedObj.Value)
    }
}

func TestMultiLineAssignStatement(t *testing.T) {
    input := `
a = 1
b = 2
c = a + b
    `

    env := testRunProgram(input)
    if resultObj := env.Get("c").(*evaluator.IntegerInst); resultObj.Value != 3 {
        t.Errorf("expected %v, got %v", 3, resultObj.Value)
    }
}

func TestPlUSASSIGN(t *testing.T) {
    input := `
val = 1
val += 2
`
    env := testRunProgram(input)
    if obj := env.Get("val").(*evaluator.IntegerInst); obj.Value != 3 {
        t.Errorf("expect 3, got %v", obj.Value)
    }
}

func TestMINUSASSIGN(t *testing.T) {
    input := `
val = 3
val -= 2
`
    env := testRunProgram(input)
    if obj := env.Get("val").(*evaluator.IntegerInst); obj.Value != 1 {
        t.Errorf("expect 1, got %v", obj.Value)
    }
}

func TestMULASSIGN(t *testing.T) {
    input := `
val = 2
val *= 2
`
    env := testRunProgram(input)
    if obj := env.Get("val").(*evaluator.IntegerInst); obj.Value != 4 {
        t.Errorf("expect 4, got %v", obj.Value)
    }
}

func TestDIVIDEASSIGN(t *testing.T) {
    input := `
val = 2
val /= 2
`
    env := testRunProgram(input)
    if obj := env.Get("val").(*evaluator.IntegerInst); obj.Value != 1 {
        t.Errorf("expect 1, got %v", obj.Value)
    }
}

func TestEQ(t *testing.T) {
    input := `
res = 2 == 2
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.BoolInst); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", obj.Py__str__())
    }
}

func TestNOT(t *testing.T) {
    input := `
res = not 1 > 2
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.BoolInst); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", obj.Py__str__())
    }
}

func TestAND(t *testing.T) {
    input := `
res = 2 > 1 and 1 > 2
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.BoolInst); obj != evaluator.Py_False {
        t.Errorf("expect False, got %v", obj.Py__str__())
    }
}

func TestOR(t *testing.T) {
    input := `
res = 2 > 1 or 1 > 2
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.BoolInst); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", obj.Py__str__())
    }
}

func TestPAREN(t *testing.T) {
    input := `
res = (not 2 > 1) or 2 > 1
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.BoolInst); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", obj.Py__str__())
    }
}

func TestIfStatement(t *testing.T) {
    input := `
res = 0
if 2 > 1:
    res = 10
if 2 < 1:
    res = res * 2
    `
    env := testRunProgram(input)
    if resultObj := env.Get("res").(*evaluator.IntegerInst); resultObj.Value != 10 {
        t.Errorf("expected %v, got %v", 10, resultObj.Value)
    }
}

func TestIfElifElseStatement(t *testing.T) {
    input := `
res = 0
if 2 > 1:
    res = 1
elif 2 < 1:
    res = 2
else:
    res = 3
    `

    env := testRunProgram(input)
    if resultObj := env.Get("res").(*evaluator.IntegerInst); resultObj.Value != 1 {
        t.Errorf("expected %v, got %v", 1, resultObj.Value)
    }
}

func TestWhileStatement(t *testing.T) {
    input := `
i = 0
total = 0
while i < 10:
    total = total + i
    i = i + 1
    `
    env := testRunProgram(input)
    if resultObj := env.Get("total").(*evaluator.IntegerInst); resultObj.Value != 45 {
        t.Errorf("expected %v, got %v", 45, resultObj.Value)
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
    if obj := env.Get("b"); obj.(*evaluator.IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestEOFLineStatement(t *testing.T) {
    // should have no error
    input := ""+
    "a = 1 + 1\n" +
    "a + 1\n" +
    "b = a + 1"
    env := testRunProgram(input)
    if obj := env.Get("b"); obj.(*evaluator.IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*evaluator.IntegerInst).Value)
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
    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 6 {
        t.Errorf("instance method wrong: expected 6, got %v", obj.(*evaluator.IntegerInst).Value)
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
    if obj := env.Get("res1"); obj.(*evaluator.IntegerInst).Value != 10 {
        t.Errorf("instance method wrong: expected 10, got %v", obj.(*evaluator.IntegerInst).Value)
    }
    if obj := env.Get("res2"); obj.(*evaluator.IntegerInst).Value != 10 {
        t.Errorf("instance method wrong: expected 10, got %v", obj.(*evaluator.IntegerInst).Value)
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
    if obj := env.Get("res1"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*evaluator.IntegerInst).Value)
    }
    if obj := env.Get("res2"); obj.(*evaluator.IntegerInst).Value != 2 {
        t.Errorf("instance attr wrong: expected 2, got %v", obj.(*evaluator.IntegerInst).Value)
    }
}

func TestClass__new__(t *testing.T) {
    input := `
class Foo:
    def __new__(cls, a):
        return object.__new__(cls)
        
    def __init__(self, a):
        self.a = a

res = Foo(1).a
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*evaluator.IntegerInst).Value)
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
    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*evaluator.IntegerInst).Value)
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
    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*evaluator.IntegerInst).Value)
    }
}

func TestTypeReturnType(t *testing.T) {
    input := `
s = "hello world"
res = type(s)
    `
    env := testRunProgram(input)
    if obj := env.Get("res"); obj != evaluator.Py_str {
        t.Errorf("type() wrong: expected %v, got %v", evaluator.Py_str.Py__str__(), obj.Py__str__())
    }
}

func TestDotGetExpression(t *testing.T) {
    input := ""+
    "class Foo:\n" +
    "   a = 1\n" +
    "res = Foo.a"
    env := testRunProgram(input)

    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("res should be %d, got %v", 1, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestDotSetExpression(t *testing.T) {
    input := ""+
    "class Foo:\n" +
    "   a = 1\n" +
    "Foo.a = 2\n" +
    "res = Foo.a"
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 2 {
        t.Errorf("res should be %d, got %v", 2, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestReturnStatement(t *testing.T) {
    // should have no error
    input := `
def foo():
    if 1 > 0:
        return 1
    return 0

res = foo()
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.IntegerInst); obj.Value != 1 {
        t.Errorf("expect 1, got %v", obj.Value)
    }
}

func TestFunctionCallStatement(t *testing.T) {
    // should have no error
    input := ""+
    "def foo(a, b):\n" +
    "    return a + b\n" +
    "res = foo(1, 1)\n"
    env := testRunProgram(input)
    if obj := env.Get("res"); obj.(*evaluator.IntegerInst).Value != 2 {
        t.Errorf("expected %v, got %v", 2, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestStringAssignStatement(t *testing.T) {
    input := `
val = "abc"
    `
    env := testRunProgram(input)
    if resultObj := env.Get("val").(*evaluator.PyStrInst); resultObj.Value != "abc" {
        t.Errorf("expected 'abc', got %v", resultObj.Value)
    }
}

func TestStringPlusStatement(t *testing.T) {
    input := `
s = 'abc' + 'def'
    `
    env := testRunProgram(input)
    if resultObj := env.Get("s").(*evaluator.PyStrInst); resultObj.Value != "abcdef" {
        t.Errorf("expected abcdef, got %v", resultObj.Value)
    }
}

func TestListStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*evaluator.ListInst
    }{
        {
            "a = 'abc'\n" +
            "c = [1, a, 'd']\n",
            map[string]*evaluator.ListInst{
                "c": &evaluator.ListInst{
                        Items: []evaluator.Object{
                                &evaluator.IntegerInst{Value: 1},
                                &evaluator.PyStrInst{Value: "abc"},
                                &evaluator.PyStrInst{Value: "d"},
                            },
                    },
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resObj := env.Get(target).(*evaluator.ListInst); len(resObj.Items) != len(expectedObj.Items) {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resObj)
            }
        }
    }
}

func TestDictStatement(t *testing.T) {
    testCases := []struct {
        input       string
        expected    map[string]*evaluator.DictInst
    }{
        {
            "a = 'abc'\n" +
            "d = {a: 'abc'}\n",
            map[string]*evaluator.DictInst{
                "d": &evaluator.DictInst{
                        Map: map[evaluator.PyStrInst]evaluator.Object{
                                evaluator.PyStrInst{Value: "abc"}: &evaluator.PyStrInst{Value: "d"},
                            },
                    },
            },
        },
    }

    for _, testCase := range testCases {
        env := testRunProgram(testCase.input)
        for target, expectedObj := range testCase.expected {
            if resObj := env.Get(target).(*evaluator.DictInst); len(resObj.Map) != len(expectedObj.Map) {
                t.Errorf("expected (%s=%v), got (%s=%v)",
                    target, *expectedObj, target, *resObj)
            }
        }
    }
}

func TestZeroDivisionError(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            expr := r.(*evaluator.ExceptionInst)
            if msg := expr.Py__str__().Value; msg != "ZeroDivisionError: division by zero" {
                t.Errorf("expected 'ZeroDivisionError' got %v", msg)
            }
        }
    } ()

    testRunProgram("1 / 0")
    
}

func TestTypeError(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            expr := r.(*evaluator.ExceptionInst)
            if msg := expr.Py__str__().Value; msg != "TypeError: two different types" {
                t.Errorf("expected 'TypeError' got %v", msg)
            }
        }
    } ()

    testRunProgram("'a' + 1")
    
}

func TestIntClass(t *testing.T) {
    input := `
res = int() + int(1) + int('2')
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.IntegerInst); obj.Value != 3 {
        t.Errorf("expect 3, got %v", obj.Value)
    }
}

func TestStrClass(t *testing.T) {
    input := `
res = str() + str(1) + str("2")
`
    env := testRunProgram(input)
    if obj := env.Get("res").(*evaluator.PyStrInst); obj.Value != "12" {
        t.Errorf("expect '12', got %v", obj.Value)
    }
}

func testRunProgram(input string) *evaluator.Environment{
    l := lexer.New(input)
    p := parser.New(l)
    stmts := p.Parsing()
    env := evaluator.NewEnvironment()
    evaluator.Exec(stmts, env)
    return env
}

