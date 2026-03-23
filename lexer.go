package main

import (
	"strconv"
)

const NULL_CHARACTER byte = '\x00'

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

	case '/':
		// comments
		if l.match('/') {
			// keep consuming (skipping) char until newline or end of source
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			l.addToken(SLASH)
		}

	case '"':
		l.string()

	// skip whitespace
	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		l.line++

	default:
		if l.isDigit(c) {
			l.number()
		} else if l.isAlpha(c) {
			l.identifier()
		} else {
			printError(l.line, "Unexpected character.")
		}
	}
}

// advance returns the current character and then advances `current` by 1.
func (l *lexer) advance() byte {
	l.current++
	return l.source[l.current-1]
}

// match checks if the current character equals to `expected`
// and advances `current` by 1 if they do match
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

// peek returns the current character without advancing `current`
func (l *lexer) peek() byte {
	if l.isAtEnd() {
		return NULL_CHARACTER // null character
	}

	return l.source[l.current]
}

// peekNext returns the character right next to `current`, without advancing `current`
func (l *lexer) peekNext() byte {
	if l.current+1 >= uint32(len(l.source)) {
		return NULL_CHARACTER
	}

	return l.source[l.current+1]
}

// string consumes the `current` to get the literal string value
func (l *lexer) string() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}

	if l.isAtEnd() {
		printError(l.line, "Unterminated string")
		return
	}

	l.advance() // the closing "

	value := l.source[l.start+1 : l.current-1]
	l.addToken(STRING, value)
}

// number checks if current token is a number
func (l *lexer) number() {
	// consume digits
	for l.isDigit(l.peek()) {
		l.advance()
	}

	// look for fractional part
	if l.peek() == '.' && l.isDigit(l.peekNext()) {
		// consume the '.'
		l.advance()

		// keep consuming digits after '.'
		for l.isDigit(l.peek()) {
			l.advance()
		}
	}

	// TODO: check this error again, see if need to return error
	value, _ := strconv.ParseFloat(l.source[l.start:l.current], 64)
	l.addToken(NUMBER, value)
}

// identifier checks if current token is an identifier
func (l *lexer) identifier() {
	for l.isAlphaNumeric(l.peek()) {
		l.advance()
	}

	value := l.source[l.start:l.current]
	t := IDENTIFIER
	if v, ok := keywords[value]; ok {
		t = v
	}

	l.addToken(t)
}

// addToken append a new token to `tokens`
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

// isDigit checks if a character is a digit
func (l *lexer) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isAlpha checks if a character is an alphabet
func (l *lexer) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

// isAlphaNumeric checks if a character is either a digit or an alphabet
func (l *lexer) isAlphaNumeric(c byte) bool {
	return l.isAlpha(c) || l.isDigit(c)
}

// isAtEnd checks if the current pointer is at the end of source
func (l *lexer) isAtEnd() bool {
	return l.current >= uint32(len(l.source))
}
