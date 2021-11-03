package main

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	// 全局变量表
	global Table
	// 当前作用域变量表
	local *Table
	// 函数调用返回值保存栈
	returnStack []interface{}
}

func _Interpreter() *Interpreter {
	global := Table{nil, map[string]interface{}{}}
	return &Interpreter{
		global:      global,
		local:       &global,
		returnStack: []interface{}{},
	}
}

// 解释器执行所有语句
func (interpreter *Interpreter) interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		stmt.exec(interpreter)
	}
}

// 进入或退出作用域
func (interpreter *Interpreter) enterScope(target *Table) {
	interpreter.local = target
}

// 检查所有操作数的类型是否正确
func checkOperands(kind reflect.Kind, operator Token, operands ...interface{}) {
	for _, operand := range operands {
		if reflect.TypeOf(operand).Kind() != kind {
			exitWithErr(operator.line, "Operator '"+operator.lexeme+"' expect right operands.")
		}
	}
}

// 真值判断
func isTrue(obj interface{}) bool {
	if obj == nil || obj == false {
		return false
	}
	return true
}

// 获得任意类型对应的字符串表示
func toString(obj interface{}) string {
	if obj == nil {
		return "nil"
	}
	if reflect.TypeOf(obj).Kind() == reflect.TypeOf(Function{}).Kind() {
		return "<fun $" + obj.(Function).declaration.name.lexeme + ">"
	}
	return fmt.Sprint(obj)
}
