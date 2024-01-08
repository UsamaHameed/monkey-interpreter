package main

import (
    "fmt"
    "os"
    "github.com/UsamaHameed/monkey-interpreter/repl"
)

func main() {
    fmt.Printf("type some code")
    repl.Start(os.Stdin, os.Stdout)
}
