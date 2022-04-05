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

func exitWithErr(line int, message string) {
	out("[line %d] %s\n", line, message)
	os.Exit(1)
}

var Buf bytes.Buffer
var multiWriter io.Writer

// Play 编译成动态链接库作为插件开放给其他程序调用
func Play(code string) {
	Buf = bytes.Buffer{}
	// 向stdout和buffer中输出程序执行结果
	multiWriter = io.MultiWriter(os.Stdout, &Buf)
	tokens := _Lexer(code).lex()
	stmts := _Parser(tokens).parse()
	_Interpreter().interpret(stmts)
}

func out(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(multiWriter, format, a...)
}
