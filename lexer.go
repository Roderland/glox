package main

import "strconv"

type Lexer struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func _Lexer(source string) *Lexer {
	return &Lexer{source: source, tokens: []Token{}, start: 0, current: 0, line: 1}
}

func (lexer *Lexer) lex() []Token {
	for !lexer.isEnd() {
		lexer.start = lexer.current
		lexer.scanToken()
	}
	lexer.tokens = append(lexer.tokens, _Token(EOF, "$EOF", nil, lexer.line))
	return lexer.tokens
}

func (lexer *Lexer) scanToken() {
	char := lexer.next()
	switch char {
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		lexer.line++
	case '(':
		lexer.addToken(LEFT_PAREN, nil)
	case ')':
		lexer.addToken(RIGHT_PAREN, nil)
	case '{':
		lexer.addToken(LEFT_BRACE, nil)
	case '}':
		lexer.addToken(RIGHT_BRACE, nil)
	case ',':
		lexer.addToken(COMMA, nil)
	case ';':
		lexer.addToken(SEMICOLON, nil)
	case '.':
		lexer.addToken(DOT, nil)
	case '+':
		lexer.addToken(PLUS, nil)
	case '-':
		lexer.addToken(MINUS, nil)
	case '*':
		lexer.addToken(STAR, nil)
	case '/':
		if lexer.match('/') {
			for !lexer.isEnd() && lexer.peek() != '\n' {
				lexer.next()
			}
		} else {
			lexer.addToken(SLASH, nil)
		}
	case '!':
		if lexer.match('=') {
			lexer.addToken(BANG_EQUAL, nil)
		} else {
			lexer.addToken(BANG, nil)
		}
	case '=':
		if lexer.match('=') {
			lexer.addToken(EQUAL_EQUAL, nil)
		} else {
			lexer.addToken(EQUAL, nil)
		}
	case '>':
		if lexer.match('=') {
			lexer.addToken(GREATER_EQUAL, nil)
		} else {
			lexer.addToken(GREATER, nil)
		}
	case '<':
		if lexer.match('=') {
			lexer.addToken(LESS_EQUAL, nil)
		} else {
			lexer.addToken(LESS, nil)
		}
	case '"':
		for !lexer.isEnd() && lexer.peek() != '"' {
			if lexer.peek() == '\n' {
				lexer.line++
			}
			lexer.next()
		}
		if lexer.isEnd() {
			exitWithErr(lexer.line, "Unterminated string.")
		}
		lexer.next()
		str := lexer.source[lexer.start+1 : lexer.current-1]
		lexer.addToken(STRING, str)
	default:
		if isDigit(char) {
			for isDigit(lexer.peek()) {
				lexer.next()
			}
			if lexer.peek() == '.' && isDigit(lexer.peekNext()) {
				lexer.next()
				for isDigit(lexer.peek()) {
					lexer.next()
				}
			}
			double, err := strconv.ParseFloat(lexer.source[lexer.start:lexer.current], 64)
			if err != nil {
				exitWithErr(lexer.line, err.Error())
			}
			lexer.addToken(NUMBER, double)
		} else if isAlpha(char) {
			for isAlphaOrDigit(lexer.peek()) {
				lexer.next()
			}
			lexer.addToken(findType(lexer.source[lexer.start:lexer.current]), nil)
		} else {
			exitWithErr(lexer.line, "Unexpected character.")
		}
	}
}

func (lexer *Lexer) addToken(tokenType uint8, literal interface{}) {
	lexeme := lexer.source[lexer.start:lexer.current]
	lexer.tokens = append(lexer.tokens, _Token(tokenType, lexeme, literal, lexer.line))
}

func (lexer *Lexer) isEnd() bool {
	return lexer.current >= len(lexer.source)
}

func (lexer *Lexer) next() byte {
	char := lexer.source[lexer.current]
	lexer.current++
	return char
}

func (lexer *Lexer) match(expected byte) bool {
	if !lexer.isEnd() && lexer.source[lexer.current] == expected {
		lexer.current++
		return true
	}
	return false
}

func (lexer *Lexer) peek() byte {
	if lexer.isEnd() {
		return 0
	}
	return lexer.source[lexer.current]
}

func (lexer *Lexer) peekNext() byte {
	if lexer.current+1 >= len(lexer.source) {
		return 0
	}
	return lexer.source[lexer.current+1]
}

func isAlphaOrDigit(c byte) bool {
	return isDigit(c) || isAlpha(c)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}
