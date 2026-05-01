package main

import (
	"fmt"
	"reflect"
	"strings"
)

type interpreter struct{}

type runtimeError struct {
	tok     *token
	message string
}

func (e *runtimeError) Error() string {
	return fmt.Sprintf("error at '%s': %s", e.tok.lexeme, e.message)
}

func (v *interpreter) visitTernary(expr *ternary) (any, error) {
	if isTruthy(v.evaluate(expr.condition)) {
		return v.evaluate(expr.thenExpr)
	} else {
		return v.evaluate(expr.elseExpr)
	}
}

func (v *interpreter) visitBinary(expr *binary) (any, error) {
	left := v.evaluate(expr.left)
	right := v.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		return mustFloat64(left) - mustFloat64(right)
	case STAR:
		return mustFloat64(left) * mustFloat64(right)
	case SLASH:
		return mustFloat64(left) / mustFloat64(right)
	case PLUS:
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l + r
			}
		}

		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		panic("operands must be two numbers or two strings")

	case GREATER:
		if res, ok := compareValues(left, right); ok {
			return res > 0
		}
		panic("operands must be two numbers or two strings")
	case GREATER_EQUAL:
		if res, ok := compareValues(left, right); ok {
			return res >= 0
		}
		panic("operands must be two numbers or two strings")
	case LESS:
		if res, ok := compareValues(left, right); ok {
			return res < 0
		}
		panic("operands must be two numbers or two strings")
	case LESS_EQUAL:
		if res, ok := compareValues(left, right); ok {
			return res <= 0
		}
		panic("operands must be two numbers or two strings")
	case BANG_EQUAL:
		// TODO: currently using deepEqual. Maybe we limit to compare only string & number ?
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	}

	// unreachable
	return nil
}

func (v *interpreter) visitUnary(expr *unary) (any, error) {
	right := v.evaluate(expr.right)
	switch expr.operator.tokenType {
	case MINUS:
		// TODO: check this, we don't know the type of right, so for now just cast it to float64
		return -mustFloat64(right)
	case BANG:
		return !isTruthy(right)
	}

	// unreachable
	return nil
}

func (v *interpreter) visitGrouping(expr *grouping) (any, error) {
	return v.evaluate(expr.expression)
}

func (v *interpreter) visitLiteral(expr *literal) (any, error) {
	return expr.value
}

func (v *interpreter) evaluate(e expr) (any, error) {
	return e.accept(v)
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func mustFloat64(v any) float64 {
	n, ok := toFloat64(v)
	if !ok {
		// TODO: check this, for now can just panic
		panic("operand must be a number")
	}
	return n
}

func isTruthy(e any) bool {
	if e == nil {
		return false
	}
	if b, ok := e.(bool); ok {
		return b
	}

	return false
}

func compareValues(a, b any) (int, bool) {
	// string compare
	if l, ok := a.(string); ok {
		if r, ok := b.(string); ok {
			return strings.Compare(l, r), true
		}
	}

	// number compare
	if l, ok := toFloat64(a); ok {
		if r, ok := toFloat64(b); ok {
			switch {
			case l < r:
				return -1, true
			case l > r:
				return 1, true
			default:
				return 0, true
			}
		}
	}

	return 0, false
}

func isEqual(a any, b any) bool {
	return reflect.DeepEqual(a, b)
}
