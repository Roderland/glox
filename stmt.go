package main

type (
	stmt interface {
		exec(*interpreter)
	}
	exprStmt struct {
		body expr
	}
	printStmt struct {
		body expr
	}
	varStmt struct {
		name        token
		initializer expr
	}
	blockStmt struct {
		stmts []stmt
	}
	ifStmt struct {
		condition  expr
		thenBranch stmt
		elseBranch stmt
	}
	whileStmt struct {
		condition expr
		body      stmt
		increment expr
	}
	breakStmt    struct{}
	continueStmt struct{}
	functionStmt struct {
		name   token
		params []token
		body   []stmt
	}
	returnStmt struct {
		keyword token
		value   expr
	}
)
