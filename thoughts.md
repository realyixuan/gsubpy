# Some Thoughts For realizing an interpreter of subset of Python

There, Python has a lot of syntax, implementing them all at once is nearly impossible. So, it's better to take a small step at a time.

### Here's the rough plans:

Along realizing it, the development is parted into several phases, then solve all, step by step, phase by phase.

##### Phase 1: most-basic features

desired outcomes:

- being able to interpreter following code:

~~~python
def return_even_or_zero(n):
    if n % 2 == 0:
        return n
    return 0
total = 0
i = 0
while i < 10:
    val = return_even_or_zero(i)
    total = total + val
    i = i + 1
~~~

First things first:

simplifying problem: take concentration on something foundamental

- data type
    - int
- features
    - assignment statement
    - arithmic of plus, multiply
    - selection (i.e. if-else)
    - iteration (i.e. while)
    - and simple function

the key point of this phase is just getting things done. Keep general direction right, other than that I shouldn't take care of trivial details, as many things havne't got very clear, at least I think so.
    
steps:

step1: `assignment identifier = <expression>`

step2: `if statement`

step3: `while-loop statement`

step4: `function-definition and function-call`


##### Phase 2: enhancement

desired outcomes:

- built function putted in
- more builtin data types
    - str
    - list
    - dict

...

##### Phase 3: superior

desired outcomes:

- class added

...
    
