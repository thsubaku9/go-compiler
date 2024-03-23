package vm

import (
	"fmt"
	"monkey-c/compiler"
	"monkey-i/ast"
	"monkey-i/lexer"
	"monkey-i/object"
	"monkey-i/parser"
	"testing"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {

		prog := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(prog)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

	}
}
