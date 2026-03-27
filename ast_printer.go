package main

import "strings"

import "fmt"

type astPrinter struct{}

func (v astPrinter) print(e expr) string {
	res := e.accept(v)
	str, ok := res.(string)
	if !ok {
		panic(fmt.Sprintf("expected string, got %T", res))
	}

	return str
}

func (v astPrinter) visitBinary(n *binary) any {
	return v.parenthesize(n.operator.lexeme, n.left, n.right)
}

func (v astPrinter) visitUnary(n *unary) any {
	return v.parenthesize(n.operator.lexeme, n.right)
}

func (v astPrinter) visitGrouping(n *grouping) any {
	return v.parenthesize("group", n.expression)
}

func (v astPrinter) visitLiteral(n *literal) any {
	switch val := n.value.(type) {
	case nil:
		return "nil"
	case string:
		return val
	case int:
		return fmt.Sprintf("%v", val)
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"

	default:
		panic(fmt.Sprintf("unknown literal type: %T", val))
	}
}

func (v astPrinter) parenthesize(name string, expressions ...expr) string {
	var s strings.Builder
	s.WriteString("(" + name)

	for _, e := range expressions {
		s.WriteString(" ")
		str, ok := e.accept(v).(string)
		if !ok {
			panic("expected string")
		}
		s.WriteString(str)
	}
	s.WriteString(")")

	return s.String()
}
