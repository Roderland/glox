package main

type (
	stmt interface {
		exec(scopeList)
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
)
