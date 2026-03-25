package main

type visitor interface {
	visitBinary(*binary) any
	visitUnary(*unary) any
	visitLiteral(*literal) any
	visitGrouping(*grouping) any
}
