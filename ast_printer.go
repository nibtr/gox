package main

import "strings"

import "fmt"

type astPrinter struct{}

func (v astPrinter) Print(e expr) (string, error) {
	res, err := e.accept(v)
	if err != nil {
		return "", err
	}

	str, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", res)
	}

	return str, nil
}

func (v astPrinter) visitTernary(n *ternary) (any, error) {
	return v.parenthesize("?:", n.condition, n.thenExpr, n.elseExpr)
}

func (v astPrinter) visitBinary(n *binary) (any, error) {
	return v.parenthesize(n.operator.lexeme, n.left, n.right)
}

func (v astPrinter) visitUnary(n *unary) (any, error) {
	return v.parenthesize(n.operator.lexeme, n.right)
}

func (v astPrinter) visitGrouping(n *grouping) (any, error) {
	return v.parenthesize("group", n.expression)
}

func (v astPrinter) visitLiteral(n *literal) (any, error) {
	switch val := n.value.(type) {
	case nil:
		return "nil", nil
	case string:
		return val, nil
	case int:
		return fmt.Sprintf("%v", val), nil
	case float64:
		return fmt.Sprintf("%g", val), nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil

	default:
		return nil, fmt.Errorf("unknown literal type: %T", val)
	}
}

func (v astPrinter) parenthesize(name string, expressions ...expr) (string, error) {
	var s strings.Builder
	s.WriteString("(" + name)

	for _, e := range expressions {
		s.WriteString(" ")

		str, err := e.accept(v)
		if err != nil {
			return "", err
		}

		st, ok := str.(string)
		if !ok {
			return "", fmt.Errorf("expected string")
		}
		s.WriteString(st)
	}
	s.WriteString(")")

	return s.String(), nil
}
