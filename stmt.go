package main

import "fmt"

type (
	Stmt interface {
		exec(interpreter *Interpreter)
		//resolve()
	}
	exprStmt struct {
		expr Expr
	}
	printStmt struct {
		expr Expr
	}
	varStmt struct {
		name        Token
		initializer Expr
	}
	blockStmt struct {
		stmts []Stmt
	}
	ifStmt struct {
		condition  Expr
		thenBranch Stmt
		elseBranch Stmt
	}
	whileStmt struct {
		condition Expr
		body      Stmt
	}
	functionStmt struct {
		name   Token
		params []Token
		stmts  []Stmt
	}
	returnStmt struct {
		keyword Token
		value   Expr
	}
)

//func (v varStmt) resolve() {
//	declare(v.name)
//	if v.initializer != nil {
//		v.initializer.resolve()
//	}
//	define(v.name)
//}

func (e exprStmt) exec(interpreter *Interpreter) {
	e.expr.eval(interpreter)
}

func (p printStmt) exec(interpreter *Interpreter) {
	value := p.expr.eval(interpreter)
	fmt.Println(toString(value))
}

func (v varStmt) exec(interpreter *Interpreter) {
	var value interface{}
	if v.initializer != nil {
		value = v.initializer.eval(interpreter)
	}
	interpreter.envir.define(v.name.lexeme, value)
}

func (b blockStmt) exec(interpreter *Interpreter) {
	father := interpreter.envir
	child := &Envir{
		enclosing: father,
		values:    map[string]interface{}{},
	}
	interpreter.enterScope(child)
	defer interpreter.enterScope(father)
	for _, stmt := range b.stmts {
		stmt.exec(interpreter)
	}
}

func (i ifStmt) exec(interpreter *Interpreter) {
	if isTrue(i.condition.eval(interpreter)) {
		i.thenBranch.exec(interpreter)
	} else {
		if i.elseBranch != nil {
			i.elseBranch.exec(interpreter)
		}
	}
}

func (w whileStmt) exec(interpreter *Interpreter) {
	for isTrue(w.condition.eval(interpreter)) {
		w.body.exec(interpreter)
	}
}

func (f functionStmt) exec(interpreter *Interpreter) {
	fun := Function{f, interpreter.envir}
	interpreter.envir.define(f.name.lexeme, fun)
}

func (r returnStmt) exec(interpreter *Interpreter) {
	var result interface{}
	if r.value != nil {
		result = r.value.eval(interpreter)
	}
	interpreter.retVal = append(interpreter.retVal, result)
}
