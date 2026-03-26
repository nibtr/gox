package main

type expr interface {
	accept(v visitor)
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

func (b *binary) accept(v visitor) any {
	return v.visitBinary(b)
}

func (u *unary) accept(v visitor) any {
	return v.visitUnary(u)
}

func (g *grouping) accept(v visitor) any {
	return v.visitGrouping(g)
}

func (l *literal) accept(v visitor) any {
	return v.visitLiteral(l)
}
