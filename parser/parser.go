package parser

import (
	"errors"
	"fmt"
	"slices"

	"github.com/nibtr/gox/ast"
	"github.com/nibtr/gox/lexer"
)

const (
	fnKindFunction = "function"
	fnKindMethod   = "method"
)

type parser struct {
	tokens    []lexer.Token
	current   uint32
	loopDepth uint32
	// errors accumulates recovered parse errors so ParseProgram can report all of them.
	errors []error
}

type ParseError struct {
	Token   *lexer.Token
	Message string
}

func (e *ParseError) Error() string {
	if e.Token.TokenType == lexer.EOF {
		return fmt.Sprintf("[line %d] error at end: %s\n", e.Token.Line-1, e.Message)
	}
	return fmt.Sprintf("[line %d] error at '%s': %s\n",
		e.Token.Line,
		e.Token.Lexeme,
		e.Message,
	)
}

func NewParser(tokens []lexer.Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) ParseExpression() (ast.Expr, error) {
	return p.expression()
}

func (p *parser) ParseProgram() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}
	for !p.IsAtEnd() {
		stmt := p.declaration()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	if len(p.errors) > 0 {
		return nil, errors.Join(p.errors...)
	}

	return statements, nil
}

// declaration recovers from parse errors at statement boundaries: on error
// it records the error, synchronizes to the next statement, and returns nil
// so the caller (ParseProgram or block) just skips it and keeps parsing.
func (p *parser) declaration() ast.Stmt {
	var stmt ast.Stmt
	var err error

	switch {
	case p.match(lexer.CLASS):
		stmt, err = p.classDeclaration()
	case p.match(lexer.FUNC):
		stmt, err = p.function(fnKindFunction)
	case p.match(lexer.VAR):
		stmt, err = p.varDeclaration()
	default:
		stmt, err = p.statement()
	}

	if err != nil {
		p.errors = append(p.errors, err)
		p.synchronize()
		return nil
	}

	return stmt
}

func (p *parser) function(kind string) (ast.Stmt, error) {
	name, err := p.consume(lexer.IDENTIFIER, fmt.Sprintf("expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(lexer.LEFT_PAREN, fmt.Sprintf("expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}

	params := []lexer.Token{}
	if !p.check(lexer.RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				return nil, p.error(p.peek(), "can't have more than 255 parameters.")
			}

			v, err := p.consume(lexer.IDENTIFIER, "expect parameter name")
			if err != nil {
				return nil, err
			}
			params = append(params, *v)

			if !p.match(lexer.COMMA) {
				break
			}
		}
	}

	_, err = p.consume(lexer.RIGHT_PAREN, "expect ')' after parameters.")
	_, err = p.consume(lexer.LEFT_BRACE, fmt.Sprintf("expect '{' before %s body.", kind))

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStmt{
		Name:   *name,
		Params: params,
		Body:   body,
	}, nil
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

func (p *parser) classDeclaration() (ast.Stmt, error) {
	name, err := p.consume(lexer.IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(lexer.LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := []ast.FunctionStmt{}
	for !p.check(lexer.RIGHT_BRACE) && !p.IsAtEnd() {
		f, err := p.function(fnKindMethod)
		if err != nil {
			return nil, err
		}

		if val, ok := f.(*ast.FunctionStmt); ok {
			methods = append(methods, *val)
		}
	}

	_, err = p.consume(lexer.RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &ast.ClassStmt{
		Name:    *name,
		Methods: methods,
	}, nil
}

func (p *parser) statement() (ast.Stmt, error) {
	if p.match(lexer.FOR) {
		return p.forStatement()
	}
	if p.match(lexer.IF) {
		return p.ifStatement()
	}
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}
	if p.match(lexer.RETURN) {
		return p.returnStatement()
	}
	if p.match(lexer.WHILE) {
		return p.whileStatement()
	}
	if p.match(lexer.BREAK) {
		return p.breakStatement()
	}
	if p.match(lexer.CONTINUE) {
		return p.continueStatement()
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

func (p *parser) forStatement() (ast.Stmt, error) {
	p.loopDepth += 1
	defer func() {
		p.loopDepth -= 1
	}()

	var initializer ast.Stmt
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		v, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}
		initializer = v
	} else {
		v, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}
		initializer = v
	}

	var condition ast.Expr
	if !p.check(lexer.SEMICOLON) {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}
		condition = v
	}

	_, err := p.consume(lexer.SEMICOLON, "expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment ast.Expr
	if !p.check(lexer.LEFT_BRACE) {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}
		increment = v
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if condition == nil {
		// if no condition, infinite loop
		condition = &ast.Literal{Value: true}
	}
	body = &ast.WhileStmt{Condition: condition, Body: body, Increment: increment}
	if initializer != nil {
		body = &ast.BlockStmt{Statements: []ast.Stmt{
			initializer,
			body,
		}}
	}

	return body, nil
}

func (p *parser) ifStatement() (ast.Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	if p.match(lexer.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *parser) whileStatement() (ast.Stmt, error) {
	p.loopDepth += 1
	defer func() {
		p.loopDepth -= 1
	}()

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{Condition: condition, Body: body}, nil
}

func (p *parser) breakStatement() (ast.Stmt, error) {
	if p.loopDepth == 0 {
		return nil, p.error(p.peek(), "break outside loop.")
	}
	_, err := p.consume(lexer.SEMICOLON, "expect ';' after break.")
	if err != nil {
		return nil, err
	}

	return &ast.BreakStmt{}, nil
}

func (p *parser) continueStatement() (ast.Stmt, error) {
	if p.loopDepth == 0 {
		return nil, p.error(p.peek(), "continue outside loop.")
	}
	_, err := p.consume(lexer.SEMICOLON, "expect ';' after continue.")
	if err != nil {
		return nil, err
	}

	return &ast.ContinueStmt{}, nil
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

func (p *parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()
	var value ast.Expr
	if !p.check(lexer.SEMICOLON) {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}
		value = v
	}

	_, err := p.consume(lexer.SEMICOLON, "expect ';' after return value")
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStmt{Keyword: *keyword, Value: value}, nil
}

func (p *parser) block() ([]ast.Stmt, error) {
	stmts := []ast.Stmt{}
	for !p.check(lexer.RIGHT_BRACE) && !p.IsAtEnd() {
		dec := p.declaration()
		if dec != nil {
			stmts = append(stmts, dec)
		}
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

func (p *parser) synchronize() {
	p.advance()

	for !p.IsAtEnd() {
		if t := p.previous(); t.TokenType == lexer.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case lexer.CLASS, lexer.FUNC, lexer.VAR, lexer.FOR, lexer.IF,
			lexer.WHILE, lexer.PRINT, lexer.RETURN:
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
		if v, ok := expr.(*ast.GetExpr); ok {
			get := v
			return &ast.SetExpr{
				Object: get.Object,
				Name:   get.Name,
				Value:  value,
			}, nil
		}

		p.error(equals, "invalid assignment target.")
	}

	return expr, nil
}

func (p *parser) ternary() (ast.Expr, error) {
	expr, err := p.or()
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

func (p *parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
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

	return p.call()
}

func (p *parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(lexer.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(lexer.DOT) {
			name, err := p.consume(lexer.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &ast.GetExpr{Object: expr, Name: *name}
			break
		} else {
			break
		}
	}

	return expr, nil
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

func (p *parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := []ast.Expr{}

	if !p.check(lexer.RIGHT_PAREN) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			if len(arguments) >= 255 {
				// TODO: currently fail-fast
				// if want panic mode, need to use with synchronize
				return nil, p.error(p.peek(), "can't have more than 255 arguments.")
			}
			arguments = append(arguments, arg)

			if !p.match(lexer.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(lexer.RIGHT_PAREN, "expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.Call{
		Callee:    callee,
		Paren:     *paren,
		Arguments: arguments,
	}, nil
}

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
	if p.IsAtEnd() {
		return false
	}

	return p.peek().TokenType == t
}

// advance consumes the token at `current` and returns it,
// then advances `current` to next token
func (p *parser) advance() *lexer.Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.previous()
}

// IsAtEnd checks if the token at `current` is an EOF
func (p *parser) IsAtEnd() bool {
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
		Token:   t,
		Message: message,
	}
}
