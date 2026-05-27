package runtime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nibtr/gox/ast"
	"github.com/nibtr/gox/lexer"
)

const (
	divisionByZeroErrMsg           = "division by zero"
	operandsMustBeTwoNumbersErrMsg = "operands must be two numbers"
	operandMustBeNumberErrMsg      = "operand must be a number"
	operandsMustBeStrOrNumErrMsg   = "operands must be two strings or numbers"
)

type interpreter struct {
	environment *Environment
}

func NewInterpreter() *interpreter {
	return &interpreter{
		environment: NewEnvironment(),
	}
}

// RuntimeError represents a runtime evaluation error tied to a token.
type RuntimeError struct {
	Token   *lexer.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] error at '%s': %s\n", e.Token.Line, e.Token.Lexeme, e.Message)
}

func (v *interpreter) Eval(expr ast.Expr) (any, error) {
	return v.evaluate(expr)
}

func (v *interpreter) Intepret(statements []ast.Stmt) error {
	for _, s := range statements {
		if err := v.execute(s); err != nil {
			return err
		}
	}
	return nil
}

// ------------ Expression section -------------------

func (v *interpreter) VisitAssignExpr(expr *ast.Assign) (any, error) {
	value, err := v.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	if err := v.environment.assign(expr.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (v *interpreter) VisitTernary(expr *ast.Ternary) (any, error) {
	val, err := v.evaluate(expr.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(val) {
		return v.evaluate(expr.ThenExpr)
	} else {
		return v.evaluate(expr.ElseExpr)
	}
}

func (v *interpreter) VisitLogical(expr *ast.Logical) (any, error) {
	left, err := v.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.TokenType == lexer.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return v.evaluate(expr.Right)
}

func (v *interpreter) VisitBinary(expr *ast.Binary) (any, error) {
	left, err := v.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := v.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case lexer.MINUS:
		l, r, err := asTwoFloat64(&expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l - r, nil

	case lexer.STAR:
		l, r, err := asTwoFloat64(&expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l * r, nil

	case lexer.SLASH:
		l, r, err := asTwoFloat64(&expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		if r == 0 {
			return nil, &RuntimeError{
				Token:   &expr.Operator,
				Message: divisionByZeroErrMsg,
			}
		}
		return l / r, nil

	case lexer.PLUS:
		// string concatenation only allowed if both operands are strings
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r, nil
			}
			return nil, &RuntimeError{
				Token:   &expr.Operator,
				Message: operandsMustBeStrOrNumErrMsg,
			}
		}

		// otherwise treat as numeric addition
		l, r, err := asTwoFloat64(&expr.Operator, left, right)
		if err != nil {
			// override error for clarity
			err.Message = operandsMustBeStrOrNumErrMsg
			return nil, err
		}
		return l + r, nil

	case lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL:
		res, err := compareOperands(left, right, &expr.Operator)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.TokenType {
		case lexer.GREATER:
			return res > 0, nil
		case lexer.GREATER_EQUAL:
			return res >= 0, nil
		case lexer.LESS:
			return res < 0, nil
		case lexer.LESS_EQUAL:
			return res <= 0, nil
		}

	case lexer.BANG_EQUAL:
		// TODO: currently using deepEqual. Maybe we limit to compare only string & number ?
		return !isEqual(left, right), nil
	case lexer.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	// unreachable
	panic("unreachable")
}

func (v *interpreter) VisitUnary(expr *ast.Unary) (any, error) {
	right, err := v.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.TokenType {
	case lexer.MINUS:
		n, err := asFloat64(&expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -n, nil
	case lexer.BANG:
		return !isTruthy(right), nil
	}

	// unreachable
	panic("unreachable")
}

func (v *interpreter) VisitGrouping(expr *ast.Grouping) (any, error) {
	return v.evaluate(expr.Expression)
}

func (v *interpreter) VisitLiteral(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (v *interpreter) VisitVariable(expr *ast.Variable) (any, error) {
	return v.environment.get(expr.Name)
}

// ----------- Statement section -------------------

func (v *interpreter) VisitVarStmt(stmt *ast.VarStmt) error {
	var value any
	if stmt.Initializer != nil {
		v, err := v.evaluate(stmt.Initializer)
		if err != nil {
			return err
		}
		value = v
	}

	v.environment.define(stmt.Name.Lexeme, value)
	return nil
}

func (v *interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) error {
	_, err := v.evaluate(stmt.Expression)
	return err
}

func (v *interpreter) VisitIfStmt(stmt *ast.IfStmt) error {
	cond, err := v.evaluate(stmt.Condition)
	if err != nil {
		return err
	}
	if isTruthy(cond) {
		return v.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return v.execute(stmt.ElseBranch)
	}

	return nil
}

func (v *interpreter) VisitWhileStmt(stmt *ast.WhileStmt) error {
	for {
		cond, err := v.evaluate(stmt.Condition)
		if err != nil {
			return err
		}

		if !isTruthy(cond) {
			break
		}

		if err := v.execute(stmt.Body); err != nil {
			return err
		}
	}

	return nil
}

func (v *interpreter) VisitPrintStmt(stmt *ast.PrintStmt) error {
	value, err := v.evaluate(stmt.Expression)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func (v *interpreter) VisitBlockStmt(stmt *ast.BlockStmt) error {
	return v.executeBlock(stmt.Statements, NewEnvironmentWithEnclosing(v.environment))
}

// ------------------- Helpers ---------------------

func (v *interpreter) execute(stmt ast.Stmt) error {
	return stmt.Accept(v)
}

// evaluate dispatches AST node evaluation
func (v *interpreter) evaluate(e ast.Expr) (any, error) {
	return e.Accept(v)
}

func (v *interpreter) executeBlock(stmts []ast.Stmt, env *Environment) error {
	// restore previous environment if error occurs
	previous := v.environment
	defer func() {
		v.environment = previous
	}()

	v.environment = env
	for _, stmt := range stmts {
		if err := v.execute(stmt); err != nil {
			return err
		}
	}

	return nil
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
func asFloat64(operator *lexer.Token, operand any) (float64, *RuntimeError) {
	n, ok := toFloat64(operand)
	if !ok {
		return 0, &RuntimeError{
			Token:   operator,
			Message: operandMustBeNumberErrMsg,
		}
	}
	return n, nil
}

// asTwoFloat64 validates and converts two operands (used for binary math ops)
func asTwoFloat64(op *lexer.Token, left, right any) (float64, float64, *RuntimeError) {
	l, lok := toFloat64(left)
	r, rok := toFloat64(right)

	if !lok || !rok {
		return 0, 0, &RuntimeError{
			Token:   op,
			Message: operandsMustBeTwoNumbersErrMsg,
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
func compareOperands(a, b any, operator *lexer.Token) (int, *RuntimeError) {
	// string compare
	if l, ok := a.(string); ok {
		if r, ok := b.(string); ok {
			return strings.Compare(l, r), nil
		}

		return 0, &RuntimeError{
			Token:   operator,
			Message: operandsMustBeStrOrNumErrMsg,
		}
	}

	// number compare
	l, r, err := asTwoFloat64(operator, a, b)
	if err != nil {
		err.Message = operandsMustBeStrOrNumErrMsg // override err message for clarity
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
