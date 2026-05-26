package parser

import (
	"fmt"
	"slices"

	"github.com/nibtr/gox/ast"
	"github.com/nibtr/gox/lexer"
)

type parser struct {
	tokens  []lexer.Token
	current uint32
}

type ParseError struct {
	tok     *lexer.Token
	message string
}

func (e *ParseError) Error() string {
	if e.tok.TokenType == lexer.EOF {
		return fmt.Sprintf("[line %d] Error at end: %s\n", e.tok.Line, e.message)
	}
	return fmt.Sprintf("[line %d] Error at '%s': %s\n",
		e.tok.Line,
		e.tok.Lexeme,
		e.message,
	)
}

func NewParser(tokens []lexer.Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) Parse() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *parser) declaration() (ast.Stmt, error) {
	// TODO: decide parser error strategy:
	// 1: fail-fast (current):
	//   - return err immediately and DO NOT call p.synchronize()
	//   - simplifies design, single error per run, no partial AST execution
	//
	// 2: recovering parser:
	//   - call p.synchronize()
	//   - continue parsing after errors
	//   - requires returning aggregated errors
	if p.match(lexer.VAR) {
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

func (p *parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(lexer.IDENTIFIER, "expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr

	if p.match(lexer.EQUAL) {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}
		initializer = v
	}

	_, err = p.consume(lexer.SEMICOLON, "expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &ast.VarStmt{Name: *name, Initializer: initializer}, nil
}

func (p *parser) statement() (ast.Stmt, error) {
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}

	if p.match(lexer.LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.BlockStmt{
			Statements: stmts,
		}, nil
	}

	return p.expressionStatement()
}

func (p *parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(lexer.SEMICOLON, "expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStmt{Expression: value}, nil
}

func (p *parser) block() ([]ast.Stmt, error) {
	stmts := []ast.Stmt{}
	for !p.check(lexer.RIGHT_BRACE) && !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, dec)
	}
	if _, err := p.consume(lexer.RIGHT_BRACE, "expect '}' after block."); err != nil {
		return nil, err
	}
	return stmts, nil
}

func (p *parser) expressionStatement() (ast.Stmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(lexer.SEMICOLON, "expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStmt{Expression: e}, nil
}

// TODO: remember to synchronize errors
func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if t := p.previous(); t.TokenType == lexer.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case lexer.CLASS:
		case lexer.FUNC:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return
		}

		p.advance()
	}
}

func (p *parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *parser) assignment() (ast.Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	if p.match(lexer.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		// expr is a Variable
		if v, ok := expr.(*ast.Variable); ok {
			name := v.Name
			return &ast.Assign{Name: name, Value: value}, nil
		}

		p.error(equals, "invalid assignment target.")
	}

	return expr, nil
}

func (p *parser) ternary() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(lexer.QUESTION) {
		thenExpr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(lexer.COLON, "expect ':' after ternary true branch")
		if err != nil {
			return nil, err
		}

		elseExpr, err := p.ternary() // right-associative
		if err != nil {
			return nil, err
		}

		return &ast.Ternary{
			Condition: expr,
			ThenExpr:  thenExpr,
			ElseExpr:  elseExpr,
		}, nil
	}

	return expr, nil
}

func (p *parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()

		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()

		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = &ast.Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (ast.Expr, error) {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		return &ast.Unary{
			Operator: *operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *parser) primary() (ast.Expr, error) {
	if p.match(lexer.FALSE) {
		return &ast.Literal{Value: false}, nil
	}
	if p.match(lexer.TRUE) {
		return &ast.Literal{Value: true}, nil
	}
	if p.match(lexer.NIL) {
		return &ast.Literal{Value: nil}, nil
	}
	if p.match(lexer.NUMBER, lexer.STRING) {
		t := p.previous()
		return &ast.Literal{Value: t.Literal}, nil
	}
	if p.match(lexer.IDENTIFIER) {
		return &ast.Variable{Name: *p.previous()}, nil
	}
	if p.match(lexer.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(lexer.RIGHT_PAREN, "expect ')' after expression."); err != nil {
			return nil, err
		}
		return &ast.Grouping{Expression: expr}, nil
	}

	return nil, p.error(p.peek(), "expect expression.")
}

//
// ======== HELPERS ========
//

// match checks whether the current token matches any of the given types.
// If a match is found, it advances the parser to the next token and returns true.
// If none of the types match, it leaves the parser unchanged and returns false.
func (p *parser) match(types ...lexer.TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}

	return false
}

// check checks if token at `current` is equal to `t`
func (p *parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == t
}

// advance consumes the token at `current` and returns it,
// then advances `current` to next token
func (p *parser) advance() *lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if the token at `current` is an EOF
func (p *parser) isAtEnd() bool {
	return p.peek().TokenType == lexer.EOF
}

// peek returns the token at `current`
func (p *parser) peek() *lexer.Token {
	return &p.tokens[p.current]
}

// previous returns the most recently consumed token,
// which is the token just before the current position (current - 1).
func (p *parser) previous() *lexer.Token {
	return &p.tokens[p.current-1]
}

// consume advances the current pointer if it's the same as `t`
func (p *parser) consume(t lexer.TokenType, message string) (*lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, p.error(p.peek(), message)
}

// error returns a parserError should any parsing errors occur
func (p *parser) error(t *lexer.Token, message string) *ParseError {
	return &ParseError{
		tok:     t,
		message: message,
	}
}
