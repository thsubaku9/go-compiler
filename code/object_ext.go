package code

import (
	"fmt"
	"monkey-i/object"
)

const COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION_OBJ"

type CompiledFunction struct {
	Instructions Instructions
}

func (cf *CompiledFunction) Type() object.ObjectType { return COMPILED_FUNCTION_OBJ }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}
