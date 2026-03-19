package main

import "fmt"

type lexer struct {
	source string
	tokens []token

	// points to 1st char currently being considered
	start uint32
	// points at the char currently being considered
	current uint32
	// the line `current` is on
	line uint32
}

// newLexer returns a new instance of lexer
func newLexer(source string) *lexer {
	return &lexer{
		source: source,
		line:   1,

		// just being explicit
		tokens:  []token{},
		start:   0,
		current: 0,
	}
}

// scanTokens scans the source and extract the tokens
func (l *lexer) scanTokens() []token {
	for !l.isAtEnd() {
		l.scanToken()
	}

	// add an EOF at the end of source
	l.tokens = append(l.tokens, newToken(EOF, "", nil, l.line))
	return l.tokens
}

// scanToken scans an individual token,
// mutates the `current` pointer in the process
func (l *lexer) scanToken() {
	// reset `start` so that we're at the beginning of the next lexeme
	l.start = l.current

	// this fn returns the current character THEN moves `current` up by 1
	c := l.advance()

	switch c {
	case '(':
		l.addToken(LEFT_PAREN)
	case ')':
		l.addToken(RIGHT_PAREN)
	case '{':
		l.addToken(LEFT_BRACE)
	case '}':
		l.addToken(RIGHT_BRACE)
	case ',':
		l.addToken(COMMA)
	case '.':
		l.addToken(DOT)
	case '-':
		l.addToken(MINUS)
	case '+':
		l.addToken(PLUS)
	case ';':
		l.addToken(SEMICOLON)
	case '*':
		l.addToken(STAR)

	case '!':
		if l.match('=') {
			l.addToken(BANG_EQUAL)
		} else {
			l.addToken(BANG)
		}
	case '=':
		if l.match('=') {
			l.addToken(EQUAL_EQUAL)
		} else {
			l.addToken(EQUAL)
		}
	case '<':
		if l.match('=') {
			l.addToken(LESS_EQUAL)
		} else {
			l.addToken(LESS)
		}
	case '>':
		if l.match('=') {
			l.addToken(GREATER_EQUAL)
		} else {
			l.addToken(GREATER)
		}

	default:
		printError(l.line, "Unexpected character.")
	}
}

// advance returns the current character and then advances `current` by 1.
func (l *lexer) advance() byte {
	l.current++
	return l.source[l.current-1]
}

// match checks if the current character equals to `expected`
// and then advances `current` by 1
func (l *lexer) match(expected byte) bool {
	if l.isAtEnd() {
		return false
	}

	if l.source[l.current] != expected {
		return false
	}

	l.current++
	return true
}

func (l *lexer) addToken(tokenType tokenType, literal ...any) {
	lexeme := l.source[l.start:l.current]
	var value any = nil

	switch len(literal) {
	case 0:
		value = nil
	case 1:
		value = literal[0]
	default:
		panic("addToken: too many literal arguments")
	}

	l.tokens = append(l.tokens, newToken(tokenType, lexeme, value, l.line))
}

// isAtEnd checks if the current pointer is at the end of source
func (l *lexer) isAtEnd() bool {
	return l.current >= uint32(len(l.source))
}

type token struct {
	tokenType tokenType
	lexeme    string
	literal   any
	line      uint32
}

// newToken returns a new token instance
func newToken(tokenType tokenType, lexeme string, literal any, line uint32) token {
	return token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t token) String() string {
	return fmt.Sprintf("%s %s %v", tokenName[t.tokenType], t.lexeme, t.literal)
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

var tokenName = map[tokenType]string{
	LEFT_PAREN:  "LEFT_PAREN",
	RIGHT_PAREN: "RIGHT_PAREN",
	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
	COMMA:       "COMMA",
	DOT:         "DOT",
	MINUS:       "MINUS",
	PLUS:        "PLUS",
	SEMICOLON:   "SEMICOLON",
	SLASH:       "SLASH",
	STAR:        "STAR",

	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",

	EOF: "EOF",
}
