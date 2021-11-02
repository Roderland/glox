package main

import "reflect"

type Parser struct {
	tokens  []Token
	current int
}

func _Parser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (parser *Parser) parse() []Stmt {
	statements := make([]Stmt, 0)
	for !parser.isEnd() {
		statements = append(statements, parser.declaration())
	}
	return statements
}

/*  ===================  Statement  ===================  */
func (parser *Parser) declaration() Stmt {
	if parser.match(FUN) {
		return parser.functionDeclaration()
	}
	if parser.match(VAR) {
		return parser.varDeclaration()
	}
	return parser.statement()
}

func (parser *Parser) functionDeclaration() Stmt {
	name := parser.consume(IDENTIFIER, "Expect function name.")

	parser.consume(LEFT_PAREN, "Expect '(' after function name.")
	params := make([]Token, 0)
	if parser.peek().tokenType != RIGHT_PAREN {
		for {
			params = append(params, parser.consume(IDENTIFIER, "Expect parameter name."))
			if !parser.match(COMMA) {
				break
			}
		}
	}
	parser.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	parser.consume(LEFT_BRACE, "Expect '{' before function body.")
	stmts := make([]Stmt, 0)
	for parser.peek().tokenType != RIGHT_BRACE {
		stmts = append(stmts, parser.declaration())
	}
	parser.consume(RIGHT_BRACE, "Expect '}' after block.")

	return functionStmt{name, params, stmts}
}

func (parser *Parser) varDeclaration() Stmt {
	name := parser.consume(IDENTIFIER, "Expect variable name.")
	var initializer Expr
	if parser.match(EQUAL) {
		initializer = parser.expression()
	}
	parser.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return varStmt{name, initializer}
}

func (parser *Parser) statement() Stmt {
	if parser.match(PRINT) {
		return parser.printStatement()
	}
	if parser.match(LEFT_BRACE) {
		return parser.blockStatement()
	}
	if parser.match(IF) {
		return parser.ifStatement()
	}
	if parser.match(WHILE) {
		return parser.whileStatement()
	}
	if parser.match(FOR) {
		return parser.forStatement()
	}
	if parser.match(RETURN) {
		return parser.returnStatement()
	}
	return parser.exprStatement()
}

func (parser *Parser) returnStatement() Stmt {
	keyword := parser.previous()
	var value Expr
	if parser.peek().tokenType != SEMICOLON {
		value = parser.expression()
	}
	parser.consume(SEMICOLON, "Expect ';' after return value.")
	return returnStmt{keyword, value}
}

func (parser *Parser) forStatement() Stmt {
	parser.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	if parser.match(SEMICOLON) {
		// no init
	} else if parser.match(VAR) {
		initializer = parser.varDeclaration()
	} else {
		initializer = parser.exprStatement()
	}

	var condition Expr
	if parser.peek().tokenType != SEMICOLON {
		condition = parser.expression()
	}
	parser.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if parser.peek().tokenType != RIGHT_PAREN {
		increment = parser.expression()
	}
	parser.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body := parser.statement()

	if increment != nil {
		body = blockStmt{[]Stmt{body, exprStmt{increment}}}
	}

	if condition == nil {
		condition = Literal{true}
	}

	var loop Stmt = whileStmt{condition, body}

	if initializer != nil {
		loop = blockStmt{[]Stmt{initializer, loop}}
	}

	return loop
}

func (parser *Parser) whileStatement() Stmt {
	parser.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := parser.expression()
	parser.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := parser.statement()
	return whileStmt{condition, body}
}

func (parser *Parser) ifStatement() Stmt {
	parser.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := parser.expression()
	parser.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch := parser.statement()
	var elseBranch Stmt
	if parser.match(ELSE) {
		elseBranch = parser.statement()
	}
	return ifStmt{condition, thenBranch, elseBranch}
}

func (parser *Parser) blockStatement() Stmt {
	stmts := make([]Stmt, 0)
	for parser.peek().tokenType != RIGHT_BRACE {
		stmts = append(stmts, parser.declaration())
	}
	parser.consume(RIGHT_BRACE, "Expect '}' after block.")
	return blockStmt{stmts}
}

func (parser *Parser) printStatement() Stmt {
	value := parser.expression()
	parser.consume(SEMICOLON, "Expect ';' after value.")
	return printStmt{expr: value}
}

func (parser *Parser) exprStatement() Stmt {
	expr := parser.expression()
	parser.consume(SEMICOLON, "Expect ';' after value.")
	return exprStmt{expr: expr}
}

/*  ===================  Expression  ===================  */
func (parser *Parser) expression() Expr {
	return parser.assignment()
}

// { "=" }
func (parser *Parser) assignment() Expr {
	left := parser.or()
	if parser.match(EQUAL) {
		equal := parser.previous()
		right := parser.assignment()
		if reflect.TypeOf(left) == reflect.TypeOf(Variable{}) {
			name := left.(Variable).name
			return Assign{
				name:  name,
				value: right,
			}
		}
		exitWithErr(equal.line, "Invalid assignment target.")
	}
	return left
}

// { "or" }
func (parser *Parser) or() Expr {
	left := parser.and()
	for parser.match(OR) {
		operator := parser.previous()
		right := parser.and()
		left = Logical{left, operator, right}
	}
	return left
}

// { "and" }
func (parser *Parser) and() Expr {
	left := parser.equality()
	for parser.match(AND) {
		operator := parser.previous()
		right := parser.equality()
		left = Logical{left, operator, right}
	}
	return left
}

// { "==", "!=" }
func (parser *Parser) equality() Expr {
	left := parser.comparison()
	for parser.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right := parser.comparison()
		left = Binary{left, operator, right}
	}
	return left
}

// { ">", ">=", "<", "<=" }
func (parser *Parser) comparison() Expr {
	left := parser.term()
	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right := parser.term()
		left = Binary{left, operator, right}
	}
	return left
}

// { "+", "-" }
func (parser *Parser) term() Expr {
	left := parser.factor()
	for parser.match(PLUS, MINUS) {
		operator := parser.previous()
		right := parser.factor()
		left = Binary{left, operator, right}
	}
	return left
}

// { "*", "/" }
func (parser *Parser) factor() Expr {
	left := parser.unary()
	for parser.match(SLASH, STAR) {
		operator := parser.previous()
		right := parser.unary()
		left = Binary{left, operator, right}
	}
	return left
}

// { "!", "-" }
func (parser *Parser) unary() Expr {
	if parser.match(BANG, MINUS) {
		operator := parser.previous()
		right := parser.unary()
		return Unary{operator, right}
	}
	return parser.call()
}

// { call-function }
func (parser *Parser) call() Expr {
	callee := parser.primary()
	for {
		if parser.match(LEFT_PAREN) {
			args := make([]Expr, 0)
			if parser.peek().tokenType != RIGHT_PAREN {
				for {
					args = append(args, parser.expression())
					if !parser.match(COMMA) {
						break
					}
				}
			}
			paren := parser.consume(RIGHT_PAREN, "Expect ')' after arguments.")
			callee = Call{callee, paren, args}
		} else {
			break
		}
	}
	return callee
}

// { "true", "false", "nil", Number, String, "(" }
func (parser *Parser) primary() Expr {
	if parser.match(TRUE) {
		return Literal{true}
	}
	if parser.match(FALSE) {
		return Literal{false}
	}
	if parser.match(NIL) {
		return Literal{nil}
	}
	if parser.match(NUMBER, STRING) {
		return Literal{parser.previous().literal}
	}
	if parser.match(LEFT_PAREN) {
		expr := parser.expression()
		parser.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}
	if parser.match(IDENTIFIER) {
		return Variable{parser.previous()}
	}
	exitWithErr(parser.peek().line, "Unexpected '"+parser.peek().lexeme+"' at here.")
	return nil
}

func (parser *Parser) peek() Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) isEnd() bool {
	return parser.peek().tokenType == EOF
}

func (parser *Parser) next() Token {
	token := parser.tokens[parser.current]
	parser.current++
	return token
}

func (parser *Parser) previous() Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) match(types ...uint8) bool {
	for _, u := range types {
		if parser.peek().tokenType == u {
			parser.next()
			return true
		}
	}
	return false
}

func (parser *Parser) consume(expected uint8, message string) Token {
	if parser.peek().tokenType != expected {
		exitWithErr(parser.peek().line, message)
	}
	return parser.next()
}
