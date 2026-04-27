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

func NewParser(tokens []token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) Parse() (expr expr, err error) {
	return p.expression()
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

func (p *parser) expression() (expr, error) {
	return p.ternary()
}

func (p *parser) ternary() (expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(QUESTION) {
		thenExpr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(COLON, "Expect ':' after ternary true branch")
		if err != nil {
			return nil, err
		}

		elseExpr, err := p.ternary() // right-associative
		if err != nil {
			return nil, err
		}

		return &ternary{
			condition: expr,
			thenExpr:  thenExpr,
			elseExpr:  elseExpr,
		}, nil
	}

	return expr, nil
}

func (p *parser) equality() (expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator, _ := p.previous()
		right, err := p.comparison()

		if err != nil {
			return nil, err
		}

		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) comparison() (expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator, _ := p.previous()
		right, err := p.term()

		if err != nil {
			return nil, err
		}

		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) term() (expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator, _ := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator, _ := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = &binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (expr, error) {
	if p.match(BANG, MINUS) {
		operator, _ := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		return &unary{
			operator: *operator,
			right:    right,
		}, nil
	}

	return p.primary()
}

func (p *parser) primary() (expr, error) {
	if p.match(FALSE) {
		return &literal{value: false}, nil
	}
	if p.match(TRUE) {
		return &literal{value: true}, nil
	}
	if p.match(NIL) {
		return &literal{value: nil}, nil
	}
	if p.match(NUMBER, STRING) {
		t, _ := p.previous()
		return &literal{value: t.literal}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &grouping{expr}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
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

// consume advances the current pointer if it's the same as `t`
func (p *parser) consume(t tokenType, message string) (*token, error) {
	if p.check(t) {
		return p.advance()
	}

	return nil, p.error(p.peek(), message)
}

// error returns a parserError should any parsing errors occur
func (p *parser) error(t *token, message string) *parseError {
	return &parseError{
		tok:     t,
		message: message,
	}
}
