package main

type ExprVisitor interface {
	// visitTernary evaluates condition ? thenExpr : elseExpr
	visitTernary(n *ternary) (any, error)
	// visitBinary evaluates binary expressions like +, -, *, /, comparisons
	visitBinary(n *binary) (any, error)
	// visitUnary evaluates unary expressions like -x and !x
	visitUnary(n *unary) (any, error)
	// visitGrouping evaluates expressions inside parentheses (...)
	visitGrouping(n *grouping) (any, error)
	// visitLiteral returns a literal value directly
	visitLiteral(n *literal) (any, error)
}

type expr interface {
	accept(v ExprVisitor) (any, error)
}

type ternary struct {
	condition expr
	thenExpr  expr
	elseExpr  expr
}

type binary struct {
	left     expr
	operator token
	right    expr
}

type unary struct {
	operator token
	right    expr
}

type grouping struct {
	expression expr
}

type literal struct {
	value any
}

func (n *ternary) accept(v ExprVisitor) (any, error) {
	return v.visitTernary(n)
}

func (n *binary) accept(v ExprVisitor) (any, error) {
	return v.visitBinary(n)
}

func (n *unary) accept(v ExprVisitor) (any, error) {
	return v.visitUnary(n)
}

func (n *grouping) accept(v ExprVisitor) (any, error) {
	return v.visitGrouping(n)
}

func (n *literal) accept(v ExprVisitor) (any, error) {
	return v.visitLiteral(n)
}
