package main

import (
	"fmt"
	"slices"
)

type parser struct {
	tokens  []token
	current uint32
}

type parseError struct {
	tok     *token
	message string
}

func (e *parseError) Error() string {
	if e.tok.tokenType == EOF {
		return fmt.Sprintf("[line %d] Error at end: %s", e.tok.line, e.message)
	}
	return fmt.Sprintf("[line %d] Error at '%s': %s",
		e.tok.line,
		e.tok.lexeme,
		e.message,
	)
}

func newParser(tokens []token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) Parse() (expr expr, err error) {
	return p.expression(), nil
}

func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if t, _ := p.previous(); t.tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
		case FUNC:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}

func (p *parser) expression() expr {
	return p.equality()
}

func (p *parser) equality() expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator, _ := p.previous()
		right := p.comparison()
		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser) comparison() expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator, _ := p.previous()
		right := p.term()
		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser) term() expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator, _ := p.previous()
		right := p.factor()
		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser) factor() expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator, _ := p.previous()
		right := p.unary()
		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr
}

func (p *parser) unary() expr {
	if p.match(BANG, MINUS) {
		operator, _ := p.previous()
		right := p.unary()
		return &unary{
			operator: *operator,
			right:    right,
		}
	}

	return p.primary()
}

func (p *parser) primary() expr {
	if p.match(FALSE) {
		return &literal{value: false}
	}
	if p.match(TRUE) {
		return &literal{value: true}
	}
	if p.match(NIL) {
		return &literal{value: nil}
	}
	if p.match(NUMBER, STRING) {
		t, _ := p.previous()
		return &literal{value: t.literal}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &grouping{expr}
	}

	// TODO: handle error return instead of panic
	panic(p.error(p.peek(), "Expect expression."))
}

//
// ======== HELPERS ========
//

// match checks whether the current token matches any of the given types.
// If a match is found, it advances the parser to the next token and returns true.
// If none of the types match, it leaves the parser unchanged and returns false.
func (p *parser) match(types ...tokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}

	return false
}

// check checks if token at `current` is equal to `t`
func (p *parser) check(t tokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().tokenType == t
}

// advance consumes the token at `current` and returns it,
// then advances `current` to next token
func (p *parser) advance() (*token, error) {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if the token at `current` is an EOF
func (p *parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

// peek returns the token at `current`
func (p *parser) peek() *token {
	return &p.tokens[p.current]
}

// previous returns the most recently consumed token,
// which is the token just before the current position (current - 1).
func (p *parser) previous() (*token, error) {
	return &p.tokens[p.current-1], nil
}

func (p *parser) consume(t tokenType, message string) (*token, error) {
	if p.check(t) {
		return p.advance()
	}

	return nil, p.error(p.peek(), message)
}

func (p *parser) error(t *token, message string) *parseError {
	return &parseError{
		tok:     t,
		message: message,
	}
}
