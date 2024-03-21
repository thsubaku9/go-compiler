package vm

import (
	"monkey-i/ast"
	"monkey-i/lexer"
	"monkey-i/parser"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
