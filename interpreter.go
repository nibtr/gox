package main

type interpreter struct{}

func (v *interpreter) visitTernary(expr *ternary) any {
	return ""
}

func (v *interpreter) visitBinary(expr *binary) any {
	left := v.evaluate(expr.left)
	right := v.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		return mustFloat64(left) - mustFloat64(right)
	case STAR:
		return mustFloat64(left) * mustFloat64(right)
	case SLASH:
		return mustFloat64(left) / mustFloat64(right)
		// TODO: case PLUS
	}

	// unreachable
	return nil
}

func (v *interpreter) visitUnary(expr *unary) any {
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

func (v *interpreter) visitGrouping(expr *grouping) any {
	return v.evaluate(expr.expression)
}

func (v *interpreter) visitLiteral(expr *literal) any {
	return expr.value
}

func (v *interpreter) evaluate(e expr) any {
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
