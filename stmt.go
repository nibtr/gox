package main

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	visitExpressionStmt(stmt *ExpressionStmt) (any, error)
	visitPrintStmt(stmt *PrintStmt) error
	visitVarStmt(stmt *VarStmt) error
}

type ExpressionStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitExpressionStmt(s)
}

func (s *PrintStmt) Accept(v StmtVisitor) (any, error) {
	err := v.visitPrintStmt(s)
	return nil, err
}

func (s *VarStmt) Accept(v StmtVisitor) (any, error) {
	err := v.visitVarStmt(s)
	return nil, err
}
