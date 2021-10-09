package main

import (
	"fmt"
	"reflect"
)

type jump uint8

const (
	J_NONE jump = iota
	J_BREAK
	J_CONTINUE
	J_RETURN
)

type interpreter struct {
	stmts []stmt
	env   scopeList
	jp    jump
}

func newInterpreter(stmts []stmt) *interpreter {
	env := make(scopeList, 1)
	env[0] = *newScope(nil)
	return &interpreter{stmts: stmts, env: env, jp: J_NONE}
}

func (i *interpreter) run() {
	for _, s := range i.stmts {
		s.exec(i)
	}
}

func (b *breakStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	itp.jp = J_BREAK
}

func (c *continueStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	itp.jp = J_CONTINUE
}

func (w *whileStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	for isTrue(w.condition.eval(itp)) {
		w.body.exec(itp)
		if itp.jp == J_BREAK {
			itp.jp = J_NONE
			break
		} else if itp.jp == J_CONTINUE {
			itp.jp = J_NONE
		}
		if w.increment != nil {
			w.increment.eval(itp)
		}
	}
}

func (i *ifStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	if isTrue(i.condition.eval(itp)) {
		i.thenBranch.exec(itp)
	} else if i.elseBranch != nil {
		i.elseBranch.exec(itp)
	}
}

func (b *blockStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	itp.env.push(*newScope(itp.env.peek()))
	for _, s := range b.stmts {
		s.exec(itp)
	}
	itp.env.pop()
}

func (v *varStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	var value interface{}
	if v.initializer != nil {
		value = v.initializer.eval(itp)
	}
	itp.env.peek().define(v.name.text, value)
}

func (p *printStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	value := p.body.eval(itp)
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
}

func (e *exprStmt) exec(itp *interpreter) {
	if itp.jp != J_NONE {
		return
	}
	e.body.eval(itp)
}

func (l *logical) eval(itp *interpreter) interface{} {
	left := l.left.eval(itp)
	if l.operator.ttype == OR {
		if isTrue(left) {
			return left
		}
	} else {
		if !isTrue(left) {
			return left
		}
	}
	return l.right.eval(itp)
}

func (a *assign) eval(itp *interpreter) interface{} {
	value := a.right.eval(itp)
	itp.env.peek().assign(a.name.text, value)
	return value
}

func (v *variable) eval(itp *interpreter) interface{} {
	return itp.env.peek().get(v.name.text)
}

func (b *binary) eval(itp *interpreter) interface{} {
	left := b.left.eval(itp)
	right := b.right.eval(itp)
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

func (u *unary) eval(itp *interpreter) interface{} {
	right := u.right.eval(itp)
	switch u.operator.ttype {
	case MINUS:
		checkNumber(u.operator, right)
		return -(right.(float64))
	case BANG:
		return !isTrue(right)
	}
	return nil
}

func (l *literal) eval(itp *interpreter) interface{} {
	return l.value
}

func (g *group) eval(itp *interpreter) interface{} {
	return g.body.eval(itp)
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
