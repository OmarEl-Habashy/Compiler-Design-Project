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
	stmt := parser.ParseProgram()
	// stmt := parser.ParseIfStatement()

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
	// -----------------
	// Semantic Analysis
	// -----------------
	sem := NewSemanticAnalyzer()
	semTree := sem.Analyze(stmt)

	if len(sem.Errors()) > 0 {
		fmt.Println("\n SEMANTIC MESSAGES:")
		for _, e := range sem.Errors() {
			fmt.Println("  ", e)
		}
	}

	fmt.Println("\n Semantic Tree:")
	fmt.Println(GetSemanticTreeJSON(semTree))

	// -----------------
	// Code Generation
	// -----------------
	fmt.Println("\n--------------------------------")
	fmt.Println("GENERATED ASSEMBLY CODE:")
	fmt.Println("--------------------------------")

	codegen := NewCodeGenerator()
	asmCode := codegen.Generate(stmt)

	fmt.Println(asmCode)

	// Write to file
	err = os.WriteFile("output.asm", []byte(asmCode), 0644)
	if err != nil {
		fmt.Println("Error writing assembly file:", err)
	} else {
		fmt.Println("\n✓ Assembly code written to output.asm")
	}
}
