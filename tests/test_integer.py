val = 10 + 20 * 10 / 2 - 50

assert val == 60


val = 10
val -= 5

assert val == 5


res = 2 == 2

assert res == True


a = 1
b = a

assert a is b


a = 1
b = 1

assert a is not b

assert (not 1 > 2) is True

assert (2 > 1 and 1 > 2) is False

