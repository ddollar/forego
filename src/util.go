package main

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(stdout, a...)
}
