package main

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	visitExpressionStmt(stmt *ExpressionStmt) (any, error)
	visitPrintStmt(stmt *PrintStmt) error
}

type ExpressionStmt struct {
	Expression expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitExpressionStmt(s)
}

type PrintStmt struct {
	Expression expr
}

func (s *PrintStmt) Accept(v StmtVisitor) (any, error) {
	err := v.visitPrintStmt(s)
	return nil, err
}
