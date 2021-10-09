package main

type (
	expr interface {
		eval(*interpreter) interface{}
	}
	binary struct {
		left     expr
		operator token
		right    expr
	}
	unary struct {
		operator token
		right    expr
	}
	literal struct {
		value interface{}
	}
	group struct {
		body expr
	}
	variable struct {
		name token
	}
	assign struct {
		name  token
		right expr
	}
	logical struct {
		left     expr
		operator token
		right    expr
	}
	call struct {
		callee    expr
		paren     token
		arguments []expr
	}
)
