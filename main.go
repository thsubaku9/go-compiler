package main

import (
	"monkey-c/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
