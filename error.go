package main

import (
	"fmt"
	"os"
)

func exitWithErr(format string, a ...interface{}) {
	fmt.Println(fmt.Errorf(format, a...))
	os.Exit(1)
}
