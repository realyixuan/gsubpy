package main

import (
    "fmt"
    "os"

    "os/exec"
    "path/filepath"
)

func main() {
    info, _ := os.Stat(os.Args[1])

    var pass bool = true

    if info.IsDir() {
        entries, _ := os.ReadDir(os.Args[1])
     
        for _, e := range entries {
            if !e.IsDir() {
                filePath := filepath.Join(os.Args[1], e.Name())
                cmd := exec.Command("go", "run", "main.go", filePath)
                _, err := cmd.Output()
                if err != nil {
                    fmt.Println("[x]", filePath)
                    pass = false
                } else {
                    fmt.Println("[v]", filePath)
                }
            }
        }
    } else {
        cmd := exec.Command("go", "run", "main.go", os.Args[1])
        _, err := cmd.Output()
        if err != nil {
            fmt.Println("[x]", os.Args[1])
        } else {
            fmt.Println("[v]", os.Args[1])
        }
    }

    if !pass {
        os.Exit(1)
    }

    fmt.Println("all passed")
}
