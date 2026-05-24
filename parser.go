package main

import (
	"fmt"
	"slices"
)

type parser struct {
	tokens  []Token
	current uint32
}

type ParseError struct {
	tok     *Token
	message string
}

func (e *ParseError) Error() string {
	if e.tok.tokenType == EOF {
		return fmt.Sprintf("[line %d] Error at end: %s\n", e.tok.line, e.message)
	}
	return fmt.Sprintf("[line %d] Error at '%s': %s\n",
		e.tok.line,
		e.tok.lexeme,
		e.message,
	)
}

func NewParser(tokens []Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) Parse() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *parser) declaration() (Stmt, error) {
	// TODO: decide parser error strategy:
	// 1: fail-fast (current):
	//   - return err immediately and DO NOT call p.synchronize()
	//   - simplifies design, single error per run, no partial AST execution
	//
	// 2: recovering parser:
	//   - call p.synchronize()
	//   - continue parsing after errors
	//   - requires returning aggregated errors
	if p.match(VAR) {
		v, err := p.varDeclaration()
		if err != nil {
			// p.synchronize()
			return nil, err
		}
		return v, nil
	}

	v, err := p.statement()
	if err != nil {
		// p.synchronize()
		return nil, err
	}
	return v, nil
}

func (p *parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr

	if p.match(EQUAL) {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}
		initializer = v
	}

	_, err = p.consume(SEMICOLON, "expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &VarStmt{name: *name, initializer: initializer}, nil
}

func (p *parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &PrintStmt{Expression: value}, nil
}

func (p *parser) expressionStatement() (Stmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMICOLON, "expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ExpressionStmt{Expression: e}, nil
}

// TODO: remember to synchronize errors
func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if t := p.previous(); t.tokenType == SEMICOLON {
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

func (p *parser) expression() (Expr, error) {
	return p.ternary()
}

func (p *parser) ternary() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(QUESTION) {
		thenExpr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(COLON, "expect ':' after ternary true branch")
		if err != nil {
			return nil, err
		}

		elseExpr, err := p.ternary() // right-associative
		if err != nil {
			return nil, err
		}

		return &Ternary{
			condition: expr,
			thenExpr:  thenExpr,
			elseExpr:  elseExpr,
		}, nil
	}

	return expr, nil
}

func (p *parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()

		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()

		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		return &Unary{
			operator: *operator,
			right:    right,
		}, nil
	}

	return p.primary()
}

func (p *parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{value: false}, nil
	}
	if p.match(TRUE) {
		return &Literal{value: true}, nil
	}
	if p.match(NIL) {
		return &Literal{value: nil}, nil
	}
	if p.match(NUMBER, STRING) {
		t := p.previous()
		return &Literal{value: t.literal}, nil
	}
	if p.match(IDENTIFIER) {
		return &Variable{name: *p.previous()}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "expect ')' after expression."); err != nil {
			return nil, err
		}
		return &Grouping{expr}, nil
	}

	return nil, p.error(p.peek(), "expect expression.")
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
func (p *parser) advance() *Token {
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
func (p *parser) peek() *Token {
	return &p.tokens[p.current]
}

// previous returns the most recently consumed token,
// which is the token just before the current position (current - 1).
func (p *parser) previous() *Token {
	return &p.tokens[p.current-1]
}

// consume advances the current pointer if it's the same as `t`
func (p *parser) consume(t tokenType, message string) (*Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, p.error(p.peek(), message)
}

// error returns a parserError should any parsing errors occur
func (p *parser) error(t *Token, message string) *ParseError {
	return &ParseError{
		tok:     t,
		message: message,
	}
}
