package main

import (
	"fmt"
	"regexp"
)

type Rule struct {
	Type    TokenType
	Pattern *regexp.Regexp
}

var rules = []Rule{
	{KEYWORD, regexp.MustCompile(`^\b(if|else|while|for|int|double|return)\b`)},
	{CONSTANT, regexp.MustCompile(`^\d+(\.\d*)?`)},
	{IDENTIFIER, regexp.MustCompile(`^[a-zA-Z_]\w*`)},
	{OPERATOR, regexp.MustCompile(`^([><=!]=?)`)}, // '?' is a Quantifier that means "ZERO or One of the previous characters, you could have [">=", ">", "==", "=", "<=", "<", "!"]"
	{SPECIAL_SYMBOL, regexp.MustCompile(`^[(){};]`)},
}

var whitespacePattern = regexp.MustCompile(`^\s+`) //skip whiteSpaces

func Lex(input string) []Token {
	var tokens []Token
	pos := 0

	for pos < len(input) {
		//1. skip whitespaces
		if loc := whitespacePattern.FindStringIndex(input[pos:]); loc != nil {
			pos += loc[1]
			continue
		}

		//2. match rules
		matchFound := false
		for _, rule := range rules {
			loc := rule.Pattern.FindStringIndex(input[pos:])
			if loc != nil {
				matchText := input[pos : pos+loc[1]]
				tokens = append(tokens, Token{Type: rule.Type, Value: matchText})
				pos += loc[1]
				matchFound = true
				break
			}
		}

		//3. catch error
		if !matchFound {
			fmt.Printf("Lexer Error : Unkonwn character %c at position %d\n", input[pos], pos)
			pos++
		}
	}
	return tokens
}
