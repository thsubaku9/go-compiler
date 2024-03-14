package compiler

import (
	"monkey-c/code"
	"monkey-i/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}
