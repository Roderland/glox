package main

import (
	"reflect"
)

type parser struct {
	tokens []token
	cur    int
}

func newParser(tokens []token) *parser {
	return &parser{tokens: tokens, cur: 0}
}

func (p *parser) run() []stmt {
	stmtList := make([]stmt, 0)
	for !p.eof() {
		stmtList = append(stmtList, p.declaration())
	}
	return stmtList
}

func (p *parser) declaration() stmt {
	if p.match(VAR) {
		return p.varStatement()
	}
	return p.statement()
}

func (p *parser) varStatement() stmt {
	name := p.consume(IDENTIFIER, "Expect a variable name")
	var initializer expr = nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration")
	return &varStmt{name, initializer}
}

func (p *parser) statement() stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LBRACE) {
		return p.blockStatement()
	}
	return p.expressionStatement()
}

func (p *parser) blockStatement() stmt {
	stmts := make([]stmt, 0)
	for !p.eof() && p.tokens[p.cur].ttype != RBRACE {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RBRACE, "Expect '}' after block")
	return &blockStmt{stmts: stmts}
}

func (p *parser) printStatement() stmt {
	value := p.expression()
	p.consume(SEMICOLON, "expect ';' after value.")
	return &printStmt{body: value}
}

func (p *parser) expressionStatement() stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return &exprStmt{body: value}
}

func (p *parser) expression() expr {
	return p.assignment()
}

func (p *parser) assignment() expr {
	left := p.equality()
	if p.match(EQUAL) {
		equals := p.tokens[p.cur-1]
		right := p.assignment()
		if reflect.TypeOf(left) == reflect.TypeOf(&variable{}) {
			name := left.(*variable).name
			return &assign{
				name:  name,
				right: right,
			}
		}
		exitWithErr("[ line %d ] Operator '=' expect a variable at left", equals.line)
	}
	return left
}

func (p *parser) equality() expr {
	e := p.comparison()
	for p.match(BEQUAL, EEQUAL) {
		operator := p.tokens[p.cur-1]
		right := p.comparison()
		e = &binary{
			left:     e,
			operator: operator,
			right:    right,
		}
	}
	return e
}

func (p *parser) comparison() expr {
	e := p.term()
	for p.match(LESS, LEQUAL, GREATER, GEQUAL) {
		operator := p.tokens[p.cur-1]
		right := p.term()
		e = &binary{
			left:     e,
			operator: operator,
			right:    right,
		}
	}
	return e
}

func (p *parser) term() expr {
	e := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.tokens[p.cur-1]
		right := p.factor()
		e = &binary{
			left:     e,
			operator: operator,
			right:    right,
		}
	}
	return e
}

func (p *parser) factor() expr {
	e := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.tokens[p.cur-1]
		right := p.unary()
		e = &binary{
			left:     e,
			operator: operator,
			right:    right,
		}
	}
	return e
}

func (p *parser) unary() expr {
	if p.match(BANG, MINUS) {
		operator := p.tokens[p.cur-1]
		right := p.unary()
		return &unary{
			operator: operator,
			right:    right,
		}
	}
	return p.primary()
}

func (p *parser) primary() expr {
	if p.match(FALSE) {
		return &literal{false}
	}
	if p.match(TRUE) {
		return &literal{true}
	}
	if p.match(NIL) {
		return &literal{nil}
	}
	if p.match(NUMBER, STRING) {
		return &literal{p.tokens[p.cur-1].literal}
	}
	if p.match(LPAREN) {
		e := p.expression()
		p.consume(RPAREN, "Expect ')' after expression")
		return &group{e}
	}
	if p.match(IDENTIFIER) {
		return &variable{p.tokens[p.cur-1]}
	}
	exitWithErr("[ line %d ] Expect expression", p.tokens[p.cur].line)
	return nil
}

func (p *parser) match(tts ...tokenType) bool {
	for _, tt := range tts {
		if !p.eof() && tt == p.tokens[p.cur].ttype {
			p.cur++
			return true
		}
	}
	return false
}

func (p *parser) eof() bool {
	return p.cur >= len(p.tokens)
}

func (p *parser) consume(tt tokenType, s string) token {
	if !p.match(tt) {
		exitWithErr("[ line %d ] %s", p.tokens[p.cur].line, s)
	}
	return p.tokens[p.cur-1]
}
