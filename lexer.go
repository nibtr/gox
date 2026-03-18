package main

import "fmt"

type lexer struct {
	source  string
	tokens  []token
	start   uint32
	current uint32
	line    uint32
}

func newLexer(source string) *lexer {
	return &lexer{
		source: source,
		line:   1,
	}
}

// scanTokens scans the source and extract the tokens
func (l *lexer) scanTokens() []token {
	for !l.isAtEnd() {
		// we're at the beginning of the next lexeme
		l.start = l.current
		l.scanTokens()
	}

	l.tokens = append(l.tokens, newToken(EOF, "", nil, l.line))
	return l.tokens
}

func (l *lexer) isAtEnd() bool {
	return l.current >= uint32(len(l.source))
}

type token struct {
	tokenType tokenType
	lexeme    string
	literal   any
	line      uint32
}

func newToken(tokenType tokenType, lexeme string, literal any, line uint32) token {
	return token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
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
