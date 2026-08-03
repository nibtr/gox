package resolver

import (
	"slices"

	"github.com/nibtr/gox/ast"
	"github.com/nibtr/gox/lexer"
	"github.com/nibtr/gox/parser"
	"github.com/nibtr/gox/runtime"
)

// Resolver is a struct that performs sematic-analysis in a new
// single pass over the tree to resolve all of the variables it contains.
//
// Implements the ExprVisitor and StmtVisitor and sits after parse step
// and before interpret step.
type Resolver struct {
	Interpreter *runtime.Interpreter
	// The scope stack used for local block scopes.
	Scopes []scope
}

// scope behaves like a linked list - the chain of Environment objects.
// The value associated with a key in the scope map represents
// whether or not we have finished resolving that variable’s initializer
type scope map[string]bool

func NewResolver(i *runtime.Interpreter) *Resolver {
	return &Resolver{
		Interpreter: i,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) error {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) error {
	r.declare(&stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitVariable(expr *ast.Variable) (any, error) {
	if len(r.Scopes) != 0 && !r.Scopes[0][expr.Name.Lexeme] {
		return nil, &parser.ParseError{
			Token:   &expr.Name,
			Message: "Can't read local variable in its own initializer.",
		}
	}

	r.resolveLocal(expr, &expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.Assign) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, &expr.Name)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) error {
	r.declare(&stmt.Name)
	r.define(&stmt.Name)
	r.resolveFunction(stmt)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) error {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) error {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) error {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) error {
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) error {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	if stmt.Increment != nil {
		r.resolveExpr(stmt.Increment)
	}

	return nil
}

func (r *Resolver) VisitBreakStmt(stmt *ast.BreakStmt) error       { return nil }
func (r *Resolver) VisitContinueStmt(stmt *ast.ContinueStmt) error { return nil }

func (r *Resolver) VisitBinary(expr *ast.Binary) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitCall(expr *ast.Call) (any, error) {
	r.resolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}
	return nil, nil
}

func (r *Resolver) VisitGrouping(expr *ast.Grouping) (any, error) {
	r.resolveExpr(expr.Expression)
	return nil, nil
}

func (r *Resolver) VisitLiteral(expr *ast.Literal) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogical(expr *ast.Logical) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitUnary(expr *ast.Unary) (any, error) {
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitTernary(expr *ast.Ternary) (any, error) {
	r.resolveExpr(expr.Condition)
	r.resolveExpr(expr.ThenExpr)
	if expr.ElseExpr != nil {
		r.resolveExpr(expr.ElseExpr)
	}
	return nil, nil
}

// ------- Helpers ---------

func (r *Resolver) resolveStmts(stmts []ast.Stmt) {
	for _, s := range stmts {
		r.resolveStmt(s)
	}
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.Scopes = append(r.Scopes, make(scope))
}

func (r *Resolver) endScope() {
	r.Scopes = r.Scopes[:len(r.Scopes)-1] // pop
}

func (r *Resolver) declare(name *lexer.Token) {
	if len(r.Scopes) == 0 {
		return
	}

	scope := r.Scopes[0] // peek
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *lexer.Token) {
	if len(r.Scopes) == 0 {
		return
	}

	r.Scopes[0][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name *lexer.Token) {
	for i := range slices.Backward(r.Scopes) {
		if _, ok := r.Scopes[i][name.Lexeme]; ok {
			r.Interpreter.Resolve(expr, len(r.Scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(f *ast.FunctionStmt) {
	r.beginScope()
	for _, param := range f.Params {
		r.declare(&param)
		r.define(&param)
	}
	r.resolveStmts(f.Body)
	r.endScope()
}
