package resolver

import (
	"github.com/nibtr/gox/ast"
	"github.com/nibtr/gox/lexer"
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
