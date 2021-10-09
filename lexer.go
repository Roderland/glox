package main

import (
	"strconv"
)

type lexer struct {
	source string
	tokens []token
	left   int
	right  int
	line   uint
}

func newLexer(source string) *lexer {
	return &lexer{
		source: source,
		left:   0,
		right:  0,
		line:   1,
	}
}

func (l *lexer) run() []token {
	for !l.eof() {
		l.left = l.right
		l.next()
	}
	return l.tokens
}

func (l *lexer) next() {
	u := l.source[l.right]
	l.right++
	switch u {
	case '+':
		l.add(PLUS, nil)
	case '-':
		l.add(MINUS, nil)
	case '*':
		l.add(STAR, nil)
	case '/':
		if l.match('/') {
			for !l.eof() && l.source[l.right] != '\n' {
				l.right++
			}
		} else {
			l.add(SLASH, nil)
		}
	case '(':
		l.add(LPAREN, nil)
	case ')':
		l.add(RPAREN, nil)
	case '{':
		l.add(LBRACE, nil)
	case '}':
		l.add(RBRACE, nil)
	case ',':
		l.add(COMMA, nil)
	case '.':
		l.add(DOT, nil)
	case ';':
		l.add(SEMICOLON, nil)
	case '!':
		if l.match('=') {
			l.add(BEQUAL, nil)
		} else {
			l.add(BANG, nil)
		}
	case '=':
		if l.match('=') {
			l.add(EEQUAL, nil)
		} else {
			l.add(EQUAL, nil)
		}
	case '<':
		if l.match('=') {
			l.add(LEQUAL, nil)
		} else {
			l.add(LESS, nil)
		}
	case '>':
		if l.match('=') {
			l.add(GEQUAL, nil)
		} else {
			l.add(GREATER, nil)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		l.line++
	case '"':
		l.string()
	default:
		if isDigit(u) {
			l.number()
		} else if isAlpha(u) {
			l.identifier()
		} else {
			exitWithErr("[ line %d ] Illegal token at '%s'", l.line, l.source[l.left:l.right])
		}
	}
}

func isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (l *lexer) identifier() {
	for !l.eof() && (isDigit(l.source[l.right]) || isAlpha(l.source[l.right])) {
		l.right++
	}
	s := l.source[l.left:l.right]
	if tt, ok := keywords[s]; ok {
		l.add(tt, nil)
		if tt == BREAK || tt == CONTINUE {
			if !l.match(';') {
				exitWithErr("[ line %d ] Expect ';' at end", l.line)
			}
		}
	} else {
		l.add(IDENTIFIER, nil)
	}
}

func (l *lexer) number() {
	for !l.eof() && isDigit(l.source[l.right]) {
		l.right++
	}
	if !l.eof() && l.source[l.right] == '.' {
		l.right++
		for !l.eof() && isDigit(l.source[l.right]) {
			l.right++
		}
	}

	s, err := strconv.ParseFloat(l.source[l.left:l.right], 64)
	if err != nil {
		exitWithErr("[ line %d ] Illegal number at '%s'", l.line, l.source[l.left:l.right])
	}
	l.add(NUMBER, s)
}

func (l *lexer) string() {
	for !l.eof() && l.source[l.right] != '"' {
		if l.source[l.right] == '\n' {
			l.line++
		}
		l.right++
	}
	if l.eof() {
		exitWithErr("[ line %d ] Expect '\"' after '%s'", l.line, l.source[l.left+1:l.right-1])
		return
	}
	l.right++
	s := l.source[l.left+1 : l.right-1]
	l.add(STRING, s)
}

func (l *lexer) match(c uint8) bool {
	if l.eof() || l.source[l.right] != c {
		return false
	}
	l.right++
	return true
}

func (l *lexer) add(ttype tokenType, literal interface{}) {
	tk := token{
		ttype:   ttype,
		text:    l.source[l.left:l.right],
		literal: literal,
		line:    l.line,
	}
	l.tokens = append(l.tokens, tk)
}

func (l *lexer) eof() bool {
	return l.right >= len(l.source)
}
