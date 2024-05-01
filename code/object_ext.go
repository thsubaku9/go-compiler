package code

import (
	"fmt"
	"monkey-i/object"
)

const (
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION_OBJ"
	CLOSURE_OBJ           = "CLOSURE"
)

type CompiledFunction struct {
	Instructions  Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() object.ObjectType { return COMPILED_FUNCTION_OBJ }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

type Closure struct {
	Fn   *CompiledFunction
	Free []object.Object
}

func (cl *Closure) Type() object.ObjectType { return CLOSURE_OBJ }
func (cl *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", cl)
}
