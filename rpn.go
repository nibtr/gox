package main

import (
	"fmt"
	"strings"
)

type rpn struct{}

func (v rpn) convert(e expr) any {
	return e.accept(v)
}

func (v rpn) visitBinary(n *binary) any {
	return v.build(n.operator.lexeme, n.left, n.right)
}

func (v rpn) visitUnary(n *unary) any {
	return v.build(n.operator.lexeme, n.right)
}

func (v rpn) visitGrouping(n *grouping) any {
	return n.expression.accept(v)
}

func (v rpn) visitLiteral(n *literal) any {
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

func (v rpn) build(lexeme string, expressions ...expr) string {
	var s strings.Builder

	for _, e := range expressions {
		str, ok := e.accept(v).(string)
		if !ok {
			panic("expected string")
		}

		s.WriteString(str)
		s.WriteString(" ")
	}

	s.WriteString(lexeme)

	return s.String()
}
