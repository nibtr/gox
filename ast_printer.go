package main

import "strings"

import "fmt"

type astPrinter struct{}

func (v astPrinter) Print(e Expr) (string, error) {
	res, err := e.Accept(v)
	if err != nil {
		return "", err
	}

	str, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", res)
	}

	return str, nil
}

func (v astPrinter) visitTernary(n *Ternary) (any, error) {
	return v.parenthesize("?:", n.condition, n.thenExpr, n.elseExpr)
}

func (v astPrinter) visitBinary(n *Binary) (any, error) {
	return v.parenthesize(n.operator.lexeme, n.left, n.right)
}

func (v astPrinter) visitUnary(n *Unary) (any, error) {
	return v.parenthesize(n.operator.lexeme, n.right)
}

func (v astPrinter) visitGrouping(n *Grouping) (any, error) {
	return v.parenthesize("group", n.expression)
}

func (v astPrinter) visitLiteral(n *Literal) (any, error) {
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

func (v astPrinter) visitVariable(n *Variable) (any, error) {
	return n.name.lexeme, nil
}

func (v astPrinter) parenthesize(name string, expressions ...Expr) (string, error) {
	var s strings.Builder
	s.WriteString("(" + name)

	for _, e := range expressions {
		s.WriteString(" ")

		str, err := e.Accept(v)
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
