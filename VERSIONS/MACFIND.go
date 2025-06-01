package main

import (
    _ "embed"
    "fmt"
)

//go:embed oui_v2.txt
var fileData string

func main() {
    fmt.Println("Embedded File Content:")
    fmt.Println(fileData)
}
