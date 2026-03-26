package main

type visitor interface {
	visitBinary(n *binary) any
	visitUnary(n *unary) any
	visitGrouping(n *grouping) any
	visitLiteral(n *literal) any
}

type expr interface {
	accept(v visitor) any
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

func (n *binary) accept(v visitor) any {
	return v.visitBinary(n)
}

func (n *unary) accept(v visitor) any {
	return v.visitUnary(n)
}

func (n *grouping) accept(v visitor) any {
	return v.visitGrouping(n)
}

func (n *literal) accept(v visitor) any {
	return v.visitLiteral(n)
}
