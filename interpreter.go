package main

import (
	"fmt"
	"reflect"
	"strings"
)

type interpreter struct{}

// RuntimeError represents a runtime evaluation error tied to a token.
type RuntimeError struct {
	tok     *token
	message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("error at '%s': %s\n", e.tok.lexeme, e.message)
}

func (v *interpreter) Intepret(e expr) (any, error) {
	return v.evaluate(e)
}

func (v *interpreter) visitTernary(expr *ternary) (any, error) {
	val, err := v.evaluate(expr.condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(val) {
		return v.evaluate(expr.thenExpr)
	} else {
		return v.evaluate(expr.elseExpr)
	}
}

func (v *interpreter) visitBinary(expr *binary) (any, error) {
	left, err := v.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := v.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.tokenType {
	case MINUS:
		l, r, err := asTwoFloat64(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return l - r, nil

	case STAR:
		l, r, err := asTwoFloat64(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return l * r, nil

	case SLASH:
		l, r, err := asTwoFloat64(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		if r == 0 {
			return nil, &RuntimeError{
				tok:     &expr.operator,
				message: "division by zero",
			}
		}
		return l / r, nil

	case PLUS:
		// string concatenation only allowed if both operands are strings
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r, nil
			}
			return nil, &RuntimeError{
				tok:     &expr.operator,
				message: "operands must be two strings",
			}
		}

		// otherwise treat as numeric addition
		l, r, err := asTwoFloat64(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return l + r, nil

	case GREATER, GREATER_EQUAL, LESS, LESS_EQUAL:
		res, err := compareOperands(left, right, &expr.operator)
		if err != nil {
			return nil, err
		}

		switch expr.operator.tokenType {
		case GREATER:
			return res > 0, nil
		case GREATER_EQUAL:
			return res >= 0, nil
		case LESS:
			return res < 0, nil
		case LESS_EQUAL:
			return res <= 0, nil
		}

	case BANG_EQUAL:
		// TODO: currently using deepEqual. Maybe we limit to compare only string & number ?
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	// unreachable
	panic("unreachable")
}

func (v *interpreter) visitUnary(expr *unary) (any, error) {
	right, err := v.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.tokenType {
	case MINUS:
		n, err := asFloat64(&expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -n, nil
	case BANG:
		return !isTruthy(right), nil
	}

	// unreachable
	panic("unreachable")
}

func (v *interpreter) visitGrouping(expr *grouping) (any, error) {
	return v.evaluate(expr.expression)
}

func (v *interpreter) visitLiteral(expr *literal) (any, error) {
	return expr.value, nil
}

// evaluate dispatches AST node evaluation
func (v *interpreter) evaluate(e expr) (any, error) {
	return e.accept(v)
}

// toFloat64 converts supported numeric types into float64
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

// asFloat64 validates and converts a single operand (used for unary ops)
func asFloat64(operator *token, operand any) (float64, *RuntimeError) {
	n, ok := toFloat64(operand)
	if !ok {
		return 0, &RuntimeError{
			tok:     operator,
			message: "operand must be a number",
		}
	}
	return n, nil
}

// asTwoFloat64 validates and converts two operands (used for binary math ops)
func asTwoFloat64(op *token, left, right any) (float64, float64, *RuntimeError) {
	l, lok := toFloat64(left)
	r, rok := toFloat64(right)

	if !lok || !rok {
		return 0, 0, &RuntimeError{
			tok:     op,
			message: "operands must be two numbers",
		}
	}

	return l, r, nil
}

// isTruthy defines language truthiness rules:
// false values: nil, false, 0, ""
func isTruthy(e any) bool {
	switch v := e.(type) {
	case nil:
		return false
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case int64:
		return v != 0
	case string:
		return v != ""
	default:
		return true
	}
}

// compareOperands compares two values if both are numbers or both are strings
// returns: -1 (a < b), 0 (a == b), 1 (a > b)
func compareOperands(a, b any, operator *token) (int, *RuntimeError) {
	defaultErrMsg := "operands must be two strings or numbers"

	// string compare
	if l, ok := a.(string); ok {
		if r, ok := b.(string); ok {
			return strings.Compare(l, r), nil
		}

		return 0, &RuntimeError{
			tok:     operator,
			message: defaultErrMsg,
		}
	}

	// number compare
	l, r, err := asTwoFloat64(operator, a, b)
	if err != nil {
		err.message = defaultErrMsg // override err message for clarity
		return 0, err
	}

	switch {
	case l < r:
		return -1, nil
	case l > r:
		return 1, nil
	default:
		return 0, nil
	}
}

// isEqual checks equality with numeric normalization + deep fallback
func isEqual(a any, b any) bool {
	// nil handling
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// numbers
	if la, ok := toFloat64(a); ok {
		if lb, ok := toFloat64(b); ok {
			return la == lb
		}
	}

	// fallback
	return reflect.DeepEqual(a, b)
}
