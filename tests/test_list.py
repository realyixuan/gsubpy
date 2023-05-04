l = [1, 2,
3, ]

assert l[1] == 2


l[1] = 4
assert l[1] == 4

assert l == [1, 4, 3]

l.append(2)
assert l == [1, 4, 3, 2]

l.pop()
assert l == [1, 4, 3]

