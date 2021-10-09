package main

type (
	scope struct {
		pre *scope
		tab map[string]interface{}
	}
	scopeList []scope
)

func newScope(pre *scope) *scope {
	return &scope{pre: pre, tab: make(map[string]interface{})}
}

func (sl *scopeList) push(node scope) {
	*sl = append(*sl, node)
}

func (sl *scopeList) pop() {
	*sl = (*sl)[:len(*sl)-1]
}

func (sl *scopeList) peek() *scope {
	return &(*sl)[len(*sl)-1]
}

func (s *scope) define(key string, value interface{}) {
	if _, ok := s.tab[key]; ok {
		exitWithErr("Variable '%s' has been defined", key)
	}
	s.tab[key] = value
}

func (s *scope) assign(key string, value interface{}) {
	if _, ok := s.tab[key]; ok {
		s.tab[key] = value
		return
	}
	if s.pre != nil {
		s.pre.assign(key, value)
		return
	}
	exitWithErr("Variable '%s' is not defined", key)
}

func (s *scope) get(key string) interface{} {
	if value, ok := s.tab[key]; ok {
		return value
	}
	if s.pre != nil {
		return s.pre.get(key)
	}
	exitWithErr("Variable '%s' is not defined", key)
	return nil
}
