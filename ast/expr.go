package ast

import (
	"github.com/nibtr/gox/lexer"
)

type ExprVisitor interface {
	// VisitTernary evaluates condition ? thenExpr : elseExpr
	VisitTernary(n *Ternary) (any, error)
	// VisitBinary evaluates binary expressions like +, -, *, /, comparisons
	VisitBinary(n *Binary) (any, error)
	// VisitUnary evaluates unary expressions like -x and !x
	VisitUnary(n *Unary) (any, error)
	// VisitGrouping evaluates expressions inside parentheses (...)
	VisitGrouping(n *Grouping) (any, error)
	// VisitLiteral returns a literal value directly
	VisitLiteral(n *Literal) (any, error)
	// VisitVariable evaluates variable expressions (identifier lookup)
	VisitVariable(n *Variable) (any, error)
}

type Expr interface {
	Accept(v ExprVisitor) (any, error)
}

type Ternary struct {
	Condition Expr
	ThenExpr  Expr
	ElseExpr  Expr
}

type Binary struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

type Unary struct {
	Operator lexer.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Variable struct {
	Name lexer.Token
}

func (n *Ternary) Accept(v ExprVisitor) (any, error) {
	return v.VisitTernary(n)
}

func (n *Binary) Accept(v ExprVisitor) (any, error) {
	return v.VisitBinary(n)
}

func (n *Unary) Accept(v ExprVisitor) (any, error) {
	return v.VisitUnary(n)
}

func (n *Grouping) Accept(v ExprVisitor) (any, error) {
	return v.VisitGrouping(n)
}

func (n *Literal) Accept(v ExprVisitor) (any, error) {
	return v.VisitLiteral(n)
}

func (n *Variable) Accept(v ExprVisitor) (any, error) {
	return v.VisitVariable(n)
}
