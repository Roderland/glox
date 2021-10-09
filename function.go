package main

type callable interface {
	arity() int
	call(*interpreter, []interface{}) interface{}
}

type function struct {
	declaration *functionStmt
}

func (f *function) call(itp *interpreter, args []interface{}) interface{} {
	itp.env.push(*newScope(itp.env.peek()))
	for i := range f.declaration.params {
		itp.env.peek().define(f.declaration.params[i].text, args[i])
	}
	for _, s := range f.declaration.body {
		s.exec(itp)
	}
	itp.env.pop()
	itp.jp = J_NONE
	return itp.res
}

func (f *function) arity() int {
	return len(f.declaration.params)
}
