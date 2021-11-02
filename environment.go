package main

type Envir struct {
	enclosing *Envir
	values    map[string]interface{}
}

//func copyMap(src map[string]interface{}) map[string]interface{} {
//	dst := map[string]interface{}{}
//	for k, v := range src {
//		dst[k] = v
//	}
//	return dst
//}

func (envir *Envir) define(name string, value interface{}) {
	//newValues := copyMap(envir.values)
	//newValues[name] = value
	//return &Envir{
	//	enclosing: envir.enclosing,
	//	values:    newValues,
	//}
	envir.values[name] = value
}

func (envir *Envir) get(name Token) interface{} {
	value, ok := envir.values[name.lexeme]
	if !ok {
		if envir.enclosing != nil {
			return envir.enclosing.get(name)
		}
		exitWithErr(name.line, "Undefined variable '"+name.lexeme+"'.")
	}
	return value
}

func (envir *Envir) assign(name Token, value interface{}) {
	_, ok := envir.values[name.lexeme]
	if !ok {
		if envir.enclosing != nil {
			envir.enclosing.assign(name, value)
			return
		}
		exitWithErr(name.line, "Undefined variable '"+name.lexeme+"'.")
	}
	envir.values[name.lexeme] = value
}
