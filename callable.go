package main

type (
	Callable interface {
		call(interpreter *Interpreter, args []interface{}) interface{}
		arity() int
	}
	Function struct {
		declaration functionStmt
		closure     *Envir
	}
)

func (f *Function) call(interpreter *Interpreter, args []interface{}) interface{} {
	functionEnvironment := &Envir{
		enclosing: f.closure,
		values:    map[string]interface{}{},
	}
	for i := 0; i < f.arity(); i++ {
		functionEnvironment.define(f.declaration.params[i].lexeme, args[i])
	}
	father := interpreter.envir
	interpreter.enterScope(functionEnvironment)
	defer interpreter.enterScope(father)
	for _, stmt := range f.declaration.stmts {
		oldNum := len(interpreter.retVal)
		stmt.exec(interpreter)
		if oldNum < len(interpreter.retVal) {
			result := interpreter.retVal[len(interpreter.retVal)-1]
			interpreter.retVal = interpreter.retVal[:len(interpreter.retVal)-1]
			return result
		}
	}
	return nil
}

func (f *Function) arity() int {
	return len(f.declaration.params)
}
