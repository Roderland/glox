package main

type Table struct {
	father *Table
	values map[string]interface{}
}

func (table *Table) define(name string, value interface{}) {
	table.values[name] = value
}

func (table *Table) get(name Token) interface{} {
	value, ok := table.values[name.lexeme]
	if !ok {
		if table.father != nil {
			return table.father.get(name)
		}
		exitWithErr(name.line, "Undefined variable '"+name.lexeme+"'.")
	}
	return value
}

func (table *Table) assign(name Token, value interface{}) {
	_, ok := table.values[name.lexeme]
	if !ok {
		if table.father != nil {
			table.father.assign(name, value)
			return
		}
		exitWithErr(name.line, "Undefined variable '"+name.lexeme+"'.")
	}
	table.values[name.lexeme] = value
}
