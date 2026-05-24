package ast

import (
	"github.com/nibtr/gox/lexer"
)

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) (any, error)
	VisitPrintStmt(stmt *PrintStmt) error
	VisitVarStmt(stmt *VarStmt) error
}

type ExpressionStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type VarStmt struct {
	Name        lexer.Token
	Initializer Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitExpressionStmt(s)
}

func (s *PrintStmt) Accept(v StmtVisitor) (any, error) {
	err := v.VisitPrintStmt(s)
	return nil, err
}

func (s *VarStmt) Accept(v StmtVisitor) (any, error) {
	err := v.VisitVarStmt(s)
	return nil, err
}
