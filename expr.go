package main

type (
	expr interface {
		eval(scopes scopeList) interface{}
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
)
