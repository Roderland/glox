package main

import (
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(64)
	}
	run(os.Args[1])
}

func run(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		exitWithErr(err.Error())
	}
	tokens := newLexer(string(bytes)).run()
	stmts := newParser(tokens).run()
	newInterpreter(stmts).run()
}
