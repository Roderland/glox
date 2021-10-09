package main

import (
	"fmt"
	"reflect"
)

type interpreter struct {
	stmts []stmt
	env   scopeList
}

func newInterpreter(stmts []stmt) *interpreter {
	env := make(scopeList, 1)
	env[0] = *newScope(nil)
	return &interpreter{stmts: stmts, env: env}
}

func (i *interpreter) run() {
	for _, s := range i.stmts {
		s.exec(i.env)
	}
}

func (w *whileStmt) exec(scopes scopeList) {
	for isTrue(w.condition.eval(scopes)) {
		w.body.exec(scopes)
	}
}

func (i *ifStmt) exec(scope scopeList) {
	if isTrue(i.condition.eval(scope)) {
		i.thenBranch.exec(scope)
	} else if i.elseBranch != nil {
		i.elseBranch.exec(scope)
	}
}

func (b *blockStmt) exec(scopes scopeList) {
	scopes.push(*newScope(scopes.peek()))
	for _, s := range b.stmts {
		s.exec(scopes)
	}
	scopes.pop()
}

func (v *varStmt) exec(scopes scopeList) {
	var value interface{}
	if v.initializer != nil {
		value = v.initializer.eval(scopes)
	}
	scopes.peek().define(v.name.text, value)
}

func (p *printStmt) exec(scopes scopeList) {
	value := p.body.eval(scopes)
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
}

func (e *exprStmt) exec(scopes scopeList) {
	e.body.eval(scopes)
}

func (l *logical) eval(scopes scopeList) interface{} {
	left := l.left.eval(scopes)
	if l.operator.ttype == OR {
		if isTrue(left) {
			return left
		}
	} else {
		if !isTrue(left) {
			return left
		}
	}
	return l.right.eval(scopes)
}

func (a *assign) eval(scopes scopeList) interface{} {
	value := a.right.eval(scopes)
	scopes.peek().assign(a.name.text, value)
	return value
}

func (v *variable) eval(scopes scopeList) interface{} {
	return scopes.peek().get(v.name.text)
}

func (b *binary) eval(scopes scopeList) interface{} {
	left := b.left.eval(scopes)
	right := b.right.eval(scopes)
	switch b.operator.ttype {
	case MINUS:
		checkNumber(b.operator, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		checkNumber(b.operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		checkNumber(b.operator, left, right)
		return left.(float64) * right.(float64)
	case PLUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}
		if reflect.TypeOf(left).Kind() == reflect.String && reflect.TypeOf(right).Kind() == reflect.String {
			return left.(string) + right.(string)
		}
	case GREATER:
		checkNumber(b.operator, left, right)
		return left.(float64) > right.(float64)
	case GEQUAL:
		checkNumber(b.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumber(b.operator, left, right)
		return left.(float64) < right.(float64)
	case LEQUAL:
		checkNumber(b.operator, left, right)
		return left.(float64) <= right.(float64)
	case BEQUAL:
		return left != right
	case EEQUAL:
		return left == right
	}
	return nil
}

func (u *unary) eval(scopes scopeList) interface{} {
	right := u.right.eval(scopes)
	switch u.operator.ttype {
	case MINUS:
		checkNumber(u.operator, right)
		return -(right.(float64))
	case BANG:
		return !isTrue(right)
	}
	return nil
}

func (l *literal) eval(scopes scopeList) interface{} {
	return l.value
}

func (g *group) eval(scopes scopeList) interface{} {
	return g.body.eval(scopes)
}

func checkNumber(operator token, objects ...interface{}) {
	for _, o := range objects {
		if reflect.TypeOf(o).Kind() != reflect.Float64 {
			exitWithErr("[ line %d ] Operator '%s' expect number as operands", operator.line, operator.text)
		}
	}
}

func isTrue(o interface{}) bool {
	if o == nil || o == false {
		return false
	}
	return true
}
