package main

import (
	"testing"
)

func TestLexer(t *testing.T) {
	s := "print \"Hello, world!\";"
	tokens := _Lexer(s).lex()
	var expect = []string{"print", "\"Hello, world!\"", ";", "$EOF"}
	if len(expect) != len(tokens) {
		t.Errorf("Expected %d tokens but get %d.\n", len(expect), len(tokens))
	}
	for idx, token := range tokens {
		if expect[idx] != token.lexeme {
			t.Errorf("Expected token: %s but get %s.\n", expect[idx], token.lexeme)
		}
	}
}

func TestParser(t *testing.T) {
	tokens := make([]Token, 0)
	tokens = append(tokens, _Token(IF, "if", nil, 1))
	tokens = append(tokens, _Token(LEFT_PAREN, "(", nil, 1))
	tokens = append(tokens, _Token(TRUE, "true", true, 1))
	tokens = append(tokens, _Token(RIGHT_PAREN, ")", nil, 1))
	tokens = append(tokens, _Token(PRINT, "print", nil, 2))
	tokens = append(tokens, _Token(STRING, "hello", "hello", 2))
	tokens = append(tokens, _Token(SEMICOLON, ";", nil, 2))
	tokens = append(tokens, _Token(EOF, "$EOF", nil, 2))
	stmts := _Parser(tokens).parse()
	if len(stmts) != 1 {
		t.Errorf("Expected 1 statement.\n")
	}
	stmt := stmts[0].(ifStmt)
	if stmt.condition.(Literal).value != true {
		t.Errorf("Expected condition equals true.\n")
	}
	if stmt.thenBranch.(printStmt).expr.(Literal).value != "hello" {
		t.Errorf("Expected the expr of thenBranch equals \"hello\".\n")
	}
	if stmt.elseBranch != nil {
		t.Errorf("Expected elseBranch is empty.\n")
	}
}

func TestInterpreter(t *testing.T) {
	stmt := ifStmt{
		condition:  Literal{value: true},
		thenBranch: printStmt{expr: Literal{value: "hello"}},
		elseBranch: nil,
	}
	stmt.exec(_Interpreter())
	if Buf.String() != "hello\n" {
		t.Errorf("Expected \"hello\" in buffer.\n")
	}
}
