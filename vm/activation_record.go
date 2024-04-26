package vm

import "monkey-c/code"

type ActivationRecord struct {
	fn                 *code.CompiledFunction
	instructionPointer int
}

func (ar *ActivationRecord) Instructions() code.Instructions {
	return ar.fn.Instructions
}

func NewRecord(fn *code.CompiledFunction) *ActivationRecord {
	return &ActivationRecord{
		fn:                 fn,
		instructionPointer: -1,
	}
}
