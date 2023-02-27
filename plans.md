# Some Thoughts For realizing an interpreter of subset of Python

There, Python has a lot of syntax, implementing them all at once is nearly impossible. So, it's better to take a small step at a time.

### Here's the rough plans:

Along realizing it, the development is parted into several phases, then solve all, step by step, phase by phase.

##### Phase 1: most-basic features

desired outcomes:

- being able to interpreter following code:

~~~python
total = 0
i = 0
while i < 10:
    total = total + i
    i = i + 1
if total > 5:
    total = 10
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

the key point of this phase is just getting things done. Keep general direction right, other than that I shouldn't take care of trivial details, as many things havne't got very clear, at least I think so.

steps:

step1: `assignment identifier = <expression>`

step2: `if statement`

step3: `while-loop statement`

phase 1 should get done as quickly as possible.

##### Phase 2: enhancement

desired outcomes:

- built function putted in
- more builtin data types
    - str
    - list
    - dict
- simple function

...

##### Phase 3: superior

desired outcomes:

- class added

...
    
