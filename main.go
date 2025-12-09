package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Reading file...")
	fileName := "test_code.txt"
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("file name " + fileName)
	fmt.Printf("file size %d bytes\n", len(data))
	fmt.Println("--------------------------------")

	code := string(data)

	fmt.Println("Input Code:\n", code)
	fmt.Println("--------------------------------")

	tokenList := Lex(code)

	for _, t := range tokenList {
		fmt.Println(t)
	}


	
	fmt.Println("--------------------------------")
    
    parser := NewParser(tokenList)
    stmt := parser.ParseIfStatement()

    // Check for parser errors
    if len(parser.Errors()) > 0 {
        fmt.Println("\n PARSER ERRORS:")
        for _, err := range parser.Errors() {
            fmt.Println("  Error:", err)
        }
        return
    }
    fmt.Println("\n Parse Tree (Abstract Syntax Tree):")
	
    fmt.Println(GetParseTreeJSON(stmt))
}