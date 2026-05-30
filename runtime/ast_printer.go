package runtime

import (
	"fmt"
	"strings"

	"github.com/nibtr/gox/ast"
)

type astPrinter struct{}

func (v astPrinter) Print(e ast.Expr) (string, error) {
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

// VisitCall implements [ast.ExprVisitor].
func (v astPrinter) VisitCall(n *ast.Call) (any, error) {
	exprs := make([]ast.Expr, 0, len(n.Arguments)+1)
	exprs = append(exprs, n.Callee)
	exprs = append(exprs, n.Arguments...)

	return v.parenthesize("call", exprs...)
}

func (v astPrinter) VisitAssignExpr(n *ast.Assign) (any, error) {
	return v.parenthesize(
		"assign "+n.Name.Lexeme,
		n.Value,
	)
}

func (v astPrinter) VisitTernary(n *ast.Ternary) (any, error) {
	return v.parenthesize("?:", n.Condition, n.ThenExpr, n.ElseExpr)
}

func (v astPrinter) VisitLogical(n *ast.Logical) (any, error) {
	return v.parenthesize(
		n.Operator.Lexeme,
		n.Left,
		n.Right,
	)
}

func (v astPrinter) VisitBinary(n *ast.Binary) (any, error) {
	return v.parenthesize(n.Operator.Lexeme, n.Left, n.Right)
}

func (v astPrinter) VisitUnary(n *ast.Unary) (any, error) {
	return v.parenthesize(n.Operator.Lexeme, n.Right)
}

func (v astPrinter) VisitGrouping(n *ast.Grouping) (any, error) {
	return v.parenthesize("group", n.Expression)
}

func (v astPrinter) VisitLiteral(n *ast.Literal) (any, error) {
	switch val := n.Value.(type) {
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

func (v astPrinter) VisitVariable(n *ast.Variable) (any, error) {
	return n.Name.Lexeme, nil
}

func (v astPrinter) parenthesize(name string, expressions ...ast.Expr) (string, error) {
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
