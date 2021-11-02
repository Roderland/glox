package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: glox [inputFile]")
		os.Exit(64)
	}
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read file '%s'.\n", os.Args[1])
		os.Exit(1)
	}
	tokens := _Lexer(string(bytes)).lex()
	statements := _Parser(tokens).parse()
	_Interpreter().interpret(statements)
}

func exitWithErr(line int, message string) {
	fmt.Printf("[line %d] %s\n", line, message)
	os.Exit(1)
}
