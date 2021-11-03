package main

type Function struct {
	declaration functionStmt
}

func (f *Function) call(interpreter *Interpreter, args []interface{}) interface{} {
	functionLocal := &Table{
		father: interpreter.local,
		values: map[string]interface{}{},
	}
	for i := 0; i < f.arity(); i++ {
		functionLocal.define(f.declaration.params[i].lexeme, args[i])
	}
	interpreter.enterScope(functionLocal)
	defer interpreter.enterScope(functionLocal.father)
	for _, stmt := range f.declaration.stmts {
		depth := len(interpreter.returnStack)
		stmt.exec(interpreter)
		if depth < len(interpreter.returnStack) {
			result := interpreter.returnStack[len(interpreter.returnStack)-1]
			interpreter.returnStack = interpreter.returnStack[:len(interpreter.returnStack)-1]
			return result
		}
	}
	return nil
}

func (f *Function) arity() int {
	return len(f.declaration.params)
}
