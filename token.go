package main

const (
	// Single-character tokens.
	LEFT_PAREN uint8 = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	TRUE
	VAR
	WHILE

	// End of file.
	EOF
)

type Token struct {
	tokenType uint8
	lexeme    string
	literal   interface{}
	line      int
}

func _Token(tokenType uint8, lexeme string, literal interface{}, line int) Token {
	return Token{tokenType: tokenType, lexeme: lexeme, literal: literal, line: line}
}

func findType(text string) uint8 {
	switch text {
	case "and":
		return AND
	case "else":
		return ELSE
	case "false":
		return FALSE
	case "for":
		return FOR
	case "fun":
		return FUN
	case "if":
		return IF
	case "nil":
		return NIL
	case "or":
		return OR
	case "print":
		return PRINT
	case "return":
		return RETURN
	case "true":
		return TRUE
	case "var":
		return VAR
	case "while":
		return WHILE
	default:
		return IDENTIFIER
	}
}
