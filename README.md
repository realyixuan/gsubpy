# gsubpy

This is an interpreter for subset of Python3 written by golang, which means it will just realize some certain features, some like import, many things will not be implementd. It's intended to be small and simple. (Maybe I will add some features I think cool. At that time, I need give it a new name.)

In a word, gsubpy, a Python interpreter for fun, not big things.

### Quickstart

- install

~~~shell
$ go install github.com/realyixuan/gsubpy@latest
~~~

- running

with repl:

~~~bash
$ gsubpy
>>> print("Hello world")
Hello world 
~~~

or with `.py` file (there are some examples under `demos/`):

~~~shell
$ gsubpy a_py_file.py
~~~

### Supporting features:

- data: `int`, `str`, `list`, `dict`

- builtin: `print`, `len`, `int`, `str`, `bool`, `hash`, `type`, `object`, `id`, `Exception`, `StopIteration`, `list`, `dict`, `isinstance`, `issubclass`, `iter`, `next`, `range`, `max`, `min`, `dir`

- statement: `if`, `while`, `def`, `class`, `return`, `break`, `for`, `break`, `continue`, `raise`, `assert`

- operations:

    - dot operation for your own defined attrs (and some special methods)

    - `+-*/`

    - `>`, `<`, `==`, `!=`

    - `not`, `in`, `not in`, `is`, `is not`, `and`, `or`

- function without keyword arguments

- class without multi-inheritance

### Supports

<a href="https://jb.gg/OpenSourceSupport">
    <img width="100" height="100" src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo.">
</a>

### Reference:

- [Writing An Interpreter In Go - Thorsten Ball](https://www.amazon.com/Writing-Interpreter-Go-Thorsten-Ball/dp/3982016118)

