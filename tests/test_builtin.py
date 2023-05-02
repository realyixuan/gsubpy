assert type("hello") == str

assert int() + int(1) + int('2') == 3

assert len('abc') == 3

assert str() + str(1) + str('2') == '12'

assert hash('.') == hash('.')

assert bool(1) == True


it = iter([1, 2])

assert next(it) == 1
assert next(it) == 2


it = iter(range(2))

assert next(it) == 0
assert next(it) == 1

