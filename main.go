package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: glox [InputFile]")
		os.Exit(64)
	}
	bts, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read file '%s'.\n", os.Args[1])
		os.Exit(65)
	}

	Play(string(bts))
}

var Buf bytes.Buffer
var writer io.Writer

func init() {
	Buf = bytes.Buffer{}
	writer = io.MultiWriter(os.Stdout, &Buf)
}

func out(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(writer, format, a...)
}

// Play 可以编译成动态链接库作为插件开放给其他程序调用
func Play(code string) {
	Buf.Reset()
	tokens := _Lexer(code).lex()
	stmts := _Parser(tokens).parse()
	_Interpreter().interpret(stmts)
}
