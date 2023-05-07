package pytest

import (
    "fmt"
    "os"

    "os/exec"
    "path/filepath"
)

func Main() {
    info, _ := os.Stat(os.Args[2])

    var pass bool = true

    if info.IsDir() {
        entries, _ := os.ReadDir(os.Args[2])
     
        for _, e := range entries {
            if !e.IsDir() {
                filePath := filepath.Join(os.Args[2], e.Name())
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
        cmd := exec.Command("go", "run", "main.go", os.Args[2])
        _, err := cmd.Output()
        if err != nil {
            fmt.Println("[x]", os.Args[2])
        } else {
            fmt.Println("[v]", os.Args[2])
        }
    }

    if !pass {
        os.Exit(1)
    }

    fmt.Println("all passed")
}
