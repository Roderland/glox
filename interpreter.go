package main

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	globals *Envir
	envir   *Envir
	retVal  []interface{}
}

func _Interpreter() *Interpreter {
	globals := &Envir{nil, map[string]interface{}{}}
	return &Interpreter{globals: globals, envir: globals, retVal: []interface{}{}}
}

func (interpreter *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		statement.exec(interpreter)
	}
}

//func (interpreter *Interpreter) execute(stmt Stmt) {
//	stmt.exec(interpreter)
//}

//func (interpreter *Interpreter) evaluate(expr Expr) interface{} {
//	return expr.eval(interpreter)
//}

func (interpreter *Interpreter) enterScope(target *Envir) {
	interpreter.envir = target
}

func checkOperands(kind reflect.Kind, operator Token, operands ...interface{}) {
	for _, operand := range operands {
		if reflect.TypeOf(operand).Kind() != kind {
			exitWithErr(operator.line, fmt.Sprintf("Operator '%s' expect right operands.", operator.lexeme))
		}
	}
}

func isTrue(obj interface{}) bool {
	if obj == nil || obj == false {
		return false
	}
	return true
}

func toString(obj interface{}) string {
	if obj == nil {
		return "nil"
	}
	if reflect.TypeOf(obj).Kind() == reflect.TypeOf(Function{}).Kind() {
		return "<fun $" + obj.(Function).declaration.name.lexeme + ">"
	}
	return fmt.Sprint(obj)
}
