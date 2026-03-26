package main

type visitor interface {
	visitBinary(n *binary) any
	visitUnary(n *unary) any
	visitLiteral(n *literal) any
	visitGrouping(n *grouping) any
}
