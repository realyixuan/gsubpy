package test

import (
    "testing"
    "gsubpy/lexer"
    "gsubpy/parser"
    "gsubpy/evaluator"
)

func TestOneLineAssignStatement(t *testing.T) {
    input := `val = 10 + 20 * 10 / 2 - 50`
    env := testRunProgram(input)

    if resultedObj := env.GetFromString("val").(*evaluator.IntegerInst); resultedObj.Value != 60 {
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
    if resultObj := env.GetFromString("c").(*evaluator.IntegerInst); resultObj.Value != 3 {
        t.Errorf("expected %v, got %v", 3, resultObj.Value)
    }
}

func TestPlUSASSIGN(t *testing.T) {
    input := `
val = 1
val += 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("val").(*evaluator.IntegerInst); obj.Value != 3 {
        t.Errorf("expect 3, got %v", obj.Value)
    }
}

func TestMINUSASSIGN(t *testing.T) {
    input := `
val = 3
val -= 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("val").(*evaluator.IntegerInst); obj.Value != 1 {
        t.Errorf("expect 1, got %v", obj.Value)
    }
}

func TestMULASSIGN(t *testing.T) {
    input := `
val = 2
val *= 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("val").(*evaluator.IntegerInst); obj.Value != 4 {
        t.Errorf("expect 4, got %v", obj.Value)
    }
}

func TestDIVIDEASSIGN(t *testing.T) {
    input := `
val = 2
val /= 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("val").(*evaluator.IntegerInst); obj.Value != 1 {
        t.Errorf("expect 1, got %v", obj.Value)
    }
}

func TestEQ(t *testing.T) {
    input := `
res = 2 == 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", obj)
    }
}

func TestNOT(t *testing.T) {
    input := `
res = not 1 > 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", evaluator.StringOf(obj))
    }
}

func TestAND(t *testing.T) {
    input := `
res = 2 > 1 and 1 > 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj != evaluator.Py_False {
        t.Errorf("expect False, got %v", evaluator.StringOf(obj))
    }
}

func TestOR(t *testing.T) {
    input := `
res = 2 > 1 or 1 > 2
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", evaluator.StringOf(obj))
    }
}

func TestPAREN(t *testing.T) {
    input := `
res = (2 < 1) or 2 > 1
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", obj)
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
    if resultObj := env.GetFromString("res").(*evaluator.IntegerInst); resultObj.Value != 10 {
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
    if resultObj := env.GetFromString("res").(*evaluator.IntegerInst); resultObj.Value != 1 {
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
    if resultObj := env.GetFromString("total").(*evaluator.IntegerInst); resultObj.Value != 45 {
        t.Errorf("expected %v, got %v", 45, resultObj.Value)
    }
}

func TestExpressionStatement(t *testing.T) {
    input := `
a = 1 + 1
1 + 1
a + 1
`
    testRunProgram(input)
}

func TestBlankLineStatement(t *testing.T) {
    input := ""+
    "a = 1 + 1\n" +
    "     \n" +
    "1 + 1\n" +
    "\n"      +
    "     \n" +
    "a + 1\n" +
    "b = a + 1\n"
    env := testRunProgram(input)
    if obj := env.GetFromString("b"); obj.(*evaluator.IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestEOFLineStatement(t *testing.T) {
    input := `
a = 1 + 1
a + 1
b = a + 1
`
    env := testRunProgram(input)
    if obj := env.GetFromString("b"); obj.(*evaluator.IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestFunctionDefStatement(t *testing.T) {
    input := `
def foo(a, b):
    c = a + b
`
    env := testRunProgram(input)
    if obj := env.GetFromString("foo"); obj == nil {
        t.Errorf("func 'foo' does not exists")
    }
}

func TestClassStatement(t *testing.T) {
    input := `
class Foo:
   a = 1
`
    env := testRunProgram(input)
    if obj := env.GetFromString("Foo"); obj == nil {
        t.Errorf("class 'Foo' does not exists")
    }
}

func TestInstanceMethod(t *testing.T) {
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
    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 6 {
        t.Errorf("instance method wrong: expected 6, got %v", obj.(*evaluator.IntegerInst).Value)
    }
}

func TestClassInheritanceForClass(t *testing.T) {
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
    if obj := env.GetFromString("res1"); obj.(*evaluator.IntegerInst).Value != 10 {
        t.Errorf("instance method wrong: expected 10, got %v", obj.(*evaluator.IntegerInst).Value)
    }
    if obj := env.GetFromString("res2"); obj.(*evaluator.IntegerInst).Value != 10 {
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
        Base.__init__(self, a)

foo = Foo(1, 2)
res1 = foo.a
res2 = foo.b
    `
    env := testRunProgram(input)
    if obj := env.GetFromString("res1"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*evaluator.IntegerInst).Value)
    }
    if obj := env.GetFromString("res2"); obj.(*evaluator.IntegerInst).Value != 2 {
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
    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 1 {
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
    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 1 {
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
    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("instance attr wrong: expected 1, got %v", obj.(*evaluator.IntegerInst).Value)
    }
}

func TestTypeReturnType(t *testing.T) {
    input := `
s = "hello world"
res = type(s)
    `
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj != evaluator.Py_str {
        t.Errorf("type() wrong: expected Py_str %v, got %v", evaluator.Py_str.Id(), obj.Id())
    }
}

func TestDotGetExpression(t *testing.T) {
    input := ""+
    "class Foo:\n" +
    "   a = 1\n" +
    "res = Foo.a"
    env := testRunProgram(input)

    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 1 {
        t.Errorf("res should be %d, got %v", 1, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestDotSetExpression(t *testing.T) {
    input := `
class Foo:
   a = 1
Foo.a = 2
res = Foo.a
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 2 {
        t.Errorf("res should be %d, got %v", 2, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestReturnStatement(t *testing.T) {
    input := `
def foo():
    if 1 > 0:
        return 1
    return 0

res = foo()
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res").(*evaluator.IntegerInst); obj.Value != 1 {
        t.Errorf("expect 1, got %v", obj.Value)
    }
}

func TestFunctionCallStatement(t *testing.T) {
    input := ""+
    "def foo(a, b):\n" +
    "    return a + b\n" +
    "res = foo(1, 1)\n"
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 2 {
        t.Errorf("expected %v, got %v", 2, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestStringAssignStatement(t *testing.T) {
    input := `
val = "abc"
    `
    env := testRunProgram(input)
    if resultObj := env.GetFromString("val").(*evaluator.StringInst); resultObj.Value != "abc" {
        t.Errorf("expected 'abc', got %v", resultObj.Value)
    }
}

func TestStringPlusStatement(t *testing.T) {
    input := `
s = 'abc' + 'def'
    `
    env := testRunProgram(input)
    if resultObj := env.GetFromString("s").(*evaluator.StringInst); resultObj.Value != "abcdef" {
        t.Errorf("expected abcdef, got %v", resultObj.Value)
    }
}

func TestZeroDivisionError(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            if e, ok := r.(*evaluator.ExceptionInst); ok == false {
                t.Errorf("expected 'ZeroDivisionError' got %v", e)
            }
        }
    } ()

    testRunProgram("1 / 0")
    
}

func TestTypeError(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            if e, ok := r.(*evaluator.ExceptionInst); ok == false {
                t.Errorf("expected 'TypeError' got %v", e)
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
    if obj := env.GetFromString("res").(*evaluator.IntegerInst); obj.Value != 3 {
        t.Errorf("expect 3, got %v", obj.Value)
    }
}

func TestBuiltinLen(t *testing.T) {
    input := `
lt = [1, 2, 3]
res = len(lt)
    `
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj.(*evaluator.IntegerInst).Value != 3 {
        t.Errorf("expected %v, got %v", 3, obj.(*evaluator.IntegerInst).Value)
    }
}

func TestStrClass(t *testing.T) {
    input := `
res = str() + str(1) + str("2")
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res").(*evaluator.StringInst); obj.Value != "12" {
        t.Errorf("expect '12', got %v", obj.Value)
    }
}

func TesthashFunction(t *testing.T) {
    input := `
res = hash(".")
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res").(*evaluator.IntegerInst); obj.Value != -2750485470804815302 {
        t.Errorf("expect -2750485470804815302, got %v", obj.Value)
    }
}

func TestboolFunction(t *testing.T) {
    input := `
res = bool(1)
`
    env := testRunProgram(input)
    if obj := env.GetFromString("res"); obj != evaluator.Py_True {
        t.Errorf("expect True, got %v", obj)
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

