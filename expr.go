package main

type visitor interface {
	visitTernary(n *ternary) (any, error)
	visitBinary(n *binary) (any, error)
	visitUnary(n *unary) (any, error)
	visitGrouping(n *grouping) (any, error)
	visitLiteral(n *literal) (any, error)
}

type expr interface {
	accept(v visitor) (any, error)
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

func (n *ternary) accept(v visitor) (any, error) {
	return v.visitTernary(n)
}

func (n *binary) accept(v visitor) (any, error) {
	return v.visitBinary(n)
}

func (n *unary) accept(v visitor) (any, error) {
	return v.visitUnary(n)
}

func (n *grouping) accept(v visitor) (any, error) {
	return v.visitGrouping(n)
}

func (n *literal) accept(v visitor) (any, error) {
	return v.visitLiteral(n)
}
