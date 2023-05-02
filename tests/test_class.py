class Foo(object):
    def __init__(self, a):
        self.a = a

assert Foo(1).a == 1


class Foo:                                                                                     
    c = 3
    def __init__(self):
        self.a = 1
        self.b = 2
    def sum(self):
        return self.a + self.b + self.c

foo = Foo()

assert foo.sum() == 6


class Base:
    x = 10
class Foo(Base):
    factor = 30

foo = Foo()

assert foo.x == 10


class Base:                                                                                    
    def __init__(self, a):
        self.a = a

class Foo(Base):
    def __init__(self, a, b):
        self.b = b
        Base.__init__(self, a)

foo = Foo(1, 2)

assert foo.a == 1


class Foo:
    def __new__(cls, a):
        return object.__new__(cls)
        
    def __init__(self, a):
        self.a = a

assert Foo(1).a == 1


class Foo:
   a = 1                                                                                       
Foo.a = 2

assert Foo.a == 2

