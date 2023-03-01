package main

import (
    "os"

    "gsubpy/repl"
)

func main() {
    if len(os.Args) == 1 {
        repl.REPLRunning()
    }
}

