a =    \
    1

assert a == 1

class Foo:
    def __iter__(self):
        return self
    def __next__(self):
        raise StopIteration


for i in Foo():
    'pass'


assert not 1 > 2 is True

