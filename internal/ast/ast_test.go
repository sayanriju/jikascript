package ast

import "testing"

// TestNodeInterfaces keeps the statement and expression marker methods in
// place. It is a compile-time check: if a node stops implementing Stmt or Expr,
// this file fails to build.
func TestNodeInterfaces(t *testing.T) {
	statements := []Stmt{
		&VarDecl{}, &ConstDecl{}, &IncludeStmt{}, &FunctionDecl{}, &ReturnStmt{},
		&ExprStmt{}, &AssignStmt{}, &IfStmt{}, &WhileStmt{}, &DoWhileStmt{},
		&ForStmt{}, &SwitchStmt{}, &BreakStmt{}, &ContinueStmt{}, &TaskDecl{},
		&BlockStmt{},
	}
	expressions := []Expr{
		&ArrowFunctionExpr{}, &Identifier{}, &NumberLiteral{}, &StringLiteral{},
		&BoolLiteral{}, &CallExpr{}, &SelectorExpr{}, &IndexExpr{},
		&BinaryExpr{}, &UnaryExpr{}, &ArrayLiteralExpr{}, &ObjectLiteralExpr{},
	}
	if len(statements) == 0 || len(expressions) == 0 {
		t.Fatal("node interface lists must not be empty")
	}
}

func TestUnaryExprPreservesPostfixAndPosition(t *testing.T) {
	pos := Position{Line: 4, Column: 9}
	expr := &UnaryExpr{Op: "++", Expr: &Identifier{Name: "i"}, Postfix: true, Pos: pos}
	if !expr.Postfix {
		t.Fatal("expected Postfix to be true")
	}
	if expr.Op != "++" {
		t.Fatalf("Op = %q, want ++", expr.Op)
	}
	if expr.Pos != pos {
		t.Fatalf("Pos = %+v, want %+v", expr.Pos, pos)
	}
}
