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
	ifStmt struct {
		condition  expr
		thenBranch stmt
		elseBranch stmt
	}
	whileStmt struct {
		condition expr
		body      stmt
	}
)
