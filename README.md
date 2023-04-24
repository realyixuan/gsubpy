# gsubpy

This is an interpreter for subset of Python3 written by golang, which means it will just realize some certain features, some like import, many things will not be implementd. It's intended to be small and simple. (Maybe I will add some features I think cool. At that time, I need give it a new name.)

In a word, gsubpy, a Python interpreter for fun, not big things.

### Quickstart

- install

~~~shell
$ go install github.com/realyixuan/gsubpy
~~~

- running

with repl:

~~~bash
$ gsubpy
>>> print("Hello world")
Hello world 
~~~

or with `.py` file:

~~~shell
$ gsubpy a_py_file.py
~~~

### Supporting features (currently):

- types: `int`, `str`, `list`, `dict`

- builtin: `print`, `len`, `int`, `str`, `bool`, `hash`, `type`, `object`

- statement: `if`, `while`, `def`, `class`, `return`

- operations: dot operation for your own defined attrs (and some special methods), `+-*/` and `> < ==` between integers

- function without keyword arguments

- class without multi-inheritance


### Reference:

- [Writing An Interpreter In Go - Thorsten Ball](https://www.amazon.com/Writing-Interpreter-Go-Thorsten-Ball/dp/3982016118)
