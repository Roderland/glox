package main

const (
	PLUS      tokenType = iota // "+"
	MINUS                      // "-"
	STAR                       // "*"
	SLASH                      // "/"
	LPAREN                     // "("
	RPAREN                     // ")"
	LBRACE                     // "{"
	RBRACE                     // "}"
	COMMA                      // ","
	DOT                        // "."
	SEMICOLON                  // ";"
	BANG                       // "!"
	EQUAL                      // "="
	BEQUAL                     // "!="
	EEQUAL                     // "=="
	LESS                       // "<"
	LEQUAL                     // "<="
	GREATER                    // ">"
	GEQUAL                     // ">="
	AND                        // "and"
	OR                         // "or"
	IF                         // "if"
	ELSE                       // "else"
	WHILE                      // "while"
	FOR                        // "for"
	VAR                        // "var"
	NIL                        // "nil"
	TRUE                       // "true"
	FALSE                      // "false"
	FUN                        // "fun"
	RETURN                     // "return"
	PRINT                      // "print"
	STRING
	NUMBER
	IDENTIFIER
)

var keywords = make(map[string]tokenType)

func init() {
	keywords["and"] = AND
	keywords["or"] = OR
	keywords["if"] = IF
	keywords["else"] = ELSE
	keywords["while"] = WHILE
	keywords["for"] = FOR
	keywords["var"] = VAR
	keywords["nil"] = NIL
	keywords["true"] = TRUE
	keywords["false"] = FALSE
	keywords["fun"] = FUN
	keywords["return"] = RETURN
	keywords["print"] = PRINT
}

type (
	tokenType uint8
	token     struct {
		ttype   tokenType
		text    string
		literal interface{}
		line    uint
	}
)
