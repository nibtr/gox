package main

type ExprVisitor interface {
	// visitTernary evaluates condition ? thenExpr : elseExpr
	visitTernary(n *Ternary) (any, error)
	// visitBinary evaluates binary expressions like +, -, *, /, comparisons
	visitBinary(n *Binary) (any, error)
	// visitUnary evaluates unary expressions like -x and !x
	visitUnary(n *Unary) (any, error)
	// visitGrouping evaluates expressions inside parentheses (...)
	visitGrouping(n *Grouping) (any, error)
	// visitLiteral returns a literal value directly
	visitLiteral(n *Literal) (any, error)
	// visitVariable evaluates variable declaration
	visitVariable(n *Variable) (any, error)
}

type Expr interface {
	Accept(v ExprVisitor) (any, error)
}

type Ternary struct {
	condition Expr
	thenExpr  Expr
	elseExpr  Expr
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

type Unary struct {
	operator Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

type Variable struct {
	name Token
}

func (n *Ternary) Accept(v ExprVisitor) (any, error) {
	return v.visitTernary(n)
}

func (n *Binary) Accept(v ExprVisitor) (any, error) {
	return v.visitBinary(n)
}

func (n *Unary) Accept(v ExprVisitor) (any, error) {
	return v.visitUnary(n)
}

func (n *Grouping) Accept(v ExprVisitor) (any, error) {
	return v.visitGrouping(n)
}

func (n *Literal) Accept(v ExprVisitor) (any, error) {
	return v.visitLiteral(n)
}

func (n *Variable) Accept(v ExprVisitor) (any, error) {
	return v.visitVariable(n)
}
