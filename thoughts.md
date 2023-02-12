For realizing an interpreter of subset of Python

There, Python has a lot of syntax, implementing them all at once is nearly impossible. So, it's better to take a small step at a time.

First things first:

Removing, except for int, all complex data structure such as string, list ..., for data structure, just focus on int. Then, take concentration on statements of sequence, selection and iteration...

So, for now, here is focus:

- first step
    - data structre:
        - int
    - simplifying indent, say
        - fixed 4 spaces, and
        - one layer of indent
    - basic control flows:
        - sequence
        - selection
        - iteration
    - simple function

- second step
    - put in builtin function

- thrid step
    - put class in


1. ANYWAY, implement "a = 1" and "a = 1 + 1" first

2. sequence execution

3. if-statement

Consequently, after first step having been finished, the interpreter should at least be able to interpret this:

~~~python
def return_even_or_zero(n):
    if i % 2 == 0:
        return i
    return 0
total = 0
i = 0
while i < 10:
    val = return_even_or_zero(i)
    total = total + val
    i = i + 1
~~~

