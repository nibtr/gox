package resolver

import (
	"slices"

	"github.com/nibtr/gox/ast"
	"github.com/nibtr/gox/lexer"
	"github.com/nibtr/gox/parser"
	"github.com/nibtr/gox/runtime"
)

type FunctionType int

const (
	None FunctionType = iota
	Function
)

// Resolver is a struct that performs sematic-analysis in a new
// single pass over the tree to resolve all of the variables it contains.
//
// Implements the ExprVisitor and StmtVisitor and sits after parse step
// and before interpret step.
type Resolver struct {
	Interpreter *runtime.Interpreter
	// The scope stack used for local block scopes.
	Scopes          []scope
	CurrentFunction FunctionType
}

// scope behaves like a linked list - the chain of Environment objects.
// The value associated with a key in the scope map represents
// whether or not we have finished resolving that variable’s initializer
type scope map[string]bool

func NewResolver(i *runtime.Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     i,
		Scopes:          make([]scope, 0),
		CurrentFunction: None,
	}
}

// ---------- Statements -----------

func (r *Resolver) VisitClassStmt(stmt *ast.ClassStmt) error {
	r.declare(&stmt.Name)
	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) error {
	r.beginScope()
	defer r.endScope()
	return r.ResolveStmts(stmt.Statements)
}

func (r *Resolver) VisitVarStmt(stmt *ast.VarStmt) error {
	if err := r.declare(&stmt.Name); err != nil {
		return err
	}
	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return err
		}
	}
	r.define(&stmt.Name)
	return nil
}

func (r *Resolver) VisitVariable(expr *ast.Variable) (any, error) {
	if len(r.Scopes) != 0 {
		if defined, ok := r.Scopes[len(r.Scopes)-1][expr.Name.Lexeme]; ok && !defined {
			return nil, &parser.ParseError{
				Token:   &expr.Name,
				Message: "Can't read local variable in its own initializer.",
			}
		}
	}

	r.resolveLocal(expr, &expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.Assign) (any, error) {
	if err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}

	r.resolveLocal(expr, &expr.Name)
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.FunctionStmt) error {
	if err := r.declare(&stmt.Name); err != nil {
		return err
	}
	r.define(&stmt.Name)
	if err := r.resolveFunction(stmt, Function); err != nil {
		return err
	}
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.ExpressionStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) error {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return err
	}

	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return err
	}

	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) error {
	if r.CurrentFunction == None {
		return &parser.ParseError{
			Token:   &stmt.Keyword,
			Message: "Can't return from top-level code.",
		}
	}
	if stmt.Value != nil {
		return r.resolveExpr(stmt.Value)
	}

	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) error {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return err
	}

	if err := r.resolveStmt(stmt.Body); err != nil {
		return err
	}

	if stmt.Increment != nil {
		if err := r.resolveExpr(stmt.Increment); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitBreakStmt(stmt *ast.BreakStmt) error       { return nil }
func (r *Resolver) VisitContinueStmt(stmt *ast.ContinueStmt) error { return nil }

// ------------ Expressions ----------------

func (r *Resolver) VisitBinary(expr *ast.Binary) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}

	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitCall(expr *ast.Call) (any, error) {
	if err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}

	for _, arg := range expr.Arguments {
		if err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGrouping(expr *ast.Grouping) (any, error) {
	return nil, r.resolveExpr(expr.Expression)
}

func (r *Resolver) VisitLiteral(expr *ast.Literal) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogical(expr *ast.Logical) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}

	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitUnary(expr *ast.Unary) (any, error) {
	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitTernary(expr *ast.Ternary) (any, error) {
	if err := r.resolveExpr(expr.Condition); err != nil {
		return nil, err
	}

	if err := r.resolveExpr(expr.ThenExpr); err != nil {
		return nil, err
	}

	if expr.ElseExpr != nil {
		if err := r.resolveExpr(expr.ElseExpr); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// ------- Helpers ---------

func (r *Resolver) ResolveStmts(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) error {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) beginScope() {
	r.Scopes = append(r.Scopes, make(scope))
}

func (r *Resolver) endScope() {
	r.Scopes = r.Scopes[:len(r.Scopes)-1] // pop
}

func (r *Resolver) declare(name *lexer.Token) error {
	if len(r.Scopes) == 0 {
		return nil
	}

	scope := r.Scopes[len(r.Scopes)-1] // peek
	if _, ok := scope[name.Lexeme]; ok {
		return &parser.ParseError{
			Token:   name,
			Message: "Already variable with this name in this scope.",
		}
	}

	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name *lexer.Token) {
	if len(r.Scopes) == 0 {
		return
	}

	r.Scopes[len(r.Scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name *lexer.Token) {
	for i := range slices.Backward(r.Scopes) {
		if _, ok := r.Scopes[i][name.Lexeme]; ok {
			r.Interpreter.Resolve(expr, len(r.Scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(f *ast.FunctionStmt, funcType FunctionType) error {
	enclosingFunc := r.CurrentFunction
	r.CurrentFunction = funcType

	r.beginScope()
	defer func() {
		r.endScope()
		r.CurrentFunction = enclosingFunc
	}()

	for _, param := range f.Params {
		if err := r.declare(&param); err != nil {
			return err
		}
		r.define(&param)
	}

	return r.ResolveStmts(f.Body)
}
