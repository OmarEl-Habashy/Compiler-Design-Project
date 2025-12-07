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
}
