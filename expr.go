package main

type expr interface {
	accept(v visitor)
}

type binary struct {
	left     expr
	operator token
	right    expr
}

func (b *binary) accept(v visitor) any {
	return v.visitBinary(b)
}

type unary struct {
	operator token
	right    expr
}

func (u *unary) accept(v visitor) any {
	return v.visitUnary(u)
}

type grouping struct {
	expression expr
}

func (g *grouping) accept(v visitor) any {
	return v.visitGrouping(g)
}

type literal struct {
	value any
}

func (l *literal) accept(v visitor) any {
	return v.visitLiteral(l)
}
