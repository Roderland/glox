package main

import (
	"fmt"
	"reflect"
)

type (
	Expr interface {
		eval(interpreter *Interpreter) interface{}
	}

	Literal struct {
		value interface{}
	}

	Unary struct {
		operator Token
		right    Expr
	}

	Binary struct {
		left     Expr
		operator Token
		right    Expr
	}

	Grouping struct {
		expression Expr
	}

	Variable struct {
		name Token
	}

	Assign struct {
		name  Token
		value Expr
	}

	Logical struct {
		left     Expr
		operator Token
		right    Expr
	}

	Call struct {
		callee Expr
		paren  Token
		args   []Expr
	}
)

func (l Literal) eval(interpreter *Interpreter) interface{} {
	return l.value
}

func (u Unary) eval(interpreter *Interpreter) interface{} {
	right := u.right.eval(interpreter)
	switch u.operator.tokenType {
	case MINUS:
		return -(right.(float64))
	case BANG:
		return !isTrue(right)
	default:
		return right
	}
}

func (b Binary) eval(interpreter *Interpreter) interface{} {
	left := b.left.eval(interpreter)
	right := b.right.eval(interpreter)
	switch b.operator.tokenType {
	case PLUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 {
			checkOperands(reflect.Float64, b.operator, right)
			return left.(float64) + right.(float64)
		} else {
			checkOperands(reflect.String, b.operator, left, right)
			return left.(string) + right.(string)
		}
	case MINUS:
		checkOperands(reflect.Float64, b.operator, left, right)
		return left.(float64) - right.(float64)
	case STAR:
		checkOperands(reflect.Float64, b.operator, left, right)
		return left.(float64) * right.(float64)
	case SLASH:
		checkOperands(reflect.Float64, b.operator, left, right)
		return left.(float64) / right.(float64)
	case GREATER:
		checkOperands(reflect.Float64, b.operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkOperands(reflect.Float64, b.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkOperands(reflect.Float64, b.operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkOperands(reflect.Float64, b.operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return left != right
	case EQUAL_EQUAL:
		return left == right
	}
	return nil
}

func (g Grouping) eval(interpreter *Interpreter) interface{} {
	return g.expression.eval(interpreter)
}

func (v Variable) eval(interpreter *Interpreter) interface{} {
	return interpreter.local.get(v.name)
}

func (a Assign) eval(interpreter *Interpreter) interface{} {
	value := a.value.eval(interpreter)
	interpreter.local.assign(a.name, value)
	return value
}

func (l Logical) eval(interpreter *Interpreter) interface{} {
	left := l.left.eval(interpreter)
	if l.operator.tokenType == OR {
		if isTrue(left) {
			return left
		}
	} else {
		if !isTrue(left) {
			return left
		}
	}
	return l.right.eval(interpreter)
}

func (c Call) eval(interpreter *Interpreter) interface{} {
	callee := c.callee.eval(interpreter)

	args := make([]interface{}, len(c.args))
	for i, arg := range c.args {
		args[i] = arg.eval(interpreter)
	}

	fun, ok := callee.(Function)

	if !ok {
		exitWithErr(c.paren.line, "Can only call functions")
	}
	if fun.arity() != len(args) {
		exitWithErr(c.paren.line, fmt.Sprintf("Expect %d arguments but get %d", fun.arity(), len(args)))
	}

	return fun.call(interpreter, args)
}
