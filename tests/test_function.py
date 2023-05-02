def foo(a, b):
    return a + b

assert foo(1, 2) == 3


def foo():
    if 1 > 0:
        return 1
    return 0

assert foo() == 1

