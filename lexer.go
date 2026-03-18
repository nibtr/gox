package main

import "fmt"

type lexer struct {
	source string
}

type token struct {
	tokenType tokenType
	lexeme    string
	literal   any
	line      uint32
}

func (t token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
}

// scanTokens scans the source and extract the tokens
func (l *lexer) scanTokens() []token {
	return []token{}
}

// tokenType enum
type tokenType int

const (
	// single character
	LEFT_PAREN tokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// 1 or 2 character
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)
