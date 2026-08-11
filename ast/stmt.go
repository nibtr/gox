package ast

import (
	"github.com/nibtr/gox/lexer"
)

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitFunctionStmt(stmt *FunctionStmt) error
	VisitReturnStmt(stmt *ReturnStmt) error
	VisitExpressionStmt(stmt *ExpressionStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
	VisitBreakStmt(stmt *BreakStmt) error
	VisitContinueStmt(stmt *ContinueStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitVarStmt(stmt *VarStmt) error
	VisitClassStmt(stmt *ClassStmt) error
}

type FunctionStmt struct {
	Name   lexer.Token
	Params []lexer.Token
	Body   []Stmt
}

type ReturnStmt struct {
	Keyword lexer.Token
	Value   Expr
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
	// Increment runs after Body on every iteration, including when Body
	// exits via continue. Only set for desugared for-loops.
	Increment Expr
}

type BreakStmt struct{}

type ContinueStmt struct{}

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

type ClassStmt struct {
	Name    lexer.Token
	Methods []FunctionStmt
}

func (s *FunctionStmt) Accept(v StmtVisitor) error {
	return v.VisitFunctionStmt(s)
}

func (s *ReturnStmt) Accept(v StmtVisitor) error {
	return v.VisitReturnStmt(s)
}

func (s *IfStmt) Accept(v StmtVisitor) error {
	return v.VisitIfStmt(s)
}

func (s *WhileStmt) Accept(v StmtVisitor) error {
	return v.VisitWhileStmt(s)
}

func (s *BreakStmt) Accept(v StmtVisitor) error {
	return v.VisitBreakStmt(s)
}

func (s *ContinueStmt) Accept(v StmtVisitor) error {
	return v.VisitContinueStmt(s)
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

func (s *ClassStmt) Accept(v StmtVisitor) error {
	return v.VisitClassStmt(s)
}
