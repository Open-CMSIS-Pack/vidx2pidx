package commands

import "fmt"

var Vidx = struct {
    Add func(name, source string) int
}{
    Add: add,
}

func add(name, source string) int {
    fmt.Println("Adding from vidx")
    fmt.Println("Name ", name)
    fmt.Println("Source ", source)
    return 0
}
