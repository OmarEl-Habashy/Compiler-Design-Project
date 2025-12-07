package main

type TokenType string

const (
	KEYWORD        TokenType = "KEYWORD"
	IDENTIFIER     TokenType = "IDENTIFIER"
	CONSTANT       TokenType = "CONSTANT"
	OPERATOR       TokenType = "OPERATOR"
	SPECIAL_SYMBOL TokenType = "SPECIAL_SYMBOL"
)

type Token struct {
	Type  TokenType
	Value string
}
