package main

import (
	"reflect"
)

type parser struct {
	tokens []token
	cur    int
	inLoop bool
}

func newParser(tokens []token) *parser {
	return &parser{tokens: tokens, cur: 0, inLoop: false}
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
	if p.match(FUN) {
		return p.functionStatement()
	}
	return p.statement()
}

func (p *parser) functionStatement() stmt {
	name := p.consume(IDENTIFIER, "Expect a function name")
	p.consume(LPAREN, "Expect '(' after function name")
	params := make([]token, 0)
	if !p.eof() && p.tokens[p.cur].ttype != RPAREN {
		params = append(params, p.consume(IDENTIFIER, "Expect parameter name"))
		for p.match(COMMA) {
			if len(params) >= 255 {
				exitWithErr("Can't have more than 255 parameters")
			}
			params = append(params, p.consume(IDENTIFIER, "Expect parameter name"))
		}
	}
	p.consume(RPAREN, "Expect ')' after parameters")
	p.consume(LBRACE, "Expect '{' before function body")
	body := p.blockStatement()
	return &functionStmt{
		name:   name,
		params: params,
		body:   body,
	}
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
		return &blockStmt{p.blockStatement()}
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(BREAK, CONTINUE) {
		token := p.tokens[p.cur-1]
		if p.inLoop {
			if token.ttype == BREAK {
				return &breakStmt{}
			} else {
				return &continueStmt{}
			}
		} else {
			exitWithErr("[ line %d ] '%s' is outside loop", token.line, token.text)
		}
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	return p.expressionStatement()
}

func (p *parser) returnStatement() stmt {
	keyword := p.tokens[p.cur-1]
	var value expr
	if !p.eof() && p.tokens[p.cur].ttype != SEMICOLON {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value")
	return &returnStmt{
		keyword: keyword,
		value:   value,
	}
}

func (p *parser) forStatement() stmt {
	p.consume(LPAREN, "Expect '(' after 'for'")
	var initializer stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varStatement()
	} else {
		initializer = p.expressionStatement()
	}
	var condition expr = &literal{true}
	if !p.eof() && p.tokens[p.cur].ttype != SEMICOLON {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition")
	var increment expr
	if !p.eof() && p.tokens[p.cur].ttype != RPAREN {
		increment = p.expression()
	}
	p.consume(RPAREN, "Expect ')' after for clauses")
	p.inLoop = true
	body := p.statement()
	p.inLoop = false
	body = &whileStmt{
		condition: condition,
		body:      body,
		increment: increment,
	}
	if initializer != nil {
		stmts := make([]stmt, 0)
		stmts = append(stmts, initializer, body)
		body = &blockStmt{stmts}
	}
	return body
}

func (p *parser) whileStatement() stmt {
	p.consume(LPAREN, "Expect '(' after while")
	condition := p.expression()
	p.consume(RPAREN, "Expect ')' after condition")
	p.inLoop = true
	body := p.statement()
	p.inLoop = false
	return &whileStmt{
		condition: condition,
		body:      body,
		increment: nil,
	}
}

func (p *parser) ifStatement() stmt {
	p.consume(LPAREN, "Expect '(' after 'if'")
	condition := p.expression()
	p.consume(RPAREN, "Expect ')' after condition")
	thenBranch := p.statement()
	var elseBranch stmt = nil
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return &ifStmt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (p *parser) blockStatement() []stmt {
	stmts := make([]stmt, 0)
	for !p.eof() && p.tokens[p.cur].ttype != RBRACE {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RBRACE, "Expect '}' after block")
	return stmts
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
	left := p.or()
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

func (p *parser) or() expr {
	left := p.and()
	for p.match(OR) {
		operator := p.tokens[p.cur-1]
		right := p.and()
		left = &logical{
			left:     left,
			operator: operator,
			right:    right,
		}
	}
	return left
}

func (p *parser) and() expr {
	left := p.equality()
	for p.match(AND) {
		operator := p.tokens[p.cur-1]
		right := p.equality()
		left = &logical{
			left:     left,
			operator: operator,
			right:    right,
		}
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
	return p.call()
}

func (p *parser) call() expr {
	ep := p.primary()
	for {
		if p.match(LPAREN) {
			ep = p.finishCall(ep)
		} else {
			break
		}
	}
	return ep
}

func (p *parser) finishCall(callee expr) expr {
	arguments := make([]expr, 0)
	if !p.eof() && p.tokens[p.cur].ttype != RPAREN {
		arguments = append(arguments, p.expression())
		for p.match(COMMA) {
			if len(arguments) >= 255 {
				exitWithErr("Can't have more than 255 arguments")
			}
			arguments = append(arguments, p.expression())
		}
	}
	paren := p.consume(RPAREN, "Expect ')' after arguments")
	return &call{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
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
