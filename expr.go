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
