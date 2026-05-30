package ast

import (
	"github.com/nibtr/gox/lexer"
)

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitFunctionStmt(stmt *FunctionStmt) error
	VisitExpressionStmt(stmt *ExpressionStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitVarStmt(stmt *VarStmt) error
}

type FunctionStmt struct {
	Name   lexer.Token
	Params []lexer.Token
	Body   []Stmt
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type ExpressionStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type BlockStmt struct {
	Statements []Stmt
}

type VarStmt struct {
	Name        lexer.Token
	Initializer Expr
}

func (s *FunctionStmt) Accept(v StmtVisitor) error {
	return v.VisitFunctionStmt(s)
}

func (s *IfStmt) Accept(v StmtVisitor) error {
	return v.VisitIfStmt(s)
}

func (s *WhileStmt) Accept(v StmtVisitor) error {
	return v.VisitWhileStmt(s)
}

func (s *ExpressionStmt) Accept(v StmtVisitor) error {
	return v.VisitExpressionStmt(s)
}

func (s *PrintStmt) Accept(v StmtVisitor) error {
	return v.VisitPrintStmt(s)
}

func (s *BlockStmt) Accept(v StmtVisitor) error {
	return v.VisitBlockStmt(s)
}

func (s *VarStmt) Accept(v StmtVisitor) error {
	return v.VisitVarStmt(s)
}
