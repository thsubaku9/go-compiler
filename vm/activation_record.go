package vm

import "monkey-c/code"

type ActivationRecord struct {
	fn                 *code.CompiledFunction
	instructionPointer int
	basePointer        int
}

func (ar *ActivationRecord) Instructions() code.Instructions {
	return ar.fn.Instructions
}

func NewRecord(fn *code.CompiledFunction, basePointer int) *ActivationRecord {
	return &ActivationRecord{
		fn:                 fn,
		instructionPointer: -1,
		basePointer:        basePointer,
	}
}
