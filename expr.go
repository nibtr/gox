package main

type expr interface {
	exprNode() // marker method
}

type binary struct {
	left     expr
	operator token
	right    expr
}

func (b *binary) exprNode()

type unary struct {
	operator token
	right    expr
}

func (u *unary) exprNode()

type grouping struct {
	expression expr
}

func (g *grouping) exprNode()

type literal struct {
	value any
}

func (l *literal) exprNode()
